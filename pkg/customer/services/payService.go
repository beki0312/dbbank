package services

import (
	"context"
	"fmt"
	"mybankcli/pkg/types"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

// PayService - Меню для оплата услуг
func PayService(conn *pgx.Conn) {
	var number string
	for {
		fmt.Print("Оплатить услуги")
		fmt.Print(types.ServiceAdd)
		fmt.Scan(&number)
		switch number {
		case "1":
			//попплнение баланса телефон
			PayServicePhone(conn)
			continue
		case "2":
			//
			PayServiceLight(conn)
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
func PayServicePhone(connect *pgx.Conn) error {
	fmt.Println("услуга для пополнение номер телефона")
	var amuntaccount, amount int64
	var accountName, phone string
	fmt.Print("Введите номер счета для снятия денег: ")
	fmt.Scan(&accountName)
	fmt.Print("Введите сумму: ")
	fmt.Scan(&amount)
	fmt.Print("Введите номер телефона: ")
	fmt.Scan(&phone)

	err := connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountName).Scan(&amuntaccount)

	if err != nil {
		fmt.Printf("can't get Balance %e", err)
		return err
	}
	if amount > amuntaccount {
		err = errors.New("Not enough amount on your balance")
		fmt.Println(err)
		return err
	}
	tx, err := connect.Begin(context.Background())
	if err != nil {
		fmt.Printf("can't open transaction %e", err)
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.Background())

		}
		err = tx.Commit(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}()

	_, err = tx.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amuntaccount-amount, accountName)
	if err != nil {
		return err
	} else {
		fmt.Println("Перевод Успешно отправлено!!!")
		fmt.Println("")
	}
	
	return nil
}
// ???????????????????
func PayServiceLight(connect *pgx.Conn) error {
	fmt.Println("Перевод по номеру счета")
	var amuntName, amountCustomer, amount int64
	var accountName, accountCustomer string
	fmt.Print("Введите номер счета для снятия денег: ")
	fmt.Scan(&accountName)
	fmt.Print("Введите сумму: ")
	fmt.Scan(&amount)
	fmt.Print("Введите номер счета получателя: ")
	fmt.Scan(&accountCustomer)

	err := connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountName).Scan(&amuntName)
	err = connect.QueryRow(context.Background(), `select amount from account where account_name = $1`, accountCustomer).Scan(&amountCustomer)

	if err != nil {
		fmt.Printf("can't get Balance %e", err)
		return err
	}
	if amount > amuntName {
		err = errors.New("Not enough amount on your balance")
		fmt.Println(err)
		return err
	}
	tx, err := connect.Begin(context.Background())
	if err != nil {
		fmt.Printf("can't open transaction %e", err)
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.Background())

		}
		err = tx.Commit(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}()

	_, err = tx.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amuntName-amount, accountName)
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), `update account set amount = $1 where account_name = $2`, amountCustomer+amount, accountCustomer)
	if err != nil {
		return err
	} else {
		fmt.Println("Перевод Успешно отправлено!!!")
		fmt.Println("")
	}

	return nil
}