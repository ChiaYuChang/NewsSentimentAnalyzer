package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"text/template"
	"time"

	rg "github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/randanGenerator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var TIME_MIN, _ = time.Parse(time.DateOnly, "2020-01-01")
var TIME_MAX = time.Now().UTC()

type Sampler[T any] struct {
	x []T
	p []float64
	c []float64
}

func NewSampler[T any](x []T, weight []float64) Sampler[T] {
	if weight == nil {
		weight = make([]float64, len(x))
		for i := range weight {
			weight[i] = 1.0 / float64(len(x))
		}
	}
	c := make([]float64, len(x))
	for i, w := range weight {
		c[i] = w
		if i > 0 {
			c[i] += c[i-1]
		}
	}

	return Sampler[T]{x, weight, c}
}

func (s Sampler[T]) Get() T {
	r := rand.Float64()
	return s.x[sort.SearchFloat64s(s.c, r)]
}

func (s Sampler[T]) GetN(n int) []T {
	rs := make([]T, n)
	for i := 0; i < n; i++ {
		rs[i] = s.Get()
	}
	return rs
}

const (
	APITypeSource string = "source"
	APITypeLLM    string = "language_model"
)

type API struct {
	Id          int
	Name        string
	Type        string
	Image       string
	Icon        string
	DocumentURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Probability float64
}

var APIs = []API{
	{
		Id:          1,
		Name:        "NEWSDATA.IO",
		Type:        APITypeSource,
		Image:       "logo_NEWSDATA.IO.png",
		Icon:        "favicon_NEWSDATA.IO.png",
		DocumentURL: "https://newsdata.io/documentation/",
		CreatedAt:   TIME_MIN,
		UpdatedAt:   TIME_MIN,
		Probability: 0.5,
	},
	{
		Id:          2,
		Name:        "GNews",
		Type:        APITypeSource,
		Image:       "logo_GNews.png",
		Icon:        "favicon_GNews.ico",
		DocumentURL: "https://gnews.io/docs/v4",
		CreatedAt:   TIME_MIN,
		UpdatedAt:   TIME_MIN,
		Probability: 0.5,
	},
	{
		Id:          3,
		Name:        "NEWS API",
		Type:        APITypeSource,
		Image:       "logo_NEWS_API.png",
		Icon:        "favicon_NEWS_API.ico",
		DocumentURL: "https://newsapi.org/docs/",
		CreatedAt:   TIME_MIN,
		UpdatedAt:   TIME_MIN,
		Probability: 0.5,
	},
	{
		Id:          4,
		Name:        "OpenAI",
		Type:        APITypeLLM,
		Image:       "logo_ChatGPT.svg",
		Icon:        "favicon_ChatGPT.ico",
		CreatedAt:   TIME_MIN,
		UpdatedAt:   TIME_MIN,
		DocumentURL: "https://openai.com/blog/introducing-chatgpt-and-whisper-apis",
		Probability: 1.0,
	},
}

