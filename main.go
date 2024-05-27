package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Peer struct {
	Name   string `json:"name"`
	PeerID string `json:"peer_id"`
}

type PassThruRequest struct {
	PeerID  string `json:"peer_id"`
	URI     string `json:"uri"`
	Payload string `json:"payload"`
}

func main() {
	r := gin.Default()

	r.GET("/v1/p2p/peers", func(c *gin.Context) {
		peers := []Peer{
			{
				// JPM Chase
				Name:   "984500653R409CC5AB28",
				PeerID: "hash1",
			},
			{
				// MAS
				Name:   "54930035WQZLGC45RZ35",
				PeerID: "hash2",
			},
			{
				// HLB
				Name:   "549300BUPYUQGB5BFX94",
				PeerID: "hash3",
			},
			{
				// BNM
				Name:   "549300NROGNBV2T1GS07",
				PeerID: "hash3",
			},
		}
		c.JSON(http.StatusOK, peers)
	})

	r.POST("/v1/p2p/passthru", func(c *gin.Context) {
		var req PassThruRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error() + " 1"})
			return
		}

		httpReq, err := http.NewRequest("POST", req.URI, bytes.NewBuffer([]byte(req.Payload)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error() + " 2"})
			return
		}

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error() + " 3"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error() + " 4"})
			return
		}

		c.Data(resp.StatusCode, "application/json", body)
	})

	err := r.Run(":5000") // Run on port 5000
	if err != nil {
		panic(err)
	}
}
