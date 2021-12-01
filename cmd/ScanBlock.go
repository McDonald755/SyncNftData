package cmd

import (
	"SyncEthData/config"
	"SyncEthData/syncData"
	"SyncEthData/utils"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"math/big"
	"sync"
	"time"
)

func ScanCmd() *cobra.Command {
	var blockNum int
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "s",
		Long:  "It will sync the latest block ",
		RunE: func(cmd *cobra.Command, args []string) error {
			scanBlock(blockNum)
			return nil
		},
	}
	scanCmd.Flags().IntVarP(&blockNum, "blockNum", "b", 0, "input blockNum")
	return scanCmd
}

func scanBlock(blockNum int) {
	height := syncData.GetBlockHeight(config.CLIENT[0])
	distance := (height - blockNum) / len(config.CLIENT)
	wg := sync.WaitGroup{}
	for i := 0; i < len(config.CLIENT)-1; i++ {
		wg.Add(1)
		go getBlock(config.CLIENT[i], i, distance, blockNum, &wg)
	}
	wg.Add(1)
	go scanNewBlock(config.CLIENT[len(config.CLIENT)-1], (len(config.CLIENT)-1)*distance+blockNum, &wg)
	wg.Wait()
}

func getBlock(client *ethclient.Client, i int, distance int, blockNum int, wg *sync.WaitGroup) {
	from := distance*i + blockNum
	end := from + distance
	for from < end {
		block, err := syncData.GetBlockByNum(client, big.NewInt(int64(from)))
		if err != nil {
			time.Sleep(time.Hour)
		} else {
			utils.TransformData(block)
			from += 1
		}
	}
	wg.Done()
}

func scanNewBlock(client *ethclient.Client, from int, wg *sync.WaitGroup) {
	for true {
		block, err := syncData.GetBlockByNum(client, big.NewInt(int64(from)))
		if err != nil {
			time.Sleep(time.Hour)
		} else {
			utils.TransformData(block)
			from += 1
		}
	}
	wg.Done()
}
