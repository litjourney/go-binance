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

type orderStatusServiceWsTestSuite struct {
	suite.Suite
	apiKey     string
	secretKey  string
	signedKey  string
	timeOffset int64

	ctrl   *gomock.Controller
	client *mock.MockClient

	requestID         string
	symbol            string
	orderId           int64
	origClientOrderId string
	recvWindow        uint16

	service *OrderStatusWsService
	request *OrderStatusWsRequest
}

func (s *orderStatusServiceWsTestSuite) SetupTest() {
	s.apiKey = "dummyApiKey"
	s.secretKey = "dummySecretKey"
	s.signedKey = "HMAC"
	s.timeOffset = 0

	s.requestID = "e2a85d9f-07a5-4f94-8d5f-789dc3deb098"
	s.symbol = "BTCUSDT"
	s.orderId = 12345
	s.origClientOrderId = "origClientOrderId"
	s.recvWindow = 5000

	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)

	s.service = &OrderStatusWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.request = NewOrderStatusWsRequest().
		Symbol(s.symbol).
		OrderID(s.orderId).
		OrigClientOrderID(s.origClientOrderId).
		RecvWindow(s.recvWindow)
}

func (s *orderStatusServiceWsTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestOrderStatusServiceWs(t *testing.T) {
	suite.Run(t, new(orderStatusServiceWsTestSuite))
}

func (s *orderStatusServiceWsTestSuite) TestDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.service.Do(s.requestID, s.request)
	s.NoError(err)
}

func (s *orderStatusServiceWsTestSuite) TestDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	err := s.service.Do("", s.request)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *orderStatusServiceWsTestSuite) TestDo_EmptyApiKey() {
	s.reset("", s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorApiKeyIsNotSet)
}

func (s *orderStatusServiceWsTestSuite) TestDo_EmptySecretKey() {
	s.reset(s.apiKey, "", s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorSecretKeyIsNotSet)
}

func (s *orderStatusServiceWsTestSuite) TestSyncDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	expectedResponse := StatusOrderWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: StatusOpenOrderResult(
			&[]Order{
				{
					Symbol:  s.symbol,
					OrderID: s.orderId,
				},
			},
		),
	}

	rawResponseData, err := json.Marshal(expectedResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	response, err := s.service.SyncDo(s.requestID, s.request)
	s.Require().NoError(err)
	// Dereference the pointer to access the slice
	resultSlice := *response.Result
	s.Equal(s.symbol, resultSlice[0].Symbol)
	s.Equal(s.orderId, resultSlice[0].OrderID)
}

func (s *orderStatusServiceWsTestSuite) TestSyncDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().
		WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)

	response, err := s.service.SyncDo("", s.request)
	s.Nil(response)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *orderStatusServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.service = &OrderStatusWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}
