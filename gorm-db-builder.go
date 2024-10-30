package goapiconfigutilis

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDbBuilder struct {
	host     string
	username string
	password string
	dbName   string
	port     int
	DB       *gorm.DB
}

func InitGormDbBuilder(host, username, password, dbName string, port int) *GormDbBuilder {
	return &GormDbBuilder{
		host:     host,
		username: username,
		password: password,
		dbName:   dbName,
		port:     port,
	}
}

func (d *GormDbBuilder) builderDbLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
}

func (d *GormDbBuilder) BuildForPostgres() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		d.host,
		d.username,
		d.password,
		d.dbName,
		d.port)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:      d.builderDbLogger(),
		NowFunc:     time.Now,
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}

	d.DB = db
	return db, nil
}

func (d *GormDbBuilder) BuildForMySQL() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.username,
		d.password,
		d.host,
		d.port,
		d.dbName)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                     dsn,
		DefaultStringSize:       256,
		DontSupportRenameIndex:  true,
		DontSupportRenameColumn: true,
	}), &gorm.Config{
		Logger:      d.builderDbLogger(),
		NowFunc:     time.Now,
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}

	d.DB = db
	return db, nil
}

func (d *GormDbBuilder) DoMigration(enterties []interface{}) error {
	for _, dest := range enterties {
		if err := d.DB.AutoMigrate(&dest); err != nil {
			return err
		}
	}

	return nil
}
