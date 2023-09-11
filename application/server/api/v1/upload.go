package v1

import (
	bc "application/blockchain"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UploadSongrongRequestBody struct {
	AccountId   string  `json:"accountId"`   //操作人ID
	Umbrella  string  `json:"umbrella"`  //是否开伞
	SongrongID   string `json:"songrongID"`   //来自哪个批次的松茸ID
	Size string `json:"size"` //大小
  	Storage string `json:"storage"` //存储情况
  	Time1 string `json:"time1"` //入库时间
  	Time2 string `json:"time2"` //出库时间
}

type QueryUploadSongrongRequestBody struct {
	BuyerID      string  `json:"buyerID"`        //购买者的ID
}

func UploadSongrong(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(UploadSongrongRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.Umbrella == "" || body.SongrongID == "" || body.Size == "" || body.Storage == "" || body.Time1 == "" || body.Time2 == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.AccountId))
	bodyBytes = append(bodyBytes, []byte(body.Umbrella))
  	bodyBytes = append(bodyBytes, []byte(body.SongrongID))
  	bodyBytes = append(bodyBytes, []byte(body.Size))
  	bodyBytes = append(bodyBytes, []byte(body.Storage))
  	bodyBytes = append(bodyBytes, []byte(body.Time1))
  	bodyBytes = append(bodyBytes, []byte(body.Time2))
	//调用智能合约
	resp, err := bc.ChannelExecute("uploadSongrong", bodyBytes)
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

func QueryUploadSongrong(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(QueryUploadSongrongRequestBody)
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
	resp, err := bc.ChannelQuery("queryUploadSongrong", bodyBytes)
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
