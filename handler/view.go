/*
 * @File: view.go
 * @Author: 张志鹏
 * @Version: 1.0
 * @Contact: zhangzhipeng@uniml.com
 * @CreatedTime: 2022/04/16 02:12:50
 * @UpdatedTime: 2022/04/16 02:12:50
 * @Description: http入口
**/
package handler

import (
	"bytes"
	"dingdong_hacker/app/assets"
	"dingdong_hacker/app/dingdong"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Msg string `json:"msg"`
}

type AdrResp struct {
	ID     string `json:"id"`
	Detail string `json:"detail"`
}

// 更新地址ID接口
func UpdateAddress(c *gin.Context) {
	err := dingdong.UpdateAddress()
	if err == nil {
		c.JSON(200, AdrResp{
			ID:     dingdong.AddressID,
			Detail: dingdong.Address,
		})
	} else {
		c.JSON(500, Resp{
			Msg: err.Error(),
		})
	}
}

type UserConfig struct {
	Header dingdong.DingDongHeader `json:"header"`
	Body   dingdong.DingDongBody   `json:"body"`
}

// 更新用户配置
func UpdateUserConfig(c *gin.Context) {
	var userConf UserConfig
	c.BindJSON(&userConf)
	dingdong.UpdateHeader(userConf.Header)
	dingdong.UpdateBody(userConf.Body)
	c.JSON(200, nil)
}

// 执行后台购买逻辑
func Purchase(c *gin.Context) {
	dingdong.StartPurchase()
	c.JSON(http.StatusNoContent, nil)
}

// 购买结果数据结构
type OrderResultCheck struct {
	PurchaseFinished bool      `json:"finished"`
	FinishReason     string    `json:"reason"`
	OrderSuccess     bool      `json:"orderSuccess"`
	OrderTime        time.Time `json:"orderTime"`
}

// 查询有无订单生成
func GetOrder(c *gin.Context) {
	c.JSON(http.StatusOK, OrderResultCheck{
		PurchaseFinished: dingdong.PurchaseFinished,
		FinishReason:     dingdong.FinishReason,
		OrderSuccess:     dingdong.OrderCreated,
		OrderTime:        dingdong.OrderCreatedTime,
	})
}

// 静态资源返回函数
func StaticFile(c *gin.Context) {
	fileName := fmt.Sprintf("/static%s", c.Request.URL.Path)
	if file, ok := assets.Assets.Files[fileName]; ok {
		tmpS := strings.Split(c.Request.URL.Path, ".")
		expand := tmpS[len(tmpS)-1]
		fileType := ""
		if expand == "webp" {
			fileType = fmt.Sprintf("image/%s", expand)
		} else if expand == "mp3" {
			fileType = "audio/mpeg"
		} else {
			fileType = fmt.Sprintf("text/%s; charset=utf-8", expand)
		}
		c.DataFromReader(http.StatusOK, file.Size(), fileType, bytes.NewReader(file.Data), nil)
	}
}

var Server *gin.Engine

func init() {
	r := gin.Default()

	r.PUT("/api/v1/address", UpdateAddress)
	r.PUT("/api/v1/user_conf", UpdateUserConfig)
	r.POST("api/v1/purchase", Purchase)
	r.GET("/api/v1/order", GetOrder)
	r.GET("/:staticFile", StaticFile)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/index.html")
	})
	Server = r

	// 开启后台购买指令监听
	go dingdong.BackendPurchase(dingdong.PurchaseChan)
}
