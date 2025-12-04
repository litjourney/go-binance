package binance

import (
	"encoding/json"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// NewOrderCancelWsService creates order
type OrderCancelWsService struct {
	c          websocket.Client
	ApiKey     string
	SecretKey  string
	KeyType    string
	TimeOffset int64
}

// NewNewOrderCancelWsService init NewOrderCancelWsService
func NewOrderCancelWsService(apiKey, secretKey string) (*OrderCancelWsService, error) {
	conn, err := websocket.NewConnection(WsApiInitReadWriteConn, WebsocketKeepalive, WebsocketTimeoutReadWriteConnection)
	if err != nil {
		return nil, err
	}

	client, err := websocket.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &OrderCancelWsService{
		c:         client,
		ApiKey:    apiKey,
		SecretKey: secretKey,
		KeyType:   common.KeyTypeHmac,
	}, nil
}

// OrderCancelWsRequest parameters for 'order.cancel' websocket API
type OrderCancelWsRequest struct {
	symbol             string
	orderId            int64
	origClientOrderId  *string
	newClientOrderId   *string
	cancelRestrictions *OrderCancelRestrictionsType
	recvWindow         *uint16
}

// NewOrderCancelWsRequest init OrderCancelWsRequest
func NewOrderCancelWsRequest() *OrderCancelWsRequest {
	return &OrderCancelWsRequest{}
}

func (s *OrderCancelWsRequest) GetParams() map[string]interface{} {
	return s.buildParams()
}

// buildParams builds params
func (s *OrderCancelWsRequest) buildParams() params {
	m := params{
		"symbol": s.symbol,
	}
	if s.orderId != 0 {
		m["orderId"] = s.orderId
	}
	if s.origClientOrderId != nil {
		m["origClientOrderId"] = *s.origClientOrderId
	}
	if s.newClientOrderId != nil {
		m["newClientOrderId"] = *s.newClientOrderId
	}
	if s.cancelRestrictions != nil {
		m["cancelRestrictions"] = *s.cancelRestrictions
	}
	if s.recvWindow != nil {
		m["recvWindow"] = *s.recvWindow
	}
	return m
}

// Do - sends 'order.cancel' request
func (s *OrderCancelWsService) Do(requestID string, request *OrderCancelWsRequest) error {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.OrderCancelSpotWsApiMethod,
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

// SyncDo - sends 'order.cancel' request and receives response
func (s *OrderCancelWsService) SyncDo(requestID string, request *OrderCancelWsRequest) (*CancelOrderWsResponse, error) {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.OrderCancelSpotWsApiMethod,
		request.buildParams(),
	)
	if err != nil {
		return nil, err
	}

	response, err := s.c.WriteSync(requestID, rawData, websocket.WriteSyncWsTimeout)
	if err != nil {
		return nil, err
	}

	cancelOrderWsResponse := &CancelOrderWsResponse{}
	if err := json.Unmarshal(response, cancelOrderWsResponse); err != nil {
		return nil, err
	}

	return cancelOrderWsResponse, nil
}

// ReceiveAllDataBeforeStop waits until all responses will be received from websocket until timeout expired
func (s *OrderCancelWsService) ReceiveAllDataBeforeStop(timeout time.Duration) {
	s.c.Wait(timeout)
}

// GetReadChannel returns channel with API response data (including API errors)
func (s *OrderCancelWsService) GetReadChannel() <-chan []byte {
	return s.c.GetReadChannel()
}

// GetReadErrorChannel returns channel with errors which are occurred while reading websocket connection
func (s *OrderCancelWsService) GetReadErrorChannel() <-chan error {
	return s.c.GetReadErrorChannel()
}

// GetReconnectCount returns count of reconnect attempts by client
func (s *OrderCancelWsService) GetReconnectCount() int64 {
	return s.c.GetReconnectCount()
}

// Symbol set symbol
func (s *OrderCancelWsRequest) Symbol(symbol string) *OrderCancelWsRequest {
	s.symbol = symbol
	return s
}

// RecvWindow set recvWindow
func (s *OrderCancelWsRequest) RecvWindow(recvWindow uint16) *OrderCancelWsRequest {
	s.recvWindow = &recvWindow
	return s
}

func (s *OrderCancelWsRequest) OrderID(orderId int64) *OrderCancelWsRequest {
	s.orderId = orderId
	return s
}

func (s *OrderCancelWsRequest) OrigClientOrderID(origClientOrderId *string) *OrderCancelWsRequest {
	s.origClientOrderId = origClientOrderId
	return s
}

func (s *OrderCancelWsRequest) NewClientOrderID(newClientOrderId *string) *OrderCancelWsRequest {
	s.newClientOrderId = newClientOrderId
	return s
}

func (s *OrderCancelWsRequest) CancelRestrictions(cancelRestrictions OrderCancelRestrictionsType) *OrderCancelWsRequest {
	s.cancelRestrictions = &cancelRestrictions
	return s
}

// CreateOrderResult define order creation result
type CancelOrderResult struct {
	CancelOrderResponse
}

// CreateOrderWsResponse define 'order.cancel' websocket API response
type CancelOrderWsResponse struct {
	Id     string            `json:"id"`
	Status int               `json:"status"`
	Result CancelOrderResult `json:"result"`

	// error response
	Error *common.APIError `json:"error,omitempty"`
}
