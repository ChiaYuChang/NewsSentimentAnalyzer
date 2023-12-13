package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	cohere "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/Cohere"
	openai "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/OpenAI"
	newsparser "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/grpc/newsParser"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/go-chi/chi/v5"
	"github.com/pemistahl/lingua-go"
	"github.com/redis/go-redis/v9"

	ldcli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/grpc/languageDetector"
	ldpf "github.com/ChiaYuChang/NewsSentimentAnalyzer/proto/language_detector"
)

func (repo APIRepo) GetAnalyzer(w http.ResponseWriter, req *http.Request) {
	// pcid := chi.URLParam(req, "pcid")
	pageData := object.AnalyzerPage{
		Page: object.Page{
			HeadConent: view.SharedHeadContent(),
			Title:      "analyzer",
		},
		Prompt:  map[string]string{},
		Version: "v1",
	}
	pageData.Prompt["openai-sentiment"] = openai.SentimentAnalysisPrompt
	pageData.Prompt["cohere-sentiment"] = cohere.SentimentAnalysisPrompt

	err := repo.View.ExecuteTemplate(w, "analyzer.gotmpl", pageData)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to execute template")
	}
	return
}

func (repo APIRepo) PostAnalyzer(w http.ResponseWriter, req *http.Request) {
	resp := pageform.PreviewPostResp{}
	pcid := chi.URLParam(req, "pcid")
	aid, _ := strconv.Atoi(req.URL.Query().Get("aid"))

	if _, err := repo.Cache.Get(context.Background(), pcid).Result(); err == redis.Nil {
		resp.WithEcError(ec.MustGetEcErr(ec.ECGone).
			WithDetails("cache already expire").
			WithDetails(pcid))
		w.WriteHeader(resp.Error.Code)
		resp.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["gone"]
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	err := req.ParseForm()
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to parse form")
		resp.WithEcError(ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("Failed to parse form", err.Error()))
		resp.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["bad-request"]
		w.WriteHeader(resp.Error.Code)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	var fdata api.AnalyzerOption
	err = repo.FormDecoder.Decode(&fdata, req.PostForm)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to decode form")
		resp.WithEcError(ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("Failed to decode form", err.Error()))
		resp.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["bad-request"]
		w.WriteHeader(resp.Error.Code)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	_, err = repo.Cache.JSONSet(pcid, ".analyzer_options", fdata)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to write cache")
		resp.WithEcError(ec.MustGetEcErr(ec.ECGone).
			WithDetails(err.Error()))
		resp.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["gone"]
		w.WriteHeader(resp.Error.Code)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	jid, err := repo.CacheToStore(pcid, aid, fdata.APIId)
	if err != nil {
		global.Logger.Error().Str("error", err.Error()).Msg("Failed to cache to store")
		var ecErr *ec.Error
		if ok := errors.As(err, &ecErr); !ok {
			ecErr = ec.MustGetEcErr(ec.ECServerError)
		}
		resp.WithEcError(ecErr)
		resp.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["internal-server-error"]
		w.WriteHeader(ecErr.HttpStatusCode)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	global.Logger.Info().
		Str("url", req.URL.String()).
		Str("pcid", pcid).
		Int("jid", jid).
		Int("aid", aid).
		Int("lid", fdata.APIId).
		Msg("done")

	resp.RedirectURL = fmt.Sprintf("/v1%s?jid=%d", global.AppVar.App.RoutePattern.Page["job"], jid)
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(resp)
	w.Write(b)
	return
}

func (repo APIRepo) CacheToStore(pcid string, aid, lid int) (int, error) {
	// save cache to premint storage
	global.Logger.Info().Msg("read cache")
	res, err := repo.Cache.JSONGet(pcid, ".")
	if err != nil {
		return 0, ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while get cache").
			WithDetails(err.Error())
	}

	global.Logger.Info().Msg("unmarshal cache")
	var cache api.PreviewCache
	err = json.Unmarshal(res.([]byte), &cache)
	if err != nil {
		return 0, ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while unmarshaling cache").
			WithDetails(err.Error())
	}

	selectedItem := cache.SelectedItems()
	global.Logger.Info().Int("n", len(selectedItem)).Msg("OK")
	ldcli.MustGetLanguageDetectorClient()

	reqChan := make(chan *ldpf.LanguageDetectRequest)
	respChan, err, errChan := ldcli.MustGetLanguageDetectorClient().DetectLanguage(context.TODO(), reqChan)
	if err != nil {
		return 0, ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while detect language").
			WithDetails(err.Error())
	}

	go func(errChan <-chan error) {
		for err := range errChan {
			global.Logger.Error().Err(err).Msg("error while detect language")
		}
	}(errChan)

	go func(reqChan chan<- *ldpf.LanguageDetectRequest) {
		defer close(reqChan)
		thr := float64(.9)
		for _, item := range selectedItem {
			reqChan <- &ldpf.LanguageDetectRequest{
				Id:        item.Id.String(),
				Text:      item.Title,
				Threshold: &thr,
			}
		}
		global.Logger.Info().Msg("detect language req go run time done")
	}(reqChan)

	cnChan := make(chan *service.NewsCreateRequest, 10)
	go func(respChan <-chan *ldpf.LanguageDetectResponse,
		cnChan chan<- *service.NewsCreateRequest) {
		defer close(cnChan)
		cli, _ := newsparser.GetNewsParserClient()
		i := int64(0)
		for resp := range respChan {
			prev := selectedItem[resp.GetId()]
			u, _ := url.Parse(prev.Link)
			// guid := parser.GetDefaultParser().ToGUID(u)
			_, guid, err := cli.GetGUID(context.TODO(), i, u.String())
			if err != nil {
				global.Logger.Error().
					Err(err).
					Str("url", u.String()).
					Msg("error while GetGUID")
				continue
			}
			lang := lingua.Language(int(resp.GetLanguage())).IsoCode639_1()
			cnChan <- prev.ToNewsCreateRequest(guid, lang.String(), u.Host, nil)
			i++
		}
		global.Logger.Info().Msg("create news req go run time done")
	}(respChan, cnChan)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cache.AnalyzerOptions.APIId = 0 // omit analyzer api id
	result, err := repo.Service.TX().DoCacheToStoreTx(
		ctx, cache.Query.UserId, strings.TrimSuffix(
			strings.TrimPrefix(pcid, global.PREVIEW_CACHE_KEY_PREFIX),
			global.PREVIEW_CACHE_KEY_SUFFIX),
		int16(aid), cache.Query.RawQuery, int16(lid),
		cache.AnalyzerOptions.ToString("", ""), cnChan)

	if err != nil {
		return 0, err
	}

	_, err = repo.Cache.Del(context.Background(), pcid).Result()
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to delete cache")
		return int(result.JobId), err
	}
	return int(result.JobId), nil
}
