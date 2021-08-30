package service

import (
	"context"
	"fmt"
	"mybankcli/pkg/customer/services"
	"mybankcli/pkg/types"
	"os"

	"github.com/jackc/pgx/v4"
)
func Auther(conn *pgx.Conn)  {
	var numberauther string
	
	for{
		fmt.Println(types.Auther)
		fmt.Scan(&numberauther)
		switch numberauther{
		case "1":
			ManagerAccount(conn)
			continue
		case "2":
			services.CustomerAccount(conn)
			continue
		case "q":
			os.Exit(0)
		}
	}
	
}
func ManagerAccount(connect *pgx.Conn) error {
	var phone,password, pass string 
	fmt.Print("Введите Лог: ")
	fmt.Scan(&phone)
	fmt.Print("Введите парол: ")
	fmt.Scan(&password)
	println("")
	ctx:=context.Background()
	err:=connect.QueryRow(ctx, `select password from managers where phone=$1`,phone).Scan(&pass)
	if err != nil {
		fmt.Printf("can't get password %e",err)
		return err
	}
	if password ==pass{
		fmt.Println("Хуш омадед Менедчер")
		println("")
	}else{
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
		fmt.Println(types.MenuManager)
		fmt.Scan(&cmd)
		switch cmd {
		case "1":
			//TODO: Добавить пользователя
			ManagerAddCustomer(con)
			continue
		case "2":
			//TODO: Добавить счет
			ManagerAddAccount(con)
			continue
		case "3":
			//TODO: Добавить услугу
			ManagerAddServices(con)
			continue
		case "10":
			//TODO: Добавить Банкоматов
			ManagerAddAtm(con)
			continue
		case "q":
			os.Exit(0)
		default:
			fmt.Println("Выбрана неверная команда")
			return
		}
	}
}

func ManagerAddCustomer(connect *pgx.Conn,)  {
	var name,surname,phone,password string 
	var amount int64
			fmt.Print("Введите Имя: ")
			fmt.Scan(&name)
			fmt.Print("Введите Фамилия: ")
			fmt.Scan(&surname)
			fmt.Print("Введите Лог: ")
			fmt.Scan(&phone)
			fmt.Print("Введите парол: ")
			fmt.Scan(&password)
			
			fmt.Print("Введите Балансе: ")
			fmt.Scan(&amount)
			println("")
	fmt.Println("Добалили клиент: Имя ",name, " фамиля ",surname," Логин ",phone," Парол ",password," Балансе",amount)
	println("")
	ctx:=context.Background()
	item:=types.Client{}
	err:=connect.QueryRow(ctx, `insert into customer (name,surname,phone,password,amount)
	values ($1,$2,$3,$4,$5) returning id,name,surname,phone,password,amount,active,created 
	`,name,surname,phone,password,amount).Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Amount,&item.Active,&item.Created)
	if err != nil {
		fmt.Printf("can't insert %e",err)
		// return 
	}
}
func ManagerAddAccount(connect *pgx.Conn)  {
	fmt.Println("Добавить Счеты ")
	var account, amount int64 
			// fmt.Print("Введите id: ")
			// fmt.Scan(&id)
			fmt.Print("Введите Счет: ")			
			fmt.Scan(&account)
			fmt.Print("Введите Баланс: ")			
			fmt.Scan(&amount)
			println("")
	fmt.Println("Добалили Счет и баланс: ",account)
	println("")
	ctx:=context.Background()
	item:=types.Account{}
	err:=connect.QueryRow(ctx, `insert into account (account,amount) values ($1,$2) returning id,account,amount 
	`,account,amount).Scan(&item.ID,&item.Account,&item.Amount)
	if err != nil {
		fmt.Printf("can't insert %e",err)
		// return 
	}
}
func ManagerAddServices(connect *pgx.Conn)  {
	fmt.Println("Добавить услуги ")
	var name string 
			fmt.Print("Введите название услуги: ")			
			fmt.Scan(&name)
			println("")
	fmt.Println("Добавили услуги : ",name)
	println("")
	ctx:=context.Background()
	item:=types.Services{}
	err:=connect.QueryRow(ctx, `insert into services (name) values ($1) returning id,name 
	`,name).Scan(&item.ID,&item.Name)
	if err != nil {
		fmt.Printf("can't insert %e",err)
		// return 
	}
}
func ManagerAddAtm(connect *pgx.Conn,)  {
	var numbers int64
	var district, address string 
			fmt.Print("Введите № Банкомата: ")
			fmt.Scan(&numbers)
			fmt.Print("ВВедите район: ")
			fmt.Scan(&district)
			fmt.Print("Введите адрес Банкомата: ")
			fmt.Scan(&address)
			println("")
	fmt.Println("Добалили список Банкомат:  № ",numbers,", Район: ",district,", Адресс: ",address)
	println("")
	ctx:=context.Background()
	item:=types.Atm{}
	err:=connect.QueryRow(ctx, `insert into atm (numbers,district,address)
	values ($1,$2,$3) returning id,numbers,district,address,active,created 
	`,numbers,district,address).Scan(&item.ID,&item.Numbers,&item.District,&item.Address,&item.Active,&item.Created)
	if err != nil {
		fmt.Printf("can't insert %e",err)
	}
}
