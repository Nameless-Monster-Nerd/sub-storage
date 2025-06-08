package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Dsn string
var MinioClient *minio.Client
var BucketName string
var Env string
var RabbitKey string

func init() {
	loadEnv()

	Env = os.Getenv("ENV")
	RabbitKey = os.Getenv("RABBIT_KEY")


	Dsn = os.Getenv("POSTGRES_DSN")
	BucketName = os.Getenv("MINIO_BUCKET")


	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := strings.ToLower(os.Getenv("MINIO_USE_SSL")) == "true"

	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("failed to initialize MinIO client: %v", err)
	}
}

func loadEnv() {
	file, err := os.ReadFile(".env")
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		os.Setenv(key, val)
	}
}

func PrefixGenerator(id string, ss *string, ep *string, isFlixhq bool) string {
	flix := "noFlix"
	type_ := "movie"
	if isFlixhq {
		flix = "flix"
	}
	objectName := fmt.Sprintf("%s/%s/%s", type_, flix, id)
	if ss != nil {
		type_ = "tv"
		objectName = fmt.Sprintf("%s/%s/%s/%s/%s", type_, flix, id, *ss, *ep)
	}
	return objectName
}
