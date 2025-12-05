package binance

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/adshao/go-binance/v2/common/websocket"
	"github.com/adshao/go-binance/v2/common/websocket/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type openOrderStatusServiceWsTestSuite struct {
	suite.Suite
	apiKey     string
	secretKey  string
	signedKey  string
	timeOffset int64

	ctrl   *gomock.Controller
	client *mock.MockClient

	requestID  string
	symbol     string
	recvWindow uint16

	service *OpenOrderStatusWsService
	request *OpenOrderStatusWsRequest
}

func (s *openOrderStatusServiceWsTestSuite) SetupTest() {
	s.apiKey = "dummyApiKey"
	s.secretKey = "dummySecretKey"
	s.signedKey = "HMAC"
	s.timeOffset = 0

	s.requestID = "e2a85d9f-07a5-4f94-8d5f-789dc3deb098"
	s.symbol = "BTCUSDT"
	s.recvWindow = 5000

	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)

	s.service = &OpenOrderStatusWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.request = NewOpenOrderStatusWsRequest().
		Symbol(&s.symbol).
		RecvWindow(s.recvWindow)
}

func (s *openOrderStatusServiceWsTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestOpenOrderStatusServiceWs(t *testing.T) {
	suite.Run(t, new(openOrderStatusServiceWsTestSuite))
}

func (s *openOrderStatusServiceWsTestSuite) TestDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.service.Do(s.requestID, s.request)
	s.NoError(err)
}

func (s *openOrderStatusServiceWsTestSuite) TestDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	err := s.service.Do("", s.request)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *openOrderStatusServiceWsTestSuite) TestDo_EmptyApiKey() {
	s.reset("", s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorApiKeyIsNotSet)
}

func (s *openOrderStatusServiceWsTestSuite) TestDo_EmptySecretKey() {
	s.reset(s.apiKey, "", s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorSecretKeyIsNotSet)
}

func (s *openOrderStatusServiceWsTestSuite) TestSyncDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	expectedResponse := StatusOpenOrderWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: StatusOpenOrdersResult{
			{
				Symbol: s.symbol,
			},
		},
	}

	rawResponseData, err := json.Marshal(expectedResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	response, err := s.service.SyncDo(s.requestID, s.request)
	s.Require().NoError(err)
	s.Equal(s.symbol, response.Result[0].Symbol)
}

func (s *openOrderStatusServiceWsTestSuite) TestSyncDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().
		WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)

	response, err := s.service.SyncDo("", s.request)
	s.Nil(response)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *openOrderStatusServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.service = &OpenOrderStatusWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}
