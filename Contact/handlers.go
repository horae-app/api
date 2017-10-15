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

	contact, errMsg := GetById(companyId, contactId)
	error_code := http.StatusNotFound
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
