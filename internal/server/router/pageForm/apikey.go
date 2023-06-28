package pageform

type APIKeyPost struct {
	ApiID int16  `form:"api-id"  validate:"min=1"`
	Key   string `form:"api-key" validate:"min=32,max=64"`
}
