package convert

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type TimeTo time.Time

func (t TimeTo) Time() time.Time {
	return time.Time(t)
}

func (t TimeTo) ToPgTimeStamp() pgtype.Timestamp {
	return pgtype.Timestamp{Time: t.Time().UTC(), Valid: true}
}

func (t TimeTo) ToPgTimeStampZ() pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t.Time().UTC(), Valid: true}
}
