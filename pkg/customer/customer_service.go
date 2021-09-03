package customer

import (
	"fmt"
	"github.com/jackc/pgx/v4"
)


type CustomerService struct {
	repository *CustomerRepository
}

func NewCustomerServicce(connect *pgx.Conn) *CustomerService{
	return &CustomerService{&CustomerRepository{connect: connect}}
}

func (s *CustomerService) CustomerAccount() error {
	var phone,password,pass string
	fmt.Print("input log: ")
	fmt.Scan(&phone)
	fmt.Print("input password: ")
	fmt.Scan(&password)
	fmt.Println("")
	customer,err:=s.repository.GetCustomerWithIdAndPass(phone,pass)
	if err != nil {
		fmt.Printf("can't get password customer %e",err)
		return err
	}
	// if customer == nil {
	// 	fmt.Println("customer login and password can't input")
	// 	return errors.New("customer not found")
	// }
	fmt.Println("Good luck customer")
	fmt.Println("")
	s.Loop(customer)
	return nil
}

func (s *CustomerService) Loop(customer Customer)  {
	
}