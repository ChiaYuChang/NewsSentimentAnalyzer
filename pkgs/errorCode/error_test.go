package errorcode_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	errCode := ec.Success
	status := http.StatusOK
	msg := "OK"

	err := ec.NewError(errCode, status, msg)
	require.NotNil(t, err)
	require.Equal(t, err.ErrorCode, errCode)
	require.Equal(t, err.HttpStatusCode, status)
	require.Equal(t, err.Message, msg)
	require.Equal(t,
		fmt.Sprintf("code: %d, status: %d, msg: %s", errCode, status, msg),
		err.Error(),
	)

	details := []string{"first line", "second line"}
	err.WithDetails(details...)
	require.Equal(t, err.Details, details)

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("code: %d, status: %d, msg: %s", errCode, status, msg))
	sb.WriteString(", details:\n")
	for i, d := range details {
		sb.WriteString(fmt.Sprintf("\t- %2d: %s\n", i, d))
	}
	require.Equal(t, sb.String(), err.Error())

	err.WithMessage("Update OK!")
	require.Equal(t, "Update OK!", err.Msg())

	err.WithMessagef("Hi, %s. The method works!")
	require.Equal(t, err.Msgf("Sam"), "Hi, Sam. The method works!")

	nerr := err.Clone()
	nerr.WithMessage("Clone err")
	nerr.Details[0] = "updated first line"
	require.NotEqual(t, err, nerr)
	require.NotEqual(t, err.Message, nerr.Message)
	require.NotEqual(t, err.Details, nerr.Details)
}
