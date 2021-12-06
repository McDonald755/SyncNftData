package db

import (
	"time"
)

type NFT_DATA struct {
	ID          int64     `gorm:"column:ID" json:"ID"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
	TokenId     int64     `gorm:"column:token_id" json:"token_id"`
	TokenUri    string    `gorm:"column:token_uri" json:"token_uri"`
	//TokenSymbol string    `gorm:"column:token_symbol" json:"token_symbol"`
	//TokenName   string    `gorm:"column:token_name" json:"token_name"`
	Owner     string `gorm:"column:owner" json:"owner"`
	OracleAdd string `gorm:"column:oracle_add" json:"oracle_add"`
}

func (NFT_DATA) TableName() string {
	return "NFT_DATA"
}

type ORACLE_DATA struct {
	ID          int64     `gorm:"column:ID" json:"ID"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
	Address     string    `gorm:"column:address" json:"address"`
	TokenSymbol string    `gorm:"column:token_symbol" json:"token_symbol"`
	TokenName   string    `gorm:"column:token_name" json:"token_name"`
}

func (ORACLE_DATA) TableName() string {
	return "ORACLE_DATA"
}
