package repo

import (
	"errors"
	"log"
	"regexp"
	"testing"
	"usermanager/app/domain"
	proto "usermanager/app/ui/protos/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createUserRepo() (*userRepo, sqlmock.Sqlmock) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to create a stub db connection: %v", err)
	}
	//defer mockDb.Close()

	dialector := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       mockDb,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open gorm db: %v", err)
	}

	return NewUserRepo(gdb), mock
}

func TestAdd_NickAlreadyExist_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arange
	expectedErr := "nickname already exist"
	expectedErrCode := codes.InvalidArgument
	pgErr := pgconn.PgError{
		Code:           UNIQUE_INDEX_VIOLATION_CODE,
		ConstraintName: domain.UniqueConstraintNickname,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnError(error(&pgErr))
	mock.ExpectRollback()

	// act
	res := userRepo.Add(domain.User{})

	// assert
	assert.NotNil(t, res)
	statusErr := status.Convert(res)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestAdd_EmailAlreadyExist_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	expectedErr := "email already exist"
	expectedErrCode := codes.InvalidArgument
	pgErr := pgconn.PgError{
		Code:           UNIQUE_INDEX_VIOLATION_CODE,
		ConstraintName: domain.UniqueConstraintEmail,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnError(error(&pgErr))
	mock.ExpectRollback()

	// act
	res := userRepo.Add(domain.User{})

	// assert
	assert.NotNil(t, res)
	statusErr := status.Convert(res)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestAdd_ErrOcurred_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	expectedErr := "test err"
	expectedErrCode := codes.Internal

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnError(errors.New(expectedErr))
	mock.ExpectRollback()

	// act
	res := userRepo.Add(domain.User{})

	// assert
	assert.NotNil(t, res)
	statusErr := status.Convert(res)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestAdd_ShouldPass(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// act
	res := userRepo.Add(domain.User{})

	// assert
	assert.Nil(t, res)
}

func TestUpdate_UserDoesntExist_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	expectedErr := "no user in database"
	expectedErrCode := codes.NotFound

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	// act
	res := userRepo.Update(domain.User{})

	// assert
	assert.NotNil(t, res)
	statusErr := status.Convert(res)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestUpdate_NickAlreadyExist_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	expectedErr := "nickname already exist"
	expectedErrCode := codes.InvalidArgument
	pgErr := pgconn.PgError{
		Code:           UNIQUE_INDEX_VIOLATION_CODE,
		ConstraintName: domain.UniqueConstraintNickname,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
		WillReturnError(error(&pgErr))
	mock.ExpectRollback()

	// act
	res := userRepo.Update(domain.User{})

	// assert
	assert.NotNil(t, res)
	statusErr := status.Convert(res)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestUpdate_EmailAlreadyExist_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	expectedErr := "email already exist"
	expectedErrCode := codes.InvalidArgument
	pgErr := pgconn.PgError{
		Code:           UNIQUE_INDEX_VIOLATION_CODE,
		ConstraintName: domain.UniqueConstraintEmail,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
		WillReturnError(error(&pgErr))
	mock.ExpectRollback()

	// act
	res := userRepo.Update(domain.User{})

	// assert
	assert.NotNil(t, res)
	statusErr := status.Convert(res)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestUpdate_ErrOcurred_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	expectedErr := "test err"
	expectedErrCode := codes.Internal

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
		WillReturnError(errors.New(expectedErr))
	mock.ExpectRollback()

	// act
	res := userRepo.Update(domain.User{})

	// assert
	assert.NotNil(t, res)
	statusErr := status.Convert(res)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestUpdate_ShouldPass(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// act
	res := userRepo.Update(domain.User{})

	// assert
	assert.Nil(t, res)
}

func TestDelete_ErrOcurred_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	userId := uuid.New()
	expectedErr := "cannot delete user TEST ERR"
	expectedErrCode := codes.Internal

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users"`)).
		WillReturnError(errors.New("TEST ERR"))

	mock.ExpectRollback()

	// act
	err := userRepo.Delete(userId)

	// assert
	assert.NotNil(t, err)
	statusErr := status.Convert(err)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestDeleteUser_UserNotFound_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	mockUserId := uuid.New()
	expectedErr := "no user in database"
	expectedErrCode := codes.NotFound

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users"`)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	// act
	res := userRepo.Delete(mockUserId)

	// assert
	assert.NotNil(t, res)
	statusErr := status.Convert(res)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestDeleteUser_ShouldPass(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	mockUserId := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// act
	res := userRepo.Delete(mockUserId)

	// assert
	assert.Nil(t, res)
}

func TestGetUserList_ErrOcurred_ShouldReturnErr(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	expectedErr := "test err"
	expectedErrCode := codes.Internal

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnError(errors.New(expectedErr))

	// act
	res, err := userRepo.GetPage(&proto.UserPageRequest_UserFilterOptions{}, 1, 1)

	// assert
	assert.Equal(t, 0, len(res))
	assert.NotNil(t, err)
	statusErr := status.Convert(err)
	assert.Equal(t, expectedErrCode, statusErr.Code())
	assert.Equal(t, expectedErr, statusErr.Message())
}

func TestGetUserList_ShouldPass(t *testing.T) {
	userRepo, mock := createUserRepo()

	// arrange
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country"}).
		AddRow("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c", "aleksa", "vasiljevic", "aki", "pass", "a@gmail.com", "SRB").
		AddRow("cea1b24d-0627-4ea0-aa2b-7af4c6c2a31c", "marko", "milanovic", "mare", "pass2", "m@gmail.com", "SRB")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(rows)

	// act
	res, err := userRepo.GetPage(&proto.UserPageRequest_UserFilterOptions{}, 1, 2)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res))
}
