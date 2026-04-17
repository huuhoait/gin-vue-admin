package system

import (
	"context"
	sysModel "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/service/system"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderAuthority = initOrderCasbin + 1

type initAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderAuthority, &initAuthority{})
}

func (i *initAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysAuthority{})
}

func (i *initAuthority) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysAuthority{})
}

func (i *initAuthority) InitializerName() string {
	return sysModel.SysAuthority{}.TableName()
}

func (i *initAuthority) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []sysModel.SysAuthority{
		{AuthorityId: 888, AuthorityName: "Normal User", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 9528, AuthorityName: "Test Role", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 8881, AuthorityName: "Normal User Sub-role", ParentId: utils.Pointer[uint](888), DefaultRouter: "dashboard"},
		{AuthorityId: 9100, AuthorityName: "KYC Reviewer", ParentId: utils.Pointer[uint](0), DefaultRouter: "kyc"},
		{AuthorityId: 9200, AuthorityName: "Accountant", ParentId: utils.Pointer[uint](0), DefaultRouter: "commission"},
	}

	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrapf(err, "%stable data initialization failed!", sysModel.SysAuthority{}.TableName())
	}
	// data authority
	if err := db.Model(&entities[0]).Association("DataAuthorityId").Replace(
		[]*sysModel.SysAuthority{
			{AuthorityId: 888},
			{AuthorityId: 9528},
			{AuthorityId: 8881},
			{AuthorityId: 9100},
			{AuthorityId: 9200},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%stable data initialization failed!",
			db.Model(&entities[0]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}
	if err := db.Model(&entities[1]).Association("DataAuthorityId").Replace(
		[]*sysModel.SysAuthority{
			{AuthorityId: 9528},
			{AuthorityId: 8881},
			{AuthorityId: 9100},
			{AuthorityId: 9200},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%stable data initialization failed!",
			db.Model(&entities[1]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}

	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("authority_id = ?", "8881").
		First(&sysModel.SysAuthority{}).Error, gorm.ErrRecordNotFound) { // Check if data exists
		return false
	}
	return true
}
