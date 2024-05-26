package helper

import (
	"gorm.io/gorm"
	"strconv"
)

func ResetAutoIncrement(db *gorm.DB, model interface{}, primaryColumn string, tableName string) error {
	var count int64
	db.Model(model).Count(&count)

	var query string
	if count == 0 {
		query = "ALTER TABLE " + tableName + " AUTO_INCREMENT = 1"
	} else {
		var LastId uint
		err := db.Raw("SELECT MAX(" + primaryColumn + ") FROM " + tableName).Scan(&LastId).Error
		if err != nil {
			return err
		}
		query = "ALTER TABLE " + tableName + " AUTO_INCREMENT = " + strconv.FormatUint(uint64(LastId+1), 10)
	}

	err := db.Exec(query).Error
	if err != nil {
		return err
	}

	return nil
}
