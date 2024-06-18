package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// client message
type P2PClientMessage struct {
	PeerID    string `json:"peer_id"`
	MessageID int    `json:"message_id"`
	Method    string `json:"method"`
	Payload   []byte `json:"payload"`
}

// server message
type P2PServerMessasge struct {
	MessageID int    `json:"message_id"`
	Method    string `json:"method"`
	Payload   []byte `json:"payload"`
}

func main() {
	var peers map[string]string = make(map[string]string)

	for i := 1; i <= 4; i++ {
		backendURI := os.Getenv("BACKEND_URI" + fmt.Sprint(i))
		peers["hash"+fmt.Sprint(i)] = backendURI
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment must be set")
	}

	r := gin.Default()

	r.GET("/peers", func(c *gin.Context) {
		peers := map[string]string{
			"984500653R409CC5AB28": "hash1",
			"54930035WQZLGC45RZ35": "hash2",
			"549300BUPYUQGB5BFX94": "hash3",
			"549300NROGNBV2T1GS07": "hash4",
		}

		c.JSON(http.StatusOK, peers)
	})

	r.POST("/passthrough", func(c *gin.Context) {
		var req P2PClientMessage
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		passThruMessage := P2PServerMessasge{
			MessageID: req.MessageID,
			Method:    req.Method,
			Payload:   req.Payload,
		}

		jsonpassThruMessage, err := json.Marshal(passThruMessage)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to marshal pass thru message object " + err.Error()})
			return
		}

		httpReq, err := http.NewRequest("POST", peers[req.PeerID], bytes.NewBuffer(jsonpassThruMessage))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create http request " + err.Error()})
			return
		}

		go func() {
			client := &http.Client{}
			resp, err := client.Do(httpReq)
			if err != nil {
				fmt.Println("failed to send http request " + err.Error())
				return
			}
			defer resp.Body.Close()

			_, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("failed to read response " + err.Error())
				return
			}
		}()

		c.Data(http.StatusOK, "application/json", []byte{})
	})

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
