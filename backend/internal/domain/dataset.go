package domain

import (
	"strings"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"
	"github.com/gocarina/gocsv"
)

type Dataset struct {
	// LayerZeroUserFees *LayerZeroUserFeesDataset `csv:"LayerZeroUserFeesDataset"`
	RankLZ             *RankLZDataset             `csv:"Rank_LayerZero"`       // https://dune.com/queries/2714717/4518403
	ZkSyncEraUserStat  *ZkSyncEraUserStatDataset  `csv:"ZkSync_Era_User_Stat"` // https://dune.com/queries/3251197/5440625
	MultiChainBalances *MultiChainBalancesDataset `csv:"Multi_Chain_Balances"` // https://dune.com/queries/3686257/6200465
}

func SliceDataToStrings(in interface{}) ([][]string, error) {
	out, err := gocsv.MarshalString(in)
	if err != nil {
		return nil, err
	}
	rows := strings.Fields(out)
	res := make([][]string, len(rows))

	for i, row := range rows {
		res[i] = strings.Split(row, ",")
	}

	return res, nil
}

type (
	Balance struct { // https://dune.com/queries/3459600/5814175
		Symbol       string  `csv:"Symbol"`
		Balance      float64 `csv:"Balance"`
		BalanceUsd   float64 `csv:"Balance_usd"`
		TokenAddress string  `csv:"Token_address"`
	}

	RankLZDataset struct {
		Rk  *int64   `csv:"Ranking#"`
		Rs  *int64   `csv:"Ranking_Score"`
		Tc  *int64   `csv:"Transactions_Count"`
		Amt *float64 `csv:"Bridged_Amount($)"`
		Cc  *string  `csv:"Interacted_Source_Chains/Destination_Chains/Contracts_Count"`
		Dwm *string  `csv:"Unique_Active_Days/Weeks/Months"`
		Ibt *string  `csv:"Initial_Active_Date"`
		Lzd *int64   `csv:"LZ_Age_In_Days"`
	}

	MultiChainBalancesDataset struct {
		Chain          *string  `csv:"Chain"`
		Txs            *int64   `csv:"Transactions"`
		Failed         *int64   `csv:"Failed"`
		Chain_Balance  *float64 `csv:"Chain_Balance"`
		Chain_Received *float64 `csv:"Chain_Received"`
		Chain_Sent     *float64 `csv:"Chain_Sent"`
		Chain_TX_Fee   *float64 `csv:"Chain_TX_Fee"`
		USD_Balance    *float64 `csv:"USD_Balance"`
		USD_Received   *float64 `csv:"USD_Received"`
		USD_Sent       *float64 `csv:"USD_Sent"`
		USD_TX_Fee     *float64 `csv:"USD_TX_Fee"`
	}

	ZkSyncEraUserStatDataset struct {
		Rk  *int64  `csv:"Rank"`
		Amt *string `csv:"Volume($)"`
		Cc  *string `csv:"Contracts"`
		Dwm *string `csv:"DistinctDays/Weeks/Months"`
		Ibt *string `csv:"InitialTxTime"`
		Tc  *string `csv:"Transactions"`
	}

	LayerZeroUserFeesDataset struct {
		Rank        *int64   `csv:"rank_lz_fees"`
		CountTx     *int64   `csv:"count_tx_lz_fees"`
		TotalAmount *int64   `csv:"total_amount_lz_fees"`
		Arbitrum    *float64 `csv:"Arbitrum_lz_fees"`
		Avalanche   *float64 `csv:"Avalanche_lz_fees"`
		BSC         *float64 `csv:"BSC_lz_fees"`
		Base        *float64 `csv:"Base_lz_fees"`
		Celo        *float64 `csv:"Celo_lz_fees"`
		Ethereum    *float64 `csv:"Ethereum_lz_fees"`
		Fantom      *float64 `csv:"Fantom_lz_fees"`
		Gnosis      *float64 `csv:"Gnosis_lz_fees"`
		Linea       *float64 `csv:"Linea_lz_fees"`
		Optimism    *float64 `csv:"Optimism_lz_fees"`
		Polygon     *float64 `csv:"Polygon_lz_fees"`
		Scroll      *float64 `csv:"Scroll_lz_fees"`
		Zkevm       *float64 `csv:"Zkevm_lz_fees"`
		Zksync      *float64 `csv:"Zksync_lz_fees"`
		Zora        *float64 `csv:"Zora_lz_fees"`
	}
)

func (d *Dataset) ToModel() []*models.Dataset {
	var res []*models.Dataset

	data, err := SliceDataToStrings([]*Dataset{d})
	if err != nil {
		return nil
	}
	if len(data) < 2 {
		return nil
	}
	if len(data[0]) != len(data[1]) {
		return nil
	}

	for i, d := range data[0] {
		res = append(res, &models.Dataset{
			Name: d,
			Data: data[1][i],
		})
	}

	return res
}
