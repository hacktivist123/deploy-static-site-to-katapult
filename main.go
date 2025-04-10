package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the S3 client
	endpoint := os.Getenv("KATAPULT_ENDPOINT")
	accessKeyID := os.Getenv("KATAPULT_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("KATAPULT_SECRET_KEY")
	bucketName := os.Getenv("BUCKET_NAME")

	// Remove port number from endpoint for minio client
	endpoint = strings.Split(endpoint, ":")[0]

	// Initialize S3 client
	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	// Check if bucket exists
	exists, err := s3Client.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalf("Error checking bucket existence: %v", err)
	}

	if !exists {
		log.Fatalf("Bucket %s does not exist", bucketName)
	}

	// Upload files from build directory
	buildDir := "../katapult-test/build" // Adjust this path with your static site build directory
	err = filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {

			objectName := strings.TrimPrefix(path, buildDir+"/")

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = s3Client.PutObject(context.Background(), bucketName, objectName, file, info.Size(), minio.PutObjectOptions{
				ContentType: getContentType(objectName),
			})
			if err != nil {
				return err
			}
			fmt.Printf("Uploaded %s\n", objectName)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error uploading files: %v", err)
	}

	fmt.Printf("Deployment complete! Your site is available at: https://%s.%s\n", bucketName, endpoint)
}

func getContentType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".html"):
		return "text/html"
	case strings.HasSuffix(filename, ".css"):
		return "text/css"
	case strings.HasSuffix(filename, ".js"):
		return "application/javascript"
	case strings.HasSuffix(filename, ".png"):
		return "image/png"
	case strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(filename, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(filename, ".json"):
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
