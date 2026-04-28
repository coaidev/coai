package channel

import (
	"chat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func reloadChannelsOnSuccess(err error) {
	if err == nil {
		ConduitInstance.Load()
	}
}

type SyncChargeForm struct {
	Overwrite bool           `json:"overwrite"`
	Data      ChargeSequence `json:"data"`
}

func GetInfo(c *gin.Context) {
	c.JSON(http.StatusOK, SystemInstance.AsInfo())
}

func AttachmentService(c *gin.Context) {
	// /attachments/:hash -> ~/storage/attachments/:hash
	hash := c.Param("hash")
	c.File(fmt.Sprintf("storage/attachments/%s", hash))
}

func DeleteChannel(c *gin.Context) {
	id := c.Param("id")
	state := ConduitInstance.DeleteChannel(utils.ParseInt(id))
	reloadChannelsOnSuccess(state)

	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}

func ActivateChannel(c *gin.Context) {
	id := c.Param("id")
	state := ConduitInstance.ActivateChannel(utils.ParseInt(id))
	reloadChannelsOnSuccess(state)

	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}

func DeactivateChannel(c *gin.Context) {
	id := c.Param("id")
	state := ConduitInstance.DeactivateChannel(utils.ParseInt(id))
	reloadChannelsOnSuccess(state)

	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}

func GetChannelList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   ConduitInstance.Sequence,
	})
}

func GetChannel(c *gin.Context) {
	id := c.Param("id")
	channel := ConduitInstance.Sequence.GetChannelById(utils.ParseInt(id))

	c.JSON(http.StatusOK, gin.H{
		"status": channel != nil,
		"data":   channel,
	})
}

func CreateChannel(c *gin.Context) {
	var channel Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	state := ConduitInstance.CreateChannel(&channel)
	reloadChannelsOnSuccess(state)
	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}

func UpdateChannel(c *gin.Context) {
	var channel Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	id := c.Param("id")
	channel.Id = utils.ParseInt(id)

	state := ConduitInstance.UpdateChannel(channel.Id, &channel)
	reloadChannelsOnSuccess(state)
	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}

func SetCharge(c *gin.Context) {
	var charge Charge
	if err := c.ShouldBindJSON(&charge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	state := ChargeInstance.SetRule(charge)
	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}

func GetChargeList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   ChargeInstance.ListRules(),
	})
}

func DeleteCharge(c *gin.Context) {
	id := c.Param("id")
	state := ChargeInstance.DeleteRule(utils.ParseInt(id))

	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}

func SyncCharge(c *gin.Context) {
	var form SyncChargeForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
	}

	state := ChargeInstance.SyncRules(form.Data, form.Overwrite)
	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}

func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   SystemInstance,
	})
}

func UpdateConfig(c *gin.Context) {
	var config SystemConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	state := SystemInstance.UpdateConfig(&config)
	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}

func GetPlanConfig(c *gin.Context) {
	c.JSON(http.StatusOK, PlanInstance)
}

func UpdatePlanConfig(c *gin.Context) {
	var config PlanManager
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	state := PlanInstance.UpdateConfig(&config)
	c.JSON(http.StatusOK, gin.H{
		"status": state == nil,
		"error":  utils.GetError(state),
	})
}
