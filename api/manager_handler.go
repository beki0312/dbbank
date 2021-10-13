package api

import (
	"context"
	"log"
	"mybankcli/pkg/manager/service"
	"mybankcli/pkg/types"

	"github.com/jackc/pgx/v4"
)
type ManagerHandler struct {
	connect 			*pgx.Conn
	managerRepository 	*service.ManagerRepository
}

func NewManagerHandler(connect *pgx.Conn,managerRepository *service.ManagerRepository) *ManagerHandler {
	return &ManagerHandler{connect:connect,managerRepository: managerRepository}
}

func (h *ManagerHandler) RegistersManagers(ctx context.Context,item *types.Registration) (*types.Manager, error) {
	manager:=&types.Manager{}
	registration,err:=h.managerRepository.Register(item)
	if err != nil {
		return nil, err
	}
	if registration==nil {
		return nil,ErrNotFound
	}
	return manager,err
}

func (h *ManagerHandler) GetManagersToken(ctx context.Context, item *types.Authers) (token string, err error) {
	token,err=h.managerRepository.Token(item.Phone,item.Password)
	if err != nil {
		return "", err
	}
	return	token,err
}
//find Id customers Token
func (s *ManagerHandler) IDByTokenManagers(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.connect.QueryRow(ctx,`SELECT manager_id FROM managers_tokens WHERE token =$1`, token).Scan(&id)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id,err 
}
//Get All Manager
func (h *ManagerHandler) GetManagersAll(ctx context.Context) ( []*types.Manager,error) {
	managers,err:=h.managerRepository.ManagersAll()
	if err != nil {
		return nil,ErrInternal
	}
	
	return managers,nil
}
//Get All Active Manager
func (h *ManagerHandler) GetManagersAllActive(ctx context.Context) ( []*types.Manager,error) {
	managers,err:=h.managerRepository.ManagersAllActive()
	if err != nil {
		return nil, ErrInternal
	}
	return managers,nil
}
//Get ById Managers
func (h *ManagerHandler) GetManagersById(ctx context.Context,id int64) (*types.Manager,error) {
	if (id<=0) {
		return nil, ErrInternal
	}
	managers,err:=h.managerRepository.ManagersById(id)
	if err != nil {
		return nil,ErrInternal
	}
	if managers==nil {
		return nil,ErrNotFound
	}
	return managers,nil
}
// Delete Manager 
func (h *ManagerHandler) GetManagersRemoveByID(ctx context.Context, id int64) (*types.Manager, error) { 
	if (id<=0) {
		return nil,ErrInternal	
	}
	managers,err:=h.managerRepository.ManagersRemoveByID(id)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	if managers==nil {
		return nil,ErrNotFound
	}
	return managers, nil
}
//Save Manager by id
func (h *ManagerHandler) PostManagers(ctx context.Context, managers *types.Manager) (*types.Manager,error) {
	if (managers.ID<=0){
		return nil,ErrInternal
	}	
	managers,err:=h.managerRepository.CreateManagers(managers)
	if err != nil {
		return nil, ErrInternal
	}
	if managers==nil{
		return nil,ErrNotFound
	}
	return managers, err
}
