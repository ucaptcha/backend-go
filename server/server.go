package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ucaptcha/backend-go/challenge"
)

type ChallengeResponse struct {
	ID string `json:"id"`
	G  string `json:"g"`
	N  string `json:"n"`
	T  int64  `json:"t"`
}

type VerifyRequest struct {
	Y string `json:"y"`
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/challenge", createChallengeHandler)
	r.POST("/challenge/:id/validation", verifyChallengeHandler)

	return r
}

func createChallengeHandler(c *gin.Context) {
	ch, err := challenge.NewChallenge()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ChallengeResponse{
		ID: ch.ID,
		G:  ch.G.String(),
		N:  ch.N.String(),
		T:  ch.T,
	})
}

func verifyChallengeHandler(c *gin.Context) {
	id := c.Param("id")

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result := challenge.VerifyChallenge(id, req.Y)
	switch result {
	case 1:
		c.JSON(http.StatusOK, gin.H{"success": true})
	case 0:
		c.JSON(http.StatusUnauthorized, gin.H{"success": false})
	case 2:
		c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
	}
}
