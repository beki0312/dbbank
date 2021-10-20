package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"mybankcli/pkg/account"
	"mybankcli/pkg/customers"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

//Регистрация Клиентов
func (h *CustomerHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var item *types.Registration
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return
	}
	_, err = h.customerRepository.Register(r.Context(), item)
	if err != nil {
		return
	}
	RespondJSON(w, item)
}

//Авторизация Клиента
func (h *CustomerHandler) CustomerTokens(w http.ResponseWriter, r *http.Request) {
	var auther *types.Authers
	err := json.NewDecoder(r.Body).Decode(&auther)
	if err != nil {
		log.Print(err)
		return
	}
	token, err := h.customerRepository.Token(r.Context(), auther.Phone, auther.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print(err)
		return

	}
	RespondJSON(w, token)
}

//Удалиение Токен Customers по их Id
func (h *CustomerHandler) DeleteTokensById(w http.ResponseWriter, r *http.Request) {

	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := h.customerRepository.CustomersTokenRemoveByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

//Вывод список клиентов по их Id
func (h *CustomerHandler) CustomerById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println("err", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := h.customerRepository.CustomerById(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

// Удалить клиента по Id
func (h *CustomerHandler) DeleteCustomerById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := h.customerRepository.CustomersDeleteById(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

// Удалить Счет по Id клиента
func (h *CustomerHandler) DeleteAccountById(w http.ResponseWriter, r *http.Request) {
	idparam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := h.customerRepository.AccountsDeleteById(r.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	RespondJSON(w, item)
}

//Таблица транзаксия
func (h *CustomerHandler) Transactions(w http.ResponseWriter, r *http.Request) {

	tansfer, err := h.customerRepository.HistoryTansfer(r.Context())
	if err != nil {
		log.Println(err)
		return
	}

	RespondJSON(w, tansfer)
}

//Список Банкоматов
func (h *CustomerHandler) GetAllAtms(w http.ResponseWriter, r *http.Request) {
	atm, err := h.customerRepository.CustomerAtm(r.Context())
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w, atm)
}

// выводит список всех клиентов
func (h *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	cust, err := h.customerRepository.Customers(r.Context())

	if err != nil {
		// w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	RespondJSON(w, cust)
}

//Id customers Token
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

//Добавление список банкомата
func (h *CustomerHandler) PostNewAtm(w http.ResponseWriter, r *http.Request) {
	var atm *types.Atm
	err := json.NewDecoder(r.Body).Decode(&atm)
	if err != nil {
		log.Print(err)
		return
	}
	item, err := h.customerRepository.CreateAtms(r.Context(), atm)
	if err != nil {
		log.Print(err)
		return
	}
	RespondJSON(w, item)
}
