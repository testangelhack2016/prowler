package main

import (
	"context"
	"fmt"
	"log"
	"time"

	driver "github.com/arangodb/go-driver"
	driverhttp "github.com/arangodb/go-driver/http"
)

func main() {
	fmt.Println("Starting Engine service...")

	// Create ArangoDB client
	conn, err := driverhttp.NewConnection(driverhttp.ConnectionConfig{
		Endpoints: []string{"http://arangodb:8529"},
	})
	if err != nil {
		log.Fatalf("Failed to create HTTP connection: %v", err)
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "prowler"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	db, err := client.Database(nil, "_system")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	for {
		fmt.Println("Engine running... executing query.")

		// Placeholder for a query
		query := "FOR d IN S3Bucket RETURN d"
		cursor, err := db.Query(context.Background(), query, nil)
		if err != nil {
			log.Printf("Failed to execute query: %v", err)
		} else {
			defer cursor.Close()
			for {
				var doc interface{}
				_, err := cursor.ReadDocument(context.Background(), &doc)
				if driver.IsNoMoreDocuments(err) {
					break
				} else if err != nil {
					log.Printf("Failed to read document: %v", err)
				}
				fmt.Printf("Got doc: %v\n", doc)
			}
		}

		time.Sleep(10 * time.Second)
	}
}
