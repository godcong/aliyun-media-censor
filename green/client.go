package green

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

// DefaultClient ...
var DefaultClient = newClient()

func newClient() *green.Client {
	//TODO add client
	cli, err := green.NewClientWithAccessKey("cn-hangzhou", "LTAIeVGE3zRrmiNm", "F6twxkASutmcZbpPdFEqe4igtpFtu4")
	if err != nil {
		panic(err)
	} else {
		return cli
	}

	return &green.Client{}
}
