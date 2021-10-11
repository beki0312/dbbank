package account

import (
	"fmt"
	"log"
	"mybankcli/pkg/utils"
	"github.com/jackc/pgx/v4"
)
type AccountService struct {
	accountRepository *AccountRepository
	// transactionRepository *TransactionRepository
}
func NewAccountServicce(connect *pgx.Conn) *AccountService{
	return &AccountService{accountRepository: &AccountRepository{connect: connect}}
}
// func NewAccountServicce(connect *pgx.Conn) *AccountService{
// 	return &AccountService{accountRepository: &AccountRepository{connect: connect},transactionRepository: &TransactionRepository{connect: connect}}
// }
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
	_,err=s.accountRepository.CreateTransactions(payerAccountId,receiverAccountId,amount)
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
