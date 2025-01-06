package request

type UpdateUserInfoRequest struct {
	Uuid      string `json:"uuid"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Birthday  string `json:"birthday"`
	Signature string `json:"signature"`
	Avatar    string `json:"avatar"`
}
