package service

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// ref: github.com/emicklei/pgtalk

func TimeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func TimeToTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t.UTC(), Valid: true}
}

func StringToText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: len(s) > 0}
}

func ToPgError(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr, true
	} else {
		return nil, false
	}
}

func ToPgTypeInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: i,
		Valid: true,
	}
}

func ToPgTypeInt2(i int16) pgtype.Int2 {
	return pgtype.Int2{
		Int16: i,
		Valid: true,
	}
}
