package main

import (
	"SyncNftData/config"
	"SyncNftData/db"
	"SyncNftData/log"
	"SyncNftData/oracle"
	"SyncNftData/syncData"
	"SyncNftData/utils"
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
	Gap := 10
	start := 100
	a := 1000000 / 13

	y := len(Oracles) / 13
	for i := 0; i < 13; i++ {
		from := a*i + start
		wg.Add(1)
		go syncData.TSyncData(config.CLIENTS[i], int64(from), Oracles[i*y:(i+1)*y], contractABI, &wg)
	}
	wg.Wait()
}
