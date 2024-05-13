package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/jwtoken"
	"github.com/jackc/pgx/v5"
)

type JWTokensRepo struct {
}

func NewJWTokensRepo() JWTokens {
	return &JWTokensRepo{}
}

func (r *JWTokensRepo) GetNumber(ctx context.Context, transaction Transaction, id int64, role int, purpose jwtoken.Purpose) (int64, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return 0, errors.New("InsertJWToken: error: type assertion failed on interface Transaction")
	}
	query := `
		SELECT number FROM jwtokens WHERE user_id=$1 AND role=$2 AND purpose=$3 ORDER BY number
	`
	rows, err := tx.Query(ctx, query, id, role, int(purpose))
	if err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var number int64

	if rowExist := rows.Next(); !rowExist {
		return number, nil
	}

	if err := rows.Scan(&number); err != nil {
		return 0, fmt.Errorf("scan/get token number")
	}

	nextNum := number + 1

	for rows.Next() {
		if err := rows.Scan(&number); err != nil {
			return 0, fmt.Errorf("scan next: %w", err)
		}

		if number != nextNum {
			return nextNum, nil
		}

		nextNum++
	}

	return nextNum, nil
}

func (r *JWTokensRepo) CheckJwt(ctx context.Context, transaction Transaction, tokenData *jwtoken.JWTokenData) (bool, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return false, errors.New("CheckJwt: error: type assertion failed on interface Transaction")
	}
	query := `
		SELECT 1 FROM jwtokens WHERE user_id=$1 AND role=$2 AND purpose=$3 AND number=$4 AND secret=$5
	`
	row := tx.QueryRow(ctx, query, tokenData.ID, tokenData.Role, tokenData.Purpose, tokenData.Number, tokenData.Secret)
	var res int
	if err := row.Scan(&res); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("CheckJwt/Scan: %w", err)
	}
	return true, nil
}

func (r *JWTokensRepo) CreateJwt(ctx context.Context, transaction Transaction, tokenData *jwtoken.JWTokenData) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("CreateJwt: error: type assertion failed on interface Transaction")
	}
	query := `
		INSERT INTO jwtokens VALUES($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.Exec(ctx, query, tokenData.ID, tokenData.Role,
		tokenData.Purpose, tokenData.Number, tokenData.ExpiresAt.UnixMilli(), tokenData.Secret)
	if err != nil {
		return fmt.Errorf("exec/insert jwt token: %w", err)
	}

	return nil
}

func (r *JWTokensRepo) Drop(ctx context.Context, transaction Transaction, userId int64, role int, number int64) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("Drop: error: type assertion failed on interface Transaction")
	}
	query := `
		DELETE FROM jwtokens WHERE user_id=$1 AND role=$2 AND number=$3
	`
	_, err := tx.Exec(ctx, query, userId, role, number)
	if err != nil {
		return fmt.Errorf("exec/drop tokens: %w", err)
	}

	return nil
}

func (r *JWTokensRepo) DropAll(ctx context.Context, transaction Transaction, userId int64, role int) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("DropAll: error: type assertion failed on interface Transaction")
	}
	query := `
		DELETE FROM jwtokens WHERE user_id=$1 AND role=$2
	`
	_, err := tx.Exec(ctx, query, userId, role)
	if err != nil {
		return fmt.Errorf("exec/drop all tokens: %w", err)
	}

	return nil
}

func (r *JWTokensRepo) DropAllExpired(ctx context.Context, transaction Transaction, purpose jwtoken.Purpose, now int64) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("DropAllExpired: error: type assertion failed on interface Transaction")
	}
	query := `
		DELETE FROM jwtokens WHERE purpose=$1 AND expires_at<=$2 RETURNING user_id, role, number
	`
	rows, err := tx.Query(ctx, query, int(purpose), now)
	if err != nil {
		return fmt.Errorf("DropAllExpired/exec/drop expired tokens: %w", err)
	}
	defer rows.Close()

	var userIds, roles, numbers []int64
	for rows.Next() {
		var userId, role, number int64
		if err = rows.Scan(&userId, &role, &number); err != nil {
			return fmt.Errorf("DropAllExpired/scan drop expired token: %w", err)
		}
		userIds, roles, numbers = append(userIds, userId), append(roles, role), append(numbers, number)
	}

	rows.Close()

	for i := range userIds {
		_, err = tx.Exec(ctx, `DELETE FROM jwtokens WHERE user_id=$1 AND number=$2 AND role=$3`, userIds[i], numbers[i], roles[i])
		if err != nil {
			return fmt.Errorf("DropAllExpired/drop expired access token: %w", err)
		}
	}

	return nil
}
