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

type allOrdersServiceWsTestSuite struct {
	suite.Suite
	apiKey     string
	secretKey  string
	signedKey  string
	timeOffset int64

	ctrl   *gomock.Controller
	client *mock.MockClient

	requestID  string
	symbol     string
	orderId    int64
	startTime  int64
	endTime    int64
	limit      int
	recvWindow uint16

	service *AllOrdersWsService
	request *AllOrdersWsRequest
}

func (s *allOrdersServiceWsTestSuite) SetupTest() {
	s.apiKey = "dummyApiKey"
	s.secretKey = "dummySecretKey"
	s.signedKey = "HMAC"
	s.timeOffset = 0

	s.requestID = "e2a85d9f-07a5-4f94-8d5f-789dc3deb098"
	s.symbol = "BTCUSDT"
	s.orderId = 12345
	s.startTime = 1600000000000
	s.endTime = 1600000001000
	s.limit = 10
	s.recvWindow = 5000

	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)

	s.service = &AllOrdersWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.request = NewAllOrdersWsRequest().
		Symbol(s.symbol).
		OrderId(s.orderId).
		StartTime(s.startTime).
		EndTime(s.endTime).
		Limit(s.limit).
		RecvWindow(s.recvWindow)
}

func (s *allOrdersServiceWsTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestAllOrdersServiceWs(t *testing.T) {
	suite.Run(t, new(allOrdersServiceWsTestSuite))
}

func (s *allOrdersServiceWsTestSuite) TestDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.service.Do(s.requestID, s.request)
	s.NoError(err)
}

func (s *allOrdersServiceWsTestSuite) TestDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	err := s.service.Do("", s.request)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *allOrdersServiceWsTestSuite) TestDo_EmptyApiKey() {
	s.reset("", s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorApiKeyIsNotSet)
}

func (s *allOrdersServiceWsTestSuite) TestDo_EmptySecretKey() {
	s.reset(s.apiKey, "", s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorSecretKeyIsNotSet)
}

func (s *allOrdersServiceWsTestSuite) TestSyncDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	expectedResponse := AllOrderWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: AllOrderResult{
			{
				Symbol:  s.symbol,
				OrderID: s.orderId,
			},
		},
	}

	rawResponseData, err := json.Marshal(expectedResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	response, err := s.service.SyncDo(s.requestID, s.request)
	s.Require().NoError(err)
	s.Equal(s.symbol, response.Result[0].Symbol)
	s.Equal(s.orderId, response.Result[0].OrderID)
}

func (s *allOrdersServiceWsTestSuite) TestSyncDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().
		WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)

	response, err := s.service.SyncDo("", s.request)
	s.Nil(response)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *allOrdersServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.service = &AllOrdersWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}
