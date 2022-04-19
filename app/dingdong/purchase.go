/*
 * @File: purchase.go
 * @Author: 张志鹏
 * @Version: 1.0
 * @Contact: zhangzhipeng@uniml.com
 * @CreatedTime: 2022/04/15 01:16:48
 * @UpdatedTime: 2022/04/15 01:16:48
 * @Description: 购买逻辑
**/
package dingdong

import (
    "dingdong_hacker/app/assets"
    "fmt"
    "sync"
    "time"
)

// 全局地址ID
var AddressID string
var Address string

// 更新地址ID
func UpdateAddress() error {
    header, body := GetRequestData()
    timeS := fmt.Sprintf("%d", (time.Now().UnixMilli() / 1000))
    header["ddmc-time"] = timeS
    body["time"] = timeS
    signRes := assets.SignForm(body)
    for k := range signRes {
        body[k] = signRes[k]
    }
    adrID, adrDetail, err := GetAddressID(header, body)
    if err != nil {
        return err
    }
    if adrID != "" {
        AddressID = adrID
        Address = adrDetail
        return nil
    } else {
        return fmt.Errorf("更新地址失败")
    }
}

// 订单生成时间
var OrderCreated bool = false
var OrderCreatedTime time.Time
var PurchaseFinished = false
var FinishReason = ""

// 购买指令管道
var PurchaseChan = make(chan int)
var purchaseStart bool = false

// 购买逻辑
func Purchase(concurrentlNum int) {
    fmt.Println("开始抢购")
    quitC := make(chan bool)
    var cartInfo *CartResp
    var reversedTime *ReversedTime
    var orderInfo *CheckOrderResp
    var cartErr, reversedTimeErr, ordInfoErr, ordCrtErr error

    wg := sync.WaitGroup{}

    // 1.轮询购物车
    wg.Add(concurrentlNum)
    for i := 0; i < concurrentlNum; i ++ {
        go func(c chan bool) {
            var newCartInfo *CartResp
            defer wg.Done()
            h, b := GetRequestData()
            for {
                select {
                case quit, ok := <-c:
                    if !ok {
                        return
                    }
                    if quit {
                        return
                    }
                default:
                    newCartInfo, cartErr = GetCart(h, b)
                    // 如果购物车中没有有效商品
                    if cartErr != nil {
                        if IsCartInvalidError(cartErr) {
                            close(quitC)
                            PurchaseFinished = true
                            FinishReason = "没有可以购买的商品"
                        } else if IsExpireError(cartErr) {
                            close(quitC)
                            PurchaseFinished = true
                            FinishReason = "登录失效，请重新填写用户信息"
                        }
                    }
                    if cartInfo == nil && newCartInfo != nil {
                        cartInfo = newCartInfo
                        for _, prod := range cartInfo.Data.NewOrderProductList[0].Products {
                            fmt.Println(prod.ProductName)
                        }
                    }
                    time.Sleep(time.Millisecond * 300)
                }
            }
        }(quitC)
    }

    // 2.轮询配送时间
    wg.Add(concurrentlNum)
    for i := 0; i < concurrentlNum; i++ {
        go func(c chan bool) {
            defer wg.Done()
            h, b := GetRequestData()
            for {
                select {
                case quit, ok := <-c:
                    if !ok {
                        return
                    }
                    if quit {
                        return
                    }
                default:
                    if cartInfo != nil {
                        reversedTime, reversedTimeErr = GetReverseTime(h, b, AddressID, cartInfo)
                        if reversedTimeErr != nil {
                            if IsExpireError(reversedTimeErr) {
                                close(quitC)
                                PurchaseFinished = true
                                FinishReason = "登录失效，请重新填写用户信息"
                            }
                        }
                        if reversedTime != nil {
                            fmt.Println(reversedTime.StartTimeStamp)
                        }
                    }
                    time.Sleep(time.Millisecond * 300)
                }
            }
        }(quitC)
    }

    // 3.轮询订单
    wg.Add(concurrentlNum)
    for i := 0; i < concurrentlNum; i++ {
        go func(c chan bool) {
            defer wg.Done()
            h, b := GetRequestData()
            for {
                select {
                case quit, ok := <-c:
                    if !ok {
                        return
                    }
                    if quit {
                        return
                    }
                default:
                    if reversedTime != nil {
                        orderInfo, ordInfoErr = CheckOrder(h, b, AddressID, cartInfo, reversedTime)
                        if ordInfoErr != nil {
                            if IsExpireError(ordInfoErr) {
                                close(quitC)
                                PurchaseFinished = true
                                FinishReason = "登录失效，请重新填写用户信息"
                            }
                        }
                    }
                    time.Sleep(time.Millisecond * 300)
                }
            }
        }(quitC)
    }

    // // 4.轮询创建订单
    wg.Add(concurrentlNum)
    for i := 0; i < concurrentlNum; i ++ {
        go func(c chan bool) {
            defer wg.Done()
            h, b := GetRequestData()
            for {
                select {
                case quit, ok := <-c:
                    if !ok {
                        return
                    }
                    if quit {
                        return
                    }
                default:
                    if orderInfo != nil  && reversedTime != nil {
                        _, ordCrtErr = CreareOrder(h, b, AddressID, cartInfo, reversedTime, orderInfo)
                        // 如果订单创建成功，则关闭管道
                        if ordCrtErr == nil {
                            OrderCreatedTime = time.Now()
                            close(c)
                            PurchaseFinished = true
                            for i := 0; i < 10; i++ {
                                fmt.Println("订单创建成功")
                            }
                        } else if IsExpireError(ordCrtErr) {
                            close(c)
                            PurchaseFinished = true
                            FinishReason = "登录失效，请重新填写用户信息"
                        } else if IsRiskError(ordCrtErr) {
                            close(c)
                            PurchaseFinished = true
                            FinishReason = ordCrtErr.Error()
                        }
                    }
                    time.Sleep(time.Millisecond * 300)
                }
            }
        }(quitC)
    }

    // 模拟用
    // wg.Add(1)
    // go func(c chan bool) {
    //     defer wg.Done()
    //     for {
    //         select {
    //             case quit, ok := <- c:
    //                 if !ok {
    //                     return
    //                 }
    //                 if quit {
    //                     return
    //                 }
    //             default:
    //                 fmt.Println("正在抢购")
    //                 time.Sleep(time.Second * 10);
    //                 PurchaseFinished = true;
    //                 OrderCreated = true;
    //                 OrderCreatedTime = time.Now()
    //                 close(c)
    //         }
    //     }
    // } (quitC)

    // 5.定时停止
    wg.Add(1)
    go func(c chan bool) {
        defer wg.Done()
        time.Sleep(time.Minute * 2)
        if !PurchaseFinished {
            close(c)
            PurchaseFinished = true
            FinishReason = "长时间没有抢到，为避免风控，停止抢菜"
        }
    }(quitC)

    // 重置购买已开始选项
    purchaseStart = false

    wg.Wait()
}

// 后台购买逻辑
func BackendPurchase(purchaseC chan int) {
    for {
        select {
            // 等待开始指令
            case concurrentNum, ok := <-purchaseC:
                if !ok {
                    return
                }
                // 开始购买
                Purchase(concurrentNum)
            default:
                time.Sleep(time.Millisecond * 100)
        }
    }
    // 结束关闭管道
    // close(purchaseC)
}

// 发送购买指令
func StartPurchase(concurrentNum int) {
    if !purchaseStart {
        PurchaseChan <- concurrentNum
        purchaseStart = true
    }
}
