package main

import (
	"SyncNftData/config"
	"SyncNftData/db"
	"SyncNftData/log"
	"SyncNftData/oracle"
	"SyncNftData/syncData"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"sync"
	"time"
)

func main() {
	log.ConfigLocalFilesystemLogger("./errorLog", "log", time.Hour*24*14, time.Hour*24)
	//cmd.Execute()

	wg := sync.WaitGroup{}

	//init oracle data
	Oracles := db.GetOracleAddrAll()
	//init standard ERC-721 contract data
	contractABI, _ := abi.JSON(strings.NewReader(oracle.OracleABI))
	startNum := 10000
	//init data
	oracleNum := len(Oracles)
	//clientLen := len(config.CLIENTS)
	num := oracleNum / 50
	gap := int64(100)
	//newNum, _ := config.CLIENTS[0].BlockNumber(context.Background())
	for i := 0; i < 50; i++ {
		time.Sleep(time.Millisecond * 300)
		wg.Add(1)
		a := i % 3
		go syncData.InitData(config.CLIENTS[a], int64(13464000), Oracles[i*num:(i+1)*num], contractABI, &wg, int64(startNum), gap, a)
	}
	wg.Wait()
}
