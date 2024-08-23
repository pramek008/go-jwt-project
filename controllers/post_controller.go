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
	// if err := c.ShouldBindJSON(&post); err != nil {
	// 	utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	post.Title = c.Request.FormValue("title")
	post.Content = c.Request.FormValue("content")

	userID, _ := c.Get("user_id")
	post.UserID = userID.(uuid.UUID)

	file, _ := c.FormFile("file")
	if file != nil {
		fileUrl, err := utils.UploadFile(c, file)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to upload file")
			return
		}
		post.FileURL = fileUrl
	}

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

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	// if err := c.ShouldBindJSON(&post); err != nil {
	// 	utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	post.Title = c.Request.FormValue("title")
	post.Content = c.Request.FormValue("content")

	file, _ := c.FormFile("file")
	if file != nil {
		fileUrl, err := utils.UploadFile(c, file)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to upload file")
			return
		}
		post.FileURL = fileUrl
	}

	// database.DB.Db.Save(&post)
	if err := database.DB.Db.Save(&post); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update post")
		return
	}
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

	// Convert posts to PostResponse format
	postResponses := []models.PostResponse{}
	for _, post := range posts {
		postResponse := models.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			FileURL:   post.FileURL,
			UserID:    post.UserID,
			User:      models.UserResponse{ID: post.User.ID, Nickname: post.User.Nickname, Email: post.User.Email, CreatedAt: post.User.CreatedAt, UpdatedAt: post.User.UpdatedAt},
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			DeletedAt: &post.DeletedAt.Time,
		}
		postResponses = append(postResponses, postResponse)
	}

	// Send response with an empty list if no posts were found
	utils.SendPaginatedResponse(c, http.StatusOK, true, "Post fetched successfully,", postResponses, int64(limit), int64(page), total)
}
