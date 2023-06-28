package model

import (
	"context"
	"errors"
	"fmt"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn QueryCallBackFun) error
	DoCheckAndUpdateUserPasswordTx(ctx context.Context, params *CheckAndUpdateUserPasswordTxParams) error
	DoCreateOrUpdateAPIKeyTx(ctx context.Context, params *CreateOrUpdateAPIKeyTxParams) (*CreateOrUpdateAPIKeyTxResults, error)
	Close() error
}

type CheckAndUpdateUserPasswordTxParams struct {
	ID          int32  `json:"id"`
	Email       string `json:"email"`
	OldPassword []byte `json:"old_password"`
	NewPassword []byte `json:"new_password"`
}

type QueryCallBackFun func(*Queries) error

type CreateOrUpdateAPIKeyTxParams struct {
	Owner int32  `json:"owner"`
	ApiID int16  `json:"api_id"`
	Key   string `json:"key"`
}

type CreateOrUpdateAPIKeyTxResults struct {
	ApiKeyId int32 `json:"apikey_id"`
	N        int64 `json:"n"`
}

type PGXStore struct {
	Querier
	Conn *pgx.Conn
}

func (s PGXStore) Close() error {
	return s.Close()
}

func NewDBConnection(ctx context.Context, connStr string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(
		context.Background(),
		connStr)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewPGXStore(conn *pgx.Conn) Store {
	return PGXStore{New(conn), conn}
}

func (s PGXStore) ExecTx(ctx context.Context, fn QueryCallBackFun) error {
	tx, err := s.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err = fn(New(tx)); err != nil {
		// rollback if queries failed
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %w, rollback err: %w", err, rbErr)
		}
		return err
	}

	// commit transaction
	return tx.Commit(ctx)
}

func (s PGXStore) DoCheckAndUpdateUserPasswordTx(ctx context.Context, params *CheckAndUpdateUserPasswordTxParams) error {
	return checkAndUpdateUserPasswordTx(s, ctx, params)
}

func (s PGXStore) DoCreateOrUpdateAPIKeyTx(ctx context.Context, params *CreateOrUpdateAPIKeyTxParams) (*CreateOrUpdateAPIKeyTxResults, error) {
	return createOrUpdateAPIKeyTx(s, ctx, params)
}

type PGXPoolStore struct {
	Querier
	Conn *pgxpool.Pool
}

func (s PGXPoolStore) Close() error {
	return s.Close()
}

func NewDBConnectionPools(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("Unable to create connection pool: %w", err)
	}

	return dbPool, nil
}

func NewPGXPoolStore(pool *pgxpool.Pool) Store {
	return PGXPoolStore{New(pool), pool}
}

func (s PGXPoolStore) ExecTx(ctx context.Context, fn QueryCallBackFun) error {
	tx, err := s.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err = fn(New(tx)); err != nil {
		// rollback if queries failed
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %w, rollback err: %w", err, rbErr)
		}
		return err
	}

	// commit transaction
	return tx.Commit(ctx)
}

func (s PGXPoolStore) DoCheckAndUpdateUserPasswordTx(ctx context.Context, params *CheckAndUpdateUserPasswordTxParams) error {
	return checkAndUpdateUserPasswordTx(s, ctx, params)
}

func (s PGXPoolStore) DoCreateOrUpdateAPIKeyTx(ctx context.Context, params *CreateOrUpdateAPIKeyTxParams) (*CreateOrUpdateAPIKeyTxResults, error) {
	return createOrUpdateAPIKeyTx(s, ctx, params)
}

func checkAndUpdateUserPasswordTx(s Store, ctx context.Context, params *CheckAndUpdateUserPasswordTxParams) error {
	err := s.ExecTx(ctx, func(q *Queries) error {
		auth, err := q.GetUserAuth(ctx, params.Email)
		if err != nil {
			return err
		}

		if err = bcrypt.CompareHashAndPassword(auth.Password, params.OldPassword); err != nil {
			ecErr := ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
			ecErr.WithDetails(err.Error())
			return ecErr
		}

		_, err = q.UpdatePassword(ctx, &UpdatePasswordParams{
			ID:       auth.ID,
			Password: params.NewPassword,
		})
		return err
	})
	return err
}

func createOrUpdateAPIKeyTx(s Store, ctx context.Context, params *CreateOrUpdateAPIKeyTxParams) (*CreateOrUpdateAPIKeyTxResults, error) {
	result := &CreateOrUpdateAPIKeyTxResults{}
	err := s.ExecTx(ctx, func(q *Queries) error {
		var err error
		var apikey *GetAPIKeyRow
		apikey, err = q.GetAPIKey(ctx, &GetAPIKeyParams{
			Owner: params.Owner,
			ApiID: params.ApiID,
		})

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				result.ApiKeyId, err = q.CreateAPIKey(ctx, &CreateAPIKeyParams{
					Owner: params.Owner,
					ApiID: params.ApiID,
					Key:   params.Key,
				})
				result.N = 1
			}
		} else {
			result.ApiKeyId = apikey.ID
			result.N, err = q.UpdateAPIKey(ctx, &UpdateAPIKeyParams{
				Owner:    params.Owner,
				Key:      params.Key,
				OldApiID: params.ApiID,
				NewApiID: params.ApiID,
			})
		}
		return err
	})
	return result, err
}
