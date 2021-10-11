package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)


func main() {
	s := "Password"
	p, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bcrypt.CompareHashAndPassword(p, []byte("Password")))
	fmt.Println(bcrypt.CompareHashAndPassword(p, []byte("Passwordadfadf")))
	fmt.Println(bcrypt.ErrMismatchedHashAndPassword)
}