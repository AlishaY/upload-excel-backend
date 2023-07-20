package database

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "sqlserver://coadmin:alisha@1234@127.0.0.1:1433?database=TODO&connection+timeout=30"

	var err error
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// Perform any additional configuration or setup for the database connection if needed
}
