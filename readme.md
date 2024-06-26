# __I. Работа с базой данных__

Проект изначально создавался под Postgres, но переделан и полностью рабочий под MySQL. Если есть отличия Postgres от MySQL по тексту добавлены коментарии. Поэтому видим Postgres, а подразумеваем MySQL.

## 1. Cхема БД и SQL-код

[diagram.io](https://dbdiagram.io/home)

    ->  Go to App
        Export PostgreSQL

doc/db.dbml
doc/schema.sql

<br>
<br>

## 2. Docker  Postgres

    $ docker ps    // список всех запущенных контейнеров
    $ docker images   // список всех имеющихся образов

__Удалить установленный Postgres__

    Uninstall the PostgreSQL application
    $ sudo apt-get --purge remove postgresql
    Remove PostgreSQL packages
    $ dpkg -l | grep postgres
    To uninstall PostgreSQL completely, you need to remove all of these packages using the following command:
    $ sudo apt-get --purge remove <package_name>
    Remove PostgreSQL directories
    $ sudo rm -rf /var/lib/postgresql/
    $ sudo rm -rf /var/log/postgresql/
    $ sudo rm -rf /etc/postgresql/
    Remove the postgres user
    $ sudo deluser postgres
    Verify uninstallation
    $ psql --version

__Скачать образ__

[hub.docker.com](https://hub.docker.com/)   
поиск postgres   
<https://hub.docker.com/_/postgres>

    docker pull <image>:<tag>
    $ docker pull postgres:16-alpine
    или
    $ docker pull mysql:8.0
<br>

__Запуск контейнера из образа__

    docker run --name <container_name> -e <environment_variable> -d <image>:<tag>

___Environment Variables:___

* POSTGRES_PASSWORD   
* POSTGRES_USER   
* POSTGRES_DB   
* POSTGRES_INITDB_ARGS   
* POSTGRES_INITDB_WALDIR   
* POSTGRES_HOST_AUTH_METHOD   
* PGDATA   

For example:
```
$ docker run -d \
    --name some-postgres \
    -e POSTGRES_PASSWORD=mysecretpassword \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v /custom/mount:/var/lib/postgresql/data \
    postgres
```

пароли можно подгружать из файла:   

    $ docker run --name some-postgres -e POSTGRES_PASSWORD_FILE=/run/secrets/postgres-passwd -d postgres

Port mapping    

    docker run --name ‹container_name> -e ‹environment_variable> -p ‹host_ports:container_ports> -d ‹image>:<tag>

    $ docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
    или
    $ docker run --name mysql8 -p 3306:3306 -e MYSQL_DATABASE=main_db -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0

    $ docker ps

    $ docker stop postgres16
    $ docker ps -a   // все контейнеры вне зависимости запущены или нет
    $ docker start postgres16   // снова запустить имеющийся контейнер
    $ docker rm postgres16   // удалить полностью имеющийся контейнер
<br>

__Запуск команды в контейнере__

    docker exec -it ‹container _name_or_id> ‹command> [args]

    $ docker exec -it postgres16 psql -U root
        select now();
        \q    - выход

    $ docker exec -it postgres16 /bin/sh     // запускаем оболочку в контейнере
<br>

__Просмотр логов контейнера__

    docker logs <container_name_or_id>
    $ docker logs postgres16
<br>
<br>

---------

__TablePlus__

Для kubuntu лучше либо pgAdmin4 либо DBeaver
При повторном подключении DBeaver к базе может возникнуть ошибка Public Key Retrieval is not allowed
В подключении к базе советуют (при разработке) в Свойствах драйвера изменить переменные allowPublicKeyRetrieval=true & useSSL=false

Для mac - [tableplus.com](https://tableplus.com/)  

    basename: root
    user: root
    password: secret
    url: localhost:5432
    для mysql url: 127.0.0.1:3306
<br>
<br>

---------   
## 3. Миграции

[github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)

    $ brew install golang-migrate
    $ migrate -version
        v4.17.0
    $ migrate -help

    $ migrate create -ext sql -dir db/migration -seq init schema

Открыть созданный файл миграции db/migration/000001_init.up.sql и скопировать содержимое вашего файла схемы базы данных doc/schema_mysql.sql в этот файл

    $ docker exec -it postgres16 /bin/sh
        # createdb -username=root -owner=root main_db
        # psql main_db
        # dropdb main_db
        # exit

    $ docker exec postgres16 createdb --username=root --owner=root main_db
    $ docker exec -it postgres16 psql -U root main_db
        \q

__Миграция в проекте__

    $ set -e
    $ source ./app.env
    // если нет базы  $ docker exec postgres16 createdb --username=root --owner=root main_db
    $ migrate -path ./db/migration -database "$DB_SOURCE" -verbose up
    или для mysql:
    $ migrate -path ./db/migration -database "mysql://root:secret@tcp(localhost:3306)/main_db" -verbose up
    $ migrate -path ./db/migration -database "mysql://root:secret@tcp(localhost:3306)/main_db" force 1
    $ migrate -path ./db/migration -database "mysql://root:secret@tcp(localhost:3306)/main_db" -verbose down

<br>

__Cоздаем Makefile__

    run-postgres: ## Start postgresql database docker image.
        docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

    start-postgres16: ## Start available postgresql database docker container.
        docker start postgres16

    stop-postgres: ## Stop postgresql database docker image.
        docker stop postgres16
<br>


---------
## 4. CRUD

* Create
* Read
* Update
* Delete

> __DATABASE/SQL__
> * Очень быстро и просто
> * Ручное сопоставление полей SQL с переменными 
> * Легко допускать ошибки, которые не обнаруживаются до выполнения

> __GORM__
> * Функции CRUD уже реализованы, очень короткий рабочий код
> * Необходимо научиться писать запросы с использованием функции gorm
> * Выполняется медленно при высокой нагрузке   
> * В 3 - 5  раз медленнее работает

> __SQLX__
> * Довольно быстрый и простой в использовании
> * Сопоставление полей с помощью тегов текста запроса и структуры
> * Сбой не произойдет до времени выполнения

> __SQLC__
> * Очень быстрый и простой в использовании
> * Автоматическая генерация кода
> * Отслеживание ошибок запроса SQL перед генерацией кодов
> * Полная поддержка Postgres. MySQL является экспериментальным

[sqlc.dev](https://sqlc.dev/)   
[github.com/sqlc-dev/sqlc](https://github.com/sqlc-dev/sqlc)

    $ brew install sqlc
    $ sqlc version
    $ sqlc help
    $ sqlc init

sqlc.yaml   
[docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#setting-up](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#setting-up)

db/query/account.sql   
[docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#schema-and-queries](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#schema-and-queries)   

__Отличия запросов mysql от postgres:__   
в запросах не может быть RETURNING *, выкручиваемся последующим запросом данных,   
в запросах меняем аргументы $1, $2, $3... на знаки "?", иначе не будет входных аргументов

    $ sqlc generate

<br>
<br>


## 5. Тесты

    _ "github.com/lib/pq"  // без драйвера работать не будет
    Для MySql драйвер mysql:
    $ go get -u github.com/go-sql-driver/mysql
    import (
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
    )

    $ go test -v   // все тесты

    $ go test -timeout 30s ./db/sqlc -run ^TestMain$             
        ok  	github.com/Country7/backend-captaincode-mysql/db/sqlc	0.433s [no tests to run]

    $ make test   // команда test из файла Makefile
<br>
<br>


## 6. Транзакции

Перевод 10 USD из банка аккаунта 1 в банк аккаунта 2:

1. Создайте запись транзакции о переводе с суммой = 10
2. Создайте учетную запись для учетной записи 1 с суммой = -10
3. Создайте учетную запись для учетной записи 2 с суммой = +10
4. Вычтите 10 из баланса учетной записи 1
5. Добавьте 10 к балансу учетной записи 2

<br>

> *    BEGIN
> *    ...
> *    COMMIT
or
> *    BEGIN
> *    ...
> *    ROLLBACK

<br>

```go
    type Store interface {
        Querier
        TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
    }

    // SQLStore provides all functions to execute SQL queries and transactions.
    type SQLStore struct {
        *Queries
        db *sql.DB
    }

    func NewStore(db *sql.DB) Store {
        return &SQLStore{db: db, Queries: New(db)}
    }

    // execTx executes a function within a database transaction.
    func (s *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
        tx, err := s.db.BeginTx(ctx, nil)
        if err != nil {
            return err
        }

        q := New(tx)
        err = fn(q)
        if err != nil {
            if rbErr := tx.Rollback(); rbErr != nil {
                return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
            }
            return err
        }

        return tx.Commit()
    }
```
<br>
<br>

## 7. Блокировка транзакции

    BEGIN;
    
    SELECT * FROM accounts WHERE id = 1;

    SELECT * FROM WHERE id = 1 FOR UPDATE;     // блокировка запросов
    UPDATE accounts SET balance = 500 WHERE id = 1;
    COMMIT;
<br>

    $ sqlc generate

__Deadlock detected__

    INSERT INTO entries (account_id, amount) VALUES ($1, $2) RETURNING *;
и
    
    SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR UPDATE;
заблокируют друг друга (Deadlock detected) несмотря на то, что обращение идет к разным таблицам

Эти две таблицы имеют связи FOREIGN KEY:   
ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");   
при обращении к accounts происходит обновление ключа id в таблице accounts по связям с entries   
чтобы этого не происходило необходима команда !!! NO KEY UPDATE

    SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

<br><br>


## 8. Взаимоблокировки

```sql
    BEGIN:
    UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
    UPDATE accounts SET balance = balance + 10 WHERE id = 2 RETURNING *;
    ROLLBACK;

    BEGIN:
    UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;
    UPDATE accounts SET balance = balance + 10 WHERE id = 1 RETURNING *;
    ROLLBACK;
```

Одновременно две эти транзакции приведут в взаимоблокировке (Deadlock detected)

```go
    if arg.FromAccountID < arg.ToAccountID {
        result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
    } else {
        result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
    }
```
<br>
<br>


## 9. Уровень изоляции транзакций


1. Чтение незафиксированных транзакций (read uncommitted)
2. Чтение зафиксированных данных (read committed)
3. Повторяемый уровень изоляции чтения (repeatable read)
4. Параллельные разрешения (serializable)

```shell
    mysql> select @@transaction_isolation;
    mysql> select @@global.transaction_isolation;

    mysql> set session transaction isolation level read uncommitted;
    mysql> set session transaction isolation level read committed;
    mysql> set session transaction isolation level repeatable read;
    mysql> set session transaction isolation level serializable;
```

```shell
    postgres=# show transaction isolation level;

    postgres=# begin;
    postgres=# set transaction isolation level read uncommitted;
    # в postgres уровень read uncommitted ведет себя как read committed, как будто его нет
    postgres=# set transaction isolation level read committed;
    postgres=# set transaction isolation level repeatable read;
    postgres=# set transaction isolation level serializable;
    postgres=# show transaction isolation level;
    postgres=# commit;
```

|                       |READ UNCOMMITTED   | READ COMMITTED    | REPEATABLE READ   | SERIALIZABLE  |
|:-:                    |:-:                |:-:                |:-:                |:-:            |
| DIRTY READ            | V                 |         -         |         -         |       -       |
| NON-REPEATABLE READ   | V                 | V                 |         -         |       -       |
| PHANTOM READ          | V                 | V                 |         -         |       -       |
| SERIALIZATION ANOMALY | V                 | V                 | V                 |       -       |
<br>
<br>


## 10. Действие на Github Go + Postgres

> __Рабочий процесс:__
> * Является автоматизированной процедурой
> * Состоит из 1+ заданий
> * Запускается по событиям, по расписанию или вручную
> * Добавьте файл .yml в репозиторий

> __Запуск (Runner)__
> * Является ли сервер для запуска заданий
> * Запускайте по 1 заданию за раз
> * Размещено на github или самостоятельно
> * Сообщайте о ходе выполнения, журналах и результатах на github

> __Задания (Job)__
> * Представляет собой набор шагов, выполняемых в одном и том же runner
> * Обычные задания выполняются параллельно
> * Зависимые задания выполняются последовательно

> __Шаг__
> * Является отдельной задачей
> * Выполняется последовательно в рамках задания
> * Содержит более 1 действия

> __Действие__
> * Является автономной командой
> * Выполняется последовательно в пределах шага
> * Может использоваться повторно

<br>
<br>

## ---------------------------------------------

## ---------------------------------------------

# __II. Создание RESTful HTTP JSON API__

## 11. Реализация RESTful HTTP API в Go с помощью Gin (2.1)

> __Стандартный пакет net/http__
<br>

> __Popular web frameworks__
> * Gin
> * Beego
> * Echo
> * Revel
> * Martini
> * Fiber
> * Buffalo
<br>

> __Popular HTTP routers__
> * FastHttp
> * Gorilla Mux
> * HttpRouter
> * Chi
<br>

Самый популярный - __Gin__ - [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)

api/server.go   
api/account.go   
main.go

в main.go обязательно добавить импорт драйвера

    _ "github.com/lib/pq"   
    или   
    _ "github.com/go-sql-driver/mysql"   

Makefile:

    server: ## Run the application server.
        go run main.go

__В целях тестирования запросов установить Postman__

    GET http://localhost:8080/accounts

<br>
<br>


## 12. Конфигурация из файла и переменных окружения - Viper (2.2)

```shell
    $ go get github.com/spf13/viper
```

app.env   
util/config.go   
main.go   

```go
    config, err := util.LoadConfig(".")
```

```shell
    $ SERVER_ADDRESS=0.0.0.0:8081 make server
```

db/sqlc/main_test.go
<br>
<br>


## 13. Mock DB (макет) - тестирование HTTP API и 100% охвата (2.3)

__Подготовка:__

```shell
    $ go get go.uber.org/mock
    $ go install go.uber.org/mock/mockgen@latest
    $ ls -l ~/go/bin
        // проверить ~/go/bin/mockgen
    $ which mockgen
        ~/go/bin/mockgen
    Если нет, то:
    $ vi ~/.zshrc           // для mac
    $ vi ~/.bash_profile    // или для другого терминала
        i
        export PATH=$PATH:~/go/bin
        esc
        :wq
    $ source ~/.zshrc
    $ which mockgen
        ~/go/bin/mockgen
    $ mockgen -help
```

```go
    // БЫЛО:

    // api/server.go
    func NewServer(config util.Config, store *db.Store) (*Server, error)
        // для подключения к реальной базе данных используется store *db.Store
        // для тестов mock (с макетом) его надо заменить интерфейсом
    
    // db/sqlc/store.go
    type Store struct {
        db *sql.DB
        *Queries
    }
    func NewStore(db *sqL.DB) *Store {
        return &Store {
            db: db,
            Queries: New(db),
        }
    }
    func (s *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error
    func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)

    // sqlc.yaml
    emit_interface: false

    // ПЕРЕПИСЫВАЕМ:

    // sqlc.yaml
    emit_interface: true    // было false
    $ make sqlc             // обновить в терминале
        // создался новый файл с интерфейсом db/sqlc/querier.go  

    // db/sqlc/store.go
    type Store interface {
        Querier             // интерфейс из нового файла  db/sqlc/querier.go
        TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
    }
    type SQLStore struct {
        *Queries
        db *sql.DB
    }
    func NewStore(db *sql.DB) Store {
        return &SQLStore{db: db, Queries: New(db)}
    }
    func (s *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error
    func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)

    // api/server.go
    type Server struct {
        config     util.Config
        store      db.Store     // убрали * у db.Store, это теперь не указатель, а интерфейс
        tokenMaker token.Maker
        router     *gin.Engine
    }
    func NewServer(config util.Config, store db.Store) (*Server, error)  // убрали * у db.Store, это теперь не указатель, а интерфейс
```

__Создаем пакет db/moc:__

создаем папку db/moc

```shell
    $ mockgen -package mockdb -destination db/mock/store.go github.com/Country7/backend-captaincode-mysql/db/sqlc Store
        // создался файл db/mock/store.go

    // добавляем команду в файл Makefile
    mock: ## Generate a store mock.
	    mockgen -package mockdb -destination db/mock/store.go github.com/Country7/backend-captaincode-mysql/db/sqlc Store
```

__Приступаем к написанию тестов:__

api/account_test.go   
<br>
<br>


## 14. Пользовательский валидатор параметров - __transfer__ (2.4)

api/transfer.go   

api/server.go   
```go
    authRoute.POST("/transfers", server.createTransfer)
```

__Postman:__

    POST http://localhost:8080/transfers
    Body raw JSON
    {
        "from_account_id": 1,
        "to_account_id": 2,
        "amount": 10,
        "currency": "USD"
    }

api/validator.go   
```go
    import ( "github.com/go-playground/validator/v10" )
    var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
        if currency, ok := fieldLevel.Field().Interface().(string); ok {
            return util.IsSupportedCurrency(currency)
        }
        return false
    }
```

util/currency.go   
```go
    const (
        USD = "USD"
        EUR = "EUR"
        CAD = "CAD" )
    // IsSupportedCurrency returns true if the currency is supported.
    func IsSupportedCurrency(currency string) bool {
        switch currency {
        case CAD, EUR, USD:
            return true
        default:
            return false
        }
    }
```

__Регистрация вадидатора на сервере__   
api/server.go
```go
    import ("github.com/gin-gonic/gin/binding")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
```

api/account.go
```go
    type createAccountRequest struct {
        ...
        Currency string `json:"currency" binding:"required,currency"`
    }
```

api/transfer.go
```go
    type transferRequest struct {
        ...
        Currency string `json:"currency" binding:"required,currency"`
    }
```

__Postman:__
```json
    POST http://localhost:8080/transfers
    Body raw JSON
    {
        "from_account_id": 1,
        "to_account_id": 2,
        "amount": 10,
        "currency": "EUR"
    }
```

util/random.go
```go
    // RandomCurrency generates a random currency code.
    func RandomCurrency() string {
        currencies := []string{EUR, USD, CAD}
        n := len(currencies)
        return currencies[rand.Intn(n)]
    }
```
<br>
<br>


## 15. Добавление таблицы users с ограничениями уникальности и внешнего ключа (2.5)

[diagram.io](https://dbdiagram.io/home)

->  Go to App   

    Table users as U {
        username varchar [pk]
        role varchar [not null, default: 'depositor']
        hashed_password varchar [not null]
        full_name varchar [not null]
        email varchar [unique, not null]
        is_email_verified bool [not null, default: false]
        password_changed_at timestamptz [not null, default: '0001-01-01']
        created_at timestamptz [not null, default: `now()`]
    }

    Table accounts as A {
        ...
        owner varchar [ref: > U.username, not null]
        ...
        Indexes {
            owner
            (owner, currency) [unique]  // в одной валюте только один счет у пользователя
                                        // в разной валюте может быть несколько счетов
        }
    }

Export PostgreSQL

```shell
    $ migrate -help
    $ migrate create -ext sql -dir db/migration -seq add_users
```

В созданный файл db/migration/000002_add_users.up.sql копируем изменения из doc/schema.sql (таблицу users и ключи)

```shell
    $ make migrateup
        // ошибка, так как данные accounts есть, а их в новой таблице users - нет
    $ make migratedown
        // ошибка, надо вручную менять значение в БД / таблице schema_migrations с TRUE на FALSE
    $ make migratedown
        // удалились все таблицы
    $ make migrateup
```

В файле db/migration/000002_add_users.down.sql грохаем ключи, грохаем таблицу users

<br>
<br>


## 16. Обработка ошибок базы данных (2.6)

db/query/user.sql
```shell
    $ make sqlc
```

Появился db/sqlc/user.sql.go   

Изменились  
db/sqlc/models.go   
db/sqlc/querier.go   

Пишем db/sqlc/user_test.go

Правим db/sqlc/account_test.go
```go
    func createRandomAccount(t *testing.T) Account {
        user := createRandomUser(t)
        arg := CreateAccountParams{
            Owner:    user.Username,
            ...
        }
```

```shell
    $ make mock
    $ make test
```

Допиливаем api/account.go
```go
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		var pqErr *pq.Error                         // добавляем от сюда
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
                    // ошибки на сервере при создании аккаунта без юзера,
                    // и создании аккаунта с одинаковой валютой счета
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}                                           // до сюда
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
```

__Postman:__
```json
    POST http://localhost:8080/users
    Body raw JSON
    {
        "username": "QuangBang",
        "password": "secret",
        "full_name": "Quang Bang",
        "email": "quang@mail.com"
    }
    POST http://localhost:8080/users/login
    Body raw JSON
    {
        "username": "QuangBang",
        "password": "secret"
    }

При работе с mysql выдает ошибку, что полю PasswordChangedAt time.Time - row.Scan не может присвоить значение []int8.   
Решается просто. В файле app.env к переменной DB_SOURCE добавляем параметр ?parseTime=true.   
?? под вопросом, нужно проверить - И, кстати, убираем из этой переменной mysql://.   
Получается: DB_SOURCE=root:secret@tcp(localhost:3306)/main_db?parseTime=true   
Лучше с кавычками: DB_SOURCE="root:secret@tcp(localhost:3306)/main_db?parseTime=true"   
А то потом, при развертывании docker контейнеров могут быть ошибки.

    POST http://localhost:8080/accounts
    Authorization  Bearer Token  ...
    Body raw JSON
    {
        "owner": "QuangBang",
        "currency": "USD"
    }
    GET http://localhost:8080/accounts?page_id=1&page_size=5
    key   page_id 1   page_size 5
    Body raw JSON
    {
        "owner":  "QuangBang",
        "limit":  5,
        "offset": 0
    }
```
<br>
<br>

## ---------------------------------------------

## 17. Безопасное хранение паролей Hash password в Go с помощью Bcrypt (2.7)

util/password.go   
api/user.go   
api/validator.go   

github.com/go-playground/validator

    alphanum	Alphanumeric
    email	    E-mail String
<br>
<br>


## 18. Модульные тесты с помощью gomock (2.8)
<br>

## 19. Почему PASETO лучше, чем JWT, для аутентификации на основе токенов (2.9)

### Token-based Authentication

|           |                                                                               |  |        |
|:-:        |:-                                                                             |:-|:-:     |
| Client    | 1. POST /users/login <br> -----------------------> <br> {username, password}  |  | Server |

|           |                                                                               |                   |           |
|:-:        |-:                                                                             |:-                 |:-:        |
| Client    | 200 OK <br> <----------------------- <br> {access_token: JWT, PASETO, ...}    | <-- Sign token    | Server    |

|           |                                                                                           |  |        |
|:-:        |:-                                                                                         |:-|:-:     |
| Client    | 2. GET /accounts <br> -----------------------> <br> Authorization: Bearer <access_token>  |  | Server |

|           |                                                                       |                   |           |
|:-:        |-:                                                                     |:-                 |:-:        |
| Client    | 200 OK <br> <----------------------- <br> [account1, account2, ...]   | <-- Verify token  | Server    |


### АЛГОРИТМЫ ПОДПИСИ JWT

```json
    header:
    {
        "typ": "JWT",
        "alg": "HS256"
    }
    payload:
    {
        "id": "1337",
        "username": "bizone",
        "iat": 1594209600,
        "role": "user"
    }
    signature:
    ZvkYYnyM929FM4NW9_hSis7_x3_9rymsDAx9yuOcc1I
```

__Алгоритм симметричной цифровой подписи__
* Для подписи и проверки используется один и тот же секретный ключ, токен
* Для локального использования: внутренние службы, где можно совместно использовать секретный ключ
* HS256, HS384, HS512  
    - HS256 = HMAC + SHA256  
    - HMAC: Hash-based Message Authentication Code - Код аутентификации сообщения на основе хэша  
    - SHA: Secure Hash Algorithm - Алгоритм безопасного хэширования  
    - 256/384/512: количество выходных битов   

__Алгоритм асимметричной цифровой подписи__
* Закрытый ключ используется для подписи токена
* Открытый ключ используется для проверки токена
* Для публичного использования: внутренняя служба подписывает токен, но внешняя служба должна его подтвердить
* RS256, RS384, RS512 || PS256, PS384, PS512 || ES256, ES384, ES512  
    - RS256 = RSA PKCSv1.5 + SHA256 [PKCS: Public-Key Cryptography Standards - Стандарты криптографии с открытым ключом]  
    - PS256 = RSA PSS + SHA256 [PSS: Probabilistic Signature Scheme - Вероятностная схема подписи]  
    - ES256 = ECDSA + SHA256 [ECDSA: Elliptic Curve Digital Signature Algorithm - Алгоритм цифровой подписи с эллиптической кривой]   

### В чем проблема JWT?

__Слабые алгоритмы__
* Дают разработчикам слишком много алгоритмов на выбор
* Известно, что некоторые алгоритмы уязвимы:
* RSA PKCSv1.5: атака на oracle с дополнением
* ECDSA: атака с недопустимой кривой

__Тривиальная подделка__
* Установите для заголовка "alg" значение "none"
* Установите для заголовка "alg" значение "HS256", в то время как сервер обычно проверяет токен с помощью открытого ключа RSA


### Platform-Agnostic SEcurity TOkens [PASETO] Независимые от платформы токены безопасности

__Более надежные алгоритмы__
* Разработчикам не нужно выбирать алгоритм 
* Нужно только выбрать версию PASETO 
* Каждая версия имеет 1 набор надежных шифров 
* Принимаются только 2 самые последние версии PASETO

__Нетривиальная подделка__
* Больше никакого заголовка "alg" или алгоритма "none"
* Все аутентифицировано
* Зашифрованная полезная нагрузка для локального использования <симметричный ключ>   

- v1 [совместима с устаревшей системой]
    + локальный: <симметричный ключ>
        - Аутентифицированное шифрование
        - AES256 CTR + HMAC SHA384
    + открытый: <асимметричный ключ>
        - Цифровая подпись
        - RSA PSS + SHA384

* v2 [рекомендуется]
    + локальный: <симметричный ключ>
        - Аутентифицированное шифрование
        - XChaCha20-Poly1305
    + открытый: <асимметричный ключ>
        - Цифровая подпись
        - Ed25519 [EdDSA + Curve25519]

```javascript
• Version: v2
• Purpose: public [asymmetric-key digital signature]
• Payload:
    • Body:
        • Encoded: [base64]
        eyJeHAiO¡IyMDM5LTAxLTAxVDAwOjAwOjAwKzAwOjAwIiwiZGFOYSI
        6InRoaXMgaXMgYSBzaWduZWQgbWVzc2FnZSJ91g
        • Decoded:
        {
        "data": "this is a signed message" ,
        "exp": "2039-01-01T00:00:00+00:00"
        }
    • Signature: [hex-encoded]
    d600bbfa3096b0dde6bf8b89699c59a746ed2c981cc95c0bfacbc90fb7
    f8207c86b5e29edc74cb8c761318723532d0aa27e1120cb36813ba2d90
    8cda985b2408
```
<br>
<br>


## 20. Создать и верифицировать токен JWT & PASETO (2.10)

token/maker.go   
token/payload.go   
token/jwt_maker.go   
token/jwt_maker_test.go   

```shell
    $ go get github.com/google/uuid
    $ go get github.com/golang-jwt/jwt/v5
```

token/paseto_maker.go   
token/paseto_maker_test.go   

```shell
    $ go get github.com/o1egl/paseto/v2
```
<br>
<br>


## 21. API для входа в систему через токен PASETO или JWT (2.11)

api/server.go   
app.env   
util/config.go   
api/main_test.go   
api/transfer_test.go   
api/user_test.go   
api/account_test.go   
main.go   

api/user.go   
api/server.go  router   

При работе с mysql выдает ошибку, что слишком длинную строку RefreshToken сервер пытается записат в БД.   
Длина токена получается 335 символов когда в базе заявлено 255. Лечится просто:   
В db/migration/000001_init.up.sql строку `refresh_token` varchar(255) NOT NULL в CREATE TABLE `sessions`   
делаем длиной 512 - `refresh_token` varchar(512) NOT NULL   

__утилита Postman:__

```json
    POST http://localhost:8080/users/login
    Body raw JSON
        {
            "username": "qwe",
            "password": "1234567"
        }

    Response 200 OK
        {
            "session_id": "e7a1856e-aaa0-4226-9681-6ff7993dd789",
            "access_token": "v2.local. ...",
            "access_token_expires_at": "2024-03-30T16:05:55.116046+03:00",
            "refresh_token": "v2.local. ...",
            "refresh_token_expires_at": "2024-03-31T15:50:55.116324+03:00",
            "user": {
                "username": "qwe",
                "full_name": "asdfgh",
                "email": "zxc@mail.com",
                "password_changed_at": "0001-01-01T00:00:00Z",
                "created_at": "2024-03-27T17:34:06.837219Z"
            }
        }
```
<br>
<br>



## 22. Middleware авторизации (2.12)

__Что такое промежуточное программное обеспечение?__

```
|        | Send request                            |         |
|        | -------------->        Route            |         |
|        |                  /accounts/create       |         |
|        |                          |              |         |
|        |                          |              |         |
|        |                          V              |         |
|        |  ctx. Abort()       Middlewares         |         |
| CLIENT | <--------------    Logger (ctx),        | SERVER  |
|        | Send response        Auth(ctx)          |         |
|        |                          |              |         |
|        |                          | ctx.Next()   |         |
|        |                          |              |         |
|        |                          V              | Авторизация:
|        | Send response         Handler  <--------- У пользователя
|        | <--------------  createAccount(ctx)     | есть разрешение?
|        |                                         |         |
```
<br>

api/middleware.go   
api/middleware_test.go   

api/server.go   

```go
    router := gin.Default()
    router.POST("/users/login", server.loginUSer)

    authRoute := router.Group("/").Use(authMiddleware(server.tokenMaker))
    authRoute.GET("/accounts", server.listAccount)
```
<br>

__ПРАВИЛА АВТОРИЗАЦИИ__

| | | |
|:-:|:-:|:-:|
| API <br> Create account | -------> | Правило <br> Авторизованный пользователь может создать <br> учетную запись только для себя |
| API <br> Get account    | -------> | Правило <br> Авторизованный пользователь может получить <br> только те учетные записи, которыми он владеет|
| API <br> List accounts  | -------> | Правило <br> Авторизованный пользователь может перечислять <br> только те учетные записи, которые принадлежат ему|
| API <br> Transfer money | -------> | Правило <br> Авторизованный пользователь может отправлять <br> деньги только со своего собственного аккаунта|
<br>

```go
api/account.go func createAccount

    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)    // добавили
    arg := db.CreateAccountParams{
        Owner:    authPayload.Username,   // было req.Owner
        Balance:  0,
        Currency: req.Currency,
    }

api/account.go func getAccount
    
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)    // добавили
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
```

```sql
db/query/account.sql

    -- name: ListAccounts :many
    SELECT * FROM accounts
    WHERE owner = $1            // добавлено условие
    ORDER BY id LIMIT $2
    OFFSET $3;

    make sqlc       // db/sqlc/account.sql.go listAccounts обновился
    make mock
```

```go
db/sqlc/account_test.go  TestListAccounts   

api/account.go  func listAccount

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)    // добавили
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,       // добавили
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

api/transfer.go  

    func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool)
    // добавили db.Account
    return account, true

    func createTransfer
        fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
        if !valid {
            return
        }
        authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
        if fromAccount.Owner != authPayload.Username {
            err := errors.New("from account doesn't belong to the authenticated user")
            ctx.JSON(http.StatusUnauthorized, errorResponse(err))
            return
        }
        _, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
        if !valid {
            return
        }

api/account_test.go
```

<br>
<br>
<br>

## ---------------------------------------------
## ---------------------------------------------

# __III. Развертывание приложения в рабочей среде (Deploying the application to production)__

## 23. Образ Golang Docker с помощью многоступенчатого файла Dockerfile (3.1)

    $ git checkout -b deploying

update go

go.mod   

    go 1.22

.github/workflows/test.yml

    go-version: '1.22'

    $ git status
    $ git add .
    $ git status
    $ git commit -m"update go to 1.22"
    $ git push -u origin deploying

Из терминала переходим по ссылке    
<https://github.com/Country7/backend-captaincode-mysql/pull/new/deploying>

    Name -> Add docker
    -> Create pull request

    $ brew upgrade golang-migrate   // для mac
    $ curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
        // для linux
    $ migrate -version
        v4.17.0
    $ make migrate-up

.github/workflows/test.yml

    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

    $ git add .
    $ git status
    $ git commit -m"upgrade golang-migrate to v4.17.0"
    $ git push

__Если тесты на github пройдены, то приложение готово к 1 запуску__

Создаем __Dockerfile__:

<https://hub.docker.com/_/golang>

    FROM golang:1.22.2-alpine3.19
    WORKDIR /app
    COPY . .
    RUN go build -o main main.go
    EXPOSE 8080
    CMD [ "/app/main" ]

    $ docker build --help
    $ docker build -t captaincode:latest .
    $ docker images
        captaincode latest 601MB

Размер образа получился 601Мб, чтобы уменьшить размер образа нужно применить многоступенчатую сборку. Нам в образе нужен только исполняемый файл.   
__Dockerfile__:

    # Build stage
        FROM golang:1.22.2-alpine3.19 AS builder
        WORKDIR /app
        COPY . .
        RUN go build -o main main.go
    # Run stage
        FROM alpine:3.19
        WORKDIR /app
        COPY --from=builder /app/main .
        EXPOSE 8080
        CMD [ "/app/main" ]

    $ docker build -t captaincode:latest .
    $ docker images
        captaincode latest 21MB

    $ docker rmi 61504815c89a   // удалить старый образ IMAGE ID = 61504815c89a

<br>
<br>

## 24. Подключить контейнеры в одной сети docker (3.2)

    $ docker run --name captaincode -p 8080:8080 captaincode:latest
        cannot load config:Config File "app" Not Found in "[/app]"
    $ docker ps -a
    $ docker rm captaincode
    $ docker images
    $ docker rmi a4809227e909

__Dockerfile__:

    COPY app.env .

    $ docker build -t captaincode:latest .
    $ docker images
    $ docker run --name captaincode -p 8080:8080 captaincode:latest
        Running in "debug" mode. Switch to "release" mode in production.
        - using env:	export GIN_MODE=release
    $ docker rm captaincode
    $ docker run --name captaincode -p 8080:8080 -e GIN_MODE=release captaincode:latest
    $ docker ps

__Postman__: "error": "dial tcp 127.0.0.1:5432: connect: connection refused"   
__Terminal__: [GIN] | 500 | 1.437083ms | 192.168.65.1 | POST "/users/login"    
__app.env__: DB_SOURCE=postgresql://root:secret@localhost:5432/main_db?sslmode=disable

    $ docker container inspect postgres16
        "NetworkSettings": "Networks": "bridge": "IPAddress": "172.17.0.2"
    $ docker container inspect captaincode
        "NetworkSettings": "Networks": "bridge": "IPAddress": "172.17.0.3"

    $ docker stop captaincode
    $ docker rm captaincode
    $ docker run --name captaincode -p 8080:8080 -e GIN_MODE=release -e "DB_SOURCE=postgresql://root:secret@172.17.0.2:5432/main_db?sslmode=disable" captaincode:latest

__Postman__: Status 200 OK

### Способ получше (подключиться к контейнеру postgres16 по имени, а не по ip адресу)

    $ docker rm captaincode
    $ docker network ls
        9a00594f4037 bridge bridge local
    $ docker network inspect bridge
        "Containers": "Name": "postgres16"
    // контейнеры в мостовой сети bridge не могут видеть друг друга по имени, как в других сетях
    // поэтому нужно создать свою сеть и подключить к ней контейнер
    $ docker network --help
    $ docker network create ww-network
    $ docker network connect --help 
    $ docker network connect ww-network postgres16
    $ docker container inspect postgres16
        "Networks":
            "bridge": "IPAddress": "172.17.0.2"
            "ww-network": "IPAddress": "172.18.0.2"
    $ docker run --name captaincode --network ww-network -p 8080:8080 -e GIN_MODE=release -e "DB_SOURCE=postgresql://root:secret@postgres16:5432/main_db?sslmode=disable" captaincode:latest

c __mysql__ это будет так:   
для mysql в app.env DB_SOURCE должна быть в кавычках:   
    DB_SOURCE="mysql://root:secret@tcp(mysql8:3306)/main_db?parseTime=true"

Кроме того при открытии соединения    
    conn, err := sql.Open(config.DBDriver, dsnDBSource)   
переменная подключения dsnDBSource должна быть не в формате URL, а в формате DSN,   
т.е. без префикса mysql://   
    
    1.
    $ docker network create cc-network
    2.
    $ docker run --name mysql8 -p 3306:3306 -e MYSQL_DATABASE=main_db -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0
    3.
    $ docker network connect cc-network mysql8
    $ docker network inspect cc-network
    $ docker container inspect mysql8
    4.
    $ migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/main_db" -verbose up
    5.
    kubuntu:
    $ sudo apt-get update
    $ sudo apt-get install -y mysql-client
    mac:
    $ brew install mysql-client
        /Users/country/.zshrc:
        export PATH=/opt/homebrew/opt/mysql-client/bin:$PATH
        export LDFLAGS="-L/opt/homebrew/opt/mysql-client/lib"
        export CPPFLAGS="-I/opt/homebrew/opt/mysql-client/include"
        export PKG_CONFIG_PATH="/opt/homebrew/opt/mysql-client/lib/pkgconfig"
    $ mysql --version
    6.
    init.sql - для создания пользователя и назначения ему прав
    $ mysql -u root -p -h localhost -P 3306 --protocol=TCP main_db --password="secret" < "./init.sql"
    7.
    $ mysql -u mysqluser -p -h localhost -P 3306 --protocol=TCP main_db --password="secretpass"
        SHOW GRANTS;
        \q
    8.
    $ mysql -u root -p -h localhost -P 3306 --protocol=TCP main_db --password="secret"
        SELECT user, host FROM mysql.user WHERE user = 'mysqluser';
        SHOW GRANTS FOR 'mysqluser'@'%';
        или удалить пользователя DROP USER 'mysqluser'@'%';

    9.
    $ docker build -t captaincode:latest .  // ubuntu 244 Mb, alpine 158 Mb
    10. 
    $ docker run --name captaincode --network cc-network -p 8080:8080 -e "DB_SOURCE=mysql://mysqluser:secretpass@tcp(mysql8:3306)/main_db?parseTime=true" captaincode:latest
    11.
    $ docker exec -it captaincode sh
        printenv  // вывод подключенных переменных
        mysql -u root -p -h mysql8 main_db --password="secret"
            SELECT user, host FROM mysql.user WHERE user = 'mysqluser';
            SHOW GRANTS FOR 'mysqluser'@'%';
            \q
        mysql -u mysqluser -p -h mysql8 main_db --password="secretpass"
            SHOW GRANTS;
            \q
        exit
    12.
    $ docker network inspect cc-network
    $ docker container inspect captaincode
    13.
    POST http://localhost:8080/users
    Content-Type: application/json
    {
        "username": "QuangBang",
        "password": "secret",
        "full_name": "Quang Bang",
        "email": "quang@mail.com"
    }
    14.
    $ docker logs captaincode
    15.
    $ docker exec -it captaincode sh
    env | grep DATABASE

    к команде можно добавить e GIN_MODE=release

MySQL_Allow_Empty_Password

__Postman__: Status 200 OK   
__Terminal__: [GIN] | 200 | 101.980875ms | 192.168.65.1 | POST "/users/login"

    $ docker network inspect ww-network
        "Containers":
            "Name": "postgres16"
            "Name": "captaincode",

Makefile

    run-postgres: ## Run postgresql database docker image.
	    docker run --name postgres16 --network ww-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

<br>
<br>

## 25. Файл docker-compose (3.3)

Docker Compose для автоматической настройки всех служб   

docker-compose.yaml     // .yaml очень чувствителен к пробелам, устанавливаем 2 пробела для Tab Size

    $ docker compose up
    $ docker images
    $ docker ps
    $ docker network inspect ww-network

__Postman__: Status: 500 Internal Server Error

    {
        "error": "pq: relation \"users\" does not exist"
    }
    // Потому как не было миграции

Допиливаем __Dockerfile__ для добавления миграции:

    # Build stage
    RUN apk add curl
    RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

    # Run stage
    COPY --from=builder /app/migrate ./migrate
    COPY start.sh .
    COPY db/migration ./migration
    ENTRYPOINT [ "/app/start.sh" ]

Создаем файл start.sh

    $ chmod +x start.sh    // делаем его исполняемым

__start.sh__:

    #!/bin/sh
    set -e
    echo "run db migration"
    source /app/app.env
    /app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up
    echo "start the app"
    exec "$@"

    $ docker compose down  // удалит все контейнеры и сети
    $ docker image ls
    $ docker rmi api
    $ docker network ls
    $ docker compose up

!!! При запуске Docker Compose на Linux (Kubuntu) переменную $DB_SOURCE   
он взял не из docker-compose.yaml, где к адресу БД обращение по имени postgres,   
а из файла app.env, где адрес указан localhost.   

    api | error: dial tcp 127.0.0.1:5432: connect: connection refused

Пришлось в app.env внести изменения:

    # DB_SOURCE=postgresql://root:secret@localhost:5432/main_db?sslmode=disable   
    DB_SOURCE=postgresql://root:secret@postgres:5432/main_db?sslmode=disable

    для mysql нужно в кавычках:
    DB_SOURCE="mysql://root:secret@tcp(mysql8:3306)/main_db?parseTime=true"
    для прохождения тестов для mysql нужно в формате DSN - без префикса mysql://
    DB_SOURCE=root:secret@tcp(localhost:3306)/main_db?parseTime=true

[GIN] | 200 | 144.777766ms | 172.20.0.1 | POST "/users"   
[GIN] | 200 | 123.707554ms | 172.20.0.1 | POST "/users/login"

<br>
<br>

## Deploy на сервер

Создаем .github/workflows/deploy.yml

Варианты рабочих деплоев:

    doc/deploy_ubuntu-latest.yml
    doc/deploy_cont_ubuntu.13.04.yml
    doc/deploy_cont_centos7.yml

### SSH ключ для подключения deploy к серверу reg.ru

    $ ssh-keygen -t rsa -b 4096 -C "vadim-vvs@mail.ru"
    $ ssh-copy-id root@123.123.123.123
    ~/.ssh/authorized_keys на вашем сервере
    GitHub -> Repositories -> Settings -> Secrets and variables -> Actions, SSH_PRIVATE_KEY.

### Миграции на сервере reg.ru

    $ ssh -t user@website-reg.ru
        mysql --version
        migrate -version
        go version
        curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
        set -o allexport; source ~/app/app.env; set +o allexport
        echo "DATABASE_URL: $DATABASE_URL"
        ls -la ~/app/db/migration

        ps aux | grep api_app

### Go и GLIBC 2.17 на сервере reg.ru

    wget https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
    sha256sum go1.22.4.linux-amd64.tar.gz
    tar -C ~/usr/local -xzf go1.22.4.linux-amd64.tar.gz
    vi $HOME/.profile 
        export PATH=$PATH:/usr/local/go/bin
    source $HOME/.profile
    go version
    go tool dist list
    go env GOOS GOARCH
        -bash-4.2$ go env GOOS GOARCH
            linux
            amd64
    ldd --version
    ldd ~/app/api_app

    $ GOOS=linux GOARCH=amd64 go build

    -bash-4.2$ ldd --version
        ldd (GNU libc) 2.17
    -bash-4.2$ ldd ~/app/api_app

### ENV на сервере reg.ru

    $ ssh -t user@website-reg.ru
        mysql --version
        migrate -version
        go version
        go env GOOS GOARCH
        ldd --version
        ldd ~/app/api_app

        source ~/app/app.env
        echo $SERVER_ADDRESS
        cd ~/app
        ./api_app
        ps aux | grep api_app   // показывает работающий процесс api_app и процесс grep который ищет api_app
        ps aux | grep api_app | grep -v grep   // показывает только работающий процесс api_app
        pkill api_app || true   // останавливает процесс

        lsof -i: 80 | grep LISTEN   // для работы нужно установить lsof
        netstat -tulpn | grep :8080  // показывает процесс, который использует порт 8080, если есть root права

### Скачивание, распаковка и установка GLIBC 2.17

    - name: Install GLIBC 2.17
        run: |
          sudo apt-get update
          sudo apt-get install -y wget build-essential
          wget http://ftp.gnu.org/gnu/libc/glibc-2.17.tar.gz
          tar -xvf glibc-2.17.tar.gz
          cd glibc-2.17
          mkdir build
          cd build
          ../configure --prefix=/opt/glibc-2.17
          make -j$(nproc)
          sudo make install

### Запуск приложения на сервере в фоновом режиме

Этот команда запустит api_app в фоновом режиме, перенаправит вывод в файл api_app.log    
и освободит консоль. Процесс будет продолжать работать даже после закрытия консоли

    nohup ./api_app &> api_app.log &

screen — это оконный менеджер, который позволяет запускать и контролировать несколько терминалов из одной консоли. Вы можете использовать screen для запуска приложения и отделения его от текущей сессии.

    sudo yum install screen
    screen -S myapp
    ./api_app
    Для отделения от сессии screen, нажмите Ctrl+A, затем D.
    Чтобы снова подключиться к сессии:
    screen -r myapp

tmux — это другой оконный менеджер, который можно использовать аналогично screen.

    sudo yum install tmux
    tmux new -s myapp
    ./api_app
    Для отделения от сессии tmux, нажмите Ctrl+B, затем D.
    Чтобы снова подключиться к сессии:
    tmux attach -t myapp

Использование systemd
Если вы хотите запустить ваше приложение как сервис, вы можете создать systemd юнит. Это предпочтительный метод для длительно работающих процессов.

    Создайте файл юнита:
    sudo nano /etc/systemd/system/api_app.service

    [Unit]
    Description=API Application

    [Service]
    ExecStart=./api_app
    WorkingDirectory=~/app
    Restart=always

    [Install]
    WantedBy=multi-user.target

    Замените /path/to/your на реальный путь к вашему приложению.
    Затем выполните следующие команды для управления сервисом:
    sudo systemctl daemon-reload
    sudo systemctl start api_app.service
    sudo systemctl enable api_app.service

Теперь ваше приложение будет запускаться при старте системы и работать в фоновом режиме.

<br>
<br>

__Сгенерируем ключ-строку из 32 символов__:

    $ openssl rand -hex 64
        // строка 128 символов
    $ openssl rand -hex 64 | head -c 32
        // строка 32 символа




<br>
<br>
<br>

# PS

### Собрать все зависимости из go.mod

    $ go mod tidy


