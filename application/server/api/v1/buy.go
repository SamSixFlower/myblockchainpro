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

type BuySongrongRequestBody struct {
	SongrongID   string  `json:"songrongID"`     //要购买的松茸的ID
	SellerID     string  `json:"sellerID"`       //采购商的ID
	BuyerID      string  `json:"buyerID"`        //购买者的ID
}

type ConfirmSongrongRequestBody struct {
	SongrongID   string  `json:"songrongID"`     //要购买的松茸的ID
	SellerID     string  `json:"sellerID"`       //采购商的ID
}

type QuerySellingBuyListRequestBody struct {
	SellerID     string  `json:"sellerID"`       //采购商的ID
}

type QuerySellingConfirmListRequestBody struct {
	BuyerID      string  `json:"buyerID"`        //购买者的ID
}

func BuySongrong(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(BuySongrongRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.SellerID == "" || body.BuyerID == "" {
		appG.Response(http.StatusBadRequest, "失败", "销售对象和Seller发起销售人不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.SongrongID))
	bodyBytes = append(bodyBytes, []byte(body.SellerID))
  	bodyBytes = append(bodyBytes, []byte(body.BuyerID))
	//调用智能合约
	resp, err := bc.ChannelExecute("buySongrong", bodyBytes)
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

func ConfirmSongrong(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(ConfirmSongrongRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.SellerID == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.SongrongID))
	bodyBytes = append(bodyBytes, []byte(body.SellerID))
	//调用智能合约
	resp, err := bc.ChannelExecute("confirmSongrong", bodyBytes)
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

func QuerySellingBuyList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(QuerySellingBuyListRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.SellerID != "" {
		bodyBytes = append(bodyBytes, []byte(body.SellerID))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("querySellingBuyList", bodyBytes)
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

func QuerySellingConfirmList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(QuerySellingConfirmListRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.BuyerID == "" {
		appG.Response(http.StatusBadRequest, "失败", "必须指定BuyerId查询")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.BuyerID))
	//调用智能合约
	resp, err := bc.ChannelQuery("querySellingConfirmList", bodyBytes)
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
