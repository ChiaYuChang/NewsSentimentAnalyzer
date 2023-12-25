package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	// http client
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	cohere "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/Cohere"
	openai "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api/OpenAI"

	// grpc client
	ldcli "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/grpc/languageDetector"
	newsparser "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/grpc/newsParser"
	ldpf "github.com/ChiaYuChang/NewsSentimentAnalyzer/proto/language_detector"

	// http server
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	"github.com/go-chi/chi/v5"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/jackc/pgerrcode"
	"github.com/pemistahl/lingua-go"
	"github.com/redis/go-redis/v9"
	// pgv "github.com/pgvector/pgvector-go"
)

func (repo APIRepo) GetAnalyzer(w http.ResponseWriter, req *http.Request) {
	pcid := chi.URLParam(req, "pcid")
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

	repo.Cache.ExpireGT(context.Background(), pcid, global.CacheExpireDefault)

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

	if res, err := repo.Cache.JSONGet(pcid, ".is_done"); err == redis.Nil {
		resp.WithEcError(ec.MustGetEcErr(ec.ECGone).
			WithDetails("cache already expire").
			WithDetails(pcid)).
			WithRedirectURL(global.AppVar.App.RoutePattern.ErrorPage["gone"])
		w.WriteHeader(resp.HttpStatusCode())
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	} else {
		if bytes.Compare([]byte("true"), res.([]byte)) == 0 {
			resp.Error = &pageform.PreviewError{
				Code:        http.StatusInternalServerError,
				PgxCode:     pgerrcode.UniqueViolation,
				Message:     "You have already submitted this request",
				RedirectURL: fmt.Sprintf("/v1%s", global.AppVar.App.RoutePattern.Page["job"]),
			}
			w.WriteHeader(resp.HttpStatusCode())
			b, _ := json.Marshal(resp)
			w.Write(b)
			return
		}

	}

	err := req.ParseForm()
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to parse form")
		resp.WithEcError(
			ec.MustGetEcErr(ec.ECBadRequest).
				WithMessage(err.Error()).
				WithDetails("Failed to parse form")).
			WithRedirectURL(global.AppVar.App.RoutePattern.ErrorPage["bad-request"])
		w.WriteHeader(resp.HttpStatusCode())
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	var fdata service.AnalyzerOption
	err = repo.FormDecoder.Decode(&fdata, req.PostForm)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to decode form")
		resp.WithEcError(ec.MustGetEcErr(ec.ECBadRequest).
			WithMessage(err.Error()).
			WithDetails("Failed to decode form")).
			WithRedirectURL(global.AppVar.App.RoutePattern.ErrorPage["bad-request"])
		w.WriteHeader(resp.HttpStatusCode())
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	_, err = repo.Cache.JSONSet(pcid, ".analyzer_options", fdata)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to write cache")
		resp.WithEcError(ec.MustGetEcErr(ec.ECGone).
			WithMessage(err.Error())).
			WithRedirectURL(global.AppVar.App.RoutePattern.ErrorPage["gone"])
		w.WriteHeader(resp.HttpStatusCode())
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	jid, ecErr := repo.CacheToStore(pcid, aid, fdata.APIId)
	if ecErr != nil {
		resp.WithEcError(ecErr).
			WithOutDetails().
			WithOutMessage().
			WithRedirectURL(global.AppVar.App.RoutePattern.ErrorPage["internal-server-error"])
		w.WriteHeader(ecErr.HttpStatusCode)
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	resp.RedirectURL = fmt.Sprintf("/v1%s?jid=%d", global.AppVar.App.RoutePattern.Page["job"], jid)
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(resp)
	w.Write(b)
	return
}

func (repo APIRepo) CacheToStore(pcid string, aid, lid int) (int, *ec.Error) {
	// save cache to premint storage
	res, err := repo.Cache.JSONGet(pcid, ".")
	if err != nil {
		return 0, ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while get cache").
			WithMessage(err.Error())
	}

	global.Logger.Info().Msg("unmarshal cache")
	var cache api.PreviewCache
	err = json.Unmarshal(res.([]byte), &cache)
	if err != nil {
		return 0, ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while unmarshaling cache").
			WithMessage(err.Error())
	}

	selectedItem := cache.SelectedItems()
	global.Logger.Info().Int("n", len(selectedItem)).Msg("OK")
	ldcli.MustGetLanguageDetectorClient()

	reqChan := make(chan *ldpf.LanguageDetectRequest)
	respChan, err, errChan := ldcli.MustGetLanguageDetectorClient().DetectLanguage(context.TODO(), reqChan)
	if err != nil {
		return 0, ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while detect language").
			WithMessage(err.Error())
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
		if ecErr, ok := err.(*ec.Error); ok {
			return 0, ecErr
		} else {
			return 0, ec.MustGetEcErr(ec.ECServerError).
				WithMessage(err.Error()).
				WithDetails("error while DoCacheToStoreTx")
		}
	}

	repo.Cache.JSONSet(pcid, ".is_done", true)
	_, err = repo.Cache.ExpireLT(context.Background(), pcid, 1*time.Minute).Result()
	if err != nil {
		return int(result.JobId), ec.MustGetEcErr(ec.ECServerError).
			WithMessage(err.Error()).
			WithDetails("Failed to delete cache")
	}

	return int(result.JobId), nil
}
