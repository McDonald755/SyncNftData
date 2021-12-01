package db

import "SyncNftData/config"

func GetToAccount(from int, limit int) *[]string {
	//todo 看看能不能一次获取到所有数据，如果太大的话就得分页获取
	var result []string
	row := config.DB.Select("TO_ACCOUNT").Where("TX_DATA <> ?", "0x").Row()
	row.Scan(&result)
	return &result
}

func SaveOracle(addr []string) {
	// todo save数据
}

func GetOracleAddr() *[]string {
	// TODO 获取Oracle的地址 看看能不能一次性获取到
	return nil
}
