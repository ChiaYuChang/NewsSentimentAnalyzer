package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
)

var ErrPageFormHandlerHasBeenRegistered = errors.New("the PageFormHandler has already been registered")
var ErrHandlerNotFound = errors.New("unregistered handler")
var ErrUnknownEndpoint = errors.New("unknown endpoint")
var ErrNotSupportedEndpoint = errors.New("currently not support endpoint")

var PageFormHandlerRepo = pageFormHandlerRepo{}

func RegisterPageForm(pf pageform.PageForm, handler PageFormHandler, mapping map[string]string) {
	PageFormHandlerRepo.RegisterPageForm(pf, handler, mapping)
}

type ClientRepoKey [2]string

func NewRepoMapKey(apiName string, endpointName string) ClientRepoKey {
	return ClientRepoKey{apiName, endpointName}
}

func (k ClientRepoKey) APIName() string {
	return k[0]
}

func (k ClientRepoKey) EndpointName() string {
	return k[1]
}

func (k ClientRepoKey) String() string {
	return fmt.Sprintf("%s-%s", k[0], k[1])
}

type PageFormHandler interface {
	Handle(apikey string, pageForm pageform.PageForm) (api.Request, error)
	Parse(response *http.Response) (api.Response, error)
}

type pageFormHandlerRepo map[ClientRepoKey]PageFormHandler

func (repo pageFormHandlerRepo) RegisterPageForm(
	pf pageform.PageForm, handler PageFormHandler, mapping map[string]string) error {
	key := NewRepoMapKey(pf.API(), pf.Endpoint())
	if _, ok := repo[key]; ok {
		return ErrPageFormHandlerHasBeenRegistered
	}
	repo[key] = handler
	if mapping != nil {
		if ep, ok := mapping[pf.Endpoint()]; ok {
			key := NewRepoMapKey(pf.API(), ep)
			repo[key] = handler
		}
	}
	return nil
}

func (repo pageFormHandlerRepo) Handle(apikey string, pf pageform.PageForm) (api.Request, error) {
	key := NewRepoMapKey(pf.API(), pf.Endpoint())
	handler, ok := repo[key]
	if !ok {
		return nil, ErrHandlerNotFound
	}

	return handler.Handle(apikey, pf)
}

// func (repo pageFormHandlerRepo) Parse(response *http.Response) (api.Response, error) {

// }

// func (repo pageFormHandlerRepo) Do(cli http.Client, apikey string, pf pageform.PageForm) error {
// 	key := newRepoMapKey(pf.API(), pf.Endpoint())
// 	handler, ok := repo[key]
// 	if !ok {
// 		return ErrHandlerNotFound
// 	}

// 	q, err := handler.Handle(apikey, pf)
// 	if err != nil {
// 		return fmt.Errorf("error while .Handle: %w", err)
// 	}

// 	req, err := q.ToHttpRequest()
// 	if err != nil {
// 		return fmt.Errorf("error while .ToRequest: %w", err)
// 	}

// 	cPars := make(chan *service.NewsCreateRequest)
// 	go func() {
// 		for p := range cPars {
// 			global.Logger.
// 				Info().
// 				Str("md5", p.Md5Hash).
// 				Str("title", p.Title).
// 				Time("publish_at", p.PublishedAt.UTC()).
// 				Msg("Create an article")
// 		}
// 	}()

// 	reqs := []*http.Request{req}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	wg := &sync.WaitGroup{}
// 	for i := 0; i < len(reqs); i++ {
// 		httpResp, err := cli.Do(reqs[i])
// 		if err != nil {
// 			cancel()
// 			return err
// 		}

// 		resp, err := handler.Parse(httpResp)
// 		if err != nil {
// 			cancel()
// 			return err
// 		}

// 		if resp.HasNext() {
// 			next, err := resp.NextPageRequest(nil)
// 			if err != nil {
// 				cancel()
// 				return err
// 			}
// 			reqs = append(reqs, next)
// 		}

// 		wg.Add(1)
// 		go func(wg *sync.WaitGroup) {
// 			defer wg.Done()
// 			resp.ToNews(ctx, wg, cPars)
// 		}(wg)
// 	}
// 	wg.Wait()
// 	close(cPars)
// 	return nil
// }
