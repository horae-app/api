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
