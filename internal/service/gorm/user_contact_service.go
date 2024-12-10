package gorm

type userContactService struct {
}

var UserContactService = new(userContactService)

func (u *userContactService) GetUserList() (string, error) {

	return "", nil
}
