package main_test

import (
	"fmt"
	"github.com/alecthomas/units"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
)

func TestMaxHeaderBytes(t *testing.T) {

	client := &http.Client{}

	req, err := http.NewRequest(
		"GET", "http://localhost:8118/api?token=123", nil)
	require.NoError(t, err)

	foo := make([]byte, int(2 * units.KiB))
	_, err = rand.Read(foo)
	require.NoError(t, err)
	req.Header.Set("X-Foo", url.QueryEscape(string(foo)))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer (func() {
		_ = resp.Body.Close()
	})()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	fmt.Println("body: ", string(body))

	require.Equal(t, http.StatusRequestHeaderFieldsTooLarge, resp.StatusCode)
}
