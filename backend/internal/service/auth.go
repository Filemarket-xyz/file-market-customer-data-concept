package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/repository"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/hash"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/jwtoken"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/rand_manager"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/time_manager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type AuthService struct {
	cfg              *config.ServiceConfig
	repoUsers        repository.Users
	repoJWTokens     repository.JWTokens
	repoTransactions repository.Transactions
	jwtManager       jwtoken.JWTokenManager
	hashManager      hash.HashManager
	timeManager      time_manager.TimeManager
	rand             rand_manager.RandManager
	logging          logger.Logger
}

func NewAuthService(
	cfg *config.ServiceConfig,
	repoUsers repository.Users,
	repoJWTokens repository.JWTokens,
	repoTransactions repository.Transactions,
	jwtManager jwtoken.JWTokenManager,
	hashManager hash.HashManager,
	timeManager time_manager.TimeManager,
	rand rand_manager.RandManager,
	logging logger.Logger,
) Auth {

	return &AuthService{
		cfg:              cfg,
		repoUsers:        repoUsers,
		repoJWTokens:     repoJWTokens,
		repoTransactions: repoTransactions,
		jwtManager:       jwtManager,
		hashManager:      hashManager,
		timeManager:      timeManager,
		rand:             rand,
		logging:          logging,
	}
}

const (
	alphabet       = "abcdefghijklmnopqrstuvwxyz1234567890"
	authMessage    = "Hello, %s! Please, sign this message with random param %s to sign in!"
	randStringSize = 20
)

func (s *AuthService) AuthByToken(
	ctx context.Context,
	user *domain.UserWithTokenNumber,
) (*models.AuthResponse, error) {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetUserByToken/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(context.Background())

	return s.getAuthRespByIdAndRole(ctx, tx, user.Id, user.Role)
}

func (s *AuthService) getAuthRespByIdAndRole(
	ctx context.Context,
	tx repository.Transaction,
	id int64,
	role domain.Role,
) (*models.AuthResponse, error) {
	var (
		resp models.AuthResponse
	)
	switch role {
	case domain.RoleClient:
		client, err := s.repoUsers.GetClientById(ctx, tx, id)
		if err != nil {
			return nil, newServiceError(code500,
				fmt.Errorf("getAuthRespByIdAndRole/GetClientById: %w", err), InternalError, "")
		}
		resp.User.Client = client.ToModel()
		resp.User.Role = domain.P(domain.RoleToModel(domain.RoleClient))
	}
	resp.ServerTime = s.timeManager.Now().UnixMilli()
	return &resp, nil
}

func (s *AuthService) GetUserByJWToken(
	ctx context.Context,
	purpose jwtoken.Purpose,
	token string,
) (*domain.UserWithTokenNumber, error) {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetUserByJWToken/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(context.Background())

	tokenData, err := s.jwtManager.ParseToken(token)
	if err != nil {
		return nil, newServiceError(code400,
			fmt.Errorf("GetUserByJWToken/ParseToken: %w", err), ParseTokenFailed, "")
	}
	if tokenData.Purpose != purpose {
		return nil, newServiceError(code400,
			fmt.Errorf("GetUserByJWToken: wrong token purpose"), ParseTokenFailed, "")
	}
	if tokenData.ExpiresAt.Before(s.timeManager.Now()) {
		return nil, newServiceError(code400,
			fmt.Errorf("GetUserByJWToken: token expired"), ParseTokenFailed, "")
	}

	correct, err := s.repoJWTokens.CheckJwt(ctx, tx, tokenData)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, newServiceError(code401,
				fmt.Errorf("GetUserByJWToken/CheckJwt: %w", err), ParseTokenFailed, "")
		}
		return nil, newServiceError(code500,
			fmt.Errorf("GetUGetUserByJWTokenserByToken/CheckJwt: %w", err), InternalError, "")
	}
	if !correct {
		return nil, newServiceError(code401,
			fmt.Errorf("GetUserByJWToken: %s", TokenWrongSecret), TokenWrongSecret, "")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetUserByJWToken/Commit: %w", err), InternalError, "")
	}

	return &domain.UserWithTokenNumber{
		Id:     tokenData.ID,
		Role:   domain.Role(tokenData.Role),
		Number: tokenData.Number,
	}, nil
}

