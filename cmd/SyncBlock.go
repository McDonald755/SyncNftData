package cmd

import (
	"SyncNftData/config"
	"SyncNftData/db"
	"SyncNftData/oracle"
	"SyncNftData/syncData"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/cobra"
	"strings"
	"sync"
	"time"
)

/**
Used to synchronize data
*/
func SyncCmd() *cobra.Command {
	var startNum int64
	var syncNum int64
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "s",
		Long:  "It will sync the latest block ",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			for i := 0; i < clientLen-1; i++ {
				time.Sleep(time.Millisecond * 300)
				wg.Add(1)
				go syncData.TSyncData(config.CLIENTS[i], int64(startNum), Oracles[i*num:(i+1)*num], contractABI, &wg, int64(number), i)
			}
			wg.Add(1)
			go syncData.SyncData(config.CLIENTS[clientLen-1], int64(syncNum), Oracles, contractABI, &wg)
			wg.Wait()
			return nil
		},
	}

	syncCmd.Flags().Int64VarP(&startNum, "startNum", "s", 0, "input blockNum")
	syncCmd.Flags().Int64VarP(&syncNum, "syncNum", "n", 0, "input blockNum")
	return syncCmd
}
