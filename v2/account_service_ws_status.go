package binance

import (
	"encoding/json"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// NewAccountStatusWsService creates account status websocket service
type AccountStatusWsService struct {
	c          websocket.Client
	ApiKey     string
	SecretKey  string
	KeyType    string
	TimeOffset int64
}

// NewAccountStatusWsService init NewAccountStatusWsService
func NewAccountStatusWsService(apiKey, secretKey string) (*AccountStatusWsService, error) {
	conn, err := websocket.NewConnection(WsApiInitReadWriteConn, WebsocketKeepalive, WebsocketTimeoutReadWriteConnection)
	if err != nil {
		return nil, err
	}

	client, err := websocket.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &AccountStatusWsService{
		c:         client,
		ApiKey:    apiKey,
		SecretKey: secretKey,
		KeyType:   common.KeyTypeHmac,
	}, nil
}

// AccountStatusWsRequest parameters for 'account.status' websocket API
type AccountStatusWsRequest struct {
	omitZeroBalances *bool
	recvWindow       *uint16
}

// NewAccountStatusWsRequest init AccountStatusWsRequest
func NewAccountStatusWsRequest() *AccountStatusWsRequest {
	return &AccountStatusWsRequest{}
}

func (s *AccountStatusWsRequest) GetParams() map[string]interface{} {
	return s.buildParams()
}

// buildParams builds params
func (s *AccountStatusWsRequest) buildParams() params {
	m := params{}
	if s.omitZeroBalances != nil {
		m["omitZeroBalances"] = *s.omitZeroBalances
	}
	if s.recvWindow != nil {
		m["recvWindow"] = *s.recvWindow
	}
	return m
}

// RecvWindow set recvWindow
func (s *AccountStatusWsRequest) RecvWindow(recvWindow uint16) *AccountStatusWsRequest {
	s.recvWindow = &recvWindow
	return s
}

func (s *AccountStatusWsRequest) OmitZeroBalances(omitZeroBalances bool) *AccountStatusWsRequest {
	s.omitZeroBalances = &omitZeroBalances
	return s
}

// Do - sends 'account.status' request
func (s *AccountStatusWsService) Do(requestID string, request *AccountStatusWsRequest) error {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.AccountStatusWsApiMethod,
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

// SyncDo - sends 'account.status' request and receives response
func (s *AccountStatusWsService) SyncDo(requestID string, request *AccountStatusWsRequest) (*AccountStatusWsResponse, error) {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		websocket.AccountStatusWsApiMethod,
		request.buildParams(),
	)
	if err != nil {
		return nil, err
	}

	response, err := s.c.WriteSync(requestID, rawData, websocket.WriteSyncWsTimeout)
	if err != nil {
		return nil, err
	}

	accountStatusWsResponse := &AccountStatusWsResponse{}
	if err := json.Unmarshal(response, accountStatusWsResponse); err != nil {
		return nil, err
	}

	return accountStatusWsResponse, nil
}

// ReceiveAllDataBeforeStop waits until all responses will be received from websocket until timeout expired
func (s *AccountStatusWsService) ReceiveAllDataBeforeStop(timeout time.Duration) {
	s.c.Wait(timeout)
}

// GetReadChannel returns channel with API response data (including API errors)
func (s *AccountStatusWsService) GetReadChannel() <-chan []byte {
	return s.c.GetReadChannel()
}

// GetReadErrorChannel returns channel with errors which are occurred while reading websocket connection
func (s *AccountStatusWsService) GetReadErrorChannel() <-chan error {
	return s.c.GetReadErrorChannel()
}

// GetReconnectCount returns count of reconnect attempts by client
func (s *AccountStatusWsService) GetReconnectCount() int64 {
	return s.c.GetReconnectCount()
}

// AccountStatusResult define account status result
type AccountStatusResult struct {
	Account
}

// AccountStatusWsResponse define 'account.status' websocket API response
type AccountStatusWsResponse struct {
	Id     string              `json:"id"`
	Status int                 `json:"status"`
	Result AccountStatusResult `json:"result"`

	// error response
	Error *common.APIError `json:"error,omitempty"`
}
