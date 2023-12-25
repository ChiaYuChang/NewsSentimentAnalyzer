package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
)

type Embeddings struct {
}

type DataFrame struct {
	Data    [][]float32
	rowAttr map[string][]any
	rowNum  int
	colAttr map[string][]any
	colNum  int
}

func (df DataFrame) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("DataFrame(%d, %d)", df.rowNum, df.colNum))
	for row := 0; row < df.rowNum; row++ {
		sb.WriteString(fmt.Sprintf("\n%10s", df.rowAttr["name"][row]))
		if df.HasRowAttr("sentiment") {
			sb.WriteString(fmt.Sprintf("%10s", df.rowAttr["sentiment"][row]))
		}
		if df.HasRowAttr("model") {
			sb.WriteString(fmt.Sprintf("%32s", df.rowAttr["model"][row]))
		}
		if df.HasRowAttr("nid") {
			sb.WriteString(fmt.Sprintf("%5d", df.rowAttr["nid"][row]))
		}

		col := 0
		for col = 0; col < df.colNum && col < 5; col++ {
			sb.WriteString(fmt.Sprintf("\t%5.4f", df.Data[row][col]))
		}

		if col < df.colNum {
			sb.WriteString(fmt.Sprintf("\t...(%d cols)", df.colNum-col))
		}
	}
	return sb.String()
}

func (df DataFrame) GetInstance(i int) []float32 {
	return df.Data[i]
}

func (df DataFrame) GetInstanceVecString(i int) string {
	s := fmt.Sprintf("%15.13f", df.GetInstance(i))
	return strings.ReplaceAll(s, " ", ", ")
}

func (df DataFrame) NRow() int {
	return df.rowNum
}

func (df DataFrame) NCol() int {
	return df.colNum
}

func (df DataFrame) Size() [2]int {
	return [2]int{df.rowNum, df.colNum}
}

func (df DataFrame) RowName(i int) any {
	if rn, ok := df.rowAttr["name"]; ok {
		return rn[i]
	}
	return ""
}

func (df DataFrame) ColName(i int) any {
	if cn, ok := df.colAttr["name"]; ok {
		return cn[i]
	}
	return ""
}

func (df DataFrame) HasRowAttr(attr string) bool {
	_, ok := df.rowAttr[attr]
	return ok
}

func (df DataFrame) HasColAttr(attr string) bool {
	_, ok := df.colAttr[attr]
	return ok
}

func (df DataFrame) RowAttr(i int, attr string) any {
	if attr, ok := df.rowAttr[attr]; ok {
		return attr[i]
	}
	return ""
}

func (df DataFrame) ColAttr(i int, attr string) any {
	if attr, ok := df.colAttr[attr]; ok {
		return attr[i]
	}
	return ""
}

func (df *DataFrame) SetRowAttr(attr string, values []any) (*DataFrame, error) {
	if len(values) != df.rowNum {
		return nil, fmt.Errorf("values length not match")
	}
	df.rowAttr[attr] = values
	return df, nil
}

func (df *DataFrame) SetColAttr(attr string, values []any) (*DataFrame, error) {
	if len(values) != df.colNum {
		return nil, fmt.Errorf("values length not match")
	}
	df.colAttr[attr] = values
	return df, nil
}

func (df DataFrame) Row(i int) DataFrameRow {
	ra := make(map[string]any)
	for k, v := range df.rowAttr {
		ra[k] = v[i]
	}

	return DataFrameRow{
		rowAttr: ra,
		data:    df.Data[i],
		colNum:  df.colNum,
		colAttr: df.colAttr,
	}
}

func (df DataFrame) Rows() []DataFrameRow {
	rows := make([]DataFrameRow, df.rowNum)
	for r := 0; r < df.rowNum; r++ {
		rows[r] = df.Row(r)
	}
	return rows
}

type DataFrameRow struct {
	rowAttr map[string]any
	data    []float32
	colNum  int
	colAttr map[string][]any
}

func (r DataFrameRow) VectorString() string {
	s := fmt.Sprintf("%15.13f", r.data)
	return strings.ReplaceAll(s, " ", ", ")
}

func (r DataFrameRow) RowName() any {
	return r.RowAttr("name")
}

func (r DataFrameRow) ColName(i int) any {
	return r.ColAttr(i, "name")
}

func (r DataFrameRow) RowAttr(attr string) any {
	return r.rowAttr[attr]
}

func (r DataFrameRow) ColAttr(i int, attr string) any {
	return r.colAttr[attr][i]
}

func NewDataFrameWithData(rowName, colName []any,
	data [][]float32) (*DataFrame, error) {
	rowNum := len(rowName)
	colNum := 0
	if rowNum > 0 {
		colNum = len(colName)
	}

	if rowNum != len(data) {
		return nil, fmt.Errorf("row number not match")
	}

	if colNum != len(data[0]) {
		return nil, fmt.Errorf("col number not match")
	}

	return &DataFrame{
		rowAttr: map[string][]any{"name": rowName},
		colAttr: map[string][]any{"name": colName},
		Data:    data,
		rowNum:  rowNum,
		colNum:  colNum,
	}, nil
}

func NewDataFrame(nrow, ncol int) *DataFrame {
	data := make([][]float32, nrow)
	for r := 0; r < nrow; r++ {
		data[r] = make([]float32, ncol)
	}

	return &DataFrame{
		rowAttr: map[string][]any{},
		colAttr: map[string][]any{},
		Data:    data,
		rowNum:  nrow,
		colNum:  ncol,
	}
}

func NewEmbedding(file string, dim int) (*DataFrame, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error while opening file: %v", err)
	}

	defer fd.Close()

	fileReader := csv.NewReader(fd)
	records, err := fileReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error while reading file: %v", err)
	}

	data := NewDataFrame(len(records), dim)
	rn := make([]any, len(records))
	sentiment := make([]any, len(records))

	for i, record := range records {
		rn[i] = record[0]
		if strings.Contains(record[0], "pos") {
			sentiment[i] = string(model.SentimentPositive)
		} else if strings.Contains(record[0], "neg") {
			sentiment[i] = string(model.SentimentNegative)
		} else {
			sentiment[i] = string(model.SentimentNeutral)
		}

		for j := 1; j < len(record); j++ {
			fmt.Sscanf(record[j], "%f", &data.Data[i][j-1])
		}
	}
	data.SetRowAttr("name", rn)
	data.SetRowAttr("sentiment", sentiment)
	return data, nil
}
