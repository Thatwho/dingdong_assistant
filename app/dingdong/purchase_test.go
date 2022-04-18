/*
 * @File: test_purchase.go
 * @Author: 张志鹏
 * @Version: 1.0
 * @Contact: zhangzhipeng@uniml.com
 * @CreatedTime: 2022/04/17 22:24:14
 * @UpdatedTime: 2022/04/17 22:24:14
 * @Description: 测试购买逻辑
**/
package dingdong

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func SimulateData() {
    Header = DingDongHeader{
        CityNumber :   "0101",
        BuildVersion : "2.83.0",
        DeviceID :     "osP8I0Y_61BdE-UqmuBc6aOQnoVA",
        StationID :    "5bc5a951716de1a94f8b6fbc",
        DChannel :     "applet",
        OSVersion :    "[object Undefined]",
        AppClientID :  "4",
        Cookie :       "DDXQSESSID=8bd3177221dd07d568c88976d4f0303c",
        IP :           "",
        Longitude :    "121.51597",
        Latitude :     "31.296135",
        ApiVersion :   "9.50.0",
        UID :          "5ff887ebbef9980001c364da",
        UserAgent :    "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat",
        Referer :      "https://servicewechat.com/wx1e113254eda17715/425/page-frame.html",
    }
    Body = DingDongBody{
        UID :          "5ff887ebbef9980001c364da",
        Longitude :    "121.51597",
        Latitude :     "31.296135",
        StationID :    "5bc5a951716de1a94f8b6fbc",
        CityNumber :   "0101",
        ApiVersion :   "9.50.0",
        AppVersion :   "2.83.0",
        AppletSource : "",
        Channel :      "applet",
        AppClientID :  "4",
        ShareUID :     "",
        OpenID :       "osP8I0Y_61BdE-UqmuBc6aOQnoVA",
        H5Source :     "",
        DeviceToken :  "WHJMrwNw1k/FKPjcOOgRd+OIAhs278G2+ynPr8kCRZg2Ip9Q8WOFIcoldDOWZxX1XDCEj5JfbCBLkkj8POSEvJgEvU03AoKnBdCW1tldyDzmauSxIJm5Txg==1487582755342",
        SID :          "8bd3177221dd07d568c88976d4f0303c",
    }
}

func TestPurchase(t *testing.T) {
    SimulateData()
    UpdateAddress()
    var cartInfo *CartResp
    var reversedTime *ReversedTime
    var orderInfo *CheckOrderResp
    var cartErr, reversedTimeErr, ordInfoErr, ordCrtErr error

    // 1.轮询购物车
    h, b := GetRequestData()
    cartInfo, cartErr = GetCart(h, b)
    // 如果购物车中没有有效商品
    if IsCartInvalidError(cartErr) {
        panic("没有商品")
    } else if IsExpireError(cartErr) {
        panic("登录过期")
    } else if cartErr != nil {
        panic(cartErr)
    }
    // fmt.Println(cartInfo.Data.NewOrderProductList[0].Products[0].ProductName)
    time.Sleep(time.Millisecond * 300)

    // 2.轮询配送时间
    if cartInfo != nil {
        h, b = GetRequestData()
        reversedTime, reversedTimeErr = GetReverseTime(h, b, AddressID, cartInfo)
        if IsExpireError(reversedTimeErr) {
            panic("登录过期")
        } else if reversedTimeErr != nil {
            panic(reversedTimeErr)
        }
    }

    // 3.轮询购物车
    if reversedTime != nil {
        h, b = GetRequestData()
        orderInfo, ordInfoErr = CheckOrder(h, b, AddressID, cartInfo, reversedTime)
        if IsExpireError(ordInfoErr) {
            t.Error("登录失效，请重新填写用户信息")
        } else if ordInfoErr != nil {
            t.Error(ordInfoErr)
        }
    }

    // 4.加购
    if orderInfo != nil {
        h, b := GetRequestData()
        _, ordCrtErr = CreareOrder(h, b, AddressID, cartInfo, reversedTime, orderInfo)
        // 如果订单创建成功，则关闭管道
        if ordCrtErr == nil {
            OrderCreatedTime = time.Now()
            for i:=0; i<10; i++ {    
                fmt.Println("订单创建成功")
            }
        } else if IsExpireError(ordCrtErr) {
            t.Error("登录失效，请重新填写用户信息")
        } else {
            t.Error(ordCrtErr)
        }
    }
}

func TestCheckOrder(t *testing.T) {
    b, _ := os.ReadFile("/home/uniml/person/golang/dingdong_hacker/data/checkOrder.json")
    var chOdrResp CheckOrderResp
    err := json.Unmarshal(b, &chOdrResp)
    if err != nil {
        t.Error(err)
    }
    t.Log(chOdrResp.Data.Order.Freights[0].Freight.FreightRealMoney)
}

func TestOrderPackages(t *testing.T) {
    var pkgs []Package
    b, _ := os.ReadFile("/home/uniml/person/golang/dingdong_hacker/data/checkOrderPackages.json")
    err := json.Unmarshal(b, &pkgs)
    if err != nil {
        t.Error(err)
    }
    // fmt.Println(pkgs[0].Products[0].ProductName)
}