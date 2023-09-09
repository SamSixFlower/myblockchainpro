//new chaincode
package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
/*
case "sellSongrong":
		return api.SellSongrong(stub, args)
	case "buySongrong":
		return api.BuySongrong(stub, args)
	case "signSongrong":
		return api.SignSongrong(stub, args)
	case "signBaozhuang":
		return api.SignBaozhuang(stub, args)
	case "querySellSongrong":
		return api.QuerySellSongrong(stub, args)
	case "queryBuySongrong":
		return api.QueryBuySongrong(stub, args)
	case "querySignSongrong":
		return api.QueryignSongrong(stub, args)
	case "querySignBaozhuang":
		return api.QuerySignBaozhuang(stub, args)
 type SongRong1 struct {
	SongRongID string  `json:"songrongeId"` //本批松茸ID
	SellerID   string  `json:"sellerID"`   //采购商ID
	Place  string    `json:"place"`  //松茸产地
	Amount   float64 `json:"amount"`    //售卖总量
	Time  string `json:"time"`  //售卖时间
	SellingStatus string  `json:"sellingStatus"` //销售状态
 	BuyerID string  `json:"buyerID"` //购买者
}
*/
// 新建松茸售卖信息(采购商)
func SellSongrong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//传入采购商id，产地，售卖数量，售卖时间
	accountId := args[0] //accountId用于验证是否为采购商
	place := args[1]
	amount := args[2]
	time := args[3]
	if accountId == ""  {
		return shim.Error("参数存在空值")
	}
	// 参数数据格式转换
	var formattedAmount float64
	if val, err := strconv.ParseFloat(amount, 64); err != nil {
		return shim.Error(fmt.Sprintf("参数格式转换出错: %s", err))
	} else {
		formattedLAmount = val
	}
	//判断是否采购商操作
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{accountId})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("操作人权限验证失败%s", err))
	}
	var account model.Account
	if err = json.Unmarshal(resultsAccount[0], &account); err != nil {
		return shim.Error(fmt.Sprintf("查询操作人信息-反序列化出错: %s", err))
	}
	if account.UserName != "采购商" {
		return shim.Error(fmt.Sprintf("操作人权限不足%s", err))
	}
	songrong1 := &model.SongRong1{
		SongRongID: stub.GetTxID()[:16],
		SellerID:   accountId,
		Place:    place,
		Amount:  formattedAmount,
		Time:  time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		SellingStatus: model.SellingStatusConstant()["saleStart"],
		BuyerID: "",
	}
	// 写入账本
	if err := utils.WriteLedger(songrong1, stub, model.SellsongrongKey, []string{songrong1.SongRongID, songrong1.SellerID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	songrong1Byte, err := json.Marshal(songrong1)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(songrong1Byte)
}

/*
// CreateRealEstate 新建房地产(管理员)
func CreateRealEstate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	accountId := args[0] //accountId用于验证是否为管理员
	proprietor := args[1]
	totalArea := args[2]
	livingSpace := args[3]
	if accountId == "" || proprietor == "" || totalArea == "" || livingSpace == "" {
		return shim.Error("参数存在空值")
	}
	if accountId == proprietor {
		return shim.Error("操作人应为管理员且与所有人不能相同")
	}
	// 参数数据格式转换
	var formattedTotalArea float64
	if val, err := strconv.ParseFloat(totalArea, 64); err != nil {
		return shim.Error(fmt.Sprintf("totalArea参数格式转换出错: %s", err))
	} else {
		formattedTotalArea = val
	}
	var formattedLivingSpace float64
	if val, err := strconv.ParseFloat(livingSpace, 64); err != nil {
		return shim.Error(fmt.Sprintf("livingSpace参数格式转换出错: %s", err))
	} else {
		formattedLivingSpace = val
	}
	//判断是否管理员操作
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{accountId})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("操作人权限验证失败%s", err))
	}
	var account model.Account
	if err = json.Unmarshal(resultsAccount[0], &account); err != nil {
		return shim.Error(fmt.Sprintf("查询操作人信息-反序列化出错: %s", err))
	}
	if account.UserName != "管理员" {
		return shim.Error(fmt.Sprintf("操作人权限不足%s", err))
	}
	//判断业主是否存在
	resultsProprietor, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{proprietor})
	if err != nil || len(resultsProprietor) != 1 {
		return shim.Error(fmt.Sprintf("业主proprietor信息验证失败%s", err))
	}
	realEstate := &model.RealEstate{
		RealEstateID: stub.GetTxID()[:16],
		Proprietor:   proprietor,
		Encumbrance:  false,
		TotalArea:    formattedTotalArea,
		LivingSpace:  formattedLivingSpace,
	}
	// 写入账本
	if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	realEstateByte, err := json.Marshal(realEstate)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(realEstateByte)
}
*/

// QueryRealEstateList 查询房地产(可查询所有，也可根据所有人查询名下房产)
func QuerySellSongrong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var songrong1List []model.SongRong1
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellsongrongKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var songrong1 model.SongRong1
			err := json.Unmarshal(v, &songrong1)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellSongrong-反序列化出错: %s", err))
			}
			realEstateList = append(songrong1List, songrong1)
		}
	}
	songrong1ListByte, err := json.Marshal(songrong1List)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellSongrong-序列化出错: %s", err))
	}
	return shim.Success(songrong1ListByte)
}
/*
// QueryRealEstateList 查询房地产(可查询所有，也可根据所有人查询名下房产)
func QueryRealEstateList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var realEstateList []model.RealEstate
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var realEstate model.RealEstate
			err := json.Unmarshal(v, &realEstate)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryRealEstateList-反序列化出错: %s", err))
			}
			realEstateList = append(realEstateList, realEstate)
		}
	}
	realEstateListByte, err := json.Marshal(realEstateList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryRealEstateList-序列化出错: %s", err))
	}
	return shim.Success(realEstateListByte)
}
*/
