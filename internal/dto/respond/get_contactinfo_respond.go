package respond

import "encoding/json"

type GetContactInfoRespond struct {
	ContactId        string          `json:"contact_id"`
	ContactName      string          `json:"contact_name"`
	ContactAvatar    string          `json:"contact_avatar"`
	ContactPhone     string          `json:"contact_phone,omitempty"`
	ContactEmail     string          `json:"contact_email,omitempty"`
	ContactGender    bool            `json:"contact_gender,omitempty"`
	ContactSignature string          `json:"contact_signature,omitempty"`
	ContactBirthday  string          `json:"contact_birthday,omitempty"`
	ContactNotice    string          `json:"contact_notice,omitempty"`
	ContactMembers   json.RawMessage `json:"contact_members,omitempty"`
	ContactMemberCnt int             `json:"contact_member_cnt,omitempty"`
	ContactOwnerId   string          `json:"contact_owner_id,omitempty"`
	ContactAddMode   bool            `json:"contact_add_mode,omitempty"`
}
