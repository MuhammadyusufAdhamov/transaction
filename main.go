package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	PostgresUser     = "postgres"
	PostgresDatabase = "crud"
	PostgresPassword = "7"
	PostgresHost     = "localhost"
	PostgresPort     = 5432
)

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s database=%s password=%s sslmode=disable",
		PostgresHost,
		PostgresPort,
		PostgresUser,
		PostgresDatabase,
		PostgresPassword,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed connect to database: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			fmt.Println("transaction rolled back")
			err := tx.Rollback()
			if err != nil {
				return
			}
		} else {
			fmt.Println("transaction committed")
			err := tx.Commit()
			if err != nil {
				log.Fatalf("failed to commit transaction: %v", err)
			}
		}
	}()

	queryFromAccount := `
			update accounts set balance=balance-$1 where id=$2
`

	queryToAccount := `
		update accounts set balance=balance+$1 where id=$2
`

	queryTransfer := `
		insert into transfer (
		                      from_account_id,
		                      to_account_id,
		                      amount,
		                      payment_type
		) values ($1,$2,$3,$4)
`

	fromAccountId := 1
	toAccountId := 2
	amount := 100
	paymentType := "click"

	//	insert transfer
	_, err = tx.Exec(
		queryTransfer,
		fromAccountId,
		toAccountId,
		amount,
		paymentType,
	)

	if err != nil {
		fmt.Println("failed to insert transfer", err)
		return
	}

	if true {
		err := errors.New("some error")
		if err != nil {
			return
		}
	}

	_, err = tx.Exec(queryFromAccount, amount, fromAccountId)
	if err != nil {
		log.Fatalf("failed to update account: %v", err)
		return
	}

	_, err = tx.Exec(queryToAccount, amount, toAccountId)
	if err != nil {
		log.Fatalf("failed to update account: %v", err)
		return
	}

}
