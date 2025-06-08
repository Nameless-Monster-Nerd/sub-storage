package proxy

import (
	"compress/gzip"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	utils "github.com/nameless-Monster-Nerd/subtitle/src/modules"
)

func Proxy(c *gin.Context) {
	key := c.Query("key")

	obj, err := utils.MinioClient.GetObject(context.Background(), utils.BucketName, key, minio.GetObjectOptions{})
	if err != nil {
		c.String(http.StatusInternalServerError, "error: %v", err)
		return
	}

	gzipReader, err := gzip.NewReader(obj)
	if err != nil {
		c.String(http.StatusInternalServerError, "gzip error: %v", err)
		return
	}
	defer gzipReader.Close()

	c.DataFromReader(http.StatusOK, -1, "text/plain; charset=utf-8", gzipReader, nil)
}
