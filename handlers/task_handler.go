package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"
	"todoList/models"

	"github.com/gin-gonic/gin"
)

var (
	tasks  []models.Task
	nextID uint = 1
	mu     sync.Mutex
)

func GetTasks(content *gin.Context) {
	statusQuery := content.Query("status")

	mu.Lock()
	defer mu.Unlock()

	if statusQuery != "" && !isValidStatus(models.Status(statusQuery)) {
		content.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	if statusQuery == "" {
		content.JSON(http.StatusOK, gin.H{"data": tasks})
		return
	}
	var filteredTasks []models.Task
	for _, task := range tasks {
		if string(task.Status) == statusQuery {
			filteredTasks = append(filteredTasks, task)
		}
	}
	content.JSON(http.StatusOK, gin.H{"data": filteredTasks})
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
	mu.Lock()
	defer mu.Unlock()

	newTask := models.Task{
		ID:          nextID,
		Title:       input.Title,
		Description: input.Description,
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
	}
	nextID++
	tasks = append(tasks, newTask)
	content.JSON(http.StatusCreated, gin.H{"data": tasks})
}

func GetTask(content *gin.Context) {
	idStr := content.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		content.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for _, task := range tasks {
		if task.ID == uint(id) {
			content.JSON(http.StatusOK, gin.H{"data": task})
			return
		}
	}
	content.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func UpdateTask(content *gin.Context) {
	idStr := content.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		content.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input models.UpdateTaskInput
	if err := content.ShouldBindJSON(&input); err != nil {
		content.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for count, task := range tasks {
		if task.ID == uint(id) {
			if input.Title != nil {
				tasks[count].Title = *input.Title
			}
			if input.Description != nil {
				tasks[count].Description = *input.Description
			}
			if input.Status != nil {
				tasks[count].Status = *input.Status
			}
			tasks[count].UpdatedAt = time.Now()
			content.JSON(http.StatusOK, gin.H{"data": tasks[count]})
			return
		}
	}
	content.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func DeleteTask(content *gin.Context) {
	idStr := content.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		content.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for count, task := range tasks {
		if task.ID == uint(id) {
			tasks = append(tasks[:count], tasks[count+1:]...)
			content.Status(http.StatusNoContent)
			return
		}
	}
	content.JSON(http.StatusNotFound, gin.H{"error": "not found"})

}
