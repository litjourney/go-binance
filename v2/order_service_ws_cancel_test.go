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

type orderCancelServiceWsTestSuite struct {
	suite.Suite
	apiKey     string
	secretKey  string
	signedKey  string
	timeOffset int64

	ctrl   *gomock.Controller
	client *mock.MockClient

	requestID          string
	symbol             string
	orderId            int64
	origClientOrderId  string
	newClientOrderId   string
	cancelRestrictions OrderCancelRestrictionsType
	recvWindow         uint16

	service *OrderCancelWsService
	request *OrderCancelWsRequest
}

func (s *orderCancelServiceWsTestSuite) SetupTest() {
	s.apiKey = "dummyApiKey"
	s.secretKey = "dummySecretKey"
	s.signedKey = "HMAC"
	s.timeOffset = 0

	s.requestID = "e2a85d9f-07a5-4f94-8d5f-789dc3deb098"
	s.symbol = "BTCUSDT"
	s.orderId = 12345
	s.origClientOrderId = "origClientOrderId"
	s.newClientOrderId = "newClientOrderId"
	s.cancelRestrictions = OrderCancelOnlyNew
	s.recvWindow = 5000

	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)

	s.service = &OrderCancelWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.request = NewOrderCancelWsRequest().
		Symbol(s.symbol).
		OrderID(s.orderId).
		OrigClientOrderID(s.origClientOrderId).
		NewClientOrderID(s.newClientOrderId).
		CancelRestrictions(s.cancelRestrictions).
		RecvWindow(s.recvWindow)
}

func (s *orderCancelServiceWsTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestOrderCancelServiceWs(t *testing.T) {
	suite.Run(t, new(orderCancelServiceWsTestSuite))
}

func (s *orderCancelServiceWsTestSuite) TestDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.service.Do(s.requestID, s.request)
	s.NoError(err)
}

func (s *orderCancelServiceWsTestSuite) TestDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	err := s.service.Do("", s.request)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *orderCancelServiceWsTestSuite) TestDo_EmptyApiKey() {
	s.reset("", s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorApiKeyIsNotSet)
}

func (s *orderCancelServiceWsTestSuite) TestDo_EmptySecretKey() {
	s.reset(s.apiKey, "", s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorSecretKeyIsNotSet)
}

func (s *orderCancelServiceWsTestSuite) TestSyncDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	expectedResponse := CancelOrderWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: CreateOrderResult{
			CreateOrderResponse{
				Symbol:        s.symbol,
				OrderID:       s.orderId,
				ClientOrderID: s.newClientOrderId,
			},
		},
	}

	rawResponseData, err := json.Marshal(expectedResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	response, err := s.service.SyncDo(s.requestID, s.request)
	s.Require().NoError(err)
	s.Equal(s.symbol, response.Result.Symbol)
	s.Equal(s.orderId, response.Result.OrderID)
	s.Equal(s.newClientOrderId, response.Result.ClientOrderID)
}

func (s *orderCancelServiceWsTestSuite) TestSyncDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().
		WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)

	response, err := s.service.SyncDo("", s.request)
	s.Nil(response)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *orderCancelServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.service = &OrderCancelWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}
