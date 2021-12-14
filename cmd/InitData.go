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
Used to init data
*/
func InitCmd() *cobra.Command {
	var startNum int64
	var endNum int64
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "i",
		Long:  "It will init data",
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
			gap := int64(1000)
			newNum, _ := config.CLIENTS[0].BlockNumber(context.Background())
			for i := 0; i < clientLen; i++ {
				time.Sleep(time.Millisecond * 300)
				wg.Add(1)
				go syncData.InitData(config.CLIENTS[i], int64(newNum), Oracles[i*num:(i+1)*num], contractABI, &wg, int64(startNum), gap, i)
			}
			wg.Wait()
			return nil
		},
	}

	initCmd.Flags().Int64VarP(&startNum, "startNum", "s", 0, "input blockNum")
	initCmd.Flags().Int64VarP(&endNum, "endNum", "e", 0, "input blockNum")
	return initCmd
}
