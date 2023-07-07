package main

import (
  "fmt"
  "github.com/syunsukeA/oreno_ramen/golang/internal"
)

func main() {
	// hotpepper APIで店を取得
	var shop_data map[string]interface{}
	shop_data = internal.GetShopList("b6507930d6c151bd")
	fmt.Println(shop_data["results"])
	
}