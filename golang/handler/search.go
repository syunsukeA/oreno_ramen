package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"

	"github.com/gin-gonic/gin"
)

const (
	HP_API_KEY = "b6507930d6c151bd"
)

// [!!注意!!] repositoryを追加したらmain.goに反映すること！！
type HSearch struct {
	Sr repository.Shop
	Ur repository.User
	Rr repository.Review
}

/*
Visitedが一番意見が食い違っていそう。
とりあえず一番低コストな"Unvisitedのコピー"として実装。
*/
func (h *HSearch) SearchVisited(c *gin.Context) {
	r := c.Request
	w := c.Writer
	// URLからusernameを取得
	username := c.Param("username")
	// usernameからUser情報を検索
	user, err := h.Ur.FindByUsername(c, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	// queryパラメータ判定
	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")
	rng := r.URL.Query().Get("rng")
	// ToDo: paramエラーハンドリング
	if lat == "" || lng == "" || rng == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
	// 外部APIから検索したshop_idとuser_idを受け取ってunvisitedなshop_idを返す
	visitedIDs, err := h.Sr.GetVisitedShopIDs(c, searchedIDs, user.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	//visitedな店舗情報を配列に格納
	visitedhpShops := []*object.HPShop{}
	for _, shop := range hpShops {
		for _, visitedID := range visitedIDs {
			if shop.ID == visitedID {
				// shopIDからreviewを (現状最大20件) 引っ張ってくる
				log.Println("shop id: ", shop.ID)
				reviews, err := h.Rr.FindReviewsByShopID(c, user.UserID, shop.ID)
				log.Println("esc", reviews)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(err)
					return
				}
				shop.Reviews = reviews
				visitedhpShops = append(visitedhpShops, shop)
			}
		}
	}
	// ResponseBodyに書き込み
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Type", "charset=utf-8")
	if err := json.NewEncoder(w).Encode(visitedhpShops); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (h *HSearch) SearchUnvisited(c *gin.Context) {
	r := c.Request
	w := c.Writer
	// URLからusernameを取得
	username := c.Param("username")
	// usernameからUser情報を検索
	user, err := h.Ur.FindByUsername(c, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	// queryパラメータ判定
	lat := r.URL.Query().Get("lat")
	lng := r.URL.Query().Get("lng")
	rng := r.URL.Query().Get("rng")
	// ToDo: paramエラーハンドリング
	if lat == "" || lng == "" || rng == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
	// 外部APIから検索したshop_idとuser_idを受け取ってunvisitedなshop_idを返す
	visitedIDs, err := h.Sr.GetVisitedShopIDs(c, searchedIDs, user.UserID)
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
	if err := json.NewEncoder(w).Encode(unvisitedhpShops); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
