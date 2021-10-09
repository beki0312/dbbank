package api
// type CustomerHandler struct {
// 	customerService *CustomerService
// }

// //Get All Customer
// func (h *CustomerHandler) GetAllCustomers(ctx context.Context) ( []*types.Customer,error) {
//   customers, err := h.customerService.GetAllCustomers()
// 	return customers, err
// }
// //Get ById customer
// func (h *CustomerHandler) GetCustomerById(ctx context.Context, id int64) (*types.Customer, error) {
// 	// validation
//   if (id <= 0) {
//     return nil, ErrBadRequest
//   }
//   customer, err := h.customerService.GetCustomerById(id)
//   if err != nil {
//     log.error("error while getting customer by id", id, err)
//     return nil, ErrInternal
//   }
//   if customer == nil {
//   	return nil, ErrNotFound
//   }
// 	return customers, err
// }
// //Get ById customer
// func (s *CustomerHandler) CustomerById(ctx context.Context,id int64) (*types.Customer,error) {
// 	customers:=&types.Customer{}
// 	err:=s.connect.QueryRow(ctx,`select id,name,surname,phone,password,active,created from customer where id=$1`,
// 	id).Scan(&customers.ID,&customers.Name,&customers.SurName,&customers.Phone,&customers.Password,&customers.Active,&customers.Created)
// 	if err != nil {
// 		log.Println(err)
// 		return nil,ErrInternal
// 	}
// 	return customers,nil
// }
// 	passHash , _:=bcrypt.GenerateFromPassword([]byte(password),14)
// 	err := s.customerRepository.connect.QueryRow(ctx, `select id from customer where phone=$1 and password = $2`, phone, passHash).
// func (s *CustomerService) CustomerAccount(phone string) error{
// 	var password string
// 	// var pass []byte
// 	phone=utils.ReadString("Введите Лог: ")
// 	password=utils.ReadString("Введите парол: ")
// 	passBytes,_:=bcrypt.GenerateFromPassword([]byte(password),14)
//   passHash := string(passBytes)
// 	cust:=types.Customer{}
// 	ctx := context.Background()
// 	err := s.customerRepository.connect.QueryRow(ctx, `select password from customer where phone=$1`,phone).
// 	Scan(&cust.ID,&cust.Name,&cust.SurName,&cust.Phone,&cust.Password,&cust.Active,&cust.Created)
// 	if err != nil {
// 		utils.ErrCheck(err)
// 		return err
// 	}
// 	if passHash == cust.Password{
// 		fmt.Println("Хуш омадед Мизоч!!!")
// 		println("")
// 	} else {
// 		fmt.Println("Шумо логин ё паролро нодуруст дохил намудед!!!")
// 		fmt.Println(err)
// 		return err
// 	}
// 	s.ServiceLoop(phone)





