package request

type SysDictionarySearch struct {
	Name string `json:"name" form:"name" gorm:"column:name;comment:dictionary name"` // dictionary name
}

type ImportSysDictionaryRequest struct {
	Json string `json:"json" binding:"required"` // JSONString
}
