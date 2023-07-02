package mytencent

import (
	"fmt"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
)

/*
go get -v -u github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn
go get github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cdn/v20180606
go get github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common

*/

func PurgeCdn(SecretId string, SecretKey string, hostname string) bool {
	credential := common.NewCredential(
		SecretId,
		SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cdn.NewClient(credential, "", cpf)
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cdn.NewPurgePathCacheRequest()

	request.Paths = common.StringPtrs([]string{"https://" + hostname})
	request.FlushType = common.StringPtr("flush")
	response, err := client.PurgePathCache(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		//fmt.Printf("An API error has returned: %s", err)
		return false
	}
	if err != nil {
		//panic(err)
		return false
	}
	// 输出json格式的字符串回包
	fmt.Println("[Tencent] " + response.ToJsonString())
	return true
}
