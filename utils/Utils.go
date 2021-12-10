package utils

import (
	"SyncNftData/config"
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
	"github.com/holiman/uint256"
	log "github.com/sirupsen/logrus"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func CheckOracleType(client *ethclient.Client, trxs types.Transactions, oracles map[string]byte, newOracles map[string]byte) int {
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
		// if to ==nil means create contract
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

			//if b1 && b2 && b3 && b4 && b5 && b6 && b7 && b8 && b9 is true means the contract is 721-contract
			if b1 && b2 && b3 && b4 && b5 && b6 && b7 && b8 && b9 {
				receipt, err := client.TransactionReceipt(context.Background(), trx.Hash())
				if err != nil {
					log.Error("TransactionReceipt err:", err)
				}

				//filter new oracle data
				if _, ok := oracles[receipt.ContractAddress.String()]; !ok {
					s := receipt.ContractAddress.String()
					newOracles[s] = byte(1)
					//save data
					go transferOracle(client, receipt.ContractAddress.String())
				}
			}
		}
	}
	return len(newOracles)
}

func OraclesToMap(oracles []string) map[string]byte {
	var result map[string]byte
	result = make(map[string]byte)
	for _, s := range oracles {
		result[s] = byte(1)
	}
	return result
}

func transferOracle(client *ethclient.Client, addres string) {
	symbol, name := getTokenNameAndSymbol(client, addres)
	data := db.ORACLE_DATA{
		Address:     addres,
		TokenSymbol: symbol,
		TokenName:   name,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	db.SaveOracle(&data)
}

func getTokenNameAndSymbol(client *ethclient.Client, addr string) (string, string) {
	var s, n string
	newOracle, err := oracle.NewOracle(common.HexToAddress(addr), client)
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
	return s, n
}

func getTokenSymbol(client *ethclient.Client, addr string) string {
	var s string
	newOracle, err := oracle.NewOracle(common.HexToAddress(addr), client)
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
	return s
}

func getTokenUrI(client *ethclient.Client, addr string, tokenId *big.Int) string {
	i := "Undefined"
	newOracle, err := oracle.NewOracle(common.HexToAddress(addr), client)
	if err != nil {
		log.Error("Init Oracle Error:", err, "Oracle Addr Is :", addr)
	}

	if tokenId != nil {
		uri, err := newOracle.TokenURI(nil, tokenId)
		if err != nil && err.Error() == "abi: attempting to unmarshall an empty string while arguments are expected" && err.Error() == "execution reverted" {
			log.Error("Get Token Uri Error:", err)
		} else {
			i = uri
		}
	}
	return i
}

func ScanLog(client *ethclient.Client, contractABI abi.ABI, addres map[string]byte, from int64) error {
	accounts := TransferAccounts(addres)
	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{
			{contractABI.Events["Approval"].ID,
				contractABI.Events["Transfer"].ID,
				contractABI.Events["ApprovalForAll"].ID},
		},
		FromBlock: big.NewInt(from),
		ToBlock:   big.NewInt(from + 1),
		Addresses: *accounts,
	}

	filterLogs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}
	go loopFilterLog(client, filterLogs)

	return nil
}

func loopFilterLog(client *ethclient.Client, filterLogs []types.Log) {
	for _, l := range filterLogs {
		dealLogMessage(client, l)
	}
}

func TScanLog(client *ethclient.Client, contractABI abi.ABI, addres []string, from int64, gap int64) error {
	accounts := TTransferAccounts(addres)
	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{
			{contractABI.Events["Approval"].ID,
				contractABI.Events["Transfer"].ID,
				contractABI.Events["ApprovalForAll"].ID},
		},
		FromBlock: big.NewInt(from),
		ToBlock:   big.NewInt(from + gap),
		Addresses: *accounts,
	}

	filterLogs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}

	for _, l := range filterLogs {
		dealLogMessage(client, l)
	}
	return nil
}

func dealLogMessage(client *ethclient.Client, l types.Log) {
	switch l.Topics[0].String() {
	case "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925":
		//Approval
		i := uint256.NewInt(0)
		if len(l.Topics) == 4 {
			//tokenID= topic【3】
			u := parsingUint256(l.Topics[3].Hex(), l.TxHash.String())
			if u != nil {
				i = u
			}
		} else {
			//tokenID = l.data
			toString := hex.EncodeToString(l.Data)
			u := parsingUint256(toString, l.TxHash.String())
			if u != nil {
				i = u
			}
		}
		data := db.NFT_DATA{
			TokenId:       i.ToBig().String(),
			OracleAdd:     strings.ToLower(l.Address.String()),
			TokenApproval: strings.ToLower(common.HexToAddress(l.Topics[2].String()).String()),
		}
		db.UpdateNftApproval(&data)
	case "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef":
		//Transfer
		data := TransferNftData(client, l)
		//fmt.Println("保存数据", from)
		db.SaveOrUpdateNftData(data)
	case "0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31":
		//ApprovalForAll
		i := uint256.NewInt(0)
		if len(l.Topics) == 4 {
			//tokenID= topic【3】
			u := parsingUint256(l.Topics[3].Hex(), l.TxHash.String())
			if u != nil {
				i = u
			}
		} else {
			//tokenID = l.data
			toString := hex.EncodeToString(l.Data)
			u := parsingUint256(toString, l.TxHash.String())
			if u != nil {
				i = u
			}
		}
		data := db.ORACLE_DATA{
			Address:     strings.ToLower(l.Address.String()),
			ApprovalAll: strings.ToLower(common.HexToAddress(l.Topics[2].String()).String()),
		}
		db.UpdateOracleApprove(&data, i.ToBig().String())
	}
}

