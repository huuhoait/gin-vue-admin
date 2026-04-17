package system

import (
	"context"

	sysModel "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenuAuthority = initOrderMenu + initOrderAuthority

type initMenuAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initMenuAuthority) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i *initMenuAuthority) InitializerName() string {
	return "sys_menu_authorities"
}

func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	initAuth := &initAuthority{}
	authorities, ok := ctx.Value(initAuth.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "Create [Menu-Permission] Associationfailed, permission table seed data not found")
	}

	allMenus, ok := ctx.Value(new(initMenu).InitializerName()).([]sysModel.SysBaseMenu)
	if !ok {
		return next, errors.Wrap(errors.New(""), "Create [Menu-Permission] Associationfailed, not foundMenutable seed data")
	}
	next = ctx

	// Buildmenu IDMap, ConvenientfastLookup
	menuMap := make(map[uint]sysModel.SysBaseMenu)
	for _, menu := range allMenus {
		menuMap[menu.ID] = menu
	}
	authorityById := make(map[uint]*sysModel.SysAuthority)
	for idx := range authorities {
		authorityById[authorities[idx].AuthorityId] = &authorities[idx]
	}

	// ForDifferentRolePartMatchDifferentPermission
	// 1. Super AdminRole(888) - OwnHaveAllMenuPermission
	superAdmin := authorityById[888]
	if superAdmin == nil {
		return next, errors.Wrap(errors.New("missing authority 888"), "ForSuper AdminPartMatchMenufailed")
	}
	if err = db.Model(superAdmin).Association("SysBaseMenus").Replace(allMenus); err != nil {
		return next, errors.Wrap(err, "ForSuper AdminPartMatchMenufailed")
	}

	// 2. Normal UserRole(8881) - OnlyOwnHaveBasicFunctionMenu
	// OnlySelectPartParentMenuAnd Itschild menus
	var menu8881 []sysModel.SysBaseMenu

	// AddDashboard, About UsAndPersonal InfoMenu
	for _, menu := range allMenus {
		if menu.ParentId == 0 && (menu.Name == "dashboard" || menu.Name == "about" || menu.Name == "person" || menu.Name == "state") {
			menu8881 = append(menu8881, menu)
		}
	}

	normalUser := authorityById[8881]
	if normalUser == nil {
		return next, errors.Wrap(errors.New("missing authority 8881"), "ForNormal UserPartMatchMenufailed")
	}
	if err = db.Model(normalUser).Association("SysBaseMenus").Replace(menu8881); err != nil {
		return next, errors.Wrap(err, "ForNormal UserPartMatchMenufailed")
	}

	// 3. Test Role(9528) - OwnHavePartMenuPermission
	var menu9528 []sysModel.SysBaseMenu

	// AddAllParentMenu
	for _, menu := range allMenus {
		if menu.ParentId == 0 {
			menu9528 = append(menu9528, menu)
		}
	}

	// AddPartchild menus - SystemUtility, Examplesetc.Moduleofchild menus
	for _, menu := range allMenus {
		parentName := ""
		if menu.ParentId > 0 && menuMap[menu.ParentId].Name != "" {
			parentName = menuMap[menu.ParentId].Name
		}

		if menu.ParentId > 0 && (parentName == "systemTools" || parentName == "example") {
			menu9528 = append(menu9528, menu)
		}
	}

	testRole := authorityById[9528]
	if testRole == nil {
		return next, errors.Wrap(errors.New("missing authority 9528"), "ForTest RolePartMatchMenufailed")
	}
	if err = db.Model(testRole).Association("SysBaseMenus").Replace(menu9528); err != nil {
		return next, errors.Wrap(err, "ForTest RolePartMatchMenufailed")
	}

	// 4. KYC Reviewer (9100) - only SkyAgent/KYC menus
	var reviewerMenus []sysModel.SysBaseMenu
	for _, menu := range allMenus {
		if menu.Name == "skyagent" || menu.Name == "kyc" || menu.Name == "kycCases" {
			reviewerMenus = append(reviewerMenus, menu)
		}
	}
	kycReviewer := authorityById[9100]
	if kycReviewer != nil {
		if err = db.Model(kycReviewer).Association("SysBaseMenus").Replace(reviewerMenus); err != nil {
			return next, errors.Wrap(err, "ForKYC ReviewerPartMatchMenufailed")
		}
	}

	// 5. Accountant (9200) - only SkyAgent/Commission menu
	var accountantMenus []sysModel.SysBaseMenu
	for _, menu := range allMenus {
		if menu.Name == "skyagent" || menu.Name == "commission" {
			accountantMenus = append(accountantMenus, menu)
		}
	}
	accountant := authorityById[9200]
	if accountant != nil {
		if err = db.Model(accountant).Association("SysBaseMenus").Replace(accountantMenus); err != nil {
			return next, errors.Wrap(err, "ForAccountantPartMatchMenufailed")
		}
	}

	return next, nil
}

func (i *initMenuAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	auth := &sysModel.SysAuthority{}
	if ret := db.Model(auth).
		Where("authority_id = ?", 9528).Preload("SysBaseMenus").Find(auth); ret != nil {
		if ret.Error != nil {
			return false
		}
		return len(auth.SysBaseMenus) > 0
	}
	return false
}
