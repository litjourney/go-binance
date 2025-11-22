package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/adshao/go-binance/v2/common"
)

// WsApiMethodType define method name for websocket API
type WsApiMethodType string

type WsExchangeInfoSymbolStatusType string

type WsExchangeInfoPermissionType string

// WsApiRequest define common websocket API request
type WsApiRequest struct {
	Id     string                 `json:"id"`
	Method WsApiMethodType        `json:"method"`
	Params map[string]interface{} `json:"params"`
}

// WriteSyncWsTimeout defines timeout for WriteSync method of client_ws
var WriteSyncWsTimeout = 5 * time.Second

const (
	// apiKey define key for websocket API parameters
	apiKey = "apiKey"

	// timestampKey define key for websocket API parameters
	timestampKey = "timestamp"

	// signatureKey define key for websocket API parameters
	signatureKey = "signature"

	// SPOT

	// UserDataStreamSubscribeSignatureSpotWsApiMethod define method for user data stream subscription via websocket API with signature
	UserDataStreamSubscribeSignatureSpotWsApiMethod WsApiMethodType = "userDataStream.subscribe.signature"

	// OrderPlaceSpotWsApiMethod define method for creation order via websocket API
	OrderPlaceSpotWsApiMethod WsApiMethodType = "order.place"

	// OrderCancelSpotWsApiMethod define method for cancel order via websocket API
	OrderCancelSpotWsApiMethod WsApiMethodType = "order.cancel"

	// OrderStatusSpotWsApiMethod define method for query order via websocket API
	OrderStatusSpotWsApiMethod WsApiMethodType = "order.status"

	TimeCheckWsApiMehod WsApiMethodType = "time"

	ExchangeInfoWsApiMehod WsApiMethodType = "exchangeInfo"

	AccountStatusWsApiMethod WsApiMethodType = "account.status"

	// OrderOpenStatusSpotWsApiMethod define method for query open orders via websocket API
	OrderOpenStatusSpotWsApiMethod WsApiMethodType = "openOrders.status"

	AllOrdersSpotWsApiMethod WsApiMethodType = "allOrders"

	// OrderListPlaceOcoSpotWsApiMethod define method for creation OCO order list via websocket API
	OrderListPlaceOcoSpotWsApiMethod WsApiMethodType = "orderList.place.oco"

	// OrderListPlaceSpotWsApiMethod define method for creation order list (deprecated OCO) via websocket API
	OrderListPlaceSpotWsApiMethod WsApiMethodType = "orderList.place"

	// OrderListPlaceOtoSpotWsApiMethod define method for creation OTO order list via websocket API
	OrderListPlaceOtoSpotWsApiMethod WsApiMethodType = "orderList.place.oto"

	// OrderListPlaceOtocoSpotWsApiMethod define method for creation OTOCO order list via websocket API
	OrderListPlaceOtocoSpotWsApiMethod WsApiMethodType = "orderList.place.otoco"

	// OrderListCancelSpotWsApiMethod define method for canceling order list via websocket API
	OrderListCancelSpotWsApiMethod WsApiMethodType = "orderList.cancel"

	// SorOrderPlaceSpotWsApiMethod define method for SOR order placement via websocket API
	SorOrderPlaceSpotWsApiMethod WsApiMethodType = "sor.order.place"

	// SorOrderTestSpotWsApiMethod define method for SOR order testing via websocket API
	SorOrderTestSpotWsApiMethod WsApiMethodType = "sor.order.test"

	// FUTURES

	// OrderPlaceFuturesWsApiMethod define method for creation order via websocket API
	OrderPlaceFuturesWsApiMethod WsApiMethodType = "order.place"

	// CancelFuturesWsApiMethod define method for cancel order via websocket API
	CancelFuturesWsApiMethod WsApiMethodType = "order.cancel"

	// OrderStatusFuturesWsApiMethod define method for query order via websocket API
	OrderStatusFuturesWsApiMethod WsApiMethodType = "order.status"

	// ExchangeInfo Permission Types
	PermissionTypeSPOT      WsExchangeInfoPermissionType = "SPOT"
	PermissionTypeMARGIN    WsExchangeInfoPermissionType = "MARGIN"
	PermissionTypeLEVERAGED WsExchangeInfoPermissionType = "LEVERAGED"
	PermissionTypeTRDGRP002 WsExchangeInfoPermissionType = "TRD_GRP_002"
	PermissionTypeTRDGRP003 WsExchangeInfoPermissionType = "TRD_GRP_003"
	PermissionTypeTRDGRP004 WsExchangeInfoPermissionType = "TRD_GRP_004"
	PermissionTypeTRDGRP005 WsExchangeInfoPermissionType = "TRD_GRP_005"
	PermissionTypeTRDGRP006 WsExchangeInfoPermissionType = "TRD_GRP_006"
	PermissionTypeTRDGRP007 WsExchangeInfoPermissionType = "TRD_GRP_007"
	PermissionTypeTRDGRP008 WsExchangeInfoPermissionType = "TRD_GRP_008"
	PermissionTypeTRDGRP009 WsExchangeInfoPermissionType = "TRD_GRP_009"
	PermissionTypeTRDGRP010 WsExchangeInfoPermissionType = "TRD_GRP_010"
	PermissionTypeTRDGRP011 WsExchangeInfoPermissionType = "TRD_GRP_011"
	PermissionTypeTRDGRP012 WsExchangeInfoPermissionType = "TRD_GRP_012"
	PermissionTypeTRDGRP013 WsExchangeInfoPermissionType = "TRD_GRP_013"
	PermissionTypeTRDGRP014 WsExchangeInfoPermissionType = "TRD_GRP_014"
	PermissionTypeTRDGRP015 WsExchangeInfoPermissionType = "TRD_GRP_015"
	PermissionTypeTRDGRP016 WsExchangeInfoPermissionType = "TRD_GRP_016"
	PermissionTypeTRDGRP017 WsExchangeInfoPermissionType = "TRD_GRP_017"
	PermissionTypeTRDGRP018 WsExchangeInfoPermissionType = "TRD_GRP_018"
	PermissionTypeTRDGRP019 WsExchangeInfoPermissionType = "TRD_GRP_019"
	PermissionTypeTRDGRP020 WsExchangeInfoPermissionType = "TRD_GRP_020"
	PermissionTypeTRDGRP021 WsExchangeInfoPermissionType = "TRD_GRP_021"
	PermissionTypeTRDGRP022 WsExchangeInfoPermissionType = "TRD_GRP_022"
	PermissionTypeTRDGRP023 WsExchangeInfoPermissionType = "TRD_GRP_023"
	PermissionTypeTRDGRP024 WsExchangeInfoPermissionType = "TRD_GRP_024"
	PermissionTypeTRDGRP025 WsExchangeInfoPermissionType = "TRD_GRP_025"

	SymbolStatusTrading  WsExchangeInfoSymbolStatusType = "TRADING"
	SymbolStatusEndOfDay WsExchangeInfoSymbolStatusType = "END_OF_DAY"
	SymbolStatusHalt     WsExchangeInfoSymbolStatusType = "HALT"
	SymbolStatusBreak    WsExchangeInfoSymbolStatusType = "BREAK"
)

