package binance

import (
	"encoding/json"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// NewTimeCheckWsService creates order
type TimeCheckWsService struct {
	c websocket.Client
}

// NewTimeCheckWsService init TimeCheckWsService
func NewTimeCheckWsService(apiKey, secretKey string) (*TimeCheckWsService, error) {
	conn, err := websocket.NewConnection(WsApiInitReadWriteConn, WebsocketKeepalive, WebsocketTimeoutReadWriteConnection)
	if err != nil {
		return nil, err
	}

	client, err := websocket.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &TimeCheckWsService{
		c: client,
	}, nil
}

// TimeCheckWsRequest parameters for 'time' websocket API
type TimeCheckWsRequest struct{}

// NewTimeCheckWsRequest init TimeCheckWsRequest
func NewTimeCheckWsRequest() *TimeCheckWsRequest {
	return &TimeCheckWsRequest{}
}

func (s *TimeCheckWsRequest) GetParams() map[string]interface{} {
	return s.buildParams()
}

// buildParams builds params
func (s *TimeCheckWsRequest) buildParams() params {
	m := params{}

	return m
}

// Do - sends 'time' request
func (s *TimeCheckWsService) Do(requestID string, request *TimeCheckWsRequest) error {
	rawData, err := websocket.CreateRequestWithSigned(
		requestID,
		websocket.TimeCheckWsApiMehod,
		request.buildParams(),
	)
	if err != nil {
		return err
	}

	if err := s.c.Write(requestID, rawData); err != nil {
		return err
	}

	return nil
}

// SyncDo - sends 'time' request and receives response
func (s *TimeCheckWsService) SyncDo(requestID string, request *TimeCheckWsRequest) (*TimeCheckWsResponse, error) {
	rawData, err := websocket.CreateRequestWithSigned(
		requestID,
		websocket.TimeCheckWsApiMehod,
		request.buildParams(),
	)
	if err != nil {
		return nil, err
	}

	response, err := s.c.WriteSync(requestID, rawData, websocket.WriteSyncWsTimeout)
	if err != nil {
		return nil, err
	}

	timeCheckWsResponse := &TimeCheckWsResponse{}
	if err := json.Unmarshal(response, timeCheckWsResponse); err != nil {
		return nil, err
	}

	return timeCheckWsResponse, nil
}

// ReceiveAllDataBeforeStop waits until all responses will be received from websocket until timeout expired
func (s *TimeCheckWsService) ReceiveAllDataBeforeStop(timeout time.Duration) {
	s.c.Wait(timeout)
}

// GetReadChannel returns channel with API response data (including API errors)
func (s *TimeCheckWsService) GetReadChannel() <-chan []byte {
	return s.c.GetReadChannel()
}

// GetReadErrorChannel returns channel with errors which are occurred while reading websocket connection
func (s *TimeCheckWsService) GetReadErrorChannel() <-chan error {
	return s.c.GetReadErrorChannel()
}

// GetReconnectCount returns count of reconnect attempts by client
func (s *TimeCheckWsService) GetReconnectCount() int64 {
	return s.c.GetReconnectCount()
}

// TimeCheckResult define order creation result
type TimeCheckResult struct {
	ServerTime int64 `json:"serverTime"`
}

// TimeCheckWsResponse define 'time' websocket API response
type TimeCheckWsResponse struct {
	Id     string          `json:"id"`
	Status int             `json:"status"`
	Result TimeCheckResult `json:"result"`

	// error response
	Error *common.APIError `json:"error,omitempty"`
}
