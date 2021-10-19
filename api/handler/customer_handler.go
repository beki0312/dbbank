package handler

import (
	"context"
	"errors"
	"log"
	"mybankcli/pkg/account"
	"mybankcli/pkg/customers"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"

	"github.com/jackc/pgx/v4"
)

// Errors
var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

//Сервис - описывает обслуживание клиентов.
type CustomerHandler struct {
	connect            *pgx.Conn
	customerRepository *customers.CustomerRepository
	accountRepository  *account.AccountRepository
}

//NewServer - функция-конструктор для создания нового сервера.
func NewCustomerHandler(connect *pgx.Conn, customerRepository *customers.CustomerRepository, accountRepository *account.AccountRepository) *CustomerHandler {
	return &CustomerHandler{connect: connect, customerRepository: customerRepository, accountRepository: accountRepository}
}



//Регистрация Менеджера
func (h *CustomerHandler) RegistersCustomers(ctx context.Context, item *types.Registration) (*types.Customer, error) {
	cust := &types.Customer{}
	registration, err := h.customerRepository.Register(item)
	if err != nil {
		return nil, err
	}
	if registration == nil {
		return nil, ErrNotFound
	}
	return cust, err
}

//Найти токена менеджера
func (h *CustomerHandler) GetCustomerToken(ctx context.Context, item *types.Authers) (token string, err error) {
	// if item.Password != item.Password {
	// 	log.Println("TOken is Error")
	// }
	token, err = h.customerRepository.Token(item.Phone, item.Password)
	if err != nil {
		return "", err
	}
	return token, err
}

//find Id customers Token
func (s *CustomerHandler) IDByTokenCustomers(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.connect.QueryRow(ctx, `SELECT customer_id FROM customers_tokens WHERE token =$1`, token).Scan(&id)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, ErrInternal
	}
	return id, err
}
func (h *CustomerHandler) PostCustomers(ctx context.Context, customer *types.Customer) (*types.Customer, error) {
	if customer.ID <= 0 {
		return nil, ErrInternal
	}
	customers, err := h.customerRepository.CreateCustomers(customer)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if customers == nil {
		return nil, ErrNotFound
	}
	return customer, nil
}

//Список всех Клиентов
func (h *CustomerHandler) GetAllCustomer(ctx context.Context) ([]*types.Customer, error) {
	customers, err := h.customerRepository.Customers()
	if err != nil {
		return nil, ErrInternal
	}
	return customers, err
}

//Список всех активный клиентов
func (h *CustomerHandler) GetAllActiveCustomers(ctx context.Context) ([]*types.Customer, error) {
	customers, err := h.customerRepository.AllActiveCustomers()
	if err != nil {
		return nil, ErrInternal
	}
	return customers, err
}

//Списко Клиентов по их Id
func (h *CustomerHandler) GetCustomerById(ctx context.Context, id int64) (*types.Customer, error) {
	if id <= 0 {
		return nil, ErrInternal
	}
	customers, err := h.customerRepository.CustomerById(id)
	if err != nil {
		// log.Printf(("error while getting customer by id %e,%e"),id,err)
		return nil, ErrInternal
	}
	if customers == nil {
		return nil, ErrNotFound
	}
	return customers, nil
}

// Удаление клиентов по их Id
func (h *CustomerHandler) GetDeleteCustomerByID(ctx context.Context, id int64) (*types.Customer, error) {
	if id <= 0 {
		return nil, ErrInternal
	}
	customers, err := h.customerRepository.CustomersDeleteById(id)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	if customers == nil {
		return nil, ErrNotFound
	}
	return customers, nil
}

// Удаление счета по Id клиента
func (h *CustomerHandler) GetDeleteAccountByID(ctx context.Context, id int64) (*types.Account, error) {
	if id <= 0 {
		return nil, ErrInternal
	}
	accounts, err := h.customerRepository.AccountsDeleteById(id)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	if accounts == nil {
		return nil, ErrNotFound
	}
	return accounts, nil
}

