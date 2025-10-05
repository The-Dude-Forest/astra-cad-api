package community

import (
	"encoding/json"
	"go-auth/internal/database"
	"go-auth/internal/response"
	"go-auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SaveLayout(db *database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var raw map[string]interface{}
		if err := c.ShouldBindJSON(&raw); err != nil {
			response.Respond(c, http.StatusBadRequest, "Error while sharing your layout [bad payload]", nil)
			return
		}

		hub := models.Hub{
			Author: raw["author"].(string),
			Title:  raw["title"].(string),
		}

		// Everything else goes into Structure
		delete(raw, "author")
		delete(raw, "title")
		hub.Structure = models.HubStructure{}
		bytes, _ := json.Marshal(raw)
		json.Unmarshal(bytes, &hub.Structure)

		if err := db.DB.Create(&hub).Error; err != nil {
			response.Respond(c, http.StatusInternalServerError, "could not create the hub [SaveLayout]", err.Error())
			return
		}

		response.Respond(c, http.StatusOK, "Successfully shared your layout!", nil)
	}
}

func GetLayouts(db *database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var hubs []models.Hub
		if err := db.DB.Table("hubs").Select("hubs.*").Scan(&hubs).Error; err != nil {
			response.Respond(c, http.StatusBadRequest, "Error while retrieving community layouts..", nil)
			return
		}

		response.Respond(c, http.StatusOK, "Successfully loaded community layouts", hubs)
	}
}
