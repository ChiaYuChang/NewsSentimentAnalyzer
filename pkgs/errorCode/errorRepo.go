package errorcode

var defaultErrorRepo ErrorRepo

func init() {
	if errRepo, err := NewDefaultErrorRepo(); err != nil {
		panic(err.Error())
	} else {
		defaultErrorRepo = errRepo
	}
}

func NewDefaultErrorRepo() (ErrorRepo, error) {
	return NewErrorRepo(
		WithSuccess(),
		WithClientErr(),
		WithServerErr())
}

func WithOptions(opts ...ErrorRepoOption) error {
	for _, opt := range opts {
		err := opt(defaultErrorRepo)
		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterErr(code ErrorCode, httpStatusCode int, message string) error {
	return defaultErrorRepo.RegisterErr(code, httpStatusCode, message)
}

func RegisterFromErr(code ErrorCode, httpStatusCode int, err error) error {
	return defaultErrorRepo.RegisterErrFromErr(err, code, httpStatusCode)
}

func GetErr(code ErrorCode) (error, bool) {
	return defaultErrorRepo.GetErr(code)
}

func MustGetErr(code ErrorCode) error {
	return defaultErrorRepo.MustGetErr(code)
}

type ErrorRepo map[ErrorCode]*Error

type ErrorRepoOption func(repo ErrorRepo) error

func NewErrorRepo(opts ...ErrorRepoOption) (ErrorRepo, error) {
	repo := make(ErrorRepo)
	repo.RegisterErr(ECCodeHasBeenUsed, 0, "the error code has been used")
	for _, opt := range opts {
		err := opt(repo)
		if err != nil {
			return repo, err
		}
	}
	return repo, nil
}

func (er ErrorRepo) RegisterErr(code ErrorCode, httpStatusCode int, message string) error {
	_, ok := er.GetErr(code)
	if ok {
		return er[ECCodeHasBeenUsed].Clone()
	}
	er[code] = NewError(code, httpStatusCode, message)
	return nil
}

func (er ErrorRepo) RegisterErrFromErr(err error, code ErrorCode, httpStatusCode int) error {
	return er.RegisterErr(code, httpStatusCode, err.Error())
}

func (er ErrorRepo) GetErr(code ErrorCode) (error, bool) {
	err, ok := er[code]
	if !ok {
		return err, false
	}
	return err.Clone(), true
}

func (er ErrorRepo) MustGetErr(code ErrorCode) error {
	e, _ := er.GetErr(code)
	return e
}
