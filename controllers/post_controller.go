package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pramek008/go-jwt-project/database"
	"github.com/pramek008/go-jwt-project/models"
	"github.com/pramek008/go-jwt-project/utils"
)

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	post.UserID = userID.(uint)

	if err := database.DB.Db.Create(&post).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// c.JSON(http.StatusCreated, post)
	utils.SendResponse(c, http.StatusCreated, true, "Post created successfully", post)
}

func GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := database.DB.Db.Preload("User").First(&post, id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		utils.SendErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	// c.JSON(http.StatusOK, post)
	utils.SendResponse(c, http.StatusOK, true, "Post fetched successfully", post)
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := database.DB.Db.First(&post, id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		utils.SendErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	userID, _ := c.Get("user_id")
	if post.UserID != userID.(uint) {
		// c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this post"})
		utils.SendErrorResponse(c, http.StatusForbidden, "Not authorized to update this post")
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	database.DB.Db.Save(&post)
	// c.JSON(http.StatusOK, post)
	utils.SendResponse(c, http.StatusOK, true, "Post updated successfully", post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := database.DB.Db.First(&post, id).Error; err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		utils.SendErrorResponse(c, http.StatusNotFound, "Post not found")
		return
	}

	userID, _ := c.Get("user_id")
	if post.UserID != userID.(uint) {
		// c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this post"})
		utils.SendErrorResponse(c, http.StatusForbidden, "Not authorized to delete this post")
		return
	}

	database.DB.Db.Delete(&post)
	// c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
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
	utils.SendResponse(c, http.StatusOK, true, "Posts fetched successfully", gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"posts": posts,
	})
}
