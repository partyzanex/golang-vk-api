package vkapi_test

import (
	"testing"
	"os"
	"github.com/partyzanex/golang-vk-api"
)

func TestVKClient_UsersGet(t *testing.T) {
	vk, err := vkapi.NewVKClientWithToken(
		os.Getenv("VK_TOKEN"),
		&vkapi.TokenOptions{
			ServiceToken: true,
		},
	)
	if err != nil {
		t.Fatalf("creating vk client failed: %s", err)
	}

	users, err := vk.UsersGet("antony_goroni")
	if err != nil {
		t.Fatalf("getting users failed: %s", err)
	}

	t.Log(users)
}
