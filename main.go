package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

// Connect to ClickHouse using the native client
func connectToClickHouse() clickhouse.Conn {
	// Use the container IP address instead of localhost
	opts := &clickhouse.Options{
		Addr: []string{"172.18.0.2:9000"}, // Update with your ClickHouse IP or hostname
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

	// Ping the database to ensure the connection is working
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	return conn
}

func main() {
	// Check if the --migration flag is provided
	if len(os.Args) > 1 && os.Args[1] == "--migration" {
		runMigration()
	} else {
		fetchAnalyticsData()
	}
}

// Function to run migration (add fake invoices data)
func runMigration() {
	conn := connectToClickHouse()
	defer conn.Close()

	// Create the invoices table
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
		// Generate a date for the invoice, distributing it across months.
		year := 2024
		month := (i % 12) + 1 // This will give a month between 1 and 12
		day := (i % 28) + 1   // This will ensure a valid day between 1 and 28
		invoiceDate := fmt.Sprintf("%d-%02d-%02d", year, month, day)

		// Insert query with the generated date
		insertQuery := fmt.Sprintf("INSERT INTO invoices VALUES (%d, %.2f, '%s')", i, float32(i)*100, invoiceDate)

		err := conn.Exec(context.Background(), insertQuery)
		if err != nil {
			log.Fatalf("Error inserting data: %v", err)
		}
	}

	log.Println("Inserted 100 fake invoices.")
}

// Function to fetch analytics data
func fetchAnalyticsData() {
	conn := connectToClickHouse()
	defer conn.Close()

	// Fetch the most revenue-generating month
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
