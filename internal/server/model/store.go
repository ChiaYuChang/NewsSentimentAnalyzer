package model

import (
	"context"
	"fmt"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	Querier
	DoCheckAndUpdateUserPasswordTx(ctx context.Context, params *CheckAndUpdateUserPasswordTxParams) error
	Close() error
}

type PGXStore struct {
	Querier
	Conn *pgx.Conn
}

func (s PGXStore) Close() error {
	return s.Close()
}

type PGXPoolStore struct {
	Querier
	Conn *pgxpool.Pool
}

func (s PGXPoolStore) Close() error {
	return s.Close()
}

func NewDBConnection(ctx context.Context, config *pgx.ConnConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(
		context.Background(),
		config.ConnString())

	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}
	return conn, nil
}

func NewPGXStore(conn *pgx.Conn) Store {
	return PGXStore{New(conn), conn}
}

func NewDBConnectionPools(ctx context.Context, config *pgx.ConnConfig) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		return nil, fmt.Errorf("Unable to create connection pool: %w", err)
	}

	return dbPool, nil
}

// func NewPGXPoolStore(pool *pgxpool.Pool) Store {
// 	return PGXPoolStore{New(pool), pool}
// }

type CheckAndUpdateUserPasswordTxParams struct {
	ID          int32  `json:"id"`
	Email       string `json:"email"`
	OldPassword []byte `json:"old_password"`
	NewPassword []byte `json:"new_password"`
}

type QueryCallBackFun func(*Queries) error

func (s PGXStore) exectTx(ctx context.Context, fn QueryCallBackFun) error {
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
	//TODO create an event record here
	return s.CheckAndUpdateUserPasswordTx(ctx, params)
}

func (s PGXStore) CheckAndUpdateUserPasswordTx(
	ctx context.Context, params *CheckAndUpdateUserPasswordTxParams) error {
	err := s.exectTx(ctx, func(q *Queries) error {
		auth, err := s.Querier.GetUserAuth(ctx, params.Email)
		if err != nil {
			return err
		}

		if err = bcrypt.CompareHashAndPassword(auth.Password, params.OldPassword); err != nil {
			ecErr := ec.MustGetErr(ec.ECUnauthorized).(*ec.Error)
			ecErr.WithDetails(err.Error())
			return ecErr
		}

		_, err = s.Querier.UpdatePassword(ctx, &UpdatePasswordParams{
			ID:       auth.ID,
			Password: params.NewPassword,
		})
		return err
	})
	return err
}
