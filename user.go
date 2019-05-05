package vkapi

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"strconv"
)

//last seen device
const (
	_               = iota
	PlatformMobile
	PlatformIPhone
	PlatfromIPad
	PlatformAndroid
	PlatformWPhone
	PlatformWindows
	PlatformWeb
)

var (
	userFields = "photo_id, verified, sex, bdate, city, country, " +
		"home_town, has_photo, photo_50, photo_100, photo_200_orig, " +
		"photo_200, photo_400_orig, photo_max, photo_max_orig, online, " +
		"domain, has_mobile, contacts, site, education, universities, " +
		"schools, status, last_seen, followers_count, " +
		"occupation, nickname, relatives, relation, personal, connections, " +
		"exports, activities, interests, music, movies, tv, books, games, " +
		"about, quotes, can_post, can_see_all_posts, can_see_audio, " +
		"can_write_private_message, can_send_friend_request, is_favorite, " +
		"is_hidden_from_feed, timezone, screen_name, maiden_name, " +
		"crop_photo, is_friend, friend_status, career, military, " +
		"blacklisted, blacklisted_by_me"
)

type User struct {
	UID                     int          `json:"id"`
	FirstName               string       `json:"first_name"`
	LastName                string       `json:"last_name"`
	IsClosed                bool         `json:"is_closed"`
	CanAccessClosed         bool         `json:"can_access_closed"`
	Sex                     int          `json:"sex"`
	Nickname                string       `json:"nickname"`
	Domain                  string       `json:"domain"`
	ScreenName              string       `json:"screen_name"`
	BDate                   string       `json:"bdate"`
	City                    *UserCity    `json:"city"`
	Country                 *UserCountry `json:"country"`
	Photo100                string       `json:"photo_100"`
	Photo200                string       `json:"photo_200_orig"`
	PhotoMax                string       `json:"photo_max"`
	PhotoID                 string       `json:"photo_id"`
	HasPhoto                int          `json:"has_photo"`
	HasMobile               int          `json:"has_mobile"`
	IsFriend                int          `json:"is_friend"`
	FriendStatus            int          `json:"friend_status"`
	Online                  int          `json:"online"`
	CanPost                 int          `json:"can_post"`
	CanSeeAllPosts          int          `json:"can_see_all_posts"`
	CanSeeAudio             int          `json:"can_see_audio"`
	CanWritePrivateMessages int          `json:"can_write_private_message"`
	CanSendFriendRequest    int          `json:"can_send_friend_request"`
	HomePhone               string       `json:"home_phone"`
	Site                    string       `json:"site"`
	Status                  string       `json:"activity"`
	LastSeen                *LastSeen    `json:"last_seen"`
	Verified                int          `json:"verified"`
	FollowersCount          int          `json:"followers_count"`
	BlackListed             int          `json:"black_listed"`
	BlackListedByMe         int          `json:"black_listed_by_me"`
	IsHiddenFromFeed        int          `json:"is_hidden_from_feed"`
	Occupation              *Occupation  `json:"occupation"`
}

type UserCity struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type UserCountry struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type LastSeen struct {
	Time     int64 `json:"time"`
	Platform int   `json:"platform"`
}

type Occupation struct {
	Type string     `json:"type"`
	ID   FloatToInt `json:"id"`
	Name string     `json:"name"`
}

type FloatToInt int

func (f *FloatToInt) UnmarshalJSON(b []byte) error {
	str := string(b)
	strSlice := strings.Split(str, ".")
	fl, err := strconv.ParseInt(strSlice[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "cannot parse float")
	}

	*f = FloatToInt(fl)

	return nil
}

func (client *VKClient) UsersGet(users ...string) ([]*User, error) {
	v := make(url.Values)
	v.Add("user_ids", strings.Join(users, ","))
	v.Add("fields", userFields)

	resp, err := client.makeRequest("users.get", v)
	if err != nil {
		return nil, errors.Wrap(err, "make request failed")
	}

	var userList []*User

	err = json.Unmarshal(resp.Response, &userList)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling failed")
	}

	return userList, nil
}
