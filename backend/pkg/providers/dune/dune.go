package dune

import (
	"context"
	"time"
)

//go:generate mockgen -source dune.go -destination ../../../mocks/pkg/dune/dune.go
type Dune interface {
	ExecuteQuery(ctx context.Context, queryId string, req *Req) (*ExecuteResponse, error)
	GetLatestQueryResult(ctx context.Context, queryId, filter string) (*QueryResultResponse, error)
	GetExecutionResult(ctx context.Context, executionId, filter string) (*QueryResultResponse, error)

	LayerZeroUserFees(ctx context.Context, address string) (*QueryResultResponse, error)
}

const (
	MultiChainBalancesQueryId = "3737145"
	WalletBalanceByBlockchain = "3737347"
	ZkSyncUsersStatQueryId    = "3736876"
	RankUserQueryId           = "3735555"

	BlockchainEth      = "ethereum"
	BlockchainOptimism = "optimism"
	BlockchainBase     = "base"
	BlockchainArbitrum = "arbitrum"

	COMPLETEDStateResult = "QUERY_STATE_COMPLETED"
	FAILEDStateResult    = "QUERY_STATE_FAILED"
)

type (
	Req struct {
		RequsetData *RequsetData `json:"query_parameters"`
	}
	RequsetData struct {
		Addr       string `json:"addr,omitempty"`
		Blockchain string `json:"blockchain,omitempty"`
	}

	ExecuteResponse struct {
		ExecutionId string `json:"execution_id"`
		State       string `json:"state"`
	}

	QueryResultResponse struct {
		ExecutionID         string    `json:"execution_id"`
		QueryID             int       `json:"query_id"`
		State               string    `json:"state"`
		IsExecutionFinished bool      `json:"is_execution_finished"`
		SubmittedAt         time.Time `json:"submitted_at"`
		ExpiresAt           time.Time `json:"expires_at"`
		ExecutionStartedAt  time.Time `json:"execution_started_at"`
		ExecutionEndedAt    time.Time `json:"execution_ended_at"`
		NextOffset          int       `json:"next_offset"`
		NextURI             string    `json:"next_uri"`
		Result              *Result   `json:"result"`
	}

	Result struct {
		Metadata *Metadata        `json:"metadata"`
		Rows     []map[string]any `json:"rows"`
	}

	Metadata struct {
		ColumnNames         []string `json:"column_names"`
		ColumnTypes         []string `json:"column_types"`
		ResultSetBytes      int      `json:"result_set_bytes"`
		RowCount            int      `json:"row_count"`
		TotalResultSetBytes int      `json:"total_result_set_bytes"`
		TotalRowCount       int      `json:"total_row_count"`
		DatapointCount      int      `json:"datapoint_count"`
		PendingTimeMillis   int      `json:"pending_time_millis"`
		ExecutionTimeMillis int      `json:"execution_time_millis"`
	}
)
