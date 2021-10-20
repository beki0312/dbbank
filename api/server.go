package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"mybankcli/api/handler"
	"mybankcli/api/middlware"
	"mybankcli/pkg/types"
	"net/http"
)

//Сервис - описывает обслуживание клиентов.
type Server struct {
	mux             *mux.Router
	customerHandler *handler.CustomerHandler
	managerHandler  *handler.ManagerHandler
	accountHandler  *handler.AccountHandler
}

//NewServer - функция-конструктор для создания нового сервера.
func NewServer(mux *mux.Router, customerHandler *handler.CustomerHandler, managerHandler *handler.ManagerHandler, accountHandler *handler.AccountHandler) *Server {
	return &Server{mux: mux, customerHandler: customerHandler, managerHandler: managerHandler, accountHandler: accountHandler}
}

//ServeHTTP - метод для запуска сервера.
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

func (s *Server) Init() {
	customerAuth := middlware.Authenticate(s.customerHandler.TokenCustomers)
	customersSubrouter := s.mux.PathPrefix("/api/customers").Subrouter()
	customersSubrouter.Use(customerAuth)
	s.mux.HandleFunc("/registration", s.customerHandler.Registration).Methods(POST)
	customersSubrouter.HandleFunc("/token", s.customerHandler.CustomerTokens).Methods(POST)
	customersSubrouter.HandleFunc("/token/{id}", s.customerHandler.DeleteTokensById).Methods(DELETE)
	customersSubrouter.HandleFunc("/{id}", s.customerHandler.CustomerById).Methods(GET)
	customersSubrouter.HandleFunc("/{id}", s.customerHandler.DeleteCustomerById).Methods(DELETE)
	customersSubrouter.HandleFunc("/tranferaccount", s.TransferMoneyByAccounts).Methods(PUT)
	customersSubrouter.HandleFunc("/tranferPhone", s.TransferMoneyByPhones).Methods(PUT)
	customersSubrouter.HandleFunc("/accounts/{id}", s.accountHandler.GetAccountById).Methods(GET)
	customersSubrouter.HandleFunc("/accounts/{id}", s.customerHandler.DeleteAccountById).Methods(DELETE)

	s.mux.HandleFunc("/transactions", s.customerHandler.Transactions).Methods(GET)
	s.mux.HandleFunc("/accounts", s.accountHandler.GetAllAccounts).Methods(GET)
	s.mux.HandleFunc("/atm", s.customerHandler.GetAllAtms).Methods(GET)

	managersAuth := middlware.Authenticate(s.managerHandler.TokenManagers)
	managersSubRouter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubRouter.Use(managersAuth)
	s.mux.HandleFunc("/ManagerRegister", s.managerHandler.Registration).Methods(POST)
	managersSubRouter.HandleFunc("/token", s.managerHandler.ManagerToken).Methods(POST)
	managersSubRouter.HandleFunc("/", s.managerHandler.GetAllManagers).Methods(GET)
	managersSubRouter.HandleFunc("/customers", s.customerHandler.GetAllCustomers).Methods(GET)
	managersSubRouter.HandleFunc("/{id}", s.managerHandler.GetManagerById).Methods(GET)
	managersSubRouter.HandleFunc("/token/{id}", s.managerHandler.DeleteTokenById).Methods(DELETE)
	managersSubRouter.HandleFunc("/{id}", s.managerHandler.DeleteManagerById).Methods(DELETE)
	managersSubRouter.HandleFunc("/accounts", s.accountHandler.PostNewAccounts).Methods(POST)
	managersSubRouter.HandleFunc("/atm", s.customerHandler.PostNewAtm).Methods(POST)
}

//Перевод по номеру счета
func (s *Server) TransferMoneyByPhones(w http.ResponseWriter, r *http.Request) {
	var accounts *types.AccountPhoneTransactions
	err := json.NewDecoder(r.Body).Decode(&accounts)
	if err != nil {
		log.Print(err)
		return
	}
	_, err = s.customerHandler.PutTransferMoneyByPhone(r.Context(), accounts)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w, accounts)
}

//Перевод по номери телефона
func (s *Server) TransferMoneyByAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts *types.AccountTransfer
	err := json.NewDecoder(r.Body).Decode(&accounts)
	if err != nil {
		log.Print(err)
		return
	}
	_, err = s.customerHandler.PostTransferMoneyByAccount(r.Context(), accounts)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w, accounts)
}

//respondJSON - ответ от JSON.
func RespondJSON(w http.ResponseWriter, item interface{}) {
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
