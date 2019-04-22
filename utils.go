package vkapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ResolveScreenName struct {
	Type     string `json:"type"`
	ObjectID int    `json:"object_id"`
}

func (client *VKClient) ResolveScreenName(name string) (ResolveScreenName, error) {
	var res ResolveScreenName
	params := url.Values{}
	params.Set("screen_name", name)

	resp, err := client.makeRequest("utils.resolveScreenName", params)
	if err == nil {
		json.Unmarshal(resp.Response, &res)
	}
	if res.ObjectID == 0 || res.Type == "" {
		err = fmt.Errorf("%s not found", name)
	}
	return res, err

}

func ArrayToStr(a []int) string {
	s := make([]string, len(a))

	for i, num := range a {
		s[i] = strconv.Itoa(num)
	}

	return strings.Join(s, ",")
}
