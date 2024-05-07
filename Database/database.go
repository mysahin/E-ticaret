package database

import (
	"ETicaret/Models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func Connect() {
	port := "5432"
	dsn := fmt.Sprintf("host=db user=postgres password=password dbname=db port=%s sslmode=disable TimeZone=Asia/Shanghai", port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected success")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")
	if err := db.AutoMigrate(Models.Login{}, Models.User{}, Models.Product{}, Models.Category{}, Models.Type{}); err != nil {
		log.Fatal(err)
	}

	DB = Dbinstance{
		Db: db,
	}

}
