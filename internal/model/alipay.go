package model

// AlipayTradeStatus provides enum type for alipay trade status const
type AlipayTradeStatus string

const (
	// AlipayWaitBuyerPay Trade created and waiting for buyer to pay
	AlipayWaitBuyerPay AlipayTradeStatus = "WAIT_BUYER_PAY"

	// AlipayTradeClosed Trade closed due to time out or refund after pay
	AlipayTradeClosed AlipayTradeStatus = "TRADE_CLOSED"

	// AlipayTradeSuccess Trade successes
	AlipayTradeSuccess AlipayTradeStatus = "TRADE_SUCCESS"

	// AlipayTradeFinished Trade finished
	AlipayTradeFinished AlipayTradeStatus = "TRADE_FINISHED"
)

// FundBill provides basic struct for bill
// type FundBill struct {
// 	Amount      string `json:"amount"`
// 	FundChannel string `json:"fundChannel"`
// }

// VoucherDetail provides basic struct for voucher
// type VoucherDetail struct {
// 	Amount             string `json:"amount"`
// 	MerchantContribute string `json:"merchantContribute"`
// 	OtherContribute    string `json:"other_contribute"`
// 	Type               string `json:"type"`
// 	Memo               string `json:"memo"`
// }

// NotifyRequest provides full notify request struct
type NotifyRequest struct {
	NotifyTime        string            `json:"notify_time"`
	NotifyType        string            `json:"notify_type"`
	NotifyID          string            `json:"notify_id"`
	AppID             string            `json:"app_id"`
	Charset           string            `json:"charset"`
	Version           string            `json:"version"`
	SignType          string            `json:"sign_type"`
	Sign              string            `json:"sign"`
	TradeNo           string            `json:"trade_no"`
	OutTradeNo        string            `json:"out_trade_no"`
	OutBizNo          string            `json:"out_biz_no,omitempty"`
	BuyerID           string            `json:"buyer_id,omitempty"`
	BuyerLogonID      string            `json:"buyer_logon_id,omitempty"`
	SellerID          string            `json:"seller_id,omitempty"`
	SellerEmail       string            `json:"seller_email,omitempty"`
	TradeStatus       AlipayTradeStatus `json:"trade_status,omitempty"`
	TotalAmount       string            `json:"total_amount,omitempty"`
	ReceiptAmount     string            `json:"receipt_amount,omitempty"`
	InvoiceAmount     string            `json:"invoice_amount,omitempty"`
	BuyerPayAmount    string            `json:"buyer_pay_amount,omitempty"`
	PointAmount       string            `json:"point_amount,omitempty"`
	RefundFee         string            `json:"refund_fee,omitempty"`
	Subject           string            `json:"subject,omitempty"`
	Body              string            `json:"body,omitempty"`
	GmtCreate         string            `json:"gmt_create,omitempty"`
	GmtPayment        string            `json:"gmt_payment,omitempty"`
	GmtRefund         string            `json:"gmt_refund,omitempty"`
	GmtClose          string            `json:"gmt_close,omitempty"`
	FundBillList      string            `json:"fund_bill_list,omitempty"`
	PassbackParams    string            `json:"passback_params,omitempty"`
	VoucherDetailList string            `json:"voucher_detail_list,omitempty"`
}