type Endpoint struct {
	Name         string
	APIId        int
	TemplateName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var Endpoints = []Endpoint{
	{
		Name:         "Latest News",
		APIId:        1,
		TemplateName: "NEWSDATA.IO-latest_news.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "News Archive",
		APIId:        1,
		TemplateName: "NEWSDATA.IO-news_archive.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "News Sources",
		APIId:        1,
		TemplateName: "NEWSDATA.IO-news_sources.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "Search",
		APIId:        2,
		TemplateName: "GNews-search.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "Top Headlines",
		APIId:        2,
		TemplateName: "GNews-top_headlines.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "Everything",
		APIId:        3,
		TemplateName: "NewsAPI-everything.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "Top Headlines",
		APIId:        3,
		TemplateName: "NewsAPI-top_headlines.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
	{
		Name:         "Sources",
		APIId:        3,
		TemplateName: "NewsAPI-sources.gotmpl",
		CreatedAt:    TIME_MIN,
		UpdatedAt:    TIME_MIN,
	},
}

var UserRole = []string{
	"user",
	"admin",
}

type Password []byte

func (pws Password) String() string {
	return string(pws)
}

type User struct {
	Id        uuid.UUID
	Password  Password
	FirstName string
	LastName  string
	Role      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUsers(n int) []User {
	n += 2
	us := make([]User, n)

	var rts []time.Time
	rts = rg.GenRdnTimes(2, TIME_MIN, TIME_MAX)
	us[0] = User{
		Id:        uuid.New(),
		Password:  rg.Must(bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)),
		FirstName: rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		LastName:  rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		Role:      UserRole[0],
		Email:     "test@example.com",
		CreatedAt: rts[0],
		UpdatedAt: rts[1],
	}

	rts = rg.GenRdnTimes(2, TIME_MIN, TIME_MAX)
	us[1] = User{
		Id:        uuid.New(),
		Password:  rg.Must(bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)),
		FirstName: rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		LastName:  rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
		Role:      UserRole[1],
		Email:     "admin@example.com",
		CreatedAt: rts[0],
		UpdatedAt: rts[1],
	}

	rs := NewSampler(UserRole, []float64{0.99, 0.01})
	for i := 2; i < n; i++ {
		rawPwd, _ := rg.GenRdmPwd(8, 30, 1, 1, 1, 1)
		encPwd, _ := bcrypt.GenerateFromPassword([]byte(rawPwd), bcrypt.DefaultCost)
		rts := rg.GenRdnTimes(2, TIME_MIN, TIME_MAX)
		u := User{
			Id:        uuid.New(),
			Password:  encPwd,
			FirstName: rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
			LastName:  rg.Must(rg.Alphabet.GenRdmString(3 + rand.Intn(17))),
			Role:      rs.Get(),
			Email:     rg.Must(rg.GenRdmEmail(rg.AlphaNum, rg.AlphabetLower)),
			CreatedAt: rts[0],
			UpdatedAt: rts[1],
		}
		us[i] = u
	}
	return us
}

type APIKey struct {
	Id        int
	Owner     uuid.UUID
	APIId     int
	Key       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewApiKeys(apis []API, users []User) []APIKey {
	apikeys := []APIKey{}

	kl := []int{16, 32, 52}

	s := NewSampler([]rune{'0', 'x', '-'}, []float64{0.6, 0.35, 0.05})
	for _, u := range users {
		for _, a := range apis {
			if rand.Float64() < a.Probability {
				apikeys = append(apikeys, APIKey{
					Owner: u.Id,
					APIId: a.Id,
					Key:   string(s.GetN(kl[rand.Intn(len(kl))])),
				})
			}
		}
	}

	rand.Shuffle(len(apikeys), func(i, j int) {
		apikeys[i], apikeys[j] = apikeys[j], apikeys[i]
	})

	rts := rg.GenRdnTimes(len(apikeys), TIME_MIN, TIME_MAX)
	for i := range apikeys {
		apikeys[i].Id = i + 1
		apikeys[i].CreatedAt = rts[i]
		apikeys[i].UpdatedAt = rg.GenRdnTime(rts[i], TIME_MAX)
	}

	return apikeys
}

var JobStatus = []string{
	"created",
	"running",
	"done",
	"failure",
	"canceled",
}

type Job struct {
	Id        int
	Owner     uuid.UUID
	Status    string
	SrcApiId  int
	SrcQuery  string
	LlmApiId  int
	LlmQuery  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt NullableTime
}

type NullableTime struct {
	time.Time
	Valid bool
}

func (t NullableTime) Format(layout string) string {
	if !t.Valid {
		return "null"
	}
	return "'" + t.Time.Format(layout) + "'"
}

func NewJobs(maxJobN int, apis []API, users []User) []Job {
	srcApis := []int{1, 2, 3}
	llmApis := []int{4}

	jobs := make([]Job, 0, maxJobN*len(users))

	jss := NewSampler(JobStatus, []float64{0.3, 0.3, 0.3, 0.07, 0.03})
	for i := 0; i < len(users); i++ {
		n := rand.Intn(maxJobN-20) + 20
		owner := users[i]
		// fmt.Println(owner.Id, n)
		for j := 0; j < n; j++ {
			srcApi := apis[srcApis[rand.Intn(len(srcApis))]-1]
			llmApi := apis[llmApis[rand.Intn(len(llmApis))]-1]
			jobs = append(jobs, Job{
				Owner:    owner.Id,
				Status:   jss.Get(),
				SrcApiId: srcApi.Id,
				SrcQuery: rg.Must(rg.AlphaNum.GenRdmString(rand.Intn(10) + 20)),
				LlmApiId: llmApi.Id,
				LlmQuery: "{}",
			})
		}
	}

	rand.Shuffle(len(jobs), func(i, j int) { jobs[i], jobs[j] = jobs[j], jobs[i] })
	cts := rg.GenRdnTimes(len(jobs), TIME_MIN, TIME_MAX)
	for i := 0; i < len(jobs); i++ {
		var ct, ut time.Time
		// var dt NullableTime
		ct = cts[i]
		switch jobs[i].Status {
		default:
			ut = ct
		case "running", "done", "failure":
			ut = rg.GenRdnTime(cts[i], TIME_MAX)
		case "canceled":
			ut = rg.GenRdnTime(cts[i], TIME_MAX)
			// dt = NullableTime{Time: ut, Valid: true}
		}
		jobs[i].Id = i + 1
		jobs[i].CreatedAt = cts[i]
		jobs[i].UpdatedAt = ut
		// jobs[i].DeletedAt = dt
	}
	return jobs
}

func main() {
	tmpls, err := template.ParseGlob("./template/*.gotmpl")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while .ParseGlob: %v", err)
		os.Exit(1)
	}

	users := NewUsers(50)
	fusers, err := os.Create("000013_add_user.up.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while .Create: %v", err)
		os.Exit(1)
	}
	defer fusers.Close()
	err = tmpls.ExecuteTemplate(fusers, "user.gotmpl", users)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while .ExecuteTemplate: %v", err)
		os.Exit(1)
	}
	apikeys := NewApiKeys(APIs, users)
	fapikeys, err := os.Create("000014_apikey.up.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while .Create: %v", err)
		os.Exit(1)
	}
	defer fapikeys.Close()
	err = tmpls.ExecuteTemplate(fapikeys, "apikey.gotmpl", apikeys)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while .ExecuteTemplate: %v", err)
		os.Exit(1)
	}

	jobs := NewJobs(50, APIs, users)
	fjobs, err := os.Create("000015_jobs.up.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while .Create: %v", err)
		os.Exit(1)
	}
	defer fjobs.Close()
	err = tmpls.ExecuteTemplate(fjobs, "job.gotmpl", jobs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while .ExecuteTemplate: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
