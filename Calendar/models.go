package Calendar

import (
	"time"
	"github.com/gocql/gocql"
	"github.com/horae-app/api/Cassandra"
	company "github.com/horae-app/api/Company"
	contact "github.com/horae-app/api/Contact"
)


type CalendarBasic struct {
	ID          gocql.UUID
	Start_at    time.Time
	End_at      time.Time
	Description string
	Value       float32
	Status      string
}

type CalendarList struct {
	CalendarBasic
	Contact_id gocql.UUID
}

type Calendar struct {
	CalendarBasic
	Company company.CompanyBasic
	Contact contact.Contact
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

	if self.Status == "" {
		self.Status = "new";
	}

	db_cmd := "UPDATE calendar SET company_id = ?, contact_id = ?, start_at = ?, end_at = ?, description = ?, value = ?, status = ? WHERE id = ?"
	query := Cassandra.Session.Query(db_cmd, self.Company.ID, self.Contact.ID, self.Start_at, self.End_at, self.Description, self.Value, self.Status, self.ID)
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

	db_cmd := "SELECT id, company_id, contact_id, start_at, end_at, description, value, status from calendar WHERE " + field + " = ?"
	query := Cassandra.Session.Query(db_cmd, value)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		if m["company_id"].(gocql.UUID).String() != companyId {
			continue
		}

		cal_contact, _ := contact.GetById(m["company_id"].(gocql.UUID).String(), m["contact_id"].(gocql.UUID).String())

		calendar = Calendar{
			CalendarBasic: CalendarBasic{
				ID:          m["id"].(gocql.UUID),
				Start_at:    m["start_at"].(time.Time),
				End_at:      m["end_at"].(time.Time),
				Description: m["description"].(string),
				Value:       m["value"].(float32),
				Status:      m["status"].(string),
			},
			Company: cal_contact.Company,
			Contact: cal_contact,
		}
	}

	if calendar.IsNew() {
		return calendar, "Not found"
	}

	return calendar, ""
}

func GetAll(companyId string) []CalendarList {
	var calendars []CalendarList
	var id, contact_id gocql.UUID
	var start_at, end_at time.Time
	var description, status string
	var value float32

	db_cmd := "SELECT id, contact_id, start_at, end_at, description, value, status from calendar WHERE company_id = ?"
	query := Cassandra.Session.Query(db_cmd, companyId)
	iterable := query.Iter()
	for iterable.Scan(&id, &contact_id, &start_at, &end_at, &description, &value, &status) {
		calendar := CalendarList{
			CalendarBasic: CalendarBasic{
				ID:          id,
				Start_at:    start_at,
				End_at:      end_at,
				Description: description,
				Value:       value,
				Status:      status,
			},
			Contact_id: contact_id,
		}
		calendars = append(calendars, calendar)
	}

	return calendars
}

func GetAllTomorrow(date time.Time) []Calendar {
	var calendars []Calendar
	var id, contact_id gocql.UUID
	var start_at, end_at time.Time
	var description, status string
	var value float32

	start_at = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	start_at = start_at.Add(time.Duration(24)*time.Hour)
	end_at = start_at.Add(time.Duration(86399)*time.Second)

	db_cmd := "SELECT id, contact_id, start_at, end_at, description, value, status from calendar WHERE start_at > ? AND start_at <= ? ALLOW FILTERING"
	query := Cassandra.Session.Query(db_cmd, start_at, end_at)
	iterable := query.Iter()
	for iterable.Scan(&id, &contact_id, &start_at, &end_at, &description, &value, &status) {
		calendar := Calendar{
			CalendarBasic: CalendarBasic{
				ID:          id,
				Start_at:    start_at,
				End_at:      end_at,
				Description: description,
				Value:       value,
				Status:      status,
			},
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
	Calendars []CalendarList
}
