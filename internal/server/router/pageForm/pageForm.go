package pageform

type AuthInfo struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type SignUpInfo struct {
	AuthInfo
	FirstName string `form:"first-name"`
	LastName  string `form:"last-name"`
}
