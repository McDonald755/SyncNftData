package utils

import (
	"SyncNftData/db"
	"SyncNftData/oracle"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CheckOracleType(client *ethclient.Client, trxs types.Transactions, oracles map[string]byte) map[string]byte {
	var (
		balanceOf         = "70a08231"
		ownerOf           = "6352211e"
		approve           = "095ea7b3"
		getApproved       = "081812fc"
		setApprovalForAll = "a22cb465"
		isApprovedForAll  = "e985e9c5"
		transferFrom      = "23b872dd"
		safeTransferFrom  = "42842e0e"
		safeTransferFrom2 = "b88d4fde"
	)
	for _, trx := range trxs {
		if trx.To() == nil {
			//encode tx_data to string
			txData := hex.EncodeToString(trx.Data())
			b1 := strings.Contains(txData, balanceOf)
			b2 := strings.Contains(txData, ownerOf)
			b3 := strings.Contains(txData, approve)
			b4 := strings.Contains(txData, getApproved)
			b5 := strings.Contains(txData, setApprovalForAll)
			b6 := strings.Contains(txData, isApprovedForAll)
			b7 := strings.Contains(txData, transferFrom)
			b8 := strings.Contains(txData, safeTransferFrom)
			b9 := strings.Contains(txData, safeTransferFrom2)
			if b1 && b2 && b3 && b4 && b5 && b6 && b7 && b8 && b9 {
				receipt, err := client.TransactionReceipt(context.Background(), trx.Hash())
				if err != nil {
					log.Error("TransactionReceipt err:", err)
				}
				if _, ok := oracles[receipt.ContractAddress.String()]; !ok {
					oracles[receipt.ContractAddress.String()] = byte(0)

					//save data
					data := transferOracle(client, receipt.ContractAddress.String())
					db.SaveOracle(data)
				}
			}
		}
	}
	return oracles
}

func transferOracle(client *ethclient.Client, addres string) *db.ORACLE_DATA {
	symbol, name, _ := getTokenNameAndSymbol(client, addres, nil)
	data := db.ORACLE_DATA{
		Address:     addres,
		TokenSymbol: symbol,
		TokenName:   name,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	return &data
}

func getTokenNameAndSymbol(client *ethclient.Client, addr string, tokenId *big.Int) (string, string, string) {
	var s, n, i string
	address := common.HexToAddress(addr)
	newOracle, err := oracle.NewOracle(address, client)
	if err != nil {
		log.Error("Init Oracle Error:", err, "Oracle Addr Is :", addr)
	}
	symbol, err := newOracle.Symbol(nil)
	if err != nil {
		log.Error("Get Token Symbol Error:", err)
		s = "Undefined"
	} else {
		s = symbol
	}

	name, err := newOracle.Name(nil)
	if err != nil {
		log.Error("Get Token Name Error:", err)
		n = "Undefined"
	} else {
		n = name
	}

	if tokenId != nil {
		uri, err := newOracle.TokenURI(nil, tokenId)
		if err != nil {
			log.Error("Get Token Uri Error:", err)
			i = "Undefined"
		} else {
			i = uri
		}
	}
	return s, n, i
}

func ScanLog(client *ethclient.Client, contractABI abi.ABI, addres map[string]byte, from int64) {

	accounts := TransferAccounts(addres)
	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{
			{contractABI.Events["Transfer"].ID},
		},
		FromBlock: big.NewInt(from),
		ToBlock:   big.NewInt(from),
		Addresses: *accounts,
	}

	filterLogs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Error("Get log error:", err)
	}

	for _, l := range filterLogs {
		data := TransferNftData(l)
		fmt.Println("保存数据", from)
		log.Info("保存数据", from)
		db.SaveOrUpdateNftData(data)
	}
}

func TScanLog(client *ethclient.Client, contractABI abi.ABI, addres []string, from int64) {

	accounts := TTransferAccounts(addres)
	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{
			{contractABI.Events["Transfer"].ID},
		},
		FromBlock: big.NewInt(from),
		ToBlock:   big.NewInt(from),
		Addresses: *accounts,
	}

	filterLogs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Error("Get log error:", err)
	}

	for _, l := range filterLogs {
		data := TransferNftData(l)
		fmt.Println("保存数据", from)
		log.Info("保存数据", from)
		db.SaveOrUpdateNftData(data)
	}
}

func TransferNftData(l types.Log) *db.NFT_DATA {
	parseInt, err := strconv.ParseInt(l.Topics[3].Hex(), 0, 16)
	if err != nil {
		log.Error(err.Error())
	}

	symbol, name, uri := getTokenNameAndSymbol(nil, l.Address.String(), big.NewInt(parseInt))
	data := db.NFT_DATA{
		TokenId:     parseInt,
		TokenSymbol: symbol,
		TokenName:   name,
		TokenUri:    uri,
		Owner:       common.HexToAddress(l.Topics[2].Hex()).String(),
		OracleAdd:   l.Address.String(),
	}
	return &data
}

func TransferAccounts(addres map[string]byte) *[]common.Address {
	result := []common.Address{}
	for k, _ := range addres {
		result = append(result, common.HexToAddress(k))
	}
	return &result
}

func TTransferAccounts(addres []string) *[]common.Address {
	result := []common.Address{}
	for _, v := range addres {
		result = append(result, common.HexToAddress(v))
	}
	return &result
}

//get 721Token message by bsc
func CrawlData(page int) {
	for i := 1; i <= page; i++ {
		datas := []db.ORACLE_DATA{}
		url := "https://bscscan.com/tokens-nft?ps=100&p=" + strconv.Itoa(i)
		get, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		reader, err := goquery.NewDocumentFromReader(get.Body)
		if err != nil {
			fmt.Println(err)
		}
		reader.Find(".text-primary").Each(func(i int, s *goquery.Selection) {
			attr, _ := s.Attr("title")
			if s.Text() != "Etherscan" {
				data := db.ORACLE_DATA{
					CreatedTime: time.Now(),
					UpdatedTime: time.Now(),
					Address:     attr,
					TokenName:   s.Text(),
				}
				datas = append(datas, data)
			}
		})
		db.SaveOracles(&datas)
		time.Sleep(time.Second)
	}
}

/**
===================================================================Methods not used yet, don't remove===================================================================
*/

/*func CalculateBlock(from *big.Int, len int, gap *big.Int) (*big.Int, *big.Int, *big.Int) {
	var startBlock, endBlock, newGap *big.Int
	if len == -1 {
		//logs > 10000 end = start+(gap/2)
		startBlock = from
		newGap = newGap.Quo(gap, big.NewInt(2))
		endBlock = endBlock.Add(startBlock, newGap)
	} else if len < 5000 {
		// logs < 5000 start=from+gap end=start +newGap
		startBlock = startBlock.Add(from, gap)
		newGap = newGap.Mul(gap, big.NewInt(2))
		endBlock = endBlock.Add(startBlock, newGap)
	}
	return startBlock, endBlock, newGap
}
*/
