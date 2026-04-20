package broadcast

import (
	"chat/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ViewBroadcastAPI(c *gin.Context) {
	c.JSON(http.StatusOK, getLatestBroadcast(c))
}

func CreateBroadcastAPI(c *gin.Context) {
	user := auth.RequireAdmin(c)
	if user == nil {
		return
	}

	var form createRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, createResponse{
			Status: false,
			Error:  err.Error(),
		})
	}

	err := createBroadcast(c, user, form.Content)
	if err != nil {
		c.JSON(http.StatusOK, createResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, createResponse{
		Status: true,
	})
}

func GetBroadcastListAPI(c *gin.Context) {
	user := auth.RequireAdmin(c)
	if user == nil {
		return
	}

	data, err := getBroadcastList(c)
	if err != nil {
		c.JSON(http.StatusOK, listResponse{
			Data: []Info{},
		})
		return
	}

	c.JSON(http.StatusOK, listResponse{
		Data: data,
	})
}

func RemoveBroadcastAPI(c *gin.Context) {
	user := auth.RequireAdmin(c)
	if user == nil {
		return
	}

	indexStr := c.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusOK, createResponse{
			Status: false,
			Error:  "invalid index",
		})
		return
	}

	if err := removeBroadcast(c, index); err != nil {
		c.JSON(http.StatusOK, createResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, createResponse{
		Status: true,
	})
}

func UpdateBroadcastAPI(c *gin.Context) {
	user := auth.RequireAdmin(c)
	if user == nil {
		return
	}

	var form updateRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, createResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	if err := updateBroadcast(c, form.ID, form.Content); err != nil {
		c.JSON(http.StatusOK, createResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, createResponse{
		Status: true,
	})
}
