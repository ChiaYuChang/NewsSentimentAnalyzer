package main

import (
	"encoding/xml"
	"fmt"
	"sync"
	"syscall/js"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/thirdparty/wasm/htmltable"
)

var document js.Value

type queryTable struct {
	sync.Once
	*htmltable.HTMLTable
	Raw string
	Err error
}

var QueryTable queryTable

func main() {
	c := make(chan bool)

	document = js.Global().Get("document")

	fmt.Println("Go Web Assembly is awesome")
	js.Global().Set("sort_table", js.FuncOf(SortTable))

	<-c
	fmt.Println("Done")
}

func SortTable(this js.Value, args []js.Value) any {
	tableEl := document.Call("getElementById", "qurey-result")
	QueryTable.Once.Do(func() {
		var tb *htmltable.HTMLTable
		QueryTable.Raw = tableEl.Get("outerHTML").String()
		QueryTable.Err = xml.Unmarshal([]byte(QueryTable.Raw), &tb)
		QueryTable.HTMLTable = tb
	})

	if QueryTable.Err != nil {
		return QueryTable.Err.Error()
	}

	QueryTable.HTMLTable.SortByWithMemory(args[0].Int())
	if b, err := xml.MarshalIndent(QueryTable.Body, "", "\t"); err != nil {
		return fmt.Sprintf("error while xml.Marshal: %w", err.Error())
	} else {
		tableEl.Get("lastElementChild").Call("replaceChildren")
		tableEl.Get("lastElementChild").Set("outerHTML", string(b))
		return nil
	}
}
