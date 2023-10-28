package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/spf13/pflag"
)

type Config struct {
	API      []APIItem      `json:"api"`
	n2i      map[string]int `json:"-"` // convert api name to api index
	Endpoint []EndpointItem `json:"endpoint"`
	User     struct {
		SpecialUser []UserItem `json:"special_user"`
		N           int        `json:"n"`
	}
	Job struct {
		SpecialJob    []JobItem `json:"special_job"`
		MaxJobPerUser int       `json:"max_job_per_user"`
	}
}

func (c Config) String() string {
	sb := strings.Builder{}
	sb.WriteString("Config:\n")
	b, _ := json.MarshalIndent(c, "  ", "    ")
	sb.Write(b)
	return sb.String()
}

func main() {
	var configPath string

	pflag.StringVarP(&configPath, "config", "c", "./config.json", "path to config file")

	f, err := os.OpenFile(configPath, os.O_RDONLY, 644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		log.Fatal(err)
	}

	mod := modifiers.New()
	for i := range config.API {
		if err := mod.Struct(context.Background(), &config.API[i]); err != nil {
			log.Fatal(err)
		}
	}

	for i := range config.Job.SpecialJob {
		if err := mod.Struct(context.Background(), &config.Job.SpecialJob[i]); err != nil {
			log.Fatal(err)
		}
	}

	for i := range config.Endpoint {
		if err := mod.Struct(context.Background(), &config.Endpoint[i]); err != nil {
			log.Fatal(err)
		}
	}

	for i := range config.API {
		if err := mod.Struct(context.Background(), &config.API[i]); err != nil {
			log.Fatal(err)
		}
	}

	config.n2i = make(map[string]int, len(config.API))
	for i := range config.API {
		config.API[i].Id = i + 1
		config.n2i[config.API[i].Name] = config.API[i].Id
	}

	for i := range config.Endpoint {
		config.Endpoint[i].APIId = config.n2i[config.Endpoint[i].APIName]
	}

	for i := range config.Job.SpecialJob {
		rts := rg.GenRdnTimes(2, TIME_MIN, TIME_MAX)
		config.Job.SpecialJob[i].Id = i + 1
		config.Job.SpecialJob[i].CreatedAt = rts[0]
		config.Job.SpecialJob[i].UpdatedAt = rts[1]
		config.Job.SpecialJob[i].SrcApiId = config.n2i[config.Job.SpecialJob[i].SrcAPIName]
		config.Job.SpecialJob[i].LlmApiId = config.n2i[config.Job.SpecialJob[i].LlmAPIName]
	}

	log.Println(config)
}
