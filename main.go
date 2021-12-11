package main

import (
	"SyncNftData/config"
	"SyncNftData/db"
	"SyncNftData/log"
	"SyncNftData/oracle"
	"SyncNftData/syncData"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"sync"
	"time"
)

func main() {
	log.ConfigLocalFilesystemLogger("./errorLog", "log", time.Hour*24*14, time.Hour*24)
	//cmd.Execute()
	startNum := 1
	syncNum := 5554071
	wg := sync.WaitGroup{}

	//init oracle data
	Oracles := db.TGetOracleAddrAll()
	//init standard ERC-721 contract data
	contractABI, _ := abi.JSON(strings.NewReader(oracle.OracleABI))

	//init data
	oracleNum := len(Oracles)
	clientLen := len(config.CLIENTS)
	num := oracleNum / clientLen
	number, _ := config.CLIENTS[0].BlockNumber(context.Background())

	wg.Add(1)
	go syncData.TSyncData(config.CLIENTS[0], int64(startNum), Oracles[:num], contractABI, &wg, int64(number), 1)

	wg.Add(1)
	go syncData.SyncData(config.CLIENTS[clientLen-1], int64(syncNum), Oracles, contractABI, &wg)
	wg.Wait()
}
