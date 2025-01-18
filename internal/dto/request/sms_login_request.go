package request

type SmsLoginRequest struct {
	Telephone string `json:"telephone"`
	SmsCode   string `json:"sms_code"`
}
