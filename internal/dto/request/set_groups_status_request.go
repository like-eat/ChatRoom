package request

type SetGroupsStatusRequest struct {
	UuidList []string `json:"uuid_list"`
	Status   int8     `json:"status"`
}
