package main

import (
	"fmt"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
)

func main() {
	if err := global.ReadSecret("./secret.json"); err != nil {
		panic(fmt.Sprintf("error while reading secret: %s", err.Error()))
	}
	fmt.Println(global.AppVar.Secret)
}
