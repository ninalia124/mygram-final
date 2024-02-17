package controllers

import (
	"encoding/json"
	"mygram-final/database"
	"mygram-final/helpers"
	"mygram-final/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// fungsi CreateSosialMedia untuk membuat baru
func CreateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	socialMediaRequest := models.CreateSocialMediaRequest{}
	userID := uint(userData["id"].(float64))
	//untuk memeriksa jenis content type yang digunakan
	if contentType == appJson {
		if err := c.ShouldBindJSON(&socialMediaRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&socialMediaRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	socialMedia := models.SocialMedia{
		Name:           socialMediaRequest.Name,
		SocialMediaUrl: socialMediaRequest.SocialMediaUrl,
		UserId:         userID,
	}

	err := db.Debug().Create(&socialMedia).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	socialMediaString, _ := json.Marshal(socialMedia)
	socialMediaResponse := models.CreateSocialMediaResponse{}
	json.Unmarshal(socialMediaString, &socialMediaResponse)

	c.JSON(http.StatusCreated, socialMediaResponse)
}

// fungsi  GetSocialMedia untuk mendapatkan sosial media yang sudah dicreate
func GetSocialMedia(c *gin.Context) {
	db := database.GetDB()

	socialMedias := []models.SocialMedia{}

	err := db.Debug().Preload("User").Order("id asc").Find(&socialMedias).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	socialMediaString, _ := json.Marshal(socialMedias)
	socialMediaResponse := []models.SocialMediaResponse{}
	json.Unmarshal(socialMediaString, &socialMediaResponse)

	c.JSON(http.StatusOK, socialMediaResponse)
}

// fungsi untuk mendapatkan data sosial media sesuai id
func GetIdSocialMedia(c *gin.Context) {
	db := database.GetDB()
	socialMediaId := c.Param("socialmediaId")
	var allSocialMedia []models.SocialMedia

	err := db.First(&allSocialMedia, "Id = ?", socialMediaId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Record not found",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Data": allSocialMedia,
	})
}

// fungsi untuk mengupdate sosial media existing
func UpdateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	socialMediaRequest := models.UpdateSocialMediaRequest{}
	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))
	userID := uint(userData["id"].(float64))
	//untuk memriksa jenis content type yang digunakan
	if contentType == appJson {
		if err := c.ShouldBindJSON(&socialMediaRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&socialMediaRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	socialMedia := models.SocialMedia{}
	socialMedia.ID = uint(socialMediaId)
	socialMedia.UserId = userID

	updateString, _ := json.Marshal(socialMediaRequest)
	updateData := models.SocialMedia{}
	json.Unmarshal(updateString, &updateData)

	//untuk update data social media hanya pada name dan social media url
	err := db.Model(&socialMedia).Updates(updateData).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	_ = db.First(&socialMedia, socialMedia.ID).Error

	socialMediaString, _ := json.Marshal(socialMedia)
	socialMediaResponse := models.UpdateSocialMediaResponse{}
	json.Unmarshal(socialMediaString, &socialMediaResponse)

	c.JSON(http.StatusOK, socialMediaResponse)
}

// untuk mengahpus Sosial Media existing
func DeleteSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)

	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))
	userID := uint(userData["id"].(float64))

	socialMedia := models.SocialMedia{}
	socialMedia.ID = uint(socialMediaId)
	socialMedia.UserId = userID

	//untuk menghapus social media sesuai ID
	err := db.Delete(&socialMedia).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
