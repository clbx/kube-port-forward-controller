package routers

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"

	"github.com/paultyng/go-unifi/unifi"
)

func CreateUnifiClient(baseurl, username, password string) (*unifi.Client, error) {
	client := &unifi.Client{}
	err := client.SetBaseURL(baseurl)
	jar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Jar:       jar,
		Transport: tr,
	}

	client.SetHTTPClient(httpClient)
	ctx := context.Background()
	err = client.Login(ctx, "clay@clbx.io", "qetjuv-nizpar-suSvo4")
	if err != nil {
		log.Fatalf("Failed to login: %v\n", err)
	}
	fmt.Printf("UniFi Controller Version: %s\n", client.Version())
	return client, nil
}

func CheckPort(client *unifi.Client, port int) (bool, error) {
	portforwards, err := client.ListPortForward(context.TODO(), "default")
	if err != nil {
		return false, err
	}

	for _, portforward := range portforwards {
		portNum, err := strconv.Atoi(portforward.FwdPort)
		if err != nil {
			return false, err
		}
		if portNum == port {
			return true, nil
		}
	}
	return false, nil
}
