package binance

import (
	"encoding/json"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// ExchangeInfoWsService creates exchange info service
type ExchangeInfoWsService struct {
	c websocket.Client
}

// NewExchangeInfoWsService init TimeCheckWsService
func NewExchangeInfoWsService(apiKey, secretKey string) (*ExchangeInfoWsService, error) {
	conn, err := websocket.NewConnection(WsApiInitReadWriteConn, WebsocketKeepalive, WebsocketTimeoutReadWriteConnection)
	if err != nil {
		return nil, err
	}

	client, err := websocket.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &ExchangeInfoWsService{
		c: client,
	}, nil
}

// TimeCheckWsRequest parameters for 'exchangeInfo' websocket API
type ExchangeInfoWsRequest struct {
	symbol             *string
	symbols            *[]string
	permissions        *[]websocket.WsExchangeInfoPermissionType
	showPermissionSets *bool
	symbolStatus       *websocket.WsExchangeInfoSymbolStatusType
}

func (s *ExchangeInfoWsRequest) SetSymbol(symbol string) *ExchangeInfoWsRequest {
	s.symbol = &symbol
	return s
}

func (s *ExchangeInfoWsRequest) SetSymbols(symbols []string) *ExchangeInfoWsRequest {
	s.symbols = &symbols
	return s
}

func (s *ExchangeInfoWsRequest) SetPermissions(permissions []websocket.WsExchangeInfoPermissionType) *ExchangeInfoWsRequest {
	s.permissions = &permissions
	return s
}

func (s *ExchangeInfoWsRequest) SetShowPermissionSets(showPermissionSets bool) *ExchangeInfoWsRequest {
	s.showPermissionSets = &showPermissionSets
	return s
}

func (s *ExchangeInfoWsRequest) SetSymbolStatus(symbolStatus websocket.WsExchangeInfoSymbolStatusType) *ExchangeInfoWsRequest {
	s.symbolStatus = &symbolStatus
	return s
}

func (s *ExchangeInfoWsRequest) GetParams() map[string]interface{} {
	return s.buildParams()
}

// buildParams builds params
func (s *ExchangeInfoWsRequest) buildParams() params {
	m := params{}

	if s.symbol != nil {
		m["symbol"] = *s.symbol
	}
	if s.symbols != nil {
		m["symbols"] = *s.symbols
	}
	if s.permissions != nil {
		m["permissions"] = *s.permissions
	}
	if s.showPermissionSets != nil {
		m["showPermissionSets"] = *s.showPermissionSets
	}
	if s.symbolStatus != nil {
		m["symbolStatus"] = *s.symbolStatus
	}

	return m
}

// NewTimeCheckWsRequest init TimeCheckWsRequest
func NewExchangeInfoWsRequest() *ExchangeInfoWsRequest {
	return &ExchangeInfoWsRequest{}
}

// Do - sends 'exchangeInfo' request
func (s *ExchangeInfoWsService) Do(requestID string, request *ExchangeInfoWsRequest) error {
	rawData, err := websocket.CreateRequestWithSigned(
		requestID,
		websocket.ExchangeInfoWsApiMehod,
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

// SyncDo - sends 'exchangeInfo' request and receives response
func (s *ExchangeInfoWsService) SyncDo(requestID string, request *ExchangeInfoWsRequest) (*ExchangeInfoWsResponse, error) {
	rawData, err := websocket.CreateRequestWithSigned(
		requestID,
		websocket.ExchangeInfoWsApiMehod,
		request.buildParams(),
	)
	if err != nil {
		return nil, err
	}

	response, err := s.c.WriteSync(requestID, rawData, websocket.WriteSyncWsTimeout)
	if err != nil {
		return nil, err
	}

	exchangeInfoWsResponse := &ExchangeInfoWsResponse{}
	if err := json.Unmarshal(response, exchangeInfoWsResponse); err != nil {
		return nil, err
	}

	return exchangeInfoWsResponse, nil
}

// ReceiveAllDataBeforeStop waits until all responses will be received from websocket until timeout expired
func (s *ExchangeInfoWsService) ReceiveAllDataBeforeStop(timeout time.Duration) {
	s.c.Wait(timeout)
}

// GetReadChannel returns channel with API response data (including API errors)
func (s *ExchangeInfoWsService) GetReadChannel() <-chan []byte {
	return s.c.GetReadChannel()
}

// GetReadErrorChannel returns channel with errors which are occurred while reading websocket connection
func (s *ExchangeInfoWsService) GetReadErrorChannel() <-chan error {
	return s.c.GetReadErrorChannel()
}

// GetReconnectCount returns count of reconnect attempts by client
func (s *ExchangeInfoWsService) GetReconnectCount() int64 {
	return s.c.GetReconnectCount()
}

// ExchangeInfoResult define order creation result
type ExchangeInfoResult struct {
	ExchangeInfo
}

// ExchangeInfoWsResponse define 'exchangeInfo' websocket API response
type ExchangeInfoWsResponse struct {
	Id     string             `json:"id"`
	Status int                `json:"status"`
	Result ExchangeInfoResult `json:"result"`

	// error response
	Error *common.APIError `json:"error,omitempty"`
}
