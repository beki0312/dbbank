package account

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"mybankcli/pkg/utils"
)

//Сервис - описывает обслуживание клиентов.
type AccountService struct {
	accountRepository *AccountRepository
}

//NewServer - функция-конструктор для создания нового сервера.
func NewAccountServicce(connect *pgx.Conn) *AccountService {
	return &AccountService{accountRepository: &AccountRepository{connect: connect}}
}

//Перевод по номеру счета
func (s *AccountService) TransferMoneyByAccountId(payerAccountId, receiverAccountId int64, amount int64) error {
	payerAmount, err := s.accountRepository.GetById(payerAccountId)
	if err != nil {
		log.Print("не уадлось вывести номера счета по Id")
		utils.ErrCheck(err)
		return err
	}
	if payerAmount.Amount <= 0 {
		log.Print("Введите сумму больше 0")
		return err
	}
	if amount > payerAmount.Amount {
		log.Printf("не достаточно баланс")
		return err
	}
	receiverAmount, err := s.accountRepository.GetById(receiverAccountId)
	if err != nil {
		log.Print("Не удалос найти номер счета получателя по Id")
		utils.ErrCheck(err)
		return err
	}
	newPayerAmount := payerAmount.Amount - amount
	newreceiverAmount := receiverAmount.Amount + amount
	err = s.accountRepository.CreateTransactions(payerAccountId, receiverAccountId, amount)
	if err != nil {
		log.Print("Не удалось создать транзакции")
		utils.ErrCheck(err)
		return err
	}
	err = s.accountRepository.SetAmountById(newPayerAmount, payerAccountId)
	if err != nil {
		log.Print("Ошибка при переводе отправителья")
		utils.ErrCheck(err)
		return err
	}
	err = s.accountRepository.SetAmountById(newreceiverAmount, receiverAccountId)
	if err != nil {
		log.Print("Ошибка при перевода получателя")
		utils.ErrCheck(err)
		return err
	}
	fmt.Println("Перевод Успешно отправлено!!!")
	return nil
}
