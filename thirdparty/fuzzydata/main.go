package main

import (
	"time"
)

var TIME_MIN, _ = time.Parse(time.DateOnly, "2020-01-01")
var TIME_MAX = time.Now().UTC()

// func main() {
// 	funMap := template.FuncMap{}
// 	funMap["add"] = func(x, y int) int { return x + y }

// 	tmpls := template.New("").Funcs(funMap)
// 	tmpls, err := tmpls.ParseGlob("./template/*.gotmpl")

// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "error while .ParseGlob: %v", err)
// 		os.Exit(1)
// 	}

// 	outputfolder := "./output"

// 	users := NewUsers(50)
// 	apikeys := NewApiKeys(APIs, users.Item)
// 	jobs := NewJobs(50, APIs, users.Item)

// 	tasks := []struct {
// 		TemplateName string
// 		Data         any
// 		DownCMD      []string
// 	}{
// 		{
// 			TemplateName: "000011_add_apis.up.sql.gotmpl",
// 			Data:         APIs,
// 			DownCMD: []string{
// 				"DELETE FROM apis;",
// 				"ALTER SEQUENCE apis_id_seq RESTART WITH 1;",
// 			},
// 		},
// 		{
// 			TemplateName: "000012_add_endpoints.up.sql.gotmpl",
// 			Data:         Endpoints,
// 			DownCMD: []string{
// 				"DELETE FROM endpoints;",
// 				"ALTER SEQUENCE endpoints_id_seq RESTART WITH 1;",
// 			},
// 		},
// 		{
// 			TemplateName: "000013_add_test_user.up.sql.gotmpl",
// 			Data:         users,
// 			DownCMD: []string{
// 				"DELETE FROM users;",
// 			},
// 		},
// 		{
// 			TemplateName: "000014_add_test_apikey.up.sql.gotmpl",
// 			Data:         apikeys,
// 			DownCMD: []string{
// 				"DELETE FROM apikeys;",
// 				"ALTER SEQUENCE apikeys_id_seq RESTART WITH 1;",
// 			},
// 		},
// 		{
// 			TemplateName: "000015_add_test_job.up.sql.gotmpl",
// 			Data:         jobs,
// 			DownCMD: []string{
// 				"DELETE FROM jobs;",
// 				"ALTER SEQUENCE jobs_id_seq RESTART WITH 1;",
// 			},
// 		},
// 	}

// 	for _, task := range tasks {
// 		fn := strings.TrimSuffix((task.TemplateName), ".gotmpl")
// 		fl, err := os.Create(outputfolder + "/" + fn)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "error while create file %s: %v", fn, err)
// 			os.Exit(1)
// 		}

// 		defer fl.Close()
// 		err = tmpls.ExecuteTemplate(fl, task.TemplateName, task.Data)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "error while .ExecuteTemplate %s: %v", task.TemplateName, err)
// 			os.Exit(1)
// 		}

// 		if len(task.DownCMD) > 0 {
// 			fn = strings.Replace(fn, ".up.", ".down.", 1)
// 			fl, err := os.Create(outputfolder + "/" + fn)
// 			if err != nil {
// 				fmt.Fprintf(os.Stderr, "error while create file %s: %v", fn, err)
// 				os.Exit(1)
// 			}

// 			defer fl.Close()
// 			for _, line := range task.DownCMD {
// 				fl.WriteString(line + "\n")
// 			}
// 		}
// 	}
// 	os.Exit(0)
// }
