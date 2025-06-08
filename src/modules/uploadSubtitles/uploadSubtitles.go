package uploadsubtitles

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/minio/minio-go/v7"
	utils "github.com/nameless-Monster-Nerd/subtitle/src/modules"
)

func UploadSubtitle(lang string, url string , id string , ss *string , ep *string , isFlixhq bool ) minio.UploadInfo {
	runes := []rune(url)
ext := ""
if len(runes) >= 3 {
	ext = string(runes[len(runes)-3:])

}

objectName := fmt.Sprintf("%s/%s.%s.gz", utils.PrefixGenerator(id, ss, ep, isFlixhq), lang, ext)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP error:", err)
	}
	defer resp.Body.Close() 

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		
	}

	// 2. Gzip compress
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err = gzipWriter.Write(body)
	if err != nil {
		fmt.Println("Gzip write error:", err)
	}
	err = gzipWriter.Close()
	if err != nil {
		fmt.Println("Gzip close error:", err)
		
	}
	gzippedBody := buf.Bytes()
 
	
	// objectName := type_ + name + ".gz" 
	contentType := "application/gzip"

	// Ensure the bucket exists
	ctx := context.Background()
	exists, errBucketExists := utils.MinioClient.BucketExists(ctx, utils.BucketName)
	if errBucketExists != nil {
		fmt.Println("Bucket check error:", errBucketExists)
	}
	if !exists {
		err = utils.MinioClient.MakeBucket(ctx, utils.BucketName, minio.MakeBucketOptions{
			
		})
		if err != nil {
			fmt.Println("Bucket creation error:", err)
		}
	}

	// Upload
		ssVal := "0"
			if ss != nil {
				ssVal = *ss
			}

			epVal := ""
			if ep != nil {
				epVal = *ep
			}
	uploadInfo, err := utils.MinioClient.PutObject(ctx, utils.BucketName, objectName, bytes.NewReader(gzippedBody), int64(len(gzippedBody)), minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"lang": lang,
			// "url":url,
			"id":id, 
			"ss":ssVal,
			"ep":epVal,
		},
	})
	if err != nil {
		fmt.Println("Upload error:", err)
		
	}

	fmt.Println("Successfully uploaded to MinIO:", uploadInfo)
	return uploadInfo
}
 