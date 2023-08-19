package handler

import (
	"io"
	"log"
	"strconv"
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

func (h *HReview) CreateReview(c *gin.Context) {
	r := c.Request
	w := c.Writer

	// ctxから認証済みuser情報を取得
	authedUo, exists := c.Get("authedUo")
	// authedUo情報がない場合はAuthミドルウェアでHTTPRessponse返しているはずなのでexists==falseはありえないが念の為チェック
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Something wrong in Auth-Middleware")
		return
	}
	uo, _ := authedUo.(*object.User)

	// reqBodyからreview情報取得
	req := new(object.CreateReviewRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	// HotPepper APIで shop_idが存在するか判定する
	params := url.Values{}
    params.Add("key", HP_API_KEY)
	params.Add("id", req.ShopID)
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
	// reviewを追加 (shopになければshopにも追加)
	// ToDo: 引数多いのでどうにかできたらしたい
	ro, err := h.Rr.AddReviewAndShop(c, req.ShopID, uo.UserID, hpShops[0].Name, req)
	if err != nil || ro == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// ToDo: 必要そうならBodyにreviewの情報を格納して返す
	w.WriteHeader(http.StatusNoContent)
}

func (h *HReview) UpdateReview(c *gin.Context) {
	r := c.Request
	w := c.Writer

	// ctxから認証済みuser情報を取得
	authedUo, exists := c.Get("authedUo")
	// authedUo情報がない場合はAuthミドルウェアでHTTPRessponse返しているはずなのでexists==falseはありえないが念の為チェック
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Something wrong in Auth-Middleware")
		return
	}
	uo, _ := authedUo.(*object.User)
	// reqBodyからreview情報取得
	ro := new(object.Review)
	if err := json.NewDecoder(r.Body).Decode(ro); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	// 認可の確認
	if uo.UserID != ro.UserID {
		w.WriteHeader(http.StatusForbidden)
		log.Println(http.StatusForbidden)
		return
	}
	
	// reviewの修正
	var err error
	ro, err = h.Rr.UpdateReview(c, ro)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if ro == nil {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	// ResponseBodyに書き込み
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Type", "charset=utf-8")
	if err := json.NewEncoder(w).Encode(ro); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (h *HReview) RemoveReview(c *gin.Context) {
	// r := c.Request
	w := c.Writer

	// ctxから認証済みuser情報を取得
	authedUo, exists := c.Get("authedUo")
	// authedUo情報がない場合はAuthミドルウェアでHTTPRessponse返しているはずなのでexists==falseはありえないが念の為チェック
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Something wrong in Auth-Middleware")
		return
	}
	uo, _ := authedUo.(*object.User)
	// URLからreview_idを取得
	reviewID, err := strconv.ParseInt(c.Param("review_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	
	// review_idからobject.Reviewを検索
	ro, err := h.Rr.FindByReviewID(c, reviewID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if ro == nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	// 認可判定
	if ro.UserID != uo.UserID {
		w.WriteHeader(http.StatusForbidden)
		log.Println(http.StatusForbidden)
		return
	}

	// reviewを削除 (reviewが0になったらshopsからも削除)
	ro, err = h.Rr.RemoveReviewAndShop(c, ro.UserID, ro.ShopID, ro.ReviewID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if ro == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}