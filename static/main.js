
// 加载完页面后，获取本地缓存的用户信息
window.addEventListener("load", function() {
  // 1.获取本地存储的数据
  userConfig = localStorage.getItem("UserConfig")
  if (userConfig !== null) {
    let cachedData = JSON.parse(userConfig);
    if (cachedData !== {}) {
      let header = cachedData["header"]
      for(let k in header) {
        console.log(k)
        if (k === "user-agent") {
          document.getElementById(k).textContent = header[k];
        } else {
          document.getElementById(k).value = header[k];
        }
      };
      let body = cachedData["body"]
      for(let k in body) {
        document.getElementById(k).value = body[k];
      }
      UserConfSetted = true;
      UpdateUserMsg("已加载缓存的抓包数据")
    }
  }
})

// 存储用户数据在本地
function StoreUserConfig(userConfigData) {
  localStorage.removeItem("UserConfig")
  localStorage.setItem("UserConfig", userConfigData)
}

// 解析表单数据为JSON
class TransFormToJson {
  constructor(){
    let userCnfData = new FormData(userCnfForm);
    this.header = {
        "ddmc-city-number": userCnfData.get("ddmc-city-number"),
        "ddmc-build-version": userCnfData.get("ddmc-build-version"),
        "ddmc-device-id": userCnfData.get("ddmc-device-id"),
        "ddmc-station-id": userCnfData.get("ddmc-station-id"),
        "ddmc-channel": userCnfData.get("ddmc-channel"),
        "ddmc-os-version": userCnfData.get("ddmc-os-version"),
        "ddmc-app-client-id": userCnfData.get("ddmc-app-client-id"),
        "cookie": userCnfData.get("cookie"),
        "ddmc-ip": userCnfData.get("ddmc-ip"),
        "ddmc-longitude": userCnfData.get("ddmc-longitude"),
        "ddmc-latitude": userCnfData.get("ddmc-latitude"),
        "ddmc-api-version": userCnfData.get("ddmc-api-version"),
        "ddmc-uid": userCnfData.get("ddmc-uid"),
        "user-agent": userCnfData.get("user-agent"),
        "referer": userCnfData.get("referer")
    };
    this.body = {
        "uid": userCnfData.get("uid"),
        "longitude": userCnfData.get("longitude"),
        "latitude": userCnfData.get("latitude"),
        "station_id": userCnfData.get("station_id"),
        "city_number": userCnfData.get("city_number"),
        "api_version": userCnfData.get("api_version"),
        "app_version": userCnfData.get("app_version"),
        "applet_source": userCnfData.get("applet_source"),
        "channel": userCnfData.get("channel"),
        "app_client_id": userCnfData.get("app_client_id"),
        "sharer_uid": userCnfData.get("sharer_uid"),
        "openid": userCnfData.get("openid"),
        "h5_source": userCnfData.get("h5_source"),
        "device_token": userCnfData.get("device_token"),
        "s_id": userCnfData.get("s_id")
    };
  }
}

// 更新用户提示信息
function UpdateUserMsg(msg) {
  UserMsg.textContent = msg;
}

// 提交用户配置
function PutUserConfig() {
  // 1.获取表单数据
  let userCnf = new TransFormToJson();

  // TODO: 2.验证表数据

  // 3.提交表单数据
  let userCnfStr = JSON.stringify(userCnf);
  // let userCnfStr = userCnf;
  fetch("/api/v1/user_conf", {
    method: 'PUT',
    body: userCnfStr
  }).then(function(response) {
    if (response.status == "200") {
      UpdateUserMsg("配置更新成功");
      StoreUserConfig(userCnfStr);
      UserConfSetted = true;
    } else {
      response.json().then(function(data) {
        UpdateUserMsg("发生异常,请重试\n" + data["msg"]);
      })
    }
  }).catch(function(err) {
    UpdateUserMsg("发生错误: ", err.message);
  })
}

// 标记用户配置是否设置过
let UserConfSetted = false;

