package delivery

import (
	"TPForum/internal/pkg/domain"
	"TPForum/internal/pkg/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type ServiceDelivery struct {
	SUcase domain.ServiceUsecase
}

func NewServiceDelivery(router *mux.Router, usecase domain.ServiceUsecase) {
	serviceDelivery := &ServiceDelivery{
		SUcase: usecase,
	}

	router.HandleFunc("/service/clear", serviceDelivery.Clear).Methods("POST", "OPTIONS")
	router.HandleFunc("/service/status", serviceDelivery.Status).Methods("GET", "OPTIONS")
}

func (s *ServiceDelivery) Clear(w http.ResponseWriter, r *http.Request)  {

	err := s.SUcase.Clear()
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	utils.WriteResponse(w, http.StatusOK, nil)
}

func (s *ServiceDelivery) Status(w http.ResponseWriter, r *http.Request) {
	service, err := s.SUcase.Status()

	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	data, err := json.Marshal(service)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}