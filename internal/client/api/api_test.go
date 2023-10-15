package api_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	ou, err := url.Parse("https://gnews.io/api/v4/search?q=google&from=2020-01-01&to=202")
	require.NoError(t, err)

	b, err := json.Marshal(ou)
	require.NoError(t, err)
	t.Log(string(b))

	nu := &url.URL{}
	json.Unmarshal(b, nu)
	t.Log(nu)
}
