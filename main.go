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

type Convert struct {
  Body   string `form:"body" json:"body" binding:"required"`
  Format string `form:"format" json:"format" binding:"required"`
}

func main() {
  port := os.Getenv("PORT")

  if port == "" {
    log.Fatal("$PORT must be set")
  }


  router := gin.Default()

  router.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H { "message": "pong" })
  })

  router.POST("/convert", func(c *gin.Context) {
    var json Convert

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
  })

  router.Run(":" + port)
}
