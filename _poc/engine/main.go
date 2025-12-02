package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	driver "github.com/arangodb/go-driver"
	driverhttp "github.com/arangodb/go-driver/http"
)

// CheckDoc defines a dynamic check for the engine to execute
type CheckDoc struct {
	Key         string `json:"_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AQL         string `json:"aql"`
}

func main() {
	fmt.Println("Starting Engine service...")

	remediationServiceURL := os.Getenv("REMEDIATION_SERVICE_URL")
	if remediationServiceURL == "" {
		log.Println("REMEDIATION_SERVICE_URL not set. Will not call remediation service.")
	}

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
		fmt.Println("Engine running... executing dynamic checks.")

		// 1. Fetch all checks from the 'Checks' collection
		checksQuery := "FOR check IN Checks RETURN check"
		cursor, err := db.Query(context.Background(), checksQuery, nil)
		if err != nil {
			log.Printf("Failed to fetch checks: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// 2. Iterate through each check and execute its AQL query
		for {
			var check CheckDoc
			_, err := cursor.ReadDocument(context.Background(), &check)
			if driver.IsNoMoreDocuments(err) {
				break // No more checks to process
			}
			if err != nil {
				log.Printf("Failed to read check document: %v", err)
				continue
			}

			fmt.Printf("Executing check: %s\n", check.Name)

			// 3. Execute the dynamic AQL query from the check
			findingsCursor, err := db.Query(context.Background(), check.AQL, nil)
			if err != nil {
				log.Printf("Failed to execute dynamic query for check '%s': %v", check.Name, err)
				continue
			}

			// 4. Process the findings from the dynamic query
			for {
				var finding interface{}
				_, err := findingsCursor.ReadDocument(context.Background(), &finding)
				if driver.IsNoMoreDocuments(err) {
					break // No more findings for this check
				}
				if err != nil {
					log.Printf("Failed to read finding document for check '%s': %v", check.Name, err)
					continue
				}
				// Log the finding from the executed check
				log.Printf("[FINDING] Check '%s': %v", check.Name, finding)

				// 5. Call the remediation service
				if remediationServiceURL != "" {
					findingBytes, err := json.Marshal(finding)
					if err != nil {
						log.Printf("Failed to marshal finding to JSON: %v", err)
						continue
					}

					resp, err := http.Post(remediationServiceURL, "application/json", bytes.NewBuffer(findingBytes))
					if err != nil {
						log.Printf("Failed to call remediation service: %v", err)
						continue
					}
					defer resp.Body.Close()

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Printf("Failed to read remediation response: %v", err)
						continue
					}

					log.Printf("[REMEDIATION] Remediation for check '%s': %s", check.Name, string(body))
				}
			}
			findingsCursor.Close() // Close the cursor for the findings query
		}
		cursor.Close() // Close the cursor for the checks query

		time.Sleep(10 * time.Second)
	}
}
