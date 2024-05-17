package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

type DatasetsRepo struct {
}

func NewDatasetsRepo() Datasets {
	return &DatasetsRepo{}
}

func (r *DatasetsRepo) GetDatasetsByClientId(ctx context.Context, transaction Transaction, clientId int64) (*domain.Dataset, error) {

	// layerZeroUserFeesDataset, err := r.GetLayerZeroUserFeesDataset(ctx, transaction, clientId)
	// if err != nil && !errors.Is(err, ErrNoRows) {
	// 	return nil, fmt.Errorf("GetDatasetsByClientId/GetLayerZeroUserFeesDataset: %w", err)
	// }

	multiChainBalancesDataset, err := r.GetMultiChainBalancesDataset(ctx, transaction, clientId)
	if err != nil && !errors.Is(err, ErrNoRows) {
		return nil, fmt.Errorf("GetDatasetsByClientId/GetMultiChainBalancesDataset: %w", err)
	}

	rankZeroUserDataset, err := r.GetRankLayerZeroUserDataset(ctx, transaction, clientId)
	if err != nil && !errors.Is(err, ErrNoRows) {
		return nil, fmt.Errorf("GetDatasetsByClientId/GetRankLayerZeroUserDataset: %w", err)
	}

	zkSyncEraUserStatDataset, err := r.GetZkSyncEraUserStatDataset(ctx, transaction, clientId)
	if err != nil && !errors.Is(err, ErrNoRows) {
		return nil, fmt.Errorf("GetDatasetsByClientId/GetZkSyncEraUserStatDataset: %w", err)
	}

	res := &domain.Dataset{
		// LayerZeroUserFees: layerZeroUserFeesDataset,
		MultiChainBalances: multiChainBalancesDataset,
		RankLZ:             rankZeroUserDataset,
		ZkSyncEraUserStat:  zkSyncEraUserStatDataset,
	}

	return res, nil
}

func (r *DatasetsRepo) DeleteByClientId(ctx context.Context, transaction Transaction, id int64) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("DeleteByClientId: error: type assertion failed on interface Transaction")
	}
	query := `
		DELETE FROM datasets WHERE client_id=$1
	`
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeleteByClientId/Exec: %w", err)
	}

	return nil
}

// func (r *DatasetsRepo) AddBalance(ctx context.Context, transaction Transaction, clientId int64, b *domain.Balance) error {
// 	tx, ok := transaction.(pgx.Tx)
// 	if !ok {
// 		return errors.New("AddBalance: error: type assertion failed on interface Transaction")
// 	}

// 	if _, err := tx.Exec(ctx, `
// 		INSERT INTO datasets AS d (
// 			client_id,
// 			rank_lz_fees,
// 			count_tx_lz_fees,
// 			total_amount_lz_fees,
// 			Arbitrum_lz_fees,
// 			Avalanche_lz_fees,
// 			BSC_lz_fees,
// 			Base_lz_fees,
// 			Celo_lz_fees,
// 			Ethereum_lz_fees,
// 			Fantom_lz_fees,
// 			Gnosis_lz_fees,
// 			Linea_lz_fees,
// 			Optimism_lz_fees,
// 			Polygon_lz_fees,
// 			Scroll_lz_fees,
// 			Zkevm_lz_fees,
// 			Zksync_lz_fees,
// 			Zora_lz_fees
// 		)
// 		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
// 		ON CONFLICT ON CONSTRAINT datasets_client_id_key
// 		DO UPDATE SET
// 			rank_lz_fees = excluded.rank_lz_fees,
// 			count_tx_lz_fees = excluded.count_tx_lz_fees,
// 			total_amount_lz_fees = excluded.total_amount_lz_fees,
// 			Arbitrum_lz_fees = excluded.Arbitrum_lz_fees,
// 			Avalanche_lz_fees = excluded.Avalanche_lz_fees,
// 			BSC_lz_fees = excluded.BSC_lz_fees,
// 			Base_lz_fees = excluded.Base_lz_fees,
// 			Celo_lz_fees = excluded.Celo_lz_fees,
// 			Ethereum_lz_fees = excluded.Ethereum_lz_fees,
// 			Fantom_lz_fees = excluded.Fantom_lz_fees,
// 			Gnosis_lz_fees = excluded.Gnosis_lz_fees,
// 			Linea_lz_fees = excluded.Linea_lz_fees,
// 			Optimism_lz_fees = excluded.Optimism_lz_fees,
// 			Polygon_lz_fees = excluded.Polygon_lz_fees,
// 			Scroll_lz_fees = excluded.Scroll_lz_fees,
// 			Zkevm_lz_fees = excluded.Zkevm_lz_fees,
// 			Zksync_lz_fees = excluded.Zksync_lz_fees,
// 			Zora_lz_fees = excluded.Zora_lz_fees
// 		WHERE d.client_id=excluded.client_id`,
// 		clientId, d.Rank, d.CountTx, d.TotalAmount, d.Arbitrum, d.Avalanche, d.BSC, d.Base, d.Celo,
// 		d.Ethereum, d.Fantom, d.Gnosis, d.Linea, d.Optimism, d.Polygon, d.Scroll, d.Zkevm, d.Zksync, d.Zora); err != nil {
// 		return fmt.Errorf("AddBalance/Exec: %w", err)
// 	}

