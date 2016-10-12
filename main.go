package main

import (
	"fmt"
	"os"
	"os/exec"

	"bytes"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" } // default for development

	router := gin.Default()

	router.GET("/", ping)
	router.POST("/convert", convert)

	router.Run(":" + port)
}

type ConvertParams struct {
	Body string `form:"body" json:"body" binding:"required"`
	Format string `form:"format" json:"format" binding:"required"`
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{ "message": "OK 🐼, Go!" })
}

func convert(c *gin.Context) {
	var json ConvertParams

	if c.BindJSON(&json) == nil {
		args := []string{"-f", "html", "-t", json.Format}

		cmd := exec.Command("pandoc", args...)
		cmd.Stdin = strings.NewReader(json.Body)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		if err != nil {
			c.JSON(500, gin.H { "error": "pandoc died for your sins 😕" })
			log.Fatal(err)
		}

		fmt.Printf(out.String())

		c.JSON(200, gin.H { "format": json.Format, "body": out.String() })
	} else {
		c.JSON(418, gin.H { "error": "body and format required" })
	}
}
