package handler

import (
	"log"
	"net/http"
	"encoding/json"

	"github.com/syunsukeA/oreno_ramen/golang/domain/object"

	"github.com/gin-gonic/gin"
)

type HSearch struct {
	RO object.Review
}

func (h *HSearch) SearchUnvisited(c *gin.Context){
	log.Println("OK!!!!!!!!!!!!!!")
	w := c.Writer
	w.Header().Set("Content-Type", "text/plain")
	if err := json.NewEncoder(w).Encode(h.RO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
}