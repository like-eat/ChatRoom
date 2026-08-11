package request

type CreateGroupRequest struct {
	OwnerId string `json:"owner_id"`
	Name    string `json:"name"`
	Notice  string `json:"notice"`
	Avatar  string `json:"avatar"`
}
