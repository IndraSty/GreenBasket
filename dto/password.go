package dto

type PasswordReq struct {
	OldPassword     string `json:"old_password" valid:"minstringlength(8)"`
	NewPassword     string `json:"new_password" valid:"required,minstringlength(8)"`
	RewritePassword string `json:"rewrite_password" valid:"required,minstringlength(8)"`
}
