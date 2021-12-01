package db

type NFT_DATA struct {
	ID          string `orm:"ID" json:"ID"`
	CREATEDTIME string `orm:"CREATED_TIME" json:"CREATED_TIME"`
	UPDATEDTIME string `orm:"UPDATED_TIME" json:"UPDATED_TIME"`
	TOKENID     string `orm:"TOKEN_ID" json:"TOKEN_ID"`
	TOKENURI    string `orm:"TOKEN_URI" json:"TOKEN_URI"`
	TOKENSYMBOL string `orm:"TOKEN_SYMBOL" json:"TOKEN_SYMBOL"`
	TOKENNAME   string `orm:"TOKEN_NAME" json:"TOKEN_NAME"`
	OWNER       string `orm:"OWNER" json:"OWNER"`
	ORACLEADD   string `orm:"ORACLE_ADD" json:"ORACLE_ADD"`
}

func (NFT_DATA) TableName() string {
	return "nft_data"
}
