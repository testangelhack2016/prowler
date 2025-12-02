package main

import (
	"context"
	"fmt"
	"log"
	"os"

	driver "github.com/arangodb/go-driver"
	driverhttp "github.com/arangodb/go-driver/http"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	fmt.Println("Starting Ingestor service...")

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

	// Create database and collections
	db, err := client.Database(nil, "_system")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// collections
	collections := []string{"S3Bucket", "Tag"}
	for _, col := range collections {
		options := &driver.CreateCollectionOptions{}
		_, err := db.CreateCollection(nil, col, options)
		if err != nil && !driver.IsConflict(err) {
			log.Fatalf("Failed to create collection %s: %v", col, err)
		}
	}

	// edge collections
	edgeCollections := []string{"has_tag"}
	for _, col := range edgeCollections {
		options := &driver.CreateCollectionOptions{Type: driver.CollectionTypeEdge}
		_, err := db.CreateCollection(nil, col, options)
		if err != nil && !driver.IsConflict(err) {
			log.Fatalf("Failed to create edge collection %s: %v", col, err)
		}
	}

	// Custom resolver for LocalStack
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpointURL := os.Getenv("AWS_ENDPOINT_URL"); awsEndpointURL != "" {
			return aws.Endpoint{
				URL:           awsEndpointURL,
				SigningRegion: region,
				Source:        aws.EndpointSourceCustom,
			}, nil
		}
		// fallback to default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// List S3 buckets
	result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("failed to list buckets, %v", err)
	}

	for _, bucket := range result.Buckets {
		fmt.Printf("bucket: %s, creation date: %s\n", *bucket.Name, *bucket.CreationDate)
	}
}
