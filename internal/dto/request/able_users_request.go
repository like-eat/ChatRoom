package request

type AbleUsersRequest struct {
	UuidList []string `json:"uuid_list"`
	IsAdmin  int8     `json:"is_admin"`
}
