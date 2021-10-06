package api

import (
	"context"
	"errors"
	"log"
	"mybankcli/pkg/types"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// Errors
var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

type CustomerHandler struct {
	connect *pgx.Conn
}
func NewCustomerHandler(connect *pgx.Conn) *CustomerHandler {
	return &CustomerHandler{connect: connect}
}
//Get All Customer
func (s *CustomerHandler) CustomerAll(ctx context.Context) ( []*types.Customer,error) {
	customers:=[]*types.Customer{}
	rows,err:=s.connect.Query(ctx,`SELECT *FROM customer`)
	if err != nil {
		return nil, ErrInternal
	}
	// defer rows.Close()
	for rows.Next(){
		item:=&types.Customer{}
		err=rows.Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
		if err != nil {
			log.Println(err)
		}
		customers = append(customers, item)
	}
	return customers,nil
}
//Get All Active Customers
func (s *CustomerHandler) CustomerAllActive(ctx context.Context) ( []*types.Customer,error) {
	customers:=[]*types.Customer{}
	rows,err:=s.connect.Query(ctx,`SELECT *FROM customer where active=true`)
	if err != nil {
		return nil, ErrInternal
	}
	// defer rows.Close()
	for rows.Next(){
		item:=&types.Customer{}
		err=rows.Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
		if err != nil {
			log.Println(err)
		}
		customers = append(customers, item)
	}
	return customers,nil
}
//Get ById customer
func (s *CustomerHandler) CustomerById(ctx context.Context,id int64) (*types.Customer,error) {
	customers:=&types.Customer{}
	err:=s.connect.QueryRow(ctx,`select id,name,surname,phone,password,active,created from customer where id=$1`,
	id).Scan(&customers.ID,&customers.Name,&customers.SurName,&customers.Phone,&customers.Password,&customers.Active,&customers.Created)
	if err != nil {
		log.Println(err)
		return nil,ErrInternal
	}
	return customers,nil
}
// Delete customer by id
func (s *CustomerHandler) CustomerRemoveByID(ctx context.Context, id int64) (*types.Customer, error) {
	cust := &types.Customer{}
	err := s.connect.QueryRow(ctx, `DELETE FROM customer WHERE id = $1`, 
	id).Scan(&cust.ID, &cust.Name, &cust.SurName,&cust.Phone,&cust.Password,&cust.Active, &cust.Created)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return cust, nil
}
//Save customers by id
func (s *CustomerHandler) CreateCustomer(ctx context.Context, customer *types.Customer) (*types.Customer,error) {
	item:=&types.Customer{}
	pass,_:=bcrypt.GenerateFromPassword([]byte(item.Password),14)
	if customer.ID==0 {
		log.Println("Вы ввели неверный номер пожалуйста введите номер с 1 ")
	}else{
		err:=s.connect.QueryRow(ctx,`insert into customer(id,name,surname,phone,password) values($1,$2,$3,$4,$5) returning id,name,surname,phone,password,active,created`,
		customer.ID,customer.Name,customer.SurName,customer.Phone,pass).Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)	
		if err != nil {
			log.Print(err)
			return nil,ErrInternal
		}
	}
	return item,nil
}

//Block and Unblock customer by his id
func (s *CustomerHandler) CustomerBlockAndUnblockById(ctx context.Context, id int64,active bool) (*types.Customer,error) {
	customers:=&types.Customer{}
	err:=s.connect.QueryRow(ctx,`update customer set active =$1 where id=$2`,active,id).Scan(
		&customers.ID,&customers.Name,&customers.SurName,&customers.Phone,&customers.Password,&customers.Active,&customers.Created)
		if err != nil {
			log.Println(err)
			return nil, ErrInternal
		}
		return customers,nil
}