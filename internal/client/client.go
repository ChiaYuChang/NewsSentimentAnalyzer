package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/google/uuid"
)

var ErrPageFormHandlerHasBeenRegistered = errors.New("the PageFormHandler has already been registered")
var ErrHandlerNotFound = errors.New("unregistered handler")
var ErrUnknownEndpoint = errors.New("unknown endpoint")
var ErrNotSupportedEndpoint = errors.New("currently not support endpoint")

var HandlerRepo = handlerRepo{}

func RegisterHandler(pf pageform.PageForm, handler Handler, mapping map[string]string) {
	HandlerRepo.RegisterHandler(pf, handler, mapping)
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

type Handler interface {
	Handle(apikey string, uid uuid.UUID, pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error)
	RequestFromCacheQuery(cq api.CacheQuery) (req api.Request, err error)
	Parse(response *http.Response) (resp api.Response, err error)
}
type handlerRepo map[ClientRepoKey]Handler

func (repo handlerRepo) RegisterHandler(
	pf pageform.PageForm, handler Handler, mapping map[string]string) error {
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

func (repo handlerRepo) Get(apiname, apiep string) (Handler, error) {
	key := NewRepoMapKey(apiname, apiep)
	handler, ok := repo[key]
	if !ok {
		return nil, ErrHandlerNotFound
	}
	return handler, nil
}

func (repo handlerRepo) GetByCacheQuery(cache api.CacheQuery) (Handler, error) {
	return repo.Get(cache.API.Name, cache.API.Endpoint)
}

func (repo handlerRepo) Do(req api.Request, handler Handler) (api.Response, error) {
	var httpReq *http.Request
	var httpResp *http.Response
	var err error

	if httpReq, err = req.ToHttpRequest(); err != nil {
		return nil, err
	}

	httpResp, err = http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return handler.Parse(httpResp)
}

func (repo handlerRepo) Handle(apikey string, uid uuid.UUID,
	pf pageform.PageForm) (ckey string, cache *api.PreviewCache, err error) {

	handler, err := repo.Get(pf.API(), pf.Endpoint())
	if err != nil {
		return "", nil, err
	}
	return handler.Handle(apikey, uid, pf)
}

func (repo handlerRepo) Parse(cq api.CacheQuery, response *http.Response) (api.Response, error) {
	handler, err := repo.GetByCacheQuery(cq)
	if err != nil {
		return nil, err
	}
	return handler.Parse(response)
}
