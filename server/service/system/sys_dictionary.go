package system

import (
	"encoding/json"
	"errors"

	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateSysDictionary
//@description: Create dictionary data
//@param: sysDictionary model.SysDictionary
//@return: err error

type DictionaryService struct{}

var DictionaryServiceApp = new(DictionaryService)

func (dictionaryService *DictionaryService) CreateSysDictionary(sysDictionary system.SysDictionary) (err error) {
	if (!errors.Is(global.GVA_DB.First(&system.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound)) {
		return errors.New("duplicate type already exists, creation not allowed")
	}
	err = global.GVA_DB.Create(&sysDictionary).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteSysDictionary
//@description: Delete dictionary data
//@param: sysDictionary model.SysDictionary
//@return: err error

func (dictionaryService *DictionaryService) DeleteSysDictionary(sysDictionary system.SysDictionary) (err error) {
	err = global.GVA_DB.Where("id = ?", sysDictionary.ID).Preload("SysDictionaryDetails").First(&sysDictionary).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}
	if err != nil {
		return err
	}
	err = global.GVA_DB.Delete(&sysDictionary).Error
	if err != nil {
		return err
	}

	if sysDictionary.SysDictionaryDetails != nil {
		return global.GVA_DB.Where("sys_dictionary_id=?", sysDictionary.ID).Delete(sysDictionary.SysDictionaryDetails).Error
	}
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateSysDictionary
//@description: Update dictionary data
//@param: sysDictionary *model.SysDictionary
//@return: err error

func (dictionaryService *DictionaryService) UpdateSysDictionary(sysDictionary *system.SysDictionary) (err error) {
	var dict system.SysDictionary
	sysDictionaryMap := map[string]interface{}{
		"Name":     sysDictionary.Name,
		"Type":     sysDictionary.Type,
		"Status":   sysDictionary.Status,
		"Desc":     sysDictionary.Desc,
		"ParentID": sysDictionary.ParentID,
	}
	err = global.GVA_DB.Where("id = ?", sysDictionary.ID).First(&dict).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return errors.New("failed to query dictionary data")
	}
	if dict.Type != sysDictionary.Type {
		if !errors.Is(global.GVA_DB.First(&system.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound) {
			return errors.New("duplicate type already exists, creation not allowed")
		}
	}

	// check if it would form a circular reference
	if sysDictionary.ParentID != nil && *sysDictionary.ParentID != 0 {
		if err := dictionaryService.checkCircularReference(sysDictionary.ID, *sysDictionary.ParentID); err != nil {
			return err
		}
	}

	err = global.GVA_DB.Model(&dict).Updates(sysDictionaryMap).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSysDictionary
//@description: Get a single dictionary record by id or type
//@param: Type string, Id uint
//@return: err error, sysDictionary model.SysDictionary

func (dictionaryService *DictionaryService) GetSysDictionary(Type string, Id uint, status *bool) (sysDictionary system.SysDictionary, err error) {
	var flag = false
	if status == nil {
		flag = true
	} else {
		flag = *status
	}
	err = global.GVA_DB.Where("(type = ? OR id = ?) and status = ?", Type, Id, flag).Preload("SysDictionaryDetails", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ? and deleted_at is null", true).Order("sort")
	}).First(&sysDictionary).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: GetSysDictionaryInfoList
//@description: Get dictionary list with pagination
//@param: info request.SysDictionarySearch
//@return: err error, list interface{}, total int64

func (dictionaryService *DictionaryService) GetSysDictionaryInfoList(c *gin.Context, req request.SysDictionarySearch) (list interface{}, err error) {
	var sysDictionarys []system.SysDictionary
	query := global.GVA_DB.WithContext(c)
	if req.Name != "" {
		query = query.Where("name LIKE ? OR type LIKE ?", "%"+req.Name+"%", "%"+req.Name+"%")
	}
	// preload child dictionaries
	query = query.Preload("Children")
	err = query.Find(&sysDictionarys).Error
	return sysDictionarys, err
}

// checkCircularReference checks if a circular reference would be formed
func (dictionaryService *DictionaryService) checkCircularReference(currentID uint, parentID uint) error {
	if currentID == parentID {
		return errors.New("cannot set a dictionary as its own parent")
	}

	// recursively check the parent chain
	var parent system.SysDictionary
	err := global.GVA_DB.Where("id = ?", parentID).First(&parent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // parent does not exist, allow setting
		}
		return err
	}

	// if the parent has a parent, continue checking
	if parent.ParentID != nil && *parent.ParentID != 0 {
		return dictionaryService.checkCircularReference(currentID, *parent.ParentID)
	}

	return nil
}

//@author: [pixelMax]
//@function: ExportSysDictionary
//@description: Export dictionary as JSON (including dictionary details)
//@param: id uint
//@return: exportData map[string]interface{}, err error

func (dictionaryService *DictionaryService) ExportSysDictionary(id uint) (exportData map[string]interface{}, err error) {
	var dictionary system.SysDictionary
	// query dictionary and all its details
	err = global.GVA_DB.Where("id = ?", id).Preload("SysDictionaryDetails", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort")
	}).First(&dictionary).Error
	if err != nil {
		return nil, err
	}

	// clear ID, created_at, updated_at and other fields from dictionary details
	var cleanDetails []map[string]interface{}
	for _, detail := range dictionary.SysDictionaryDetails {
		cleanDetail := map[string]interface{}{
			"label":  detail.Label,
			"value":  detail.Value,
			"extend": detail.Extend,
			"status": detail.Status,
			"sort":   detail.Sort,
			"level":  detail.Level,
			"path":   detail.Path,
		}
		cleanDetails = append(cleanDetails, cleanDetail)
	}

	// construct export data
	exportData = map[string]interface{}{
		"name":                 dictionary.Name,
		"type":                 dictionary.Type,
		"status":               dictionary.Status,
		"desc":                 dictionary.Desc,
		"sysDictionaryDetails": cleanDetails,
	}

	return exportData, nil
}

//@author: [pixelMax]
//@function: ImportSysDictionary
//@description: Import dictionary from JSON (including dictionary details)
//@param: jsonStr string
//@return: err error

func (dictionaryService *DictionaryService) ImportSysDictionary(jsonStr string) error {
	// parse directly into SysDictionary struct
	var importData system.SysDictionary
	if err := json.Unmarshal([]byte(jsonStr), &importData); err != nil {
		return errors.New("invalid JSON format: " + err.Error())
	}

	// validate required fields
	if importData.Name == "" {
		return errors.New("dictionary name cannot be empty")
	}
	if importData.Type == "" {
		return errors.New("dictionary type cannot be empty")
	}

	// check if dictionary type already exists
	if !errors.Is(global.GVA_DB.First(&system.SysDictionary{}, "type = ?", importData.Type).Error, gorm.ErrRecordNotFound) {
		return errors.New("duplicate type already exists, import not allowed")
	}

	// create dictionary (clear imported data IDs and timestamps)
	dictionary := system.SysDictionary{
		Name:   importData.Name,
		Type:   importData.Type,
		Status: importData.Status,
		Desc:   importData.Desc,
	}

	// start transaction
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// create dictionary
		if err := tx.Create(&dictionary).Error; err != nil {
			return err
		}

		// process dictionary details
		if len(importData.SysDictionaryDetails) > 0 {
			// create a mapping to track old ID to new ID correspondence
			idMap := make(map[uint]uint)

			// first pass: create all detail records
			for _, detail := range importData.SysDictionaryDetails {
				// validate required fields
				if detail.Label == "" || detail.Value == "" {
					continue
				}

				// record old ID
				oldID := detail.ID

				// create new detail record (ID will be auto-set by GORM)
				detailRecord := system.SysDictionaryDetail{
					Label:           detail.Label,
					Value:           detail.Value,
					Extend:          detail.Extend,
					Status:          detail.Status,
					Sort:            detail.Sort,
					Level:           detail.Level,
					Path:            detail.Path,
					SysDictionaryID: int(dictionary.ID),
				}

				// create detail record
				if err := tx.Create(&detailRecord).Error; err != nil {
					return err
				}

				// record old ID to new ID mapping
				if oldID > 0 {
					idMap[oldID] = detailRecord.ID
				}
			}

			// second pass: update parent_id relationships
			for _, detail := range importData.SysDictionaryDetails {
				if detail.ParentID != nil && *detail.ParentID > 0 && detail.ID > 0 {
					if newID, exists := idMap[detail.ID]; exists {
						if newParentID, parentExists := idMap[*detail.ParentID]; parentExists {
							if err := tx.Model(&system.SysDictionaryDetail{}).
								Where("id = ?", newID).
								Update("parent_id", newParentID).Error; err != nil {
								return err
							}
						}
					}
				}
			}
		}

		return nil
	})
}
