package controllers

import (
	"encoding/json"
	"mygram-final/database"
	"mygram-final/helpers"
	"mygram-final/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	appJson = "Application/json"
)

// untuk mendaftarkan user baru
func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}
	//untuk memeriksa jenis content type yang digunakan
	if contentType == appJson {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	err := db.Debug().Create(&User).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":              User.ID,
		"email":           User.Email,
		"username":        User.Username,
		"age":             User.Age,
		"urlImageProfile": User.ProfileImageUrl,
	})
}

// untuk login user
func UserLogin(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}
	password := ""
	//untuk memriksa jenis content type yang digunakan
	if contentType == appJson {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}
	password = User.Password

	//untuk memeriksa kecocokan email dan password untuk proses login
	err := db.Debug().Where("email = ? ", User.Email).Take(&User).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid Email/Password",
		})
		return
	}
	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))
	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid Email/Password",
		})
		return
	}
	token := helpers.GenerateToken(User.ID, User.Email)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// untuk mengubah user existing
func UpdateUser(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	//Hanya update email dan username saja
	userRequest := models.UpdateUserRequest{}
	userID := uint(userData["id"].(float64))
	//pengecekan jenis content type yang digunakan
	if contentType == appJson {
		if err := c.ShouldBindJSON(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&userRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
	}

	user := models.User{}
	user.ID = userID

	updateString, _ := json.Marshal(userRequest)
	updateData := models.User{}
	json.Unmarshal(updateString, &updateData)

	err := db.Model(&user).Updates(updateData).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}
	_ = db.First(&user, user.ID).Error

	userString, _ := json.Marshal(user)
	userResponse := models.CreateUserResponse{}
	json.Unmarshal(userString, &userResponse)

	//untuk memanmpilkan email dan username yang sudah diupdate
	c.JSON(http.StatusCreated, userResponse)
}

// untuk menghapus user exsiting
func DeleteUser(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)

	userID := uint(userData["id"].(float64))

	user := models.User{}
	user.ID = userID

	//untuk menghapus user register
	err := db.Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
