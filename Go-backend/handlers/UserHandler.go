/*
 * @Author: Jeffrey Zhu 1624410543@qq.com
 * @Date: 2025-03-14 23:10:57
 * @LastEditors: Jeffrey Zhu 1624410543@qq.com
 * @LastEditTime: 2025-03-15 23:48:57
 * @FilePath: \Smart-Snap-AI\Go-backend\handlers\UserHandlers.go
 * @Description: File Description Here...
 *
 * Copyright (c) 2025 by JeffreyZhu, All Rights Reserved.
 */
package handlers

import (
	"Go-backend/models"
	"Go-backend/utils"
	"net/http"

	"Go-backend/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration by creating a new user in the database
func Register(c *gin.Context) {
	var user models.User

	// log.Fatalf("%s", user.Email)
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: "Invalid request"})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := utils.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, models.Response{Code: 409, Message: "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := utils.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Code: 200, Message: "User registered successfully"})
}

// Login handles user authentication and returns a JWT token upon successful login
func Login(c *gin.Context) {
	var user models.User
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: "Invalid request"})
		return
	}

	if err := utils.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{Code: 401, Message: "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{Code: 401, Message: "Invalid email or password"})
		return
	}
	token, err := middleware.GenerateJWT(user) // 修改此行
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Code: 200, Message: "Login successful", Data: token})
}

func GetUserInfo(c *gin.Context) {
	var user models.User
	_, claims, _ := middleware.ValidToken(c.GetHeader("Authorization"))

	utils.DB.Where("id = ?", claims.UserID).First(&user)
	c.JSON(http.StatusOK, models.Response{Code: 200, Message: "Success", Data: user})
}
