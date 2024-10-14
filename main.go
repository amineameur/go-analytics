package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func connectToClickHouse() clickhouse.Conn {
	opts := &clickhouse.Options{
		Addr: []string{"172.18.0.2:9001"},
		Auth: clickhouse.Auth{
			Database: "analytics",
			Username: "default",
			Password: "",
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	}

	conn, err := clickhouse.Open(opts)
	if err != nil {
		log.Fatalf("Could not connect to ClickHouse: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	return conn
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--migration" {
		runMigration()
	} else {
		fetchAnalyticsData()
	}
}

func runMigration() {
	conn := connectToClickHouse()
	defer conn.Close()

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS invoices (
			id UInt32,
			amount Float32,
			created_at Date
		) ENGINE = MergeTree()
		ORDER BY id;
	`
	err := conn.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration complete: Created 'invoices' table.")

	for i := 1; i <= 100; i++ {
		year := 2024
		month := (i % 12) + 1
		day := (i % 28) + 1
		invoiceDate := fmt.Sprintf("%d-%02d-%02d", year, month, day)

		insertQuery := fmt.Sprintf("INSERT INTO invoices VALUES (%d, %.2f, '%s')", i, float32(i)*100, invoiceDate)

		err := conn.Exec(context.Background(), insertQuery)
		if err != nil {
			log.Fatalf("Error inserting data: %v", err)
		}
	}

	log.Println("Inserted 100 fake invoices.")
}

func fetchAnalyticsData() {
	conn := connectToClickHouse()
	defer conn.Close()

	query := `
		SELECT formatDateTime(created_at, '%Y-%m') AS month, sum(amount) AS total_revenue
		FROM invoices
		GROUP BY month
		ORDER BY total_revenue DESC
		LIMIT 1;
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Error fetching analytics: %v", err)
	}
	defer rows.Close()

	var month string
	var totalRevenue float64
	for rows.Next() {
		err := rows.Scan(&month, &totalRevenue)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		log.Printf("Most revenue-generating month: %s, Total Revenue: %.2f", month, totalRevenue)
	}
}
