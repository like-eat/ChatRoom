package request

type CreateGroupRequest struct {
	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
	Notice  string `json:"notice"`
	AddMode bool   `json:"add_mode"`
	Avatar  string `json:"avatar"`
}
