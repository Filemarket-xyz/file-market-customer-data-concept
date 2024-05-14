package repository

import (
	"context"
	"math/big"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/jwtoken"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoRows = pgx.ErrNoRows

type Users interface {
	InsertClient(ctx context.Context, tx Transaction, c *domain.Client) (int64, error)
	UpdateClient(ctx context.Context, tx Transaction, c *domain.Client) error
	UpdateDatasetPurchaseData(ctx context.Context, transaction Transaction, val bool, id int64) error
	GetClientById(ctx context.Context, tx Transaction, id int64) (*domain.Client, error)
	GetUserIdByAddress(ctx context.Context, tx Transaction, addr string) (int64, domain.Role, error)

	InsertAuthMessage(ctx context.Context, tx Transaction, authMsg *domain.AuthMessage) error
	GetAuthMessageByAddress(ctx context.Context, tx Transaction, address string) (*domain.AuthMessage, error)
	DeleteAuthMessage(ctx context.Context, tx Transaction, address string) error
}

type Datasets interface {
	GetDatasetsByClientId(ctx context.Context, tx Transaction, id int64) ([]*domain.Dataset, error)
}

type BlockCounterRepo interface {
	GetLastBlock(ctx context.Context) (*big.Int, error)
	SetLastBlock(ctx context.Context, lastBlock *big.Int) error
}

type EthTransactions interface {
	GetTransaction(ctx context.Context, tx Transaction, id string) (*domain.Transaction, error)
	InsertTransaction(ctx context.Context, tx Transaction, ethTx *domain.Transaction) error
}

type JWTokens interface {
	GetNumber(ctx context.Context, tx Transaction, id int64, role int, purpose jwtoken.Purpose) (int64, error)
	CheckJwt(ctx context.Context, tx Transaction, tokenData *jwtoken.JWTokenData) (bool, error)
	CreateJwt(ctx context.Context, tx Transaction, tokenData *jwtoken.JWTokenData) error
	Drop(ctx context.Context, tx Transaction, userId int64, role int, number int64) error
	DropAll(ctx context.Context, tx Transaction, userId int64, role int) error
	DropAllExpired(ctx context.Context, tx Transaction, purpose jwtoken.Purpose, now int64) error
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
type Transactions interface {
	BeginTransaction(ctx context.Context) (Transaction, error)
}

type Repository struct {
	Users
	Datasets
	BlockCounterRepo
	EthTransactions
	JWTokens

	Transactions
}

func NewRepository(cfg *config.Config, rdb *redis.Client, pool *pgxpool.Pool) (*Repository, error) {
	return &Repository{
		Users:            NewUsersRepo(cfg.Redis, rdb),
		Datasets:         NewDatasetsRepo(),
		BlockCounterRepo: NewBlockCounter(rdb),
		EthTransactions:  NewEthTransactionsRepo(),
		JWTokens:         NewJWTokensRepo(),
		Transactions:     NewTransactionsRepo(pool),
	}, nil
}
