package db

import (
	"SyncEthData/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SaveData(block *BLOCK, header *HEADER, trx *[]TRANSACTION) {
	dbErr := gorm.DB{}
	tx := config.DB.Begin()
	dbErr = *tx.Create(block)
	dbErr = *tx.Create(header)
	if len(*trx) > 0 {
		dbErr = *tx.Create(trx)
	}

	if dbErr.Error != nil {
		log.Error(dbErr.Error)
		log.Error("----------------------Error Num is:", block.BLOCKNUM)
		tx.Rollback()
	} else {
		//fmt.Println("save:", block.BLOCKNUM)
		tx.Commit()
	}
}
