package model

// Account 账户，虚拟管理员和若干业主账号
type Account struct {
	AccountId string  `json:"accountId"` //账号ID
	UserName  string  `json:"userName"`  //账号名
}

type SongRong1 struct {
	SongRongID string  `json:"songrongeId"` //本批松茸ID
	SellerID   string  `json:"sellerID"`   //采购商ID
	Place  string    `json:"place"`  //松茸产地
	Amount   float64 `json:"amount"`    //售卖总量
	Time  string `json:"time"`  //售卖时间
	SellingStatus string  `json:"sellingStatus"` //销售状态
	BuyerID string  `json:"buyerID"` //购买者
}

type SongRong2 struct {
	SongRong2ID string  `json:"songrongeId"` //分装后松茸ID
	Umbrella   string  `json:"umbrella"`   //开伞情况
	Place  string    `json:"place"`  //松茸产地
	Size   string `json:"size"`    //松茸大小
	Storage string `json:"storage"`    //库存时状态
	Time1  string `json:"time"`  //入库时间
	Time2  string `json:"time"`  //出库时间
}

type BaoZhuang struct {
	ProductID string  `json:"productId"` //产品号ID
	SongRong2ID string  `json:"songrongeId"` //产品装的松茸的ID
	Time  string `json:"time"`  //包装时间
}

// SellingBuy 买家参与销售
// 销售对象不能是买家发起的
// Buyer和CreateTime作为复合键,保证可以通过buyer查询到名下所有参与的销售
type SellingBuy struct {
	BuyerID      string  `json:"buyerOD"`      //参与销售人、买家(买家AccountId)
	CreateTime string  `json:"createTime"` //创建时间
	Selling    SongRong1 `json:"selling"`    //销售对象
}
/*
// RealEstate 房地产作为担保出售、捐赠或质押时Encumbrance为true，默认状态false。
// 仅当Encumbrance为false时，才可发起出售、捐赠或质押
// Proprietor和RealEstateID一起作为复合键,保证可以通过Proprietor查询到名下所有的房产信息
type RealEstate struct {
	RealEstateID string  `json:"realEstateId"` //房地产ID
	Proprietor   string  `json:"proprietor"`   //所有者(业主)(业主AccountId)
	Encumbrance  bool    `json:"encumbrance"`  //是否作为担保
	TotalArea    float64 `json:"totalArea"`    //总面积
	LivingSpace  float64 `json:"livingSpace"`  //生活空间
}
// Selling 销售要约
// 需要确定ObjectOfSale是否属于Seller
// 买家初始为空
// Seller和ObjectOfSale一起作为复合键,保证可以通过seller查询到名下所有发起的销售
type Selling struct {
	ObjectOfSale  string  `json:"objectOfSale"`  //销售对象(正在出售的房地产RealEstateID)
	Seller        string  `json:"seller"`        //发起销售人、卖家(卖家AccountId)
	Buyer         string  `json:"buyer"`         //参与销售人、买家(买家AccountId)
	Price         float64 `json:"price"`         //价格
	CreateTime    string  `json:"createTime"`    //创建时间
	SalePeriod    int     `json:"salePeriod"`    //智能合约的有效期(单位为天)
	SellingStatus string  `json:"sellingStatus"` //销售状态
}

// SellingStatusConstant 销售状态
var SellingStatusConstant = func() map[string]string {
	return map[string]string{
		"saleStart": "销售中", //正在销售状态,等待买家光顾
		"cancelled": "已取消", //被卖家取消销售或买家退款操作导致取消
		"expired":   "已过期", //销售期限到期
		"delivery":  "交付中", //买家买下并付款,处于等待卖家确认收款状态,如若卖家未能确认收款，买家可以取消并退款
		"done":      "完成",  //卖家确认接收资金，交易完成
	}
}

// Donating 捐赠要约
// 需要确定ObjectOfDonating是否属于Donor
// 需要指定受赠人Grantee，并等待受赠人同意接收
type Donating struct {
	ObjectOfDonating string `json:"objectOfDonating"` //捐赠对象(正在捐赠的房地产RealEstateID)
	Donor            string `json:"donor"`            //捐赠人(捐赠人AccountId)
	Grantee          string `json:"grantee"`          //受赠人(受赠人AccountId)
	CreateTime       string `json:"createTime"`       //创建时间
	DonatingStatus   string `json:"donatingStatus"`   //捐赠状态
}

// DonatingStatusConstant 捐赠状态
var DonatingStatusConstant = func() map[string]string {
	return map[string]string{
		"donatingStart": "捐赠中", //捐赠人发起捐赠合约，等待受赠人确认受赠
		"cancelled":     "已取消", //捐赠人在受赠人确认受赠之前取消捐赠或受赠人取消接收受赠
		"done":          "完成",  //受赠人确认接收，交易完成
	}
}

// DonatingGrantee 供受赠人查询的
type DonatingGrantee struct {
	Grantee    string   `json:"grantee"`    //受赠人(受赠人AccountId)
	CreateTime string   `json:"createTime"` //创建时间
	Donating   Donating `json:"donating"`   //捐赠对象
}

const (
	AccountKey         = "account-key"
	RealEstateKey      = "real-estate-key"
	SellingKey         = "selling-key"
	SellingBuyKey      = "selling-buy-key"
	DonatingKey        = "donating-key"
	DonatingGranteeKey = "donating-grantee-key"
)
*/

const (
	AccountKey         = "account-key"
	SellSongrongKey      = "sell-songrong-key"
	
	SellingKey         = "selling-key"
	SellingBuyKey      = "selling-buy-key"
	DonatingKey        = "donating-key"
	DonatingGranteeKey = "donating-grantee-key"
)
