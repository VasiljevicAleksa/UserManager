package mocks

import (
	"usermanager/app/domain"
	proto "usermanager/app/ui/protos/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (r *UserRepoMock) Add(user domain.User) error {
	args := r.Called(user)

	var r0 error
	if rf, ok := args.Get(0).(func(domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = args.Error(0)
	}

	return r0
}

func (r *UserRepoMock) Update(user domain.User) error {
	args := r.Called(user)

	var r0 error
	if rf, ok := args.Get(0).(func(domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = args.Error(0)
	}

	return r0
}

// not implemented
// we don't need get page mock at the moment
func (r *UserRepoMock) GetPage(filter *proto.UserPageRequest_UserFilterOptions,
	offset int32, limit int32) (users []domain.User, err error) {
	return users, nil
}

// not implemented
// we don't need delete mock at the moment
func (r *UserRepoMock) Delete(id uuid.UUID) error {
	return nil
}
