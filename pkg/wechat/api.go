package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

const host = "https://api.weixin.qq.com"

type BasicAPI interface { //基础接口，需提供appid和secret
	Code2Session(ctx context.Context, code string) (*Code2SessionResp, error)
	GetAccessToken(ctx context.Context) (*GetAccessTokenResp, error)
}

type ServerAPI interface { //服务端接口，需提供access_token
	GetUserPhoneNumber(ctx context.Context, code string) (*UserPhoneNumberResp, error)
}

type FullAPI interface { //全部接口
	BasicAPI
	ServerAPI
}

type basic struct {
	appid  string
	secret string
	client *http.Client
}

type TokenFunc func(ctx context.Context) (string, error)

type server struct {
	client *http.Client
	token  TokenFunc
}

type full struct {
	*basic
	*server
}

func NewBasicAPI(appid, secret string, client *http.Client) BasicAPI {
	if client == nil {
		client = http.DefaultClient
	}
	return &basic{
		appid:  appid,
		secret: secret,
		client: client,
	}
}

func NewServerAPI(client *http.Client, token TokenFunc) ServerAPI {
	if client == nil {
		client = http.DefaultClient
	}
	return &server{
		client: client,
		token:  token,
	}
}

func NewFullAPI(appid, secret string, client *http.Client, token TokenFunc) FullAPI {
	if client == nil {
		client = http.DefaultClient
	}
	return &full{
		basic: &basic{
			appid:  appid,
			secret: secret,
			client: client,
		},
		server: &server{
			client: client,
			token:  token,
		},
	}
}

func (api *basic) get(ctx context.Context, path string, data url.Values, result any) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+path+"?"+data.Encode(), http.NoBody)
	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(result)
	return err
}

func (api *server) post(ctx context.Context, path string, data any, result any) error {
	tk, err := api.token(ctx)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(data)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, host+path+"?access_token="+tk, bytes.NewReader(b))
	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(result)
	return err
}

type respErr struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}
