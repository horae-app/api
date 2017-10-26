package Contact

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["companyId"]

	contact, errMsg := ContactForm(r, companyId)
	if errMsg == "" {
		contact, errMsg = contact.Save()
	}

	if errMsg == "" {
		log.Println("[New Contact] Success:", contact.ID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(NewContactResponse{ID: contact.ID})
	} else {
		log.Println("[New Contact] Error:", errMsg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["companyId"]
	contactId := vars["contactId"]

	error_code := http.StatusNotFound
	contact, errMsg := GetById(companyId, contactId)
	if errMsg == "" {
		error_code = http.StatusBadRequest
		_, errMsg = contact.Delete()
	}

	if errMsg == "" {
		log.Println("[Delete Contact] Success:", contact.ID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Message: "Contact removed"})
	} else {
		log.Println("[Delete Contact] Error:", errMsg)
		w.WriteHeader(error_code)
		json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
	}
}

func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["companyId"]

	contacts := GetAll(companyId)

	log.Println("[List Contact] Success:", companyId)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListResponse{Contacts: contacts})
}

func UserAuth(w http.ResponseWriter, r *http.Request) {
	errMsg, contact := ContactAuth(r)
	if errMsg != "" {
		log.Println("[Contact Auth] Error:", errMsg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
		return
	}

	log.Println("[Contact Auth] Success:", contact.Email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{Contact: contact})
}

func Calendar(w http.ResponseWriter, r *http.Request) {
	email, errMsg := CalendarForm(r)
	if errMsg != "" {
		log.Println("[Contact Calendar] Error:", errMsg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
		return
	}

	calendars := GetAllByEmail(email)

	log.Println("[Contact Calendar] Success:", email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListCalendarResponse{Calendars: calendars})
}
