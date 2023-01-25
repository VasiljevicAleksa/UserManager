package validation

import (
	"testing"
	proto "usermanager/app/ui/protos/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserReq_WithValidReq_ShouldPass(t *testing.T) {
	req := &proto.CreateUserRequest{
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
		Country:   "RS",
	}
	err := ValidateCreateUserReq(req)

	assert.Nil(t, err)
}

func TestCreateUserReq_FirstNameMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.CreateUserRequest{
		Lastname: "test",
		Nickname: "test",
		Password: "testPass",
		Email:    "test@test.com",
		Country:  "SRB",
	}
	expectedErr := "firstname is required"

	err := ValidateCreateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestCreateUserReq_LastNameMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.CreateUserRequest{
		Firstname: "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
		Country:   "SRB",
	}
	expectedErr := "lastname is required"

	err := ValidateCreateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestCreateUserReq_NicknameMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.CreateUserRequest{
		Firstname: "test",
		Lastname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
		Country:   "SRB",
	}
	expectedErr := "nickname is required"

	err := ValidateCreateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestCreateUserReq_PasswordMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.CreateUserRequest{
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Email:     "test@test.com",
		Country:   "SRB",
	}
	expectedErr := "password is required"

	err := ValidateCreateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestCreateUserReq_EmailMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.CreateUserRequest{
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Country:   "SRB",
	}
	expectedErr := "email is required"

	err := ValidateCreateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestCreateUserReq_EmailWrongFormat_ShouldReturnErr(t *testing.T) {
	req := &proto.CreateUserRequest{
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "wrongFormat",
		Country:   "SRB",
	}
	expectedErr := "email bad format"

	err := ValidateCreateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestCreateUserReq_CountryMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.CreateUserRequest{
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
	}
	expectedErr := "country is required"

	err := ValidateCreateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestCreateUserReq_CountryWrongFormat_ShouldReturnErr(t *testing.T) {
	req := &proto.CreateUserRequest{
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
		Country:   "SRBSRBSRB",
	}
	expectedErr := "country should have 2 letters"

	err := ValidateCreateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_WithValidReq_ShouldPass(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:        uuid.NewString(),
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
		Country:   "RS",
	}
	err := ValidateUpdateUserReq(req)

	assert.Nil(t, err)
}

func TestUpdateUserReq_IdMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Lastname: "test",
		Nickname: "test",
		Password: "testPass",
		Email:    "test@test.com",
		Country:  "SRB",
	}
	expectedErr := "id is required"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_IdWrongFormat_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:       "wrong-format",
		Lastname: "test",
		Nickname: "test",
		Password: "testPass",
		Email:    "test@test.com",
		Country:  "SRB",
	}
	expectedErr := "id wrong format"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_FirstNameMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:       uuid.NewString(),
		Lastname: "test",
		Nickname: "test",
		Password: "testPass",
		Email:    "test@test.com",
		Country:  "SRB",
	}
	expectedErr := "firstname is required"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_LastNameMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:        uuid.NewString(),
		Firstname: "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
		Country:   "SRB",
	}
	expectedErr := "lastname is required"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_NicknameMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:        uuid.NewString(),
		Firstname: "test",
		Lastname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
		Country:   "SRB",
	}
	expectedErr := "nickname is required"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_PasswordMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:        uuid.NewString(),
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Email:     "test@test.com",
		Country:   "SRB",
	}
	expectedErr := "password is required"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_EmailMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:        uuid.NewString(),
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Country:   "SRB",
	}
	expectedErr := "email is required"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_EmailWrongFormat_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:        uuid.NewString(),
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "wrongFormat",
		Country:   "SRB",
	}
	expectedErr := "email bad format"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_CountryMissing_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:        uuid.NewString(),
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
	}
	expectedErr := "country is required"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}

func TestUpdateUserReq_CountryWrongFormat_ShouldReturnErr(t *testing.T) {
	req := &proto.UpdateUserRequest{
		Id:        uuid.NewString(),
		Firstname: "test",
		Lastname:  "test",
		Nickname:  "test",
		Password:  "testPass",
		Email:     "test@test.com",
		Country:   "SRBSRBSRB",
	}
	expectedErr := "country should have 2 letters"

	err := ValidateUpdateUserReq(req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), expectedErr)
}
