package Calendar

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["companyId"]
	contactId := vars["contactId"]

	calendar, errMsg := CalendarForm(r, companyId, contactId)
	if errMsg == "" {
		calendar, errMsg = calendar.Save()
	}

	if errMsg == "" {
		log.Println("[New Calendar] Success:", calendar.ID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(NewCalendarResponse{ID: calendar.ID})
	} else {
		log.Println("[New Calendar] Error:", errMsg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
	}
}
