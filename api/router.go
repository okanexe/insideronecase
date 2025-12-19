package api

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"insider/internal/service"

	"github.com/gorilla/mux"
)

type API struct {
	s *service.MessageService
}

func New(s *service.MessageService) *API {
	return &API{
		s: s,
	}
}

func (a *API) SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/start", a.StartMessageSending).Methods("POST")
	router.HandleFunc("/stop", a.StopMessageSending).Methods("POST")
	router.HandleFunc("/messages", a.GetSentMessages).Methods("GET")
	return router
}
