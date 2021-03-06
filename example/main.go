package main

import (
	"github.com/Titanarthas/docauto"
)

type CouponTakeListParams struct {
	UserID     int64  `json:"user_id" valid:"optional" docComment:"内置类型测试"`                                       // 应用内领券必传
	UserID2     *int64  `json:"user_id2" valid:"optional" docComment:"空指针测试"`
	ActivityID []int64  `json:"activity_id" valid:"optional" docComment:"数组测试"`                                   // 优惠券ID
	Type       *int64    `json:"type" valid:"required~参数type必填" docComment:"指针测试"`            // 券类型，0：普通优惠券；1:纸质优惠券；2，会员专属
	CTLP CouponTakeListParams2 `json:"ctlp" docComment:"2级测试"`
	CTLP2 []CouponTakeListParams2 `json:"ctlp2" docComment:"2级测试2"`
	CTLP3 *CouponTakeListParams2 `json:"ctlp3" docComment:"*测试"`
}

type CouponTakeListParams2 struct {
	UserID     int64  `json:"user_id" valid:"optional" docComment:"应用内领券必传"`                                       // 应用内领券必传
	ActivityID []int64  `json:"activity_id" valid:"optional" docComment:"优惠券ID"`                                   // 优惠券ID
	Type       int    `json:"type" valid:"required~参数type必填" docComment:"券类型，0：普通优惠券；1:纸质优惠券；2，会员专属"`            // 券类型，0：普通优惠券；1:纸质优惠券；2，会员专属
}

func main() {
	docauto.Init(&docauto.Config{
		On:       true,                 //是否开启自动生成API文档功能
		DocTitle: "Iris",
		DocPath:  "apidoc.html",        //生成API文档名称存放路径
		BaseUrls: map[string]string{"Production": "", "Staging": ""},
	})

	req := CouponTakeListParams{}
	req.UserID = 1
	req.Type = &req.UserID
	req.CTLP2 = []CouponTakeListParams2{CouponTakeListParams2{}}
	req.CTLP3 = &CouponTakeListParams2{}
	docauto.GenerateDocStruct("POST", "/hello/", "test", req, req)
}
