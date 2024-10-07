package vcenter

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Model struct {
	User string
	Pass string
	Host string
	Port string
	DC   string
}
type Datastore struct {
	Datastore string `json:"datastore"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	FreeSpace int64  `json:"free_space"`
	Capacity  int64  `json:"capacity"`
}

type Host struct {
	Host            string `json:"host"`
	Name            string `json:"name"`
	ConnectionState string `json:"connection_state"`
	PowerState      string `json:"power_state"`
}

type VcenterDatacenter struct {
	Name       string `json:"name"`
	Datacenter string `json:"datacenter"`
}

func (m *Model) GetDatastores(sessionId string) ([]Datastore, error) {
	var msg []Datastore

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	req, err := http.NewRequest("GET", "https://"+addr+"/api/vcenter/datastore?datacenters="+m.DC, nil)
	if err != nil {
		return msg, err
	}

	req.Close = true

	sessionId = strings.Trim(sessionId, "\"")
	req.Header.Set("vmware-api-session-id", sessionId)

	resp, err := client.Do(req)
	if err != nil {
		return msg, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("http response error: %v", resp.StatusCode)
		return msg, err
	}

	bodyText, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(bodyText, &msg)
	if err != nil {

		return msg, err
	}

	return msg, nil
}

func (m *Model) GetDatacenter(sessionId string) ([]VcenterDatacenter, error) {
	var msg []VcenterDatacenter

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	req, err := http.NewRequest("GET", "https://"+addr+"/api/vcenter/datacenter?datacenters="+m.DC, nil)
	if err != nil {
		return msg, err
	}

	req.Close = true

	sessionId = strings.Trim(sessionId, "\"")
	req.Header.Set("vmware-api-session-id", sessionId)

	resp, err := client.Do(req)
	if err != nil {
		return msg, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("http response error: %v", resp.StatusCode)
		return msg, err
	}

	bodyText, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(bodyText, &msg)
	if err != nil {

		return msg, err
	}

	return msg, nil
}

func (m *Model) GetHosts(sessionId string) ([]Host, error) {

	var host []Host

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	req, err := http.NewRequest("GET", "https://"+addr+"/api/vcenter/host?datacenters="+m.DC, nil)
	if err != nil {
		return host, err
	}

	req.Close = true

	sessionId = strings.Trim(sessionId, "\"")
	req.Header.Set("vmware-api-session-id", sessionId)

	resp, err := client.Do(req)
	if err != nil {
		return host, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("http response error: %v", resp.StatusCode)
		return host, err
	}

	bodyText, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(bodyText, &host)
	if err != nil {

		return host, err
	}

	return host, nil
}

func (m *Model) Authenticate() (string, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	req, err := http.NewRequest("POST", "https://"+addr+"/api/session", nil)
	if err != nil {
		return "", err
	}

	req.Close = true

	req.Header.Add("Authorization", "Basic "+basicAuth(m.User, m.Pass))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		err := fmt.Errorf("http response error: %v", resp.StatusCode)
		return "", err
	}

	bodyText, err := io.ReadAll(resp.Body)
	sessionId := string(bodyText)

	return sessionId, nil
}

func (m *Model) LogOut(sessionId string) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	req, err := http.NewRequest("DELETE", "https://"+addr+"/api/session", nil)
	if err != nil {
		return
	}

	req.Close = true

	sessionId = strings.Trim(sessionId, "\"")
	req.Header.Set("vmware-api-session-id", sessionId)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	return
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
