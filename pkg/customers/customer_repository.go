package customers

import (
	"context"
	"errors"
	"fmt"
	"mybankcli/pkg/account"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"github.com/jackc/pgx/v4"
)
type CustomerRepository struct {
	connect *pgx.Conn
}
// CustomerPerevod - Перевести деньги другому клиенту
func(s *CustomerRepository) CustomerTransfer() {
	// var number string
	for {
		fmt.Print("Переводы")
		num:=utils.ReadString(types.MenuMoneyTransfer)
		// fmt.Scan(&number)
		switch num {
		case "1":
			//Перевод по номер счета
			s.CustomerTransferAccount()
			continue
		case "2":
			//перевод по номеру телефона
			s.PhoneTransaction()
			continue
		case "q":
			return
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}
	}
}
// CustomerPerevodAccount - перевод по номеру счета
func (s *CustomerRepository) CustomerTransferAccount() error {
	var payerAccountId, receiverAccountId int64
	accountService:=account.NewAccountServicce(s.connect)
	fmt.Println("Перевод по номеру счета")
	payerAccount:=utils.ReadString("введите номер счета для снятия денег: ")
	amount:=utils.ReadInt("Введите сумму: ")
	receiverAccount:=utils.ReadString("Введите номер счета получателя: ")
	err := s.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, payerAccount).Scan(&payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	err = s.connect.QueryRow(context.Background(), `select id from account where account_name = $1`, receiverAccount).Scan(&receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return 	accountService.TransferMoneyByAccountId(payerAccountId,receiverAccountId,amount)
}
// CustomerPerevodPhone - перевод по номеру телефона
func (s *CustomerRepository) PhoneTransaction() error {
	var payerAccountId,receiverAccountId int64
	accountService:=account.NewAccountServicce(s.connect)
	fmt.Println("Перевод по номеру телефона")
	payerPhone:=utils.ReadString("Input payerPhone: ")
	amount:=utils.ReadInt("Input amount: ")
	receiverPhone:=utils.ReadString("Input receiverPhone: ")
	ctx:=context.Background()
	selectSql:=`select account.id from account left join customer on customer.id=account.customer_id where customer.phone=$1`
	err:=s.connect.QueryRow(ctx,selectSql,payerPhone).Scan(&payerAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	err=s.connect.QueryRow(ctx,selectSql,receiverPhone).Scan(&receiverAccountId)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	return accountService.TransferMoneyByAccountId(payerAccountId,receiverAccountId,amount)
}
// PayService - Меню для оплата услуг
func (s *CustomerRepository) PayService() {
	// var number string
	for {
		fmt.Print("Оплатить услуги")
		fmt.Print(types.ServiceAdd)
		num:=utils.ReadString(types.MenuMoneyTransfer)
		// fmt.Scan(&number)
		switch num {
		case "1":
			//попплнение баланса телефон
			s.PayServicePhone()
			continue
		case "q":
			return
		default:
			fmt.Println("Выбрана неверная команда")
			continue
		}
	}
}
// PayServicePhone - попплнение баланса телефон
func (s *CustomerRepository) PayServicePhone() error {
	fmt.Println("услуга для пополнение номер телефона")
	var amuntaccount, amount int64
	var accountName, phone string
	accountName=utils.ReadString("Введите номер счета для снятия денег: ")
	amount=utils.ReadInt("Введите сумму: ")
	fmt.Print("Введите номер телефона: ")
	fmt.Scan(&phone)
	err := s.connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountName).Scan(&amuntaccount)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	if amount > amuntaccount {
		err = errors.New("Not enough amount on your balance")
		fmt.Println(err)
		return err
	}
	
	_, err = s.connect.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amuntaccount-amount, accountName)
	if err != nil {
		utils.ErrCheck(err)
		return err
	}
	fmt.Println("Успешно!!!")
	
	return nil
}
