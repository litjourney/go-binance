package binance

import (
	"encoding/json"
	"testing"

	"github.com/adshao/go-binance/v2/common/websocket/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type timeCheckServiceWsTestSuite struct {
	suite.Suite

	ctrl   *gomock.Controller
	client *mock.MockClient

	requestID string

	service *TimeCheckWsService
	request *TimeCheckWsRequest
}

func (s *timeCheckServiceWsTestSuite) SetupTest() {
	s.requestID = "e2a85d9f-07a5-4f94-8d5f-789dc3deb098"

	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)

	s.service = &TimeCheckWsService{
		c: s.client,
	}

	s.request = NewTimeCheckWsRequest()
}

func (s *timeCheckServiceWsTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestTimeCheckServiceWs(t *testing.T) {
	suite.Run(t, new(timeCheckServiceWsTestSuite))
}

func (s *timeCheckServiceWsTestSuite) TestDo() {
	s.reset()

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.service.Do(s.requestID, s.request)
	s.NoError(err)
}

func (s *timeCheckServiceWsTestSuite) TestSyncDo() {
	s.reset()

	expectedResponse := TimeCheckWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: TimeCheckResult{
			ServerTime: 123456789,
		},
	}

	rawResponseData, err := json.Marshal(expectedResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	response, err := s.service.SyncDo(s.requestID, s.request)
	s.Require().NoError(err)
	s.Equal(expectedResponse.Result.ServerTime, response.Result.ServerTime)
}

func (s *timeCheckServiceWsTestSuite) reset() {
	s.service = &TimeCheckWsService{
		c: s.client,
	}
}
