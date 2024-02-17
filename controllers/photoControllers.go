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

// fungsi CreatePhoto untuk membuat foto baru
func CreatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	photoRequest := models.CreatePhotoRequest{}
	userID := uint(userData["id"].(float64))
	//untuk meemeriksa jenis content type
	if contentType == appJson {
		if err := c.ShouldBindJSON(&photoRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&photoRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	photo := models.Photo{
		Title:    photoRequest.Title,
		Caption:  photoRequest.Caption,
		PhotoUrl: photoRequest.PhotoUrl,
		UserId:   userID,
	}

	err := db.Debug().Create(&photo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	_ = db.First(&photo, photo.ID).Error

	photoString, _ := json.Marshal(photo)
	photoResponse := models.CreatePhotoResponse{}
	json.Unmarshal(photoString, &photoResponse)

	c.JSON(http.StatusCreated, photoResponse)
}

// fungsi GetPhoto untuk mendapatkan semua photo yang sudah dicreate
func GetPhoto(c *gin.Context) {
	db := database.GetDB()

	photos := []models.Photo{}

	err := db.Debug().Preload("User").Order("id asc").Find(&photos).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	photosString, _ := json.Marshal(photos)
	photosResponse := []models.PhotoResponse{}
	json.Unmarshal(photosString, &photosResponse)

	c.JSON(http.StatusOK, photosResponse)
}

// fungsi GetbyIDphoto untuk mengambil foto yang sudah dicreate
func GetByIdPhoto(c *gin.Context) {
	db := database.GetDB()
	photoId := c.Param("photoId")
	var photoById []models.Photo

	err := db.First(&photoById, "Id = ?", photoId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Record not found",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Photo": photoById,
	})
}

// fungsi UpdatePhoto untuk mengupdate photo existing
func UpdatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	photoRequest := models.UpdatePhotoRequest{}
	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userID := uint(userData["id"].(float64))
	//untuk memeriksa jenis content type yang digunakan
	if contentType == appJson {
		if err := c.ShouldBindJSON(&photoRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&photoRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	photo := models.Photo{}
	photo.ID = uint(photoId)
	photo.UserId = userID

	updateString, _ := json.Marshal(photoRequest)
	updateData := models.Photo{}
	json.Unmarshal(updateString, &updateData)

	//Update yang dilakukan yaitu pada Title, Caption, dan Photo URL
	err := db.Model(&photo).Updates(updateData).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	photoString, _ := json.Marshal(photo)
	photoResponse := models.UpdatePhotoResponse{}
	json.Unmarshal(photoString, &photoResponse)

	c.JSON(http.StatusOK, photoResponse)
}

func DeletePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userID := uint(userData["id"].(float64))

	photo := models.Photo{}
	photo.ID = uint(photoId)
	photo.UserId = userID

	//untuk menghapus foto yang diminta sesuai ID yang sesuai
	err := db.Delete(&photo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}
