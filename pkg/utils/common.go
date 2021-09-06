package utils

import "fmt"

// 
func ReadString(lable string) string {
	var input string

	fmt.Print(lable)
	fmt.Scan(&input)

	return input
}

// phone:= ReadString("input phone:")