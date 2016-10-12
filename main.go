package main

import (
	"bytes"
	"log"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", ping)
	router.POST("/convert", convert)
	router.Run()
}

type ConvertParams struct {
	Body string `json:"body" binding:"required"`
	From string `json:"from" binding:"required"`
	To   string `json:"to"   binding:"required"`
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "OK 🐼, Go!"})
}

func convert(c *gin.Context) {
	var json ConvertParams

	validation_err := c.BindJSON(&json)
	if validation_err != nil {
		log.Println(validation_err)
		c.JSON(400, gin.H{"error": "body, to, and from are required params"})
		return
	}

	args := []string{"-f", json.From, "-t", json.To}

	cmd := exec.Command("pandoc", args...)
	cmd.Stdin = strings.NewReader(json.Body)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd_err := cmd.Run()

	if cmd_err != nil {
		log.Panic(cmd_err)
	}

	c.JSON(200, gin.H{"format": json.To, "body": out.String()})
}
