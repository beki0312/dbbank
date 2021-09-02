package customer

import "time"

type Customer struct {
	ID       		int64  			`json:"id"`
	Name     		string 			`json:"name"`
	SurName  		string 			`json:"surname"`
	Phone    		string 			`json:"phone"`
	Active   		bool   			`json:"active"`
	Created  		time.Time 		`json:"created"`
}