package green

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

var DefaultClient = newClient()

func newClient() *green.Client {
	//TODO add client
	return &green.Client{}
}