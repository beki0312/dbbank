package api
// import (
// 	"github.com/gorilla/mux"
// 	"github.com/jackc/pgx/v4"
// )
// type BankServer struct {
// 	mux  			*mux.Router
// 	customerHandler *CustomerHandler
// }
// func NewBankServer(connect *pgx.Conn) *BankServer {
// 	return &BankServer{customerHandler: NewCustomerHandler(connect)}
// }
// func (s *BankServer) Start(port string)  {
// 	s.mux.HandleFunc("/customers",s.GetAllCustomers).Methods("GET")



	
// }