package Contact

import (
	"encoding/json"
	company "github.com/horae-app/api/Company"
	"net/http"
)

func ContactForm(r *http.Request, company_id string) (Contact, string) {
	var contact Contact

	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		return contact, err.Error()
	}

	cont_company, errMsg := company.GetById(company_id)
	if errMsg != "" {
		return contact, "No company for " + company_id
	}
	contact.Company = cont_company

	return contact, contact.Validate()
}

func CalendarForm(r *http.Request) (string, string) {
	var auth AuthRequest
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		return "", err.Error()
	}

	if Auth(auth.Email, auth.Token) {
		return auth.Email, "Incorrect username and/or password"
	}

	return auth.Email, ""
}
