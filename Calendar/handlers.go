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

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["companyId"]
	calendarId := vars["calendarId"]

	error_code := http.StatusNotFound
	calendar, errMsg := GetById(companyId, calendarId)
	if errMsg == "" {
		error_code = http.StatusBadRequest
		_, errMsg = calendar.Delete()
	}

	if errMsg == "" {
		log.Println("[Delete Calendar] Success:", calendar.ID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessResponse{Message: "Calendar removed"})
	} else {
		log.Println("[Delete Calendar] Error:", errMsg)
		w.WriteHeader(error_code)
		json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
	}
}

func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["companyId"]

	calendars := GetAll(companyId)

	log.Println("[List Calendar] Success:", companyId)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListResponse{Calendars: calendars})
}
