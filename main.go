package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "github.com/antonholmquist/jason"
)

func main() {
	url := "https://get.geojs.io/v1/ip/geo.json"

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
  	str := string(byteArray)

	v, err := jason.NewObjectFromBytes([]byte(str))
	if err != nil{
		panic(err)
	}

	latitude, _ := v.GetString("latitude")
	longitude, _ := v.GetString("longitude")

	fmt.Println("latitude:" + latitude)
	fmt.Println("longitude:" + longitude)
}