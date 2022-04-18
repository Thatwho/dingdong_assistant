/*
 * @File: user_config.go
 * @Author: 张志鹏
 * @Version: 1.0
 * @Contact: zhangzhipeng@uniml.com
 * @CreatedTime: 2022/04/15 00:03:47
 * @UpdatedTime: 2022/04/15 00:03:47
 * @Description: 配置用户信息
**/
package dingdong

// 请求头
type DingDongHeader struct {
    CityNumber   string `json:"ddmc-city-number"`
    BuildVersion string `json:"ddmc-build-version"`
    DeviceID     string `json:"ddmc-device-id"`
    StationID    string `json:"ddmc-station-id"`
    DChannel     string `json:"ddmc-channel"`
    OSVersion    string `json:"ddmc-os-version"`
    AppClientID  string `json:"ddmc-app-client-id"`
    Cookie       string `json:"cookie"`
    IP           string `json:"ddmc-ip"`
    Longitude    string `json:"ddmc-longitude"`
    Latitude     string `json:"ddmc-latitude"`
    ApiVersion   string `json:"ddmc-api-version"`
    UID          string `json:"ddmc-uid"`
    UserAgent    string `json:"user-agent"`
    Referer      string `json:"referer"`
}

func TurnHeader(h DingDongHeader) map[string]string {
    header := make(map[string]string)
    header["ddmc-city-number"] = h.CityNumber
    header["ddmc-build-version"] = h.BuildVersion
    header["ddmc-device-id"] = h.DeviceID
    header["ddmc-station-id"] = h.StationID
    header["ddmc-channel"] = h.DChannel
    header["ddmc-os-version"] = h.OSVersion
    header["ddmc-app-client-id"] = h.AppClientID
    header["cookie"] = h.Cookie
    header["ddmc-ip"] = h.IP
    header["ddmc-longitude"] = h.Longitude
    header["ddmc-latitude"] = h.Latitude
    header["ddmc-api-version"] = h.ApiVersion
    header["ddmc-uid"] = h.UID
    header["user-agent"] = h.UserAgent
    header["referer"] = h.Referer
    // header["accept-encoding"] = "gzip, deflate, br"
    return header
}

var Header DingDongHeader

// 请求体
type DingDongBody struct {
    UID          string `json:"uid"`
    Longitude    string `json:"longitude"`
    Latitude     string `json:"latitude"`
    StationID    string `json:"station_id"`
    CityNumber   string `json:"city_number"`
    ApiVersion   string `json:"api_version"`
    AppVersion   string `json:"app_version"`
    AppletSource string `json:"applet_source"`
    Channel      string `json:"channel"`
    AppClientID  string `json:"app_client_id"`
    ShareUID     string `json:"sharer_uid"`
    OpenID       string `json:"openid"`
    H5Source     string `json:"h5_source"`
    DeviceToken  string `json:"device_token"`
    SID          string `json:"s_id"`
}

func TurnBody(b DingDongBody) map[string]string {
    body := make(map[string]string)
    body["uid"] = b.UID
    body["longitude"] = b.Longitude
    body["latitude"] = b.Latitude
    body["station_id"] = b.StationID
    body["city_number"] = b.CityNumber
    body["api_version"] = b.ApiVersion
    body["app_version"] = b.AppVersion
    body["applet_source"] = b.AppletSource
    body["channel"] = b.Channel
    body["app_client_id"] = b.AppClientID
    body["sharer_uid"] = b.ShareUID
    body["openid"] = b.OpenID
    body["h5_source"] = b.H5Source
    body["device_token"] = b.DeviceToken
    body["s_id"] = b.SID
    return body
}

var Body DingDongBody

// 更新Header配置
func UpdateHeader(header DingDongHeader) {
    Header = header
}

// 更新Body配置
func UpdateBody(body DingDongBody) {
    Body = body
}

func GetRequestData() (map[string]string, map[string]string) {
    return TurnHeader(Header), TurnBody(Body)
}

// 登录过期异常
type ExpireError struct {}

func (ExpireError) Error() string {
    return "您的登录信息已过期"
}

func IsExpireError(e error) bool {
    _, ok := e.(ExpireError)
    return ok
}

// 购物车异常
type CartInvalidError struct{}

func (CartInvalidError) Error() string {
    return "购物车中没有有效的商品"
}

func IsCartInvalidError(e error) bool {
    _, ok := e.(CartInvalidError)
    return ok
}

// 账号风控异常
type RiskError struct{}

func(RiskError) Error() string {
    return "账号被风控，请联系叮咚客服"
}

func IsRiskError(e error) bool {
    _, ok := e.(RiskError)
    return ok
}