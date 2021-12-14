package db

import "time"

type NFT_DATA struct {
	ID            string    `gorm:"column:ID" json:"ID"`
	CreatedTime   time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime   time.Time `gorm:"column:updated_time" json:"updated_time"`
	TokenId       string    `gorm:"column:token_id" json:"token_id"`
	TokenUri      string    `gorm:"column:token_uri" json:"token_uri"`
	Owner         string    `gorm:"column:owner" json:"owner"`
	OracleAddr    string    `gorm:"column:oracle_addr" json:"oracle_addr"`
	TokenApproval string    `gorm:"token_approval" json:"token_approval"`
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
	ApprovalAll string    `gorm:"column:approval_all" json:"approval_all"`
	Index       string    `gorm:"column:index" json:"index"`
}

func (ORACLE_DATA) TableName() string {
	return "ORACLE_DATA"
}
