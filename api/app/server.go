package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"mybankcli/api"
	"mybankcli/api/app/middlware"
	"mybankcli/pkg/types"
	"net/http"
	"strconv"
)

//Сервис - описывает обслуживание клиентов.
type Server struct {
	mux             *mux.Router
	customerHandler *api.CustomerHandler
	managerHandler  *api.ManagerHandler
	accountHandler  *api.AccountHandler
}

//NewServer - функция-конструктор для создания нового сервера.
func NewServer(mux *mux.Router, customerHandler *api.CustomerHandler, managerHandler *api.ManagerHandler, accountHandler *api.AccountHandler) *Server {
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
	customersSubrouter.HandleFunc("", s.CustomerRegistration).Methods(POST)
	customersSubrouter.HandleFunc("/token", s.GetCustomerTokens).Methods(POST)
	customersSubrouter.HandleFunc("token/{id}", s.GetDeleteCustomersTokensById).Methods(DELETE)
	customersSubrouter.HandleFunc("/", s.GetAllCustomers).Methods(GET)
	customersSubrouter.HandleFunc("/{id}", s.GetCustomersById).Methods(GET)
	customersSubrouter.HandleFunc("/{id}", s.GetDeleteCustomerById).Methods(DELETE)
	customersSubrouter.HandleFunc("/tranferaccount", s.PutTransferMoneyByAccounts).Methods(PUT)
	customersSubrouter.HandleFunc("/tranferPhone", s.PutTransferMoneyByPhones).Methods(PUT)
	customersSubrouter.HandleFunc("/accounts/{id}", s.GetAccountById).Methods(GET)
	customersSubrouter.HandleFunc("/accounts/{id}", s.GetDeleteAccountById).Methods(DELETE)
	customersSubrouter.HandleFunc("/atm", s.PostNewAtm).Methods(POST)

	s.mux.HandleFunc("/transactions", s.GetTransactions).Methods(GET)
	s.mux.HandleFunc("/accounts", s.GetAccountsAll).Methods(GET)
	s.mux.HandleFunc("/atm", s.GetAtmsAll).Methods(GET)

	managersAuth := middlware.Authenticate(s.managerHandler.IDByTokenManagers)
	managersSubRouter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubRouter.Use(managersAuth)
	managersSubRouter.HandleFunc("/", s.ManagerRegistration).Methods(POST)
	managersSubRouter.HandleFunc("/token", s.GetManagersTokens).Methods(POST)
	managersSubRouter.HandleFunc("/", s.GetAllManagers).Methods(GET)
	managersSubRouter.HandleFunc("/{id}", s.GetManagersById).Methods(GET)
	managersSubRouter.HandleFunc("/token/{id}", s.GetDeleteManagerTokensById).Methods(DELETE)
	managersSubRouter.HandleFunc("/{id}", s.GetDeleteManagerById).Methods(DELETE)
	managersSubRouter.HandleFunc("/accounts", s.PostNewAccounts).Methods(POST)

}

//Регистрация Клиентов
func (s *Server) CustomerRegistration(w http.ResponseWriter, r *http.Request) {
	var item *types.Registration
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return
	}
	_, err = s.customerHandler.RegistersCustomers(r.Context(), item)
	if err != nil {
		return
	}
	RespondJSON(w, item)
}

//Авторизация Клиента
func (s *Server) GetCustomerTokens(w http.ResponseWriter, r *http.Request) {
	var auther *types.Authers
	err := json.NewDecoder(r.Body).Decode(&auther)
	if err != nil {
		log.Print(err)
		return
	}
	token, err := s.customerHandler.GetCustomerToken(r.Context(), auther)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print(err)
		return

	}
	RespondJSON(w, token)
}

//Удалиение Токен Customers по их Id
func (s *Server) GetDeleteCustomersTokensById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := s.customerHandler.GetCustomersTokensRemoveByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

// выводит список всех клиентов
func (s *Server) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	cust, err := s.customerHandler.GetAllCustomer(r.Context())

	if err != nil {
		// w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	RespondJSON(w, cust)
}

//Вывод список клиентов по их Id
func (s *Server) GetCustomersById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println("err", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.customerHandler.GetCustomerById(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

// Удалить клиента по Id
func (s *Server) GetDeleteCustomerById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		fmt.Println("хато")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := s.customerHandler.GetDeleteCustomerByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

// Удалить Счет по Id клиента
func (s *Server) GetDeleteAccountById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := s.customerHandler.GetDeleteAccountByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
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

//Таблица транзаксия
func (s *Server) GetTransactions(w http.ResponseWriter, r *http.Request) {

	tansfer, err := s.customerHandler.GetTransaction(r.Context())
	if err != nil {
		log.Println(err)
		return
	}

	RespondJSON(w, tansfer)
}

//список счетов
func (s *Server) GetAccountsAll(w http.ResponseWriter, r *http.Request) {
	account, err := s.accountHandler.GetAccountsAll(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, account)
}

//Список счетов по Id
func (s *Server) GetAccountById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println("err", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.accountHandler.GetCustomerAccountById(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

//Список Банкоматов
func (s *Server) GetAtmsAll(w http.ResponseWriter, r *http.Request) {
	atm, err := s.customerHandler.GetAllAtm(r.Context())
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w, atm)
}

//Добавление счет клиента
func (s *Server) PostNewAccounts(w http.ResponseWriter, r *http.Request) {
	var account *types.Account
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		log.Print(err)
		return
	}
	item, err := s.customerHandler.PostAccounts(r.Context(), account)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w, item)
}

//Добавление список банкомата
func (s *Server) PostNewAtm(w http.ResponseWriter, r *http.Request) {
	var atm *types.Atm
	err := json.NewDecoder(r.Body).Decode(&atm)
	if err != nil {
		log.Print(err)
		return
	}
	item, err := s.customerHandler.PostAtm(r.Context(), atm)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w, item)
}

//Регистрация
func (s *Server) ManagerRegistration(w http.ResponseWriter, r *http.Request) {
	var managers *types.Registration
	err := json.NewDecoder(r.Body).Decode(&managers)
	if err != nil {
		return
	}
	_, err = s.managerHandler.RegistersManagers(r.Context(), managers)
	if err != nil {
		return
	}
	RespondJSON(w, managers)
}

//Авторизация Менеджера
func (s *Server) GetManagersTokens(w http.ResponseWriter, r *http.Request) {
	var auther *types.Authers
	err := json.NewDecoder(r.Body).Decode(&auther)
	if err != nil {
		log.Print(err)
		return
	}
	token, err := s.managerHandler.GetManagersToken(r.Context(), auther)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print(err)
		return
	}
	RespondJSON(w, token)
}

//Список Всех Менеджеров
func (s *Server) GetAllManagers(w http.ResponseWriter, r *http.Request) {
	managers, err := s.managerHandler.GetManagersAll(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, managers)
}

//Список Менеджеров по их Id
func (s *Server) GetManagersById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		fmt.Println("khato")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println("err", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.managerHandler.GetManagersById(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

//Удалиение менеджеров по их Id
func (s *Server) GetDeleteManagerById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := s.managerHandler.GetManagersRemoveByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

//Удалиение Токен менеджера по их Id
func (s *Server) GetDeleteManagerTokensById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := s.managerHandler.GetManagersTokensRemoveByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
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
