package pageform

type AuthInfo struct {
	Email    string `mod:"trim,lcase" form:"email"    validate:"email"`
	Password string `                 form:"password" validate:"password"`
}

type ChangePassword struct {
	OldPassword string `form:"old-password" validate:"password"`
	NewPassword string `form:"new-password" validate:"password"`
}

type SignUpInfo struct {
	AuthInfo
	FirstName string `mod:"trim" form:"first-name" validate:"required"`
	LastName  string `mod:"trim" form:"last-name"  validate:"required"`
}
