/*
 * @File: dingdong_api.go
 * @Author: 张志鹏
 * @Version: 1.0
 * @Contact: zhangzhipeng@uniml.com
 * @CreatedTime: 2022/04/15 00:19:27
 * @UpdatedTime: 2022/04/15 00:19:27
 * @Description: 叮咚买菜的API
**/
package dingdong

import (
	"dingdong_hacker/app/assets"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"unsafe"

	// "net/url"
	"strings"

	"github.com/google/uuid"
)

// 地址请求的响应
type AdrResp struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ValidAddress []struct {
			ID         string `json:"id"`
			IsDefault  bool   `json:"is_default"`
			AddrDetail string `json:"addr_detail"`
			StationId  string `json:"station_id"`
		} `json:"valid_address"`
	} `json:"data"`
}

// GetAddressID 获取用户的默认地址ID
func GetAddressID(header, queryStr map[string]string) (string, string, error) {
	// 1.构造get请求
	req, _ := http.NewRequest(http.MethodGet, "https://sunquan.api.ddxq.mobi/api/v1/user/address/", nil)
	// 设置请求头
	for key := range header {
		req.Header.Add(key, header[key])
	}
	// 设置query参数
	query := req.URL.Query()
	for key := range queryStr {
		query.Add(key, queryStr[key])
	}
	req.URL.RawQuery = query.Encode()

	// 2.发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("创建地址请求失败: %s\n", err.Error())
		return "", "", fmt.Errorf("创建地址请求失败: %s", err.Error())
	}
	respData, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// 3.解析请求
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("请求地址失败，%s\n", respData)
		return "", "", fmt.Errorf("请求地址失败，%s", respData)
	}
	var adrResp AdrResp
	err = json.Unmarshal(respData, &adrResp)
	if err != nil {
		fmt.Printf("解析地址响应失败, %s\n", err.Error())
		return "", "", fmt.Errorf("解析地址响应失败, %s", err.Error())
	}

	// 4.比对参数
	if !adrResp.Success {
		return "", "", fmt.Errorf("请求失败, %s，请更新抓包数据", adrResp.Message)
	}
	var defaultAdrID, defaultAddrDetail string
	for _, adr := range adrResp.Data.ValidAddress {
		if adr.IsDefault {
			if adr.StationId != header["ddmc-station-id"] {
				fmt.Printf("参数中的stationId和响应不一致, %s\n", adr.StationId)
			}
			defaultAdrID = adr.ID
			defaultAddrDetail = adr.AddrDetail
		} else {
			continue
		}
	}
	return defaultAdrID, defaultAddrDetail, nil
}

// 响应通用结构
type NormalResp struct {
	Msg string `json:"msg"`
}

// 据响应判断是否过期
func IsExpireResponse(respData []byte) bool {
	var resp NormalResp
	err := json.Unmarshal(respData, &resp)
	if err == nil && resp.Msg == "您的访问已过期" {
		return true
	}
	return false
}

// 商品结构
type CartProducts struct {
	ID                 string `json:"id"`
	CartID             string `json:"cart_id"`
	CategoryPath       string `json:"category_path"`
	IsBooking          int    `json:"is_booking"`
	ProductName        string `json:"product_name"`
	Count              int    `json:"count"`
	Price              string `json:"price"`
	TotalMoney         string `json:"total_money"`
	TotalPrice         string `json:"total_price"`
	InstantRebateMoney string `json:"instant_rebate_money"`
	ActivityID         string `json:"activity_id"`
	ConditionsNum      string `json:"conditions_num"`
	ProductType        int    `json:"product_type"`
	// ProductName        *string       `json:"product_name"`
	Sizes            []interface{} `json:"sizes"`
	Type             int           `json:"type"`
	TotalOriginMoney string        `json:"total_origin_money"`
	TotalOriginPrice string        `json:"total_origin_price"`
	PriceType        int           `json:"price_type"`
	BatchType        int           `json:"batch_type"`
	SubList          []interface{} `json:"sub_list"`
	OrderSort        int           `json:"order_sort"`
	OriginPrice      string        `json:"origin_price"`
}

