// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	pgv "github.com/pgvector/pgvector-go"
)

type ApiType string

const (
	ApiTypeLanguageModel ApiType = "language_model"
	ApiTypeSource        ApiType = "source"
)

func (e *ApiType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ApiType(s)
	case string:
		*e = ApiType(s)
	default:
		return fmt.Errorf("unsupported scan type for ApiType: %T", src)
	}
	return nil
}

type NullApiType struct {
	ApiType ApiType `json:"api_type"`
	Valid   bool    `json:"valid"` // Valid is true if ApiType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullApiType) Scan(value interface{}) error {
	if value == nil {
		ns.ApiType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ApiType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullApiType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ApiType), nil
}

type EventType string

const (
	EventTypeSignIn        EventType = "sign-in"
	EventTypeSignOut       EventType = "sign-out"
	EventTypeAuthorization EventType = "authorization"
	EventTypeApiKey        EventType = "api-key"
	EventTypeQuery         EventType = "query"
)

func (e *EventType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EventType(s)
	case string:
		*e = EventType(s)
	default:
		return fmt.Errorf("unsupported scan type for EventType: %T", src)
	}
	return nil
}

type NullEventType struct {
	EventType EventType `json:"event_type"`
	Valid     bool      `json:"valid"` // Valid is true if EventType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEventType) Scan(value interface{}) error {
	if value == nil {
		ns.EventType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EventType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEventType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EventType), nil
}

type JobStatus string

const (
	JobStatusCreated  JobStatus = "created"
	JobStatusRunning  JobStatus = "running"
	JobStatusDone     JobStatus = "done"
	JobStatusFailed   JobStatus = "failed"
	JobStatusCanceled JobStatus = "canceled"
)

func (e *JobStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = JobStatus(s)
	case string:
		*e = JobStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for JobStatus: %T", src)
	}
	return nil
}

type NullJobStatus struct {
	JobStatus JobStatus `json:"job_status"`
	Valid     bool      `json:"valid"` // Valid is true if JobStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullJobStatus) Scan(value interface{}) error {
	if value == nil {
		ns.JobStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.JobStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullJobStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.JobStatus), nil
}

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (e *Role) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Role(s)
	case string:
		*e = Role(s)
	default:
		return fmt.Errorf("unsupported scan type for Role: %T", src)
	}
	return nil
}

type NullRole struct {
	Role  Role `json:"role"`
	Valid bool `json:"valid"` // Valid is true if Role is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRole) Scan(value interface{}) error {
	if value == nil {
		ns.Role, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Role.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Role), nil
}

type Sentiment string

const (
	SentimentPositive Sentiment = "positive"
	SentimentNeutral  Sentiment = "neutral"
	SentimentNegative Sentiment = "negative"
)

func (e *Sentiment) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Sentiment(s)
	case string:
		*e = Sentiment(s)
	default:
		return fmt.Errorf("unsupported scan type for Sentiment: %T", src)
	}
	return nil
}

type NullSentiment struct {
	Sentiment Sentiment `json:"sentiment"`
	Valid     bool      `json:"valid"` // Valid is true if Sentiment is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullSentiment) Scan(value interface{}) error {
	if value == nil {
		ns.Sentiment, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Sentiment.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullSentiment) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Sentiment), nil
}

type Api struct {
	ID          int16              `json:"id"`
	Name        string             `json:"name"`
	Type        ApiType            `json:"type"`
	Image       string             `json:"image"`
	Icon        string             `json:"icon"`
	DocumentUrl string             `json:"document_url"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	DeletedAt   pgtype.Timestamptz `json:"deleted_at"`
}

type Apikey struct {
	ID        int32              `json:"id"`
	Owner     uuid.UUID          `json:"owner"`
	ApiID     int16              `json:"api_id"`
	Key       string             `json:"key"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

type Embedding struct {
	ID        int64              `json:"id"`
	Model     string             `json:"model"`
	NewsID    int64              `json:"news_id"`
	Embedding pgv.Vector         `json:"embedding"`
	Sentiment Sentiment          `json:"sentiment"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

type Endpoint struct {
	ID           int32              `json:"id"`
	Name         string             `json:"name"`
	ApiID        int16              `json:"api_id"`
	TemplateName string             `json:"template_name"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	DeletedAt    pgtype.Timestamptz `json:"deleted_at"`
}

type Job struct {
	ID        int64              `json:"id"`
	Ulid      string             `json:"ulid"`
	Owner     uuid.UUID          `json:"owner"`
	Status    JobStatus          `json:"status"`
	SrcApiID  int16              `json:"src_api_id"`
	SrcQuery  string             `json:"src_query"`
	LlmApiID  int16              `json:"llm_api_id"`
	LlmQuery  []byte             `json:"llm_query"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

type Keyword struct {
	ID      int64  `json:"id"`
	NewsID  int64  `json:"news_id"`
	Keyword string `json:"keyword"`
}

type Log struct {
	ID        int64              `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	Type      EventType          `json:"type"`
	Message   string             `json:"message"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type News struct {
	ID          int64              `json:"id"`
	Md5Hash     string             `json:"md5_hash"`
	Guid        string             `json:"guid"`
	Author      []string           `json:"author"`
	Title       string             `json:"title"`
	Link        string             `json:"link"`
	Description string             `json:"description"`
	Language    pgtype.Text        `json:"language"`
	Content     []string           `json:"content"`
	Category    string             `json:"category"`
	Source      string             `json:"source"`
	RelatedGuid []string           `json:"related_guid"`
	PublishAt   pgtype.Timestamptz `json:"publish_at"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
}

type Newsjob struct {
	ID     int64 `json:"id"`
	JobID  int64 `json:"job_id"`
	NewsID int64 `json:"news_id"`
}

type SchemaMigration struct {
	Version int64 `json:"version"`
	Dirty   bool  `json:"dirty"`
}

type User struct {
	ID                uuid.UUID          `json:"id"`
	Password          []byte             `json:"password"`
	FirstName         string             `json:"first_name"`
	LastName          string             `json:"last_name"`
	Role              Role               `json:"role"`
	Email             string             `json:"email"`
	Opt               pgtype.Text        `json:"opt"`
	CreatedAt         pgtype.Timestamptz `json:"created_at"`
	UpdatedAt         pgtype.Timestamptz `json:"updated_at"`
	DeletedAt         pgtype.Timestamptz `json:"deleted_at"`
	PasswordUpdatedAt pgtype.Timestamptz `json:"password_updated_at"`
}
