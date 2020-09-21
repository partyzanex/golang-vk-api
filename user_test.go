package vkapi_test

import (
	"github.com/partyzanex/golang-vk-api"
	"os"
	"testing"
)

func TestVKClient_UsersGet(t *testing.T) {
	vk, err := vkapi.NewVKClientWithToken(
		os.Getenv("VK_TOKEN"),
		&vkapi.TokenOptions{
			ServiceToken: true,
		},
		true,
	)
	if err != nil {
		t.Fatalf("creating vk client failed: %s", err)
	}

	users, err := vk.UsersGet([]int{294110051})
	if err != nil {
		t.Fatalf("getting users failed: %s", err)
	}

	t.Log(users)
}
