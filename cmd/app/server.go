package app

import (
	"encoding/json"
	"fmt"
	"log"
	"mybankcli/api"
	"mybankcli/pkg/types"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// "golang.org/x/crypto/bcrypt"
)
type Server struct {
	mux 				*mux.Router
	customerHandler		*api.CustomerHandler
	managerHandler 		*api.ManagerHandler
	accountHandler 		*api.AccountHandler
}
//NewServer - функция-конструктор для создания нового сервера.
func NewServer(mux *mux.Router,customerHandler *api.CustomerHandler,managerHandler *api.ManagerHandler,accountHandler *api.AccountHandler) *Server{
	return &Server{mux: mux,customerHandler: customerHandler,managerHandler: managerHandler,accountHandler: accountHandler}
}
//ServeHTTP - метод для запуска сервера.
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}
const (
	GET = "GET"
	POST = "POST"
	DELETE = "DELETE"
	PUT="PUT"
)
func (s *Server) Init()  {
	s.mux.HandleFunc("/customers",s.GetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers",s.PostCustomers).Methods(POST)
	s.mux.HandleFunc("/customers/{id}",s.GetCustomersById).Methods(GET)
	s.mux.HandleFunc("/customers/{id}",s.GetDeleteCustomerById).Methods(DELETE)
	s.mux.HandleFunc("/customers/account/tranferaccount",s.PutTransferMoneyByAccounts).Methods(PUT)
	s.mux.HandleFunc("/customers/account/tranferPhone",s.PutTransferMoneyByPhones).Methods(PUT)
	s.mux.HandleFunc("/transaction",s.GetTransactions).Methods(GET)
	s.mux.HandleFunc("/accounts",s.GetAccountsAll).Methods(GET)
	s.mux.HandleFunc("/accounts/",s.PostNewAccounts).Methods(POST)
	s.mux.HandleFunc("/accounts/{id}",s.GetAccountById).Methods(GET)

	s.mux.HandleFunc("/atm",s.GetAtmsAll).Methods(GET)
	s.mux.HandleFunc("/atm",s.PostNewAtm).Methods(POST)




	s.mux.HandleFunc("/managers",s.GetAllManagers).Methods(GET)
	s.mux.HandleFunc("/managers",s.PutCreateManager).Methods(POST)
	s.mux.HandleFunc("/managers/{id}",s.GetManagersById).Methods(GET)
	s.mux.HandleFunc("/managers/{id}",s.GetDeleteManagerById).Methods(GET)
}
// func (s *Server) handleCustomerRegistration(w http.ResponseWriter, r *http.Request) {
// 	var item *types.Registration
// 	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
// 		errWriter(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	hashed, err := bcrypt.GenerateFromPassword([]byte(item.Password),14)
// 	if err != nil {
// 		errWriter(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	item.Password = string(hashed)
// 	saved, err := s.customerHandler.RegistersCustomers(r.Context(), item)
// 	if err != nil {
// 		errWriter(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	RespondJSON(w, saved)
// }

// func (s *Server) handleCustomerGetToken(w http.ResponseWriter, r *http.Request) {
// 	// var item *customers.Auth 
// 	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
// 		errWriter(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	token, err := s.customerHandler.Token(r.Context(), item.Login, item.Password)
// 		if err != nil {
// 		errWriter(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	RespondJSON(w, map[string]interface{}{"status": "ok", "token": token})
// }

// выводит список всех клиентов
func (s *Server) GetAllCustomers(w http.ResponseWriter, r *http.Request)  {
	cust,err:=s.customerHandler.GetAllCustomer(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,cust)
}
//Вывод список по их Id
func (s *Server) GetCustomersById(w http.ResponseWriter, r *http.Request)  {
	idparam,ok:=mux.Vars(r)["id"]
	if  !ok {
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return 
	}
	id,err:=strconv.ParseInt(idparam,10,64)
	if err != nil {
		log.Println("err",err)
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	item,err:=s.customerHandler.GetCustomerById(r.Context(),id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,item)
}
// Удалить клиента по Id
func (s *Server) GetDeleteCustomerById(w http.ResponseWriter, r *http.Request)  {
	idparam,ok:=mux.Vars(r)["id"]
	if  !ok {
		fmt.Println("хато")
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return 
	}	
	id,err:=strconv.ParseInt(idparam,10,64)
	if err != nil {
		log.Println(err)
		return
	}
	item,err:=s.customerHandler.GetDeleteCustomerByID(r.Context(),id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,item)
}
//Регистрация нового клиента
func (s *Server) PostCustomers(w http.ResponseWriter, r *http.Request)  {
	var customer *types.Customer
	err:=json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		log.Print(err)
		return
	}
	item,err:=s.customerHandler.PostCustomers(r.Context(),customer)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w,item)
}
//Перевод по номеру счета
func(s *Server) PutTransferMoneyByPhones(w http.ResponseWriter, r *http.Request)  {
	var accounts *types.AccountPhoneTransactions
	err:=json.NewDecoder(r.Body).Decode(&accounts)
	if err != nil {
		log.Print(err)
		return
	}
	_,err=s.customerHandler.PutTransferMoneyByPhone(r.Context(),accounts)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w,accounts)	
}
//Перевод по номери телефона
func (s *Server) PutTransferMoneyByAccounts(w http.ResponseWriter, r *http.Request)  {
	var accounts *types.AccountTransfer
	err:=json.NewDecoder(r.Body).Decode(&accounts)
	if err != nil {
		log.Print(err)
		return
	}
	_,err=s.customerHandler.PostTransferMoneyByAccount(r.Context(),accounts)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w,accounts)	
}
//Таблица транзаксия
func (s *Server) GetTransactions(w http.ResponseWriter, r *http.Request) {
	tansfer,err:=s.customerHandler.GetTransaction(r.Context())
	if err != nil {
		log.Println(err)
		return
	}

	RespondJSON(w,tansfer)
}
//список счетов
func (s *Server) GetAccountsAll(w http.ResponseWriter, r *http.Request)  {
	account,err:=s.accountHandler.GetAccountsAll(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,account)
}
//Список счетов по Id
func (s *Server) GetAccountById(w http.ResponseWriter, r *http.Request)  {
	idparam,ok:=mux.Vars(r)["id"]
	if  !ok {
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return 
	}
	id,err:=strconv.ParseInt(idparam,10,64)
	if err != nil {
		log.Println("err",err)
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	item,err:=s.accountHandler.GetCustomerAccountById(r.Context(),id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,item)
}
//Список Банкоматов
func (s *Server) GetAtmsAll(w http.ResponseWriter, r *http.Request)  {
	atm,err:=s.customerHandler.GetAllAtm(r.Context())
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w,atm)
}
//Добавление счет клиента
func (s *Server) PostNewAccounts(w http.ResponseWriter, r *http.Request)  {
	var account *types.Account
	err:=json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		log.Print(err)
		return
	}
	item,err:=s.customerHandler.PostAccounts(r.Context(),account)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w,item)
}
//Регистрация нового клиента
func (s *Server) PostNewAtm(w http.ResponseWriter, r *http.Request)  {
	var atm *types.Atm
	err:=json.NewDecoder(r.Body).Decode(&atm)
	if err != nil {
		log.Print(err)
		return
	}
	item,err:=s.customerHandler.PostAtm(r.Context(),atm)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w,item)
}