func (s *AuthService) RefreshJWTokens(
	ctx context.Context,
	user *domain.UserWithTokenNumber,
) (*models.AuthResponse, *domain.PairOfTokens, error) {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return nil, nil, newServiceError(code500,
			fmt.Errorf("RefreshJWTokens/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(context.Background())

	if err := s.repoJWTokens.Drop(ctx, tx, user.Id, int(user.Role), user.Number); err != nil {
		return nil, nil, newServiceError(code500,
			fmt.Errorf("RefreshJWTokens/DropJWTokens: %w", err), InternalError, "")
	}
	accessToken, accessExpiresAt, refreshToken, refreshExpiresAt, err := s.generateAndSaveTokens(ctx, tx, user.Id, int(user.Role), user.Number)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.getAuthRespByIdAndRole(ctx, tx, user.Id, user.Role)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, newServiceError(code500,
			fmt.Errorf("RefreshJWTokens/Commit: %w", err), InternalError, "")
	}
	return resp,
		&domain.PairOfTokens{
			RefreshToken: &domain.Token{
				Token:     refreshToken,
				ExpiresAt: refreshExpiresAt,
			},
			AccessToken: &domain.Token{
				Token:     accessToken,
				ExpiresAt: accessExpiresAt,
			},
		}, nil
}

func (s *AuthService) Logout(ctx context.Context, user *domain.UserWithTokenNumber) error {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return newServiceError(code500,
			fmt.Errorf("Logout/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(context.Background())

	if err := s.repoJWTokens.Drop(ctx, tx, user.Id, int(user.Role), user.Number); err != nil {
		return newServiceError(code500,
			fmt.Errorf("Logout/Drop: %w", err), InternalError, "")
	}

	if err := tx.Commit(ctx); err != nil {
		return newServiceError(code500,
			fmt.Errorf("Logout/Commit: %w", err), InternalError, "")
	}
	return nil
}

func (s *AuthService) FullLogout(ctx context.Context, user *domain.UserWithTokenNumber) error {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return newServiceError(code500,
			fmt.Errorf("FullLogout/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(context.Background())

	if err := s.repoJWTokens.DropAll(ctx, tx, user.Id, int(user.Role)); err != nil {
		return newServiceError(code500,
			fmt.Errorf("FullLogout/DropAll: %w", err), InternalError, "")
	}

	if err := tx.Commit(ctx); err != nil {
		return newServiceError(code500,
			fmt.Errorf("FullLogout/Commit: %w", err), InternalError, "")
	}
	return nil
}

func (s *AuthService) GetAuthMessage(
	ctx context.Context,
	req *models.AuthMessageRequest,
) (*models.AuthMessageResponse, error) {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetAuthMessage/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(context.Background())

	msg, err := s.repoUsers.GetAuthMessageByAddress(ctx, tx, *req.Address)
	if err != nil && !errors.Is(err, repository.ErrNoRows) {
		return nil, newServiceError(code500,
			fmt.Errorf("GetAuthMessage/GetAuthMessageByAddress: %w", err), InternalError, "")
	}
	if msg != nil {
		if s.timeManager.Now().Sub(s.timeManager.MillisecondsToTime(msg.CreatedAt)) < 5*time.Minute {
			return &models.AuthMessageResponse{
				Message: &msg.Message,
			}, nil
		}
	}
	if err := s.repoUsers.DeleteAuthMessage(ctx, tx, *req.Address); err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetAuthMessage/DeleteAuthMessage: %w", err), InternalError, "")
	}
	message := fmt.Sprintf(authMessage, strings.ToLower(*req.Address), s.randomString(64))
	if err := s.repoUsers.InsertAuthMessage(ctx, tx, &domain.AuthMessage{
		Address:   strings.ToLower(*req.Address),
		Message:   message,
		CreatedAt: s.timeManager.Now().UnixMilli(),
	}); err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetAuthMessage/InsertAuthMessage: %w", err), InternalError, "")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, newServiceError(code500,
			fmt.Errorf("GetAuthMessage/Commit: %w", err), InternalError, "")
	}
	return &models.AuthMessageResponse{
		Message: &message,
	}, nil
}

func (s *AuthService) AuthByMessage(
	ctx context.Context,
	req *models.AuthBySignatureRequest,
) (*models.AuthResponse, *domain.PairOfTokens, error) {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return nil, nil, newServiceError(code500,
			fmt.Errorf("AuthByMessage/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(context.Background())

	msg, err := s.repoUsers.GetAuthMessageByAddress(ctx, tx, *req.Address)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, nil, newServiceError(code400,
				fmt.Errorf("AuthByMessage/GetAuthMessageByAddress: %w", err), AuthMessageNotExist, "")
		}
		return nil, nil, newServiceError(code500,
			fmt.Errorf("AuthByMessage/GetAuthMessageByAddress: %w", err), InternalError, "")
	}
	if s.timeManager.Now().Sub(s.timeManager.MillisecondsToTime(msg.CreatedAt)) >= 5*time.Minute {
		return nil, nil, newServiceError(code400,
			fmt.Errorf("AuthByMessage: %s", AuthMessageExpired), AuthMessageExpired, "")
	}

	hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg.Message), msg.Message)))
	sig := common.FromHex(*req.Signature)
	if len(sig) > 0 && sig[len(sig)-1] > 4 {
		sig[len(sig)-1] -= 27
	}
	pubKey, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		return nil, nil, newServiceError(code401,
			fmt.Errorf("AuthByMessage: %s", EcrecoverFailed), EcrecoverFailed, "")
	}
	pkey, err := crypto.UnmarshalPubkey(pubKey)
	if err != nil {
		return nil, nil, newServiceError(code401,
			fmt.Errorf("AuthByMessage/UnmarshalPubkey: %w", err), EcrecoverFailed, "")
	}
	signedAddress := crypto.PubkeyToAddress(*pkey)

	if !strings.EqualFold(signedAddress.Hex(), *req.Address) {
		s.logging.Infof("AuthByMessage: signedAddress: %s not equal address request: %s",
			signedAddress.Hex(), *req.Address)
		return nil, nil, newServiceError(code401,
			fmt.Errorf("AuthByMessage: %s", WrongSignature), WrongSignature, "")
	}

	userId, userRole, err := s.repoUsers.GetUserIdByAddress(ctx, tx, *req.Address)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			client, err2 := s.createClient(ctx, tx, *req.Address)
			if err2 != nil {
				return nil, nil, newServiceError(code500,
					fmt.Errorf("AuthByMessage/createClient: %w", err2), InternalError, "")
			}
			userId = client.Id
			userRole = domain.RoleClient
		} else {
			return nil, nil, newServiceError(code500,
				fmt.Errorf("AuthByMessage/GetUserIdByAddress: %w", err), InternalError, "")
		}
	}
	if userId <= 0 {
		return nil, nil, newServiceError(code400,
			fmt.Errorf("AuthByMessage/: user not exist"), UserNotExist, "")
	}

	number, err := s.repoJWTokens.GetNumber(ctx, tx, userId, int(userRole), jwtoken.PurposeAccess)
	if err != nil {
		return nil, nil, newServiceError(code500,
			fmt.Errorf("AuthByMessage/GetNumber: %w", err), InternalError, "")
	}

	accessToken, accessExpiresAt, refreshToken, refreshExpiresAt, err := s.generateAndSaveTokens(ctx, tx, userId, int(userRole), number)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.getAuthRespByIdAndRole(ctx, tx, userId, userRole)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, newServiceError(code500,
			fmt.Errorf("AuthByMessage/Commit: %w", err), InternalError, "")
	}

	return resp, &domain.PairOfTokens{
		RefreshToken: &domain.Token{
			Token:     refreshToken,
			ExpiresAt: refreshExpiresAt,
		},
		AccessToken: &domain.Token{
			Token:     accessToken,
			ExpiresAt: accessExpiresAt,
		},
	}, nil
}

