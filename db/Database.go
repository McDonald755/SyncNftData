package db

import (
	"SyncNftData/config"
)

func SaveOracle(oracle *ORACLE_DATA) {
	config.DB.Create(oracle)
}

func SaveOrUpdateNftData(nft *NFT_DATA) {
	var id int64
	result := config.DB.Table("NFT_DATA").Select("ID").Where("oracle_add,token_id", nft.OracleAdd, nft.TokenId).Find(&id)
	if result.RowsAffected == 0 {
		config.DB.Create(nft)
	} else {
		nft.ID = id
		config.DB.Save(nft)
	}
}

func GetOracleAddrAll() map[string]byte {
	var (
		addres []string
		result map[string]byte
	)
	config.DB.Table("ORACLE_DATA").Select("address").Find(&addres)
	for _, addre := range addres {
		result[addre] = byte(0)
	}
	return result
}

/**
===================================================================Methods not used yet, don't remove===================================================================
*/

/*func GetToAccount(from int64) *[]ResultVo {
	var r []ResultVo
	config.DB.Table("TRANSACTION").Select("id,to_account , tx_data, block_number").Where("tx_data <> '0x'").Limit(5000000).Offset(int(from)).Find(&r)
	return &r
}

func GetTrxByNum(block int64) *[]ResultVo {
	var r []ResultVo
	config.DB.Table("TRANSACTION").Select(" id,to_account , tx_data, block_number").Where("BLOCK_NUMBER", block).Find(&r)
	return &r
}
func SaveOracles(addr *[]ORACLE_DATA) {
	result := config.DB.Create(addr)
	if result.Error != nil {
		log.Error("Save Oracle Data Error:",  result.Error)
	}
func GetTrxTotalNum() int64 {
	var count int64
	config.DB.Table("TRANSACTION").Select("ID").Order("ID desc").Limit(1).Find(&count)
	return count
}

func GetOracleNum() int64 {
	var count int64
	config.DB.Table("ORACLE_DATA").Select("ID").Order("ID desc").Limit(1).Find(&count)
	return count
}

}*/
