package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ucaptcha/backend-go/challenge"
	"github.com/ucaptcha/backend-go/config"
	"github.com/ucaptcha/backend-go/types"
)

type ChallengeResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	G       string `json:"g"`
	N       string `json:"n"`
	T       int64  `json:"t"`
}

type VerifyRequest struct {
	Y string `json:"y"`
}

type ChallengeRequest struct {
	Difficulty *int64 `json:"difficulty,omitempty"`
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/challenge", createChallengeHandler)
	r.POST("/challenge/:id/validation", verifyChallengeHandler)
	r.PUT("/difficulty", updateDifficultyHandler)

	return r
}

func createChallengeHandler(c *gin.Context) {
	var req ChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ch, err := challenge.NewChallenge()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, ChallengeResponse{
			Success: true,
			ID:      ch.ID,
			G:       ch.G.String(),
			N:       ch.N.String(),
			T:       ch.T,
		})
		return
	}

	var ch *types.Challenge
	var err error

	if req.Difficulty != nil {
		ch, err = challenge.NewChallenge(*req.Difficulty)
	} else {
		ch, err = challenge.NewChallenge()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ChallengeResponse{
		Success: true,
		ID:      ch.ID,
		G:       ch.G.String(),
		N:       ch.N.String(),
		T:       ch.T,
	})
}

func verifyChallengeHandler(c *gin.Context) {
	id := c.Param("id")

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}

	result, err := challenge.VerifyChallenge(id, req.Y)
	if err != nil {
		switch result {
		case 2:
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		case 3:
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		case 4:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		}
		return
	}
	switch result {
	case 1:
		c.JSON(http.StatusOK, gin.H{"success": true})
	case 0:
		c.JSON(http.StatusUnauthorized, gin.H{"success": false})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Unknown error"})
	}
}

type DifficultyRequest struct {
	Difficulty int64 `json:"difficulty"`
}

func updateDifficultyHandler(c *gin.Context) {
	var req DifficultyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}

	config.GlobalConfig.Difficulty = req.Difficulty
	c.JSON(http.StatusOK, gin.H{"success": true, "difficulty": req.Difficulty})
}
