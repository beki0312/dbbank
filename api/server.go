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
	customerAuth := middlware.Authenticate(s.customerHandler.IDByTokenCustomers)
	customersSubrouter := s.mux.PathPrefix("/api/customers").Subrouter()
	customersSubrouter.Use(customerAuth)
	s.mux.HandleFunc("/registration", s.customerHandler.CustomerRegistration).Methods(POST)
	customersSubrouter.HandleFunc("/token", s.customerHandler.GetCustomerTokens).Methods(POST)
	customersSubrouter.HandleFunc("/token/{id}", s.customerHandler.GetDeleteCustomersTokensById).Methods(DELETE)
	customersSubrouter.HandleFunc("/{id}", s.customerHandler.GetCustomersById).Methods(GET)
	customersSubrouter.HandleFunc("/{id}", s.customerHandler.GetDeleteCustomerById).Methods(DELETE)
	customersSubrouter.HandleFunc("/tranferaccount", s.PutTransferMoneyByAccounts).Methods(PUT)
	customersSubrouter.HandleFunc("/tranferPhone", s.PutTransferMoneyByPhones).Methods(PUT)
	customersSubrouter.HandleFunc("/accounts/{id}", s.accountHandler.GetAccountById).Methods(GET)
	customersSubrouter.HandleFunc("/accounts/{id}", s.customerHandler.GetDeleteAccountById).Methods(DELETE)

	s.mux.HandleFunc("/transactions", s.customerHandler.GetTransactions).Methods(GET)
	s.mux.HandleFunc("/accounts", s.accountHandler.GetAccountsAll).Methods(GET)
	s.mux.HandleFunc("/atm", s.customerHandler.GetAtmsAll).Methods(GET)

	managersAuth := middlware.Authenticate(s.managerHandler.IDByTokenManagers)
	managersSubRouter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubRouter.Use(managersAuth)
	s.mux.HandleFunc("/ManagerRegister", s.managerHandler.ManagerRegistration).Methods(POST)
	managersSubRouter.HandleFunc("/token", s.managerHandler.GetManagersTokens).Methods(POST)
	managersSubRouter.HandleFunc("/", s.managerHandler.GetAllManagers).Methods(GET)
	managersSubRouter.HandleFunc("/customers", s.customerHandler.GetAllCustomers).Methods(GET)
	managersSubRouter.HandleFunc("/{id}", s.managerHandler.GetManagersById).Methods(GET)
	managersSubRouter.HandleFunc("/token/{id}", s.managerHandler.GetDeleteManagerTokensById).Methods(DELETE)
	managersSubRouter.HandleFunc("/{id}", s.managerHandler.GetDeleteManagerById).Methods(DELETE)
	managersSubRouter.HandleFunc("/accounts", s.accountHandler.PostNewAccounts).Methods(POST)
	managersSubRouter.HandleFunc("/atm", s.customerHandler.PostNewAtm).Methods(POST)
}

//Перевод по номеру счета
func (s *Server) PutTransferMoneyByPhones(w http.ResponseWriter, r *http.Request) {
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
func (s *Server) PutTransferMoneyByAccounts(w http.ResponseWriter, r *http.Request) {
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
