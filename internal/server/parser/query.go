package parser

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/PuerkitoBio/goquery"
)

// Querier is a struct that contains query options.
type Querier struct {
	header  map[string]string
	client  *http.Client
	handler map[string]ContentEncodeHandler
}

func NewQuerier(opts ...QuerierOpt) (*Querier, error) {
	q := &Querier{
		header:  map[string]string{},
		client:  &http.Client{},
		handler: map[string]ContentEncodeHandler{},
	}

	var err error
	for _, opt := range opts {
		if q, err = opt(q); err != nil {
			return nil, err
		}
	}
	return q, nil
}

// HasContentEncodingHandler returns true if the querier has a handler for a content encoding.
func (q *Querier) HasContentEncodingHandler(ct string) bool {
	_, ok := q.handler[ct]
	return ok
}

// ContentEncodeHandler is a function for decoding content.
type ContentEncodeHandler func(r io.ReadCloser) (io.ReadCloser, error)

// DefaultTypeHandler is a default handler for content encoding.
// It simply returns the original io.ReadCloser.
func DefaultTypeHandler(rc io.ReadCloser) (io.ReadCloser, error) {
	return rc, nil
}

// QuerierOpt is a function that modifies a Querier.
type QuerierOpt func(q *Querier) (*Querier, error)

type QuerierStep func(*Query) *Query

// WithDefaultHeader is a QuerierOpt that sets default header.
// It sets the following header:
// For all requests:
// - User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/
// - Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8
// - Accept-Encoding: gzip
// For responses:
// - Content-Encoding: gzip
// - Content-Encoding Handler: gzip.NewReader
// For client:
// - Client: http.DefaultClient
// - Client Timeout: 10 seconds
// It returns an error if any of the above fails.
func WithDefaultHeader() QuerierOpt {
	return func(q *Querier) (*Querier, error) {
		var err error
		for _, opt := range []QuerierOpt{
			WithUserAgentLinuxChrome(),
			WithAccept("text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8"),
			WithGzipCompression(),
		} {
			if q, err = opt(q); err != nil {
				return nil, err
			}
		}
		return q, nil
	}
}

// WithTestHeader is a QuerierOpt that sets test header.
// It sets the following header:
// For all requests:
// - User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/
// - Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8
func WithTestHeader() QuerierOpt {
	return func(q *Querier) (*Querier, error) {
		var err error
		for _, opt := range []QuerierOpt{
			WithUserAgentLinuxChrome(),
			WithAccept("text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8"),
		} {
			if q, err = opt(q); err != nil {
				return nil, err
			}
		}
		return q, nil
	}
}

// WithUserAgent is a QuerierOpt that sets User-Agent header.
func WithUserAgent(v string) QuerierOpt {
	return func(q *Querier) (*Querier, error) {
		q.header["User-Agent"] = v
		return q, nil
	}
}

// WithUserAgentLinuxChrome is a QuerierOpt that sets User-Agent header to a Linux Chrome:
// Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/
func WithUserAgentLinuxChrome() QuerierOpt {
	return WithUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
}

// WithHeader is a QuerierOpt that sets a header.
func WithHeader(key, val string) QuerierOpt {
	return func(q *Querier) (*Querier, error) {
		q.header[key] = val
		return q, nil
	}
}

// WithAccept is a QuerierOpt that sets Accept header.
func WithAccept(val string) QuerierOpt {
	return WithHeader("Accept", val)
}

// WithAcceptEncoding is a QuerierOpt that sets Accept-Encoding header.
func WithAcceptEncoding(val string) QuerierOpt {
	return WithHeader("Accept-Encoding", val)
}

// WithContentEncodingHandler is a QuerierOpt that sets a handler for a content encoding.
func WithContentEncodingHandler(ce string, h ContentEncodeHandler) QuerierOpt {
	return func(q *Querier) (*Querier, error) {
		q.handler[ce] = h
		return q, nil
	}
}

// WithCompression is a QuerierOpt that sets Accept-Encoding header and a handler for a content encoding.
func WithCompression(tag string, handler ContentEncodeHandler) QuerierOpt {
	return func(q *Querier) (*Querier, error) {
		var err error

		for _, opt := range []QuerierOpt{
			WithAcceptEncoding(tag),
			WithContentEncodingHandler(tag, handler),
		} {
			q, err = opt(q)
			if err != nil {
				return q, err
			}
		}
		return q, nil
	}
}

// WithGzipCompression is a QuerierOpt that sets Accept-Encoding header to gzip and a handler for gzip.
func WithGzipCompression() QuerierOpt {
	return WithCompression("gzip", func(r io.ReadCloser) (io.ReadCloser, error) {
		return gzip.NewReader(r)
	})
}

