package system

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type SysExportTemplateService struct {
}

var SysExportTemplateServiceApp = new(SysExportTemplateService)

// CreateSysExportTemplate creates an export template record
// Author [piexlmax](https://github.com/piexlmax)
func (sysExportTemplateService *SysExportTemplateService) CreateSysExportTemplate(sysExportTemplate *system.SysExportTemplate) (err error) {
	err = global.GVA_DB.Create(sysExportTemplate).Error
	return err
}

// DeleteSysExportTemplate deletes an export template record
// Author [piexlmax](https://github.com/piexlmax)
func (sysExportTemplateService *SysExportTemplateService) DeleteSysExportTemplate(sysExportTemplate system.SysExportTemplate) (err error) {
	err = global.GVA_DB.Delete(&sysExportTemplate).Error
	return err
}

// DeleteSysExportTemplateByIds batch deletes export template records
// Author [piexlmax](https://github.com/piexlmax)
func (sysExportTemplateService *SysExportTemplateService) DeleteSysExportTemplateByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysExportTemplate{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateSysExportTemplate updates an export template record
// Author [piexlmax](https://github.com/piexlmax)
func (sysExportTemplateService *SysExportTemplateService) UpdateSysExportTemplate(sysExportTemplate system.SysExportTemplate) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		conditions := sysExportTemplate.Conditions
		e := tx.Delete(&[]system.Condition{}, "template_id = ?", sysExportTemplate.TemplateID).Error
		if e != nil {
			return e
		}
		sysExportTemplate.Conditions = nil

		joins := sysExportTemplate.JoinTemplate
		e = tx.Delete(&[]system.JoinTemplate{}, "template_id = ?", sysExportTemplate.TemplateID).Error
		if e != nil {
			return e
		}
		sysExportTemplate.JoinTemplate = nil

		e = tx.Updates(&sysExportTemplate).Error
		if e != nil {
			return e
		}
		if len(conditions) > 0 {
			for i := range conditions {
				conditions[i].ID = 0
			}
			e = tx.Create(&conditions).Error
		}
		if len(joins) > 0 {
			for i := range joins {
				joins[i].ID = 0
			}
			e = tx.Create(&joins).Error
		}
		return e
	})
}

// GetSysExportTemplate gets an export template record by id
// Author [piexlmax](https://github.com/piexlmax)
func (sysExportTemplateService *SysExportTemplateService) GetSysExportTemplate(id uint) (sysExportTemplate system.SysExportTemplate, err error) {
	err = global.GVA_DB.Where("id = ?", id).Preload("JoinTemplate").Preload("Conditions").First(&sysExportTemplate).Error
	return
}

