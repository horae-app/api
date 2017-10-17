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

type Contact struct {
	ID      gocql.UUID
	Company company.CompanyBasic
	Name    string
	Email   string
	Phone   string
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

func (self Contact) Invite() {
	var token int

	db_cmd := "select token FROM contact WHERE id = ?"
	query := Cassandra.Session.Query(db_cmd, self.ID)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		token = m["token"].(int)
	}

	if token == 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		token = rand.Intn(999999)
	}

	_, errMsg := util.Invite(self.Email, self.Name, token)
	if errMsg != "" {
		log.Println("[Error] Could not send invite to " + self.Email)
	}

	db_cmd = "UPDATE contact SET \"token\" = ? WHERE id = ?"
	query = Cassandra.Session.Query(db_cmd, token, self.ID)
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
			ID:      m["id"].(gocql.UUID),
			Name:    m["name"].(string),
			Company: cont_company,
			Email:   m["email"].(string),
			Phone:   m["phone"].(string),
		}
	}

	if contact.IsNew() {
		return contact, "Not found"
	}

	return contact, ""
}

func GetAll(companyId string) []Contact {
	var contacts []Contact
	cont_company, _ := company.GetById(companyId)

	db_cmd := "SELECT id, name, email, phone from contact WHERE company_id = ?"
	query := Cassandra.Session.Query(db_cmd, cont_company.ID)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		contact := Contact{
			ID:      m["id"].(gocql.UUID),
			Name:    m["name"].(string),
			Company: cont_company,
			Email:   m["email"].(string),
			Phone:   m["phone"].(string),
		}
		contacts = append(contacts, contact)
	}

	return contacts
}

func GetByEmail(company company.CompanyBasic, email string) (Contact, string) {
	return GetBy(company.ID.String(), "email", email)
}

func GetById(companyId string, contactId string) (Contact, string) {
	return GetBy(companyId, "id", contactId)
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
	Contacts []Contact
}
