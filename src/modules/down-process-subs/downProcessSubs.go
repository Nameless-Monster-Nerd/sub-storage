package downprocesssubs

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	utils "github.com/nameless-Monster-Nerd/subtitle/src/modules"
)

// DownloadSub loads and decompresses a GZIP subtitle file from MinIO and returns its content as string
func DownProcessSubs(bucketName, objectName string) (string, error) {
	// Get the object from MinIO
	obj, err := utils.MinioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get object: %w", err)
	}
	defer obj.Close()

	// Wrap object reader in gzip reader
	gzr, err := gzip.NewReader(obj)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	// Read decompressed content
	var sb strings.Builder
	_, err = io.Copy(&sb, gzr)
	if err != nil {
		return "", fmt.Errorf("failed to read decompressed content: %w", err)
	}

	return sb.String(), nil
}
