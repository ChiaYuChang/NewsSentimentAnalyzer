package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

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

	pcid := chi.URLParam(req, "pcid")
	var pf pageform.PreviewPostForm
	repo.FormDecoder.Decode(&pf, req.Form)

	res, err := repo.Cache.JSONGet(pcid, ".news_item")
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while get cache").
			WithDetails(err.Error())
		resp.Error = &pageform.PreviewError{
			Code:        ecErr.HttpStatusCode,
			Message:     ecErr.Message,
			Detail:      ecErr.Details,
			RedirectURL: global.AppVar.App.RoutePattern.ErrorPage["internal-server-error"],
		}

		b, _ := json.Marshal(resp)
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	sort.Sort(sort.StringSlice(pf.Item))
	var prev []api.NewsPreview
	err = json.Unmarshal(res.([]byte), &prev)
	if err != nil {
		ecErr := ec.MustGetEcErr(ec.ECBadRequest).
			WithDetails("error while unmarshaling cache").
			WithDetails(err.Error())
		resp.Error = &pageform.PreviewError{
			Code:        ecErr.HttpStatusCode,
			Message:     ecErr.Message,
			Detail:      ecErr.Details,
			RedirectURL: global.AppVar.App.RoutePattern.ErrorPage["internal-server-error"],
		}

		b, _ := json.Marshal(resp)
		w.WriteHeader(ecErr.HttpStatusCode)
		w.Write(b)
		return
	}

	var selectedPrev []api.NewsPreview
	if pf.SelectAll {
		selectedPrev = prev
	} else {
		selectedPrev := make([]api.NewsPreview, 0, len(prev))
		for _, p := range prev {
			_, ok := sort.Find(len(pf.Item), func(i int) int {
				return strings.Compare(p.Id.String(), pf.Item[i])
			})
			if ok {
				selectedPrev = append(selectedPrev, p)
			}
		}
	}

	// for _, p := range selectedPrev {
	// 	href, err := url.Parse(p.Link)
	// 	if err != nil {
	// 		global.Logger.Error().
	// 			Err(err).
	// 			Str("link", p.Link).
	// 			Msg("error while parsing link")
	// 		continue
	// 	}
	// 	guid := parser.GetDefaultParser().ToGUID(href)
	// 	lang, ok := global.LanguageDetector().DetectLanguageOf(p.Title)
	// 	if !ok {
	// 		lang = lingua.Unknown
	// 	}
	// 	req := p.ToNewsCreateRequest(guid, lang.IsoCode639_1().String(), href.Host)
	// 	repo.Service.News().Create(context.Background())
	// }

	global.Logger.Info().Int("n", len(selectedPrev)).Msg("OK")
	resp.RedirectURL = fmt.Sprintf("/%s%s", repo.Version, global.AppVar.App.RoutePattern.Page["welcome"])
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

	if cq, ecErr = getCacheQuery(repo, pcid); ecErr == nil {
		if resp, ecErr = fetchNextPage(cq); ecErr == nil {
			prev, hasNext, ecErr = getPreviewsAndUpdateCache(repo, pcid, resp)
		}
	}

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
		return nil, ec.MustGetEcErr(ec.ECServerError).
			WithDetails("error while .Do").
			WithDetails(err.Error())
	}
	return resp, nil
}

func getPreviewsAndUpdateCache(repo APIRepo, pcid string, resp api.Response) ([]api.NewsPreview, bool, *ec.Error) {
	// append prev to cache
	next, prev := resp.ToNewsItemList()

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
