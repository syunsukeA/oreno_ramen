package main

import (
  _"fmt"
  "github.com/syunsukeA/oreno_ramen/golang/internal"
  "github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", internal.GetShoplist)
    router.Run("localhost:8080")
}