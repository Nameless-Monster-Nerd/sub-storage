package psql

import (
	utils "github.com/nameless-Monster-Nerd/subtitle/src/modules"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Sub struct {
	ID     string  `gorm:"column:id"`
	SS     *string `gorm:"column:ss"`
	EP     *string `gorm:"column:ep"`
	Key    string  `gorm:"primaryKey;column:key"`
	Bucket string  `gorm:"column:bucket"`
	Lang 	string `gorm:"column:lang"`
	Flix   bool    `gorm:"column:flix"`
	
}

func BatchUpload(subs []Sub) {
	db, err := gorm.Open(postgres.Open(utils.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Sub{})
	if err != nil {
		panic(err)
	}

	if len(subs) == 0 {
		return
	}

	// Bulk insert
	result := db.Create(&subs)
	if result.Error != nil {
		panic(result.Error)
	}
}

func BatchSearch(id string, ss *string, ep *string, flix bool) ([]Sub, error) {
	db, err := gorm.Open(postgres.Open(utils.Dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var subs []Sub
	query := db.Model(&Sub{}).Where("id = ? AND flix = ?", id, flix)

	if ss != nil {
		query = query.Where("ss = ?", *ss)
	}
	if ep != nil {
		query = query.Where("ep = ?", *ep)
	}

	if err := query.Find(&subs).Error; err != nil {
		return nil, err
	}

	if len(subs) == 0 {
		return nil, nil
	}

	return subs, nil
}



