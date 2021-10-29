package service

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"mybankcli/pkg/types"
	"mybankcli/pkg/utils"
)

//Сервис - описывает обслуживание клиентов.
type ManagerRepository struct {
	connect *pgx.Conn
}

//NewServer - функция-конструктор для создания нового сервера.
func NewManagerRepository(connect *pgx.Conn) *ManagerRepository {
	return &ManagerRepository{connect: connect}
}

//Регистрация клиента
func (s *ManagerRepository) Register(ctx context.Context, reg *types.Registration) (*types.Manager, error) {
	item := &types.Manager{}
	hash, err := utils.HashPassword(reg.Password)
	if err != nil {
		log.Print("Ошибка хеширование парола")
		return nil, err
	}
	err = s.connect.QueryRow(ctx, `INSERT INTO managers (name,surname,phone, password)
	VALUES ($1,$2,$3,$4) ON CONFLICT (phone) DO NOTHING RETURNING id,name,surname,phone,password,active, created`, reg.FirstName, reg.LastName, reg.Phone, hash).Scan(
		&item.ID, &item.Name, &item.SurName, &item.Phone, &item.Password, &item.Active, &item.Created)
	if err == pgx.ErrNoRows {
		return nil, err
	}
	if err != nil {
		log.Print("Не удалось регистрировать менеджера")
		return nil, err
	}
	return item, err
}

//   метод генерации токенов менеджеров
func (s *ManagerRepository) Token(ctx context.Context, phone string, password string) (token string, err error) {
	var hash string
	var id int64
	err = s.connect.QueryRow(ctx, `SELECT id,password FROM managers WHERE phone =$1`, phone).Scan(&id, &hash)
	if err == pgx.ErrNoRows {
		return "", err
	}
	if err != nil {
		return "", err
	}
	err = utils.CheckPasswordHass(password, hash)
	if err != nil {
		log.Print("неправильно логин или пароль")
		return "", err
	}
	token, _ = utils.HashPassword(password)
	_, err = s.connect.Exec(context.Background(), `INSERT INTO managers_tokens(token,manager_id) VALUES($1,$2)`, token, id)
	if err != nil {
		log.Print("Не удалось вставить токень менеджера в таблицу")
		return "", err
	}
	return token, err
}

//Список всех Менеджеров
func (s *ManagerRepository) ManagersAll(ctx context.Context) ([]*types.Manager, error) {
	managers := []*types.Manager{}
	rows, err := s.connect.Query(ctx, `SELECT *FROM managers`)
	if err != nil {
		log.Print("Не удалось вывести список менеджеров")
		return nil, err
	}
	// defer rows.Close()
	for rows.Next() {
		item := &types.Manager{}
		err = rows.Scan(&item.ID, &item.Name, &item.SurName, &item.Phone, &item.Password, &item.Active, &item.Created)
		if err != nil {
			log.Println(err)
		}
		managers = append(managers, item)
	}
	return managers, err
}

//Вывод всех активный Менеджеров
func (s *ManagerRepository) ManagersAllActive() ([]*types.Manager, error) {
	managers := []*types.Manager{}
	rows, err := s.connect.Query(context.Background(), `SELECT *FROM managers where active=true`)
	if err != nil {
		log.Print("Не возможно вывести активный менеджеров")
		return nil, err
	}
	// defer rows.Close()
	for rows.Next() {
		item := &types.Manager{}
		err = rows.Scan(&item.ID, &item.Name, &item.SurName, &item.Phone, &item.Password, &item.Active, &item.Created)
		if err != nil {
			log.Println("Не удалось сканировать ")
		}
		managers = append(managers, item)
	}
	return managers, err
}

//Список Менеджеров по Id
func (s *ManagerRepository) ManagersById(ctx context.Context, id int64) (*types.Manager, error) {
	managers := &types.Manager{}
	if id <= 0 {
		log.Print("id начинается с 1")
		return nil, pgx.ErrNoRows
	}
	err := s.connect.QueryRow(ctx, `select id,name,surname,phone,password,active,created from managers where id=$1`,
		id).Scan(&managers.ID, &managers.Name, &managers.SurName, &managers.Phone, &managers.Password, &managers.Active, &managers.Created)
	if err != nil {
		log.Println("Ошибка при выводе список менеджера")
		return nil, err
	}

	return managers, err
}

// Удалит Менеджера по их Id
func (s *ManagerRepository) ManagersRemoveByID(ctx context.Context, id int64) (*types.Manager, error) {
	managers := &types.Manager{}
	if id <= 0 {
		log.Print("id начинается с 1")
		return nil, pgx.ErrNoRows
	}
	err := s.connect.QueryRow(context.Background(), `DELETE FROM managers WHERE id = $1`,
		id).Scan(&managers.ID, &managers.Name, &managers.SurName, &managers.Phone, &managers.Password, &managers.Active, &managers.Created)
	if err != nil {
		log.Print("Не удалось удалить менеджера по ID")
		return nil, err
	}
	return managers, err
}

//Удаление токен менеджера по их Id
func (s *ManagerRepository) ManagersTokenRemoveByID(ctx context.Context, id int64) (*types.Tokens, error) {
	tokens := &types.Tokens{}
	err := s.connect.QueryRow(ctx, `delete from managers_tokens where manager_id=$1`,
		id).Scan(&tokens.Id)
	if err != nil {
		log.Print("не удалось удалить токена менеджера")
		return nil, err
	}
	return tokens, err
}

//Создание нового менеджера
func (s *ManagerRepository) CreateManagers(managers *types.Manager) (*types.Manager, error) {
	item := &types.Manager{}
	pass, err := utils.HashPassword(managers.Password)
	if err != nil {
		return nil, err
	}
	if managers.ID == 0 {
		log.Println("Вы ввели неверный номер пожалуйста введите номер с 1 ")
	} else {
		err := s.connect.QueryRow(context.Background(), `insert into managers(id,name,surname,phone,password) values($1,$2,$3,$4,$5) returning id,name,surname,phone,password,active,created`,
			managers.ID, managers.Name, managers.SurName, managers.Phone, pass).Scan(&item.ID, &item.Name, &item.SurName, &item.Phone, &item.Password, &item.Active, &item.Created)
		if err != nil {
			log.Print("Ошибка при создание менеджера")
			log.Print(err)
			return nil, err
		}
	}
	return item, err
}
