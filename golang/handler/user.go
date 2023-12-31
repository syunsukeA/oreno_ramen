package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"
	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"

	"github.com/gin-gonic/gin"
)

type HUser struct {
	Ur repository.User
	Rr repository.Review
}

type UserProfileResponse struct {
	User            *object.User     `json:"user"`
	LatestReviews   []*object.Review `json:"latestReviews"`
	BookmarkReviews []*object.Review `json:"bookmarkReviews"`
}

func (h *HUser) UserProfile(c *gin.Context) {
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

	// ユーザー情報取得
	uo, err := h.Ur.FindByUsername(c, uo.UserName)
	if err != nil {
		log.Printf("Internal server err")
		w.WriteHeader(http.StatusInternalServerError)
		c.Abort()
		return
	}

	// userIDから関連する最新レビューを取得
	latestReviews, err := h.Rr.GetLatestReviewByUserID(c, uo.UserID, 3)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// レビューがない場合，404を返す
	if len(latestReviews) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// userIDから関連するブックマーク付きレビューを取得
	bookmarkReviews, err := h.Rr.GetBookmarkReviewByUserID(c, uo.UserID, 3)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// レビューがない場合，404を返す
	if len(bookmarkReviews) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// ユーザー情報とレビューを結合したレスポンスを生成
	response := UserProfileResponse{
		User:            uo,
		LatestReviews:   latestReviews,
		BookmarkReviews: bookmarkReviews,
	}

	// ResponseBodyに書き込み
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Type", "charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
