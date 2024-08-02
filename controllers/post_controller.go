package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pramek008/go-jwt-project/database"
	"github.com/pramek008/go-jwt-project/models"
	"github.com/pramek008/go-jwt-project/utils"
)

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	post.UserID = userID.(uuid.UUID)

	if err := database.DB.Db.Create(&post).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	utils.SendResponse(c, http.StatusCreated, true, "Post created successfully", post)
}

func GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := database.DB.Db.Preload("User").First(&post, "id = ?", id).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	utils.SendResponse(c, http.StatusOK, true, "Post fetched successfully", post)
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := database.DB.Db.First(&post, "id = ?", id).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	userID, _ := c.Get("user_id")
	if post.UserID != userID.(uuid.UUID) {
		utils.SendErrorResponse(c, http.StatusForbidden, "Not authorized to update this post")
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	database.DB.Db.Save(&post)
	utils.SendResponse(c, http.StatusOK, true, "Post updated successfully", post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := database.DB.Db.First(&post, "id = ?", id).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	userID, _ := c.Get("user_id")
	if post.UserID != userID.(uuid.UUID) {
		utils.SendErrorResponse(c, http.StatusForbidden, "Not authorized to delete this post")
		return
	}

	database.DB.Db.Delete(&post)
	utils.SendResponse[map[string]interface{}](c, http.StatusOK, true, "Post deleted successfully", nil)
}

func ListPosts(c *gin.Context) {
	var posts []models.Post
	var total int64

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	// Count total posts
	database.DB.Db.Model(&models.Post{}).Count(&total)

	// Fetch paginated posts with preload
	database.DB.Db.Preload("User").Offset(offset).Limit(limit).Find(&posts)

	// Send response
	// utils.SendResponse(c, http.StatusOK, true, "Posts fetched successfully", gin.H{
	// 	"page":  page,
	// 	"limit": limit,
	// 	"total": total,
	// 	"posts": posts,
	// })

	utils.SendPaginatedResponse(c, http.StatusOK, true, "Post fetched successfully,", posts, int64(limit), int64(page), total)
}
