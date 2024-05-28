package routers

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strconv"

	"github.com/paultyng/go-unifi/unifi"
)

type UnifiRouter struct {
	SiteID string
	Client *unifi.Client
}

func CreateUnifiRouter(baseurl, username, password, site string) (*UnifiRouter, error) {
	client := &unifi.Client{}
	err := client.SetBaseURL(baseurl)
	if err != nil {
		return nil, err
	}
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
	err = client.Login(ctx, username, password)
	if err != nil {
		return nil, err
	}
	fmt.Printf("UniFi Controller Version: %s\n", client.Version())

	siteId, err := getSiteID(client, site)
	if err != nil {
		return nil, err
	}

	router := &UnifiRouter{
		SiteID: siteId,
		Client: client,
	}

	return router, nil
}

func (router *UnifiRouter) CheckPort(port int) (bool, error) {
	portforwards, err := router.Client.ListPortForward(context.TODO(), router.SiteID)
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

func (router *UnifiRouter) AddPort(config PortConfig) error {
	portforward := &unifi.PortForward{
		SiteID: router.SiteID,
		//TODO: Create functionality for this, but not needed for MVP
		DestinationIP: "any",
		DstPort:       strconv.Itoa(config.DstPort),
		Enabled:       config.Enabled,
		Fwd:           config.SrcIp,
		FwdPort:       strconv.Itoa(config.DstPort),
		Name:          config.Name,
		PfwdInterface: config.Interface,
		Proto:         config.Protocol,
		//TODO: Create functionality for this, but not needed for MVP
		Src: "any",
	}

	_, err := router.Client.CreatePortForward(context.TODO(), router.SiteID, portforward)
	if err != nil {
		return err
	}
	return nil
}

func getSiteID(client *unifi.Client, siteName string) (string, error) {
	sites, err := client.ListSites(context.TODO())
	if err != nil {
		return "", err
	}

	for _, site := range sites {
		if site.Name == siteName {
			return site.ID, nil
		}
	}

	return "", fmt.Errorf("No site found by name %s", siteName)
}
