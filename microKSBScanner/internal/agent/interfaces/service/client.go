package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	agent "ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/agent/usecases"
	"ksb-dev.keysystems.local/intgrsrv/microKSBScanner/internal/common/structs"
)

type Client struct {
	client *http.Client
}

var _ agent.Service = &Client{}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		client: httpClient,
	}
}

func (c *Client) Ping(settings structs.AgentSettings, ping structs.PingRequestData) (structs.PingResponseData, error) {
	var rv structs.PingResponseData

	pingRequest := structs.PingRequest{
		Version: 1,
		Request: ping,
	}

	serverUrl := settings.GetServerAddress()
	serverUrl.Path = strings.TrimRight(serverUrl.Path, "/") + "/ping"

	pingBytes, err := json.Marshal(pingRequest)
	if err != nil {
		return rv, fmt.Errorf("json.Marshal ping err: %w", err)
	}

	return c.requestServer(serverUrl, pingBytes)
}

func (c *Client) SendData(settings structs.AgentSettings, data structs.AddDataPacket) (structs.PingResponseData, error) {
	var rv structs.PingResponseData

	serverUrl := settings.GetServerAddress()
	serverUrl.Path = strings.TrimRight(serverUrl.Path, "/") + "/data"

	request := structs.AddDataRequset{
		Version: 1,
		ID:      "",
		Request: data,
	}

	dataBytes, err := json.Marshal(request)
	if err != nil {
		return rv, fmt.Errorf("json.Marhsal data err: %w", err)
	}

	return c.requestServer(serverUrl, dataBytes)
}

func (c *Client) requestServer(serverUrl url.URL, dataBytes []byte) (structs.PingResponseData, error) {
	var rv structs.PingResponseData

	resp, err := c.client.Post(serverUrl.String(), "contentType/json", bytes.NewReader(dataBytes))
	if err != nil {
		return rv, fmt.Errorf("client.Get err: %w", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rv, fmt.Errorf("ioutil.ReadAll err: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return rv, fmt.Errorf("server [%s] response error: %s - %s", serverUrl.String(), resp.Status, string(b))
	}

	//log.Printf("response bytes = %s", string(b))
	response := structs.ResponseBody{}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return rv, fmt.Errorf("json.Unmarshal PingResponse err: %w. Bytes: %s", err, string(b))
	}

	if response.StatusCode != 200 {
		return rv, fmt.Errorf("server [%s] response logical status %d error: %s", serverUrl.String(), response.StatusCode, response.Message)
	}
	//log.Printf("response = %#v", response)

	var re struct {
		Body structs.PingResponse `json:"body"`
	}
	err = json.Unmarshal(b, &re)
	if err != nil {
		return rv, fmt.Errorf("json.Unmarshal PingResponse err: %w. Bytes: %s", err, string(b))
	}
	//log.Printf("re = %#v", re)

	return re.Body.Response, nil
}
