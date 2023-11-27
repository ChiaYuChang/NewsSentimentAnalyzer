package errorcode

import "net/http"

// Register Success to ErrorReop
func WithSuccess() ErrorRepoOption {
	return func(repo ErrorRepo) error {
		for _, e := range []struct {
			code   ErrorCode
			status int
			msg    string
		}{
			{Success, http.StatusOK, "success"},
			{SuccessNoContent, http.StatusNoContent, "no content"},
		} {
			err := repo.RegisterErr(e.code, e.status, e.msg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Register Client side errors to ErrorReop
func WithClientErr() ErrorRepoOption {
	return func(repo ErrorRepo) error {
		for _, e := range []struct {
			code   ErrorCode
			status int
			msg    string
		}{
			{ECBadRequest, http.StatusBadRequest, "the request was unacceptable"},
			{ECUnauthorized, http.StatusUnauthorized, "the request has not been completed because it lacks valid authentication credentials for the requested resource"},
			{ECForbidden, http.StatusForbidden, "the client's identity is known to the server"},
			{ECNotFound, http.StatusNotFound, "the server cannot find the requested resource"},
			{ECRequestTimeout, http.StatusRequestTimeout, "request timeout"},
			{ECTooManyRequests, http.StatusTooManyRequests, "too may requests"},
			{ECConflict, http.StatusConflict, "the request conflicts with the current state of the server"},
			{ECUnsupportedMediaType, http.StatusUnsupportedMediaType, "the media format of the requested data is not supported by the server"},
			{ECUnprocessableContent, http.StatusUnprocessableEntity, "the request could not be process by server"},
			{ECInvalidParams, http.StatusBadRequest, "invalid parameters"},
			{ECGone, http.StatusGone, "the requested content has been permanently deleted from server"},
		} {
			err := repo.RegisterErr(e.code, e.status, e.msg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func WithServerErr() ErrorRepoOption {
	return func(repo ErrorRepo) error {
		for _, e := range []struct {
			code   ErrorCode
			status int
			msg    string
		}{
			{ECServerError, http.StatusInternalServerError, "internal server error"},
			{ECServiceUnavailable, http.StatusServiceUnavailable, "the server is not ready to handle the request"},
		} {
			err := repo.RegisterErr(e.code, e.status, e.msg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func WithPgxError() ErrorRepoOption {
	return func(repo ErrorRepo) error {
		err := repo.RegisterErr(ECPgxError, http.StatusInternalServerError, "pgx error")
		if err != nil {
			return nil
		}
		return nil
	}
}

func WithgRPCError() ErrorRepoOption {
	return func(repo ErrorRepo) error {
		for _, e := range []struct {
			code   ErrorCode
			status int
			msg    string
		}{
			{ECgRPCClientError, http.StatusInternalServerError, "gRPC client error"},
			{ECgRPCServerError, http.StatusInternalServerError, "gRPC server error"},
		} {
			err := repo.RegisterErr(e.code, e.status, e.msg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
