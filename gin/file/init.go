package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env"

	. "web-layout/utils/aliyun/oss"
)

var client *Client

func init() {
	os.Setenv("ALIYUN_AK", "xxx")
	os.Setenv("ALIYUN_SK", "xxx")
	os.Setenv("ALIYUN_ENDPOINT", "oss-cn-shanghai.aliyuncs.com")
	os.Setenv("ALIYUN_BUCKET", "xxx")

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	cli, err := NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	client = cli
}

func upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	filename := header.Filename

	data := make([]byte, 10000)
	_, err = file.Read(data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": "read error",
		})
	}

	if err := client.PutObject(filename, string(data)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "put to aliyun oss fail",
			"msg":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func download(c *gin.Context) {
	content := "hello world, 我是一个文件"

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=hello.txt")
	c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Writer.Write([]byte(content))
}

func main() {
	router := gin.Default()

	router.POST("/upload", upload)
	router.GET("/download", download)

	router.Run(":8080")
}
