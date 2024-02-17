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

// fungsi CreateComment untuk membuat komentar pada foto
func CreateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	commentRequest := models.CreateCommentRequest{}
	userID := uint(userData["id"].(float64))

	//memeriksa content-type yang direquest melalui API
	if contentType == appJson {
		if err := c.ShouldBindJSON(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	comment := models.Comment{
		PhotoId: commentRequest.PhotoId,
		Message: commentRequest.Message,
		UserId:  userID,
	}

	err := db.Debug().Create(&comment).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	commentString, _ := json.Marshal(comment)
	commentResponse := models.CreateCommentResponse{}
	json.Unmarshal(commentString, &commentResponse)

	c.JSON(http.StatusCreated, commentResponse)
}

// fungsi GetComment untuk get all comment yang telah dibuat
func GetComment(c *gin.Context) {
	db := database.GetDB()

	comments := []models.Comment{}

	err := db.Debug().Preload("User").Preload("Photo").Order("id asc").Find(&comments).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	commentsString, _ := json.Marshal(comments)
	commentsResponse := []models.CommentResponse{}
	json.Unmarshal(commentsString, &commentsResponse)

	c.JSON(http.StatusOK, commentsResponse)
}

// fungsi GetByIdComment untuk get satu komentar dengan user id
func GetByIdComment(c *gin.Context) {
	db := database.GetDB()
	commentId := c.Param("commentId")
	var commentById []models.Comment

	err := db.First(&commentById, "Id = ?", commentId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Record no found",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"Comment": commentById,
	})
}

// fungsi UpdateComment untuk update komentar
func UpdateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	commentRequest := models.UpdateCommentRequest{}
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := uint(userData["id"].(float64))

	//fungsi contentType untuk melihat komentar yang direquest melalui API
	if contentType == appJson {
		if err := c.ShouldBindJSON(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&commentRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	comment := models.Comment{}
	comment.ID = uint(commentId)
	comment.UserId = userID

	updateString, _ := json.Marshal(commentRequest)
	updateData := models.Comment{}
	json.Unmarshal(updateString, &updateData)
	//update koemntar hanya dapat dilakukan pada pesan
	err := db.Model(&comment).Updates(updateData).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	_ = db.First(&comment, comment.ID).Error

	commentString, _ := json.Marshal(comment)
	commentResponse := models.UpdateCommentResponse{}
	json.Unmarshal(commentString, &commentResponse)

	c.JSON(http.StatusOK, commentResponse)
}

// fungsi DeleteComment untuk menghapus komentar existing
func DeleteComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := uint(userData["id"].(float64))

	comment := models.Comment{}
	comment.ID = uint(commentId)
	comment.UserId = userID
	//untuk menghapus sesuai ID yang diminta
	err := db.Delete(&comment).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
