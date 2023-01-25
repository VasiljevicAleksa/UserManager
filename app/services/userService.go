package services

import (
	"strings"
	"usermanager/app/domain"
	notif "usermanager/app/infrastructure/notification"
	repo "usermanager/app/infrastructure/repositories"
	proto "usermanager/app/ui/protos/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Add(req *proto.CreateUserRequest) (uuid.UUID, error)
	Update(req *proto.UpdateUserRequest) error
	Delete(id string) error
	GetPage(req *proto.UserPageRequest) ([]domain.User, error)
}

type userService struct {
	repo  repo.UserRepo
	notif notif.NotificationService
}

func NewUserService(r repo.UserRepo, n notif.NotificationService) *userService {
	return &userService{
		repo:  r,
		notif: n,
	}
}

// Create new user. Returns user id or error if occured.
func (u userService) Add(req *proto.CreateUserRequest) (uuid.UUID, error) {
	// create user domain model and add to db
	user := userFromCreateReq(req)
	if err := u.repo.Add(user); err != nil {
		return uuid.Nil, err
	}

	return user.Id, nil
}

// Update provided user. Returns error if occured.
func (u *userService) Update(req *proto.UpdateUserRequest) error {
	// create user domain model and update
	user := userFromUpdateReq(req)
	if err := u.repo.Update(user); err != nil {
		return err
	}

	// notify services subscribed to user change notifications
	// this can go asynchronously, we don't need the results
	go u.notif.NotifyAboutUserChange(user.Id)

	return nil
}

// Delete user with provided id. Returns error if occured.
func (u *userService) Delete(id string) error {
	return u.repo.Delete(uuid.MustParse(id))
}

// Get user page method. Returns list of users or error if ocurred.
func (u *userService) GetPage(req *proto.UserPageRequest) ([]domain.User, error) {
	return u.repo.GetPage(req.Filter, req.Offset, req.Limit)
}

// Create user domain model from CreateUserRequest.
func userFromCreateReq(req *proto.CreateUserRequest) domain.User {
	return domain.User{
		Id:        uuid.New(),
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Nickname:  req.Nickname,
		Password:  hashPassword(req.Password),
		Email:     req.Email,
		Country:   strings.ToUpper(req.Country),
	}
}

// Create user update model from UpdateUserRequest
func userFromUpdateReq(req *proto.UpdateUserRequest) domain.User {
	return domain.User{
		Id:        uuid.MustParse(req.Id),
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Nickname:  req.Nickname,
		Password:  hashPassword(req.Password),
		Email:     req.Email,
		Country:   strings.ToUpper(req.Country),
	}
}

// Returns hash value of password
func hashPassword(password string) string {
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes)
}
