package handler

import (
	"io"
	"os"
	"log"
	"fmt"
	"time"
	"net/http"
	"path/filepath"

	"github.com/syunsukeA/oreno_ramen/golang/domain/repository"

	"github.com/gin-gonic/gin"
)

const (
	img_dir_path = ".data/img"
)

type HImg struct {
	Ur repository.User
}

func (h *HImg)ImgHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		w := c.Writer
		// ToDo: 最大サイズ指定方法の調査
		err := r.ParseMultipartForm(10 << 20) // 最大10MBのファイルサイズを許容
		if err != nil {
			log.Println("Unable to parse form data")
			w.WriteHeader(http.StatusBadRequest)
			c.Abort()
			return
		}
		// 1. Fromのファイル情報を取得
		file, header, err := r.FormFile("review_img")
		if err != nil {
			log.Println("Error retrieving the file")
			w.WriteHeader(http.StatusBadRequest)
			c.Abort()
			return
		}
		defer file.Close()
		// 2. 保存用のディレクトリを作成する（存在していなければ、保存用のディレクトリを新規作成）
		err = os.MkdirAll(img_dir_path, os.ModePerm)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			c.Abort()
			return
		}
		// 3. 保存するファイルを作成する
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))
		filepath := fmt.Sprintf("%s/%s", img_dir_path, filename)
		dst, err := os.Create(filepath)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			c.Abort()
			return
		}
		defer dst.Close()
		// アップロードしたファイルを保存用のディレクトリにコピーする
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("upload OK")
		c.Set("img_url", filename)
		// c.Abort()
	}
}

func (h *HImg) ShowImg(c *gin.Context) {
	w := c.Writer

	// URLからパラメータ取得
	filename := c.Param("filename")

	// ファイルシステムから画像データを取得
	imgFilePath := fmt.Sprintf("%s/%s", img_dir_path, filename)
	imgFile, err := os.Open(imgFilePath)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer imgFile.Close()

	// 画像データをByteスライスに変換
	imgData, err := io.ReadAll(imgFile)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	
	// ToDo: RespBodyにバイナリ？画像データの格納
	// ResponseBodyに書き込み
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Add("Content-Type", "charset=utf-8")
	c.Data(http.StatusOK, "image/jpeg", imgData)
}