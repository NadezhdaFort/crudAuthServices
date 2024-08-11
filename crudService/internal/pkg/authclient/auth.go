package authclient

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

var c *fasthttp.HostClient

func Init(host string) {
	c = &fasthttp.HostClient{
		Addr: host,
	}
}

func ValidateToken(token string) (*UserIdRole, error) {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("http://" + c.Addr + "/v3/get_user_info")
	req.Header.Set(fasthttp.HeaderAuthorization, token)
	req.Header.SetHost(c.Addr)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := c.Do(req, resp)
	if err != nil {
		return nil, err
	}
	log.Println(resp)
	log.Println(resp.StatusCode())

	if resp.StatusCode() != http.StatusOK {
		return nil, err
	}

	body := resp.Body()
	log.Printf("JSON response body: %s\n", string(body))

	var authResp HTTPResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return nil, err
	}
	var userInfo UserIdRole
	userInfo = authResp.UserInfo

	log.Println("GetUserInfo: ", userInfo, " from crudServer validateToken")

	return &userInfo, nil

}
