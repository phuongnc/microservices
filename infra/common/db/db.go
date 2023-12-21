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
	db, err := gorm.Open(mysql.Open(c.Uri), &gorm.Config{})
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
