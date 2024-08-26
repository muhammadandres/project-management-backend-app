package app

import (
	"fmt"
	"log"
	"manajemen_tugas_master/model/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "admin:andres12@tcp(manajemen-tugas-master.c5auaowcutch.ap-southeast-3.rds.amazonaws.com:3306)/manajementugasdb?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %v", err)
	}

	log.Println("Running migrations")
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Board{},
		&domain.Task{},
		&domain.Owner{},
		&domain.Manager{},
		&domain.Employee{},
		&domain.PlanningFile{},
		&domain.ProjectFile{},
		&domain.PlanningDescriptionFile{},
	); err != nil {
		return nil, err
	}

	return db, err
}
