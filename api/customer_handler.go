package api

import (
	"context"
	"errors"
	"log"
	"mybankcli/pkg/customers"
	"mybankcli/pkg/types"

	"github.com/jackc/pgx/v4"
)

// Errors
var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

type CustomerHandler struct {
	customerRepository  *customers.CustomerRepository
	connect *pgx.Conn
}
// func NewCustomerHandler(connect *pgx.Conn) *CustomerHandler {
// 	return &CustomerHandler{customerRepository: }
// }
func NewCustomerHandler(connect *pgx.Conn,customerRepository *customers.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{connect: connect,customerRepository: customerRepository}
}

//Get All Customer
func (h *CustomerHandler) GetAllCustomer(ctx context.Context) ( []*types.Customer,error) {
	customers,err:=h.customerRepository.Customers()
	if err != nil {
		return nil, ErrInternal
	}
	return customers,err
}
//Get All Active Customers
func (h *CustomerHandler) GetAllActiveCustomers(ctx context.Context) ( []*types.Customer,error) {
customers,err:=h.customerRepository.AllActiveCustomers()
if err != nil {
	return nil, ErrInternal
}
return customers,err
}
//Get ById customer
func (h *CustomerHandler) GetCustomerById(ctx context.Context,id int64) (*types.Customer,error) {
	if (id<=0) {
		return nil, ErrInternal
	}
	customers,err:=h.customerRepository.CustomerById(id)
	if err != nil {
		// log.Printf(("error while getting customer by id %e,%e"),id,err)
		return nil,ErrInternal
	}
	if customers==nil {
		return nil,ErrNotFound
	}
	return customers,nil
}
// Delete customer by id
func (h *CustomerHandler) GetDeleteCustomerByID(ctx context.Context, id int64) (*types.Customer, error) {
	if (id<=0) {
		return nil,ErrInternal	
	}
	customers,err:=h.customerRepository.CustomersDeleteById(id)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	if customers==nil {
		return nil,ErrNotFound
	}
	return customers, nil
}
//Save customers by id
func (h *CustomerHandler) PostCustomers(ctx context.Context, customer *types.Customer) (*types.Customer,error) {
	if (customer.ID<=0) {
		return nil,ErrInternal
	}
	customers,err:=h.customerRepository.CreateCustomers(customer)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if customers==nil {
		return nil,ErrNotFound
	}
	return customers,nil
}
//Block and Unblock customer by his id
func (h *CustomerHandler) CustomerBlockAndUnblockById(ctx context.Context, id int64,active bool) (*types.Customer,error) {
	customers,err:=h.customerRepository.CustomerBlockAndUnblockById(id,active)
		if err != nil {
			log.Println(err)
			return nil, ErrInternal
		}
		return customers,nil
}