/*
 * @File: sign.go
 * @Author: 张志鹏
 * @Version: 1.0
 * @Contact: zhangzhipeng@uniml.com
 * @CreatedTime: 2022/04/17 20:59:27
 * @UpdatedTime: 2022/04/17 20:59:27
 * @Description: 签名表单
**/
package assets

import (
	"encoding/json"
	"unsafe"

	"github.com/robertkrimen/otto"
)

func SignForm(form map[string]string) map[string]string {
    vm := otto.New()

    _, err := vm.Run(*(*string)(unsafe.Pointer(&Assets.Files["/static/sign.js"].Data)))
    if err != nil {
        panic("解析签名算法失败")
    }

    jsonStr, _ := json.Marshal(form)
    value, err := vm.Call("sign", nil, *(*string)(unsafe.Pointer(&jsonStr)))
    if err != nil {
        panic(err)
    }
    var signRes map[string]string
    err = json.Unmarshal([]byte(value.String()), &signRes)
    if err != nil {
        panic("解析结果失败")
    }
    return signRes
}
