package cmd

import (
	"SyncNftData/syncData"
	"github.com/spf13/cobra"
)

func ScanCmd() *cobra.Command {
	var blockNum int64
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "s",
		Long:  "It will sync the latest block ",
		RunE: func(cmd *cobra.Command, args []string) error {
			syncData.SyncData(blockNum)
			return nil
		},
	}

	scanCmd.Flags().Int64VarP(&blockNum, "blockNum", "b", 0, "input blockNum")
	return scanCmd
}
