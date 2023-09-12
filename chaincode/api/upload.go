//提交分批后的松茸信息，也就是新建songrong2
//new chaincode
package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
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
 type SongRong2 struct {
	SongRong2ID string  `json:"songrong2Id"` //分装后松茸ID
	FactoryID string  `json:"factoryId"` //厂家ID
	Umbrella   string  `json:"umbrella"`   //开伞情况
	SongRongID string  `json:"songrongId"` //来自哪批次的松茸ID
	Size   string `json:"size"`    //松茸大小
	Storage string `json:"storage"`    //库存时状态
	Time  string `json:"time"`  //上链时间
	Time1  string `json:"time"`  //入库时间
	Time2  string `json:"time"`  //出库时间
}
*/
// 新建松茸上传信息(厂家)
func UploadSongrong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//传入厂家id，开伞情况，来源的松茸id，松茸大小，库存时状态，入库时间，出库时间
	accountId := args[0] //accountId用于验证是否为厂家
	umbrella := args[1]
	songrongID := args[2]
  	size := args[3]
  	storage := args[4]
  	time1 := args[5]
  	time2 := args[6]
	if accountId == ""  {
		return shim.Error("参数存在空值")
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
	if account.UserName != "松茸厂家" {
		return shim.Error(fmt.Sprintf("操作人权限不足%s", err))
	}
	createTime, _ := stub.GetTxTimestamp()
	songrong2 := &model.SongRong2{
		SongRong2ID: stub.GetTxID()[:16],
		FactoryID:   accountId,
		Umbrella:    umbrella,
		SongRongID:  songrongID,
    		Size:  size,
    		Storage:  storage,
		Time:  time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
	  	Time1: time1,
		Time2: time2,
	}
	// 写入账本
	if err := utils.WriteLedger(songrong2, stub, model.SellSongrong2Key, []string{songrong2.SongRong2ID, songrong2.FactoryID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	songrong2Byte, err := json.Marshal(songrong2)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(songrong2Byte)
}

// QueryUploadSongrong 查询分批后上链的松茸(传入类型和厂家ID)
func QueryUploadSongrong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var songrong2List []model.SongRong2
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellSongrong2Key, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var songrong2 model.SongRong2
			err := json.Unmarshal(v, &songrong2)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryUploadSongrong-反序列化出错: %s", err))
			}
			songrong2List = append(songrong2List, songrong2)
		}
	}
	songrong2ListByte, err := json.Marshal(songrong2List)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryUploadSongrong-序列化出错: %s", err))
	}
	return shim.Success(songrong2ListByte)
}
