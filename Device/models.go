package Device

import (
	"errors"

	"github.com/horae-app/api/Cassandra"
)

type Device struct {
	Email string
	Token string
	Device string
}

func (self Device) Save() error {
	err := self.Validate()
	if err != nil {
		return err
	}

	db_cmd := "UPDATE device SET \"token\" = ?, device = ? WHERE email = ?"
	query := Cassandra.Session.Query(db_cmd, self.Token, self.Device, self.Email)
	return query.Exec()
}

func (self Device) Validate() error {
	if self.Email == "" {
		return errors.New("email is required")
	}

	if self.Token == "" {
		return errors.New("token is required")
	}

	if self.Device == "" {
		return errors.New("device is required")
	}

	return nil
}

func GetTokenByEmail(email string) string {
	var token string

	db_cmd := "SELECT \"token\" from device WHERE email = ?"
	query := Cassandra.Session.Query(db_cmd, email)
	iterable := query.Iter()
	for iterable.Scan(&token) {
		return token
	}

	return ""
}

type SuccessResponse struct {
	Message string
}

type ErrorResponse struct {
	Error string
}