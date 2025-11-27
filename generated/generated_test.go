package generated

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
	"time"
)

func TestInsertCustomer(t *testing.T) {

	//db, err := sql.Open("postgres", "pgx:pgxpass@localhost:5432/pgx_test?sslmode=disable")

	db, err := sql.Open("postgres", "postgres://pgx:pgxpass@localhost:5432/pgx_test?sslmode=disable")

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	billto := NullAddressType{
		AddressType: AddressType{
			StreetLine1: sql.NullString{"123 Main St", true},
		},
		Valid: true,
	}

	customer := &Customer{
		Name:          "John Doe",
		CustomerType:  CustomerTypeCommercial,
		CreatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		BilltoAddress: billto,
	}

	err = customer.Insert(context.Background(), db)
	if err != nil {
		t.Fatalf("failed to insert customer: %v", err)
	}

}
