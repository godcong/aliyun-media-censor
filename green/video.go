package green

import (
	"net/url"
	"strings"
)

type ResultData struct {
	Code int `json:"code"`
	Data []struct {
		Code    int    `json:"code"`
		DataID  string `json:"dataId"`
		Msg     string `json:"msg"`
		Results []struct {
			Frames []struct {
				Label     string  `json:"label"`
				Offset    int     `json:"offset"`
				Rate      float64 `json:"rate"`
				SfaceData []struct {
					Faces []struct {
						ID   string  `json:"id"`
						Name string  `json:"name"`
						Rate float64 `json:"rate"`
					} `json:"faces"`
					H int `json:"h"`
					W int `json:"w"`
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"sfaceData"`
				URL string `json:"url"`
			} `json:"frames"`
			Label      string  `json:"label"`
			Rate       float64 `json:"rate"`
			Scene      string  `json:"scene"`
			Suggestion string  `json:"suggestion"`
		} `json:"results"`
		TaskID string `json:"taskId"`
	} `json:"data"`
	Msg string `json:"msg"`
	//DataID    string `json:"dataId"`
	//TaskID    string `json:"taskId"`
	RequestID string `json:"requestId"`
}

type ClientInfo struct {
	SDKVersion  string `json:"sdkVersion"`  //否	SDK版本，通过SDK调用时，需提供该字段。
	CFGVersion  string `json:"cfgVersion"`  //否	配置信息版本，通过SDK调用时，需提供该字段。
	UserType    string `json:"userType"`    //否	用户账号类型，取值为：	taobao	others
	UserId      string `json:"userId"`      //否	用户ID，唯一标识一个用户。
	UserNick    string `json:"userNick"`    //否	用户昵称。
	Avatar      string `json:"avatar"`      //否	用户头像。
	Imei        string `json:"imei"`        //否	硬件设备码。
	Imsi        string `json:"imsi"`        //否	运营商设备码。
	Umid        string `json:"umid"`        //否	设备指纹。
	Ip          string `json:"ip"`          //否	该IP应该为公网IP。如果请求中不填写，服务端会尝试从链接或者从HTTP头中获取。如果请求是从设备端发起的，该字段通常不填写；如果是从后台发起的，该IP为用户的login IP或者设备的公网IP。
	Os          string `json:"os"`          //否	设备的操作系统，如：Android 6.0
	Channel     string `json:"channel"`     //否	渠道号。
	HostAppName string `json:"hostAppName"` //否	宿主应用名称。
	HostPackage string `json:"hostPackage"` //否	宿主应用包名。
	HostVersion string `json:"hostVersion"` //否	宿主应用版本。
}

type VideoRequest struct {
	BizType     string   `json:"bizType,omitempty"`
	Scenes      []string `json:"scenes"`
	AudioScenes []string `json:"audioScenes,omitempty"`
	Callback    string   `json:"callback,omitempty"`
	Seed        string   `json:"seed,omitempty"`
	Tasks       []struct {
		ClientInfo []ClientInfo `json:"clientInfo,omitempty"`
		DataID     string       `json:"dataId,omitempty"`
		URL        string       `json:"url,omitempty"`
		Frames     []struct {
			URL    string `json:"url,omitempty"`
			Offset int    `json:"offset,omitempty"`
		} `json:"frames,omitempty"`
		FramePrefix string `json:"framePrefix,omitempty"`
		Interval    int    `json:"interval,omitempty"`
		MaxFrames   int    `json:"maxFrames,omitempty"`
	} `json:"tasks"`
}

const AliSite = "http://green.cn-shanghai.aliyuncs.com"

func URL(values url.Values) string {
	return AliSite + "?" + values.Encode()
}

func Link(url string) string {
	if strings.Index(url, "/") == 0 {
		return AliSite + url
	}

	return AliSite + "/" + url
}

func GreenVideoSyncscan(request *VideoRequest) (*ResultData, error) {
	url := Link("green/video/syncscan")

}

func GreenVideoAsyncscan(request *VideoRequest) (*ResultData, error) {
	url := Link("green/video/asyncscan")
}
