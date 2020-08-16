package controllers_test

import "github.com/nikolausreza131192/pos/entity"

type fakeItemRepo struct {
	GetAllResults []entity.Item
	GetByIDResult entity.Item
}

func (r *fakeItemRepo) GetAll() []entity.Item {
	return r.GetAllResults
}

func (r *fakeItemRepo) GetByID(id int) entity.Item {
	return r.GetByIDResult
}

type fakeUserRepo struct {
	GetByUsernameResult   entity.User
	GetUserPasswordResult string
	GetUserPasswordError  error
	CreateUserResult      string
	CreateUserError       error
}

func (r *fakeUserRepo) GetByUsername(username string) entity.User {
	return r.GetByUsernameResult
}

func (r *fakeUserRepo) GetUserPassword(username string) (string, error) {
	return r.GetUserPasswordResult, r.GetUserPasswordError
}

func (r *fakeUserRepo) CreateUser(user entity.User) (string, error) {
	return r.CreateUserResult, r.CreateUserError
}
