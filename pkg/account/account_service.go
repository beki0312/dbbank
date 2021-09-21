package account

import (
	"context"
	"fmt"
	"log"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"

	"github.com/jackc/pgx/v4"
)
type AccountService struct {
	accountRepository *AccountRepository
	transactionRepository *TransactionRepository
}
func NewAccountServicce(connect *pgx.Conn) *AccountService{
	return &AccountService{accountRepository: &AccountRepository{connect: connect},transactionRepository: &TransactionRepository{connect: connect}}
}
//Перевод 
func (s *AccountService) TransferMoneyByAccountId(payerAccountId,receiverAccountId int64, amount int64) error {
	payerAmount,err:=s.accountRepository.GetById(payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	if amount > payerAmount.Amount {
		log.Printf("не достаточно баланс")
		return err
	}
	receiverAmount,err:=s.accountRepository.GetById(receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	newPayerAmount:=payerAmount.Amount-amount
	newreceiverAmount:=receiverAmount.Amount+amount
	err=s.CreateTransaction(payerAccountId,receiverAccountId,amount)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	err=s.accountRepository.SetAmountById(newPayerAmount,payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	err=s.accountRepository.SetAmountById(newreceiverAmount,receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	fmt.Println("Перевод Успешно отправлено!!!")
	
	return nil
}
//Таблица транзаксия
func (s *AccountService) CreateTransaction(payerAccountId,receiverAccountId,amount int64) error {
	ctx:=context.Background()
	item:=types.Transactions{}
	err:=s.transactionRepository.connect.QueryRow(ctx, `insert into transactions (debet_account_id,credit_account_id,amount) values ($1,$2,$3) returning id,debet_account_id,credit_account_id,amount,date 
	`,payerAccountId,receiverAccountId,amount).Scan(&item.ID,&item.Debet_account_id,&item.Credit_account_id,&item.Amount,&item.Date)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return err
}