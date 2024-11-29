package local

import (
	"apylee_chat_server/internal/service/redis"
	"apylee_chat_server/pkg/log"
	"apylee_chat_server/pkg/util/random"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	console "github.com/alibabacloud-go/tea-console/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strconv"
	"time"
)

const (
	accessKeyID     = "your accessKeyID in alibaba cloud"
	accessKeySecret = "your accessKeySecret in alibaba cloud"
)

var smsClient *dysmsapi20170525.Client

// createClient 使用AK&SK初始化账号Client
func createClient() (result *dysmsapi20170525.Client, err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	if smsClient == nil {
		config := &openapi.Config{
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
			AccessKeyId: tea.String(accessKeyID),
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
			AccessKeySecret: tea.String(accessKeySecret),
		}
		// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
		config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
		smsClient, err = dysmsapi20170525.NewClient(config)
	}
	return smsClient, err
}

func VerificationCode(telephone string) error {
	client, err := createClient()
	if err != nil {
		return err
	}
	code, _ := redis.GetKey("auto_code")
	if code != "" {
		// 直接返回，验证码还没过期，用户应该去输验证码
		log.GetZapLogger().Info("please input auth code")
		return nil
	}
	// 验证码过期，重新生成
	code = strconv.Itoa(random.GetRandomInt(6))
	fmt.Println(code)
	err = redis.SetKeyEx("auth_code", code, time.Minute*10) // 10分钟有效
	if err != nil {
		return err
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"), // 短信模板
		PhoneNumbers:  tea.String(telephone),
		TemplateParam: tea.String("{\"code\":\"" + code + "\"}"),
	}

	runtime := &util.RuntimeOptions{}
	// 目前使用的是测试专用签名，签名必须是“阿里云短信测试”，模板code为“SMS_154950909”
	rsp, err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		return err
	}
	console.Log(util.ToJSONString(rsp))
	return nil
}