func (s *AuthService) createClient(
	ctx context.Context,
	tx repository.Transaction,
	address string,
) (*domain.Client, error) {
	c := &domain.Client{
		Address:      common.HexToAddress(address),
		PointBalance: decimal.NewFromInt(0),
	}
	var err error
	c.Id, err = s.repoUsers.InsertClient(ctx, tx, c)
	if err != nil {
		return nil, fmt.Errorf("createUser/InsertUser: %w", err)
	}

	return c, nil
}

func (s *AuthService) generateAndSaveTokens(
	ctx context.Context,
	tx repository.Transaction,
	userId int64,
	role int,
	number int64,
) (string, time.Time, string, time.Time, error) {
	now := s.timeManager.Now()
	accessExpiresAt := now.Add(s.cfg.AccessTokenTTL)
	refreshExpiresAt := now.Add(s.cfg.RefreshTokenTTL)

	accessTokenData := &jwtoken.JWTokenData{
		Purpose:   jwtoken.PurposeAccess,
		Role:      role,
		ID:        userId,
		Number:    number,
		ExpiresAt: accessExpiresAt,
		Secret:    s.generateSecret(role, userId, number, jwtoken.PurposeAccess),
	}
	refreshTokenData := &jwtoken.JWTokenData{
		Purpose:   jwtoken.PurposeRefresh,
		Role:      role,
		ID:        userId,
		Number:    number,
		ExpiresAt: refreshExpiresAt,
		Secret:    s.generateSecret(role, userId, number, jwtoken.PurposeRefresh),
	}

	err := s.repoJWTokens.CreateJwt(ctx, tx, accessTokenData)
	if err != nil {
		return "", time.Time{}, "", time.Time{}, fmt.Errorf("generateJWTokens/CreateJwt: %w", err)
	}
	err = s.repoJWTokens.CreateJwt(ctx, tx, refreshTokenData)
	if err != nil {
		return "", time.Time{}, "", time.Time{}, fmt.Errorf("generateJWTokens/CreateJwt: %w", err)
	}

	accessToken, err := s.jwtManager.GenerateToken(accessTokenData)
	if err != nil {
		return "", time.Time{}, "", time.Time{}, fmt.Errorf("generateJWTokens/GenerateToken: %w", err)
	}
	refreshToken, err := s.jwtManager.GenerateToken(refreshTokenData)
	if err != nil {
		return "", time.Time{}, "", time.Time{}, fmt.Errorf("generateJWTokens/GenerateToken: %w", err)
	}

	return accessToken, accessExpiresAt, refreshToken, refreshExpiresAt, nil
}

