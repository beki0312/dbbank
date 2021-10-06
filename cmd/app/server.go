package app
import (
	"encoding/json"
	"fmt"
	"log"
	"mybankcli/pkg/api"
	"mybankcli/pkg/types"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)
type Server struct {
	mux 				*mux.Router
	customerHdr			*api.CustomerHandler
	managerHandler 		*api.ManagerHandler
}
//NewServer - функция-конструктор для создания нового сервера.
func NewServer(mux *mux.Router,customerHdr *api.CustomerHandler,managerHandler *api.ManagerHandler) *Server{
	return &Server{mux: mux,customerHdr: customerHdr,managerHandler: managerHandler}
}
//ServeHTTP - метод для запуска сервера.
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}
const (
	GET = "GET"
	POST = "POST"
	DELETE = "DELETE"
)
func (s *Server) Init()  {
	s.mux.HandleFunc("/customers",s.GetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers",s.PutCreateCustomer).Methods(POST)
	s.mux.HandleFunc("/customers/{id}",s.GetCustomersById).Methods(GET)
	s.mux.HandleFunc("/customers/{id}",s.GetDeleteCustomerById).Methods(DELETE)

	s.mux.HandleFunc("/managers",s.GetAllManagers).Methods(GET)
	s.mux.HandleFunc("/managers",s.PutCreateManager).Methods(POST)
	s.mux.HandleFunc("/managers/{id}",s.GetManagersById).Methods(GET)
	s.mux.HandleFunc("/managers/{id}",s.GetDeleteManagerById).Methods(GET)
}
//выводит список всех клиентов
func (s *Server) GetAllCustomers(w http.ResponseWriter, r *http.Request)  {
	cust,err:=s.customerHdr.CustomerAll(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,cust)
}
func (s *Server) GetCustomersById(w http.ResponseWriter, r *http.Request)  {
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
	item,err:=s.customerHdr.CustomerById(r.Context(),id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,item)
}
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
	item,err:=s.customerHdr.CustomerRemoveByID(r.Context(),id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,item)
}
func (s *Server) PutCreateCustomer(w http.ResponseWriter, r *http.Request)  {
	var customer *types.Customer
	err:=json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		log.Print(err)
		return
	}
	item,err:=s.customerHdr.CreateCustomer(r.Context(),customer)
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