package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func main() {
	router := gin.Default()
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.GET("/ping", ping)
	router.POST("/convert", convert)
	router.Run()
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "OK 🐼, Go!"})
}

func convert(c *gin.Context) {
	payload, _, err := c.Request.FormFile("payload")
	if err != nil {
		log.Panic(err)
	}

	inputFile := createTempFile("source_")
	defer os.Remove(inputFile.Name())
	io.Copy(inputFile, payload)

	outputFile := createTempFile("converted_")
	defer os.Remove(outputFile.Name())

	args := []string{
		"--standalone",
		"--from", c.PostForm("from"),
		"--to", c.PostForm("to"),
		"--output", outputFile.Name(),
		inputFile.Name(),
	}

	err = exec.Command("pandoc", args...).Run()
	if err != nil {
		log.Panic(err)
	}

	data, err := ioutil.ReadAll(outputFile)
	if err != nil {
		log.Panic(err)
	}
	c.Render(200, render.Data{ContentType: contentType(c.PostForm("to")), Data: data})
}

func createTempFile(prefix string) *os.File {
	tempFile, err := ioutil.TempFile("", prefix)
	if err != nil {
		log.Fatal(err)
	}
	return tempFile
}

func contentType(format string) string {
	switch format {
	case "markdown":
		return "text/markdown; charset=UTF-8"
	case "html":
		return "text/html; charset=utf-8"
	case "docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		log.Panic("Unsupported format")
	}
	return ""
}
