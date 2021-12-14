package syncData

import (
	"SyncNftData/utils"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"math/big"
	"sync"
	"time"
)

func SyncDataMain(client *ethclient.Client, from int64, oracles []string, contractABI abi.ABI, wg *sync.WaitGroup) {
	oraclesToMap := utils.OraclesToMap(oracles)
	var newOracles map[string]byte
	newOracles = make(map[string]byte)
	for true {
		//get block data
		blockNum, err := client.BlockByNumber(context.Background(), big.NewInt(from))
		if err != nil && err.Error() != "not found" {
			log.Error("BlockByNumber:", err)
			continue
		}
		if blockNum == nil && err.Error() == "not found" {
			time.Sleep(time.Second)
			continue
		}

		//Analyze the transaction
		newLen := utils.CheckOracleType(client, blockNum.Transactions(), oraclesToMap, newOracles)
		if newLen <= 0 {
			from += 1
			continue
		}

		err = utils.ScanLog(client, contractABI, newOracles, from)
		if err != nil && err.Error() == "too many requests" {
			time.Sleep(time.Second * 10)
			continue
		} else if err != nil && err.Error() != "too many requests" {
			log.Error("Main func get log error", err, "block num is :", from)
		}

		oraclesToMap = newOracles
		fmt.Println(from)
		from += 1
	}
	wg.Done()
}

//init data
func InitData(client *ethclient.Client, from int64, oracles []string, contractABI abi.ABI, wg *sync.WaitGroup, endNum int64, gap int64, i int) {
	for endNum < from {
		err := utils.ScanLogByInitData(client, contractABI, oracles, from, gap)
		//rate limit
		if err != nil && err.Error() == "too many requests" {
			time.Sleep(time.Second * 10)
			continue
		} else if err != nil && err.Error() != "too many requests" {
			log.Error("Get log error :", i, ":", err)
			time.Sleep(time.Second * 10)
			continue
		}
		fmt.Println(i, ":", from)
		from -= gap
	}
	log.Infoln(i, ":End")
	wg.Done()
}

func SyncData(client *ethclient.Client, from int64, oracles []string, contractABI abi.ABI, wg *sync.WaitGroup, gap int64, i int) {
	for true {
		number, err := client.BlockNumber(context.Background())
		if err != nil {
			log.Error("Get Block Num Error:", err)
		}

		if from < int64(number) {
			err := utils.ScanLogByInitData(client, contractABI, oracles, from, gap)
			//rate limit
			if err != nil && err.Error() == "too many requests" {
				time.Sleep(time.Second * 10)
				continue
			} else if err != nil && err.Error() != "too many requests" {
				log.Error("Get log error :", i, ":", err)
				continue
			}
			fmt.Println(i, ":", from)
			from += gap
		}
	}
	wg.Done()
}
