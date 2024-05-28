package dto

type ValidateOtpReq struct {
	UserID string `json:"user_id"`
	OTP    string `json:"otp"`
}
