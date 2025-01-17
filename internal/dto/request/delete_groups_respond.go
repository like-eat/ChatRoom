package request

type DeleteGroupsRequest struct {
	UuidList []string `json:"uuid_list"`
}
