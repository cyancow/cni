package client

import (
	"bytes"
	"encoding/json"
	"istio.io/cni/pkg/istioproxyagent/api"
	"net/http"
)

type proxyAgentClient struct {
	httpClient *http.Client
	URL        string
}

func NewProxyAgentClient(URL string) (*proxyAgentClient, error) {
	return &proxyAgentClient{
		httpClient: http.DefaultClient,
		URL:        URL,
	}, nil
}

func (p *proxyAgentClient) StartProxy(podIP, infraContainerID string, secretData map[string][]byte, podJSON, meshConfig, sidecarTemplate string) error {
	return p.callAgent("/start", api.StartRequest{
		podIP,
		infraContainerID,
		secretData,
		podJSON,
		meshConfig,
		sidecarTemplate,
	}, nil)
}

func (p *proxyAgentClient) StopProxy(podName, podSandboxID string) error {
	return p.callAgent("/stop", api.StopRequest{
		podName,
		podSandboxID,
	}, nil)
}

func (p *proxyAgentClient) IsReady(podName string, podNamespace string, podIP string, netNS string) (bool, error) {

	readinessResponse := api.ReadinessResponse{}

	err := p.callAgent("/readiness", api.ReadinessRequest{
		podName,
		podNamespace,
		podIP,
		netNS,
	}, &readinessResponse)

	if err != nil {
		return false, err
	}

	return readinessResponse.Ready, nil
}

func (p *proxyAgentClient) callAgent(path string, request interface{}, responseObj interface{}) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, p.URL+path, bytes.NewReader(b))
	if err != nil {
		return err
	}

	response, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if responseObj != nil {
		decoder := json.NewDecoder(response.Body)
		return decoder.Decode(responseObj)
	}

	// TODO: return error on HTTP error codes
	return nil
}
