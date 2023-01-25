package services

import (
	"errors"
	"testing"
	"time"
	"usermanager/app/domain"
	notifMock "usermanager/app/infrastructure/notification/mocks"
	repoMock "usermanager/app/infrastructure/repositories/mocks"
	proto "usermanager/app/ui/protos/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// create user service with mocked objects
func createUserService() (UserService, *repoMock.UserRepoMock, *notifMock.NotificationServiceMock) {
	mockedUserRepo := &repoMock.UserRepoMock{}
	mockedNotifService := &notifMock.NotificationServiceMock{}

	userService := NewUserService(mockedUserRepo, mockedNotifService)
	return userService, mockedUserRepo, mockedNotifService
}

func TestAdd_RepoAddErr_ShouldReturnErr(t *testing.T) {
	userService, mockedUserRepo, _ := createUserService()

	// arrange
	expectedRes := uuid.Nil
	expectedErr := errors.New("test error")
	req := &proto.CreateUserRequest{}

	mockedUserRepo.
		On("Add", mock.AnythingOfType("User")).
		Return(expectedErr)

	// act
	res, err := userService.Add(req)

	// assert
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr.Error())
	assert.Equal(t, res, expectedRes)
}

func TestAdd_RepoAddPass_ShouldReturnResult(t *testing.T) {
	userService, mockedUserRepo, _ := createUserService()

	// arrange
	req := &proto.CreateUserRequest{
		Firstname: "test-firstname",
		Lastname:  "test-lastname",
		Nickname:  "test",
		Password:  "test-pass",
		Email:     "test@test.com",
		Country:   "SRB",
	}

	// used to match a mock call based on only certain properties from a complex struct
	userParamMatcher := mock.MatchedBy(func(user domain.User) bool {
		firstnameMatched := user.Firstname == req.Firstname
		lastnameMatched := user.Lastname == req.Lastname
		nickMatched := user.Nickname == req.Nickname
		passMatched := compareHashAndPass(req.Password)
		emailMatched := user.Email == req.Email
		countryMached := user.Country == req.Country

		return firstnameMatched && lastnameMatched && nickMatched &&
			passMatched && emailMatched && countryMached
	})

	mockedUserRepo.
		On("Add", userParamMatcher).
		Return(nil)

	// act
	res, err := userService.Add(req)

	// assert
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestUpdate_RepoUpdateErr_ShouldReturnErr(t *testing.T) {
	userService, mockedUserRepo, _ := createUserService()

	// arrange
	expectedErr := errors.New("test error")
	req := &proto.UpdateUserRequest{
		Id:        "eb24efdf-0043-4df7-b736-1486068abf03",
		Firstname: "test-firstname",
		Lastname:  "test-lastname",
		Nickname:  "test",
		Password:  "test-pass",
		Email:     "test@test.com",
		Country:   "SRB",
	}

	mockedUserRepo.
		On("Update", mock.AnythingOfType("User")).
		Return(expectedErr)

	// act
	err := userService.Update(req)

	// assert
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr.Error(), err.Error())
}

func TestUpdate_RepoUpdatePass_ShouldReturnResult(t *testing.T) {
	userService, mockedUserRepo, mockedNotifService := createUserService()

	// arrange
	req := &proto.UpdateUserRequest{
		Id:        uuid.NewString(),
		Firstname: "test-firstname",
		Lastname:  "test-lastname",
		Nickname:  "test",
		Password:  "test-pass",
		Email:     "test@test.com",
		Country:   "SRB",
	}

	// used to match a mock call based on only certain properties from a complex struct
	userParamMatcher := mock.MatchedBy(func(user domain.User) bool {
		firstnameMatched := user.Firstname == req.Firstname
		lastnameMatched := user.Lastname == req.Lastname
		nickMatched := user.Nickname == req.Nickname
		passMatched := compareHashAndPass(req.Password)
		emailMatched := user.Email == req.Email
		countryMached := user.Country == req.Country

		return firstnameMatched && lastnameMatched && nickMatched &&
			passMatched && emailMatched && countryMached
	})

	mockedUserRepo.
		On("Update", userParamMatcher).
		Return(nil)

	mockedNotifService.
		On("NotifyAboutUserChange", uuid.MustParse(req.Id)).
		Return(nil)

	// act
	err := userService.Update(req)

	// assert
	assert.Nil(t, err)
}

func TestUpdate_RepoUpdatePass_ShouldSendNotification(t *testing.T) {
	userService, mockedUserRepo, mockedNotifService := createUserService()

	// arrange
	req := &proto.UpdateUserRequest{Id: "eb24efdf-0043-4df7-b736-1486068abf03"}

	mockedNotifService.
		On("NotifyAboutUserChange", uuid.MustParse(req.Id)).
		Return(nil)

	mockedUserRepo.
		On("Update", mock.AnythingOfType("User")).
		Return(nil)

	// act
	err := userService.Update(req)

	// because NotifyAboutUserChange will run async in another goroutine,
	// we need make sure that NotifyAboutUserChange is called before we
	// finish with test otherwise assertCalled() could failed
	time.Sleep(time.Second * 2)

	// assert
	assert.Nil(t, err)
	mockedNotifService.AssertCalled(t, "NotifyAboutUserChange", uuid.MustParse(req.Id))
}

func compareHashAndPass(password string) bool {
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err := bcrypt.CompareHashAndPassword(hashedBytes, []byte(password))
	return err == nil
}