// 通用购物车订单数据
type CartCommonInfo struct {
	Products           []CartProducts `json:"products"`
	TotalMoney         string         `json:"total_money"`
	TotalOriginMoney   string         `json:"total_origin_money"`
	GoodsRealMoney     string         `json:"goods_real_money"`
	TotalCount         int            `json:"total_count"`
	CartCount          int            `json:"cart_count"`
	IsPresale          int            `json:"is_presale"`
	InstantRebateMoney string         `json:"instant_rebate_money"`
	// CouponRebateMoney      string        `json:"coupon_rebate_money"`
	TotalRebateMoney       string        `json:"total_rebate_money"`
	UsedBalanceMoney       string        `json:"used_balance_money"`
	CanUsedBalanceMoney    string        `json:"can_used_balance_money"`
	UsedPointNum           int           `json:"used_point_num"`
	UsedPointMoney         string        `json:"used_point_money"`
	CanUsedPointNum        int           `json:"can_used_point_num"`
	CanUsedPointMoney      string        `json:"can_used_point_money"`
	IsSahreStation         int           `json:"is_share_station"`
	OnlyTodayProducts      []interface{} `json:"only_today_products"`
	OnlyTomorrowProducts   []interface{} `json:"only_tomorrow_products"`
	PackageType            int           `json:"package_type"`
	PackageID              int           `json:"package_id"`
	FrontPackageText       string        `json:"front_package_text"`
	FrontPackageType       int           `json:"front_package_type"`
	FrontPackageStockColor string        `json:"front_package_stock_color"`
	FrontPackageDGColor    string        `json:"front_package_bg_color"`
}

// 购物车查询响应
type CartResp struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
	Data    struct {
		NewOrderProductList []CartCommonInfo `json:"new_order_product_list"`
		ParentOrderInfo     struct {
			ParentOrderSign string `json:"parent_order_sign"`
		} `json:"parent_order_info"`
	} `json:"data"`
}

type CartInfo struct {
	Data struct {
		NewOrderProductList []CartCommonInfo `json:"new_order_product_list"`
		ParentOrderInfo     struct {
			ParentOrderSign string `json:"parent_order_sign"`
		} `json:"parent_order_info"`
	} `json:"data"`
}

// 获取购物车信息
//
// @return: 购物车信息
// @return: 查询购物车是否成功
func GetCart(header, queryStr map[string]string) (*CartResp, error) {
	// 1.构造get请求
	// 设置请求时间和加密
	timeS := fmt.Sprintf("%d", time.Now().UnixMilli()/1000)
	header["ddmc-time"] = timeS
	queryStr["time"] = timeS
	queryStr["is_load"] = "1"
	queryStr["ab_config"] = "{\"key_onion\":\"D\",\"key_cart_discount_price\":\"C\"}"

	// 加密
	m := assets.SignForm(queryStr)
	for k := range m {
		queryStr[k] = m[k]
	}

	req, _ := http.NewRequest(http.MethodGet, "https://maicai.api.ddxq.mobi/cart/index", nil)
	// 设置请求头
	for key := range header {
		req.Header.Add(key, header[key])
	}
	// 设置query参数
	query := req.URL.Query()
	for key := range queryStr {
		query.Add(key, queryStr[key])
	}

	req.URL.RawQuery = query.Encode()

	// 2.发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("创建购物车请求失败: %s\n", err.Error())
		return nil, fmt.Errorf("创建购物车请求失败: %s", err.Error())
	}
	respData, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// 3.解析请求
	if resp.StatusCode != http.StatusOK {
		if IsExpireResponse(respData) {
			return nil, ExpireError{}
		}
		fmt.Printf("请求购物车失败，%s\n", respData)
		return nil, fmt.Errorf("请求购物车失败，%s", respData)
	}
	var cartResp CartResp
	err = json.Unmarshal(respData, &cartResp)
	if err != nil {
		fmt.Printf("解析购物车响应失败, %s\n%s\n", err.Error(), respData)
		return nil, fmt.Errorf("解析购物车响应失败, %s", err.Error())
	}

	// 4.比对参数
	if len(cartResp.Data.NewOrderProductList) == 0 {
		return nil, CartInvalidError{}
	}
	for idx, product := range cartResp.Data.NewOrderProductList[0].Products {
		cartResp.Data.NewOrderProductList[0].Products[idx].TotalMoney = product.TotalPrice
		cartResp.Data.NewOrderProductList[0].Products[idx].TotalOriginMoney = product.TotalOriginPrice
	}
	return &cartResp, nil
}

// 配送时间
type ReversedTime struct {
	EndTimestamp   int64   `json:"end_timestamp"`
	StartTimeStamp int64   `json:"start_timestamp"`
	DisableType    *int    `json:"disableType,omitempty"`
	SelectMsg      *string `json:"select_msg,omitempty"`
}

// 配送时间响应
type ReverseTimeResp struct {
	Data []struct {
		Time []struct {
			Times []ReversedTime `json:"times"`
		} `json:"time"`
	} `json:"data"`
}

