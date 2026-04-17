package system

import (
	"context"
	sysModel "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderDict = initOrderCasbin + 1

type initDict struct{}

// auto run
func init() {
	system.RegisterInit(initOrderDict, &initDict{})
}

func (i *initDict) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysDictionary{})
}

func (i *initDict) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysDictionary{})
}

func (i *initDict) InitializerName() string {
	return sysModel.SysDictionary{}.TableName()
}

func (i *initDict) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	True := true
	entities := []sysModel.SysDictionary{
		{Name: "Gender", Type: "gender", Status: &True, Desc: "GenderDictionary"},
		{Name: "Database int type", Type: "int", Status: &True, Desc: "inttypeCorrespondingDatabasetype"},
		{Name: "database date/timetype", Type: "time.Time", Status: &True, Desc: "database date/timetype"},
		{Name: "database float", Type: "float64", Status: &True, Desc: "database float"},
		{Name: "database string", Type: "string", Status: &True, Desc: "database string"},
		{Name: "Database bool type", Type: "bool", Status: &True, Desc: "Database bool type"},
	}

	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysDictionary{}.TableName()+"table data initialization failed!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initDict) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("type = ?", "bool").First(&sysModel.SysDictionary{}).Error, gorm.ErrRecordNotFound) { // Check if data exists
		return false
	}
	return true
}
