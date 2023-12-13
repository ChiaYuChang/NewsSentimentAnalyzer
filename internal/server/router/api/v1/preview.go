package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/global"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/client/api"
	pageform "github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/pageForm"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/view/object"
	ec "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/errorCode"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func (repo APIRepo) GetPreview(w http.ResponseWriter, req *http.Request) {
	pageData := object.ResultSecectorPage{
		Page: object.Page{
			HeadConent: view.SharedHeadContent(),
			Title:      "Result Selector",
		},
		Version: repo.Version,
	}

	_ = req.ParseForm()
	err := repo.View.ExecuteTemplate(w, "preview.gotmpl", pageData)
	if err != nil {
		global.Logger.Error().
			Err(err).
			Msg("error executing template preivew.gotmpl")
	}
}

func (repo APIRepo) PostPreview(w http.ResponseWriter, req *http.Request) {
	global.Logger.Info().Msg("hit post preview")

	global.Logger.Info().Msg("parse form")
	resp := pageform.PreviewPostResp{}
	err := req.ParseForm()
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while ParseForm").
			WithDetails(err.Error())
		resp.Error = &pageform.PreviewError{
			Code:        ecErr.HttpStatusCode,
			Message:     ecErr.Message,
			Detail:      ecErr.Details,
			RedirectURL: global.AppVar.App.RoutePattern.ErrorPage["bad-request"],
		}

		b, _ := json.Marshal(resp)
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	global.Logger.Info().Msg("parse pcid")
	pcid := chi.URLParam(req, "pcid")
	aid, _ := strconv.Atoi(req.URL.Query().Get("aid"))
	eid, _ := strconv.Atoi(req.URL.Query().Get("eid"))

	global.Logger.Info().
		Str("pcid", pcid).
		Int("aid", aid).
		Int("eid", eid).
		Msg("read params")

	var pf pageform.PreviewPostForm
	err = repo.FormDecoder.Decode(&pf, req.Form)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while Decode").
			WithDetails(err.Error())
		resp.Error = &pageform.PreviewError{
			Code:        ecErr.HttpStatusCode,
			Message:     ecErr.Message,
			Detail:      ecErr.Details,
			RedirectURL: global.AppVar.App.RoutePattern.ErrorPage["bad-request"],
		}

		b, _ := json.Marshal(resp)
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	// global.Logger.Info().Msg("read cache")
	// res, err := repo.Cache.JSONGet(pcid, ".")
	// if err != nil {
	// 	ecErr := ec.MustGetEcErr(ec.ECBadRequest).
	// 		WithDetails("error while get cache").
	// 		WithDetails(err.Error())
	// 	resp.Error = &pageform.PreviewError{
	// 		Code:        ecErr.HttpStatusCode,
	// 		Message:     ecErr.Message,
	// 		Detail:      ecErr.Details,
	// 		RedirectURL: global.AppVar.App.RoutePattern.ErrorPage["internal-server-error"],
	// 	}

	// 	b, _ := json.Marshal(resp)
	// 	w.WriteHeader(ecErr.HttpStatusCode)
	// 	w.Write(b)
	// 	return
	// }

	global.Logger.Info().Strs("item", pf.Item).Msg("select preview")

	// var cache api.PreviewCache
	// err = json.Unmarshal(res.([]byte), &cache)
	// if err != nil {
	// 	ecErr := ec.MustGetEcErr(ec.ECBadRequest).
	// 		WithDetails("error while unmarshaling cache").
	// 		WithDetails(err.Error())
	// 	resp.Error = &pageform.PreviewError{
	// 		Code:        ecErr.HttpStatusCode,
	// 		Message:     ecErr.Message,
	// 		Detail:      ecErr.Details,
	// 		RedirectURL: global.AppVar.App.RoutePattern.ErrorPage["internal-server-error"],
	// 	}

	// 	b, _ := json.Marshal(resp)
	// 	w.WriteHeader(ecErr.HttpStatusCode)
	// 	w.Write(b)
	// 	return
	// }

	// prev := cache.NewsItem
	sort.Sort(sort.StringSlice(pf.Item))
	// sort.Sort(api.SortById(prev))
	// selectedPrev := make(map[string]api.NewsPreview)
	// if pf.SelectAll {
	// 	for i := range prev {
	// 		selectedPrev[prev[i].Id.String()] = prev[i]
	// 	}
	// } else {
	// 	for i, j := 0, 0; i < len(prev) && j < len(pf.Item); {
	// 		if prev[i].Id.String() < pf.Item[j] {
	// 			i++
	// 		} else if prev[i].Id.String() > pf.Item[j] {
	// 			j++
	// 		} else {
	// 			selectedPrev[prev[i].Id.String()] = prev[i]
	// 			j++
	// 			i++
	// 		}
	// 	}
	// }
	// global.Logger.Info().Int("n", len(selectedPrev)).Msg("OK")

	if pf.SelectAll {
		repo.Cache.JSONSet(pcid, ".selected_all", pf.SelectAll)
	} else {
		repo.Cache.JSONSet(pcid, ".selected_nid", pf.Item)
	}
	repo.Cache.ExpireGT(context.Background(), pcid, 10*time.Minute)

	// reqChan := make(chan *languageDetector.LanguageDetectRequest)
	// respChan, err, errChan := languagedetector.MustGetLanguageDetectorClient().DetectLanguage(context.TODO(), reqChan)
	// if err != nil {
	// 	ecErr := ec.MustGetEcErr(ec.ECBadRequest).
	// 		WithDetails("error while detect language").
	// 		WithDetails(err.Error())
	// 	resp.Error = &pageform.PreviewError{
	// 		Code:        ecErr.HttpStatusCode,
	// 		Message:     ecErr.Message,
	// 		Detail:      ecErr.Details,
	// 		RedirectURL: global.AppVar.App.RoutePattern.ErrorPage["internal-server-error"],
	// 	}
	// 	b, _ := json.Marshal(resp)
	// 	w.WriteHeader(ecErr.HttpStatusCode)
	// 	w.Write(b)
	// 	return
	// }

	// go func(reqChan chan<- *languageDetector.LanguageDetectRequest) {
	// 	defer close(reqChan)
	// 	thr := float64(.9)
	// 	for key, val := range selectedPrev {
	// 		reqChan <- &languageDetector.LanguageDetectRequest{
	// 			Id:        key,
	// 			Text:      val.Title,
	// 			Threshold: &thr,
	// 		}
	// 	}
	// 	global.Logger.Info().Msg("detect language req go run time done")
	// }(reqChan)

	// cnChan := make(chan *service.NewsCreateRequest, 10)
	// go func(respChan <-chan *languageDetector.LanguageDetectResponse,
	// 	cnChan chan<- *service.NewsCreateRequest) {
	// 	defer close(cnChan)
	// 	cli, _ := newsparser.GetNewsParserClient()
	// 	i := int64(0)
	// 	for resp := range respChan {
	// 		prev := selectedPrev[resp.GetId()]
	// 		u, _ := url.Parse(prev.Link)
	// 		// guid := parser.GetDefaultParser().ToGUID(u)
	// 		_, guid, err := cli.GetGUID(context.TODO(), i, u.String())
	// 		if err != nil {
	// 			global.Logger.Error().
	// 				Err(err).
	// 				Str("url", u.String()).
	// 				Msg("error while GetGUID")
	// 			continue
	// 		}
	// 		lang := lingua.Language(int(resp.GetLanguage())).IsoCode639_1()
	// 		cnChan <- prev.ToNewsCreateRequest(guid, lang.String(), u.Host, nil)
	// 		i++
	// 	}
	// 	global.Logger.Info().Msg("to create news req go run time done")
	// }(respChan, cnChan)

	// go func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 	defer cancel()
	// 	result, err := repo.Service.TX().DoCacheToStoreTx(
	// 		ctx, cache.Query.UserId, strings.TrimSuffix(
	// 			strings.TrimPrefix(pcid, global.PREVIEW_CACHE_KEY_PREFIX),
	// 			global.PREVIEW_CACHE_KEY_SUFFIX),
	// 		int16(aid), cache.Query.RawQuery, 5, "{}", cnChan)
	// 	if err != nil {
	// 		global.Logger.Error().Err(err).Msg("error while DoCacheToStoreTx")
	// 	} else {
	// 		if result == nil {
	// 			global.Logger.Error().Msg("result is nil")
	// 			return
	// 		}
	// 		global.Logger.Info().
	// 			Int64("jid", result.JobId).
	// 			Msg("cache to store tx done")
	// 		for _, r := range result.NewsJobCreateResults {
	// 			if r.Error != nil {
	// 				global.Logger.Error().
	// 					Err(r.Error).
	// 					Msg("error while cache to store tx")
	// 			} else {
	// 				global.Logger.Info().
	// 					Str("Md5Hash", r.Md5Hash).
	// 					Int64("news_id", r.NewsID).
	// 					Int64("news_job_id", r.NewsJobId).
	// 					Msg("cache to store tx result")
	// 			}
	// 		}
	// 	}
	// }()

	// go func(errChan <-chan error) {
	// 	for err := range errChan {
	// 		global.Logger.Error().Err(err).Msg("error while detect language")
	// 	}
	// }(errChan)

	resp.RedirectURL = fmt.Sprintf("/%s/%s/%s?aid=%d&eid=%d", repo.Version, "analyzer", pcid, aid, eid)
	b, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

func (repo APIRepo) GetFetchNextPage(w http.ResponseWriter, req *http.Request) {
	var prev []api.NewsPreview
	var hasNext bool
	var ecErr *ec.Error
	var resp api.Response
	var pcid string
	var cq api.CacheQuery

	_ = req.ParseForm()
	pcid = chi.URLParam(req, "pcid")
	isFirst := req.URL.Query().Get("first") == "true"
	if isFirst {
		res, err := repo.Cache.JSONGet(pcid, ".news_item")
		if err == nil && res != nil {
			_ = json.Unmarshal(res.([]byte), &prev)
		}
	}

	if cq, ecErr = getCacheQuery(repo, pcid); ecErr == nil {
		if len(prev) > 0 {
			hasNext = !api.IsLastPageToken(cq.NextPage)
		} else {
			if resp, ecErr = fetchNextPage(cq); ecErr == nil {
				prev, hasNext, ecErr = getPreviewsAndUpdateCache(repo, pcid, resp)
			}
		}
	}
	repo.Cache.ExpireGT(context.Background(), pcid, 5*time.Minute)

	for i := range prev {
		prev[i].Content = ""
	}

	respObj := pageform.PreviewResponse{Items: prev, HasNext: hasNext}
	if ecErr != nil {
		respObj.Error = &pageform.PreviewError{
			Code:    ecErr.HttpStatusCode,
			Message: ecErr.Message,
			Detail:  ecErr.Details,
		}
		switch ecErr.HttpStatusCode {
		case http.StatusBadRequest:
			respObj.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["bad-request"]
		case http.StatusGone:
			respObj.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["gone"]
		case http.StatusInternalServerError:
			respObj.Error.RedirectURL = global.AppVar.App.RoutePattern.ErrorPage["internal-server-error"]
		}
		respObj.Error.RedirectURL = ""
	}

	bprev, _ := json.Marshal(respObj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bprev)
}

func getCacheQuery(repo APIRepo, pcid string) (api.CacheQuery, *ec.Error) {
	var cq api.CacheQuery

	b, err := repo.Cache.JSONGet(pcid, ".query")
	if err != nil {
		var ecErr *ec.Error
		if err == redis.Nil {
			ecErr = ec.MustGetEcErr(ec.ECGone).
				WithDetails("cache expired")
		} else {
			ecErr = ec.MustGetEcErr(ec.ECServerError).
				WithDetails("error getting query cache").
				WithDetails(err.Error())
		}
		return cq, ecErr
	}

	// read cache.query
	err = json.Unmarshal(b.([]byte), &cq)
	if err != nil {
		global.Logger.Debug().Err(err).Msg("error unmarshal cq")
		ceErr := ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error unmarshal cq").
			WithDetails(err.Error())
		return cq, ceErr
	}
	return cq, nil
}

func fetchNextPage(cq api.CacheQuery) (api.Response, *ec.Error) {
	handler, err := client.HandlerRepo.GetByCacheQuery(cq)
	if err != nil {
		return nil, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while get handler from client.HandlerRepo").
			WithDetails(err.Error())
	}

	// rebuild query from cache
	req, err := handler.RequestFromCacheQuery(cq)
	if err != nil {
		return nil, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while .RequestFromCacheQuery").
			WithDetails(err.Error())
	}

	global.Logger.Info().
		Str("uid", cq.UserId.String()).
		Str("salt", cq.Salt).
		Str("rawQuery", cq.RawQuery).
		Str("nextPage", cq.NextPage.String()).
		Msg("rebuild cache query ok")

	// do request
	resp, err := client.HandlerRepo.Do(req, handler)
	if err != nil {
		var ecErr *ec.Error
		var ok bool

		ecErr, ok = err.(*ec.Error)
		if !ok {
			ecErr = ec.MustGetEcErr(ec.ECServerError)
		}
		ecErr.WithDetails("error while .Do")
		return nil, ecErr
	}

	return resp, nil
}

func getPreviewsAndUpdateCache(repo APIRepo, pcid string, resp api.Response) ([]api.NewsPreview, bool, *ec.Error) {
	// append prev to cache
	next, prev := resp.ToNewsItemList()

	if len(prev) == 0 {
		return prev, !api.IsLastPageToken(next), nil
	}

	values := make([]any, len(prev))
	for i, p := range prev {
		values[i] = p
	}
	if _, err := repo.Cache.JSONArrAppend(pcid, ".news_item", values...); err != nil {
		return nil, false, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while appending preview items to cache").
			WithDetails("error append prev to cache").
			WithDetails(err.Error())
	}

	// set next page token to cache
	if _, err := repo.Cache.JSONSet(pcid, ".query.next_page", next); err != nil {
		global.Logger.Debug().Err(err).Msg("error ")
		return nil, false, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while set next page token to cache").
			WithDetails("error append prev to cache").
			WithDetails(err.Error())
	}

	return prev, !api.IsLastPageToken(next), nil
}
