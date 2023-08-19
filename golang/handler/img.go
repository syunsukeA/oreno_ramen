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
	img_path = ".data/img"
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
		err = os.MkdirAll(img_path, os.ModePerm)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			c.Abort()
			return
		}
		// 3. 保存するファイルを作成する
		filepath := fmt.Sprintf("%s/%d%s", img_path, time.Now().UnixNano(), filepath.Ext(header.Filename))
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

		// fmt.Fprintf(w, "アップロード成功！\n")

		// ファイルの情報を表示
		// fmt.Fprintf(w, "File PATH: %s\n", filepath)
		// fmt.Fprintf(w, "File Uploaded: %s\n", header.Filename)
		// fmt.Fprintf(w, "File size: %d\n", header.Size)
		// w.WriteHeader(http.StatusOK)
		log.Printf("upload OK")
		c.Set("img_url", filepath)
		// c.Abort()
	}
}