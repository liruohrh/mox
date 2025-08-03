package mox

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/samber/mo"
)

func TestJsoniter(t *testing.T) {
	jsoniter.RegisterExtension(&OptionExtension{})
	jsoniterApi := jsoniter.Config{}.Froze()
	type User struct {
		Name mo.Option[string] `json:"name,omitempty"`
	}
	user1 := User{
		Name: mo.Some("sb"),
	}
	user2 := User{}

	b1, _ := jsoniterApi.Marshal(user1)
	require.Equal(t, "{\"name\":\"sb\"}", string(b1))
	b2, _ := jsoniterApi.Marshal(user2)
	require.Equal(t, "{}", string(b2))
}

func TestGoJson(t *testing.T) {
	type User struct {
		Name mo.Option[string] `json:"name,omitempty"`
	}
	user1 := User{
		Name: mo.Some("sb"),
	}

	b1, _ := json.Marshal(user1)
	require.Equal(t, "{\"name\":\"sb\"}", string(b1))
	var user11 User
	require.NoError(t, json.Unmarshal(b1, &user11))
	require.Equal(t, user1, user11)

	user2 := User{}
	b2, _ := json.Marshal(user2)
	require.Equal(t, "{\"name\":null}", string(b2))
	var user22 User
	require.NoError(t, json.Unmarshal(b2, &user22))
	require.NotEqual(t, user2, user22)
	require.False(t, user2.Name.IsPresent())
	require.True(t, user22.Name.IsPresent())
	require.Equal(t, "", user2.Name.OrEmpty())
	require.Equal(t, "", user22.Name.OrEmpty())
}
