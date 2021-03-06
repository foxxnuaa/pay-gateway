package sign

import (
	"encoding/json"
	"fmt"
	"github.com/blademainer/commons/pkg/field"
	"github.com/fatih/structs"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"gopkg.in/go-playground/assert.v1"
	"reflect"
	"testing"
	"time"
)

type AType struct {
	Version string `protobuf:"bytes,100,opt,name=version,proto3" json:"version,omitempty"`
	// 业务订单号
	OutTradeNo string `protobuf:"bytes,1,opt,name=out_trade_no,json=outTradeNo,proto3" json:"out_trade_no,omitempty"`
	// 支付金额（分）
	PayAmount uint32 `protobuf:"varint,3,opt,name=pay_amount,json=payAmount,proto3" json:"pay_amount,omitempty"`
	// 币种
	Currency string `protobuf:"bytes,20,opt,name=currency,proto3" json:"currency,omitempty"`
	// 接收通知的地址，不能带参数（即：不能包含问号）
	NotifyURL string `protobuf:"bytes,4,opt,name=notify_url,json=notifyURL,proto3" json:"notify_url,omitempty"`
	// 支付后跳转的前端地址
	ReturnURL string `protobuf:"bytes,5,opt,name=return_url,json=returnURL,proto3" json:"return_url,omitempty"`
	// 系统给商户分配的app_id
	AppID string `protobuf:"bytes,6,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
	// 加密方法，RSA和MD5，默认RSA
	SignType string `protobuf:"bytes,7,opt,name=sign_type,json=signType,proto3" json:"sign_type,omitempty"`
	// 签名
	Sign string `protobuf:"bytes,14,opt,name=sign,proto3" json:"sign,omitempty"`
	// 业务方下单时间，时间格式: 年年年年-月月-日日 时时:分分:秒秒，例如: 2006-01-02 15:04:05
	OrderTime string `protobuf:"bytes,8,opt,name=order_time,json=orderTime,proto3" json:"order_time,omitempty"`
	// 发起支付的用户ip
	UserIP string `protobuf:"bytes,9,opt,name=user_ip,json=userIP,proto3" json:"user_ip,omitempty"`
	// 用户在业务系统的id
	UserID string `protobuf:"bytes,18,opt,name=user_id,json=userID,proto3" json:"user_id,omitempty"`
	// 支付者账号，可选
	PayerAccount string `protobuf:"bytes,10,opt,name=payer_account,json=payerAccount,proto3" json:"payer_account,omitempty"`
	// 业务系统的产品id
	ProductID string `protobuf:"bytes,11,opt,name=product_id,json=productID,proto3" json:"product_id,omitempty"`
	// 商品名称
	ProductName string `protobuf:"bytes,12,opt,name=product_name,json=productName,proto3" json:"product_name,omitempty"`
	// 商品描述
	ProductDescribe string `protobuf:"bytes,13,opt,name=product_describe,json=productDescribe,proto3" json:"product_describe,omitempty"`
	// 参数编码，只允许utf-8编码；签名时一定要使用该编码获取字节然后再进行签名
	Charset string `protobuf:"bytes,15,opt,name=charset,proto3" json:"charset,omitempty"`
	// 回调业务系统时需要带上的字符串
	CallbackJSON string `protobuf:"bytes,16,opt,name=callback_json,json=callbackJSON,proto3" json:"callback_json,omitempty"`
	// 扩展json
	ExtJSON string `protobuf:"bytes,17,opt,name=ext_json,json=extJSON,proto3" json:"ext_json,omitempty"`
	// 渠道id（非必须），如果未指定method，系统会根据method来找到可用的channel_id
	ChannelID string `protobuf:"bytes,19,opt,name=channel_id,json=channelID,proto3" json:"channel_id,omitempty"`
	// 例如：二维码支付，银联支付等。
	Method string               `protobuf:"bytes,98,opt,name=method,proto3" json:"method,omitempty"`
	Time   *timestamp.Timestamp `protobuf:"bytes,8,opt,name=order_time,json=time,proto3" json:"time,omitempty"`
}

func TestType(t *testing.T) {
	of := reflect.TypeOf(AType{})
	fmt.Println(of)
	names := structs.Names(AType{})
	fmt.Println("Field names: ", names)
	aType := AType{Version: "hello", PayAmount: 23}
	params := structs.Map(aType)
	for k, v := range params {
		fmt.Printf("k: %v v: %v", k, v)
	}

	a := structs.New(aType)
	for _, f := range a.Fields() {
		name := f.Name()
		tag := f.Tag("json")
		fmt.Printf("name: %s tag: %s values: %v \n", name, tag, f.Value())
	}
}

func TestConvert(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339Nano, "2020-12-14T13:03:35.025484056+08:00")
	pt, _ := ptypes.TimestampProto(tm)

	aType := AType{
		Version: "hello", PayAmount: 23, Sign: "sssss",
		Time: pt,
	}
	bytes, _ := json.Marshal(aType)
	fmt.Println("json: ", string(bytes))
	compacter := NewParamsCompacter(AType{}, "json", []string{"sign"}, true, "&", "=")
	s := compacter.ParamsToString(aType)
	fmt.Println(s)
	// testutil.AssertEqual(t, "pay_amount=23&version=hello", s)

	assert.Equal(t, "pay_amount=23&time=2020-12-14T05:03:35.025484056Z&version=hello", s)
}

func TestParamsCompacter_ParamsToString(t *testing.T) {
	now := time.Now()

	aType := AType{
		Version: "hello", PayAmount: 23, Sign: "sssss", Time: &timestamp.Timestamp{
			Seconds: int64(now.Second()),
		},
	}
	s2 := aType.Time.String()
	fmt.Println("time: ", s2)
	bytes, _ := json.Marshal(aType)
	fmt.Println("json: ", string(bytes))
	compacter := NewParamsCompacter(AType{}, "json", []string{"sign"}, true, "&", "=")
	s := compacter.ParamsToString(aType)
	fmt.Println(s)
	// testutil.AssertEqual(t, "pay_amount=23&version=hello", s)

	fp := &field.Parser{
		Tag:                 "json",
		Escape:              false,
		GroupDelimiter:      '&',
		PairDelimiter:       '=',
		Sort:                true,
		IgnoreNilValueField: true,
	}
	marshal, err := fp.Marshal(aType)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(string(marshal))
}

func Benchmark(b *testing.B) {
	now := time.Now()

	aType := AType{
		Version: "hello", PayAmount: 23, Sign: "sssss", Time: &timestamp.Timestamp{
			Seconds: int64(now.Second()),
		},
	}

	compacter := NewParamsCompacter(AType{}, "json", []string{"sign"}, true, "&", "=")
	for i := 0; i < b.N; i++ {
		compacter.ParamsToString(aType)
	}
}

func Test_valueFunc(t *testing.T) {
	f := valueFunc(reflect.TypeOf(&timestamp.Timestamp{}))
	s, err := f(
		ptypes.TimestampNow(),
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(s)

	f2 := valueFunc(reflect.TypeOf(uint32(13)))
	s2, err := f2(uint32(12))
	if err != nil{
		t.Fatal(err.Error())
	}
	fmt.Println(s2)
	f2 = valueFunc(reflect.TypeOf(int32(13)))
	s2, err = f2(int32(13))
	if err != nil{
		t.Fatal(err.Error())
	}
	fmt.Println(s2)
	f2 = valueFunc(reflect.TypeOf(int64(13)))
	s2, err = f2(int64(13))
	if err != nil{
		t.Fatal(err.Error())
	}
	fmt.Println(s2)
	f2 = valueFunc(reflect.TypeOf(uint64(13)))
	s2, err = f2(uint64(13))
	if err != nil{
		t.Fatal(err.Error())
	}
	fmt.Println(s2)
}


func Test_intKind(t *testing.T) {
	type A int32
	value := A(12)
	of := reflect.TypeOf(value)
	fmt.Printf("type: %T kind: %v\n", value, of.Kind())
	valueOf := reflect.ValueOf(value)
	i := valueOf.Int()
	fmt.Println(i)
}
