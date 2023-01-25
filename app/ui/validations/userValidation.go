package validation

import (
	"errors"
	"net/mail"
	proto "usermanager/app/ui/protos/user"

	"github.com/google/uuid"
)

// CreateUserRequest proto message validation
func ValidateCreateUserReq(p *proto.CreateUserRequest) error {
	if p.Firstname == "" {
		return errors.New("firstname is required")
	}
	if p.Lastname == "" {
		return errors.New("lastname is required")
	}
	if p.Nickname == "" {
		return errors.New("nickname is required")
	}
	if p.Password == "" {
		return errors.New("password is required")
	}
	if err := emailValidation(p.Email); err != nil {
		return err
	}
	if err := countryValidation(p.Country); err != nil {
		return err
	}

	return nil
}

// UpdateUserRequest proto message validation
func ValidateUpdateUserReq(p *proto.UpdateUserRequest) error {
	if err := validateId(p.Id); err != nil {
		return err
	}
	if p.Firstname == "" {
		return errors.New("firstname is required")
	}
	if p.Lastname == "" {
		return errors.New("lastname is required")
	}
	if p.Nickname == "" {
		return errors.New("nickname is required")
	}
	if p.Password == "" {
		return errors.New("password is required")
	}
	if err := emailValidation(p.Email); err != nil {
		return err
	}
	if err := countryValidation(p.Country); err != nil {
		return err
	}
	return nil
}

// UserPageRequest proto message validation
func ValidateUserPageReq(p *proto.UserPageRequest) error {
	if p.Filter != nil {
		if p.Filter.Country != "" {
			if len(p.Filter.Country) > 2 {
				return errors.New("country should have 2 letters")
			}
		}

		if p.Filter.CreatedFrom != nil && p.Filter.CreatedTo != nil {
			if p.Filter.CreatedTo.AsTime().Before(p.Filter.CreatedFrom.AsTime()) {
				return errors.New("'Created to' time is before 'created from'")
			}
		}
	}

	return nil
}

// DeleteUserRequest proto message validation
func ValidateDeleteUserReq(p *proto.DeleteUserRequest) error {
	return validateId(p.Id)
}

func emailValidation(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email bad format")
	}
	return nil
}

// country input validation, should consists of a 2 letters.
func countryValidation(country string) error {
	if country == "" {
		return errors.New("country is required")
	}
	if len(country) != 2 {
		return errors.New("country should have 2 letters")
	}
	return nil
}

// validation of string in uuid format
func validateId(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("id wrong format")
	}
	return nil
}
