package request

type RemoveGroupMembersRequest struct {
	GroupId  string   `json:"group_id"`
	OwnerId  string   `json:"owner_id"`
	UuidList []string `json:"uuid_list"`
}
