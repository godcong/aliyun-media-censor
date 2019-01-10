package green

import (
	"bytes"
	"github.com/godcong/aliyun-media-censor/util"
	"github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const ContentTypeJSON = "application/json"
const GMTTimeFormmat = "Mon, 02 Jan 2006 15:04:05 -0700 GMT"

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

type VideoRequest struct {
	BizType     string   `json:"bizType,omitempty"`
	Scenes      []string `json:"scenes"`
	AudioScenes []string `json:"audioScenes,omitempty"`
	Callback    string   `json:"callback,omitempty"`
	Seed        string   `json:"seed,omitempty"`
	Tasks       []Task   `json:"tasks"`
}

const AliSite = "https://green.cn-shanghai.aliyuncs.com"

func URL(values url.Values) string {
	return AliSite + "?" + values.Encode()
}

func Link(url string) string {
	if strings.Index(url, "/") == 0 {
		return AliSite + url
	}

	return AliSite + "/" + url
}

func Header() *http.Header {
	h := http.Header{}
	h.Set("Accept", "application/json")
	h.Set("Content-Type", "application/json")
	h.Set("Date", time.Now().Format(GMTTimeFormmat))
	h.Set("x-acs-version", "2018-05-09")
	h.Set("x-acs-signature-nonce", util.GenerateRandomString(32))
	h.Set("x-acs-signature-version", "1.0")
	h.Set("x-acs-signature-method", "HMAC-SHA1")
	h.Set("Authorization", "...")
	return &h
}

func GreenVideoSyncscan(request *VideoRequest) (*ResultData, error) {
	url := Link("green/video/syncscan")
	log.Println("url:", url)
	resp, err := http.Post(url, ContentTypeJSON, request.Reader())
	if err != nil {
		return nil, err
	}
	var data ResultData
	err = util.UnmarshalJSON(resp.Body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil

}

func GreenVideoAsyncscan(request *VideoRequest) (*ResultData, error) {
	url := Link("green/video/asyncscan")
	log.Println("url:", url)
	resp, err := http.Post(url, ContentTypeJSON, request.Reader())
	if err != nil {
		return nil, err
	}
	var data ResultData
	err = util.UnmarshalJSON(resp.Body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func GreenVideoResults(request ...string) (*ResultData, error) {
	url := Link("green/video/results")
	log.Println("url:", url)
	marshaled, err := jsoniter.Marshal(request)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(marshaled)
	resp, err := http.Post(url, ContentTypeJSON, reader)
	if err != nil {
		return nil, err
	}
	var data ResultData
	err = util.UnmarshalJSON(resp.Body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *VideoRequest) Reader() io.Reader {
	b, err := jsoniter.Marshal(r)
	if err != nil {
		log.Println(err)
		return nil
	}
	return bytes.NewBuffer(b)

}
