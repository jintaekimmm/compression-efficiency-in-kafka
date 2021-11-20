package models

type Music struct {
	TotalRows      string `json:"-"`
	BuySeq         string `json:"buy_seq"`
	ProductCd      string `json:"product_cd"`
	BuyId          string `json:"buy_id"`
	ProductSeq     string `json:"product_seq"`
	ProductNM      string `json:"product_nm"`
	BuyNM          string `json:"buy_nm"`
	BeforeMyMoney  string `json:"before_mymoney"`
	SendStatus     string `json:"send_status"`
	ProductCnt     string `json:"product_cnt"`
	SendCnt        string `json:"send_cnt"`
	DenyCnt        string `json:"deny_cnt"`
	ReceiveCnt     string `json:"receive_cnt"`
	SendTm         string `json:"send_tm"`
	ProductMyMoney string `json:"product_mymoney"`
	BuyMyMoney     string `json:"buy_mymoney"`
	DiscountRate   string `json:"discount_rate"`
	BuyIpAddr      string `json:"buy_ipaddr"`
	BuyLoginId     string `json:"buy_login_id"`
	BuyLoginName   string `json:"buy_login_name"`
	ReceiveTm      string `json:"receive_tm"`
}
