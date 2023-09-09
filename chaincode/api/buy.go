//厂家买松茸
func BuySongrong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//传入松茸id，采购商id，购买者id
	songrongID := args[0]
	sellerID := args[1]
	buyerID := args[2]
	if songrongID == "" || sellerID == "" || buyerID == "" {
		return shim.Error("参数存在空值")
	}
	if sellerID == buyerID {
		return shim.Error("买家和卖家不能同一人")
	}
	//取出要购买的松茸批次
  	resultssongrong1, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellsongrongKey, []string{songrongID, sellerID})
	var songrong1 model.SongRong1
	if err = json.Unmarshal(resultssongrong1[0], &songrong1); err != nil {
		return shim.Error(fmt.Sprintf("BuySongrong-反序列化出错: %s", err))
	}
	//判断songrong1的状态是否为销售中
	if songrong1.SellingStatus != "saleStart" {
		return shim.Error("此交易不属于销售中状态，已经无法购买")
	}
	//根据buyer获取买家信息
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{buyerID})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("buyer买家信息验证失败%s", err))
	}
	var buyerAccount model.Account
	if err = json.Unmarshal(resultsAccount[0], &buyerAccount); err != nil {
		return shim.Error(fmt.Sprintf("查询buyer买家信息-反序列化出错: %s", err))
	}
	if buyerAccount.UserName == "采购商" {
		return shim.Error(fmt.Sprintf("采购商不能购买%s", err))
	}
	//将buyer写入交易selling,修改交易状态
	songrong1.BuyerID = buyerID
	songrong1.SellingStatus = "delivery"
	if err := utils.WriteLedger(songrong1, stub, model.SellsongrongKey, []string{songrongID, sellerID}); err != nil {
		return shim.Error(fmt.Sprintf("将buyer写入交易songrong1,修改交易状态 失败%s", err))
	}
	createTime, _ := stub.GetTxTimestamp()
	//将本次购买交易写入账本,可供买家查询
	sellingBuy := &model.SellingBuy{
		BuyerID:      buyerID,
		CreateTime: time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		Selling:    songrong1,
	}
	if err := utils.WriteLedger(sellingBuy, stub, model.SellingBuyKey, []string{sellingBuy.BuyerID, sellingBuy.CreateTime}); err != nil {
		return shim.Error(fmt.Sprintf("将本次购买交易写入账本失败%s", err))
	}
	sellingBuyByte, err := json.Marshal(sellingBuy)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(sellingBuyByte)
}
