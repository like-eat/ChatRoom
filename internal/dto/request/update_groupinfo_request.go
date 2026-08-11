package request

type UpdateGroupInfoRequest struct {
	OwnerId string `json:"owner_id"`
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Notice  string `json:"notice"`
}
