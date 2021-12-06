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

func SyncData(client *ethclient.Client, from int64, oracles map[string]byte, contractABI abi.ABI, wg *sync.WaitGroup) {
	for true {
		//get block data
		blockNum, err := client.BlockByNumber(context.Background(), big.NewInt(from))
		if err != nil {
			log.Error("BlockByNumber:", err)
		}

		if blockNum == nil && err.Error() == "not found" {
			continue
		}

		//Analyze the transaction
		oracleType := utils.CheckOracleType(client, blockNum.Transactions(), oracles)

		//filter and save nft data
		utils.ScanLog(client, contractABI, oracleType, from)
		fmt.Println(from)
		from += 1
	}
	wg.Done()
}

func TSyncData(client *ethclient.Client, from int64, oracles []string, contractABI abi.ABI, wg *sync.WaitGroup) {
	for true {
		time.Sleep(time.Millisecond * 300)
		//get block data
		blockNum, err := client.BlockByNumber(context.Background(), big.NewInt(from))
		if err != nil {
			log.Error("BlockByNumber:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		if blockNum == nil && err.Error() == "not found" {
			continue
		}

		//Analyze the transaction
		//oracleType := utils.TCheckOracleType(client,blockNum.Transactions(), oracles)

		//filter and save nft data
		err = utils.TScanLog(client, contractABI, oracles, from)
		if err != nil && err.Error() == "too many requests" {
			continue
		}
		fmt.Println(from)
		from += 1
	}
	wg.Done()
}

/**
===================================================================Methods not used yet, don't remove===================================================================
*/

/*func SyncOracleAddr() {
	wg := sync.WaitGroup{}
	totalNum := db.GetTrxTotalNum()
	distance := math.Ceil((float64(totalNum)) / 1.0)
	page := math.Ceil(float64(distance) / 5000000.0)

	for i := 0; i < 1; i++ {
		saveOracleData(int64(distance), int64(i), page)
	}
	wg.Wait()
}

func saveOracleData(distance int64, i int64, page float64) {
	startNum := distance * i
	for y := 0; y < int(page); {
		vos := db.GetToAccount(startNum)
		utils.CheckOracleType(vos)
		y += 1
		startNum = int64(y)*50000 + startNum
	}
}

func ScanOracleData(from int64) {
	for true {
		number, err := config.CLIENT[0].BlockNumber(context.Background())
		if err != nil {
			log.Error("Get Block Number Error:", err)
		}
		for from < int64(number) {
			vos := db.GetTrxByNum(from)
			utils.CheckOracleType(vos)
			syncNftData(from)
			from += 1
			time.Sleep(time.Second)
		}
	}
}

func SyncNftData() {
	gap := big.NewInt(6000)
	oracles := db.GetOracleNum()
	distance := int(oracles) / len(config.CLIENT)

	contractABI, err := abi.JSON(strings.NewReader(oracle.OracleABI))
	if err != nil {
		log.Error("Read Contract Error:", err)
	}

	for i, client := range config.CLIENT {
		addr := db.GetOracleAddr(i, i+distance)
		go scanLog(client, contractABI, addr, gap)
	}
}

func scanLog(client *ethclient.Client, contractABI abi.ABI, addr *[]string, gap *big.Int) {
	var (
		from, to *big.Int
		newGap   = gap
	)

	accounts := utils.TransferAccounts(addr)
	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{
			{contractABI.Events["Transfer"].ID},
		},
		FromBlock: from,
		ToBlock:   to.Add(from, newGap),
		Addresses: *accounts,
	}

	filterLogs, err := client.FilterLogs(context.Background(), query)
	if err != nil && err.Error() != "" {
		log.Error("Get log error:", err)
	}
	if err.Error() == "" {
		from, to, newGap = utils.CalculateBlock(from, -1, gap)
	}
	from, to, newGap = utils.CalculateBlock(from, len(filterLogs), gap)

	for _, l := range filterLogs {
		data := utils.TransferNftData(l)
		db.SaveOrUpdateNftData(data)
	}
}

func syncNftData(from int64) {
	var to *big.Int
	all := db.GetOracleAddrAll()
	accounts := utils.TransferAccounts(all)
	contractABI, err := abi.JSON(strings.NewReader(oracle.OracleABI))

	if err != nil {
		log.Error("Read Contract Error:", err)
	}
	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{
			{contractABI.Events["Transfer"].ID},
		},
		FromBlock: big.NewInt(from),
		ToBlock:   to.Add(big.NewInt(from), big.NewInt(1)),
		Addresses: *accounts,
	}

	filterLogs, err := config.CLIENT[0].FilterLogs(context.Background(), query)
	if err != nil && err.Error() != "" {
		log.Error("Get log error:", err)
	}
	for _, l := range filterLogs {
		data := utils.TransferNftData(l)
		db.SaveOrUpdateNftData(data)
	}
}*/
