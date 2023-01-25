package repo

import (
	"usermanager/app/domain"
	proto "usermanager/app/ui/protos/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const UNIQUE_INDEX_VIOLATION_CODE = "23505"

type UserRepo interface {
	Add(user domain.User) error
	Update(user domain.User) error
	GetPage(filter *proto.UserPageRequest_UserFilterOptions, offset int32, limit int32) (users []domain.User, err error)
	Delete(id uuid.UUID) error
}

type userRepo struct {
	db *gorm.DB
}

// Create new user repository with Gorm ORM library.
func NewUserRepo(gormDb *gorm.DB) *userRepo {
	return &userRepo{db: gormDb}
}

// Add user method. Returns an error if ocurred.
func (r *userRepo) Add(user domain.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return handleErr(err)
	}
	return nil
}

// Update user method. Returns an error if ocurred.
func (r *userRepo) Update(user domain.User) error {
	res := r.db.
		Model(&user).
		Where("id = ?", user.Id).
		Updates(domain.User{
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Nickname:  user.Nickname,
			Password:  user.Password,
			Email:     user.Email,
			Country:   user.Country})

	if res.Error != nil {
		return handleErr(res.Error)
	}
	if res.RowsAffected == 0 {
		return status.Error(codes.NotFound, "no user in database")
	}
	return nil
}

// Delete user method. Returns error if ocurred.
func (r *userRepo) Delete(id uuid.UUID) error {
	res := r.db.Delete(domain.User{}, id)

	if res.Error != nil {
		return status.Errorf(codes.Internal, "cannot delete user %v", res.Error)
	}
	if res.RowsAffected == 0 {
		return status.Error(codes.NotFound, "no user in database")
	}
	return nil
}

// Get user page method based on the provided filter params.
// Returns list of users or error if ocurred.
func (r *userRepo) GetPage(filter *proto.UserPageRequest_UserFilterOptions,
	offset int32, limit int32) (users []domain.User, err error) {

	// set limit and offset
	selectQuery := r.db.
		Limit(int(limit)).
		Offset(int(offset))

	// based on filter params create where condition
	if filter != nil {
		if filter.Country != "" {
			selectQuery = selectQuery.
				Where("country = ?", filter.Country)
		}
		if filter.CreatedFrom != nil {
			selectQuery = selectQuery.
				Where("created_at > ?", filter.CreatedFrom.AsTime())
		}
		if filter.CreatedTo != nil {
			selectQuery = selectQuery.
				Where("created_at < ?", filter.CreatedTo.AsTime())
		}
	}

	// check for the users
	if err := selectQuery.Find(&users).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return users, nil
}

// If there was an error it is important to handle it
// and check the uniqueness of the name and email.
func handleErr(err error) error {
	if isUniqueConstraintError(err, domain.UniqueConstraintNickname) {
		return status.Error(codes.InvalidArgument, "nickname already exist")
	}
	if isUniqueConstraintError(err, domain.UniqueConstraintEmail) {
		return status.Error(codes.InvalidArgument, "email already exist")
	}
	return status.Error(codes.Internal, err.Error())
}

// Based on a constraint and pg error code check
// whether the condition of uniqueness is violated.
func isUniqueConstraintError(err error, constraintName string) bool {
	if pqErr, ok := err.(*pgconn.PgError); ok {
		return pqErr.Code == UNIQUE_INDEX_VIOLATION_CODE &&
			pqErr.ConstraintName == constraintName
	}
	return false
}
