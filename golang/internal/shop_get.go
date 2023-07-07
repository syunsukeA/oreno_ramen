package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"io"
	"encoding/json"
	"github.com/antonholmquist/jason"
	"github.com/gin-gonic/gin"
  )

func GetShoplist(c *gin.Context) {
	fmt.Println(c)
	var shop_data map[string]interface{}
	shop_data = getShopjson("b6507930d6c151bd")
	fmt.Println(shop_data["results"])
} 

func getPosition() (string, string) {
	urls := "https://get.geojs.io/v1/ip/geo.json"

	resp, _ := http.Get(urls)
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
  	str := string(byteArray)
	// fmt.Printf("str: \n",str)

	v, err := jason.NewObjectFromBytes([]byte(str))
	if err != nil{
		panic(err)
	}

	latitude, _ := v.GetString("latitude")
	longitude, _ := v.GetString("longitude")

	fmt.Println("latitude:" + latitude)
	fmt.Println("longitude:" + longitude)
	return latitude, longitude
}

func getShopjson(api_key string) map[string]interface{} {
	// 現在地の取得
	lat, lng := getPosition()
	// hotpepper APIで店を取得
	params := url.Values{}
    params.Add("key", api_key)
	params.Add("keyword", "ラーメン")
    params.Add("lat", lat)
	params.Add("lng", lng)
	params.Add("range", "4")
	params.Add("format", "json")

    // パラメータ情報を付加したURLを作成
	urls := "http://webservice.recruit.co.jp/hotpepper/gourmet/v1/?" + params.Encode()
	resp, _ := http.Get(urls)
	fmt.Printf("resp: %v\n", resp)
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
  	str := string(byteArray)

	var shop_data map[string]interface{}
	json.Unmarshal([]byte(str), &shop_data)
	return shop_data

	// v, err := jason.NewObjectFromBytes([]byte(str))
	// if err != nil{
	// 	panic(err)
	// }
	// fmt.Println("v: ", v)
	// shop, _ := v.GetString("results")
	// fmt.Println("shop: ", shop)
}
