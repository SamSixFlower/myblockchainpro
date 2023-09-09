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
	//取出
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
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{buyer})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("buyer买家信息验证失败%s", err))
	}
	var buyerAccount model.Account
	if err = json.Unmarshal(resultsAccount[0], &buyerAccount); err != nil {
		return shim.Error(fmt.Sprintf("查询buyer买家信息-反序列化出错: %s", err))
	}
	if buyerAccount.UserName == "管理员" {
		return shim.Error(fmt.Sprintf("管理员不能购买%s", err))
	}
	//判断余额是否充足
	if buyerAccount.Balance < selling.Price {
		return shim.Error(fmt.Sprintf("房产售价为%f,您的当前余额为%f,购买失败", selling.Price, buyerAccount.Balance))
	}
	//将buyer写入交易selling,修改交易状态
	selling.Buyer = buyer
	selling.SellingStatus = model.SellingStatusConstant()["delivery"]
	if err := utils.WriteLedger(selling, stub, model.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
		return shim.Error(fmt.Sprintf("将buyer写入交易selling,修改交易状态 失败%s", err))
	}
	createTime, _ := stub.GetTxTimestamp()
	//将本次购买交易写入账本,可供买家查询
	sellingBuy := &model.SellingBuy{
		Buyer:      buyer,
		CreateTime: time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		Selling:    selling,
	}
	if err := utils.WriteLedger(sellingBuy, stub, model.SellingBuyKey, []string{sellingBuy.Buyer, sellingBuy.CreateTime}); err != nil {
		return shim.Error(fmt.Sprintf("将本次购买交易写入账本失败%s", err))
	}
	sellingBuyByte, err := json.Marshal(sellingBuy)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	//购买成功，扣取余额，更新账本余额，注意，此时需要卖家确认收款，款项才会转入卖家账户，此处先扣除买家的余额
	buyerAccount.Balance -= selling.Price
	if err := utils.WriteLedger(buyerAccount, stub, model.AccountKey, []string{buyerAccount.AccountId}); err != nil {
		return shim.Error(fmt.Sprintf("扣取买家余额失败%s", err))
	}
	// 成功返回
	return shim.Success(sellingBuyByte)
}
