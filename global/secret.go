package global

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var Secrets Secret

type Secret struct {
	Database []Database `json:"db"`
	API      []APIKey   `json:"api"`
}

func (s Secret) String() string {
	jsm, _ := json.MarshalIndent(s, "", "\t")
	return string(jsm)
}

type Database struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type APIType int8

const (
	NewsSource APIType = iota + 1
	LanguageModel
)

func (t APIType) String() string {
	switch t {
	case NewsSource:
		return "NewsSource"
	case LanguageModel:
		return "LanguageModel"
	}
	panic("unknown API type")
}

func (t *APIType) UnmarshalJSON(data []byte) error {
	switch string(data[1 : len(data)-1]) {
	case "NewsSource", "src", "source":
		*t = NewsSource
	case "LanguageModel", "lm", "language model":
		*t = LanguageModel
	default:
		return fmt.Errorf("unknown api type %s", string(data))
	}
	return nil
}

func (t *APIType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.String())), nil
}

type APIKey struct {
	Name string  `json:"name"`
	Type APIType `json:"type"`
	Key  string  `json:"key"`
}

func ReadSecret(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error while opening secret.json: %w", err)
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error while reading secrets: %w", err)
	}

	var secret Secret
	err = json.Unmarshal(bs, &secret)
	if err != nil {
		return fmt.Errorf("error while unmarshal secrets: %w", err)
	}
	AppVar.Secret = secret
	return nil
}
