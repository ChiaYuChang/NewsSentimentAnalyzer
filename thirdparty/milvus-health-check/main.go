package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.StringP("host", "h", "localhost", "Milvus host")
	pflag.IntP("port", "p", 19530, "Milvus port")
	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	cli, err := client.NewGrpcClient(
		context.Background(),
		fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port")),
	)
	if err != nil {
		log.Fatal("failed to connect to Milvus:", err.Error())
		os.Exit(1)
	}

	err = cli.Close()
	if err != nil {
		log.Fatal("failed to close Milvus client:", err.Error())
		os.Exit(1)
	}

	log.Default().Println("Milvus server is healthy")
	os.Exit(0)
}