// GetSysExportTemplateInfoList gets export template records with pagination
// Author [piexlmax](https://github.com/piexlmax)
func (sysExportTemplateService *SysExportTemplateService) GetSysExportTemplateInfoList(info systemReq.SysExportTemplateSearch) (list []system.SysExportTemplate, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// create db
	db := global.GVA_DB.Model(&system.SysExportTemplate{})
	var sysExportTemplates []system.SysExportTemplate
	// if there are search conditions, the search statement will be automatically created below
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.TableName != "" {
		db = db.Where("table_name = ?", info.TableName)
	}
	if info.TemplateID != "" {
		db = db.Where("template_id = ?", info.TemplateID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&sysExportTemplates).Error
	return sysExportTemplates, total, err
}

// ExportExcel exports to Excel
// Author [piexlmax](https://github.com/piexlmax)
func (sysExportTemplateService *SysExportTemplateService) ExportExcel(templateID string, values url.Values) (file *bytes.Buffer, name string, err error) {
	var params = values.Get("params")
	paramsValues, err := url.ParseQuery(params)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse params parameter: %v", err)
	}
	var template system.SysExportTemplate
	err = global.GVA_DB.Preload("Conditions").Preload("JoinTemplate").First(&template, "template_id = ?", templateID).Error
	if err != nil {
		return nil, "", err
	}
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	var templateInfoMap = make(map[string]string)
	columns, err := utils.GetJSONKeys(template.TemplateInfo)
	if err != nil {
		return nil, "", err
	}
	err = json.Unmarshal([]byte(template.TemplateInfo), &templateInfoMap)
	if err != nil {
		return nil, "", err
	}
	var tableTitle []string
	var selectKeyFmt []string
	for _, key := range columns {
		selectKeyFmt = append(selectKeyFmt, key)
		tableTitle = append(tableTitle, templateInfoMap[key])
	}

	selects := strings.Join(selectKeyFmt, ", ")
	var tableMap []map[string]interface{}
	db := global.GVA_DB
	if template.DBName != "" {
		db = global.MustGetGlobalDBByDBName(template.DBName)
	}

	// if there is custom SQL, use it first
	if template.SQL != "" {
		// convert url.Values to map[string]interface{} to support GORM named parameters
		sqlParams := make(map[string]interface{})
		for k, v := range paramsValues {
			if len(v) > 0 {
				sqlParams[k] = v[0]
			}
		}

		// execute raw SQL with @key named parameter support
		err = db.Raw(template.SQL, sqlParams).Scan(&tableMap).Error
		if err != nil {
			return nil, "", err
		}
	} else {
		if len(template.JoinTemplate) > 0 {
			for _, join := range template.JoinTemplate {
				db = db.Joins(join.JOINS + " " + join.Table + " ON " + join.ON)
			}
		}

		db = db.Select(selects).Table(template.TableName)

		filterDeleted := false

		filterParam := paramsValues.Get("filterDeleted")
		if filterParam == "true" {
			filterDeleted = true
		}

		if filterDeleted {
			// automatically filter soft-deleted records from the main table
			db = db.Where(fmt.Sprintf("%s.deleted_at IS NULL", template.TableName))

			// filter soft-deleted records from joined tables (if any)
			if len(template.JoinTemplate) > 0 {
				for _, join := range template.JoinTemplate {
					// check if the joined table has a deleted_at column
					hasDeletedAt := sysExportTemplateService.hasDeletedAtColumn(join.Table)
					if hasDeletedAt {
						db = db.Where(fmt.Sprintf("%s.deleted_at IS NULL", join.Table))
					}
				}
			}
		}

		if len(template.Conditions) > 0 {
			for _, condition := range template.Conditions {
				sql := fmt.Sprintf("%s %s ?", condition.Column, condition.Operator)
				value := paramsValues.Get(condition.From)

				if condition.Operator == "IN" || condition.Operator == "NOT IN" {
					sql = fmt.Sprintf("%s %s (?)", condition.Column, condition.Operator)
				}

				if condition.Operator == "BETWEEN" {
					sql = fmt.Sprintf("%s BETWEEN ? AND ?", condition.Column)
					startValue := paramsValues.Get("start" + condition.From)
					endValue := paramsValues.Get("end" + condition.From)
					if startValue != "" && endValue != "" {
						db = db.Where(sql, startValue, endValue)
					}
					continue
				}

				if value != "" {
					if condition.Operator == "LIKE" {
						value = "%" + value + "%"
					}
					db = db.Where(sql, value)
				}
			}
		}
		// pass limit via parameters
		limit := paramsValues.Get("limit")
		if limit != "" {
			l, e := strconv.Atoi(limit)
			if e == nil {
				db = db.Limit(l)
			}
		}
		// template default limit
		if limit == "" && template.Limit != nil && *template.Limit != 0 {
			db = db.Limit(*template.Limit)
		}

		// pass offset via parameters
		offset := paramsValues.Get("offset")
		if offset != "" {
			o, e := strconv.Atoi(offset)
			if e == nil {
				db = db.Offset(o)
			}
		}

		// get all columns of the current table
		table := template.TableName
		orderColumns, err := db.Migrator().ColumnTypes(table)
		if err != nil {
			return nil, "", err
		}

		// create a map to store column names
		fields := make(map[string]bool)

		for _, column := range orderColumns {
			fields[column.Name()] = true
		}

		// pass order via parameters
		order := paramsValues.Get("order")

		if order == "" && template.Order != "" {
			// if no order parameter, use the template's default order
			order = template.Order
		}

		if order != "" {
			checkOrderArr := strings.Split(order, " ")
			orderStr := ""
			// check if the requested order field is in the column list
			if _, ok := fields[checkOrderArr[0]]; !ok {
				return nil, "", fmt.Errorf("order by %s is not in the fields", order)
			}
			orderStr = checkOrderArr[0]
			if len(checkOrderArr) > 1 {
				if checkOrderArr[1] != "asc" && checkOrderArr[1] != "desc" {
					return nil, "", fmt.Errorf("order by %s is not secure", order)
				}
				orderStr = orderStr + " " + checkOrderArr[1]
			}
			db = db.Order(orderStr)
		}

		err = db.Debug().Find(&tableMap).Error
		if err != nil {
			return nil, "", err
		}
	}

	var rows [][]string
	rows = append(rows, tableTitle)
	for _, exTable := range tableMap {
		var row []string
		for _, column := range columns {
			column = strings.ReplaceAll(column, "\"", "")
			column = strings.ReplaceAll(column, "`", "")
			if len(template.JoinTemplate) > 0 {
				columnAs := strings.Split(column, " as ")
				if len(columnAs) > 1 {
					column = strings.TrimSpace(strings.Split(column, " as ")[1])
				} else {
					columnArr := strings.Split(column, ".")
					if len(columnArr) > 1 {
						column = strings.Split(column, ".")[1]
					}
				}
			}
			// special handling for time types
			if t, ok := exTable[column].(time.Time); ok {
				row = append(row, t.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, fmt.Sprintf("%v", exTable[column]))
			}
		}
		rows = append(rows, row)
	}
	for i, row := range rows {
		for j, colCell := range row {
			cell := fmt.Sprintf("%s%d", getColumnName(j+1), i+1)

			var sErr error
			if v, err := strconv.ParseFloat(colCell, 64); err == nil {
				sErr = f.SetCellValue("Sheet1", cell, v)
			} else if v, err := strconv.ParseInt(colCell, 10, 64); err == nil {
				sErr = f.SetCellValue("Sheet1", cell, v)
			} else {
				sErr = f.SetCellValue("Sheet1", cell, colCell)
			}

			if sErr != nil {
				return nil, "", sErr
			}
		}
	}
	f.SetActiveSheet(index)
	file, err = f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	return file, template.Name, nil
}

