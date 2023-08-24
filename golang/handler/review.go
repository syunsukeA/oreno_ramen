package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"

	"github.com/gin-gonic/gin"
)

type HReview struct {
	Ur repository.User
	Rr repository.Review
}

func (h *HReview) CreateReview(c *gin.Context) {
	// ミドルウェアのサーバーな内画像トランザクションのためにパニックをStatusCode:500に変換
	defer func() {
		// パニック終了した場合はStatusCodeを500にする
		if r := recover(); r != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
	}()
	
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
	uint64_bookmark, err := strconv.ParseUint(r.FormValue("bookmark"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	req.Bookmark = uint(uint64_bookmark)
	// req.Bookmark = (uint64_bookmark != 0) // boolの場合

	log.Println(req)
	// ctxからimg_urlを取得
	filename, exists := c.Get("imgFilename")
	// imgURLがない場合はerr
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Something wrong on img-uploading")
		return
	}
	// 画像取得エンドポイントのURLを格納
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
	// ミドルウェアのサーバーな内画像トランザクションのためにパニックをStatusCode:500に変換
	defer func() {
		// パニック終了した場合はStatusCodeを500にする
		if r := recover(); r != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
	}()

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

	req := new(object.CreateReviewRequest)
	// formから値を取得
	// ToDo: 変更がない場合は空文字列で送るような実装になっているならば修正が必要
	// req.ShopID = r.FormValue("shop_id")
	ro.DishName = r.FormValue("dishname")
	ro.Content = r.FormValue("content")
	uint64_eval, err := strconv.ParseUint(r.FormValue("evaluate"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	ro.Evaluate = uint(uint64_eval)
	uint64_bookmark, err := strconv.ParseUint(r.FormValue("bookmark"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	ro.Bookmark = uint(uint64_bookmark)
	// ro.Bookmark = (uint64_bookmark != 0) // boolの場合

	log.Println(req)
	// ctxからimg_urlを取得
	filename, exists := c.Get("imgFilename")
	// imgURLがない場合はerr
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Something wrong on img-uploading")
		return
	}
	// 空文字列ではない場合のみ画像取得エンドポイントのURLを上書き
	log.Println("filename: ", filename)
	var deleteFilename string
	if len(filename.(string)) > 0 {
		deleteFilename = ro.ReviewImg[4:] // 'img/' を取り除くためのハードコーディング
		ro.ReviewImg = fmt.Sprintf("img/%s", filename.(string))
	}

	// 認可の確認
	if uo.UserID != ro.UserID {
		w.WriteHeader(http.StatusForbidden)
		log.Println(http.StatusForbidden)
		return
	}

	// reviewの修正
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
	// 最後までerrがない&filenameが空文字列でないなら画像を削除
	if len(deleteFilename) > 0 {
		// 該当ファイルを削除
		removeFilePath := fmt.Sprintf("%s/%s", img_dir_path, deleteFilename)
		err := os.Remove(removeFilePath)
		// errが発生した場合はfilepathをlogに吐くようにする
		// ToDo: 削除に失敗した場合のさらに良い対処法を考える
		if err != nil {
			log.Printf("Remove reeor: '%s' ", removeFilePath)
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
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
	ro, err = h.Rr.RemoveReviewAndShop(c, ro)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// ToDo: 返すべきステータスコードちゃんと考える
	if ro == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	// ToDo: 最後までerrがない&filenameが空文字列でないなら画像を削除
	log.Println("reviewImg: ", ro.ReviewImg)
	deleteFilename := ro.ReviewImg[4:] // 'img/' を取り除くためのハードコーディング
	if len(deleteFilename) > 0 {
		// 該当ファイルを削除
		removeFilePath := fmt.Sprintf("%s/%s", img_dir_path, deleteFilename)
		err := os.Remove(removeFilePath)
		// errが発生した場合はfilepathをlogに吐くようにする
		// ToDo: 削除に失敗した場合のさらに良い対処法を考える
		if err != nil {
			log.Printf("Remove reeor: '%s' ", removeFilePath)
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	}
}
