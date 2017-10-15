package Company

import (
	"encoding/json"
	"net/http"
)

func UserForm(r *http.Request) (Company, string) {
	var company Company

	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		return company, err.Error()
	}

	return company, company.Validate()
}

func AuthForm(r *http.Request) (bool, string) {
	var auth AuthRequest
	var company Company

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		return false, err.Error()
	}

	company, errMsg := GetByEmail(auth.Email)
	if errMsg != "" || company.Password != auth.Password {
		return false, "Incorrect username and/or password"
	}

	return true, company.ID.String()
}
