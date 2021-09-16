package services
import (
	"context"
	"fmt"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
	"github.com/pkg/errors"
)

// PayService - Меню для оплата услуг
func (s *MoneyService) PayService() {
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
func (s *MoneyService) PayServicePhone() error {
	fmt.Println("услуга для пополнение номер телефона")
	var amuntaccount, amount int64
	var accountName, phone string
	accountName=utils.ReadString("Введите номер счета для снятия денег: ")
	amount=utils.ReadInt("Введите сумму: ")
	fmt.Print("Введите номер телефона: ")
	fmt.Scan(&phone)
	err := s.connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountName).Scan(&amuntaccount)
	if err != nil {
		fmt.Printf("can't get Balance %e", err)
		return err
	}
	if amount > amuntaccount {
		err = errors.New("Not enough amount on your balance")
		fmt.Println(err)
		return err
	}
	
	_, err = s.connect.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amuntaccount-amount, accountName)
	if err != nil {
		return err
	} else {
		fmt.Println("Успешно!!!")
	}
	return nil
}
