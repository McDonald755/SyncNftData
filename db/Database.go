package db

import (
	"SyncNftData/config"
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

func GetOracleAddrAll() []string {
	var (
		addres []string
	)
	config.DB.Table("ORACLE_DATA").Select("address").Limit(9000).Find(&addres)
	return addres
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
	config.DB.Table("ORACLE_DATA").Select("approval_all", "updated_time").Where("ID = ?", data.ID).Updates(data)
}

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
