package client

import (
	"context"
	"fmt"
	"net/http"
)

type Query interface {
	ToHTTPRequest(ctx context.Context) (*http.Request, error)
}

type Response interface {
	fmt.Stringer
	Status()
	NTotalResutl() int
}
