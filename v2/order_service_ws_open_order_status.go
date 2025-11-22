package binance

import (
	"encoding/json"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// OpenOrderStatusWsService query open order
type OpenOrderStatusWsService struct {
	c          websocket.Client
	ApiKey     string
	SecretKey  string
	KeyType    string
	TimeOffset int64
}

// NewNewOrderStatusWsService init NewOrderStatusWsService
func NewOpenOrderStatusWsService(apiKey, secretKey string) (*OpenOrderStatusWsService, error) {
	conn, err := websocket.NewConnection(WsApiInitReadWriteConn, WebsocketKeepalive, WebsocketTimeoutReadWriteConnection)
	if err != nil {
		return nil, err
	}

	client, err := websocket.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &OpenOrderStatusWsService{
		c:         client,
		ApiKey:    apiKey,
		SecretKey: secretKey,
		KeyType:   common.KeyTypeHmac,
	}, nil
}

// OpenOpenOrderStatusWsRequest parameters for 'openOrders.status' websocket API
type OpenOrderStatusWsRequest struct {
	symbol     string
	recvWindow *uint16
}

// NewOpenOrderStatusWsRequest init OpenOrderStatusWsRequest
func NewOpenOpenOrderStatusWsRequest() *OpenOrderStatusWsRequest {
	return &OpenOrderStatusWsRequest{}
}

func (s *OpenOrderStatusWsRequest) GetParams() map[string]interface{} {
	return s.buildParams()
}

// buildParams builds params
func (s *OpenOrderStatusWsRequest) buildParams() params {
	m := params{
		"symbol": s.symbol,
	}
	if s.recvWindow != nil {
		m["recvWindow"] = *s.recvWindow
	}
	return m
}

// Do - sends 'openOrder.status' request
func (s *OpenOrderStatusWsService) Do(requestID string, request *OpenOrderStatusWsRequest) error {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.OrderOpenStatusSpotWsApiMethod,
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

// SyncDo - sends 'openOrder.status' request and receives response
func (s *OpenOrderStatusWsService) SyncDo(requestID string, request *OpenOrderStatusWsRequest) (*StatusOpenOrderWsResponse, error) {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.OrderOpenStatusSpotWsApiMethod,
		request.buildParams(),
	)
	if err != nil {
		return nil, err
	}

	response, err := s.c.WriteSync(requestID, rawData, websocket.WriteSyncWsTimeout)
	if err != nil {
		return nil, err
	}

	statusOpenOrderWsResponse := &StatusOpenOrderWsResponse{}
	if err := json.Unmarshal(response, statusOpenOrderWsResponse); err != nil {
		return nil, err
	}

	return statusOpenOrderWsResponse, nil
}

// ReceiveAllDataBeforeStop waits until all responses will be received from websocket until timeout expired
func (s *OpenOrderStatusWsService) ReceiveAllDataBeforeStop(timeout time.Duration) {
	s.c.Wait(timeout)
}

// GetReadChannel returns channel with API response data (including API errors)
func (s *OpenOrderStatusWsService) GetReadChannel() <-chan []byte {
	return s.c.GetReadChannel()
}

// GetReadErrorChannel returns channel with errors which are occurred while reading websocket connection
func (s *OpenOrderStatusWsService) GetReadErrorChannel() <-chan error {
	return s.c.GetReadErrorChannel()
}

// GetReconnectCount returns count of reconnect attempts by client
func (s *OpenOrderStatusWsService) GetReconnectCount() int64 {
	return s.c.GetReconnectCount()
}

// Symbol set symbol
func (s *OpenOrderStatusWsRequest) Symbol(symbol string) *OpenOrderStatusWsRequest {
	s.symbol = symbol
	return s
}

// RecvWindow set recvWindow
func (s *OpenOrderStatusWsRequest) RecvWindow(recvWindow uint16) *OpenOrderStatusWsRequest {
	s.recvWindow = &recvWindow
	return s
}

type StatusOpenOrdersResult []*Order

// CreateOrderWsResponse define 'openOrders.status' websocket API response
type StatusOpenOrderWsResponse struct {
	Id     string                 `json:"id"`
	Status int                    `json:"status"`
	Result StatusOpenOrdersResult `json:"result"`

	// error response
	Error *common.APIError `json:"error,omitempty"`
}
