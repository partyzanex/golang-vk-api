package vkapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
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

	resp, err := client.MakeRequest("utils.resolveScreenName", params)
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

func BoolToInt(a bool) int {
	if a {
		return 1
	}

	return 0
}

func IntToBool(a int) bool {
	if a > 0 {
		return true
	}

	return false
}

func GetFilesSizeMB(files []string) (int, error) {
	var size int64

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return 0, err
		}
		fi, err := file.Stat()
		if err != nil {
			return 0, err
		}

		size += fi.Size()
		file.Close()
	}

	return int(size / 1048576), nil
}
