package main

import (
	"chaincode/api"
	"chaincode/model"
	"chaincode/pkg/utils"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type BlockChainSuYuan struct {
}

// Init 链码初始化
func (t *BlockChainSuYuan) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("链码初始化")
	//初始化默认数据
	var accountIds = [6]string{
		"5feceb66ffc8",
		"6b86b273ff34",
	}
	var userNames = [6]string{"采购商", "松茸厂家"}
	//初始化账号数据
	for i, val := range accountIds {
		account := &model.Account{
			AccountId: val,
			UserName:  userNames[i],
		}
		// 写入账本
		if err := utils.WriteLedger(account, stub, model.AccountKey, []string{val}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
	}
	return shim.Success(nil)
}

// Invoke 实现Invoke接口调用智能合约
func (t *BlockChainSuYuan) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	switch funcName {
	case "hello":
		return api.Hello(stub, args)
	case "queryAccountList":
		return api.QueryAccountList(stub, args)
	case "sellSongrong":
		return api.SellSongrong(stub, args)
	case "buySongrong":
		return api.BuySongrong(stub, args)
	case "confirmSongrong":
		return api.ConfirmSongrong(stub, args)
	case "uploadSongrong":
		return api.UploadSongrong(stub, args)
	case "packingSongrong":
		return api.PackingSongrong(stub, args)
	case "querySellSongrong":
		return api.QuerySellSongrong(stub, args)
	case "querySellingBuyList":
		return api.QuerySellingBuyList(stub, args)
	case "querySellingConfirmList":
		return api.QuerySellingConfirmList(stub, args)
	case "queryUploadSongrong":
		return api.QueryUploadSongrong(stub, args)
	case "queryPackingSongrong":
		return api.QueryPackingSongrong(stub, args)
	default:
		return shim.Error(fmt.Sprintf("没有该功能: %s", funcName))
	}
}

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = timeLocal
	err = shim.Start(new(BlockChainSuYuan))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
