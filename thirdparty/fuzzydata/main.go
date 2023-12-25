package main

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/model"
	"github.com/oklog/ulid"
)

func main() {
	funMap := template.FuncMap{}
	funMap["add"] = func(x, y int) int { return x + y }

	tmpls := template.New("").Funcs(funMap)
	tmpls, err := tmpls.ParseGlob("./template/*.gotmpl")

	if err != nil {
		fmt.Fprintf(os.Stderr, "error while .ParseGlob: %v", err)
		os.Exit(1)
	}

	outputfolder := "./output"

	users := NewUsers(NUM_USER)
	apikeys := NewApiKeys(APIs, users.Item)
	jobs := NewJobs(MAX_NUM_JOBS, APIs, users.Item)
	jobs.Item[0] = JobItem{
		Id:         1,
		Owner:      TEST_ADMIN_USER_UID,
		ULID:       ulid.MustNew(ulid.Timestamp(time.Now()), crand.Reader).String(),
		Status:     model.JobStatusDone,
		SrcAPIName: "GNews",
		SrcApiId:   2,
		SrcQuery:   "covid",
		LlmAPIName: "Cohere",
		LlmApiId:   6,
		LlmQuery:   `{"is_test_data": true}`,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	jids := make([]int64, 0, len(jobs.Item))
	for _, j := range jobs.Item {
		jids = append(jids, int64(j.Id))
	}

	news := NewRandomNewsItems(NUM_NEWSS)
	nids := make([]int64, 0, len(news))
	for i := range news {
		nids = append(nids, int64(i+1))
	}

	nj := NewRandomNewsJob(NUM_NEWSS*AVG_EDGE, jids, nids)

	df, err := NewEmbedding("./data/data.csv", 1536)
	if err != nil {
		panic(err)
	}

	embdmodel := make([]any, df.rowNum)
	for i := 0; i < df.rowNum; i++ {
		embdmodel[i] = "embed-multilingual-light-v3.0"
	}
	df.SetRowAttr("model", embdmodel)

	nid := make([]any, df.rowNum)
	for i := 0; i < df.rowNum; i++ {
		nid[i] = i + 1
	}
	rand.Shuffle(len(nid), func(i, j int) { nid[i], nid[j] = nid[j], nid[i] })
	df.SetRowAttr("nid", nid)

	for i := 0; i < df.rowNum; i++ {
		nj = append(nj, NewsJobItem{JobId: 1, NewsId: int64(df.rowAttr["nid"][i].(int))})
	}

	tasks := []struct {
		TemplateName string
		Data         any
		DownCMD      []string
	}{
		{
			TemplateName: "000011_add_apis.up.sql.gotmpl",
			Data:         APIs,
			DownCMD: []string{
				"DELETE FROM apis;",
				"ALTER SEQUENCE apis_id_seq RESTART WITH 1;",
			},
		},
		{
			TemplateName: "000012_add_endpoints.up.sql.gotmpl",
			Data:         Endpoints,
			DownCMD: []string{
				"DELETE FROM endpoints;",
				"ALTER SEQUENCE endpoints_id_seq RESTART WITH 1;",
			},
		},
		{
			TemplateName: "000013_add_test_user.up.sql.gotmpl",
			Data:         users,
			DownCMD: []string{
				"DELETE FROM users;",
			},
		},
		{
			TemplateName: "000014_add_test_apikey.up.sql.gotmpl",
			Data:         apikeys,
			DownCMD: []string{
				"DELETE FROM apikeys;",
				"ALTER SEQUENCE apikeys_id_seq RESTART WITH 1;",
			},
		},
		{
			TemplateName: "000015_add_test_job.up.sql.gotmpl",
			Data:         jobs,
			DownCMD: []string{
				"DELETE FROM jobs;",
				"ALTER SEQUENCE jobs_id_seq RESTART WITH 1;",
			},
		},
		{
			TemplateName: "000017_add_test_news.up.sql.gotmpl",
			Data:         news,
			DownCMD: []string{
				"DELETE FROM news;",
				"ALTER SEQUENCE news_id_seq RESTART WITH 1;",
			},
		},
		{
			TemplateName: "000018_add_test_newsjob.up.sql.gotmpl",
			Data:         nj,
			DownCMD: []string{
				"DELETE FROM newsjobs;",
				"ALTER SEQUENCE newsjobs_id_seq RESTART WITH 1;",
			},
		},
		{
			TemplateName: "000019_add_test_embedding.up.sql.gotmpl",
			Data:         df,
			DownCMD: []string{
				"DELETE FROM embeddings;",
				"ALTER SEQUENCE embeddings_id_seq RESTART WITH 1;",
			},
		},
	}

	for _, task := range tasks {
		fn := strings.TrimSuffix((task.TemplateName), ".gotmpl")
		fl, err := os.Create(outputfolder + "/" + fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while create file %s: %v", fn, err)
			os.Exit(1)
		}

		defer fl.Close()
		err = tmpls.ExecuteTemplate(fl, task.TemplateName, task.Data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while .ExecuteTemplate %s: %v\n", task.TemplateName, err)
			os.Exit(1)
		}

		if len(task.DownCMD) > 0 {
			fn = strings.Replace(fn, ".up.", ".down.", 1)
			fl, err := os.Create(outputfolder + "/" + fn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error while create file %s: %v", fn, err)
				os.Exit(1)
			}

			defer fl.Close()
			for _, line := range task.DownCMD {
				fl.WriteString(line + "\n")
			}
		}
	}
	os.Exit(0)
}
