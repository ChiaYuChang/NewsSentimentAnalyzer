package validator_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/validator"
	val "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateEnmus(t *testing.T) {
	type testCase struct {
		validator.Validator
		Name      string
		OKCase    []string
		ErrorCase []string
	}

	tcs := []testCase{
		{
			Name:      "Role",
			Validator: validator.EnmusRole,
			OKCase:    []string{"user", "admin"},
			ErrorCase: []string{"visitor", "User", "ADMIN"},
		},
		{
			Name:      "Job Status",
			Validator: validator.EnmusJobStatus,
			OKCase:    []string{"created", "running", "done", "failure", "canceled"},
			ErrorCase: []string{"unknown", "panding"},
		},
		{
			Name:      "API type",
			Validator: validator.EnmusApiType,
			OKCase:    []string{"language_model", "source"},
			ErrorCase: []string{"unknown", "database"},
		},
		{
			Name:      "Event type",
			Validator: validator.EnmusEventType,
			OKCase:    []string{"sign-in", "sign-out", "authorization", "api-key", "query"},
			ErrorCase: []string{"unknown", "surfing"},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("case %d-%s", i+1, tc.Name),
			func(t *testing.T) {
				var err error
				validate := val.New()
				err = validate.RegisterValidation(
					tc.Validator.Tag(),
					tc.Validator.ValFun(),
				)
				require.NoError(t, err)

				for _, e := range tc.OKCase {
					err = validate.Var(e, tc.Tag())
					require.NoError(t, err)
				}

				for _, e := range tc.ErrorCase {
					err = validate.Var(e, tc.Tag())
					require.Error(t, err)
				}
			},
		)
	}
}

func TestValidateEnmusList(t *testing.T) {
	type testCase struct {
		validator.Validator
		Name       string
		OKCase     []string
		ErrorCases [][]string
	}

	tcs := []testCase{
		{
			Name:      "Role",
			Validator: validator.EnmusRole,
			OKCase:    []string{"user", "admin"},
			ErrorCases: [][]string{
				{"user", "visitor", "admin"},
				{"user", "user", "admin", "admin", "admin"},
			},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("case %d-%s", i+1, tc.Name),
			func(t *testing.T) {
				var err error
				validate := val.New()
				err = validate.RegisterValidation(
					tc.Validator.Tag(),
					tc.Validator.ValFun(),
				)
				require.NoError(t, err)

				err = validate.Var(tc.OKCase, fmt.Sprintf("max=%d,%s", 5, tc.Tag()))
				require.NoError(t, err)

				for _, ec := range tc.ErrorCases {
					err = validate.Var(ec, fmt.Sprintf("max=%d,%s", 3, tc.Tag()))
					require.Error(t, err)
				}
			},
		)
	}
}

func TestValidateLeakyPassword(t *testing.T) {
	validate := val.New()
	leakeyPwdVal := validator.NewPasswordValidator(false, 0, 80, 0, 0, 0, 0)
	err := validator.RegisterValidator(validate, leakeyPwdVal)
	require.NoError(t, err)
	defer validator.RegisterValidator(validate, validator.NewDefaultPasswordValidator())

	type testCase struct {
		Name     string
		Password string
		IsValid  bool
	}

	tcs := []testCase{
		{"OK long", "Xm2pIcsyYfmpyM51KgqrdNtQQgOgQsuoaaKmFsQur7SgLox0MDFlKK5EG0vywAlfLGZNX5RBtgWLlveH", true},
		{"OK empty", "", true},
		{"too long", "Xm2pIcsyYfmpyM51KgqrdNtQQgOgQsuoaaKmFsQur7SgLox0MDFlKK5EG0vywAlfLGZNX5RBtgWLlveHs", false},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-%s", i, tc.Name),
			func(t *testing.T) {
				err := validate.Var(tc.Password, leakeyPwdVal.Tag())
				if tc.IsValid {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
				}
			},
		)
	}
}

func TestValidatePassword(t *testing.T) {
	type testCase struct {
		Name     string
		Password string
		IsValid  bool
	}

	tcs := []testCase{
		{"OK long", "h&tf97hrZCukDg*", true},
		{"OK short", "hAA6#6%y", true},
		{"OK", "hA6%yAAA", true},
		{"empty", "", false},
		{"too short", "abc", false},
		{"too long", "xKCn9uL5D&iCXf&ABf@ATFPCq#B8o6X", false},
		{"without Upper", "9bqgtny$ps", false},
		{"without Lower", "ST6QYP#5H4", false},
		{"without Digit", "gPsDFrDU%b", false},
		{"without Special", "YX4aCFU7hg", false},
		{"with not ASCII", "亂數密碼", false},
	}

	validate := val.New()
	pwdVal := validator.NewDefaultPasswordValidator()

	err := validate.RegisterValidation(pwdVal.Tag(), pwdVal.ValFun())
	require.NoError(t, err)
	for i := range tcs {
		tc := tcs[i]
		t.Run(
			fmt.Sprintf("Case %d-%s", i, tc.Name),
			func(t *testing.T) {
				err := validate.Var(tc.Password, pwdVal.Tag())
				if tc.IsValid {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
				}
			},
		)
	}
}

func TestPast(t *testing.T) {
	validate := val.New()
	validator.RegisterValidator(
		validate,
		validator.NotBeforeNow,
	)

	require.NoError(t, validate.Var(time.Now(), validator.NotBeforeNow.Tag()))
	require.NoError(t, validate.Var(time.Now().Add(-1*time.Hour), validator.NotBeforeNow.Tag()))
	require.Error(t, validate.Var(time.Now().Add(1*time.Hour), validator.NotBeforeNow.Tag()))
}

func TestValidateUUID(t *testing.T) {
	validate := val.New()
	validator.RegisterUUID(validate)

	type testCase struct {
		Name         string
		UID          uuid.UUID
		Tags         []string
		ExpectResult []bool
	}

	uid4 := uuid.New()
	uid3 := uuid.NewMD5(uuid.Nil, []byte("UUID version 3"))

	tcs := []testCase{
		{
			Name:         "UUID ver4",
			UID:          uid4,
			Tags:         []string{"uuid", "uuid3", "uuid4", "not_uuid_nil"},
			ExpectResult: []bool{true, false, true, true},
		},
		{
			Name:         "UUID ver3",
			UID:          uid3,
			Tags:         []string{"uuid", "uuid3", "uuid4", "not_uuid_nil"},
			ExpectResult: []bool{true, true, false, true},
		},
		{
			Name:         "uuid.Nil",
			UID:          uuid.Nil,
			Tags:         []string{"uuid", "uuid3", "uuid4", "not_uuid_nil"},
			ExpectResult: []bool{true, false, false, false},
		},
	}

	for i := range tcs {
		tc := tcs[i]
		t.Run(
			tc.Name,
			func(t *testing.T) {
				for j, tag := range tc.Tags {
					err := validate.Var(tc.UID, tag)
					if tc.ExpectResult[j] {
						require.NoError(t, err)
					} else {
						require.Error(t, err)
					}
				}
			},
		)
	}
}
