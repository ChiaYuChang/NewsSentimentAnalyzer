package testtool

import (
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/model"
	"github.com/ChiaYuChang/NewsSentimentAnalyzer/internal/server/service"
	"github.com/golang/mock/gomock"
)

type UserCreateReqMatcher struct {
	param *model.CreateUserParams
}

func NewUserCreateReqMatcher(req *service.UserCreateRequest) (gomock.Matcher, error) {
	if param, err := req.ToParams(); err != nil {
		return nil, err
	} else {
		return UserCreateReqMatcher{
			param: param,
		}, nil
	}
}

func (m UserCreateReqMatcher) Matches(x interface{}) bool {
	if p, ok := x.(*model.CreateUserParams); !ok {
		return false
	} else {
		p.Password = m.param.Password
		return gomock.Eq(m.param).Matches(p)
	}
}

func (m UserCreateReqMatcher) String() string {
	return "matcher ofr user create request"
}

type UserUpdatePasswordReqMatcher struct {
	param *model.UpdatePasswordParams
}

func NewUserUpdatePasswordReqMatcher(req *service.UserUpdatePasswordRequest) (gomock.Matcher, error) {
	if param, err := req.ToParams(); err != nil {
		return nil, err
	} else {
		return UserUpdatePasswordReqMatcher{
			param: param,
		}, nil
	}
}

func (m UserUpdatePasswordReqMatcher) Matches(x interface{}) bool {
	if p, ok := x.(*model.UpdatePasswordParams); !ok {
		return false
	} else {
		p.Password = m.param.Password
		return gomock.Eq(m.param).Matches(p)
	}
}

func (m UserUpdatePasswordReqMatcher) String() string {
	return "matcher ofr user update password request"
}