// 获取配送时间
//
// @return: 配送时间，是否有有效的配送时间
func GetReverseTime(
	header, queryStr map[string]string,
	adrID string,
	cartResp *CartResp,
) (*ReversedTime, error) {
	// 1.构造请求
	// 设置请求时间和加密
	timeS := fmt.Sprintf("%d", time.Now().UnixMilli()/1000)
	header["ddmc-time"] = timeS
	queryStr["time"] = timeS

	// 专属字段
	queryStr["address_id"] = adrID
	queryStr["group_config_id"] = ""
	queryStr["isBridge"] = "false"
	productsBytes, _ := json.Marshal(cartResp.Data.NewOrderProductList[0].Products)
	queryStr["products"] = fmt.Sprintf("[%s]", *(*string)(unsafe.Pointer(&productsBytes)))

	// 加密
	m := assets.SignForm(queryStr)
	for k := range m {
		queryStr[k] = m[k]
	}

	// www-form-urlencoded参数
	formSlice := make([]string, len(queryStr))
	idx := 0
	for key := range queryStr {
		formSlice[idx] = fmt.Sprintf("%s=%s", key, queryStr[key])
		idx += 1
	}
	payload := strings.NewReader(strings.Join(formSlice, "&"))

	req, _ := http.NewRequest(http.MethodPost, "https://maicai.api.ddxq.mobi/order/getMultiReserveTime", payload)
	// 设置请求头
	for key := range header {
		req.Header.Add(key, header[key])
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 2.发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("创建配送时间请求失败: %s\n", err.Error())
		return nil, fmt.Errorf("创建配送时间请求失败: %s", err.Error())
	}
	respData, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// 3.解析请求
	if resp.StatusCode != http.StatusOK {
		if IsExpireResponse(respData) {
			return nil, ExpireError{}
		}
		fmt.Printf("请求配送时间失败，%s\n", respData)
		return nil, fmt.Errorf("请求配送时间失败，%s", respData)
	}

	var reverseTimeResp ReverseTimeResp
	err = json.Unmarshal(respData, &reverseTimeResp)
	if err != nil {
		fmt.Printf("解析配送时间按响应失败, %s\n", err.Error())
		return nil, fmt.Errorf("解析配送时间按响应失败, %s", err.Error())
	}

	// 4.比对参数
	var reverseTime ReversedTime
	for _, time := range reverseTimeResp.Data[0].Time[0].Times {
		if *time.DisableType == 0 && !strings.Contains(*time.SelectMsg, "尽快") {
			reverseTime = ReversedTime{
				EndTimestamp:   time.EndTimestamp,
				StartTimeStamp: time.StartTimeStamp,
				DisableType:    new(int),
				SelectMsg:      time.SelectMsg,
			}
			return &reverseTime, nil
		}
	}
	fmt.Println("没有可用的配送时间")
	return &reverseTime, fmt.Errorf("没有可用的配送时间")
}

// 生成订单响应
type CheckOrderResp struct {
	Message string `json:"msg"`
	Success bool   `json:"success"`
	Data    struct {
		Order struct {
			FreightDiscountMoney string `json:"freight_discount_money"`
			TotalMoney           string `json:"total_money"`
			FreightMoney         string `json:"freight_money"`
			Freights             []struct {
				Freight struct {
					FreightRealMoney string `json:"freight_real_money"`
				} `json:"freight"`
			} `json:"freights"`
			DefaultCoupon struct {
				ID string `json:"_id"`
			} `json:"default_coupon"`
		} `json:"order"`
	} `json:"data"`
}

// 提交订单需要的package数据
type Package struct {
	CartCommonInfo
	REversedTime ReversedTime `json:"reserved_time"`
}

