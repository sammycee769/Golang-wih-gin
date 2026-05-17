package handlers

import (
	"net/http"
	"todoList/db"
	"todoList/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(content *gin.Context) {
	var tasks []models.Task

	query := db.DB.Model(&models.Task{})
	if status := content.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tasks).Error; err != nil {
		content.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	content.JSON(http.StatusOK, gin.H{"data": tasks})
}
func isValidStatus(status models.Status) bool {
	switch status {
	case models.StatusCompleted, models.StatusInProgress, models.StatusPending:
		return true
	default:

		return false
	}
}

func CreateTasks(content *gin.Context) {
	var input models.CreateTaskInput

	if err := content.ShouldBindJSON(&input); err != nil {
		content.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, _ := content.MustGet("currentUser").(models.User)
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      models.StatusPending,
		DueDate:     input.DueDate,
		UserID:      user.ID,
	}
	if err := db.DB.Create(&task).Error; err != nil {
		content.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	content.JSON(http.StatusCreated, gin.H{"data": task})
}

func GetTask(content *gin.Context) {
	var task models.Task
	user, _ := content.MustGet("currentUser").(models.User)
	if err := db.DB.Where("user_id = ?", user.ID).First(&task, content.Param("id")).Error; err != nil {
		content.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	content.JSON(http.StatusOK, gin.H{"data": task})
	if task.UserID != user.ID {
		content.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this task"})
		return
	}
}

func UpdateTask(content *gin.Context) {
	var task models.Task
	user, _ := content.MustGet("currentUser").(models.User)
	if err := db.DB.Where("user_id = ?", user.ID).First(&task, content.Param("id")).Error; err != nil {
		content.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	var input models.UpdateTaskInput
	if err := content.ShouldBindJSON(&input); err != nil {
		content.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Status != nil {
		if !isValidStatus(*input.Status) {
			content.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
	}

	if err := db.DB.Model(&task).Updates(input).Error; err != nil {
		content.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	content.JSON(http.StatusOK, gin.H{"data": task})
	if task.UserID != user.ID {
		content.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this task"})
		return
	}
}

func DeleteTask(content *gin.Context) {
	var task models.Task
	user, _ := content.MustGet("currentUser").(models.User)
	if err := db.DB.Where("user_id = ?", user.ID).First(&task, content.Param("id")).Error; err != nil {
		content.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	db.DB.Delete(&task)
	content.JSON(http.StatusOK, gin.H{"data": task})
		if task.UserID != user.ID {
		content.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this task"})
		return
	}
}

func PatchTask(content *gin.Context) {
	var task models.Task
	user, _ := content.MustGet("currentUser").(models.User)
	if err := db.DB.Where("user_id = ?", user.ID).First(&task, content.Param("id")).Error; err != nil {
		content.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	var input models.UpdateTaskInput
	if err := content.ShouldBindJSON(&input); err != nil {
		content.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}
	if input.Status != nil {
		if !isValidStatus(*input.Status) {
			content.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
		task.Status = *input.Status
	}
	if err := db.DB.Save(&task).Error; err != nil {
		content.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	content.JSON(http.StatusOK, gin.H{"data": task})
	if task.UserID != user.ID {
		content.JSON(http.StatusForbidden, gin.H{"error": "you don't have access to this task"})
		return
	}
}
