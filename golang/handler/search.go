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
	// HotPepper APIのresoponseをGo構造体に変換
	hpResp := new(object.HPResponse)
	json.Unmarshal(byteData, &hpResp)
	hpShops := hpResp.Result.Shop
	// 店舗情報がなかったら404を返す
	if len(hpShops) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	searchedIDs := []string{}
	for _, shop := range hpShops {
		searchedIDs = append(searchedIDs, shop.ID)
	}
	// 検索結果+DB情報からvisitedを消去
	// ToDo: repository.shopに、外部APIから検索したshop_idを受け取ってunvisitedなshop_idを返す関数を宣言
	visitedIDs, err := h.Sr.GetVisitedShopIDs(c, searchedIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// ToDo: visitedIDsの中にある店舗情報を削除する処理の実装
	// ToDo: もう少し簡潔に書けないものか...
	unvisitedhpShops := []*object.HPShop{}
	for _, shop := range hpShops {
		unvisited := true
		for _, visitedID := range visitedIDs {
			if shop.ID == visitedID {
				unvisited = false
			}
		}
		if unvisited {
			unvisitedhpShops = append(unvisitedhpShops, shop)
		}
	}
	
	// 絞り込んだ結果、店舗数が0になってしまったら404を返す
	if len(unvisitedhpShops) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// ResponseBodyに書き込み
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Type", "charset=utf-8")
	if err := json.NewEncoder(w).Encode(hpShops); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}