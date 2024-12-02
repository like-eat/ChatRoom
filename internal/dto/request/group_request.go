package request

type GroupRequest struct {
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	Notice  string `json:"notice"`
	AddMode bool   `json:"add_mode"`
	Avatar  string `json:"avatar"`
}
