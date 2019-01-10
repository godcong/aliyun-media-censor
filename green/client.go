package green

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

var DefaultClient = newClient()

func newClient() *green.Client {
	cli, err := green.NewClientWithAccessKey("cn-hangzhou", "LTAIeVGE3zRrmiNm", "F6twxkASutmcZbpPdFEqe4igtpFtu4")
	if err != nil {
		panic(err)
	}
	return cli
}

//
//func (cli Client) GetResponse(path string, clinetInfo *ClientInfo, data string) string {
//	clientInfoJson, _ := json.Marshal(clinetInfo)
//
//	req, err := http.NewRequest(method, path+"?clientInfo="+url.QueryEscape(string(clientInfoJson)), strings.NewReader(data))
//
//	if err != nil {
//		// handle error
//		return ErrorResult(err)
//	} else {
//		addRequestHeader(data, req, string(clientInfoJson), path, cli.AccessKeyId, cli.AccessKeySecret)
//
//		response, _ := http.DefaultClient.Do(req)
//
//		defer response.Body.Close()
//
//		body, err := ioutil.ReadAll(response.Body)
//		log.Println(string(body))
//		if err != nil {
//			// handle error
//			return ErrorResult(err)
//		} else {
//			return string(body)
//		}
//
//	}
//}

//type IAliYunClient interface {
//	GetResponse(path string, clinetInfo ClientInfo, bizData BizData) string
//}
