package Controllers

import (
	"fmt"
	"golang-assesment/Database"
	"golang-assesment/Middleware"
	"golang-assesment/Models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type InsertTask struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"not null" json:"title" binding:"required"`
	Description string `gorm:"not null" json:"description" binding:"required"`
	Status      string `gorm:"default:'pending'" json:"status" binding:"required,oneof=pending in-progress completed"`
	CreatedAt   string `json:"CreatedAt"`
	UpdatedAt   string `json:"UpdatedAt"`
}

func GetTaskList(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		log.Error("  @GetTaskList Invalid page number")
		ValidationResponse(c, "Invalid page number.")
		return
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		log.Error("  @GetTaskList Invalid limit")
		ValidationResponse(c, "Invalid limit.")
		return
	}

	offset := (pageNum - 1) * limitNum

	var tasks []Models.AmplTaskList
	result := Database.DB.Limit(limitNum).Offset(offset).Find(&tasks)
	if result.Error != nil {
		log.Error("  @GetTaskList error while fetch tasklist from Database", result.Error.Error())
		ValidationResponse(c, "Something went wrong while fetch Tasklist.")
		return
	}

	var total int64
	Database.DB.Model(&Models.AmplTaskList{}).Count(&total)
	successResponse(c, "Retrieve tasks successfully!!", map[string]interface{}{
		"data":  tasks,
		"page":  pageNum,
		"limit": limitNum,
		"total": total,
	})

}

func GetTask(c *gin.Context) {
	id := c.Query("id")
	fetchResp, status, msg := FetchTaskFromDB(id)

	if !status {
		log.Error(id, "  @GetTask 0 records fetch while fetch from DB:", msg)
		ValidationResponse(c, "Something went wrong while fetch task.")
		return
	}
	log.Info(id, "  @GetTask Retrive task successfully")
	successResponse(c, "Retrive task successfully!!", map[string]interface{}{"task": fetchResp})
}

func CreateTask(c *gin.Context) {
	var insertNewEntity Models.AmplTaskList

	if err := c.ShouldBindJSON(&insertNewEntity); err != nil {
		log.Error(insertNewEntity.ID, " @CreateTask Validation Error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := Database.DB.Create(&insertNewEntity).Error
	if err != nil {
		log.Error(insertNewEntity.ID, " @CreateTask Error while fetch data from db and Error: ", err)
		ValidationResponse(c, "Something went wrong while insertion")
		return

	}
	log.Info(insertNewEntity.ID, "  @CreateTask Task created successfully")
	successResponse(c, "Task created successfully!!", map[string]interface{}{"task": insertNewEntity})
}

func FetchTaskFromDB(id string) (Models.AmplTaskList, bool, string) {
	var task Models.AmplTaskList
	fetchResp := Database.DB.Table("AmplTaskList atl").Select("*").Where("atl.ID=?", id).Find(&task)
	if fetchResp.Error != nil {
		log.Error(id, " @GetTask Error while fetch data from db and Error: ", fetchResp.Error)
		return task, false, "Error while fetch data from db"

	}
	if fetchResp.RowsAffected > 0 {
		log.Error(id, "  @GetTask 0 records fetch while fetch from DB:")
		return task, true, ""
	}
	return task, false, "Something went wrong  while fetch task from db"
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task Models.AmplTaskList
	_, status, msg := FetchTaskFromDB(id)

	if !status {
		log.Error(id, "  @GetTask 0 records fetch while fetch from DB:", msg)
		ValidationResponse(c, "Something went wrong while fetch task.")
		return
	}
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskUpdate := Database.DB.Table("AmplTaskList").Where("id =?", id).Updates(map[string]interface{}{
		"title":       task.Title,
		"status":      task.Status,
		"description": task.Description,
	})
	fmt.Println("RO", taskUpdate.RowsAffected)

	if taskUpdate.Error != nil {
		log.Error(task.ID, " @UpdateTask Something went wrong while update into DB  ", taskUpdate.Error)
		Database.ResetDBPoolConnection("taskManagementDB", taskUpdate.Error.Error())
		ValidationResponse(c, "Something went wrong while update into DB")
		return
	}
	if taskUpdate.RowsAffected == 0 {
		log.Error(task.ID, " @UpdateTask 0 record affected while update into DB  , taskUpdate.RowsAffected  :", taskUpdate.RowsAffected)
		ValidationResponse(c, "0 records affected while update into DB")
		return
	}

	log.Info(task.ID, "  @UpdateTask Task Updated successfully")
	successResponse(c, "Task Updated successfully!!", map[string]interface{}{})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := Database.DB.Delete(&Models.AmplTaskList{}, id).Error; err != nil {
		log.Error(id, " @DeleteTask Something went wrong while update into DB  ", err)
		NoDataFoundResponse(c, "Data not found.")
		return
	}
	log.Info(id, "  @DeleteTask Task deleted successfully")
	successResponse(c, "Task deleted successfully!!", map[string]interface{}{})
}

func LoginAuth(c *gin.Context) {

	var input struct {
		UserId int `json:"userid" validate:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	var task Models.AmplTaskList
	if err := Database.DB.Where("ID = ?", input.UserId).First(&task).Error; err != nil {
		log.Error(task.ID, " @LoginAuth Error while check userId into Db", err)
		ValidationResponse(c, "Please enter valid userid")
		return
	}

	userID := strconv.FormatUint(uint64(task.ID), 10)
	token, err := Middleware.JWTTokenGenerate(userID)
	if err != nil {
		log.Error(userID, " @LoginAuth Error generating token")
		ValidationResponse(c, "Error generating token")
		return
	}

	log.Info(userID, "  @DeleteTask Task deleted successfully")
	successResponse(c, "Token generated successfully!!", map[string]interface{}{"token": token})
}
