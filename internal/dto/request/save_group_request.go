package request

type SaveGroupRequest struct {
	Uuid    string `json:"uuid"`
	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
	Notice  string `json:"notice"`
	AddMode int    `json:"add_mode"`
	Avatar  string `json:"avatar"`
}
