package Calendar

import (
	"github.com/gocql/gocql"
	"github.com/horae-app/api/Cassandra"
	company "github.com/horae-app/api/Company"
	contact "github.com/horae-app/api/Contact"
	"time"
)

type Calendar struct {
	ID          gocql.UUID
	Company     company.CompanyBasic
	Contact     contact.Contact
	Start_at    time.Time
	End_at      time.Time
	Description string
	Value       float32
}

func (self Calendar) IsNew() bool {
	return self.ID.String() == "00000000-0000-0000-0000-000000000000"
}

func (self Calendar) Save() (Calendar, string) {
	errMsg := self.Validate()
	if errMsg != "" {
		return self, errMsg
	}

	if self.IsNew() {
		self.ID = gocql.TimeUUID()
	}

	db_cmd := "UPDATE calendar SET company_id = ?, contact_id = ?, start_at = ?, end_at = ?, description = ?, value = ? WHERE id = ?"
	query := Cassandra.Session.Query(db_cmd, self.Company.ID, self.Contact.ID, self.Start_at, self.End_at, self.Description, self.Value, self.ID)
	err := query.Exec()
	if err != nil {
		return self, err.Error()
	}

	return self, ""
}

func (self Calendar) Validate() string {
	if self.Company.IsNew() {
		return "company is required"
	}

	if self.Contact.IsNew() {
		return "contact is required"
	}

	if self.Start_at.IsZero() {
		return "start at required"
	}

	if self.End_at.IsZero() {
		return "end at required"
	}

	return ""
}

func (self Calendar) Delete() (bool, string) {
	if self.IsNew() {
		return false, "Calendar not found"
	}

	db_cmd := "DELETE FROM calendar WHERE id = ?"
	query := Cassandra.Session.Query(db_cmd, self.ID)
	err := query.Exec()
	if err != nil {
		return false, err.Error()
	}

	return true, ""
}

func GetById(companyId string, id string) (Calendar, string) {
	return GetBy(companyId, "id", id)
}

func GetBy(companyId string, field string, value string) (Calendar, string) {
	var calendar Calendar

	db_cmd := "SELECT id, company_id, contact_id, start_at, end_at, description, value from calendar WHERE " + field + " = ?"
	query := Cassandra.Session.Query(db_cmd, value)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		if m["company_id"].(gocql.UUID).String() != companyId {
			continue
		}

		cal_contact, _ := contact.GetById(m["company_id"].(gocql.UUID).String(), m["contact_id"].(gocql.UUID).String())

		calendar = Calendar{
			ID:          m["id"].(gocql.UUID),
			Company:     cal_contact.Company,
			Contact:     cal_contact,
			Start_at:    m["start_at"].(time.Time),
			End_at:      m["end_at"].(time.Time),
			Description: m["description"].(string),
			Value:       m["value"].(float32),
		}
	}

	if calendar.IsNew() {
		return calendar, "Not found"
	}

	return calendar, ""
}

func GetAll(companyId string) []Calendar {
	var calendars []Calendar

	db_cmd := "SELECT id, company_id, contact_id, start_at, end_at, description, value from calendar WHERE company_id = ?"
	query := Cassandra.Session.Query(db_cmd, companyId)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		cal_contact, _ := contact.GetById(m["company_id"].(gocql.UUID).String(), m["contact_id"].(gocql.UUID).String())
		calendar := Calendar{
			ID:          m["id"].(gocql.UUID),
			Company:     cal_contact.Company,
			Contact:     cal_contact,
			Start_at:    m["start_at"].(time.Time),
			End_at:      m["end_at"].(time.Time),
			Description: m["description"].(string),
			Value:       m["value"].(float32),
		}
		calendars = append(calendars, calendar)
	}

	return calendars
}

type NewCalendarResponse struct {
	ID gocql.UUID
}

type SuccessResponse struct {
	Message string
}

type ErrorResponse struct {
	Error string
}

type ListResponse struct {
	Calendars []Calendar
}