func (s *Server) GetAllManagers(w http.ResponseWriter, r *http.Request)  {
	managers,err:=s.managerHandler.ManagersAll(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,managers)
}
func (s *Server) GetManagersById(w http.ResponseWriter, r *http.Request)  {
	idparam,ok:=mux.Vars(r)["id"]
	if  !ok {
		fmt.Println("khato")
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return 
	}
	id,err:=strconv.ParseInt(idparam,10,64)
	if err != nil {
		log.Println("err",err)
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return
	}
	item,err:=s.managerHandler.ManagersById(r.Context(),id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,item)
}
func (s *Server) GetDeleteManagerById(w http.ResponseWriter, r *http.Request)  {
	idparam,ok:=mux.Vars(r)["id"]
	if  !ok {
		fmt.Println("khato")
		http.Error(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		return 
	}
	id,err:=strconv.ParseInt(idparam,10,64)
	if err != nil {
		log.Println(err)
		return
	}
	item,err:=s.managerHandler.ManagersRemoveByID(r.Context(),id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,item)
}
func (s *Server) PutCreateManager(w http.ResponseWriter, r *http.Request)  {
	var managers *types.Manager
	err:=json.NewDecoder(r.Body).Decode(&managers)
	if err != nil {
		log.Print(err)
		return
	}
	item,err:=s.managerHandler.CreateManagers(r.Context(),managers)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w,item)
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


func errWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}