// WithClient is a QuerierOpt that sets a client.
func WithClient(cli *http.Client) QuerierOpt {
	return func(q *Querier) (*Querier, error) {
		q.client = cli
		return q, nil
	}
}

// WithDefaultClient is a QuerierOpt that sets http.DefaultClient as client.
func WithDefaultClient() QuerierOpt {
	return WithClient(http.DefaultClient)
}

// WithClientTimeout is a QuerierOpt that sets client timeout.
func WithClientTimeout(t time.Duration) QuerierOpt {
	return func(q *Querier) (*Querier, error) {
		q.client.Timeout = t
		return q, nil
	}
}

// WithClientDefaultTimeout is a QuerierOpt that sets client timeout to 10 seconds.
func WithClientDefaultTimeout() QuerierOpt {
	return WithClientTimeout(10 * time.Second)
}

// NewQuery is a helper function that creates a new query.
// It parses rawURL, does a http request, checks the http response, and sets a handler for the response.
func (querier Querier) DoQuery(q *Query) *Query {
	q = querier.ParseRawURL(q)
	if q.Error != nil {
		return q
	}

	q = querier.DoHttpRequest(q)
	if q.Error != nil {
		return q
	}

	q = querier.CheckHttpResponse(q)
	if q.Error != nil {
		return q
	}

	return querier.SetHandleResponse(q)
}

// NewQueryPipeline is a concurrent version of NewQuery.
// ctx with cancel can be used to cancel the pipeline.
func (querier Querier) DoQueryPipeline(ctx context.Context, inputChan <-chan *Query) (<-chan *Query, <-chan error) {
	errChan := make(chan error)
	type Step struct {
		Step    func(*Query) *Query
		ChanIn  <-chan *Query
		ChanOut chan<- *Query
		Context context.Context
	}

	steps := []Step{}
	prevChan := inputChan
	for _, step := range []QuerierStep{
		querier.ParseRawURL,
		querier.DoHttpRequest,
		querier.CheckHttpResponse,
		querier.SetHandleResponse,
	} {
		nextChan := make(chan *Query)
		steps = append(steps, Step{
			Step:    step,
			ChanIn:  prevChan,
			ChanOut: nextChan,
			Context: ctx,
		})
		prevChan = nextChan
	}
	outputChan := prevChan

	// start pipeline workers
	for i := range steps {
		go func(step Step, isLastStep bool) {
			defer close(step.ChanOut)
			if isLastStep {
				defer close(errChan)
			}

			for {
				select {
				case <-step.Context.Done():
					return
				case q, ok := <-step.ChanIn:
					if !ok {
						return
					}
					q = step.Step(q)
					if q.Error != nil {
						errChan <- fmt.Errorf("%d-th query: %w", q.Id(), q.Error)
					} else {
						step.ChanOut <- q
					}
				}
			}
		}(steps[i], i == len(steps)-1)
	}

	return outputChan, errChan
}

// ParseRawURL is a helper function that parses rawURL string into a *url.URL.
func (querier Querier) ParseRawURL(q *Query) *Query {
	if u, err := url.Parse(q.RawURL); err != nil {
		q.Error = fmt.Errorf("error while url.Parse: %w", err)
	} else {
		q.URL = u
	}
	return q
}

// DoHttpRequest is a helper function that does a http request.
func (querier Querier) DoHttpRequest(q *Query) *Query {
	req, err := http.NewRequest(http.MethodGet, q.URL.String(), nil)
	if err != nil {
		q.Error = fmt.Errorf("error while .NewRequest: %w", err)
		return q
	}
	q.req = req

	for key, val := range querier.header {
		req.Header.Set(key, val)
	}

	resp, err := querier.client.Do(req)
	if err != nil {
		q.Error = fmt.Errorf("error while .Do: %w", err)
	}
	q.resp = resp
	return q
}

// CheckResponse is a helper function that checks if the response is valid.
func (querier Querier) CheckHttpResponse(q *Query) *Query {
	if q.resp.StatusCode != http.StatusOK {
		q.Error = fmt.Errorf("request error with error code %d", q.resp.StatusCode)
		return q
	}

	return q
}

// SetHandleResponse is a helper function that sets a handler for a response.
// It sets the handler based on the Content-Encoding header.
func (querier Querier) SetHandleResponse(q *Query) *Query {
	if handler, ok := querier.handler[q.resp.Header.Get("Content-Encoding")]; ok {
		fmt.Println("use selected handler")
		q.handler = handler
	} else {
		fmt.Println("use default handler")
		q.handler = DefaultTypeHandler
	}
	return q
}

