package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mybankcli/pkg/account"
	"mybankcli/pkg/customers"
	"mybankcli/pkg/types"
	"github.com/jackc/pgx/v4"
)

// Errors
var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

type CustomerHandler struct {
	connect  *pgx.Conn
	customerRepository  	*customers.CustomerRepository
	accountRepository 		*account.AccountRepository
}
func NewCustomerHandler(connect *pgx.Conn,customerRepository *customers.CustomerRepository,accountRepository *account.AccountRepository) *CustomerHandler {
	return &CustomerHandler{connect: connect,customerRepository: customerRepository,accountRepository: accountRepository}
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
	return customer,nil
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

func (h *CustomerHandler) PostTransferMoneyByAccount(ctx context.Context,payerAccountId,receiverAccountId, amount int64) (*types.Account,error) {
	// var payerAccountId, receiverAccountId int64
	// accountService:=account.NewAccountServicce(s.connect)
	// fmt.Println("Перевод по номеру счета")
	// err := h.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, payerAccount).Scan(&payerAccountId)
	// if err != nil {
	// 	return err
	// }
	// err = h.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, receiverAccount).Scan(&receiverAccountId)
	// if err != nil {
	// 	return err
	// }

	var Currency string
	var AccountName string
	accounts:=&types.Account{
		Currency_code: Currency,
		Account_Name: AccountName,
	}
	payerAmount,err:=h.accountRepository.GetById(payerAccountId)
	if err != nil {
		log.Print(err)
		return nil,err
	}
	if amount > payerAmount.Amount {
		log.Printf("не достаточно баланс")
		return nil,err
	}
	receiverAmount,err:=h.accountRepository.GetById(receiverAccountId)
	if err != nil {
		log.Print(err)
		return nil,err
	}
	newPayerAmount:=payerAmount.Amount-amount
	newreceiverAmount:=receiverAmount.Amount+amount
	_,err=h.accountRepository.CreateTransactions(payerAccountId,receiverAccountId,amount)
	if err != nil {
		log.Print(err)
		return nil,err
	}
	
	err=h.accountRepository.SetAmountById(newPayerAmount,payerAccountId)
	if err != nil {
		log.Print(err)
		return nil,err
	}
	err=h.accountRepository.SetAmountById(newreceiverAmount,receiverAccountId)
	if err != nil {
		log.Print(err)
		return nil,err
	}
	fmt.Println("Перевод Успешно отправлено!!!")
	return accounts,err
}

func (h *CustomerHandler) GetTransaction(ctx context.Context) ([]*types.Transactions,error) {
	transaction,err:=h.accountRepository.HistoryTansfer()
	if err != nil {
		return nil, err
	}

	return transaction,err
}