package system

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"gorm.io/gorm"
	"sort"
)

const (
	Mysql           = "mysql"
	Pgsql           = "pgsql"
	Sqlite          = "sqlite"
	Mssql           = "mssql"
	InitSuccess     = "\n[%v] --> Initial data setup succeeded!\n"
	InitDataExist   = "\n[%v] --> Initial data for %v already exists!\n"
	InitDataFailed  = "\n[%v] --> %v initial data setup failed! \nerr: %+v\n"
	InitDataSuccess = "\n[%v] --> %v initial data setup succeeded!\n"
)

const (
	InitOrderSystem   = 10
	InitOrderInternal = 1000
	InitOrderExternal = 100000
)

var (
	ErrMissingDBContext        = errors.New("missing db in context")
	ErrMissingDependentContext = errors.New("missing dependent value in context")
	ErrDBTypeMismatch          = errors.New("db type mismatch")
)

// SubInitializer provides the interface used by source/*/init(), each initializer completes one initialization process
type SubInitializer interface {
	InitializerName() string // does not necessarily represent a single table, hence the broader semantics
	MigrateTable(ctx context.Context) (next context.Context, err error)
	InitializeData(ctx context.Context) (next context.Context, err error)
	TableCreated(ctx context.Context) bool
	DataInserted(ctx context.Context) bool
}

// TypedDBInitHandler executes the provided initializers
type TypedDBInitHandler interface {
	EnsureDB(ctx context.Context, conf *request.InitDB) (context.Context, error) // create database; failure is a fatal error
	WriteConfig(ctx context.Context) error                                       // write back configuration
	InitTables(ctx context.Context, inits initSlice) error                       // create tables handler
	InitData(ctx context.Context, inits initSlice) error                         // seed data handler
}

// orderedInitializer combines an order field for sorting
type orderedInitializer struct {
	order int
	SubInitializer
}

// initSlice used for sorting initializer dependencies
type initSlice []*orderedInitializer

var (
	initializers initSlice
	cache        map[string]*orderedInitializer
)

// RegisterInit registers an initialization process to be executed when InitDB() is called
func RegisterInit(order int, i SubInitializer) {
	if initializers == nil {
		initializers = initSlice{}
	}
	if cache == nil {
		cache = map[string]*orderedInitializer{}
	}
	name := i.InitializerName()
	if _, existed := cache[name]; existed {
		panic(fmt.Sprintf("Name conflict on %s", name))
	}
	ni := orderedInitializer{order, i}
	initializers = append(initializers, &ni)
	cache[name] = &ni
}

/* ---- * service * ---- */

type InitDBService struct{}

// InitDB creates the database and initializes data - main entry point
func (initDBService *InitDBService) InitDB(conf request.InitDB) (err error) {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, "adminPassword", conf.AdminPassword)
	if len(initializers) == 0 {
		return errors.New("no available initialization process, please check if initialization has already completed")
	}
	sort.Sort(&initializers) // ensure initializers with dependencies are executed later
	// Note: if an initializer has a single dependency, e.g. B=A+1, C=A+1; since BC have no dependency on each other, order doesn't matter
	// For multiple dependencies, e.g. C=A+B, D=A+B+C, E=A+1;
	// C is necessarily > A|B so it runs after AB, D > A|B|C so after ABC, E only depends on A so its order relative to CD doesn't matter
	var initHandler TypedDBInitHandler
	switch conf.DBType {
	case "mysql":
		initHandler = NewMysqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "mysql")
	case "pgsql":
		initHandler = NewPgsqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "pgsql")
	case "sqlite":
		initHandler = NewSqliteInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "sqlite")
	case "mssql":
		initHandler = NewMssqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "mssql")
	default:
		initHandler = NewMysqlInitHandler()
		ctx = context.WithValue(ctx, "dbtype", "mysql")
	}
	ctx, err = initHandler.EnsureDB(ctx, &conf)
	if err != nil {
		return err
	}

	db := ctx.Value("db").(*gorm.DB)
	global.GVA_DB = db

	if err = initHandler.InitTables(ctx, initializers); err != nil {
		return err
	}
	if err = initHandler.InitData(ctx, initializers); err != nil {
		return err
	}

	if err = initHandler.WriteConfig(ctx); err != nil {
		return err
	}
	initializers = initSlice{}
	cache = map[string]*orderedInitializer{}
	return nil
}

// createDatabase creates a database (called from EnsureDB())
func createDatabase(dsn string, driver string, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}

// createTables creates tables (default dbInitHandler.initTables behavior)
func createTables(ctx context.Context, inits initSlice) error {
	next, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, init := range inits {
		if init.TableCreated(next) {
			continue
		}
		if n, err := init.MigrateTable(next); err != nil {
			return err
		} else {
			next = n
		}
	}
	return nil
}

/* -- sortable interface -- */

func (a initSlice) Len() int {
	return len(a)
}

func (a initSlice) Less(i, j int) bool {
	return a[i].order < a[j].order
}

func (a initSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