// 生成订单
//
// @return: 订单确认信息
// @return: 订单确认是否成功, true: 成功, false: 失败
func CheckOrder(
	header, queryStr map[string]string,
	adrID string,
	cartResp *CartResp,
	reverseTime *ReversedTime,
) (*CheckOrderResp, error) {
	// 1.构造请求
	// 设置请求时间和加密
	timeS := fmt.Sprintf("%d", time.Now().UnixMilli()/1000)
	header["ddmc-time"] = timeS
	queryStr["time"] = timeS

	// 专属字段
	queryStr["address_id"] = adrID
	queryStr["user_ticket_id"] = "default"
	queryStr["freight_ticket_id"] = "default"
	queryStr["is_use_point"] = "0"
	queryStr["is_use_balance"] = "0"
	queryStr["is_buy_vip"] = "0"
	queryStr["coupons_id"] = ""
	queryStr["is_buy_coupons"] = "0"
	queryStr["check_order_type"] = "0"
	queryStr["is_support_merge_payment"] = "1"
	queryStr["showData"] = "true"
	queryStr["showMsg"] = "false"
	// 添加请求体中的packages字段
	cartInfo := cartResp.Data.NewOrderProductList[0]
	packages := make([]Package, 1)
	packages[0] = Package{
		CartCommonInfo: cartInfo,
		REversedTime: ReversedTime{
			EndTimestamp:   reverseTime.EndTimestamp,
			StartTimeStamp: reverseTime.StartTimeStamp,
		},
	}
	pkgBytes, _ := json.Marshal(packages)
	queryStr["packages"] = *(*string)(unsafe.Pointer(&pkgBytes))

	// 加密
	m := assets.SignForm(queryStr)
	for k := range m {
		queryStr[k] = m[k]
	}

	// www-form-urlencoded参数
	formSlice := make([]string, len(queryStr))
	idx := 0
	for key := range queryStr {
		formSlice[idx] = fmt.Sprintf("%s=%s", key, queryStr[key])
		idx += 1
	}
	payload := strings.NewReader(strings.Join(formSlice, "&"))

	req, _ := http.NewRequest(http.MethodPost, "https://maicai.api.ddxq.mobi/order/checkOrder", payload)
	// 设置请求头
	for key := range header {
		req.Header.Add(key, header[key])
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 2.发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("创建订单请求失败: %s\n", err.Error())
		return nil, fmt.Errorf("创建订单请求失败: %s", err.Error())
	}
	respData, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// 3.解析请求
	if resp.StatusCode != http.StatusOK {
		if IsExpireResponse(respData) {
			return nil, ExpireError{}
		}
		fmt.Printf("请求创建订单失败，%s\n", respData)
		return nil, fmt.Errorf("请求创建订单失败，%s", respData)
	}
	var checkOrderResp CheckOrderResp
	err = json.Unmarshal(respData, &checkOrderResp)
	if err != nil {
		fmt.Printf("解析创建订单响应失败, %s\n", err.Error())
		return nil, fmt.Errorf("解析创建订单响应失败, %s", err.Error())
	} else if !checkOrderResp.Success {
		fmt.Printf("请求订单失败: %s", checkOrderResp.Message)
		return nil, fmt.Errorf("请求订单失败: %s", checkOrderResp.Message)
	}

	// 4.返回
	return &checkOrderResp, nil
}

// PaymentInfo 支付信息
type PaymentInfo struct {
	ReservedTimeStart    int64  `json:"reserved_time_start"`
	ReservedTimeEnd      int64  `json:"reserved_time_end"`
	Price                string `json:"price"`
	FreightDiscountMoney string `json:"freight_discount_money"`
	FreightMoney         string `json:"freight_money"`
	OrderFreight         string `json:"order_freight"`
	ParentOrderSign      string `json:"parent_order_sign"`
	ProductType          int    `json:"product_type"` // 1
	AddressID            string `json:"address_id"`
	FormID               string `json:"form_id"`             // UUID
	ReceiptWithoutSKU    *int   `json:"receipt_without_sku"` // null
	PayType              int    `json:"pay_type"`            // 6
	UserTicketID         string `json:"user_ticket_id"`
	VIPMoney             string `json:"vip_money"`              // ""
	VIPBuyUserTicketID   string `json:"vip_buy_user_ticket_id"` // ""
	CouponsMoney         string `json:"coupons_money"`          // ""
	CouponsID            string `json:"coupons_id"`             // ""
}

// 创建订单时的商品信息
type OrderPackage struct {
	CartCommonInfo
	EndTimestamp         int64  `json:"reserved_time_end"`
	StartTimeStamp       int64  `json:"reserved_time_start"`
	ETATraceID           string `json:"eta_trace_id"`            // ""
	SoonArrival          string `json:"soon_arrival"`            // ""
	FirstSelectedBigTime int    `json:"first_selected_big_time"` // 0
	ReceiptWithoutSKU    int    `json:"receipt_without_sku"`     // 0
}

type PackageOrder struct {
	PaymentOrder PaymentInfo    `json:"payment_order"`
	Packages     []OrderPackage `json:"packages"`
}

// 提交订单请求响应
type CreateOrderResp struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
	Data    struct {
		TradeTag string `json:"tradeTag,omitempty"`
		PayURL   string `json:"pay_url,omitempty"`
	} `json:"data"`
}

