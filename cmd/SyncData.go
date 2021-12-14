package cmd

import (
	"SyncNftData/config"
	"SyncNftData/db"
	"SyncNftData/oracle"
	"SyncNftData/syncData"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/cobra"

	"strings"
	"sync"
	"time"
)

/**
Used to sync data
*/
func SyncCmd() *cobra.Command {
	var startNum int64
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "s",
		Long:  "It will sync the latest block ",
		RunE: func(cmd *cobra.Command, args []string) error {
			wg := sync.WaitGroup{}

			//init oracle data
			Oracles := db.GetOracleAddrAll()
			//init standard ERC-721 contract data
			contractABI, _ := abi.JSON(strings.NewReader(oracle.OracleABI))

			//init data
			oracleNum := len(Oracles)
			clientLen := len(config.CLIENTS)
			num := oracleNum / clientLen
			gap := int64(1)
			for i := 0; i < clientLen-1; i++ {
				time.Sleep(time.Millisecond * 300)
				wg.Add(1)
				go syncData.SyncData(config.CLIENTS[i], int64(startNum), Oracles[i*num:(i+1)*num], contractABI, &wg, gap, i)
			}
			wg.Add(1)
			//sync new contract data
			go syncData.SyncDataMain(config.CLIENTS[clientLen-1], int64(startNum), Oracles, contractABI, &wg)
			wg.Wait()
			return nil
		},
	}

	syncCmd.Flags().Int64VarP(&startNum, "startNum", "s", 0, "input blockNum")
	return syncCmd
}
