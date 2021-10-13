package api

import (
	"context"
	"errors"
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

func (h *CustomerHandler) RegistersCustomers(ctx context.Context,item *types.Registration) (*types.Customer, error) {
	// item:=types.Registration{}
	registration,err:=h.customerRepository.Register(item.Name,item.Phone,item.Password)
	if err != nil {
		return nil, err
	}
	return registration,err
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
//Перевод денег по номеру счета
func (h *CustomerHandler) PostTransferMoneyByAccount(ctx context.Context,item *types.AccountTransfer) (*types.Account,error) {
	var payerAccountId, receiverAccountId int64
	err := h.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, item.Payer_Accont).Scan(&payerAccountId)
	if err != nil {
		return nil,err
	}
	err = h.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, item.Receiver_Account).Scan(&receiverAccountId)
	if err != nil {
		return nil,err
	}
	accounts:=&types.Account{}
	payerAmount,err:=h.accountRepository.GetById(payerAccountId)
	if err != nil {
		log.Print(err)
		return nil,err
	}
	if item.Amount > payerAmount.Amount {
		log.Printf("не достаточно баланс")
		return nil,err
	}
	receiverAmount,err:=h.accountRepository.GetById(receiverAccountId)
	if err != nil {
		log.Print(err)
		return nil,err
	}
	newPayerAmount:=payerAmount.Amount-item.Amount
	newreceiverAmount:=receiverAmount.Amount+item.Amount
	err=h.accountRepository.CreateTransactions(payerAccountId,receiverAccountId,item.Amount)
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
	return accounts,err
}
//Перевод денег по номеру телефона
func(h *CustomerHandler) PutTransferMoneyByPhone(ctx context.Context, item *types.AccountPhoneTransactions) (*types.Account,error)  {
	var payerAccountId, receiverAccountId int64
	sql:=`select account.id from account left join customer on customer.id=account.customer_id where customer.phone=$1`
	err:=h.connect.QueryRow(ctx,sql,item.Payer_phone).Scan(&payerAccountId)
	if err != nil {
		return nil,err
	}
	err=h.connect.QueryRow(ctx,sql,item.Receiver_Phone).Scan(&receiverAccountId)
	if err != nil {
		return nil,err
	}
	accounts:=&types.Account{}
	payerAmount,err:=h.accountRepository.GetById(payerAccountId)
	if err != nil {
		log.Print(err)
		return nil,err
	}
	if item.Amount > payerAmount.Amount {
		log.Printf("не достаточно баланс")
		return nil,err
	}
	receiverAmount,err:=h.accountRepository.GetById(receiverAccountId)
	if err != nil {
		log.Print(err)
		return nil,err
	}
	newPayerAmount:=payerAmount.Amount-item.Amount
	newreceiverAmount:=receiverAmount.Amount+item.Amount
	err=h.accountRepository.CreateTransactions(payerAccountId,receiverAccountId,item.Amount)
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
	return accounts,err
}
//Транзаксия
func (h *CustomerHandler) GetTransaction(ctx context.Context) ([]*types.Transactions,error) {
	transaction,err:=h.accountRepository.HistoryTansfer()
	if err != nil {
		return nil, err
	}
return transaction,err
}
//ВЫвод список банкоматов
func (h *CustomerHandler) GetAllAtm(ctx context.Context) ([]*types.Atm,error) {
	atm,err:=h.customerRepository.CustomerAtm()
	if err != nil {
		return nil, ErrInternal
	}
	return atm,err
}
//Save customers by id
func (h *CustomerHandler) PostAtm(ctx context.Context, atm *types.Atm) (*types.Atm,error) {
	// if (atm.ID<=0) {
	// 	return nil,ErrInternal
	// }
	atms,err:=h.customerRepository.CreateAtms(atm)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if atms==nil {
		return nil,ErrNotFound
	}
	return atm,nil
}