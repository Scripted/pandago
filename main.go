package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func main() {
	router := gin.Default()
	router.GET("/", ping)
	router.POST("/convert", convert)
	router.Run()
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "OK 🐼, Go!"})
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

func convert(c *gin.Context) {
	var err error

	outFile, err := ioutil.TempFile("", "converted_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(outFile.Name())

	payload, _, err := c.Request.FormFile("payload")
	if err != nil {
		log.Fatal(err)
	}

	args := []string{"-f", c.PostForm("from"), "-t", c.PostForm("to"), "-o", outFile.Name()}

	cmd := exec.Command("pandoc", args...)
	cmd.Stdin = payload
	err = cmd.Run()
	if err != nil {
		log.Panic(err)
	}

	data, err := ioutil.ReadAll(outFile)
	if err != nil {
		log.Panic(err)
	}
	c.Render(200, render.Data{ContentType: contentType(c.PostForm("to")), Data: data})
}
