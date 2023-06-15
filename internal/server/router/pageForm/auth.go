package pageform

type AuthInfo struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type ChangePassword struct {
	OldPassword string `form:"old-password"`
	NewPassword string `form:"new-password"`
}

type SignUpInfo struct {
	AuthInfo
	FirstName string `form:"first-name"`
	LastName  string `form:"last-name"`
}
