# Requirements

Necessary programs to run the api

```
- docker
- go
```

# Database

Set up container running mysql server

```
 » make create-mysql-container
 ```

Wait a few seconds for container to start. Expect this message

```
   docker exec lemon_mysql mysql -uroot -pyour_secret -e "create database lemon";
   mysql: [Warning] Using a password on the command line interface can be insecure.
```

Create database

```
» make database/create   
``` 

Run migrations and fill with dump data

```
» make database/migrations/up         
» make database/fill
```

https://github.com/golang-migrate/migrate was used to manage migrations.

# Observations

- Every deposit, extraction or transfer's amount will be interpreted in the currency `cents`.

Example: a deposit of `1 ars` is interpreted as 1 cent of an `ars`. The associated transaction will show the amount
as `0,01`. Same with other currencies

# Run

```
» make run
```

**Or**

````
» make build 
» bin/lemon serve --loglevel=debug
````

**Optional parameters:**

- `--loglevel`: sets the log level. Defaults to `debug`
- `--database`: database name. Defaults to `lemon`
- `--sql_user`: mysql user. Defaults to `root`
- `--sql_password`: mysql user's password. Defaults to `your_secret`
- `--sql_max_open_conns`: max open connections to db. Defaults to `10`

# Testing

```
» make test
```

# Usage

There is a `postman` folder where you can import a set of pre-made requests. A postman environment is also provided. The
request are prepared to work with the environment given. The environment is loaded with the users and wallets created by
the `dump.sql` file (used by the command `make database/fill` )

Also, the following requests should work if the command `make database/fill` was executed

**Create user**

```
curl --location --request POST 'http://localhost:8080/api/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "firstname": "Juan", 
    "lastname": "Noli",
    "alias": "somealsias",
    "email": "juannsoli@gmail.com"
}'
```

**Get user**

```
curl --location --request GET 'http://localhost:8080/api/users/2c909e4e-88c5-4cfb-8a4e-b41bdbacd55c'
```

**Deposit**

```
curl --location --request POST 'http://localhost:8080/api/users/2c909e4e-88c5-4cfb-8a4e-b41bdbacd55c/wallets/9edbc638-4c9f-4b3a-9b54-25299206993d/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount": 100000000
}'
```

**Extraction**

```
curl --location --request POST 'http://localhost:8080/api/users/2c909e4e-88c5-4cfb-8a4e-b41bdbacd55c/wallets/9edbc638-4c9f-4b3a-9b54-25299206993d/extraction' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount": 100
}'
```

**User movements**

```
curl --location --request GET 'http://localhost:8080/api/users/2c909e4e-88c5-4cfb-8a4e-b41bdbacd55c/wallets/9edbc638-4c9f-4b3a-9b54-25299206993d/transactions?offset=0&limit=2&transaction_type=deposit&currency=USDT'
```

# Security

Database password can be changed. Search for `DB_PASSWORD` in makefile to change it.

# Future work

- Make server run in docker (only mysql server runs in docker)
- User and wallets creation should be done in a transaction. The current solution is not so bad though, wallets can be
  added later
- Add more coverage (only user's handler and user's repository are tested to give and example of how I would test the
  other files)
- Improve errors (add more custom errors, be more precise with errors)
- Paginate list users
- Extraction and deposit responses could be improved 