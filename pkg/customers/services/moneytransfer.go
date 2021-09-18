package services
import (
	"context"
	"fmt"
	"mybankcli/pkg/account"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"github.com/jackc/pgx/v4"
)
type MoneyService struct {
	connect *pgx.Conn
}
func NewMoneyServicce(connect *pgx.Conn) *MoneyService{
	return &MoneyService{connect: connect}
}
// CustomerPerevod - Перевести деньги другому клиенту
func(s *MoneyService) CustomerPerevod() {
	// var number string
	for {
		fmt.Print("Переводы")
		num:=utils.ReadString(types.MenuMoneyTransfer)
		// fmt.Scan(&number)
		switch num {
		case "1":
			//Перевод по номер счета
			s.CustomerPerevodAccount()
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
func (s *MoneyService)  CustomerPerevodAccount() error {
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
	s.Transactions(payerAccountId,receiverAccountId,amount)
	return 	accountService.TransferMoneyByAccountId(payerAccountId,receiverAccountId,amount)
}

// CustomerPerevodPhone - перевод по номеру телефона
func (s *MoneyService) PhoneTransaction() error {
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
	s.Transactions(payerAccountId,receiverAccountId,amount)
	return accountService.TransferMoneyByAccountId(payerAccountId,receiverAccountId,amount)
}
//Таблица транзаксия
func (s *MoneyService) Transactions(payerAccountId,receiverAccountId,amount int64)  {
	ctx:=context.Background()
	item:=types.Transactions{}
	err:=s.connect.QueryRow(ctx, `insert into transactions (debet_account_id,credit_account_id,amount) values ($1,$2,$3) returning id,debet_account_id,credit_account_id,amount,date 
	`,payerAccountId,receiverAccountId,amount).Scan(&item.ID,&item.Debet_account_id,&item.Credit_account_id,&item.Amount,&item.Date)
	utils.ErrCheck(err)
}
