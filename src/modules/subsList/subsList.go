package subslist

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	utils "github.com/nameless-Monster-Nerd/subtitle/src/modules"
)

type SubtitleObject struct {
	Key      string
	Metadata map[string]string
}
func Subslist(id string, ss *string, ep *string) ([]SubtitleObject, error) {
	prefix := utils.PrefixGenerator(id, ss, ep, true)

	var files []SubtitleObject
	ctx := context.Background()

	objectCh := utils.MinioClient.ListObjects(ctx, utils.BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			fmt.Println("Error:", object.Err)
			continue
		}

		// Get metadata for this object
		stat, err := utils.MinioClient.StatObject(ctx, utils.BucketName, object.Key, minio.StatObjectOptions{})
		if err != nil {
			fmt.Println("StatObject error:", err)
			continue
		}

		files = append(files, SubtitleObject{
			Key:      object.Key,
			Metadata: stat.UserMetadata,
		})
	}

	return files, nil
}
