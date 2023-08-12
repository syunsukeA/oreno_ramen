package handler

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"encoding/json"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"

	"github.com/gin-gonic/gin"
)

type HReview struct {
	Ur repository.User
	Rr repository.Review
}

type CreateReviewRequest struct {
	Content string
	Evaluate uint
	ReviewImg string // ToDo画像の取り扱いについては後ほど実装

}

func (h *HReview) CreateReview(c *gin.Context) {
	r := c.Request
	w := c.Writer
	// reqBodyからreview情報取得
	req := new(CreateReviewRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	// URLからusername、shop_idを取得
	username := c.Param("username")
	shop_id := c.Param("shop_id")
	// HotPepper APIで shop_idが存在するか判定する
	params := url.Values{}
    params.Add("key", HP_API_KEY)
	params.Add("id", shop_id)
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
	// shop_idに対応する店舗がなかったらparamater指定ミスなので400を返す (uri情報なので404かも...？)
	if len(hpShops) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	log.Println(shop_id)
	user, err := h.Ur.FindByUsername(c, username) // ToDo: usernameからUser情報を検索して返すDBコマンドの作成
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
	// reviewを追加 (shopになければshopにも追加)
	// ToDo: 引数多いのでどうにかできたらしたい
	ro, err := h.Rr.AddReviewAndShop(c, shop_id, user.UserID, hpShops[0].Name, req.Content, req.Evaluate, req.ReviewImg)
	if err != nil || ro == nil {
		log.Println("============================")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// ToDo: 必要そうならBodyにreviewの情報を格納して返す
	w.WriteHeader(http.StatusNoContent)
}