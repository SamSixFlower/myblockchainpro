//把最终打包商品上链
//提交分批后的松茸信息，也就是新建songrong2
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
type SongRong3 struct {
	ProductID string  `json:"productId"` //产品号ID
	FactoryID string  `json:"factoryId"` //厂家ID
	SongRong2ID string  `json:"songrongeId"` //产品装的松茸的ID
	Time  string `json:"time"`  //包装时间
}
*/
// 新建松茸上传信息(厂家)
func PackingSongrong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//传入厂家id，开伞情况，来源的松茸id，松茸大小，库存时状态，入库时间，出库时间
	accountId := args[0] //accountId用于验证是否为厂家
	songrong2ID := args[1]
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
	songrong3 := &model.SongRong3{
		ProductID: stub.GetTxID()[:16],
		FactoryID:   accountId,
		SongRong2ID:  songrong2ID
		Time:  time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
	}
	// 写入账本
	if err := utils.WriteLedger(songrong3, stub, model.Sellsongrong3Key, []string{songrong3.ProductID, songrong3.FactoryID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	songrong3Byte, err := json.Marshal(songrong3)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(songrong3Byte)
}

// QueryPackingSongrong 查询售卖的松茸(可以供采购商和厂家查询，采购商只能传入自己的ID，厂家可以选择采购商ID进行传入)
func QueryPackingSongrong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var songrong3List []model.SongRong3
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.Sellsongrong3Key, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var songrong3 model.SongRong3
			err := json.Unmarshal(v, &songrong3)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellSongrong-反序列化出错: %s", err))
			}
			songrong3List = append(songrong3List, songrong3)
		}
	}
	songrong3ListByte, err := json.Marshal(songrong3List)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellSongrong-序列化出错: %s", err))
	}
	return shim.Success(songrong3ListByte)
}
