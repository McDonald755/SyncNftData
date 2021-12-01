package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

func CheckOracleType(code []byte) bool {
	var (
		balanceOf         = "0x70a08231"
		ownerOf           = "0x6352211e"
		approve           = "0x095ea7b3"
		getApproved       = "0x081812fc"
		setApprovalForAll = "0xa22cb465"
		isApprovedForAll  = "0xe985e9c5"
		transferFrom      = "0x23b872dd"
		safeTransferFrom  = "0x42842e0e"
		safeTransferFrom2 = "0xb88d4fde"
		result            = false
	)

	b1 := strings.Contains(string(code), balanceOf)
	b2 := strings.Contains(string(code), ownerOf)
	b3 := strings.Contains(string(code), approve)
	b4 := strings.Contains(string(code), getApproved)
	b5 := strings.Contains(string(code), setApprovalForAll)
	b6 := strings.Contains(string(code), isApprovedForAll)
	b7 := strings.Contains(string(code), transferFrom)
	b8 := strings.Contains(string(code), safeTransferFrom)
	b9 := strings.Contains(string(code), safeTransferFrom2)
	if b1 && b2 && b3 && b4 && b5 && b6 && b7 && b8 && b9 {
		result = true
	}
	return result
}

func TransferAddr(addrs *[]string) *[]common.Address {
	return nil
}
