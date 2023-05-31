package service

import (
	"time"

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
