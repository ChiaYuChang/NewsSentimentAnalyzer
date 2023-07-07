package nullable_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/nullable"
	"github.com/stretchr/testify/require"
)

func TestNullableString(t *testing.T) {
	type TestStruct struct {
		Field nullable.String[string] `json:"field"`
	}

	type testCast struct {
		Name     string
		JsonStr  string
		IsValid  bool
		ExpValue string
	}

	tcs := []testCast{
		{
			Name:     "null",
			JsonStr:  `{"field": null}`,
			IsValid:  false,
			ExpValue: "",
		},
		{
			Name:     "OK, string null",
			JsonStr:  `{"field": "null"}`,
			IsValid:  true,
			ExpValue: "null",
		},
		{
			Name:     "OK",
			JsonStr:  `{"field": "hello"}`,
			IsValid:  true,
			ExpValue: "hello",
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-%s", i+1, tc.Name),
			func(t *testing.T) {
				var jsonObj TestStruct
				err := json.Unmarshal([]byte(tc.JsonStr), &jsonObj)
				require.NoError(t, err)

				require.NotNil(t, jsonObj)
				require.Equal(t, tc.IsValid, jsonObj.Field.Valid)
				require.Equal(t, tc.ExpValue, jsonObj.Field.Value)

				jsonBytes, err := json.Marshal(&jsonObj)
				require.NoError(t, err)
				if tc.IsValid {
					require.Equal(t, string(jsonBytes), fmt.Sprintf(`{"field":"%s"}`, tc.ExpValue))
				} else {
					require.Equal(t, string(jsonBytes), fmt.Sprintf(`{"field":%s}`, "null"))
				}
			},
		)
	}
}

func TestNullableStringT(t *testing.T) {
	type JSON string

	type TestStruct struct {
		Field nullable.String[JSON] `json:"field"`
	}

	type testCast struct {
		Name     string
		JsonStr  string
		Err      error
		IsValid  bool
		ExpValue JSON
	}

	tcs := []testCast{
		{
			Name:     "null",
			JsonStr:  `{"field": null}`,
			IsValid:  false,
			ExpValue: JSON(""),
		},
		{
			Name:     "OK, string null",
			JsonStr:  `{"field": "null"}`,
			IsValid:  true,
			ExpValue: JSON("null"),
		},
		{
			Name:     "OK",
			JsonStr:  `{"field": "hello"}`,
			IsValid:  true,
			ExpValue: JSON("hello"),
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-%s", i+1, tc.Name),
			func(t *testing.T) {
				var jsonObj TestStruct
				err := json.Unmarshal([]byte(tc.JsonStr), &jsonObj)
				if tc.Err == nil {
					require.NoError(t, err)
					require.Equal(t, tc.IsValid, jsonObj.Field.Valid)
					require.Equal(t, tc.ExpValue, jsonObj.Field.Value)

					jsonBytes, err := json.Marshal(&jsonObj)
					require.NoError(t, err)
					if tc.IsValid {
						require.Equal(t, string(jsonBytes), fmt.Sprintf(`{"field":"%s"}`, tc.ExpValue))
					} else {
						require.Equal(t, string(jsonBytes), fmt.Sprintf(`{"field":%s}`, "null"))
					}
				} else {
					require.ErrorAs(t, strconv.ErrSyntax, err)
				}
			},
		)
	}
}

func TestNullableInt(t *testing.T) {
	type TestStruct struct {
		Field nullable.Int[int] `json:"field"`
	}

	type testCast struct {
		Name     string
		JsonStr  string
		Err      error
		IsValid  bool
		ExpValue int
	}

	tcs := []testCast{
		{
			Name:     "null",
			JsonStr:  `{"field": null}`,
			IsValid:  false,
			ExpValue: 0,
		},
		{
			Name:     "OK",
			JsonStr:  `{"field": 10}`,
			IsValid:  true,
			ExpValue: 10,
		},
		{
			Name:    "Not a number",
			JsonStr: `{"field": "ef1"}`,
			Err:     strconv.ErrSyntax,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-%s", i+1, tc.Name),
			func(t *testing.T) {
				var jsonObj TestStruct
				err := json.Unmarshal([]byte(tc.JsonStr), &jsonObj)
				if tc.Err == nil {
					require.NoError(t, err)
					require.Equal(t, tc.IsValid, jsonObj.Field.Valid)
					require.Equal(t, tc.ExpValue, jsonObj.Field.Value)

					jsonBytes, err := json.Marshal(&jsonObj)
					require.NoError(t, err)
					if tc.IsValid {
						require.Equal(t, string(jsonBytes), fmt.Sprintf(`{"field":%d}`, tc.ExpValue))
					} else {
						require.Equal(t, string(jsonBytes), fmt.Sprintf(`{"field":%s}`, "null"))
					}
				} else {
					require.ErrorIs(t, err, tc.Err)
				}
			},
		)
	}
}

func TestNullableIntT(t *testing.T) {
	type RoleType int
	const (
		User  RoleType = 0
		Admin RoleType = 1
	)

	type TestStruct struct {
		Role nullable.Int[RoleType] `json:"field"`
	}

	type testCast struct {
		Name     string
		JsonStr  string
		IsValid  bool
		ExpValue RoleType
	}

	tcs := []testCast{
		{
			Name:     "null",
			JsonStr:  `{"field": null}`,
			IsValid:  false,
			ExpValue: 0,
		},
		{
			Name:     "OK",
			JsonStr:  `{"field": 0}`,
			IsValid:  true,
			ExpValue: User,
		},
		{
			Name:     "OK",
			JsonStr:  `{"field": 1}`,
			IsValid:  true,
			ExpValue: Admin,
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-%s", i+1, tc.Name),
			func(t *testing.T) {
				var jsonObj TestStruct
				err := json.Unmarshal([]byte(tc.JsonStr), &jsonObj)
				require.NoError(t, err)
				require.Equal(t, tc.IsValid, jsonObj.Role.Valid)
				require.Equal(t, tc.ExpValue, jsonObj.Role.Value)

				jsonBytes, err := json.Marshal(&jsonObj)
				require.NoError(t, err)
				if tc.IsValid {
					require.Equal(t, string(jsonBytes), fmt.Sprintf(`{"field":%d}`, tc.ExpValue))
				} else {
					require.Equal(t, string(jsonBytes), fmt.Sprintf(`{"field":%s}`, "null"))
				}
			},
		)
	}
}
