package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//для ввода строки
func ReadString(lable string) string {
	var input string
	fmt.Print(lable)
	fmt.Scan(&input)
	return input
}
//Для ввода число
func ReadInt(lable string) int64  {
	var input int64
	fmt.Print(lable)
	fmt.Scan(&input)
	return input
}
// ошибка
func ErrCheck(err error)  {
	if err != nil {
		fmt.Print(err)
		return
	}
}

//Хеш для парола
func HashPassword(password string) (string,error)  {
	bytes,err:=bcrypt.GenerateFromPassword([]byte(password),14)
	return string(bytes),err
}

func CheckPasswordHass(password,hash string) error  {
	err:=bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err
}