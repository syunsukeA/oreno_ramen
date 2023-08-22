package handler

import (
	"encoding/json"
	"io"
	"log"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"

	"github.com/gin-gonic/gin"
)

type HReview struct {
	Ur repository.User
	Rr repository.Review
}

func (h *HReview) HomeReview(c *gin.Context) {
	w := c.Writer

	// リクエストボディからオフセットを取得
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.Writer.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	// ctxから認証済みuser情報を取得
	authedUo, exists := c.Get("authedUo")
	// authedUo情報がない場合はAuthミドルウェアでHTTPRessponse返しているはずなのでexists==falseはありえないが念の為チェック
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Something wrong in Auth-Middleware")
		return
	}
	uo, _ := authedUo.(*object.User)

	// userIDから関連するレビューを取得
	number_reviews := 10 + offset
	reviews, err := h.Rr.GetLatestReviewByUserID(c, uo.UserID, int64(number_reviews))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// レビューがない場合
	if len(reviews) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// ResponseBodyに書き込み
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Type", "charset=utf-8")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
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

	req := new(object.CreateReviewRequest)

	// formから値を取得
	var err error
	req.ShopID = r.FormValue("shop_id")
	req.DishName = r.FormValue("dishname")
	req.Content = r.FormValue("content")
	uint64_eval, err := strconv.ParseUint(r.FormValue("evaluate"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	req.Evaluate = uint(uint64_eval)

	log.Println(req)
	// ctxからimg_urlを取得
	filename, exists := c.Get("img_url")
	// imgURLがない場合はerr
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Something wrong on img-uploading")
		return
	}
	// 画像取得エンドポイントのURLを格納
	/*
	ToDo: URL設計についてもう少し考える。
		・サーバーのフォルダ構造をそのままURLにする？
			・安全性に難あり？
		・img/<filename>にする？
			・現状この形
		・etc...
	*/
	req.ReviewImg = fmt.Sprintf("img/%s", filename.(string))

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
	log.Println("req: ", req)
	ro, err := h.Rr.AddReviewAndShop(c, req.ShopID, uo.UserID, hpShops[0].Name, req)
	log.Println("ro: ", ro)
	if err != nil || ro == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
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
	if err != nil {
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
