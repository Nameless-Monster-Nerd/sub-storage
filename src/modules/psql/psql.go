package psql

import (
	"log"
	"fmt"
	utils "github.com/nameless-Monster-Nerd/subtitle/src/modules"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Sub struct {
	ID     string  `gorm:"column:id"`
	SS     *string `gorm:"column:ss"`
	EP     *string `gorm:"column:ep"`
	Key    string  `gorm:"primaryKey;column:key"`
	Bucket string  `gorm:"column:bucket"`
	Lang   string  `gorm:"column:lang"`
	Flix   bool    `gorm:"column:flix"`
}

// BatchUpload inserts subtitles in bulk, skipping any duplicate primary keys.
func BatchUpload(subs []Sub) {


	// AutoMigrate ensures the table exists
	if err := utils.Db.AutoMigrate(&Sub{}); err != nil {
		fmt.Println(err)
	}

	if len(subs) == 0 {
		return
	}

	// Insert and skip duplicates (ON CONFLICT DO NOTHING)
	result := utils.Db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&subs)

	if result.Error != nil {
		log.Printf("BatchUpload error: %v\n", result.Error)
	}
}

// BatchSearch looks up subtitles based on ID, season, episode, and flix flag.
func BatchSearch(id string, ss *string, ep *string, flix bool) ([]Sub, error) {

	var subs []Sub
	query := utils.Db.Model(&Sub{}).Where("id = ? AND flix = ?", id, flix)

	if ss != nil {
		query = query.Where("ss = ?", *ss)
	}
	if ep != nil {
		query = query.Where("ep = ?", *ep)
	}

	if err := query.Find(&subs).Error; err != nil {
		return nil, err
	} 
	return subs, nil
}
