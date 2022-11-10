package main

import (
	"database/sql"
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

type Transfer struct {
	ToAccountName   string
	FromAccountName string
	Amount          float64
	PaymentType     string
}

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

	query := `
		select 
		    to_account_name,
			from_account_name,
			amount,
			payment_type
		from transfer_view
`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("failed to execute query: %v", err)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	transfers := make([]Transfer, 0)
	for rows.Next() {
		var t Transfer
		err := rows.Scan(
			&t.ToAccountName,
			&t.FromAccountName,
			&t.Amount,
			&t.PaymentType,
		)
		if err != nil {
			return
		}
		transfers = append(transfers, t)
	}

	fmt.Println(transfers)

	var accountsCount int

	queryAccountCount := `select get_accounts_count();`
	err = db.QueryRow(queryAccountCount).Scan(&accountsCount)
	if err != nil {
		fmt.Printf("failed to scan accounts count: %v", err)
		return
	}

	fmt.Println("Accounts count: ", accountsCount)

	queryInsertAccount := `call insert_account($1,$2);`
	_, err = db.Exec(queryInsertAccount, "Harry", "Kane")
	if err != nil {
		fmt.Printf("failed to insert account: %v", err)
		return
	}
}
