package green

import (
	"github.com/json-iterator/go"
	"log"
)

type ClientInfo struct {
	SDKVersion  string `json:"sdkVersion,omitempty"`  //否	SDK版本，通过SDK调用时，需提供该字段。
	CFGVersion  string `json:"cfgVersion,omitempty"`  //否	配置信息版本，通过SDK调用时，需提供该字段。
	UserType    string `json:"userType,omitempty"`    //否	用户账号类型，取值为：	taobao	others
	UserID      string `json:"userId,omitempty"`      //否	用户ID，唯一标识一个用户。
	UserNick    string `json:"userNick,omitempty"`    //否	用户昵称。
	Avatar      string `json:"avatar,omitempty"`      //否	用户头像。
	Imei        string `json:"imei,omitempty"`        //否	硬件设备码。
	Imsi        string `json:"imsi,omitempty"`        //否	运营商设备码。
	Umid        string `json:"umid,omitempty"`        //否	设备指纹。
	Ip          string `json:"ip,omitempty"`          //否	该IP应该为公网IP。如果请求中不填写，服务端会尝试从链接或者从HTTP头中获取。如果请求是从设备端发起的，该字段通常不填写；如果是从后台发起的，该IP为用户的login IP或者设备的公网IP。
	Os          string `json:"os,omitempty"`          //否	设备的操作系统，如：Android 6.0
	Channel     string `json:"channel,omitempty"`     //否	渠道号。
	HostAppName string `json:"hostAppName,omitempty"` //否	宿主应用名称。
	HostPackage string `json:"hostPackage,omitempty"` //否	宿主应用包名。
	HostVersion string `json:"hostVersion,omitempty"` //否	宿主应用版本。
}

type Frame struct {
	URL    string `json:"url,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

type Task struct {
	ClientInfo  []ClientInfo `json:"clientInfo,omitempty"`
	DataID      string       `json:"dataId,omitempty"`
	URL         string       `json:"url,omitempty"`
	Frames      []Frame      `json:"frames,omitempty"`
	FramePrefix string       `json:"framePrefix,omitempty"`
	Interval    int          `json:"interval,omitempty"`
	MaxFrames   int          `json:"maxFrames,omitempty"`
}

type BizData struct {
	BizType     string   `json:"bizType,omitempty"`
	Scenes      []string `json:"scenes"`
	AudioScenes []string `json:"audioScenes,omitempty"`
	Callback    string   `json:"callback,omitempty"`
	Seed        string   `json:"seed,omitempty"`
	Tasks       []Task   `json:"tasks"`
}

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
		URL    string `json:"url"`
	} `json:"data"`
	Msg       string `json:"msg"`
	RequestID string `json:"requestId"`
}

func (data *BizData) JSON() []byte {
	bytes, err := jsoniter.Marshal(data)
	if err != nil {
		log.Println(err)
		return nil
	}
	return bytes
}
