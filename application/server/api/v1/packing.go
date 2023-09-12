package v1

import (
	bc "application/blockchain"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PackingSongrongRequestBody struct {
	AccountId   string  `json:"accountId"`   //操作人ID
	Songrong2ID   string  `json:"songrong2ID"`   //操作人ID
}

type QueryPackingSongrongRequestBody struct {
	BuyerID string `json:"buyerID"` //所有者(业主)(业主AccountId)
}

func PackingSongrong(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(PackingSongrongRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.AccountId))
	bodyBytes = append(bodyBytes, []byte(body.Songrong2ID))
	//调用智能合约
	resp, err := bc.ChannelExecute("packingSongrong", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func QueryPackingSongrong(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(QueryPackingSongrongRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.BuyerID != "" {
		bodyBytes = append(bodyBytes, []byte(body.BuyerID))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryPackingSongrongt", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}
