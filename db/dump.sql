INSERT INTO managers (name,surname,phone,password)
VALUES ('Бекмурод','Мустонов','917030100','033333');

INSERT INTO customer(name,surname,phone,password,amount)
VALUES ('Ab','Ba','917030303','0301',13)


SELECT customer.name,customer.phone,account.currency_code, account.account_name,account.amount 
FROM account 
JOIN customer ON account.customer_phone = customer.id
where account.customer_phone=customer.id;


SELECT *FROM atm;

select amount from account where  id=1;

SELECT account.currency_code, account.account_name,account.amount FROM account 
	JOIN customer ON account.customer_id = customer.id
    where customer.phone='915030102';

SELECT id,customer_id,currency_code,account_name,amount from account
-- where customer.phoe ='915030102'


INSERT INTO account(amount) VALUES ('Сотовые операторы')

select *from customer where phone='915030102' and password='0301'

 select account.id, amount 
        from account
        left join customer on customer.id = account.customer_id
        where account.is_main = true and customer.phone = :раками_телефон


        update account a set amount = '100' 
	from customer c 
	where c.id=a.customer_id and c.phone='917030101'

    update account set amount='12000' where account_name='1122331'
    

     select account.id,customer_id,currency_code,account_name,amount 
        from account
        left join customer on customer.id = account.customer_id
        where customer.phone='915030102'