// PreviewSQL previews the final generated SQL (does not execute the query, only returns the SQL string)
// Author [piexlmax](https://github.com/piexlmax) & [trae-ai]
func (sysExportTemplateService *SysExportTemplateService) PreviewSQL(templateID string, values url.Values) (sqlPreview string, err error) {
	// parse params (consistent with export logic)
	var params = values.Get("params")
	paramsValues, _ := url.ParseQuery(params)

	// load template
	var template system.SysExportTemplate
	err = global.GVA_DB.Preload("Conditions").Preload("JoinTemplate").First(&template, "template_id = ?", templateID).Error
	if err != nil {
		return "", err
	}

	// parse template columns
	var templateInfoMap = make(map[string]string)
	columns, err := utils.GetJSONKeys(template.TemplateInfo)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(template.TemplateInfo), &templateInfoMap)
	if err != nil {
		return "", err
	}
	var selectKeyFmt []string
	for _, key := range columns {
		selectKeyFmt = append(selectKeyFmt, key)
	}
	selects := strings.Join(selectKeyFmt, ", ")

	// generate FROM and JOIN clauses
	var sb strings.Builder
	sb.WriteString("SELECT ")
	sb.WriteString(selects)
	sb.WriteString(" FROM ")
	sb.WriteString(template.TableName)

	if len(template.JoinTemplate) > 0 {
		for _, join := range template.JoinTemplate {
			sb.WriteString(" ")
			sb.WriteString(join.JOINS)
			sb.WriteString(" ")
			sb.WriteString(join.Table)
			sb.WriteString(" ON ")
			sb.WriteString(join.ON)
		}
	}

	// WHERE conditions
	var wheres []string

	// soft delete filter
	filterDeleted := false
	if paramsValues != nil {
		filterParam := paramsValues.Get("filterDeleted")
		if filterParam == "true" {
			filterDeleted = true
		}
	}
	if filterDeleted {
		wheres = append(wheres, fmt.Sprintf("%s.deleted_at IS NULL", template.TableName))
		if len(template.JoinTemplate) > 0 {
			for _, join := range template.JoinTemplate {
				if sysExportTemplateService.hasDeletedAtColumn(join.Table) {
					wheres = append(wheres, fmt.Sprintf("%s.deleted_at IS NULL", join.Table))
				}
			}
		}
	}

	// template conditions (keeping parsing rules synchronized with ExportExcel)
	if len(template.Conditions) > 0 {
		for _, condition := range template.Conditions {
			op := strings.ToUpper(strings.TrimSpace(condition.Operator))
			col := strings.TrimSpace(condition.Column)

			// preview shows passed-in values first, otherwise show placeholders
			val := ""
			if paramsValues != nil {
				val = paramsValues.Get(condition.From)
			}

			switch op {
			case "BETWEEN":
				startValue := ""
				endValue := ""
				if paramsValues != nil {
					startValue = paramsValues.Get("start" + condition.From)
					endValue = paramsValues.Get("end" + condition.From)
				}
				if startValue != "" && endValue != "" {
					wheres = append(wheres, fmt.Sprintf("%s BETWEEN '%s' AND '%s'", col, startValue, endValue))
				} else {
					wheres = append(wheres, fmt.Sprintf("%s BETWEEN {start%s} AND {end%s}", col, condition.From, condition.From))
				}
			case "IN", "NOT IN":
				if val != "" {
					// simple display of comma-separated values
					parts := strings.Split(val, ",")
					for i := range parts {
						parts[i] = strings.TrimSpace(parts[i])
					}
					wheres = append(wheres, fmt.Sprintf("%s %s ('%s')", col, op, strings.Join(parts, "','")))
				} else {
					wheres = append(wheres, fmt.Sprintf("%s %s ({%s})", col, op, condition.From))
				}
			case "LIKE":
				if val != "" {
					wheres = append(wheres, fmt.Sprintf("%s LIKE '%%%s%%'", col, val))
				} else {
					wheres = append(wheres, fmt.Sprintf("%s LIKE {%%%s%%}", col, condition.From))
				}
			default:
				if val != "" {
					wheres = append(wheres, fmt.Sprintf("%s %s '%s'", col, op, val))
				} else {
					wheres = append(wheres, fmt.Sprintf("%s %s {%s}", col, op, condition.From))
				}
			}
		}
	}

	if len(wheres) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(wheres, " AND "))
	}

	// order
	order := ""
	if paramsValues != nil {
		order = paramsValues.Get("order")
	}
	if order == "" && template.Order != "" {
		order = template.Order
	}
	if order != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(order)
	}

	// limit/offset (do not generate if passed-in or default value is 0)
	limitStr := ""
	offsetStr := ""
	if paramsValues != nil {
		limitStr = paramsValues.Get("limit")
		offsetStr = paramsValues.Get("offset")
	}

	// handle template default limit (only when non-zero)
	if limitStr == "" && template.Limit != nil && *template.Limit != 0 {
		limitStr = strconv.Itoa(*template.Limit)
	}

	// parse as numbers to determine whether to generate
	limitInt := 0
	offsetInt := 0
	if limitStr != "" {
		if v, e := strconv.Atoi(limitStr); e == nil {
			limitInt = v
		}
	}
	if offsetStr != "" {
		if v, e := strconv.Atoi(offsetStr); e == nil {
			offsetInt = v
		}
	}

	if limitInt > 0 {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(limitInt))
		if offsetInt > 0 {
			sb.WriteString(" OFFSET ")
			sb.WriteString(strconv.Itoa(offsetInt))
		}
	} else {
		// when limit is not set or is 0, only generate OFFSET when offset>0
		if offsetInt > 0 {
			sb.WriteString(" OFFSET ")
			sb.WriteString(strconv.Itoa(offsetInt))
		}
	}

	return sb.String(), nil
}

