package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Peer struct {
	Name   string `json:"name"`
	PeerID string `json:"peer_id"`
}

type PassThruRequest struct {
	PeerID    string `json:"peer_id"`
	MessageID int    `json:"message_id"`
	Method    string `json:"uri"`
	Payload   []byte `json:"payload"`
}

type PassThruMessage struct {
	MessageID int    `json:"message_id"`
	Method    string `json:"uri"`
	Payload   []byte `json:"payload"`
}

func main() {
	backendURI := os.Getenv("BACKEND_URI")
	if backendURI == "" {
		log.Fatal("BACKEND_URI environment must be set")
	}

	port := os.Getenv("PORT")
	if backendURI == "" {
		log.Fatal("PORT environment must be set")
	}

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
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		passThruMessage := PassThruMessage{
			MessageID: req.MessageID,
			Method:    req.Method,
			Payload:   req.Payload,
		}

		jsonpassThruMessage, err := json.Marshal(passThruMessage)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to marshal pass thru message object " + err.Error()})
			return
		}

		httpReq, err := http.NewRequest("POST", backendURI, bytes.NewBuffer(jsonpassThruMessage))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create http request " + err.Error()})
			return
		}

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to send http request " + err.Error()})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read response " + err.Error()})
			return
		}

		c.Data(resp.StatusCode, "application/json", body)
	})

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
