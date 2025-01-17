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

type SongrongRequestBody struct {
	AccountId   string  `json:"accountId"`   //操作人ID
	Place  string  `json:"place"`  //产地
	Amount   float64 `json:"amount"`   //售卖总量
}

type SongrongQueryRequestBody struct {
	AccountId   string  `json:"accountId"`   //操作人ID
}

func SellSongrong(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SongrongRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.Amount  < 0  {
		appG.Response(http.StatusBadRequest, "失败", "售卖量必须大于0")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.AccountId))
	bodyBytes = append(bodyBytes, []byte(body.Place))
	bodyBytes = append(bodyBytes, []byte(strconv.FormatFloat(body.Amount, 'E', -1, 64)))
	//调用智能合约
	resp, err := bc.ChannelExecute("sellSongrong", bodyBytes)
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

func QuerySellSongrong(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SongrongQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.AccountId != "" {
		bodyBytes = append(bodyBytes, []byte(body.AccountId))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("querySellSongrong", bodyBytes)
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
