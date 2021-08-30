package services
import (
	"context"
	"fmt"
	"mybankcli/pkg/types"
	"os"
	"github.com/jackc/pgx/v4"
)


func CustomerAccount(connect *pgx.Conn) error{
	var phone, password, pass string
	fmt.Print("Введите Лог: ")
	fmt.Scan(&phone)
	fmt.Print("Введите парол: ")
	fmt.Scan(&password)
	println("")
	ctx := context.Background()
	err := connect.QueryRow(ctx, `select password from customer where phone=$1`, phone).Scan(&pass)
	if err != nil {
		fmt.Printf("can't get password Customer %e", err)
		return err
	}

	if password == pass {
		fmt.Println("Хуш омадед Мизоч!!!")
		println("")
	} else {
		fmt.Println("Шумо паролро нодуруст дохил намудед!!!")
		fmt.Println(err)
		return err
	}
	Loop(connect)
	return nil
}

func Loop(con *pgx.Conn) {
	var cmd string
	for {
		fmt.Println(types.MenuCustomer)
		fmt.Scan(&cmd)
		switch cmd {
		case "1":
			//TODO: Добавить пользователя
			// ManagerAddCustomer(con)
			continue
		case "2":
			//TODO: Добавить счет
			// ManagerAddAccount(con)
			continue
		case "3":
			//TODO: Добавить услугу
			// ManagerAddServices(con)
			continue
		case "10":
			//TODO: Добавить Банкоматов
			// ManagerAddAtm(con)
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			return
		}
	}
}