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
	if err != nil || len(resultssongrong1) != 1 {
		return shim.Error(fmt.Sprintf("所要购买的松茸验证失败%s", err))
	}
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
	songrong1.SellingStatus = model.SellingStatusConstant()["delivery"]
	if err := utils.WriteLedger(songrong1, stub, model.SellsongrongKey, []string{songrongID, sellerID}); err != nil {
		return shim.Error(fmt.Sprintf("将buyer写入交易songrong1,修改交易状态为待确认 失败%s", err))
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

// QuerySellingBuyList 查询销售(可查询所有，也可根据发起销售人查询)(发起的)(供卖家查询)
func QuerySellingBuyList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var sellingBuyList []model.SellingBuy
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingBuyKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var sellingBuy model.SellingBuy
			err := json.Unmarshal(v, &sellingBuy)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellingList-反序列化出错: %s", err))
			}
			sellingBuyList = append(sellingBuyList, sellingBuy)
		}
	}
	sellingBuyListByte, err := json.Marshal(sellingBuyList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellingList-序列化出错: %s", err))
	}
	return shim.Success(sellingBuyListByte)
}

// ConfirmSongrong 采购商确认，参数传入采购商ID和售卖的松茸ID
func ConfirmSongrong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	songrongID := args[0]
	sellerID := args[1]
	if songrongID == "" || sellerID == "" {
		return shim.Error("参数存在空值")
	}
	//判断是否采购商操作
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{sellerID})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("操作人权限验证失败%s", err))
	}
	var account model.Account
	if err = json.Unmarshal(resultsAccount[0], &account); err != nil {
		return shim.Error(fmt.Sprintf("查询操作人信息-反序列化出错: %s", err))
	}
	//根据songrongID和sellerID获取想要购买的房产信息，确认存在该房产
	resultssongrong1, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellsongrongKey, []string{songrongID, sellerID})
	if err != nil || len(resultssongrong1) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取想要购买的松茸信息失败: %s", osongrongID, sellerID, err))
	}
	var songrong1 model.SongRong1
	if err = json.Unmarshal(resultssongrong1[0], &songrong1); err != nil {
		return shim.Error(fmt.Sprintf("ConfirmSongrong-反序列化出错: %s", err))
	}
	songrong1.SellingStatus = model.SellingStatusConstant()["confirm"]
	if err := utils.WriteLedger(songrong1, stub, model.SellsongrongKey, []string{songrongID, sellerID}); err != nil {
		return shim.Error(fmt.Sprintf("将confirm写入交易songrong1,修改交易状态为已确认 失败%s", err))
	}
	createTime, _ := stub.GetTxTimestamp()
	//将本次购买交易写入账本,可供买家查询
	sellingConfirm:= &model.SellingConfirm{
		BuyerID:      buyerID,
		CreateTime: time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		Selling:    songrong1,
	}
	if err := utils.WriteLedger(sellingConfirm, stub, model.SellingConfirmKey, []string{sellingConfirm.BuyerID, sellingConfirm.CreateTime}); err != nil {
		return shim.Error(fmt.Sprintf("将本次购买交易写入账本失败%s", err))
	}
	sellingConfirmByte, err := json.Marshal(sellingConfirm)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(sellingConfirmByte)
}

// QuerySellingConfirmList 查询销售(可查询所有，也可根据发起销售人查询)(发起的)(供卖家查询)
func QuerySellingConfirmList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var sellingConfirmList []model.SellingConfirm
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.SellingConfirmKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var sellingConfirm model.SellingConfirm
			err := json.Unmarshal(v, &sellingConfirm)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellingList-反序列化出错: %s", err))
			}
			sellingConfirmList = append(sellingConfirmList, sellingConfirm)
		}
	}
	sellingConfirmListByte, err := json.Marshal(sellingConfirmList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellingList-序列化出错: %s", err))
	}
	return shim.Success(sellingConfirmListByte)
}
