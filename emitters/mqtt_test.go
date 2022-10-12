package emitters

import (
	_mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jarcoal/httpmock"
	"github.com/myoung34/bluey/bluey"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type BlueMQTTTest struct {
	Type          string
	Payload       blue.Payload
	DeviceUUIDMap string
	Enabled       bool
	Template      string
}

type convMQTTTest struct {
	name string
	in   BlueMQTTTest
	out  BlueTest
}

type MockToken struct {
	err error
}

func (mt MockToken) Done() <-chan struct{} {
	return nil
}

func (mt MockToken) Wait() bool {
	return true
}

func (mt MockToken) WaitTimeout(time.Duration) bool {
	return true
}

func (mt MockToken) Error() error {
	return mt.err
}

func (mc MockClient) AddRoute(topic string, callback _mqtt.MessageHandler) {
}

func (mc MockClient) IsConnected() bool {
	return true
}

func (mc MockClient) IsConnectionOpen() bool {
	return true
}

func (mc MockClient) OptionsReader() _mqtt.ClientOptionsReader {
	return _mqtt.ClientOptionsReader{}
}

func (mc MockClient) Connect() _mqtt.Token {
	return MockToken{}
}

func (mc MockClient) Disconnect(quiesce uint) {
}

func (mc MockClient) Publish(topic string, qos byte, retained bool, payload interface{}) _mqtt.Token {
	return mc.token
}

func (mc MockClient) Subscribe(topic string, qos byte, callback _mqtt.MessageHandler) _mqtt.Token {
	return mc.token
}

func (mc MockClient) SubscribeMultiple(filters map[string]byte, callback _mqtt.MessageHandler) _mqtt.Token {
	return mc.token
}

func (mc MockClient) Unsubscribe(topics ...string) _mqtt.Token {
	return mc.token
}

func TestMQTT(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	blue.EnableLogging()

	theTests := []convMQTTTest{
		{
			name: "Tilt",
			in: BlueMQTTTest{
				Type: "mqtt",
				Payload: blue.Payload{
					ID:        "a495bb10c5b14b44b5121370f02d74de",
					Mac:       "11:22:33:44:55",
					Nickname:  "RED",
					Major:     90,
					Minor:     1024,
					Rssi:      -67,
					Timestamp: 1661445284,
				},
				DeviceUUIDMap: "a495bb30c5b14b44b5121370f02d74de=BLACK,a495bb60c5b14b44b5121370f02d74de=BLUE,a495bb20c5b14b44b5121370f02d74de=GREEN,a495bb50c5b14b44b5121370f02d74de=ORANGE,a495bb80c5b14b44b5121370f02d74de=PINK,a495bb40c5b14b44b5121370f02d74de=PURPLE,a495bb10c5b14b44b5121370f02d74de=RED,a495bb70c5b14b44b5121370f02d74de=YELLOW,a495bb90c5b14b44b5121370f02d74de=TEST,25cc0b60914de76ead903f903bfd5e53=MIGHTY",
				Enabled:       true,
			},
			out: BlueTest{
				Response:      "{\"name\": \"Device RED\", \"Minor\": 1024, \"Major\": 90}",
				CallCount:     1,
				CallSignature: "",
			},
		},
	}
	for _, theT := range theTests {
		t.Run(theT.name, func(t *testing.T) {
			sampleConfig := blue.ParseConfig("some/file/somewhere.toml")

			sampleConfig.ConfigData.Set("mqtt.device_uuids", theT.in.DeviceUUIDMap)
			sampleConfig.ConfigData.Set("mqtt.url", "tcp://localhost:1883")
			resp, err := MQTTEmitWithClient(theT.in.Payload, sampleConfig.ConfigData.Get(theT.in.Type), MockClient{})
			assert.Equal(t, nil, err)
			assert.Equal(t, theT.out.Response, resp)
		})
	}
}
