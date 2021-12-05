package main

import (
	"SyncNftData/config"
	"SyncNftData/db"
	"SyncNftData/log"
	"SyncNftData/oracle"
	"SyncNftData/syncData"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"sync"
	"time"
)

func main() {
	log.ConfigLocalFilesystemLogger("./errorLog", "log", time.Hour*24*14, time.Hour*24)
	//cmd.Execute()
	wg := sync.WaitGroup{}
	Oracles := db.TGetOracleAddrAll()

	//init standard ERC-721 contract data
	contractABI, err := abi.JSON(strings.NewReader(oracle.OracleABI))
	if err != nil {
		fmt.Println(err)
	}
	start := 1326
	a := 1000000 / 10

	y := len(Oracles) / 10
	for i := 0; i < 10; i++ {
		from := a*i + start
		wg.Add(1)
		go syncData.TSyncData(config.CLIENTS[0], int64(from), Oracles[y*i:y*(i+1)], contractABI, &wg)
	}
	wg.Wait()
}
