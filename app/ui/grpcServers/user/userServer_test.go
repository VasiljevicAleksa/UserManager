package server

import (
	"context"
	"errors"
	"testing"
	"time"
	"usermanager/app/domain"
	"usermanager/app/services/mocks"
	proto "usermanager/app/ui/protos/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	_ "github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/grpc"
)

var createUserReq = &proto.CreateUserRequest{
	Firstname: "test",
	Lastname:  "test",
	Nickname:  "test",
	Password:  "testPass",
	Email:     "test@test.com",
	Country:   "RS",
}

var updateUserReq = &proto.UpdateUserRequest{
	Id:        uuid.NewString(),
	Firstname: "test",
	Lastname:  "test",
	Nickname:  "test",
	Password:  "testPass",
	Email:     "test@test.com",
	Country:   "RS",
}

var deleteUserReq = &proto.DeleteUserRequest{
	Id: uuid.NewString(),
}

var getPageReq = &proto.UserPageRequest{
	Offset: 0,
	Limit:  1,
}

func createServer() (*userServer, *mocks.UserServiceMock) {
	mockUserService := &mocks.UserServiceMock{}
	grpcServer := NewUserGrpcServer(grpc.NewServer(), mockUserService)
	return grpcServer, mockUserService
}

func TestCreateUser_UserServiceReturnErr_ResponseShouldBeErr(t *testing.T) {
	grpcServer, mockedUserService := createServer()

	expectedErr := errors.New("error ocurred")
	ctx := context.Background()

	mockedUserService.
		On("Add", createUserReq).
		Return(uuid.Nil, expectedErr).
		Once()

	result, err := grpcServer.CreateUser(ctx, createUserReq)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr.Error())
}

func TestCreateUser_UserServiceReturnsValidRes_ResponseShouldValid(t *testing.T) {
	grpcServer, mockedUserService := createServer()
	ctx := context.Background()

	expectedId := uuid.New()
	mockedUserService.
		On("Add", createUserReq).
		Return(expectedId, nil).
		Once()

	result, err := grpcServer.CreateUser(ctx, createUserReq)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Id, expectedId.String())
}

func TestUpdateUser_UserServiceReturnErr_ResponseShouldBeErr(t *testing.T) {
	grpcServer, mockedUserService := createServer()

	expectedErr := errors.New("error ocurred")
	ctx := context.Background()

	mockedUserService.
		On("Update", updateUserReq).
		Return(expectedErr).
		Once()

	result, err := grpcServer.UpdateUser(ctx, updateUserReq)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr.Error())
}

func TestUpdate_UserServiceReturnsValidRes_ResponseShouldValid(t *testing.T) {
	grpcServer, mockedUserService := createServer()
	ctx := context.Background()

	mockedUserService.
		On("Update", updateUserReq).
		Return(nil).
		Once()

	result, err := grpcServer.UpdateUser(ctx, updateUserReq)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Id, updateUserReq.Id)
}

func TestDeleteUser_UserServiceReturnErr_ResponseShouldBeErr(t *testing.T) {
	grpcServer, mockedUserService := createServer()

	expectedErr := errors.New("error ocurred")
	ctx := context.Background()

	mockedUserService.
		On("Delete", deleteUserReq.Id).
		Return(expectedErr).
		Once()

	result, err := grpcServer.DeleteUser(ctx, deleteUserReq)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr.Error())
}

func TestDelete_UserServiceReturnsValidRes_ResponseShouldValid(t *testing.T) {
	grpcServer, mockedUserService := createServer()
	ctx := context.Background()

	mockedUserService.
		On("Delete", deleteUserReq.Id).
		Return(nil).
		Once()

	result, err := grpcServer.DeleteUser(ctx, deleteUserReq)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Id, deleteUserReq.Id)
}

func TestGetUserPageUser_UserServiceReturnErr_ResponseShouldBeErr(t *testing.T) {
	grpcServer, mockedUserService := createServer()

	expectedErr := errors.New("error ocurred")
	ctx := context.Background()

	mockedUserService.
		On("GetPage", getPageReq).
		Return([]domain.User{}, expectedErr).
		Once()

	result, err := grpcServer.GetUserPage(ctx, getPageReq)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr.Error())
}

func TestGetUserPageUser_UserServiceReturnsValidRes_ResponseShouldValid(t *testing.T) {
	grpcServer, mockedUserService := createServer()
	ctx := context.Background()
	expectedUserList := []domain.User{
		{
			Id:        uuid.New(),
			Firstname: "test",
			Lastname:  "test",
			Nickname:  "test",
			Email:     "test@test.com",
			Country:   "SRB",
			CreatedAt: time.Now(),
		},
		{
			Id:        uuid.New(),
			Firstname: "test2",
			Lastname:  "test2",
			Nickname:  "test2",
			Email:     "test2@test.com",
			Country:   "SRB",
			CreatedAt: time.Now(),
		},
	}

	mockedUserService.
		On("GetPage", getPageReq).
		Return(expectedUserList, nil).
		Once()

	result, err := grpcServer.GetUserPage(ctx, getPageReq)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Users), len(expectedUserList))

	for i := 0; i < len(expectedUserList); i++ {
		assert.Equal(t, result.Users[0].Id, expectedUserList[0].Id.String())
		assert.Equal(t, result.Users[0].Firstname, expectedUserList[0].Firstname)
		assert.Equal(t, result.Users[0].Lastname, expectedUserList[0].Lastname)
		assert.Equal(t, result.Users[0].Nickname, expectedUserList[0].Nickname)
		assert.Equal(t, result.Users[0].Email, expectedUserList[0].Email)
		assert.Equal(t, result.Users[0].Country, expectedUserList[0].Country)
		assert.Equal(t, result.Users[0].Created.AsTime(), expectedUserList[0].CreatedAt.UTC())
	}
}
