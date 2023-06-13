package global

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Endpoint struct {
	Name         string            `json:"name"`
	DocumentURL  string            `json:"document_url"`
	Image        []string          `json:"image"`
	TemplateName map[string]string `json:"templates"`
}

func ReadEndpoints(path string) ([]Endpoint, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error while opening option.json: %w", err)
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error while reading token maker option: %w", err)
	}

	var eps = []Endpoint{}
	err = json.Unmarshal(bs, &eps)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshal option: %w", err)
	}
	return eps, nil
}

func (ep Endpoint) ToString(indent string) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%sDocument URL: %s\n", indent, ep.DocumentURL))

	sb.WriteString(fmt.Sprintf("%sImage path:\n", indent))
	for _, img := range ep.Image {
		sb.WriteString(fmt.Sprintf("%s     - %s\n", indent, img))
	}

	for epName, tpName := range ep.TemplateName {
		sb.WriteString(fmt.Sprintf("%s     - %s: %s\n", indent, epName, tpName))
	}
	return sb.String()
}

func (ep Endpoint) String() string {
	return ep.ToString("")
}
