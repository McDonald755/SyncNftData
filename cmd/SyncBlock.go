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
			number, _ := config.CLIENT.BlockNumber(context.Background())

			wg.Add(1)
			go syncData.TSyncData(config.CLIENT, int64(startNum), Oracles, contractABI, &wg, int64(number), 1)

			wg.Add(1)
			go syncData.SyncData(config.CLIENT, int64(syncNum), Oracles, contractABI, &wg)
			wg.Wait()
			return nil
		},
	}

	syncCmd.Flags().Int64VarP(&startNum, "startNum", "s", 0, "input blockNum")
	syncCmd.Flags().Int64VarP(&syncNum, "syncNum", "n", 0, "input blockNum")
	return syncCmd
}
