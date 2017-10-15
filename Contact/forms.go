package Contact

import (
	"encoding/json"
	company "github.com/horae-app/api/Company"
	"net/http"
)

func ContactForm(r *http.Request, company_id string) (Contact, string) {
	var contact Contact

	if r.Body == nil {
		return contact, "Please send a request body"
	}
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
