-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE public.datasets
(
    client_id BIGINT UNIQUE NOT NULL REFERENCES clients (id) ON DELETE CASCADE,

    -- MultiChainBalances
    chain_multi_balance VARCHAR(256),
    txs_multi_balance BIGINT,
    Chain_Balance_multi_balance REAL,
    Chain_Received_multi_balance REAL,
    Chain_Sent_multi_balance REAL,
    Chain_TX_Fee_multi_balance REAL,
    USD_Balance_multi_balance REAL,
    USD_Received_multi_balance REAL,
    USD_Sent_multi_balance REAL,
    USD_TX_Fee_multi_balance REAL,
    Failed_multi_balance BIGINT,

    -- Rank LayerZero User
    rk_rank BIGINT,
    rs_rank BIGINT,
    tc_rank BIGINT,
    amt_rank REAL,
    cc_rank VARCHAR(256),
    dwm_rank VARCHAR(256),
    ibt_rank VARCHAR(256),
    lzd_rank BIGINT,

    -- ZkSyncEraUserStat
    rk_zk_sync_era BIGINT,
    amt_zk_sync_era VARCHAR(256),
    cc_zk_sync_era VARCHAR(256),
    dwm_zk_sync_era VARCHAR(256),
    ibt_zk_sync_era VARCHAR(256),
    tc_zk_sync_era VARCHAR(256),

    -- LayerZero User Fees
    rank_lz_fees BIGINT,
    count_tx_lz_fees BIGINT,
    total_amount_lz_fees REAL,
    Arbitrum_lz_fees REAL,
    Avalanche_lz_fees REAL,
    BSC_lz_fees REAL,
    Base_lz_fees REAL,
    Celo_lz_fees REAL,
    Ethereum_lz_fees REAL,
    Fantom_lz_fees REAL,
    Gnosis_lz_fees REAL,
    Linea_lz_fees REAL,
    Optimism_lz_fees REAL,
    Polygon_lz_fees REAL,
    Scroll_lz_fees REAL,
    Zkevm_lz_fees REAL,
    Zksync_lz_fees REAL,
    Zora_lz_fees REAL
);

-- Отложим
-- CREATE TABLE public.balances
-- (
--     client_id BIGINT NOT NULL REFERENCES clients (id) ON DELETE CASCADE,
--     symbol  VARCHAR(256) NOT NULL,
--     balance REAL NOT NULL,
--     balance_usd REAL NOT NULL,
--     token_address VARCHAR(42) NOT NULL
-- );
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE public.datasets;
