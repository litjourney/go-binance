package binance

import (
	"encoding/json"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// NewOrderStatusWsService creates order
type OrderStatusWsService struct {
	c          websocket.Client
	ApiKey     string
	SecretKey  string
	KeyType    string
	TimeOffset int64
}

// NewNewOrderStatusWsService init NewOrderStatusWsService
func NewOrderStatusWsService(apiKey, secretKey string) (*OrderStatusWsService, error) {
	conn, err := websocket.NewConnection(WsApiInitReadWriteConn, WebsocketKeepalive, WebsocketTimeoutReadWriteConnection)
	if err != nil {
		return nil, err
	}

	client, err := websocket.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &OrderStatusWsService{
		c:         client,
		ApiKey:    apiKey,
		SecretKey: secretKey,
		KeyType:   common.KeyTypeHmac,
	}, nil
}

// OrderStatusWsRequest parameters for 'order.status' websocket API
type OrderStatusWsRequest struct {
	symbol            string
	orderId           int64
	origClientOrderId *string
	recvWindow        *uint16
}

// NewOrderStatusWsRequest init OrderStatusWsRequest
func NewOrderStatusWsRequest() *OrderStatusWsRequest {
	return &OrderStatusWsRequest{}
}

func (s *OrderStatusWsRequest) GetParams() map[string]interface{} {
	return s.buildParams()
}

// buildParams builds params
func (s *OrderStatusWsRequest) buildParams() params {
	m := params{
		"symbol": s.symbol,
	}
	if s.orderId != 0 {
		m["orderId"] = s.orderId
	}
	if s.origClientOrderId != nil {
		m["origClientOrderId"] = *s.origClientOrderId
	}
	if s.recvWindow != nil {
		m["recvWindow"] = *s.recvWindow
	}
	return m
}

// Do - sends 'order.status' request
func (s *OrderStatusWsService) Do(requestID string, request *OrderStatusWsRequest) error {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.OrderStatusSpotWsApiMethod,
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

// SyncDo - sends 'order.status' request and receives response
func (s *OrderStatusWsService) SyncDo(requestID string, request *OrderStatusWsRequest) (*StatusOrderWsResponse, error) {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.OrderStatusSpotWsApiMethod,
		request.buildParams(),
	)
	if err != nil {
		return nil, err
	}

	response, err := s.c.WriteSync(requestID, rawData, websocket.WriteSyncWsTimeout)
	if err != nil {
		return nil, err
	}

	statusOrderWsResponse := &StatusOrderWsResponse{}
	if err := json.Unmarshal(response, statusOrderWsResponse); err != nil {
		return nil, err
	}

	return statusOrderWsResponse, nil
}

// ReceiveAllDataBeforeStop waits until all responses will be received from websocket until timeout expired
func (s *OrderStatusWsService) ReceiveAllDataBeforeStop(timeout time.Duration) {
	s.c.Wait(timeout)
}

// GetReadChannel returns channel with API response data (including API errors)
func (s *OrderStatusWsService) GetReadChannel() <-chan []byte {
	return s.c.GetReadChannel()
}

// GetReadErrorChannel returns channel with errors which are occurred while reading websocket connection
func (s *OrderStatusWsService) GetReadErrorChannel() <-chan error {
	return s.c.GetReadErrorChannel()
}

// GetReconnectCount returns count of reconnect attempts by client
func (s *OrderStatusWsService) GetReconnectCount() int64 {
	return s.c.GetReconnectCount()
}

// Symbol set symbol
func (s *OrderStatusWsRequest) Symbol(symbol string) *OrderStatusWsRequest {
	s.symbol = symbol
	return s
}

// RecvWindow set recvWindow
func (s *OrderStatusWsRequest) RecvWindow(recvWindow uint16) *OrderStatusWsRequest {
	s.recvWindow = &recvWindow
	return s
}

func (s *OrderStatusWsRequest) OrderID(orderId int64) *OrderStatusWsRequest {
	s.orderId = orderId
	return s
}

func (s *OrderStatusWsRequest) OrigClientOrderID(origClientOrderId string) *OrderStatusWsRequest {
	s.origClientOrderId = &origClientOrderId
	return s
}

// CreateOrderResult define order creation result
type StatusOpenOrderResult *[]Order

// CreateOrderWsResponse define 'order.status' websocket API response
type StatusOrderWsResponse struct {
	Id     string                `json:"id"`
	Status int                   `json:"status"`
	Result StatusOpenOrderResult `json:"result"`

	// error response
	Error *common.APIError `json:"error,omitempty"`
}
