package model

import (
	"context"
	"errors"
	"fmt"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	DoCacheToStoreTx(ctx context.Context, params *CacheToStoreTXParams) (*CacheToStoreTXResult, error)
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

func (s PGXStore) DoCacheToStoreTx(ctx context.Context, params *CacheToStoreTXParams) (*CacheToStoreTXResult, error) {
	return doCacheToStoreTx(s, ctx, params)
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

func (s PGXPoolStore) DoCacheToStoreTx(ctx context.Context, params *CacheToStoreTXParams) (*CacheToStoreTXResult, error) {
	return doCacheToStoreTx(s, ctx, params)
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

type CacheToStoreTXParams struct {
	CreateJobParams  *CreateJobParams
	CreateNewsParams <-chan *CreateNewsParams
}

type CacheToStoreTXResult struct {
	JobId                int64
	JobCreateError       *ec.Error
	NewsJobCreateResults []NewsJobCreateResult
}

func (r *CacheToStoreTXResult) Error() error {
	if err := r.JobCreateError; err != nil {
		return r.JobCreateError
	}
	for _, r := range r.NewsJobCreateResults {
		if err := r.Error; err != nil {
			return err
		}
	}
	return nil
}

type NewsJobCreateResult struct {
	Md5Hash   string    `json:"md5_hash"`
	NewsID    int64     `json:"news_id"`
	NewsJobId int64     `json:"news_job_id"`
	Error     *ec.Error `json:"error"`
}

func cacheToStoreTx(s Store, ctx context.Context, params *CacheToStoreTXParams) *CacheToStoreTXResult {
	var err error

	result := &CacheToStoreTXResult{}
	result.JobId, err = s.CreateJob(ctx, params.CreateJobParams)
	if err != nil {
		var pgErr *pgconn.PgError
		var ecErr *ec.Error
		if errors.As(err, &pgErr) {
			ecErr = ec.NewErrorFromPgErr(pgErr)
		} else {
			ecErr = ec.MustGetEcErr(ec.ECServerError).
				WithMessage(err.Error())
		}
		ecErr.WithDetails("error while created job")
		result.JobCreateError = ecErr
		return result
	}

	result.NewsJobCreateResults = []NewsJobCreateResult{}
	for param := range params.CreateNewsParams {
		row, err := s.GetNewsByMD5Hash(ctx, param.Md5Hash)

		r := NewsJobCreateResult{}
		r.Md5Hash = param.Md5Hash
		if err == nil {
			r.NewsID = row.ID
		} else {
			nid, err := s.CreateNews(ctx, param)
			if err != nil {
				var pgErr *pgconn.PgError
				var ecErr *ec.Error
				if errors.As(err, &pgErr) {
					ecErr = ec.NewErrorFromPgErr(pgErr)
				} else {
					ecErr = ec.MustGetEcErr(ec.ECServerError).
						WithMessage(err.Error())
				}
				ecErr.WithDetails("error while CreateNews")
				r.Error = ecErr
			}
			r.NewsID = nid
		}
		result.NewsJobCreateResults = append(result.NewsJobCreateResults, r)
	}

	for i, r := range result.NewsJobCreateResults {
		if r.Error == nil {
			if njId, err := s.CreateNewsJob(ctx, &CreateNewsJobParams{
				NewsID: r.NewsID, JobID: int64(result.JobId),
			}); err != nil {
				var pgErr *pgconn.PgError
				var ecErr *ec.Error
				if errors.As(err, &pgErr) {
					ecErr = ec.NewErrorFromPgErr(pgErr)
				} else {
					ecErr = ec.MustGetEcErr(ec.ECServerError).
						WithMessage(err.Error())
				}
				ecErr.WithDetails("error while CreateNewsJob")
				result.NewsJobCreateResults[i].Error = ecErr
			} else {
				result.NewsJobCreateResults[i].NewsJobId = njId
			}
		}
	}
	return result
}

func doCacheToStoreTx(s Store, ctx context.Context, params *CacheToStoreTXParams) (*CacheToStoreTXResult, error) {
	var result *CacheToStoreTXResult
	err := s.ExecTx(ctx, func(q *Queries) error {
		result = cacheToStoreTx(s, ctx, params)
		return result.Error()
	})
	return result, err
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
		JobStatusDone, JobStatusCanceled, JobStatusFailed} {
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
