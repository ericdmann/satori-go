package rtm

import (
	"encoding/json"
	"reflect"
	"strconv"
)

// ConvertToRawMessage allows for any interface to be converted to a common `json.RawMessage` data
// structure to ensure that all data is marshall-able when being published or subscribed.
func (rtm *RTMClient) ConvertToRawMessage(message interface{}) (json.RawMessage, error) {
	switch message.(type) {
	case string:
		rawMessage, err := json.Marshal(message.(string))
		return json.RawMessage(rawMessage), err
	case int, int8, int16, int32, int64:
		return json.RawMessage(strconv.FormatInt(reflect.ValueOf(message).Int(), 10)), nil
	case uint, uint8, uint16, uint32, uint64:
		return json.RawMessage(strconv.FormatUint(reflect.ValueOf(message).Uint(), 10)), nil
	case float32, float64:
		return json.RawMessage(strconv.FormatFloat(reflect.ValueOf(message).Float(), 'f', -1, 64)), nil
	case bool:
		return json.RawMessage(strconv.FormatBool(message.(bool))), nil
	case json.RawMessage:
		return message.(json.RawMessage), nil
	case nil:
		return json.RawMessage(`null`), nil
	default:
		rawMessage, err := json.Marshal(message)
		if err != nil {
			return nil, err
		}
		return json.RawMessage(rawMessage), nil
	}
}
