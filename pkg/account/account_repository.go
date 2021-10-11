package account

import (
	"context"
	"errors"
	// "fmt"
	"log"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"

	"github.com/jackc/pgx/v4"
)
type AccountRepository struct {
	connect *pgx.Conn
}
func NewAccountRepository(connect *pgx.Conn) *AccountRepository {
	return &AccountRepository{connect: connect}
}

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

func (s *AccountRepository) GetById(id int64) (types.AccountTransfer, error) {
    var account types.AccountTransfer
    err:=s.connect.QueryRow(context.Background(),`select id, customer_id,currency_code,account_name, amount from account where id=$1`,id).Scan(
		&account.Payer_Id,&account.Receiver_Id,&account.Amount)
	if err != nil {
    utils.ErrCheck(err)
    return account,err
	}
    return account, err
}
func (s *AccountRepository) SetAmountById(amount,id int64)  error{
	_,err:=s.connect.Exec(context.Background(),`update account set amount = $1 where id = $2`,amount,id)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return nil
}

//Таблица транзаксия
func (s *AccountRepository) CreateTransactions(payerAccountId,receiverAccountId,amount int64) (*types.Transactions,error) {
	ctx:=context.Background()
	item:=&types.Transactions{}
	err:=s.connect.QueryRow(ctx, `insert into transactions (debet_account_id,credit_account_id,amount) values ($1,$2,$3) returning id,debet_account_id,credit_account_id,amount,date 
	`,payerAccountId,receiverAccountId,amount).Scan(&item.ID,&item.Debet_account_id,&item.Credit_account_id,&item.Amount,&item.Date)
	if err != nil {
		utils.ErrCheck(err)
		return nil,err
	}
	return item,err
}
func(s *AccountRepository) HistoryTansfer() ([]*types.Transactions,error) {
	ctx:=context.Background()
	accounts:=[]*types.Transactions{}
	rows,err:=s.connect.Query(ctx,`select * from transactions`)
	if err != nil {
		return nil, ErrInternal
	}
	for rows.Next(){
		account:=&types.Transactions{}
		err=rows.Scan(&account.ID,&account.Debet_account_id,&account.Credit_account_id,&account.Amount,&account.Date)
		if err != nil {
			log.Println(err)
		}
		accounts=append(accounts,account)
	}
	return accounts,nil

}
func (s *AccountRepository) Accounts() ([]*types.Account,error) {
	ctx:=context.Background()
	accounts:=[]*types.Account{}
	rows,err:=s.connect.Query(ctx,`select * from account`)
	if err != nil {
		return nil, ErrInternal
	}
	for rows.Next(){
		account:=&types.Account{}
		err=rows.Scan(&account.ID,&account.Customer_Id,&account.Currency_code,&account.Account_Name,&account.Amount)
		if err != nil {
			log.Println(err)
		}
		accounts=append(accounts,account)
	}
	return accounts,nil
}

func (s *AccountRepository) GetAccountByCustomerPhone(customerId int64) (*types.Account,error) {
	accounts:=&types.Account{}
	ctx :=context.Background()
	err:=s.connect.QueryRow(ctx,`SELECT account.id,account.customer_id,account.currency_code, account.account_name,account.amount FROM account 
	JOIN customer ON account.customer_id = customer.id
	where customer.id=$1`,customerId).Scan(&accounts.ID,&accounts.Customer_Id,&accounts.Currency_code,&accounts.Account_Name,&accounts.Amount)
	if err != nil {
		// utils.ErrCheck(err)
		return nil,err
	}
	return accounts,nil
}