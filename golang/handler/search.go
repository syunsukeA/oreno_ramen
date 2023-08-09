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
	r := c.Request
	w := c.Writer
	lat := r.URL.Query().Get("lat") 
	lng := r.URL.Query().Get("lng")
	rng := r.URL.Query().Get("rng")
	test_resp := []*object.Review{}
	test_resp = append(test_resp, &object.Review{})

	log.Printf("lat: %s, lng: %s, rng: %s\n", lat, lng, rng)
	w.Header().Set("Content-Type", "text/plain")
	if err := json.NewEncoder(w).Encode(test_resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
}