func (s *AuthService) generateSecret(role int, id int64, number int64, purpose jwtoken.Purpose) string {
	toHashElems := []string{
		strconv.Itoa(role),
		fmt.Sprintf("%d", id),
		fmt.Sprintf("%d", number),
		strconv.Itoa(int(purpose)),
		s.randomString(randStringSize),
	}

	toHash := strings.Join(toHashElems, "_")
	hash := sha256.Sum256([]byte(toHash))

	return hex.EncodeToString(hash[:])
}

func (s *AuthService) randomString(l int) string {
	res := make([]byte, l)
	for i := 0; i < l; i++ {
		res[i] = alphabet[s.rand.Intn(len(alphabet))]
	}

	return string(res)
}

func (s *AuthService) DropExpiredTokens(ctx context.Context) error {
	tx, err := s.repoTransactions.BeginTransaction(ctx)
	if err != nil {
		return newServiceError(code500,
			fmt.Errorf("DropExpiredTokens/BeginTransaction: %w", err), InternalError, "")
	}
	defer tx.Rollback(context.Background())

	err = s.repoJWTokens.DropAllExpired(ctx, tx, jwtoken.PurposeAccess, s.timeManager.Now().UnixMilli())
	if err != nil {
		return newServiceError(code500,
			fmt.Errorf("DropExpiredTokens/DropAllExpired: %w", err), InternalError, "")
	}

	if err := tx.Commit(ctx); err != nil {
		return newServiceError(code500,
			fmt.Errorf("DropExpiredTokens/Commit: %w", err), InternalError, "")
	}

	return nil
}
