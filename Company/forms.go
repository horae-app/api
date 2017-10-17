package Company

import (
	"encoding/json"
	"net/http"
)

func UserForm(r *http.Request) (CompanyFull, string) {
	var company CompanyFull

	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		return company, err.Error()
	}

	return company, company.Validate()
}

func AuthForm(r *http.Request) (bool, string) {
	var auth AuthRequest
	var company CompanyBasic

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		return false, err.Error()
	}

	company, errMsg := GetByEmail(auth.Email)
	if errMsg != "" {
		return false, "Incorrect username and/or password"
	}

	password, errMsg := GetPassword(auth.Email)
	if errMsg != "" || password != auth.Password {
		return false, "Incorrect username and/or password"
	}

	return true, company.ID.String()
}
