package middleware_test

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/ratelimit"
)

func TestUberRateLimiter(t *testing.T) {
	// add one token into bucket per second
	rl := middleware.NewRateLimiter(1, ratelimit.Per(1*time.Second))

	mux := chi.NewRouter()
	mux.Use(rl.RateLimit)
	mux.Get("/time", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	})

	srv := httptest.NewServer(mux)
	url := srv.URL

	n := 5
	startTime := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(t *testing.T, i int, url string, wg *sync.WaitGroup) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Int63n(30)) * time.Millisecond)
			resp, err := http.Get(fmt.Sprintf("%s/time", url))
			require.NoError(t, err)
			_, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
			t.Logf("Case %d ...done!\n", i)
		}(t, i, url, wg)
	}
	wg.Wait()
	endTime := time.Now()
	require.True(t, endTime.Sub(startTime).Seconds() > float64(n-1))
	require.True(t, endTime.Sub(startTime).Seconds() < float64(n))
}

// func TestRedisRateLimiter(t *testing.T) {
// 	username := "username"
// 	role := tokenmaker.RUser
// 	uid := int32(1)

// 	maker := middleware.NewJWTTokenMaker(opt)
// 	maker.AllowFromHTTPCookie = true
// 	bearer, err := maker.TokenMaker.MakeToken(username, uid, role)
// 	require.NoError(t, err)

// 	_, err = maker.TokenMaker.ValidateToken(bearer)
// 	require.NoError(t, err)

// 	cli := redis.NewClient(&redis.Options{
// 		Network:    "tcp",
// 		Addr:       "localhost:6379",
// 		MaxRetries: 3,
// 		PoolSize:   1,
// 	})
// 	require.NoError(t, cli.Ping(context.Background()).Err())

// 	// add two token into bucket per second
// 	rl := middleware.NewRedisRateLimiter(cli, 2, 1*time.Second)
// 	require.NotNil(t, rl)

// 	mux := chi.NewRouter()
// 	mux.Use(maker.BearerAuthenticator)
// 	mux.Use(rl.RateLimit)
// 	mux.Get("/time", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(time.Now().Format(time.RFC3339Nano)))
// 	})

// 	srv := httptest.NewServer(mux)
// 	url := srv.URL
// 	n := 5

// 	t.Run(
// 		"OK",
// 		func(t *testing.T) {
// 			wg := &sync.WaitGroup{}
// 			for i := 0; i < n; i++ {
// 				wg.Add(1)
// 				time.Sleep(500 * time.Millisecond)
// 				go func(t *testing.T, i int, url string, wg *sync.WaitGroup) {
// 					defer wg.Done()
// 					req, err := http.NewRequest(
// 						http.MethodGet,
// 						fmt.Sprintf("%s/time", url), nil)
// 					require.NoError(t, err)
// 					require.NotNil(t, req)
// 					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))

// 					resp, err := http.DefaultClient.Do(req)
// 					require.NoError(t, err)
// 					require.Equal(t, http.StatusOK, resp.StatusCode)

// 					_, err = io.ReadAll(resp.Body)
// 					require.NoError(t, err)
// 					t.Logf("Case %d ...done!\n", i)
// 				}(t, i, url, wg)
// 			}
// 			wg.Wait()
// 		},
// 	)

// 	time.Sleep(1000 * time.Millisecond)
// 	t.Run(
// 		"Too many queries per second",
// 		func(t *testing.T) {
// 			for i := 0; i < n; i++ {
// 				req, err := http.NewRequest(
// 					http.MethodGet,
// 					fmt.Sprintf("%s/time", url), nil)
// 				require.NoError(t, err)
// 				require.NotNil(t, req)
// 				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))

// 				resp, err := http.DefaultClient.Do(req)
// 				require.NoError(t, err)
// 				if i > 1 {
// 					require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
// 				} else {
// 					require.Equal(t, http.StatusOK, resp.StatusCode)
// 				}
// 				_, err = io.ReadAll(resp.Body)
// 				require.NoError(t, err)
// 				t.Logf("Case %d ...done!\n", i)
// 			}
// 		},
// 	)
// }
