package customer

import (
	"context"
	"log"
	"github.com/jackc/pgx/v4"
)


type CustomerRepository struct {
	connect *pgx.Conn
}
//
func (r *CustomerRepository) GetCustomerWithIdAndPass (phone,pass string) (Customer,error) {
	var customer Customer
	ctx:=context.Background()
	err:=r.connect.QueryRow(ctx,`select *from customer where phone=$1 and pass=$2`,
	phone).Scan(&customer.ID,&customer.Name,&customer.SurName,&customer.Phone,&customer.Active,&customer.Created)
	if err != nil {
		log.Printf("You entered the wrong username or password %e",err)
	}

	return customer,err
}