package global

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

var Secrets Secret

type Secret struct {
	Database map[string]Database `json:"db"`
	API      API                 `json:"api"`
}

type API struct {
	NewsSource    map[string]string `json:"news_source"`
	LanguageModel map[string]string `json:"language_model"`
}

func (api API) toString(indent string) string {
	sb := strings.Builder{}
	sb.WriteString(indent + "News Source:\n")
	for k, v := range api.NewsSource {
		sb.WriteString(fmt.Sprintf("%s - %-12s: %s\n", indent, k, v))
	}
	sb.WriteString(indent + "Language Model:\n")
	for k, v := range api.LanguageModel {
		sb.WriteString(fmt.Sprintf("%s - %-12s: %s\n", indent, k, v))
	}
	return sb.String()
}

func (api API) String() string {
	return api.toString("")
}

func (s Secret) String() string {
	sb := strings.Builder{}
	sb.WriteString("Secrets:\n")
	sb.WriteString("DataBase:\n")
	for dbName, db := range s.Database {
		sb.WriteString("    - " + dbName + "\n")
		sb.WriteString(db.toString("        "))
	}
	sb.WriteString("APIs:\n")
	sb.WriteString(s.API.toString("    "))
	return sb.String()
}

type Database struct {
	UserName string            `json:"username"`
	Password string            `json:"password"`
	DBName   string            `json:"db_name"`
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	Options  map[string]string `json:"options"`
}

func (db Database) toString(indent string) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%sUsername   : %s\n", indent, db.UserName))
	sb.WriteString(fmt.Sprintf("%sPassword   : %s\n", indent, db.Password))
	sb.WriteString(fmt.Sprintf("%sDB name    : %s\n", indent, db.DBName))
	sb.WriteString(indent + "Options    :\n")
	for optName, optValue := range db.Options {
		sb.WriteString(fmt.Sprintf("%s   - %-7s: %s\n", indent, optName, optValue))
	}
	return sb.String()
}

func (db Database) String() string {
	return db.toString("")
}

func ReadSecret(path string) (*Secret, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error while opening secret.json: %w", err)
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error while reading secrets: %w", err)
	}

	var secret Secret
	err = json.Unmarshal(bs, &secret)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshal secrets: %w", err)
	}

	return &secret, nil
}
