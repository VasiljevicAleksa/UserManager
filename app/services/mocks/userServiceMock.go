package mocks

import (
	"usermanager/app/domain"
	proto "usermanager/app/ui/protos/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (u *UserServiceMock) Add(req *proto.CreateUserRequest) (uuid.UUID, error) {
	args := u.Called(req)

	var r0 uuid.UUID
	if rf, ok := args.Get(0).(func(*proto.CreateUserRequest) uuid.UUID); ok {
		r0 = rf(req)
	} else {
		r0 = args.Get(0).(uuid.UUID)
	}

	var r1 error
	if rf, ok := args.Get(1).(func(*proto.CreateUserRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

func (u *UserServiceMock) Update(req *proto.UpdateUserRequest) error {
	args := u.Called(req)

	var r0 error
	if rf, ok := args.Get(0).(func(*proto.UpdateUserRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = args.Error(0)
	}

	return r0
}

func (u *UserServiceMock) Delete(id string) error {
	args := u.Called(id)

	var r0 error
	if rf, ok := args.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = args.Error(0)
	}

	return r0
}

func (u *UserServiceMock) GetPage(req *proto.UserPageRequest) ([]domain.User, error) {
	args := u.Called(req)

	var r0 []domain.User
	if rf, ok := args.Get(0).(func(*proto.UserPageRequest) []domain.User); ok {
		r0 = rf(req)
	} else {
		r0 = args.Get(0).([]domain.User)
	}

	var r1 error
	if rf, ok := args.Get(1).(func(*proto.UserPageRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}
