package Company

import (
	"github.com/gocql/gocql"
	"github.com/horae-app/api/Cassandra"
)

type CompanyBasic struct {
	ID    gocql.UUID
	Email string
	Name  string
	City  string
	State string
}

type CompanyFull struct {
	CompanyBasic
	Password string
}

func (self CompanyBasic) IsNew() bool {
	return self.ID.String() == "00000000-0000-0000-0000-000000000000"
}

func (self CompanyFull) Save() (CompanyFull, string) {
	errMsg := self.Validate()
	if errMsg != "" {
		return self, errMsg
	}

	if self.IsNew() {
		self.ID = gocql.TimeUUID()
	}

	db_cmd := "UPDATE company SET id = ?, name = ?, password = ?, city = ?, state = ? WHERE email = ?"
	query := Cassandra.Session.Query(db_cmd, self.ID, self.Name, self.Password, self.City, self.State, self.Email)
	err := query.Exec()
	if err != nil {
		return self, err.Error()
	}

	return self, ""
}

func (self CompanyFull) Validate() string {
	if self.Name == "" {
		return "name is required"
	}

	if self.Email == "" {
		return "email is required"
	}

	if self.Password == "" {
		return "password is required"
	}

	if self.IsNew() {
		added_company, _ := GetByEmail(self.Email)
		if added_company.Email == self.Email {
			return "email already used"
		}
	}

	return ""
}

func GetByEmail(email string) (CompanyBasic, string) {
	return GetBy("email", email)
}

func GetById(id string) (CompanyBasic, string) {
	return GetBy("id", id)
}

func GetPassword(email string) (string, string) {
	db_cmd := "SELECT password from company WHERE email = ?"
	query := Cassandra.Session.Query(db_cmd, email)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		return m["password"].(string), ""
	}

	return "", "Not found"
}

func GetBy(field string, value string) (CompanyBasic, string) {
	var company CompanyBasic

	db_cmd := "SELECT id, name, city, state, email from company WHERE " + field + " = ?"
	query := Cassandra.Session.Query(db_cmd, value)
	iterable := query.Iter()
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		company = CompanyBasic{
			ID:    m["id"].(gocql.UUID),
			Name:  m["name"].(string),
			City:  m["city"].(string),
			State: m["state"].(string),
			Email: m["email"].(string),
		}
	}

	if company.IsNew() {
		return company, "Not found"
	}

	return company, ""
}

type AuthRequest struct {
	Email    string
	Password string
}

type AuthResponse struct {
	Token string
}

type NewCompanyResponse struct {
	ID gocql.UUID
}

type ErrorResponse struct {
	Error string
}