// ExportTemplate exports an Excel template
// Author [piexlmax](https://github.com/piexlmax)
func (sysExportTemplateService *SysExportTemplateService) ExportTemplate(templateID string) (file *bytes.Buffer, name string, err error) {
	var template system.SysExportTemplate
	err = global.GVA_DB.First(&template, "template_id = ?", templateID).Error
	if err != nil {
		return nil, "", err
	}
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	var templateInfoMap = make(map[string]string)

	columns, err := utils.GetJSONKeys(template.TemplateInfo)

	err = json.Unmarshal([]byte(template.TemplateInfo), &templateInfoMap)
	if err != nil {
		return nil, "", err
	}
	var tableTitle []string
	for _, key := range columns {
		tableTitle = append(tableTitle, templateInfoMap[key])
	}

	for i := range tableTitle {
		fErr := f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", getColumnName(i+1), 1), tableTitle[i])
		if fErr != nil {
			return nil, "", fErr
		}
	}
	f.SetActiveSheet(index)
	file, err = f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	return file, template.Name, nil
}

// helper function: check if a table has a deleted_at column
func (s *SysExportTemplateService) hasDeletedAtColumn(tableName string) bool {
	var count int64
	global.GVA_DB.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = ? AND COLUMN_NAME = 'deleted_at'", tableName).Count(&count)
	return count > 0
}

