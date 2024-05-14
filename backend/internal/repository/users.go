package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

type UsersRepo struct {
	cfg *config.RedisConfig
	rdb *redis.Client
}

func NewUsersRepo(
	cfg *config.RedisConfig,
	rdb *redis.Client,
) Users {
	return &UsersRepo{
		cfg: cfg,
		rdb: rdb,
	}
}

func (r *UsersRepo) InsertClient(ctx context.Context, transaction Transaction, c *domain.Client) (int64, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return 0, errors.New("InsertClient: error: type assertion failed on interface Transaction")
	}

	row := tx.QueryRow(ctx, `INSERT INTO clients (id, agreement, bought, point_balance) 
		VALUES (DEFAULT, $1,$2,$3) RETURNING id`,
		c.Agreement, c.Bought, c.PointBalance.String())
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("InsertClient/Scan: %w", err)
	}

	if _, err := tx.Exec(ctx, `INSERT INTO users (id, role, address) 
		VALUES ($1,$2,$3)`,
		id, domain.RoleClient, strings.ToLower(c.Address.String())); err != nil {
		return 0, fmt.Errorf("InsertClient/Exec: %w", err)
	}
	return id, nil
}

func (r *UsersRepo) UpdateClient(ctx context.Context, transaction Transaction, c *domain.Client) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("UpdateClient: error: type assertion failed on interface Transaction")
	}
	if _, err := tx.Exec(ctx, `UPDATE clients SET agreement=$1, bought=$2, point_balance=$3 WHERE id=$4`,
		c.Agreement, c.Bought, c.PointBalance.String(), c.Id); err != nil {
		return fmt.Errorf("UpdateClient/Exec: %w", err)
	}
	return nil
}

func (r *UsersRepo) UpdateDatasetPurchaseData(ctx context.Context, transaction Transaction, val bool, id int64) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("UpdateDatasetPurchaseData: error: type assertion failed on interface Transaction")
	}
	if _, err := tx.Exec(ctx, `UPDATE clients SET bought=$1 WHERE id=$2`,
		val, id); err != nil {
		return fmt.Errorf("UpdateDatasetPurchaseData/Exec: %w", err)
	}
	return nil
}

func (r *UsersRepo) GetClientById(ctx context.Context, transaction Transaction, id int64) (*domain.Client, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, errors.New("GetClientById: error: type assertion failed on interface Transaction")
	}

	row := tx.QueryRow(ctx, `SELECT c.id, u.address, c.agreement, c.bought, c.point_balance
		FROM clients AS c 
		JOIN users AS u ON u.id = c.id AND c.id=$1
	`, id)

	var (
		c    = &domain.Client{}
		addr string
	)
	err := row.Scan(&c.Id, &addr, &c.Agreement, &c.Bought, &c.PointBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("GetClientById/Scan: %w", err)
	}

	c.Address = common.HexToAddress(addr)

	return c, nil
}

func (r *UsersRepo) GetUserIdByAddress(ctx context.Context, transaction Transaction, addrQ string) (int64, domain.Role, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return 0, 0, errors.New("GetUserIdByAddress: error: type assertion failed on interface Transaction")
	}

	row := tx.QueryRow(ctx, `SELECT id, role FROM users WHERE address=$1`, strings.ToLower(addrQ))

	var (
		id   int64
		role int
	)
	err := row.Scan(&id, &role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, ErrNoRows
		}
		return 0, 0, fmt.Errorf("GetUserIdByAddress/Scan: %w", err)
	}

	return id, domain.Role(role), nil
}

func (r *UsersRepo) GetAuthMessageByAddress(
	ctx context.Context,
	transaction Transaction,
	address string,
) (*domain.AuthMessage, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, errors.New("GetAuthMessageByAddress: error: type assertion failed on interface Transaction")
	}
	row := tx.QueryRow(ctx, `SELECT address, created_at, code FROM auth_messages WHERE address = $1`, strings.ToLower(address))
	res := &domain.AuthMessage{}
	if err := row.Scan(&res.Address, &res.CreatedAt, &res.Message); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("GetAuthMessageByAddress/Scan: %w", err)
	}
	return res, nil
}

func (r *UsersRepo) InsertAuthMessage(
	ctx context.Context,
	transaction Transaction,
	msg *domain.AuthMessage,
) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("InsertAuthMessage: error: type assertion failed on interface Transaction")
	}
	if _, err := tx.Exec(ctx, `INSERT INTO auth_messages (address, code, created_at) VALUES ($1,$2,$3)`,
		strings.ToLower(msg.Address), msg.Message, msg.CreatedAt); err != nil {
		return fmt.Errorf("InsertAuthMessage/Exec: %w", err)
	}
	return nil
}

func (r *UsersRepo) DeleteAuthMessage(
	ctx context.Context,
	transaction Transaction,
	address string,
) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("DeleteAuthMessage: error: type assertion failed on interface Transaction")
	}
	if _, err := tx.Exec(ctx, `DELETE FROM auth_messages WHERE address=$1`, strings.ToLower(address)); err != nil {
		return fmt.Errorf("DeleteAuthMessage/Exec: %w", err)
	}
	return nil
}
