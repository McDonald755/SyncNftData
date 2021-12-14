package ES

import (
	"time"
)

type EsEnity struct {
	ID            string    `json:"id"`
	CreatedTime   time.Time `json:"created_time"`
	UpdatedTime   time.Time `json:"updated_time"`
	TokenId       string    `json:"token_id"`
	TokenUri      string    `json:"token_uri"`
	Owner         string    `json:"owner"`
	OracleAddr    string    `json:"oracle_addr"`
	TokenApproval string    `json:"token_approval"`
}
