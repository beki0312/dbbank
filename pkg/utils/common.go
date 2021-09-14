package utils

import (
	"fmt"
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