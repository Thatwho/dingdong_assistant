/*
 * @File: main.go
 * @Author: 张志鹏
 * @Version: 1.0
 * @Contact: zhangzhipeng@uniml.com
 * @CreatedTime: 2022/04/15 01:06:05
 * @UpdatedTime: 2022/04/15 01:06:05
 * @Description: 入口
**/
package main

import (
	"dingdong_hacker/handler"
	"os/exec"
	"sync"
)

func main() {
    wg := sync.WaitGroup{}
    wg.Add(1)
    go handler.Server.Run(":8000")
    cmd := exec.Command("explorer", "http://localhost:8000")
    _ = cmd.Start()
    wg.Wait()
}
