package binance

import (
	"encoding/json"
	"testing"

	"github.com/adshao/go-binance/v2/common/websocket"
	"github.com/adshao/go-binance/v2/common/websocket/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type exchangeInfoServiceWsTestSuite struct {
	suite.Suite

	ctrl   *gomock.Controller
	client *mock.MockClient

	requestID          string
	symbol             string
	symbols            []string
	permissions        []websocket.WsExchangeInfoPermissionType
	showPermissionSets bool
	symbolStatus       websocket.WsExchangeInfoSymbolStatusType

	service *ExchangeInfoWsService
	request *ExchangeInfoWsRequest
}

func (s *exchangeInfoServiceWsTestSuite) SetupTest() {
	s.requestID = "e2a85d9f-07a5-4f94-8d5f-789dc3deb098"
	s.symbol = "BTCUSDT"
	s.symbols = []string{"BTCUSDT", "ETHUSDT"}
	s.permissions = []websocket.WsExchangeInfoPermissionType{websocket.PermissionTypeSPOT}
	s.showPermissionSets = true
	s.symbolStatus = websocket.SymbolStatusTrading

	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)

	s.service = &ExchangeInfoWsService{
		c: s.client,
	}

	s.request = NewExchangeInfoWsRequest().
		SetSymbol(s.symbol).
		SetSymbols(s.symbols).
		SetPermissions(s.permissions).
		SetShowPermissionSets(s.showPermissionSets).
		SetSymbolStatus(s.symbolStatus)
}

func (s *exchangeInfoServiceWsTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestExchangeInfoServiceWs(t *testing.T) {
	suite.Run(t, new(exchangeInfoServiceWsTestSuite))
}

func (s *exchangeInfoServiceWsTestSuite) TestDo() {
	s.reset()

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.service.Do(s.requestID, s.request)
	s.NoError(err)
}

func (s *exchangeInfoServiceWsTestSuite) TestSyncDo() {
	s.reset()

	expectedResponse := ExchangeInfoWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: ExchangeInfo{
			Timezone:   "UTC",
			ServerTime: 123456789,
			Symbols: []Symbol{
				{
					Symbol: s.symbol,
				},
			},
		},
	}

	rawResponseData, err := json.Marshal(expectedResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	response, err := s.service.SyncDo(s.requestID, s.request)
	s.Require().NoError(err)
	s.Equal(expectedResponse.Result.Timezone, response.Result.Timezone)
	s.Equal(expectedResponse.Result.Symbols[0].Symbol, response.Result.Symbols[0].Symbol)
}

func (s *exchangeInfoServiceWsTestSuite) reset() {
	s.service = &ExchangeInfoWsService{
		c: s.client,
	}
}
