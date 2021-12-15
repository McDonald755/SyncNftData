package db

import (
	"SyncNftData/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func SaveOracle(oracle *ORACLE_DATA) {
	config.DB.Create(oracle)
}

func SaveOracles(oracle *[]ORACLE_DATA) {
	config.DB.Create(oracle)
}

func SaveOrUpdateNftData(nft *NFT_DATA) {
	var data NFT_DATA
	result := config.DB.Table("NFT_DATA").Where("oracle_addr = ? and token_id = ?", nft.OracleAddr, nft.TokenId).Find(&data)
	if result.RowsAffected == 0 {
		nft.CreatedTime = time.Now()
		nft.UpdatedTime = time.Now()
		config.DB.Create(nft)
	} else {
		data.TokenId = nft.TokenId
		data.TokenUri = nft.TokenUri
		data.Owner = nft.Owner
		data.OracleAddr = nft.OracleAddr
		data.UpdatedTime = time.Now()
		save := config.DB.Save(data)

		//uri maybe error
		if save.Error != nil {
			fmt.Println("“11111111111111111”")
			data.TokenId = nft.TokenId
			data.TokenUri = "Undefined"
			data.Owner = nft.Owner
			data.OracleAddr = nft.OracleAddr
			data.UpdatedTime = time.Now()
			config.DB.Save(data)
		}
	}
}

func GetOracleAddrAll() (map[string]byte, int) {
	var (
		addres []string
		result map[string]byte
	)
	result = make(map[string]byte)
	config.DB.Table("ORACLE_DATA").Select("address").Find(&addres)
	for _, addre := range addres {
		result[addre] = byte(1)
	}
	return result, len(addres)
}

func TGetOracleAddrAll() []string {
	var (
		addres []string
	)
	config.DB.Table("ORACLE_DATA").Select("address").Find(&addres)
	return addres
}

func UpdateNftApproval(nft *NFT_DATA) {
	var data NFT_DATA
	s := config.DB.Table("NFT_DATA").Where("oracle_addr = ? and token_id = ?", nft.OracleAddr, nft.TokenId).Find(&data)
	if s.Error != nil {
		log.Error("UpdateNftApproval error", s.Error)
	}
	if s.RowsAffected == 0 {
		nft.CreatedTime = time.Now()
		nft.UpdatedTime = time.Now()
		//Splicing results
		nft.TokenApproval = nft.TokenApproval + ","
		config.DB.Create(nft)
	} else {
		//Splicing results
		data.TokenApproval = nft.TokenApproval + "," + data.TokenApproval
		data.UpdatedTime = time.Now()
		config.DB.Save(data)
	}
}

func UpdateOracleApprove(oracle *ORACLE_DATA, s string) {
	data := ORACLE_DATA{}
	find := config.DB.Table("ORACLE_DATA").Where("address = ?", oracle.Address).Find(&data)
	if find.Error != nil {
		log.Error("UpdateOracleApprove error", find.Error)
	}
	if s == "0" {
		//cancel approval
		toArray := stringToArray(data.ApprovalAll)
		for _, account := range toArray {
			if account != oracle.ApprovalAll {
				toArray = append(toArray, account)
			}
		}
		data.ApprovalAll = removeSameValue(toArray)

	} else if s == "1" {
		//set approval
		array := stringToArray(data.ApprovalAll)
		array = append(array, oracle.ApprovalAll)
		data.ApprovalAll = removeSameValue(array)
	}
	data.UpdatedTime = time.Now()
	save := config.DB.Save(data)
	if save.Error != nil {
		log.Error(save.Error)
	}
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
func removeSameValue(s []string) string {
	if len(s) != 0 {
		var (
			m      map[string]byte
			result string
		)
		m = make(map[string]byte)

		for _, s2 := range s {
			m[s2] = byte(1)
		}

		for k, _ := range m {
			if k != "" {
				result = k + "," + result
			}
		}
		return result
	}
	return ""
}

func stringToArray(s string) []string {
	if s != "" {
		return strings.Split(s, ",")
	}
	return nil
}
