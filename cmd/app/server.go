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
}
//NewServer - функция-конструктор для создания нового сервера.
func NewServer(mux *mux.Router,customerHdr *api.CustomerHandler) *Server{
	return &Server{mux: mux,customerHdr: customerHdr}
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
	s.mux.HandleFunc("/customers",s.CreateCustomer).Methods(POST)
	s.mux.HandleFunc("/customers/{id}",s.GetCustomersById).Methods(GET)
	s.mux.HandleFunc("/customers/{id}",s.GetDeleteById).Methods(DELETE)
}
//выводит список всех клиентов
func (s *Server) GetAllCustomers(w http.ResponseWriter, r *http.Request)  {
	cust,err:=s.customerHdr.All(r.Context())
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
	item,err:=s.customerHdr.ById(r.Context(),id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,item)
}
func (s *Server) GetDeleteById(w http.ResponseWriter, r *http.Request)  {
	idParam:=r.URL.Query().Get("id")
	id,err:=strconv.ParseInt(idParam,10,64)
	if err != nil {
		log.Println(err)
		return
	}
	item,err:=s.customerHdr.RemoveByID(r.Context(),id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w,item)
}
func (s *Server) CreateCustomer(w http.ResponseWriter, r *http.Request)  {
	var customer *types.Customer
	err:=json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		log.Print(err)
		return
	}
	item,err:=s.customerHdr.Save(r.Context(),customer)
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