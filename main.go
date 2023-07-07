package main

import (
  "fmt"
  "github.com/syunsukeA/oreno_ramen/golang/internal"
)

func main() {
	// hotpepper APIで店を取得
	shop_data := internal.getShopList("b6507930d6c151bd")
	fmt.Println(shop_data["results"])
	
}