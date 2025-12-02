
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	driver "github.com/arangodb/go-driver"
	driverhttp "github.com/arangodb/go-driver/http"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3BucketDoc defines the structure for a bucket document in ArangoDB
type S3BucketDoc struct {
	Key      string `json:"_key"`
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}

// TagDoc defines the structure for a tag document in ArangoDB
type TagDoc struct {
	Key  string `json:"_key"`
	Name string `json:"name"`
}

// EdgeDoc defines the structure for an edge document in ArangoDB
type EdgeDoc struct {
	From string `json:"_from"`
	To   string `json:"_to"`
}

// CheckDoc defines a dynamic check for the engine to execute
type CheckDoc struct {
	Key         string `json:"_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AQL         string `json:"aql"`
}

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

	// Get database
	db, err := client.Database(nil, "_system")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Ensure collections exist
	s3Col, err := ensureCollection(db, "S3Bucket", driver.CollectionTypeDocument)
	if err != nil {
		log.Fatalf("Failed to ensure collection S3Bucket: %v", err)
	}
	tagCol, err := ensureCollection(db, "Tag", driver.CollectionTypeDocument)
	if err != nil {
		log.Fatalf("Failed to ensure collection Tag: %v", err)
	}
	hasTagCol, err := ensureCollection(db, "has_tag", driver.CollectionTypeEdge)
	if err != nil {
		log.Fatalf("Failed to ensure edge collection has_tag: %v", err)
	}
	checksCol, err := ensureCollection(db, "Checks", driver.CollectionTypeDocument)
	if err != nil {
		log.Fatalf("Failed to ensure collection Checks: %v", err)
	}

	// Create and ingest the dynamic check for public sensitive buckets
	createDefaultCheck(db, checksCol)

	// Custom resolver for LocalStack
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpointURL := os.Getenv("AWS_ENDPOINT_URL"); awsEndpointURL != "" {
			return aws.Endpoint{
				URL:           awsEndpointURL,
				SigningRegion: region,
				Source:        aws.EndpointSourceCustom,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create S3 client with path-style addressing
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	fmt.Println("Listing S3 buckets...")
	result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("failed to list buckets, %v", err)
	}

	if len(result.Buckets) == 0 {
		fmt.Println("No buckets found.")
		return
	}

	for _, bucket := range result.Buckets {
		bucketName := *bucket.Name
		fmt.Printf("Checking bucket: %s\n", bucketName)

		isPublic := checkPublicAccess(s3Client, bucketName)
		hasSensitiveTag := checkSensitiveTag(s3Client, bucketName)

		if isPublic && hasSensitiveTag {
			fmt.Printf("Found toxic combination: Public and sensitive bucket '%s'. Ingesting into ArangoDB.\n", bucketName)

			// Create S3 Bucket document
			bucketDoc := S3BucketDoc{
				Key:      bucketName,
				Name:     bucketName,
				IsPublic: true,
			}
			_, err := s3Col.CreateDocument(nil, bucketDoc)
			if err != nil && !driver.IsConflict(err) {
				log.Printf("Failed to create bucket document for %s: %v", bucketName, err)
				continue
			}

			// Create Tag document
			tagKey := "sensitivity:high"
			tagDocKey := strings.Replace(tagKey, ":", "-", -1)
			tagDoc := TagDoc{
				Key:  tagDocKey,
				Name: tagKey,
			}
			_, err = tagCol.CreateDocument(nil, tagDoc)
			if err != nil && !driver.IsConflict(err) {
				log.Printf("Failed to create tag document for %s: %v", tagKey, err)
				continue
			}

			// Create Edge document linking bucket and tag
			edgeDoc := EdgeDoc{
				From: s3Col.Name() + "/" + bucketName,
				To:   tagCol.Name() + "/" + tagDocKey,
			}
			_, err = hasTagCol.CreateDocument(nil, edgeDoc)
			if err != nil && !driver.IsConflict(err) {
				log.Printf("Failed to create edge from %s to %s: %v", bucketName, tagKey, err)
			}

			fmt.Printf("Successfully ingested bucket %s and its sensitive tag.\n", bucketName)

		} else {
			fmt.Printf("Bucket %s is not a toxic combination. Public: %v, Sensitive Tag: %v. Skipping.\n", bucketName, isPublic, hasSensitiveTag)
		}
	}
	fmt.Println("Ingestor service finished.")
}

func ensureCollection(db driver.Database, name string, colType driver.CollectionType) (driver.Collection, error) {
	exists, err := db.CollectionExists(nil, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check if collection %s exists: %w", name, err)
	}
	if exists {
		return db.Collection(nil, name)
	}

	options := &driver.CreateCollectionOptions{Type: colType}
	col, err := db.CreateCollection(nil, name, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection %s: %w", name, err)
	}
	return col, nil
}

func createDefaultCheck(db driver.Database, checksCol driver.Collection) {
	check := CheckDoc{
		Key:         "public-sensitive-s3-bucket",
		Name:        "Public Sensitive S3 Bucket",
		Description: "Finds S3 buckets that are public and have a 'sensitivity:high' tag.",
		AQL: `
FOR bucket IN S3Bucket
    FILTER bucket.is_public == true
    FOR tag, edge IN 1..1 OUTBOUND bucket has_tag
        FILTER tag.name == "sensitivity:high"
        RETURN {
            "resource_id": bucket.name,
            "message": CONCAT("S3 bucket '", bucket.name, "' is public and tagged as sensitive.")
        }
`,
	}

	_, err := checksCol.CreateDocument(nil, check)
	if err != nil && !driver.IsConflict(err) {
		log.Printf("Failed to create default check document: %v", err)
	} else if err == nil {
		fmt.Println("Successfully ingested default check definition.")
	}
}

func checkPublicAccess(s3Client *s3.Client, bucketName string) bool {
	aclOutput, err := s3Client.GetBucketAcl(context.TODO(), &s3.GetBucketAclInput{
		Bucket: &bucketName,
	})
	if err != nil {
		log.Printf("Could not get ACL for bucket %s: %v", bucketName, err)
		return false
	}
	for _, grant := range aclOutput.Grants {
		if grant.Grantee.Type == types.TypeGroup && *grant.Grantee.URI == "http://acs.amazonaws.com/groups/global/AllUsers" {
			fmt.Printf("Bucket %s has a public grant.\n", bucketName)
			return true
		}
	}
	return false
}

func checkSensitiveTag(s3Client *s3.Client, bucketName string) bool {
	taggingOutput, err := s3Client.GetBucketTagging(context.TODO(), &s3.GetBucketTaggingInput{
		Bucket: &bucketName,
	})
	if err != nil {
		// This can happen if the bucket has no tags. For this PoC, we'll just log and return false.
		log.Printf("Could not get tags for bucket %s: %v. Assuming no sensitive tag.", bucketName, err)
		return false
	}
	for _, tag := range taggingOutput.TagSet {
		if *tag.Key == "sensitivity" && *tag.Value == "high" {
			fmt.Printf("Bucket %s has a 'sensitivity:high' tag.\n", bucketName)
			return true
		}
	}
	return false
}