// 	return nil
// }

func (r *DatasetsRepo) AddLayerZeroUserFees(ctx context.Context, transaction Transaction, clientId int64, d *domain.LayerZeroUserFeesDataset) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("AddLayerZeroUserFees: error: type assertion failed on interface Transaction")
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO datasets AS d (
			client_id, 
			rank_lz_fees,
			count_tx_lz_fees,
			total_amount_lz_fees,
			Arbitrum_lz_fees,
			Avalanche_lz_fees,
			BSC_lz_fees,
			Base_lz_fees,
			Celo_lz_fees,
			Ethereum_lz_fees,
			Fantom_lz_fees,
			Gnosis_lz_fees,
			Linea_lz_fees,
			Optimism_lz_fees,
			Polygon_lz_fees,
			Scroll_lz_fees,
			Zkevm_lz_fees,
			Zksync_lz_fees,
			Zora_lz_fees
		) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		ON CONFLICT ON CONSTRAINT datasets_client_id_key
		DO UPDATE SET
			rank_lz_fees = excluded.rank_lz_fees,
			count_tx_lz_fees = excluded.count_tx_lz_fees,
			total_amount_lz_fees = excluded.total_amount_lz_fees,
			Arbitrum_lz_fees = excluded.Arbitrum_lz_fees,
			Avalanche_lz_fees = excluded.Avalanche_lz_fees,
			BSC_lz_fees = excluded.BSC_lz_fees,
			Base_lz_fees = excluded.Base_lz_fees,
			Celo_lz_fees = excluded.Celo_lz_fees,
			Ethereum_lz_fees = excluded.Ethereum_lz_fees,
			Fantom_lz_fees = excluded.Fantom_lz_fees,
			Gnosis_lz_fees = excluded.Gnosis_lz_fees,
			Linea_lz_fees = excluded.Linea_lz_fees,
			Optimism_lz_fees = excluded.Optimism_lz_fees,
			Polygon_lz_fees = excluded.Polygon_lz_fees,
			Scroll_lz_fees = excluded.Scroll_lz_fees,
			Zkevm_lz_fees = excluded.Zkevm_lz_fees,
			Zksync_lz_fees = excluded.Zksync_lz_fees,
			Zora_lz_fees = excluded.Zora_lz_fees
		WHERE d.client_id=excluded.client_id`,
		clientId, d.Rank, d.CountTx, d.TotalAmount, d.Arbitrum, d.Avalanche, d.BSC, d.Base, d.Celo,
		d.Ethereum, d.Fantom, d.Gnosis, d.Linea, d.Optimism, d.Polygon, d.Scroll, d.Zkevm, d.Zksync, d.Zora); err != nil {
		return fmt.Errorf("AddLayerZeroUserFees/Exec: %w", err)
	}

	return nil
}

func (r *DatasetsRepo) GetLayerZeroUserFeesDataset(ctx context.Context, transaction Transaction, clientId int64) (*domain.LayerZeroUserFeesDataset, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, errors.New("GetLayerZeroUserFeesDataset: error: type assertion failed on interface Transaction")
	}

	row := tx.QueryRow(ctx, `
	SELECT 
		rank_lz_fees,
		count_tx_lz_fees,
		total_amount_lz_fees,
		Arbitrum_lz_fees,
		Avalanche_lz_fees,
		BSC_lz_fees,
		Base_lz_fees,
		Celo_lz_fees,
		Ethereum_lz_fees,
		Fantom_lz_fees,
		Gnosis_lz_fees,
		Linea_lz_fees,
		Optimism_lz_fees,
		Polygon_lz_fees,
		Scroll_lz_fees,
		Zkevm_lz_fees,
		Zksync_lz_fees,
		Zora_lz_fees
	FROM datasets WHERE client_id = $1
	`, clientId)

	var (
		d = &domain.LayerZeroUserFeesDataset{}
	)
	err := row.Scan(&d.Rank, &d.CountTx, &d.TotalAmount, &d.Arbitrum, &d.Avalanche, &d.BSC, &d.Base, &d.Celo,
		&d.Ethereum, &d.Fantom, &d.Gnosis, &d.Linea, &d.Optimism, &d.Polygon, &d.Scroll, &d.Zkevm, &d.Zksync, &d.Zora)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("GetLayerZeroUserFeesDataset/Scan: %w", err)
	}

	return d, nil
}

func (r *DatasetsRepo) AddRankLayerZeroUser(ctx context.Context, transaction Transaction, clientId int64, d *domain.RankLZDataset) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("AddRankLayerZeroUser: error: type assertion failed on interface Transaction")
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO datasets AS d (
			client_id, 
			rk_rank,
			rs_rank,
			tc_rank,
			amt_rank,
			cc_rank,
			dwm_rank,
			ibt_rank,
			lzd_rank
		) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT ON CONSTRAINT datasets_client_id_key
		DO UPDATE SET
			rk_rank = excluded.rk_rank,
			rs_rank = excluded.rs_rank,
			tc_rank = excluded.tc_rank,
			amt_rank = excluded.amt_rank,
			cc_rank = excluded.cc_rank,
			dwm_rank = excluded.dwm_rank,
			ibt_rank = excluded.ibt_rank,
			lzd_rank = excluded.lzd_rank
		WHERE d.client_id=excluded.client_id`,
		clientId, d.Rk, d.Rs, d.Tc, d.Amt, d.Cc, d.Dwm, d.Ibt, d.Lzd); err != nil {
		return fmt.Errorf("AddRankLayerZeroUser/Exec: %w", err)
	}

	return nil
}

