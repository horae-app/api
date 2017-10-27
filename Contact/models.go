package Contact

import (
	"github.com/gocql/gocql"
	"github.com/horae-app/api/Cassandra"
	company "github.com/horae-app/api/Company"
	util "github.com/horae-app/api/Util"
	"log"
	"math/rand"
	"time"
)

type ContactBasic struct {
	ID    gocql.UUID
	Name  string
	Email string
	Phone string
}

type Contact struct {
	ContactBasic
	Company company.CompanyBasic
}

func (self Contact) IsNew() bool {
	return self.ID.String() == "00000000-0000-0000-0000-000000000000"
}

func (self Contact) Save() (Contact, string) {
	errMsg := self.Validate()
	if errMsg != "" {
		return self, errMsg
	}

	if self.IsNew() {
		self.ID = gocql.TimeUUID()
	}

	db_cmd := "UPDATE contact SET company_id = ?, name = ?, email = ?, phone = ? WHERE id = ?"
	query := Cassandra.Session.Query(db_cmd, self.Company.ID, self.Name, self.Email, self.Phone, self.ID)
	err := query.Exec()
	if err != nil {
		return self, err.Error()
	}

	self.Invite()

	return self, ""
}

func (self Contact) Validate() string {
	if self.Name == "" {
		return "name is required"
	}

	if self.Email == "" {
		return "email is required"
	}

	if self.Company.IsNew() {
		return "company is required"
	}

	if self.IsNew() {
		added_contact, _ := GetByEmail(self.Company, self.Email)
		if added_contact.Email == self.Email {
			return "email already used"
		}
	}

	return ""
}

func (self Contact) Delete() (bool, string) {
	if self.IsNew() {
		return false, "Contact not found"
	}

	db_cmd := "DELETE FROM contact WHERE id = ?"
	query := Cassandra.Session.Query(db_cmd, self.ID)
	err := query.Exec()
	if err != nil {
		return false, err.Error()
	}

	return true, ""
}

func (self Contact) GetToken() int {
	db_cmd := "select token FROM contact WHERE email = ?"
	query := Cassandra.Session.Query(db_cmd, self.Email)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		return m["token"].(int)
	}

	return 0
}

func (self Contact) Invite() {
	token := self.GetToken()
	if token == 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		token = rand.Intn(999999)
	}

	go util.Invite(self.Email, self.Name, token)

	db_cmd := "UPDATE contact SET \"token\" = ? WHERE id = ?"
	query := Cassandra.Session.Query(db_cmd, token, self.ID)
	err := query.Exec()
	if err != nil {
		log.Println("[Error] Could save token to " + self.ID.String())
		log.Println("[Error] Reason " + err.Error())
	}

}

func GetBy(companyId string, field string, value string) (Contact, string) {
	var contact Contact

	db_cmd := "SELECT company_id, id, name, email, phone from contact WHERE " + field + " = ?"
	query := Cassandra.Session.Query(db_cmd, value)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		if m["company_id"].(gocql.UUID).String() != companyId {
			continue
		}

		cont_company, _ := company.GetById(companyId)

		contact = Contact{
			ContactBasic: ContactBasic{
				ID:    m["id"].(gocql.UUID),
				Name:  m["name"].(string),
				Email: m["email"].(string),
				Phone: m["phone"].(string),
			},
			Company: cont_company,
		}
	}

	if contact.IsNew() {
		return contact, "Not found"
	}

	return contact, ""
}

func GetAll(companyId string) []ContactBasic {
	var contacts []ContactBasic
	var id gocql.UUID
	var name, email, phone string
	cont_company, _ := company.GetById(companyId)

	db_cmd := "SELECT id, name, email, phone from contact WHERE company_id = ?"
	query := Cassandra.Session.Query(db_cmd, cont_company.ID)
	iterable := query.Iter()
	for iterable.Scan(&id, &name, &email, &phone) {
		contact := ContactBasic{
			ID:    id,
			Name:  name,
			Email: email,
			Phone: phone,
		}
		contacts = append(contacts, contact)
	}

	return contacts
}

func CalendarList(contact_id gocql.UUID) []ContactCalendar {
	var calendars []ContactCalendar
	var id, company_id gocql.UUID
	var start_at, end_at time.Time
	var description string
	var value float32

	db_cmd := "SELECT id, start_at, end_at, description, value, company_id from calendar WHERE contact_id = ?"
	iterable := Cassandra.Session.Query(db_cmd, contact_id).Iter()
	for iterable.Scan(&id, &start_at, &end_at, &description, &value, &company_id) {
		company, _ := company.GetById(company_id.String())

		calendar := ContactCalendar{
			ID:          id,
			StartAt:     start_at,
			EndAt:       end_at,
			Description: description,
			Value:       value,
			Company:     company,
		}
		calendars = append(calendars, calendar)
	}
	return calendars
}

func GetAllByEmail(email string) []ContactCalendar {
	var calendars []ContactCalendar
	db_cmd := "SELECT id from contact WHERE email = ?"
	query := Cassandra.Session.Query(db_cmd, email)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		contact_calendars := CalendarList(m["id"].(gocql.UUID))
		calendars = append(calendars, contact_calendars...)
	}
	return calendars
}

func GetByEmail(company company.CompanyBasic, email string) (Contact, string) {
	return GetBy(company.ID.String(), "email", email)
}

func GetById(companyId string, contactId string) (Contact, string) {
	return GetBy(companyId, "id", contactId)
}

func Auth(email string, token int) (bool, ContactBasic) {
	var id gocql.UUID
	var name, contact_email, phone string
	var contact_token int

	db_cmd := "SELECT \"id\", \"name\", \"email\", \"phone\", \"token\" from contact WHERE email = ?"
	query := Cassandra.Session.Query(db_cmd, email)
	iterable := query.Iter()
	for iterable.Scan(&id, &name, &contact_email, &phone, &contact_token) {
		if contact_token == token {
			return true, ContactBasic{
				ID:    id,
				Name:  name,
				Email: contact_email,
				Phone: phone,
			}
		}
	}

	return false, ContactBasic{}
}

type NewContactResponse struct {
	ID gocql.UUID
}

type SuccessResponse struct {
	Message string
}

type ErrorResponse struct {
	Error string
}

type ListResponse struct {
	Contacts []ContactBasic
}

type AuthResponse struct {
	Contact ContactBasic
}

type AuthRequest struct {
	Email string
	Token int
}

type ContactCalendar struct {
	ID          gocql.UUID
	StartAt     time.Time
	EndAt       time.Time
	Description string
	Value       float32
	Company     company.CompanyBasic
}

type ListCalendarResponse struct {
	Calendars []ContactCalendar
}
