package Calendar

import (
	"encoding/json"
	contact "github.com/horae-app/api/Contact"
	"net/http"
)

func CalendarForm(r *http.Request, company_id string, contact_id string) (Calendar, string) {
	var calendar Calendar

	err := json.NewDecoder(r.Body).Decode(&calendar)
	if err != nil {
		return calendar, err.Error()
	}

	cal_contact, errMsg := contact.GetById(company_id, contact_id)
	if errMsg != "" {
		return calendar, "No contact for " + contact_id
	}
	calendar.Company = cal_contact.Company
	calendar.Contact = cal_contact

	return calendar, calendar.Validate()
}