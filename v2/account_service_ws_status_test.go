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

type accountStatusServiceWsTestSuite struct {
	suite.Suite
	apiKey     string
	secretKey  string
	signedKey  string
	timeOffset int64

	ctrl   *gomock.Controller
	client *mock.MockClient

	requestID        string
	omitZeroBalances bool
	recvWindow       uint16

	service *AccountStatusWsService
	request *AccountStatusWsRequest
}

func (s *accountStatusServiceWsTestSuite) SetupTest() {
	s.apiKey = "dummyApiKey"
	s.secretKey = "dummySecretKey"
	s.signedKey = "HMAC"
	s.timeOffset = 0

	s.requestID = "e2a85d9f-07a5-4f94-8d5f-789dc3deb098"
	s.omitZeroBalances = true
	s.recvWindow = 5000

	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)

	s.service = &AccountStatusWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.request = NewAccountStatusWsRequest().
		OmitZeroBalances(s.omitZeroBalances).
		RecvWindow(s.recvWindow)
}

func (s *accountStatusServiceWsTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestAccountStatusServiceWs(t *testing.T) {
	suite.Run(t, new(accountStatusServiceWsTestSuite))
}

func (s *accountStatusServiceWsTestSuite) TestDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.service.Do(s.requestID, s.request)
	s.NoError(err)
}

func (s *accountStatusServiceWsTestSuite) TestDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	err := s.service.Do("", s.request)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *accountStatusServiceWsTestSuite) TestDo_EmptyApiKey() {
	s.reset("", s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorApiKeyIsNotSet)
}

func (s *accountStatusServiceWsTestSuite) TestDo_EmptySecretKey() {
	s.reset(s.apiKey, "", s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)

	err := s.service.Do(s.requestID, s.request)
	s.ErrorIs(err, websocket.ErrorSecretKeyIsNotSet)
}

func (s *accountStatusServiceWsTestSuite) TestSyncDo() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	expectedResponse := AccountStatusWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: AccountStatusResult{
			Account: Account{
				MakerCommission:  15,
				TakerCommission:  15,
				BuyerCommission:  0,
				SellerCommission: 0,
				CanTrade:         true,
				CanWithdraw:      true,
				CanDeposit:       true,
				UpdateTime:       123456789,
				AccountType:      "SPOT",
				Balances: []Balance{
					{
						Asset:  "BTC",
						Free:   "4723846.89208129",
						Locked: "0.00000000",
					},
					{
						Asset:  "LTC",
						Free:   "4763368.68006011",
						Locked: "0.00000000",
					},
				},
			},
		},
	}

	rawResponseData, err := json.Marshal(expectedResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	response, err := s.service.SyncDo(s.requestID, s.request)
	s.Require().NoError(err)
	s.Equal(expectedResponse.Result.MakerCommission, response.Result.MakerCommission)
	s.Equal(expectedResponse.Result.TakerCommission, response.Result.TakerCommission)
	s.Equal(expectedResponse.Result.Balances[0].Asset, response.Result.Balances[0].Asset)
}

func (s *accountStatusServiceWsTestSuite) TestSyncDo_EmptyRequestID() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().
		WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)

	response, err := s.service.SyncDo("", s.request)
	s.Nil(response)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *accountStatusServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.service = &AccountStatusWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}