// ImportExcel imports from Excel
// Author [piexlmax](https://github.com/piexlmax)
func (sysExportTemplateService *SysExportTemplateService) ImportExcel(templateID string, file *multipart.FileHeader) (err error) {
	var template system.SysExportTemplate
	err = global.GVA_DB.First(&template, "template_id = ?", templateID).Error
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return err
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}
	if len(rows) < 2 {
		return errors.New("Excel data is not enough.\nIt should contain title row and data")
	}

	var templateInfoMap = make(map[string]string)
	err = json.Unmarshal([]byte(template.TemplateInfo), &templateInfoMap)
	if err != nil {
		return err
	}

	db := global.GVA_DB
	if template.DBName != "" {
		db = global.MustGetGlobalDBByDBName(template.DBName)
	}

	items, err := sysExportTemplateService.parseExcelToMap(rows, templateInfoMap)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if template.ImportSQL != "" {
			return sysExportTemplateService.importBySQL(tx, template.ImportSQL, items)
		}
		return sysExportTemplateService.importByGORM(tx, template.TableName, items)
	})
}

func (sysExportTemplateService *SysExportTemplateService) parseExcelToMap(rows [][]string, templateInfoMap map[string]string) ([]map[string]interface{}, error) {
	var titleKeyMap = make(map[string]string)
	for key, title := range templateInfoMap {
		titleKeyMap[title] = key
	}

	excelTitle := rows[0]
	for i, str := range excelTitle {
		excelTitle[i] = strings.TrimSpace(str)
	}
	values := rows[1:]
	items := make([]map[string]interface{}, 0, len(values))
	for _, row := range values {
		var item = make(map[string]interface{})
		for ii, value := range row {
			if ii >= len(excelTitle) {
				continue
			}
			if _, ok := titleKeyMap[excelTitle[ii]]; !ok {
				continue // extra title in excel with no corresponding field in template info, key is empty, must skip
			}
			key := titleKeyMap[excelTitle[ii]]
			item[key] = value
		}
		items = append(items, item)
	}
	return items, nil
}

func (sysExportTemplateService *SysExportTemplateService) importBySQL(tx *gorm.DB, sql string, items []map[string]interface{}) error {
	for _, item := range items {
		if err := tx.Exec(sql, item).Error; err != nil {
			return err
		}
	}
	return nil
}

func (sysExportTemplateService *SysExportTemplateService) importByGORM(tx *gorm.DB, tableName string, items []map[string]interface{}) error {
	needCreated := tx.Migrator().HasColumn(tableName, "created_at")
	needUpdated := tx.Migrator().HasColumn(tableName, "updated_at")

	for _, item := range items {
		if item["created_at"] == nil && needCreated {
			item["created_at"] = time.Now()
		}
		if item["updated_at"] == nil && needUpdated {
			item["updated_at"] = time.Now()
		}
	}
	return tx.Table(tableName).CreateInBatches(&items, 1000).Error
}

func getColumnName(n int) string {
	columnName := ""
	for n > 0 {
		n--
		columnName = string(rune('A'+n%26)) + columnName
		n /= 26
	}
	return columnName
}
