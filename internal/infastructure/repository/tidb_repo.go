package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RepoConfig struct {
	DSN          string
	Debug        bool
	MaxIdleConns int
	MaxOpenConns int
}

func ConnPoolConfig(config RepoConfig, db *gorm.DB) error {
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}
	if config.MaxIdleConns != 0 {
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDb.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxOpenConns != 0 {
		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDb.SetMaxOpenConns(config.MaxOpenConns)
	}
	return nil
}

// ConnectDb open connection to db
func NewTiDb(config RepoConfig) error {
	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  config.DSN,
			PreferSimpleProtocol: true,
		}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}
	ConnPoolConfig(config, db)
	return nil
}
