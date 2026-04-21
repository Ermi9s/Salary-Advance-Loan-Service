package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStore struct {
	DB    *gorm.DB
	sqlDB *sql.DB
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetMaxOpenConns(15)
	sqlDB.SetMaxIdleConns(5)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	if err := db.WithContext(ctx).AutoMigrate(&UserModel{}, &VerifiedCustomerModel{}); err != nil {
		return nil, err
	}

	return &PostgresStore{DB: db, sqlDB: sqlDB}, nil
}

func (p *PostgresStore) Close() error {
	if p == nil || p.sqlDB == nil {
		return nil
	}
	return p.sqlDB.Close()
}
