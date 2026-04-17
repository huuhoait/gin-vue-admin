package system

import (
	"fmt"
	"strconv"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateSysDictionaryDetail
//@description: Create dictionary detail data
//@param: sysDictionaryDetail model.SysDictionaryDetail
//@return: err error

type DictionaryDetailService struct{}

var DictionaryDetailServiceApp = new(DictionaryDetailService)

func (dictionaryDetailService *DictionaryDetailService) CreateSysDictionaryDetail(sysDictionaryDetail system.SysDictionaryDetail) (err error) {
	// calculate level and path
	if sysDictionaryDetail.ParentID != nil {
		var parent system.SysDictionaryDetail
		err = global.GVA_DB.First(&parent, *sysDictionaryDetail.ParentID).Error
		if err != nil {
			return err
		}
		sysDictionaryDetail.Level = parent.Level + 1
		if parent.Path == "" {
			sysDictionaryDetail.Path = strconv.Itoa(int(parent.ID))
		} else {
			sysDictionaryDetail.Path = parent.Path + "," + strconv.Itoa(int(parent.ID))
		}
	} else {
		sysDictionaryDetail.Level = 0
		sysDictionaryDetail.Path = ""
	}

	err = global.GVA_DB.Create(&sysDictionaryDetail).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteSysDictionaryDetail
//@description: Delete dictionary detail data
//@param: sysDictionaryDetail model.SysDictionaryDetail
//@return: err error

func (dictionaryDetailService *DictionaryDetailService) DeleteSysDictionaryDetail(sysDictionaryDetail system.SysDictionaryDetail) (err error) {
	// check if there are child items
	var count int64
	err = global.GVA_DB.Model(&system.SysDictionaryDetail{}).Where("parent_id = ?", sysDictionaryDetail.ID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("this dictionary detail has child items and cannot be deleted")
	}

	err = global.GVA_DB.Delete(&sysDictionaryDetail).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateSysDictionaryDetail
//@description: Update dictionary detail data
//@param: sysDictionaryDetail *model.SysDictionaryDetail
//@return: err error

func (dictionaryDetailService *DictionaryDetailService) UpdateSysDictionaryDetail(sysDictionaryDetail *system.SysDictionaryDetail) (err error) {
	// if parent ID is updated, recalculate level and path
	if sysDictionaryDetail.ParentID != nil {
		var parent system.SysDictionaryDetail
		err = global.GVA_DB.First(&parent, *sysDictionaryDetail.ParentID).Error
		if err != nil {
			return err
		}

		// check circular reference
		if dictionaryDetailService.checkCircularReference(sysDictionaryDetail.ID, *sysDictionaryDetail.ParentID) {
			return fmt.Errorf("cannot set a dictionary detail as the parent of itself or its children")
		}

		sysDictionaryDetail.Level = parent.Level + 1
		if parent.Path == "" {
			sysDictionaryDetail.Path = strconv.Itoa(int(parent.ID))
		} else {
			sysDictionaryDetail.Path = parent.Path + "," + strconv.Itoa(int(parent.ID))
		}
	} else {
		sysDictionaryDetail.Level = 0
		sysDictionaryDetail.Path = ""
	}

	err = global.GVA_DB.Save(sysDictionaryDetail).Error
	if err != nil {
		return err
	}

	// update level and path of all child items
	return dictionaryDetailService.updateChildrenLevelAndPath(sysDictionaryDetail.ID)
}

// checkCircularReference checks for circular references
func (dictionaryDetailService *DictionaryDetailService) checkCircularReference(id, parentID uint) bool {
	if id == parentID {
		return true
	}

	var parent system.SysDictionaryDetail
	err := global.GVA_DB.First(&parent, parentID).Error
	if err != nil {
		return false
	}

	if parent.ParentID == nil {
		return false
	}

	return dictionaryDetailService.checkCircularReference(id, *parent.ParentID)
}

// updateChildrenLevelAndPath updates the level and path of child items
func (dictionaryDetailService *DictionaryDetailService) updateChildrenLevelAndPath(parentID uint) error {
	var children []system.SysDictionaryDetail
	err := global.GVA_DB.Where("parent_id = ?", parentID).Find(&children).Error
	if err != nil {
		return err
	}

	var parent system.SysDictionaryDetail
	err = global.GVA_DB.First(&parent, parentID).Error
	if err != nil {
		return err
	}

	for _, child := range children {
		child.Level = parent.Level + 1
		if parent.Path == "" {
			child.Path = strconv.Itoa(int(parent.ID))
		} else {
			child.Path = parent.Path + "," + strconv.Itoa(int(parent.ID))
		}

		err = global.GVA_DB.Save(&child).Error
		if err != nil {
			return err
		}

		// recursively update children of children
		err = dictionaryDetailService.updateChildrenLevelAndPath(child.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSysDictionaryDetail
//@description: Get a single dictionary detail record by id
//@param: id uint
//@return: sysDictionaryDetail system.SysDictionaryDetail, err error

func (dictionaryDetailService *DictionaryDetailService) GetSysDictionaryDetail(id uint) (sysDictionaryDetail system.SysDictionaryDetail, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&sysDictionaryDetail).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSysDictionaryDetailInfoList
//@description: Get dictionary detail list with pagination
//@param: info request.SysDictionaryDetailSearch
//@return: list interface{}, total int64, err error

func (dictionaryDetailService *DictionaryDetailService) GetSysDictionaryDetailInfoList(info request.SysDictionaryDetailSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// create db
	db := global.GVA_DB.Model(&system.SysDictionaryDetail{})
	var sysDictionaryDetails []system.SysDictionaryDetail
	// if there are search conditions, the search statement will be automatically created below
	if info.Label != "" {
		db = db.Where("label LIKE ?", "%"+info.Label+"%")
	}
	if info.Value != "" {
		db = db.Where("value = ?", info.Value)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	if info.SysDictionaryID != 0 {
		db = db.Where("sys_dictionary_id = ?", info.SysDictionaryID)
	}
	if info.ParentID != nil {
		db = db.Where("parent_id = ?", *info.ParentID)
	}
	if info.Level != nil {
		db = db.Where("level = ?", *info.Level)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("sort").Order("id").Find(&sysDictionaryDetails).Error
	return sysDictionaryDetails, total, err
}

// GetDictionaryList gets all dictionary content by dictionary id
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryList(dictionaryID uint) (list []system.SysDictionaryDetail, err error) {
	var sysDictionaryDetails []system.SysDictionaryDetail
	err = global.GVA_DB.Find(&sysDictionaryDetails, "sys_dictionary_id = ?", dictionaryID).Error
	return sysDictionaryDetails, err
}

// GetDictionaryTreeList gets dictionary tree structure list
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryTreeList(dictionaryID uint) (list []system.SysDictionaryDetail, err error) {
	var sysDictionaryDetails []system.SysDictionaryDetail
	// only get top-level items (parent_id is null)
	err = global.GVA_DB.Where("sys_dictionary_id = ? AND parent_id IS NULL", dictionaryID).Order("sort").Find(&sysDictionaryDetails).Error
	if err != nil {
		return nil, err
	}

	// recursively load child items and set disabled property
	for i := range sysDictionaryDetails {
		// set disabled property: when status is false, disabled is true
		if sysDictionaryDetails[i].Status != nil {
			sysDictionaryDetails[i].Disabled = !*sysDictionaryDetails[i].Status
		} else {
			sysDictionaryDetails[i].Disabled = false // not disabled by default
		}
		
		err = dictionaryDetailService.loadChildren(&sysDictionaryDetails[i])
		if err != nil {
			return nil, err
		}
	}

	return sysDictionaryDetails, nil
}

// loadChildren recursively loads child items
func (dictionaryDetailService *DictionaryDetailService) loadChildren(detail *system.SysDictionaryDetail) error {
	var children []system.SysDictionaryDetail
	err := global.GVA_DB.Where("parent_id = ?", detail.ID).Order("sort").Find(&children).Error
	if err != nil {
		return err
	}

	for i := range children {
		// set disabled property: when status is false, disabled is true
		if children[i].Status != nil {
			children[i].Disabled = !*children[i].Status
		} else {
			children[i].Disabled = false // not disabled by default
		}
		
		err = dictionaryDetailService.loadChildren(&children[i])
		if err != nil {
			return err
		}
	}

	detail.Children = children
	return nil
}

// GetDictionaryDetailsByParent gets dictionary details by parent ID
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryDetailsByParent(req request.GetDictionaryDetailsByParentRequest) (list []system.SysDictionaryDetail, err error) {
	db := global.GVA_DB.Model(&system.SysDictionaryDetail{}).Where("sys_dictionary_id = ?", req.SysDictionaryID)

	if req.ParentID != nil {
		db = db.Where("parent_id = ?", *req.ParentID)
	} else {
		db = db.Where("parent_id IS NULL")
	}

	err = db.Order("sort").Find(&list).Error
	if err != nil {
		return list, err
	}

	// set disabled property
	for i := range list {
		if list[i].Status != nil {
			list[i].Disabled = !*list[i].Status
		} else {
			list[i].Disabled = false // not disabled by default
		}
	}

	// if children need to be included, recursively load all levels of child items
	if req.IncludeChildren {
		for i := range list {
			err = dictionaryDetailService.loadChildren(&list[i])
			if err != nil {
				return list, err
			}
		}
	}

	return list, err
}

// GetDictionaryListByType gets all dictionary content by dictionary type
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryListByType(t string) (list []system.SysDictionaryDetail, err error) {
	var sysDictionaryDetails []system.SysDictionaryDetail
	db := global.GVA_DB.Model(&system.SysDictionaryDetail{}).Joins("JOIN sys_dictionaries ON sys_dictionaries.id = sys_dictionary_details.sys_dictionary_id")
	err = db.Find(&sysDictionaryDetails, "type = ?", t).Error
	return sysDictionaryDetails, err
}

// GetDictionaryTreeListByType gets tree structure by dictionary type
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryTreeListByType(t string) (list []system.SysDictionaryDetail, err error) {
	var sysDictionaryDetails []system.SysDictionaryDetail
	db := global.GVA_DB.Model(&system.SysDictionaryDetail{}).
		Joins("JOIN sys_dictionaries ON sys_dictionaries.id = sys_dictionary_details.sys_dictionary_id").
		Where("sys_dictionaries.type = ? AND sys_dictionary_details.parent_id IS NULL", t).
		Order("sys_dictionary_details.sort")

	err = db.Find(&sysDictionaryDetails).Error
	if err != nil {
		return nil, err
	}

	// recursively load child items and set disabled property
	for i := range sysDictionaryDetails {
		// set disabled property: when status is false, disabled is true
		if sysDictionaryDetails[i].Status != nil {
			sysDictionaryDetails[i].Disabled = !*sysDictionaryDetails[i].Status
		} else {
			sysDictionaryDetails[i].Disabled = false // not disabled by default
		}
		
		err = dictionaryDetailService.loadChildren(&sysDictionaryDetails[i])
		if err != nil {
			return nil, err
		}
	}

	return sysDictionaryDetails, nil
}

// GetDictionaryInfoByValue gets a single dictionary detail by dictionary id and value
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryInfoByValue(dictionaryID uint, value string) (detail system.SysDictionaryDetail, err error) {
	var sysDictionaryDetail system.SysDictionaryDetail
	err = global.GVA_DB.First(&sysDictionaryDetail, "sys_dictionary_id = ? and value = ?", dictionaryID, value).Error
	return sysDictionaryDetail, err
}

// GetDictionaryInfoByTypeValue gets a single dictionary detail by dictionary type and value
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryInfoByTypeValue(t string, value string) (detail system.SysDictionaryDetail, err error) {
	var sysDictionaryDetails system.SysDictionaryDetail
	db := global.GVA_DB.Model(&system.SysDictionaryDetail{}).Joins("JOIN sys_dictionaries ON sys_dictionaries.id = sys_dictionary_details.sys_dictionary_id")
	err = db.First(&sysDictionaryDetails, "sys_dictionaries.type = ? and sys_dictionary_details.value = ?", t, value).Error
	return sysDictionaryDetails, err
}

// GetDictionaryPath gets the full path of a dictionary detail
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryPath(id uint) (path []system.SysDictionaryDetail, err error) {
	var detail system.SysDictionaryDetail
	err = global.GVA_DB.First(&detail, id).Error
	if err != nil {
		return nil, err
	}

	path = append(path, detail)

	if detail.ParentID != nil {
		parentPath, err := dictionaryDetailService.GetDictionaryPath(*detail.ParentID)
		if err != nil {
			return nil, err
		}
		path = append(parentPath, path...)
	}

	return path, nil
}

// GetDictionaryPathByValue gets the full path of a dictionary detail by value
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryPathByValue(dictionaryID uint, value string) (path []system.SysDictionaryDetail, err error) {
	detail, err := dictionaryDetailService.GetDictionaryInfoByValue(dictionaryID, value)
	if err != nil {
		return nil, err
	}

	return dictionaryDetailService.GetDictionaryPath(detail.ID)
}
