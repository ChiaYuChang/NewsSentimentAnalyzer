package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/router/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	tokenmaker "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/tokenMaker"
	val "github.com/go-playground/validator/v10"
)

var ErrEndpointNotFount = errors.New("unregistered endpoint")

type repoMapKey [2]string

func newRepoMapKey(apiName string, endpointName string) repoMapKey {
	return repoMapKey{apiName, endpointName}
}

func (k repoMapKey) APIName() string {
	return k[0]
}

func (k repoMapKey) EndpointName() string {
	return k[1]
}

func (k repoMapKey) String() string {
	return fmt.Sprintf("%s-%s", k[0], k[1])
}

type EndpointRepo struct {
	apiRepo  APIRepo
	val      *val.Validate
	PageView map[repoMapKey][]byte
	PageForm map[repoMapKey]pageform.PageForm
}

func NewEndpointRepo(apiRepo APIRepo, v *val.Validate) EndpointRepo {
	return EndpointRepo{
		apiRepo:  apiRepo,
		val:      v,
		PageView: make(map[repoMapKey][]byte),
	}
}

func (repo *EndpointRepo) RegisterEndpointsPageView(apiName, endpointName, templateName string) error {
	pageData := object.APIEndpointPage{
		Page: object.Page{
			HeadConent: view.SharedHeadContent,
			Title:      endpointName,
		},
		API:      apiName,
		Version:  global.AppVar.Server.APIVersion,
		Endpoint: endpointName,
	}

	buffer := bytes.NewBufferString("")
	err := repo.apiRepo.View.ExecuteTemplate(buffer, templateName, pageData)
	if err != nil {
		return err
	}

	key := newRepoMapKey(apiName, endpointName)
	repo.PageView[key] = buffer.Bytes()
	return nil
}

func (repo *EndpointRepo) WritePageViewTo(apiName, endpointName string, w io.Writer) error {
	var view []byte
	var ok bool
	var err error
	if view, ok = repo.PageView[newRepoMapKey(apiName, endpointName)]; !ok {
		return ErrEndpointNotFount
	}
	_, err = w.Write(view)
	return err
}

func (repo EndpointRepo) GetAPIEndpoints(key repoMapKey) (http.HandlerFunc, error) {
	page, ok := repo.PageView[key]
	if !ok {
		return nil, ec.MustGetEcErr(ec.ECNotFound)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(page)
	}, nil
}

func (repo EndpointRepo) PostAPIEndpoints(key repoMapKey) (http.HandlerFunc, error) {
	pf, err := pageform.PageFormRepo.Get(key.APIName(), key.EndpointName())
	if err != nil || pf == nil {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		if err != nil {
			ecErr.WithDetails(err.Error())
			ecErr.WithDetails(key.String() + " not found")
		}
		if pf == nil {
			ecErr.WithDetails("pf is nil")
		}
		return nil, ecErr
	}
	return func(w http.ResponseWriter, req *http.Request) {
		postEndpoints(repo, pf, w, req)
	}, nil
}

func postEndpoints(repo EndpointRepo, obj pageform.PageForm, w http.ResponseWriter, req *http.Request) {
	userInfo, ok := req.Context().Value(global.CtxUserInfo).(tokenmaker.Payload)
	if !ok {
		ecErr := ec.MustGetEcErr(ec.ECServerError)
		ecErr.WithDetails("user information not found")
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(ecErr.MustToJson())
		return
	}
	fmt.Println("Get User Info OK")

	if err := req.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Parse form OK")

	obj, err := obj.FormDecodeAndValidate(repo.apiRepo.FormDecoder, repo.val, req.PostForm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Decode and Validate OK")

	fmt.Fprintf(w,
		"User: %s(%s)\nAPI: %s\nEndpoint: %s\n%s\n",
		userInfo.GetUsername(),
		userInfo.GetRole(),
		obj.API(),
		obj.Endpoint(),
		obj.String(),
	)

	return
}