// Query is a struct that contains query result.
// It is used to pass query result between different parser.
type Query struct {
	id      int
	RawURL  string
	URL     *url.URL
	News    *News
	Error   error
	handler ContentEncodeHandler
	req     *http.Request
	resp    *http.Response
}

var ErrNilHandler = errors.New("handler is nil")

func NewQuery(rawURL string) *Query {
	return &Query{
		RawURL: rawURL,
	}
}

func NewQueryWithId(id int, rawURL string) *Query {
	return &Query{
		id:     id,
		RawURL: rawURL,
	}
}

// NewTestQuery is a helper function that creates a new query for testing.
func NewTestQuery(status int, rawURL string, body io.ReadCloser) *Query {
	q := &Query{}
	if rawURL != "" {
		q = NewQuery(rawURL)
	}

	q.handler = DefaultTypeHandler
	q.resp = &http.Response{
		StatusCode: status,
		Body:       body,
	}
	return q
}

// SetId is a helper function that sets id for a query.
func (q *Query) SetId(id int) *Query {
	q.id = id
	return q
}

// Id returns the id of a query.
func (q Query) Id() int {
	return q.id
}

// RespHttpStatusCode returns the http status code of a query.
func (q Query) RespHttpStatusCode() int {
	return q.resp.StatusCode
}

// Content returns the content of a query.
// It returns an ErrNilHandler if the handler is nil.
func (q Query) Content() (io.ReadCloser, error) {
	if q.handler == nil {
		return nil, ErrNilHandler
	}
	return q.handler(q.resp.Body)
}

func (q Query) ToNewsCreateParam(md5hash string) *service.NewsCreateRequest {
	return &service.NewsCreateRequest{
		Md5Hash:     md5hash,
		Guid:        q.News.GUID,
		Author:      q.News.Author,
		Title:       q.News.Title,
		Link:        q.News.Link.String(),
		Description: q.News.Description,
		Language:    q.News.Language,
		Content:     q.News.Content,
		Category:    q.News.Category,
		Source:      q.News.Link.Host,
		RelatedGuid: q.News.RelatedGUID,
		PublishedAt: q.News.PubDate,
	}
}

// the ToDoc method is a helper function that convert io.ReadCloser to *goquery.Document
func ToDoc(rc io.ReadCloser) (*goquery.Document, error) {
	defer rc.Close()
	return goquery.NewDocumentFromReader(rc)
}

// FmtReq and FmtResp are for debugging, they format request and response into human readable string
// they should be removed in production
func FmtReq(req *http.Request) string {
	sb := &strings.Builder{}
	fmt.Fprintln(sb, "HTTP REQUEST:")
	fmt.Fprintf(sb, "- URL: %s\n", req.URL.String())
	fmt.Fprintf(sb, "- Content Length: %d\n", req.ContentLength)

	fmt.Fprintln(sb, "- Header:")
	for key, val := range req.Header {
		fmt.Fprintf(sb, "\t %s: %s\n", key, strings.Join(val, ", "))
	}

	if len(req.Form) > 0 {
		fmt.Fprintln(sb, "- Form:")
		for key, val := range req.Form {
			fmt.Fprintf(sb, "\t %s: %s\n", key, strings.Join(val, ", "))
		}
	}
	return sb.String()
}

// FmtReq and FmtResp are for debugging, they format request and response into human readable string
// they should be remove
func FmtResp(resp *http.Response) string {
	sb := &strings.Builder{}

	fmt.Fprintln(sb, "HTTP RESPONSE:")
	fmt.Fprintf(sb, "- Status: %s (%d)\n", resp.Status, resp.StatusCode)
	fmt.Fprintf(sb, "- Content Length: %d\n", resp.ContentLength)

	fmt.Fprintln(sb, "- Header:")
	for key, val := range resp.Header {
		fmt.Fprintf(sb, "\t %s: %s\n", key, strings.Join(val, ", "))
	}

	fmt.Fprintln(sb, "- Body:")
	if body, err := io.ReadAll(resp.Body); err == nil {
		fmt.Fprintln(sb, "\t nil")
	} else {
		bodyStr := string(body)
		if len(bodyStr) > 30 {
			fmt.Fprintf(sb, "\t %s...(%d characters)\n", bodyStr[:30], len(bodyStr)-30)
		} else {
			fmt.Fprintln(sb, bodyStr)
		}
	}

	return sb.String()
}
