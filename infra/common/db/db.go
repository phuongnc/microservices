package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Uri                   string
	Driver                string
	Dialect               string
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionLifetime time.Duration
}

func NewSQL(c *DatabaseConfig) (*gorm.DB, error) {
	//db, err := gorm.Open(mysql.Open(c.Uri), &gorm.Config{})

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       c.Uri, // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(c.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(c.MaxConnectionLifetime)

	return db, nil
}
