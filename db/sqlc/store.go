package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions.
type SQLStore struct {
	*Queries
	db *sql.DB
}

const maxRetries = 3

func NewStore(db *sql.DB) Store {
	fmt.Println(">> db.NewStore Initializing Store with DB: ", db) // Выводим информацию о базе данных
	// Проверяем строку подключения для отладки
	dbStats := db.Stats()
	fmt.Printf(">> db.NewStore DB Stats: %+v\n", dbStats)

	// Проверяем базу данных
	if err := db.Ping(); err != nil {
		fmt.Printf(">> db.NewStore Error pinging DB: %v\n", err)
	}

	return &SQLStore{db: db, Queries: New(db)}
}

// execTx executes a function within a database transaction.
func (s *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
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
			// Check if error is a deadlock error
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1213 {
				// Deadlock found, retry
				continue
			}
			return err
		}
		err = tx.Commit()
		if err == nil {
			return nil
		}
		// Check if error is a deadlock error
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1213 {
			// Deadlock found, retry
			continue
		}
		return err
	}
	return fmt.Errorf("execTx: maximum retries reached: %v", err)
}

// TransferTxParams contains the input parameters of the transfer transaction.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries,
// and update accounts' balance within a single database transaction.
func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// Определение порядка блокировок
		// Determining the order of locks
		// Получение аккаунтов с блокировкой
		// Getting blocked accounts
		if arg.FromAccountID < arg.ToAccountID {
			_, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				return err
			}
			_, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				return err
			}
		} else {
			_, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
			if err != nil {
				return err
			}
			_, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
			if err != nil {
				return err
			}
		}

		// Создание записи о переводе
		// Creating a transfer record
		err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}
		result.Transfer, err = q.GetLastTransfer(ctx)
		if err != nil {
			return err
		}

		// Создание записи о списании и зачислении
		// Creating a record of debiting and crediting
		err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.GetLastEntry(ctx)
		if err != nil {
			return err
		}
		err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.GetLastEntry(ctx)
		if err != nil {
			return err
		}

		// Обновление баланса аккаунтов
		// Updating account balance
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64) (account1 Account, account2 Account, err error) {

	err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}
	account1, err = q.GetAccount(ctx, accountID1)
	if err != nil {
		return
	}

	err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}
	account2, err = q.GetAccount(ctx, accountID2)

	return
}
