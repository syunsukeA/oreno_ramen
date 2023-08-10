package handler

import (
	"log"
	"io"
	"net/http"
	"net/url"
	"encoding/json"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"

	"github.com/gin-gonic/gin"
)

const (
	HP_API_KEY = "b6507930d6c151bd"
)

type HSearch struct {
	Sr repository.Shop
	// Rr repository.Review
}

func (h *HSearch) SearchUnvisited(c *gin.Context){
	// queryパラメータ判定
	r := c.Request
	w := c.Writer
	lat := r.URL.Query().Get("lat") 
	lng := r.URL.Query().Get("lng")
	rng := r.URL.Query().Get("rng")
	// ToDo: paramエラーハンドリング

	// HotPepper API呼び出し
	params := url.Values{}
    params.Add("key", HP_API_KEY)
	params.Add("keyword", "ラーメン")
    params.Add("lat", lat)
	params.Add("lng", lng)
	params.Add("range", rng)
	params.Add("format", "json")

	urls := "http://webservice.recruit.co.jp/hotpepper/gourmet/v1/?" + params.Encode()
	resp, err := http.Get(urls)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	hpPesp := new(object.HPResponse)
	json.Unmarshal(byteData, &hpPesp)
	hpShops := hpPesp.Result.Shop

	// 検索結果+DB情報からvisitedを消去
	// ToDo: repository.shopに、外部APIから検索したshop_idを受け取ってunvisitedなshop_idを返す関数を宣言

	// ToDo: unvisited_idの中にない店舗情報を削除するような実装
	for _, shop := range hpShops {
		log.Println(shop.ID)
	}

	// ResponseBodyに書き込み
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hpShops); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}