//func T2ScanLog(client *ethclient.Client, contractABI abi.ABI, addres []string, from int64, to int64) error {
//	gap := big.NewInt(100)
//	for from < to {
//		s := big.NewInt(from)
//		e := s.Add(s, gap)
//
//		accounts := TTransferAccounts(addres)
//		query := ethereum.FilterQuery{
//			Topics: [][]common.Hash{
//				{contractABI.Events["Transfer"].ID},
//			},
//			FromBlock: s,
//			ToBlock:   e,
//			Addresses: *accounts,
//		}
//
//		filterLogs, err := client.FilterLogs(context.Background(), query)
//		if err != nil {
//			switch err.Error() {
//			//rate limit
//			case "too many requests":
//				time.Sleep(time.Second * 10)
//			default:
//				log.Error("Get log error:", err)
//			}
//		}
//
//		if err != nil && err.Error() == "" {
//			s, e, gap = CalculateBlock(gap, -1, gap)
//			continue
//		}
//		s, e, gap = CalculateBlock(gap, len(filterLogs), gap)
//
//		for _, l := range filterLogs {
//			data := TransferNftData(client, l)
//			fmt.Println("保存数据", from)
//			log.Infoln("保存数据", from)
//			db.SaveOrUpdateNftData(data)
//		}
//		from += 1
//	}
//	return nil
//}

func TransferNftData(client *ethclient.Client, l types.Log) *db.NFT_DATA {
	var (
		i   = uint256.NewInt(0)
		uri = "Undefined"
	)

	if len(l.Topics) == 4 {
		//tokenID= topic【3】
		u := parsingUint256(l.Topics[3].Hex(), l.TxHash.String())
		if u != nil {
			i = u
		}
	} else {
		//tokenID = l.data
		toString := hex.EncodeToString(l.Data)
		u := parsingUint256(toString, l.TxHash.String())
		if u != nil {
			i = u
		}
	}

	uri = getTokenUrI(client, l.Address.String(), i.ToBig())
	data := db.NFT_DATA{
		TokenId:   i.ToBig().String(),
		TokenUri:  uri,
		Owner:     strings.ToLower(common.HexToAddress(l.Topics[2].Hex()).String()),
		OracleAdd: strings.ToLower(l.Address.String()),
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
func CrawlData(from int64, page int64) {
	wg := sync.WaitGroup{}
	for i := from; i <= page; i++ {
		url := "https://bscscan.com/tokens-nft?ps=100&p=" + strconv.Itoa(int(i))
		get, err := http.Get(url)
		if err != nil {
			log.Error(err)
			i -= 1
			continue
		}
		wg.Add(1)
		go r(get.Body, config.CLIENTS[i%13], &wg, i)
	}
	wg.Wait()
}

//func CalculateBlock(from *big.Int, len int, gap *big.Int) (*big.Int, *big.Int, *big.Int) {
//	var startBlock, endBlock, newGap *big.Int
//	if len == -1 {
//		//logs > 5000 end = start+(gap/2)
//		startBlock = from
//		newGap = newGap.Quo(gap, big.NewInt(2))
//		endBlock = endBlock.Add(startBlock, newGap)
//	} else if len < 2500 {
//		// logs < 2500 start=from+gap end=start +newGap
//		startBlock = startBlock.Add(from, gap)
//		newGap = newGap.Mul(gap, big.NewInt(2))
//		endBlock = endBlock.Add(startBlock, newGap)
//	}
//	return startBlock, endBlock, newGap
//}

//func filterString(s string) string {
//	compile := regexp.MustCompile("^[0]+")
//	findString := compile.FindString(s[2:])
//	if findString == "0000000000000000000000000000000000000000000000000000000000000000" {
//		return "0x0"
//	}
//	return "0x" + s[len(findString)+2:]
//}

func parsingUint256(s string, hash string) *uint256.Int {
	var result string
	if s == "0000000000000000000000000000000000000000000000000000000000000000" {
		result = "0x0"
	} else {
		compile := regexp.MustCompile("^[0]+")
		findString := compile.FindString(s[2:])
		if findString == "0000000000000000000000000000000000000000000000000000000000000000" {
			result = "0x0"
		} else {
			result = "0x" + s[len(findString)+2:]
		}
	}
	fromHex, err := uint256.FromHex(result)
	if err != nil {
		log.Error("Hash is:", hash, "value is:", s, "\n", err)
		return nil
	}
	return fromHex
}

func r(r io.Reader, client *ethclient.Client, wg *sync.WaitGroup, i int64) {
	datas := []db.ORACLE_DATA{}
	reader, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Error(err)
	}
	reader.Find(".text-primary").Each(func(i int, s *goquery.Selection) {
		//var symbol string
		attr, _ := s.Attr("href")
		if attr != "https://etherscan.io/" && s.Text() != "Etherscan" {
			attr = attr[7:]
			//symbol = getTokenSymbol(client, attr)
			data := db.ORACLE_DATA{
				CreatedTime: time.Now(),
				UpdatedTime: time.Now(),
				Address:     strings.ToLower(attr),
				TokenName:   s.Text(),
				TokenSymbol: "Undefined",
			}
			datas = append(datas, data)
		}
	})
	db.SaveOracles(&datas)
	fmt.Println(i)
	wg.Done()
}