// CreateOrder 创建新订单
//
// @return: 新建订单响应
// @return: 订单是否创建成功，true: 创建成功; false: 创建失败
func CreareOrder(
	header, queryStr map[string]string,
	adrID string,
	cartResp *CartResp,
	reverseTime *ReversedTime,
	orderInfo *CheckOrderResp,
) (*CreateOrderResp, error) {
	// 1.设置请求时间和加密
	timeS := fmt.Sprintf("%d", time.Now().UnixMilli()/1000)
	header["ddmc-time"] = timeS
	queryStr["time"] = timeS

	// 专属字段
	queryStr["showData"] = "true"
	queryStr["showMsg"] = "false"
	queryStr["ab_config"] = "{\"key_onion\":\"C\"}"
	// 添加请求体中的package_order字段
	UID, _ := uuid.NewUUID()
	formUID := strings.Replace(UID.String(), "-", "", -1)
	pkgOdr := PackageOrder{
		PaymentOrder: PaymentInfo{
			ReservedTimeStart:    reverseTime.StartTimeStamp,
			ReservedTimeEnd:      reverseTime.EndTimestamp,
			Price:                orderInfo.Data.Order.TotalMoney,
			FreightDiscountMoney: orderInfo.Data.Order.FreightDiscountMoney,
			FreightMoney:         orderInfo.Data.Order.FreightMoney,
			OrderFreight:         orderInfo.Data.Order.Freights[0].Freight.FreightRealMoney,
			ParentOrderSign:      cartResp.Data.ParentOrderInfo.ParentOrderSign,
			ProductType:          1,
			AddressID:            adrID,
			FormID:               formUID,
			ReceiptWithoutSKU:    nil,
			PayType:              6,
			UserTicketID:         orderInfo.Data.Order.DefaultCoupon.ID,
			VIPMoney:             "",
			VIPBuyUserTicketID:   "",
			CouponsMoney:         "",
			CouponsID:            "",
		},
		Packages: []OrderPackage{{
			CartCommonInfo:       cartResp.Data.NewOrderProductList[0],
			StartTimeStamp:       reverseTime.StartTimeStamp,
            EndTimestamp:         reverseTime.EndTimestamp,
			ETATraceID:           "",
			SoonArrival:          "",
			FirstSelectedBigTime: 0,
			ReceiptWithoutSKU:    0,
		}},
	}
	pkgOdrByte, _ := json.Marshal(pkgOdr)
	queryStr["package_order"] = *(*string)(unsafe.Pointer(&pkgOdrByte))

	// 加密
	m := assets.SignForm(queryStr)
	for k := range m {
		queryStr[k] = m[k]
	}

	// www-form-urlencoded参数
	formSlice := make([]string, len(queryStr))
	idx := 0
	for key := range queryStr {
		formSlice[idx] = fmt.Sprintf("%s=%s", key, queryStr[key])
		idx += 1
	}
	payload := strings.NewReader(strings.Join(formSlice, "&"))

	req, _ := http.NewRequest(http.MethodPost, "https://maicai.api.ddxq.mobi/order/addNewOrder", payload)
	// 设置请求头
	for key := range header {
		req.Header.Add(key, header[key])
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 2.发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("新建订单请求失败: %s\n", err.Error())
		return nil, fmt.Errorf("新建订单请求失败: %s", err.Error())
	}
	respData, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// 3.解析请求
	if resp.StatusCode != http.StatusOK {
		if IsExpireResponse(respData) {
			return nil, ExpireError{}
		}
		fmt.Printf("请求新建订单失败，%s\n", respData)
		return nil, fmt.Errorf("请求新建订单失败，%s", respData)
	}

	var crtOdrResp CreateOrderResp
	err = json.Unmarshal(respData, &crtOdrResp)
	if err != nil {
		fmt.Printf("解析创建订单响应失败, %s\n", err.Error())
		return nil, fmt.Errorf("解析创建订单响应失败, %s", err.Error())
	}

	// 4.对比数据
	// 4.1.如果下单成功
	if !crtOdrResp.Success {
		if crtOdrResp.Data.TradeTag == "RISK_ORDER_BEFORE" {
			return nil, RiskError{}
		}
		return nil, fmt.Errorf("创建订单失败, %s", crtOdrResp.Message)
	}
	if len(crtOdrResp.Data.PayURL) > 0 {
		return &crtOdrResp, nil
	} else {
		return &crtOdrResp, fmt.Errorf("订单创建失败")
	}
}
