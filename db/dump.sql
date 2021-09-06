INSERT INTO managers (name,surname,phone,password)
VALUES ('Бекмурод','Мустонов','917030100','033333');

INSERT INTO customer(name,surname,phone,password,amount)
VALUES ('Ab','Ba','917030303','0301',13)


SELECT customer.name,customer.phone,account.currency_code, account.account_name,account.amount 
FROM account 
JOIN customer ON account.customer_phone = customer.id
where account.customer_phone=customer.id;


SELECT *FROM atm;


SELECT account.currency_code, account.account_name,account.amount FROM account 
	JOIN customer ON account.customer_id = customer.id
    where customer.phone='915030102';

SELECT id,customer_id,currency_code,account_name,amount from account
-- where customer.phoe ='915030102'




select *from customer where phone='915030102' and password='0301'

