package service

import (
	"ekgo/api/boot/config"
	"ekgo/api/lib/response"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime"
	"path/filepath"
)

type UpLoadOSS struct {
	FileName string `form:"file_name" json:"file_name" binding:"required"`
}

//获取阿里云token
func (p *UpLoadOSS) GetToken() *response.Write {
	Endpoint := config.Get.Oss.Endpoint
	KeyId := config.Get.Oss.KeyId
	KeySecret := config.Get.Oss.KeySecret
	Bucket := config.Get.Oss.Bucket
	client, err := oss.New(Endpoint, KeyId, KeySecret)

	if err != nil {
		return &response.Write{
			Code:  403,
			Error: err.Error(),
		}
	}
	bucket, err := client.Bucket(Bucket)
	if err != nil {
		return &response.Write{
			Code:  403,
			Error: err.Error(),
		}
	}
	// 带可选参数的签名直传
	ext := filepath.Ext(p.FileName)
	options := []oss.Option{
		oss.ContentType(mime.TypeByExtension(ext)),
	}
	key := p.FileName
	signUrl, err := bucket.SignURL(key, oss.HTTPPut, 600, options...)
	if err != nil {
		return &response.Write{
			Code:  403,
			Error: err.Error(),
		}
	}
	return &response.Write{
		Code: 200,
		Data: map[string]string{
			"key": key,
			"put": signUrl,
		},
	}
}
