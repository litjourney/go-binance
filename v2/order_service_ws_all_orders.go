package binance

import (
	"encoding/json"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// AllOrdersWsService query all orders
type AllOrdersWsService struct {
	c          websocket.Client
	ApiKey     string
	SecretKey  string
	KeyType    string
	TimeOffset int64
}

// NewNewOrderStatusWsService init NewOrderStatusWsService
func NewAllOrdersWsService(apiKey, secretKey string) (*AllOrdersWsService, error) {
	conn, err := websocket.NewConnection(WsApiInitReadWriteConn, WebsocketKeepalive, WebsocketTimeoutReadWriteConnection)
	if err != nil {
		return nil, err
	}

	client, err := websocket.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &AllOrdersWsService{
		c:         client,
		ApiKey:    apiKey,
		SecretKey: secretKey,
		KeyType:   common.KeyTypeHmac,
	}, nil
}

// AllOrdersWsRequest parameters for 'allOrders' websocket API
type AllOrdersWsRequest struct {
	symbol     string
	orderId    *int64
	startTime  *int64
	endTime    *int64
	limit      *int
	recvWindow *uint16
}

// NewAllOrdersWsRequest init AllOrdersWsRequest
func NewAllOrdersWsRequest() *AllOrdersWsRequest {
	return &AllOrdersWsRequest{}
}

// Symbol set symbol
func (s *AllOrdersWsRequest) Symbol(symbol string) *AllOrdersWsRequest {
	s.symbol = symbol
	return s
}

// RecvWindow set recvWindow
func (s *AllOrdersWsRequest) RecvWindow(recvWindow uint16) *AllOrdersWsRequest {
	s.recvWindow = &recvWindow
	return s
}

func (s *AllOrdersWsRequest) OrderId(orderId *int64) *AllOrdersWsRequest {
	s.orderId = orderId
	return s
}

func (s *AllOrdersWsRequest) StartTime(startTime *int64) *AllOrdersWsRequest {
	s.startTime = startTime
	return s
}

func (s *AllOrdersWsRequest) EndTime(endTime *int64) *AllOrdersWsRequest {
	s.endTime = endTime
	return s
}

func (s *AllOrdersWsRequest) Limit(limit *int) *AllOrdersWsRequest {
	s.limit = limit
	return s
}

func (s *AllOrdersWsRequest) GetParams() map[string]interface{} {
	return s.buildParams()
}

// buildParams builds params
func (s *AllOrdersWsRequest) buildParams() params {
	m := params{
		"symbol": s.symbol,
	}
	if s.orderId != nil {
		m["orderId"] = *s.orderId
	}
	if s.startTime != nil {
		m["startTime"] = *s.startTime
	}
	if s.endTime != nil {
		m["endTime"] = *s.endTime
	}
	if s.limit != nil {
		m["limit"] = *s.limit
	}
	if s.recvWindow != nil {
		m["recvWindow"] = *s.recvWindow
	}
	return m
}

// Do - sends 'allOrders' request
func (s *AllOrdersWsService) Do(requestID string, request *AllOrdersWsRequest) error {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.AllOrdersSpotWsApiMethod,
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

// SyncDo - sends 'allOrders' request and receives response
func (s *AllOrdersWsService) SyncDo(requestID string, request *AllOrdersWsRequest) (*AllOrderWsResponse, error) {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.AllOrdersSpotWsApiMethod,
		request.buildParams(),
	)
	if err != nil {
		return nil, err
	}

	response, err := s.c.WriteSync(requestID, rawData, websocket.WriteSyncWsTimeout)
	if err != nil {
		return nil, err
	}

	allOrderWsResponse := &AllOrderWsResponse{}
	if err := json.Unmarshal(response, allOrderWsResponse); err != nil {
		return nil, err
	}

	return allOrderWsResponse, nil
}

// ReceiveAllDataBeforeStop waits until all responses will be received from websocket until timeout expired
func (s *AllOrdersWsService) ReceiveAllDataBeforeStop(timeout time.Duration) {
	s.c.Wait(timeout)
}

// GetReadChannel returns channel with API response data (including API errors)
func (s *AllOrdersWsService) GetReadChannel() <-chan []byte {
	return s.c.GetReadChannel()
}

// GetReadErrorChannel returns channel with errors which are occurred while reading websocket connection
func (s *AllOrdersWsService) GetReadErrorChannel() <-chan error {
	return s.c.GetReadErrorChannel()
}

// GetReconnectCount returns count of reconnect attempts by client
func (s *AllOrdersWsService) GetReconnectCount() int64 {
	return s.c.GetReconnectCount()
}

// CreateOrderResult define order creation result
type AllOrderResult []*Order

// CreateOrderWsResponse define 'allOrders' websocket API response
type AllOrderWsResponse struct {
	Id     string         `json:"id"`
	Status int            `json:"status"`
	Result AllOrderResult `json:"result"`

	// error response
	Error *common.APIError `json:"error,omitempty"`
}