//Перевод денег по номеру счета
func (h *CustomerHandler) PostTransferMoneyByAccount(ctx context.Context, item *types.AccountTransfer) (*types.Account, error) {
	var payerAccountId, receiverAccountId int64
	tx, err := h.connect.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	err = h.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, item.Payer_Accont).Scan(&payerAccountId)
	if err != nil {
		return nil, err
	}
	err = h.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, item.Receiver_Account).Scan(&receiverAccountId)
	if err != nil {
		return nil, err
	}
	accounts := &types.Account{}
	payerAmount, err := h.accountRepository.GetById(payerAccountId)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if item.Amount > payerAmount.Amount {
		log.Printf("не достаточно баланс")
		return nil, err
	}
	receiverAmount, err := h.accountRepository.GetById(receiverAccountId)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	newPayerAmount := payerAmount.Amount - item.Amount
	newreceiverAmount := receiverAmount.Amount + item.Amount
	err = h.accountRepository.CreateTransactionstx(tx, payerAccountId, receiverAccountId, item.Amount)
	if err != nil {
		tx.Rollback(context.Background())
		log.Print(err)
		return nil, err
	}
	err = h.accountRepository.SetAmountByIdtx(tx, newPayerAmount, payerAccountId)
	if err != nil {
		tx.Rollback(context.Background())
		log.Print(err)
		return nil, err
	}
	err = h.accountRepository.SetAmountById(newreceiverAmount, receiverAccountId)
	if err != nil {
		tx.Rollback(context.Background())
		log.Print(err)
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		utils.ErrCheck(err)
		return nil, err
	}

	return accounts, err
}

//Перевод денег по номеру телефона
func (h *CustomerHandler) PutTransferMoneyByPhone(ctx context.Context, item *types.AccountPhoneTransactions) (*types.Account, error) {
	tx, err := h.connect.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	var payerAccountId, receiverAccountId int64
	sql := `select account.id from account left join customer on customer.id=account.customer_id where customer.phone=$1`
	err = h.connect.QueryRow(ctx, sql, item.Payer_phone).Scan(&payerAccountId)
	if err != nil {
		return nil, err
	}
	err = h.connect.QueryRow(ctx, sql, item.Receiver_Phone).Scan(&receiverAccountId)
	if err != nil {
		return nil, err
	}
	accounts := &types.Account{}
	payerAmount, err := h.accountRepository.GetById(payerAccountId)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if item.Amount > payerAmount.Amount {
		log.Printf("не достаточно баланс")
		return nil, err
	}
	receiverAmount, err := h.accountRepository.GetById(receiverAccountId)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	newPayerAmount := payerAmount.Amount - item.Amount
	newreceiverAmount := receiverAmount.Amount + item.Amount
	err = h.accountRepository.CreateTransactionstx(tx, payerAccountId, receiverAccountId, item.Amount)
	if err != nil {
		tx.Rollback(context.Background())
		log.Print(err)
		return nil, err
	}
	err = h.accountRepository.SetAmountByIdtx(tx, newPayerAmount, payerAccountId)
	if err != nil {
		tx.Rollback(context.Background())
		log.Print(err)
		return nil, err
	}
	err = h.accountRepository.SetAmountByIdtx(tx, newreceiverAmount, receiverAccountId)
	if err != nil {
		tx.Rollback(context.Background())
		log.Print(err)
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		utils.ErrCheck(err)
		return nil, err
	}

	return accounts, err
}

//Транзаксия
func (h *CustomerHandler) GetTransaction(ctx context.Context) ([]*types.Transactions, error) {
	transaction, err := h.customerRepository.HistoryTansfer()
	if err != nil {
		return nil, err
	}
	return transaction, err
}

// Удалеит Токена по их Id
func (h *CustomerHandler) GetCustomersTokensRemoveByID(ctx context.Context, id int64) (*types.Tokens, error) {
	if id <= 0 {
		return nil, ErrInternal
	}
	tokens, err := h.customerRepository.CustomersTokenRemoveByID(id)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	if tokens == nil {
		return nil, ErrNotFound
	}
	return tokens, nil
}

//ВЫвод список банкоматов
func (h *CustomerHandler) GetAllAtm(ctx context.Context) ([]*types.Atm, error) {
	atm, err := h.customerRepository.CustomerAtm()
	if err != nil {
		return nil, ErrInternal
	}
	return atm, err
}

//Создание список Банкоматов
func (h *CustomerHandler) PostAtm(ctx context.Context, atm *types.Atm) (*types.Atm, error) {
	// if (atm.ID<=0) {
	// 	return nil,ErrInternal
	// }
	atms, err := h.customerRepository.CreateAtms(atm)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if atms == nil {
		return nil, ErrNotFound
	}
	return atm, nil
}
