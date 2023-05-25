package global

var Secrets Secret

// func init() {
// 	f, err := os.Open("./secret.json")
// 	if err != nil {
// 		panic(fmt.Sprintf("error while opening secret.json: %s", err.Error()))
// 	}
// 	defer f.Close()

// 	srcts, err := io.ReadAll(f)
// 	if err != nil {
// 		panic(fmt.Sprintf("error while reading secrets: %s", err.Error()))
// 	}

// 	err = json.Unmarshal(srcts, &Secrets)
// 	if err != nil {
// 		panic(fmt.Sprintf("error while unmarshal secrets: %s", err.Error()))
// 	}
// }

type Secret struct {
	Database []Database `json:"db"`
	API      []APIKey   `json:"api"`
}

type Database struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type APIKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}
