package model

import (
	"context"
	"errors"
	"fmt"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrNewsHasAlreadyExist = errors.New("news already exist")
var ErrUpdateMoreThenOneRow = errors.New("update more then one row")

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn QueryCallBackFun) error
	DoCheckAndUpdateUserPasswordTx(ctx context.Context, params *CheckAndUpdateUserPasswordTxParams) error
	DoCreateOrUpdateAPIKeyTx(ctx context.Context, params *CreateOrUpdateAPIKeyTxParams) (*CreateOrUpdateAPIKeyTxResults, error)
	DoCountUserJobTx(ctx context.Context, owner uuid.UUID) (*CountUserJobTxResult, error)
	Close(ctx context.Context) error
}

type CheckAndUpdateUserPasswordTxParams struct {
	ID          int32  `json:"id"`
	Email       string `json:"email"`
	OldPassword []byte `json:"old_password"`
	NewPassword []byte `json:"new_password"`
}

type QueryCallBackFun func(*Queries) error

type CreateOrUpdateAPIKeyTxParams struct {
	Owner uuid.UUID `json:"owner"`
	ApiID int16     `json:"api_id"`
	Key   string    `json:"key"`
}

type CreateOrUpdateAPIKeyTxResults struct {
	ApiKeyId int32 `json:"apikey_id"`
	N        int64 `json:"n"`
}

type PGXStore struct {
	Querier
	Conn *pgx.Conn
}

func (s PGXStore) Close(ctx context.Context) error {
	return s.Conn.Close(ctx)
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

func (s PGXStore) DoCountUserJobTx(ctx context.Context, owner uuid.UUID) (*CountUserJobTxResult, error) {
	return countUserJobTx(s, ctx, owner)
}

type PGXPoolStore struct {
	Querier
	Conn *pgxpool.Pool
}

func (s PGXPoolStore) Close(ctx context.Context) error {
	s.Conn.Close()
	return nil
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

func (s PGXPoolStore) DoCountUserJobTx(ctx context.Context, owner uuid.UUID) (*CountUserJobTxResult, error) {
	return countUserJobTx(s, ctx, owner)
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

		var n int64
		n, err = q.UpdatePassword(ctx, &UpdatePasswordParams{
			ID:       auth.ID,
			Password: params.NewPassword,
		})

		if err != nil {
			return err
		}

		if n > 1 {
			return ErrUpdateMoreThenOneRow
		}
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

func createNewsTx(s Store, ctx context.Context, params *CreateNewsParams) (int64, error) {
	var newsId int64

	err := s.ExecTx(ctx, func(q *Queries) error {
		var err error
		_, err = q.GetNewsByMD5Hash(ctx, params.Md5Hash)

		if err == nil {
			return ErrNewsHasAlreadyExist
		}

		if err != pgx.ErrNoRows {
			return err
		}
		newsId, err = q.CreateNews(ctx, params)
		return err
	})
	return newsId, err
}

type CountUserJobTxResult struct {
	JobGroup  map[JobStatus]JobGroup `json:"job_group"`
	TotalJob  int                    `json:"total_job"`
	LastJobId int                    `json:"last_job_id"`
}

type JobGroup struct {
	Id   int `json:"jid"`
	NJob int `json:"n_job"`
}

func countUserJobTx(s Store, ctx context.Context, owner uuid.UUID) (*CountUserJobTxResult, error) {
	var countByGroup []*CountJobRow
	var lastJIdByGroup []*GetLastJobIdRow
	err := s.ExecTx(ctx, func(q *Queries) error {
		var err error
		if countByGroup, err = s.CountJob(ctx, owner); err != nil {
			return err
		}

		if lastJIdByGroup, err = s.GetLastJobId(ctx, owner); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	gm := map[JobStatus]int{}
	im := map[JobStatus]int{}
	for _, g := range countByGroup {
		gm[g.Status] = int(g.NJob)
	}

	for _, g := range lastJIdByGroup {
		im[g.Status] = int(g.ID)
	}

	result := &CountUserJobTxResult{
		JobGroup:  make(map[JobStatus]JobGroup, 5),
		TotalJob:  0,
		LastJobId: 0,
	}

	for _, key := range []JobStatus{JobStatusCreated, JobStatusRunning,
		JobStatusDone, JobStatusCanceled, JobStatusFailure} {
		njob, id := gm[key], im[key]

		result.JobGroup[key] = JobGroup{
			Id:   id,
			NJob: njob,
		}

		result.TotalJob += njob
		if result.LastJobId < id {
			result.LastJobId = id
		}
	}

	return result, nil
}
