package service

type UserCleanUpRequest struct{}

func (r UserCleanUpRequest) RequestName() string {
	return "admin-user-clean-up-req"
}

type UserHardDeleteRequest struct{}

func (r UserHardDeleteRequest) RequestName() string {
	return "admin-user-hard-delete-req"
}

type CreateAPIRequest struct {
	Name string `validate:"required"`
	Type string `validate:"required,api_type"`
}

func (r CreateAPIRequest) RequestName() string {
	return "admin-api-create-req"
}

type JobCleanUpRequest struct{}

func (r JobCleanUpRequest) RequestName() string {
	return "admin-job-clean-up-req"
}
