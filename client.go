package vkapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	tokenURL = "https://oauth.vk.com/token"
	apiURL   = "https://api.vk.com/method/%s"
)

type VKClient struct {
	Self   Token
	Client *http.Client
	rl     *rateLimiter
	cb     *callbackHandler
}

type TokenOptions struct {
	ServiceToken    bool
	ValidateOnStart bool
}

func newVKClientBlank() *VKClient {
	return &VKClient{
		Client: &http.Client{},
		rl:     &rateLimiter{},
		cb: &callbackHandler{
			events: make(map[string]func(*LongPollMessage)),
		},
	}
}

func NewVKClient(user string, password string, config AuthConfig) (*VKClient, error) {
	vkClient := newVKClientBlank()

	token, err := vkClient.auth(user, password, config)
	if err != nil {
		return nil, err
	}

	vkClient.Self = token

	return vkClient, nil
}

func NewVKClientWithToken(token string, options *TokenOptions) (*VKClient, error) {
	vkClient := newVKClientBlank()
	vkClient.Self.AccessToken = token

	if options != nil {
		if options.ValidateOnStart {
			uid, err := vkClient.requestSelfID()
			if err != nil {
				return nil, err
			}
			vkClient.Self.UID = uid

			if !options.ServiceToken {
				if err := vkClient.updateSelfUser(); err != nil {
					return nil, err
				}
			}
		}
	}

	return vkClient, nil
}

func (client *VKClient) updateSelfUser() error {
	me, err := client.UsersGet([]int{client.Self.UID})
	if err != nil {
		return err
	}

	client.Self.FirstName = me[0].FirstName
	client.Self.LastName = me[0].LastName
	client.Self.PicSmall = me[0].Photo
	client.Self.PicMedium = me[0].PhotoMedium
	client.Self.PicBig = me[0].PhotoBig

	return nil
}

type AuthConfig struct {
	ClientID, ClientSecret string
	GrantType              string
	Version                string
}

func (client *VKClient) auth(user string, password string, config AuthConfig) (Token, error) {
	client.rl.Wait()
	req, err := http.NewRequest("GET", tokenURL, nil)
	if err != nil {
		return Token{}, err
	}
	client.rl.Update()

	q := req.URL.Query()
	q.Add("grant_type", config.GrantType)
	q.Add("client_id", config.ClientID)
	q.Add("client_secret", config.ClientSecret)
	q.Add("username", user)
	q.Add("password", password)
	q.Add("v", config.Version)
	req.URL.RawQuery = q.Encode()

	client.rl.Wait()
	resp, err := client.Client.Do(req)
	if err != nil {
		return Token{}, err
	}
	client.rl.Update()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Token{}, err
	}

	var token Token
	json.Unmarshal(body, &token)

	if token.Error != "" {
		return Token{}, errors.New(token.Error + ": " + token.ErrorDescription)
	}

	return token, nil
}

func (client *VKClient) requestSelfID() (uid int, err error) {
	resp, err := client.makeRequest("users.get", url.Values{})
	if err != nil {
		return 0, err
	}

	rawdata, err := resp.Response.MarshalJSON()
	if err != nil {
		return 0, err
	}

	data := make([]struct {
		ID int `json:"id"`
	}, 1)

	if err := json.Unmarshal(rawdata, &data); err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, nil
	}

	return data[0].ID, nil
}

func (client *VKClient) makeRequest(method string, params url.Values) (APIResponse, error) {
	client.rl.Wait()

	endpoint := fmt.Sprintf(apiURL, method)
	if params == nil {
		params = url.Values{}
	}

	params.Set("access_token", client.Self.AccessToken)
	params.Set("v", "5.71")

	resp, err := client.Client.PostForm(endpoint, params)
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	client.rl.Update()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	if apiResp.ResponseError.ErrorCode != 0 {
		return APIResponse{}, errors.New("Error code: " + strconv.Itoa(apiResp.ResponseError.ErrorCode) + ", " + apiResp.ResponseError.ErrorMsg)
	}
	return apiResp, nil
}
