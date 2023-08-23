package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"

	"github.com/gin-gonic/gin"
)

type HHome struct {
	Ur repository.User
	Rr repository.Review
}

func (h *HHome) HomeReview(c *gin.Context) {
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

	// レビューがない場合，404を返す
	if len(reviews) == 0 {
		w.WriteHeader(http.StatusNotFound)
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

func (h *HHome) HomeBookmarkReview(c *gin.Context) {
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
	reviews, err := h.Rr.GetBookmarkReviewByUserID(c, uo.UserID, int64(number_reviews))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// レビューがない場合，404を返す
	if len(reviews) == 0 {
		w.WriteHeader(http.StatusNotFound)
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

func (h *HHome) HomeEvaluateReview(c *gin.Context) {
	w := c.Writer

	// リクエストボディからオフセットを取得
	upperStr := c.DefaultQuery("upper", "0")
	upper, err := strconv.ParseInt(upperStr, 10, 64)
	if err != nil || upper < 0 {
		c.Writer.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	lowerStr := c.DefaultQuery("lower", "0")
	lower, err := strconv.ParseInt(lowerStr, 10, 64)
	if err != nil || lower < 0 {
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

	// upperとlowerによってフィルターを掛けたレビュー取得
	reviews, err := h.Rr.GetEvaluateReviewByUserID(c, uo.UserID, upper, lower)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// レビューがない場合，404を返す
	if len(reviews) == 0 {
		w.WriteHeader(http.StatusNotFound)
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
