package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var client http.Client

func TestAddFriend(t *testing.T) {
	//build a test server
	ts := httptest.NewServer(SetupServer())
	//close the server after all test is done
	defer ts.Close()

	loginResp, err := http.PostForm(fmt.Sprintf("%s/login", ts.URL), url.Values{
		"email":    {"bakcardPing@deadCommunist.com"},
		"password": {"1234"},
		"origin":   {"app"},
	})
	if err != nil {
		t.Fatalf("login error")
	}
	cookies := loginResp.Cookies()

	formBody := url.Values{}
	formBody.Add("fd[]", "2")
	formBody.Add("fd[]", "3")
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/friends", ts.URL), nil)
	for _, cookie := range cookies {
		log.Println("cookie", cookie)
		req.AddCookie(cookie)
	}
	req.Form = formBody
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}
