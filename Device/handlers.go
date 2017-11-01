package Device

import (
	"encoding/json"
	"log"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request) {
	device, err := DeviceForm(r)

	if err == nil {
		log.Println("[Device Register] Success:", device.Email)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(SuccessResponse{Message: "Success"})
	} else {
		log.Println("[Device Register] Error:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
	}
}