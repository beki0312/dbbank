package api


import (
	"context"
	"log"
	"mybankcli/pkg/types"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)
type ManagerHandler struct {
	connect *pgx.Conn
}
func NewManagerHandler(connect *pgx.Conn) *ManagerHandler {
	return &ManagerHandler{connect: connect}
}
//Get All Manager
func (s *ManagerHandler) ManagersAll(ctx context.Context) ( []*types.Manager,error) {
	managers:=[]*types.Manager{}
	rows,err:=s.connect.Query(ctx,`SELECT *FROM managers`)
	if err != nil {
		return nil, ErrInternal
	}
	// defer rows.Close()
	for rows.Next(){
		item:=&types.Manager{}
		err=rows.Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
		if err != nil {
			log.Println(err)
		}
		managers = append(managers, item)
	}
	return managers,nil
}
//Get All Active Manager
func (s *ManagerHandler) ManagersAllActive(ctx context.Context) ( []*types.Manager,error) {
	managers:=[]*types.Manager{}
	rows,err:=s.connect.Query(ctx,`SELECT *FROM managers where active=true`)
	if err != nil {
		return nil, ErrInternal
	}
	// defer rows.Close()
	for rows.Next(){
		item:=&types.Manager{}
		err=rows.Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)
		if err != nil {
			log.Println(err)
		}
		managers = append(managers, item)
	}
	return managers,nil
}
//Get ById Managers
func (s *ManagerHandler) ManagersById(ctx context.Context,id int64) (*types.Manager,error) {
	managers:=&types.Manager{}
	err:=s.connect.QueryRow(ctx,`select id,name,surname,phone,password,active,created from managers where id=$1`,
	id).Scan(&managers.ID,&managers.Name,&managers.SurName,&managers.Phone,&managers.Password,&managers.Active,&managers.Created)
	if err != nil {
		log.Println(err)
		return nil,ErrInternal
	}
	return managers,nil
}
// Delete Manager by id
func (s *ManagerHandler) ManagersRemoveByID(ctx context.Context, id int64) (*types.Manager, error) {
	managers := &types.Manager{}
	err := s.connect.QueryRow(ctx, `DELETE FROM managers WHERE id = $1`, 
	id).Scan(&managers.ID, &managers.Name, &managers.SurName,&managers.Phone,&managers.Password,&managers.Active, &managers.Created)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return managers, nil
}
//Save Manager by id
func (s *ManagerHandler) CreateManagers(ctx context.Context, managers *types.Manager) (*types.Manager,error) {
	item:=&types.Manager{}
	pass,_:=bcrypt.GenerateFromPassword([]byte(item.Password),14)
	if managers.ID==0 {
		log.Println("Вы ввели неверный номер пожалуйста введите номер с 1 ")
	}else{
		err:=s.connect.QueryRow(ctx,`insert into managers(id,name,surname,phone,password) values($1,$2,$3,$4,$5) returning id,name,surname,phone,password,active,created`,
		managers.ID,managers.Name,managers.SurName,managers.Phone,pass).Scan(&item.ID,&item.Name,&item.SurName,&item.Phone,&item.Password,&item.Active,&item.Created)	
		if err != nil {
			log.Print(err)
			return nil,ErrInternal
		}
	}
	return item,nil
}

//Block and Unblock Managers by his id
func (s *ManagerHandler) ManagersBlockAndUnblockById(ctx context.Context, id int64,active bool) (*types.Manager,error) {
	managers:=&types.Manager{}
	err:=s.connect.QueryRow(ctx,`update managers set active =$1 where id=$2`,active,id).Scan(
		&managers.ID,&managers.Name,&managers.SurName,&managers.Phone,&managers.Password,&managers.Active,&managers.Created)
		if err != nil {
			log.Println(err)
			return nil, ErrInternal
		}
		return managers,nil
}