package Company

import (
	"encoding/json"
	"log"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request) {
	company, errMsg := UserForm(r)
	if errMsg == "" {
		company, errMsg = company.Save()
	}

	if errMsg == "" {
		log.Println("[New Company] Success:", company.ID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(NewCompanyResponse{ID: company.ID})
	} else {
		log.Println("[New Company] Error:", errMsg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	company, msg := AuthForm(r)
	if msg == "" {
		log.Println("[Auth] Success:", msg)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AuthResponse{Token: company.ID.String(), City: company.City, State: company.State})
	} else {
		log.Println("[Auth] Error:", msg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
	}
}
