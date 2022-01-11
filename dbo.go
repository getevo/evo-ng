package evo

import (
	"github.com/getevo/evo-ng/lib/generic"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"strings"
	"time"
)

var Database *gorm.DB

func setupDatabase() {
	Events.Trigger("database.setup")
	config := Config.Database
	var err error
	if config.Enabled == false {
		return
	}
	var logLevel = logger.Silent
	config.Debug = strings.ToLower(config.Debug)

	switch config.Debug {
	case "true", "all", "*", "any":
		logLevel = logger.Info
	case "warn", "warning":
		logLevel = logger.Warn
	case "err", "error":
		logLevel = logger.Error
	default:
		logLevel = logger.Silent
	}

	if config.Debug == "true" || config.Debug == "all" {
		logLevel = logger.Info
	}
	if config.Debug == "warn" || config.Debug == "warning" {
		logLevel = logger.Warn
	}
	if config.Debug == "err" || config.Debug == "error" {
		logLevel = logger.Error
	}
	var newLog = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logLevel,    // Log level
			Colorful:      true,        // Disable color
		},
	)
	cfg := &gorm.Config{
		Logger:                                   newLog,
		PrepareStmt:                              config.StmtCache,
		SkipDefaultTransaction:                   config.SkipDefaultTransaction,
		FullSaveAssociations:                     config.FullSaveAssociations,
		DisableAutomaticPing:                     config.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: config.DisableForeignKeyConstraintWhenMigrating,
		QueryFields:                              config.QueryFields,
		CreateBatchSize:                          config.CreateBatchSize,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.TablePrefix,
		},
	}

	switch strings.ToLower(config.Dialect) {
	case "mysql":
		Database, err = gorm.Open(mysql.Open(config.DSN), cfg)
	case "postgres":
		Database, err = gorm.Open(postgres.Open(config.DSN), cfg)
	case "mssql":
		Database, err = gorm.Open(sqlserver.Open(config.DSN), cfg)
	case "sqlite":
		Database, err = gorm.Open(sqlite.Open(config.DSN), cfg)
	default:
		Panic("invalid database dialect passed to configuration")
	}

	if err != nil {
		Panic(err)
		return
	}
	var connMaxLifeTime, connMaxIdleTime time.Duration
	var db, _ = Database.DB()

	connMaxLifeTime, err = generic.Parse(config.ConnMaxLifeTime).Duration()
	if err != nil {
		Panic("invalid db connection_max_lifetime")
		return
	}
	if connMaxLifeTime < 1 {
		connMaxLifeTime = 1 * time.Hour
	}

	connMaxIdleTime, err = generic.Parse(config.ConnMaxIdleTime).Duration()
	if err != nil {
		Panic("invalid db connection_max_idle_time")
		return
	}
	if connMaxIdleTime < 1 {
		connMaxIdleTime = 1 * time.Hour
	}
	if config.MaxIdleConns < 1 {
		config.MaxIdleConns = 1
	}
	if config.MaxOpenConns < 1 {
		config.MaxOpenConns = 1
	}
	if config.MaxIdleConns > config.MaxOpenConns {
		config.MaxOpenConns = config.MaxIdleConns
	}
	db.SetConnMaxLifetime(connMaxLifeTime)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	if len(config.Replicas) > 0 {
		var dialectics []gorm.Dialector
		for _, item := range config.Replicas {
			dialectics = append(dialectics, mysql.Open(item))
		}
		err := Database.Use(dbresolver.Register(dbresolver.Config{
			Replicas: dialectics,
		}))
		if err != nil {
			Panic(err)
		}
	}
	Events.Trigger("database.connected")

}

// GetDBO return database object instance
func GetDBO(name ...string) *gorm.DB {
	if Database == nil {
		setupDatabase()
	}
	return Database
}

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}
