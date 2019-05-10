/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package yz

import (
	"errors"
	"fmt"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	youzan "github.com/xu42/youzan-sdk-go"
	"github.com/xu42/youzan-sdk-go/api"
	"github.com/xu42/youzan-sdk-go/auth"
)

func New(clientID, clientSecret, kdtID string, versions ...string) *YouZan {
	c := &YouZan{
		clientID:     clientID,
		clientSecret: clientSecret,
		kdtID:        kdtID,
		version:      `3.0.0`,
	}
	if len(versions) > 0 {
		c.version = versions[0]
	}
	return c
}

type YouZan struct {
	clientID     string
	clientSecret string
	kdtID        string //授权店铺id
	version      string
	resp         auth.GenSelfTokenResponse
	Debug        bool
}

type QRInfo struct {
	URL  string `json:"qr_url"`
	Code string `json:"qr_code"`
	ID   int64  `json:"qr_id"`
}

func (y *YouZan) Request() (resp auth.GenSelfTokenResponse, err error) {
	return youzan.GenSelfToken(y.clientID, y.clientSecret, y.kdtID)
}

func (y *YouZan) Call(api string, params map[string]string, versions ...string) (result api.CallResponse, err error) {
	resp, err := y.Request()
	if err != nil {
		return result, err
	}
	var version string
	if len(versions) > 0 {
		version = versions[0]
	}
	if len(version) == 0 {
		version = y.version
	}
	result, err = youzan.Call(resp.AccessToken, api, version, params)
	if y.Debug {
		fmt.Println(result.Success, result.Result, result.Error, err)
	}
	return
}

func (y *YouZan) CreateQROrder(data echo.Store) (echo.Store, error) {
	price := data.Float64("price")
	if price <= 0 {
		return data, errors.New(echo.T(`价格必须大于0`))
	}
	var ids string
	switch v := data.Get(`ids`).(type) {
	case string:
		ids = `["` + v + `"]`
	case nil:
	case []string, []int64, []uint64, []uint, []uint32, []int, []int32, []interface{}:
		b, e := com.JSONEncode(v)
		if e != nil {
			return data, e
		}
		ids = com.Bytes2str(b)
	default:
		r := []interface{}{v}
		b, e := com.JSONEncode(r)
		if e != nil {
			return data, e
		}
		ids = com.Bytes2str(b)
	}
	params := map[string]string{
		"label_ids": ids,                                    //标签 json格式字符串,例如:[1,2],表示有两个标签
		"qr_name":   data.String("name"),                    //收款理由
		"qr_price":  tplfunc.NumberFormat(price*100, 2, ``), //价格（单位 分）。qr_type 为 QR_TYPE_FIXED 时，价格可为空。
		"qr_type":   "QR_TYPE_DYNAMIC",                      //二维码类型. QR_TYPE_FIXED_BY_PERSON ：无金额二维码，扫码后用户需自己输入金额； QR_TYPE_NOLIMIT ： 确定金额二维码，可以重复支付; QR_TYPE_DYNAMIC：确定金额二维码，只能被支付一次
		"qr_source": data.String("source"),                  //二维码创建类型标记
	}
	result, err := y.Call("youzan.pay.qrcode.create", params)
	if !result.Success {
		return data, errors.New(result.Error.Msg)
	}
	//qr_url: 需要登录才能支付
	data.DeepMerge(result.Result)
	return data, err
}

func (y *YouZan) HandleNotify(ctx echo.Context) (echo.Store, error) {
	result := echo.Store{}
	err := ctx.MustBind(&result)
	if err != nil {
		return result, err
	}
	typ := result.String(`type`)
	if typ != `trade_TradePaid` {
		return result, nil
	}
	msg := result.String(`msg`)
	if len(msg) == 0 {
		return result, ctx.E(`Invalid message with empty msg field`)
	}
	signed := com.Md5(y.clientID + msg + y.clientSecret)
	if signed != result.String(`sign`) {
		return result, ctx.E(`签名错误`)
	}
	order := echo.Store{}
	if m, e := com.URLDecode(msg); e == nil {
		msg = m
	}
	err = com.JSONDecode(com.Str2bytes(msg), &order)
	if err != nil {
		return result, err
	}
	tid := order.String(`tid`)
	params := map[string]string{
		`tid`: tid,
	}
	response, err := y.Call("youzan.trade.get", params, "4.0.0")
	if !response.Success {
		return result, errors.New(response.Error.Msg)
	}
	result.DeepMerge(response.Result)
	var uid, qrID string
	orderInfo := result.Store(`full_order_info`)
	buyerInfo := orderInfo.Store(`buyer_info`)
	buyerID := buyerInfo.Int64(`buyer_id`)
	if buyerID <= 0 {
		buyerID = buyerInfo.Int64(`fans_id`)
	}
	if buyerID <= 0 {
		uid = buyerInfo.String(`outer_user_id`)
	} else {
		uid = fmt.Sprint(buyerID)
	}
	if qrInfo, ok := result.Get(`qr_info`).(echo.Store); ok {
		qrID = qrInfo.String(`qr_id`)
	}
	result.Set(`qr_id`, qrID)
	result.Set(`uid`, uid)
	return result, err
}