var (
	// ErrorRequestIDNotSet defines that request ID is not set
	ErrorRequestIDNotSet = errors.New("ws service: request id is not set")

	// ErrorApiKeyIsNotSet defines that ApiKey is not set
	ErrorApiKeyIsNotSet = errors.New("ws service: api key is not set")

	// ErrorSecretKeyIsNotSet defines that SecretKey is not set
	ErrorSecretKeyIsNotSet = errors.New("ws service: secret key is not set")
)

func NewRequestData(
	requestID string,
	apiKey string,
	secretKey string,
	timeOffset int64,
	keyType string,
) RequestData {
	return RequestData{
		requestID:  requestID,
		apiKey:     apiKey,
		secretKey:  secretKey,
		timeOffset: timeOffset,
		keyType:    keyType,
	}
}

type RequestData struct {
	requestID  string
	apiKey     string
	secretKey  string
	timeOffset int64
	keyType    string
}

// CreateRequest creates signed ws request
func CreateRequest(reqData RequestData, method WsApiMethodType, params map[string]interface{}) ([]byte, error) {
	if reqData.requestID == "" {
		return nil, ErrorRequestIDNotSet
	}

	if reqData.apiKey == "" {
		return nil, ErrorApiKeyIsNotSet
	}

	if reqData.secretKey == "" {
		return nil, ErrorSecretKeyIsNotSet
	}

	params[apiKey] = reqData.apiKey
	params[timestampKey] = timestamp(reqData.timeOffset)

	sf, err := common.SignFunc(reqData.keyType)
	if err != nil {
		return nil, err
	}
	signature, err := sf(reqData.secretKey, encodeParams(params))
	if err != nil {
		return nil, err
	}
	params[signatureKey] = signature

	req := WsApiRequest{
		Id:     reqData.requestID,
		Method: method,
		Params: params,
	}

	rawData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

// CreateRequest creates  ws request
func CreateRequestWithSigned(requestID string, method WsApiMethodType, params map[string]interface{}) ([]byte, error) {
	req := WsApiRequest{
		Id:     requestID,
		Method: method,
		Params: params,
	}

	rawData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

// encode encodes the parameters to a URL encoded string
func encodeParams(p map[string]interface{}) string {
	queryValues := url.Values{}
	for key, value := range p {
		queryValues.Add(key, fmt.Sprintf("%v", value))
	}
	return queryValues.Encode()
}

func timestamp(offsetMilli int64) int64 {
	return time.Now().UnixMilli() - offsetMilli
}
