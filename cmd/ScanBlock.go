package cmd

import (
	"SyncNftData/config"
	"SyncNftData/db"
	"SyncNftData/oracle"
	"SyncNftData/syncData"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/spf13/cobra"
	"strings"
	"sync"
)

func ScanCmd() *cobra.Command {
	var blockNum int64
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "s",
		Long:  "It will sync the latest block ",
		RunE: func(cmd *cobra.Command, args []string) error {
			wg := sync.WaitGroup{}
			Oracles := db.TGetOracleAddrAll()
			//init standard ERC-721 contract data
			contractABI, err := abi.JSON(strings.NewReader(oracle.OracleABI))
			if err != nil {
				fmt.Println(err)
			}
			clientNum := len(config.CLIENTS)
			oracleNum := len(Oracles)
			num := oracleNum / clientNum

			for i := 0; i < clientNum; i++ {
				from := num*i + int(blockNum)
				wg.Add(1)
				go syncData.TSyncData(config.CLIENTS[0], int64(from), Oracles[num*i:num*(i+1)], contractABI, &wg)
			}

			return nil
		},
	}

	scanCmd.Flags().Int64VarP(&blockNum, "blockNum", "b", 0, "input blockNum")
	return scanCmd
}