// 信息展示框
let UserMsg = document.getElementById("message").querySelector("p");

// 定位用户配置表单
let userCnfForm = document.getElementById("user-config")
// 禁用表单默认提交行为
userCnfForm.addEventListener("submit", function(event) {
  event.preventDefault();
})

// 更新用户地址ID
function UpdateUserAddress() {
  // 1.判断用户配置是否设置过
  if (!UserConfSetted) {
    alert("请先更新抓包数据！");
    return
  }
  // 请求地址
  fetch("/api/v1/address", {
    method: 'PUT'
  }).then(function(response) {
    if (response.status == 200) {
      response.json().then(function(json) {
        UpdateUserMsg("您的默认地址是: " + json["detail"]);
      })
    } else {
      response.json().then(function(json) {
        UpdateUserMsg("服务异常: " + json["msg"]);
      })
    }
  }).catch(function(err) {
    UpdateUserMsg("发生错误: ", + err.message);
  })
}

// 抢购
function Purchase() {
  // 1.判断用户配置是否设置过
  if (!UserConfSetted) {
    alert("请先更新抓包数据！");
    return
  }
  // 请求地址
  let userCnfData = new FormData(userCnfForm);
  conNum = userCnfData.get("concurrent_num")
  fetch("api/v1/purchase?conNum=" + conNum, {
    method: "POST"
  }).then(function(response) {
    if (response.status == 204) {
      // 轮询抢购结果
      PollingStatus();
    } else {
      response.json().then(
        function(json){
          UpdateUserMsg("抢购失败! " + json["msg"]);
        }
      )
    }
  }).catch(function(err) {
    UpdateUserMsg("发生错误: ", + err.message);
  })
}

// 轮询抢购状态
let purchaseFinished = false

let purchaseStatus = ["\\", "|", "/", "-"]

let checkNum = 0;

// 查询异常的次数
let UnhealthRespCnt = 0;

function UpdatePurchasStatus() {
  checkNum += 1;
  fetch("/api/v1/order", {
    method: "GET"
  }).then(function(response) {
    if (response.status == 200) {
      response.json().then(function(json) {
        if (json["finished"]) {
          purchaseFinished = true;
          if (json["orderSuccess"]) {
            UpdateUserMsg("下单成功, 成功时间: " + json["orderTime"] + ", 快去付款!!!");
            AwakeUser();
          } else {
            UpdateUserMsg("抢菜已停止，很抱歉没抢到。" + json["reason"])
          }
        } else {
          UpdateUserMsg("抢菜中 " + purchaseStatus[(checkNum % purchaseStatus.length)]);
        }
      })
    } else {
      UnhealthRespCnt += 1;
      // 如果异常响应超过10次，则退出
      if (UnhealthRespCnt >= 10) {
        purchaseFinished = true;
        UpdateUserMsg("服务异常! 已停止！")
      } else {
        UpdateUserMsg("服务异常！")
      }
    }
  }).catch(function(err) {
    UpdateUserMsg("发生错误: ", + err.message);
  })
}

function PollingStatus() {
  this.timeId = setInterval(() => {
    if (purchaseFinished) {
      // 重置购买结束状态
      purchaseFinished = false;
      clearInterval(this.timeId);
    };
    UpdatePurchasStatus();
  }, 300)
}

// 用户配置更新按钮
let updateUserInfnBtn = document.querySelector("#update-user-info");
updateUserInfnBtn.addEventListener("click", function(){
  PutUserConfig();
})

// 用户地址ID获取按钮
let updateAdrBtn = document.getElementById("update-address");
updateAdrBtn.addEventListener("click", function() {
  UpdateUserAddress();
})

// 抢购按钮
let PurchaseBtn = document.getElementById("purchase");
PurchaseBtn.addEventListener("click", function() {
  Purchase();
})

// 播放音乐，提醒用户
function AwakeUser() {
  let player = document.createElement("audio");
  player.src = "/the_internationable.mp3";
  player.play();
}
