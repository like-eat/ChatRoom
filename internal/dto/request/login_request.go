package request

type LoginRequest struct {
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}