func (r *DatasetsRepo) GetRankLayerZeroUserDataset(ctx context.Context, transaction Transaction, clientId int64) (*domain.RankLZDataset, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, errors.New("GetRankLayerZeroUserDataset: error: type assertion failed on interface Transaction")
	}

	row := tx.QueryRow(ctx, `
	SELECT 
		rk_rank,
		rs_rank,
		tc_rank,
		amt_rank,
		cc_rank,
		dwm_rank,
		ibt_rank,
		lzd_rank
	FROM datasets WHERE client_id = $1
	`, clientId)

	var (
		d = &domain.RankLZDataset{}
	)
	err := row.Scan(&d.Rk, &d.Rs, &d.Tc, &d.Amt, &d.Cc, &d.Dwm, &d.Ibt, &d.Lzd)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("GetRankLayerZeroUserDataset/Scan: %w", err)
	}

	return d, nil
}

func (r *DatasetsRepo) AddZkSyncEraUserStat(ctx context.Context, transaction Transaction, clientId int64, d *domain.ZkSyncEraUserStatDataset) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("AddZkSyncEraUserStat: error: type assertion failed on interface Transaction")
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO datasets AS d (
			client_id, 
			rk_zk_sync_era,
			amt_zk_sync_era,
			cc_zk_sync_era,
			dwm_zk_sync_era,
			ibt_zk_sync_era,
			tc_zk_sync_era
		) 
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT ON CONSTRAINT datasets_client_id_key
		DO UPDATE SET
		rk_zk_sync_era = excluded.rk_zk_sync_era,
		amt_zk_sync_era = excluded.amt_zk_sync_era,
		cc_zk_sync_era = excluded.cc_zk_sync_era,
		dwm_zk_sync_era = excluded.dwm_zk_sync_era,
		ibt_zk_sync_era = excluded.ibt_zk_sync_era,
		tc_zk_sync_era = excluded.tc_zk_sync_era
		WHERE d.client_id=excluded.client_id`,
		clientId, d.Rk, d.Amt, d.Cc, d.Dwm, d.Ibt, d.Tc); err != nil {
		return fmt.Errorf("AddZkSyncEraUserStat/Exec: %w", err)
	}

	return nil
}

func (r *DatasetsRepo) GetZkSyncEraUserStatDataset(ctx context.Context, transaction Transaction, clientId int64) (*domain.ZkSyncEraUserStatDataset, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, errors.New("GetZkSyncEraUserStatDataset: error: type assertion failed on interface Transaction")
	}

	row := tx.QueryRow(ctx, `
	SELECT 
		rk_zk_sync_era,
		amt_zk_sync_era,
		cc_zk_sync_era,
		dwm_zk_sync_era,
		ibt_zk_sync_era,
		tc_zk_sync_era
	FROM datasets WHERE client_id = $1
	`, clientId)

	var (
		d = &domain.ZkSyncEraUserStatDataset{}
	)
	err := row.Scan(&d.Rk, &d.Amt, &d.Cc, &d.Dwm, &d.Ibt, &d.Tc)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("GetZkSyncEraUserStatDataset/Scan: %w", err)
	}

	return d, nil
}

func (r *DatasetsRepo) AddMultiChainBalancesDataset(ctx context.Context, transaction Transaction, clientId int64, d *domain.MultiChainBalancesDataset) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("AddMultiChainBalancesDataset: error: type assertion failed on interface Transaction")
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO datasets AS d (
			client_id, 
			chain_multi_balance,
			txs_multi_balance,
			Chain_Balance_multi_balance,
			Chain_Received_multi_balance,
			Chain_Sent_multi_balance,
			Chain_TX_Fee_multi_balance,
			USD_Balance_multi_balance,
			USD_Received_multi_balance,
			USD_Sent_multi_balance,
			USD_TX_Fee_multi_balance,
			Failed_multi_balance
		) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT ON CONSTRAINT datasets_client_id_key
		DO UPDATE SET
			chain_multi_balance = excluded.chain_multi_balance,
			txs_multi_balance = excluded.txs_multi_balance,
			Chain_Balance_multi_balance = excluded.Chain_Balance_multi_balance,
			Chain_Received_multi_balance = excluded.Chain_Received_multi_balance,
			Chain_Sent_multi_balance = excluded.Chain_Sent_multi_balance,
			Chain_TX_Fee_multi_balance = excluded.Chain_TX_Fee_multi_balance,
			USD_Balance_multi_balance = excluded.USD_Balance_multi_balance,
			USD_Received_multi_balance = excluded.USD_Received_multi_balance,
			USD_Sent_multi_balance = excluded.USD_Sent_multi_balance,
			USD_TX_Fee_multi_balance = excluded.USD_TX_Fee_multi_balance,
			Failed_multi_balance = excluded.Failed_multi_balance
		WHERE d.client_id=excluded.client_id`,
		clientId, d.Chain, d.Txs, d.Chain_Balance, d.Chain_Received, d.Chain_Sent, d.Chain_TX_Fee,
		d.USD_Balance, d.USD_Received, d.USD_Sent, d.USD_TX_Fee, d.Failed); err != nil {
		return fmt.Errorf("AddMultiChainBalancesDataset/Exec: %w", err)
	}

	return nil
}

func (r *DatasetsRepo) GetMultiChainBalancesDataset(ctx context.Context, transaction Transaction, clientId int64) (*domain.MultiChainBalancesDataset, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, errors.New("GetMultiChainBalancesDataset: error: type assertion failed on interface Transaction")
	}

	row := tx.QueryRow(ctx, `
	SELECT 
		chain_multi_balance,
		txs_multi_balance,
		Chain_Balance_multi_balance,
		Chain_Received_multi_balance,
		Chain_Sent_multi_balance,
		Chain_TX_Fee_multi_balance,
		USD_Balance_multi_balance,
		USD_Received_multi_balance,
		USD_Sent_multi_balance,
		USD_TX_Fee_multi_balance,
		Failed_multi_balance
	FROM datasets WHERE client_id = $1
	`, clientId)

	var (
		d = &domain.MultiChainBalancesDataset{}
	)
	err := row.Scan(&d.Chain, &d.Txs, &d.Chain_Balance, &d.Chain_Received, &d.Chain_Sent, &d.Chain_TX_Fee,
		&d.USD_Balance, &d.USD_Received, &d.USD_Sent, &d.USD_TX_Fee, &d.Failed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("GetMultiChainBalancesDataset/Scan: %w", err)
	}

	return d, nil
}
