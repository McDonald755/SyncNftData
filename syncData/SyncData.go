package syncData

import (
	"SyncNftData/db"
	"SyncNftData/oracle"
	"SyncNftData/utils"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strconv"
	"strings"
)

//todo 多线程分开跑
func GetOracleAddr(client *ethclient.Client, from int) {
	number, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Error("Get block number error:", err)
	}
	//i就是从有721的那个区快开始
	for i := from; i < int(number); i++ {
		accounts := db.GetToAccount(i, 1000)
		var addrs []string
		for _, addr := range *accounts {
			code, err := client.CodeAt(context.Background(), common.HexToAddress(addr), nil)
			if err != nil {
				fmt.Println(err)
			}

			ok := utils.CheckOracleType(code)
			if ok {
				addrs = append(addrs, addr)
			}
		}
		db.SaveOracle(addrs)
	}
	from = int(number)
}

func GetNftData(client *ethclient.Client, from *big.Int) {
	number, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Error("Get block number error:", err)
	}
	addrs := db.GetOracleAddr()
	contractABI, _ := abi.JSON(strings.NewReader(oracle.OracleABI))
	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{
			{contractABI.Events["Transfer"].ID},
		},
		FromBlock: from,
		ToBlock:   big.NewInt(int64(number)),
		Addresses: *utils.TransferAddr(addrs),
	}

	filterLogs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		fmt.Println(err)
	}
	for _, l := range filterLogs {
		fmt.Println(l.Address)
		fmt.Println("tx", l.TxHash.Hex())
		fmt.Println("from", common.HexToAddress(l.Topics[1].Hex()).String())
		fmt.Println("to", common.HexToAddress(l.Topics[2].Hex()).String())
		parseInt, err := strconv.ParseInt(l.Topics[3].Hex(), 0, 16)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("tokenId", parseInt)
	}
	from = big.NewInt(int64(number))
}
