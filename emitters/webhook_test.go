package emitters

import (
	"bytes"
	"github.com/jarcoal/httpmock"
	"github.com/myoung34/bluey/bluey"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type BlueWebhookTest struct {
	Type          string
	Payload       blue.Payload
	DeviceUUIDMap string
	Enabled       bool
	URL           string
	Headers       string
	Template      string
	Method        string
}

type BlueTest struct {
	Response      string
	CallCount     int
	CallSignature string
}

type convTest struct {
	name string
	in   BlueWebhookTest
	out  BlueTest
}

func TestWebhook(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	blue.EnableLogging()

	theTests := []convTest{
		{
			name: "POST",
			in: BlueWebhookTest{
				Type: "webhook",
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
				URL:           "http://something.com",
				Headers:       "{\"Content-Type\": \"application/json\", \"Foo\": \"bar\"}",
				Template:      "{\"color\": \"{{.Nickname}}\", \"gravity\": {{.Minor}}, \"mac\": \"{{.Mac}}\", \"temp\": {{.Major}}, \"timestamp\": \"{{.Timestamp}}\"}",
				Method:        "POST",
			},
			out: BlueTest{
				Response:      "{\"Response\":\"{\\\"color\\\": \\\"RED\\\", \\\"gravity\\\": 1024, \\\"mac\\\": \\\"11:22:33:44:55\\\", \\\"temp\\\": 90, \\\"timestamp\\\": \\\"1661445284\\\"}\"}",
				CallCount:     1,
				CallSignature: "POST http://something.com",
			},
		},
		{
			name: "GET",
			in: BlueWebhookTest{
				Type: "webhook",
				Payload: blue.Payload{
					ID:        "61020304620663030064086507660367",
					Mac:       "66:77:88:99:00",
					Nickname:  "Puck",
					Major:     1,
					Minor:     1,
					Rssi:      -7,
					Timestamp: 1661445284,
				},
				DeviceUUIDMap: "61020304620663030064086507660367=puck",
				Enabled:       true,
				URL:           "http://fake.com",
				Headers:       "{\"Content-Type\": \"application/json\"}",
				Template:      "{\"device\": \"{{.Nickname}}\", \"minor\": {{.Minor}}, \"mac\": \"{{.Mac}}\", \"major\": {{.Major}}, \"timestamp\": \"{{.Timestamp}}\"}",
				Method:        "GET",
			},
			out: BlueTest{
				Response:      "{\"Response\":\"{\\\"device\\\": \\\"puck\\\", \\\"minor\\\": 1, \\\"mac\\\": \\\"66:77:88:99:00\\\", \\\"major\\\": 1, \\\"timestamp\\\": \\\"1661445284\\\"}\"}",
				CallCount:     2,
				CallSignature: "GET http://fake.com",
			},
		},
		{
			name: "GET",
			in: BlueWebhookTest{
				Type: "webhook",
				Payload: blue.Payload{
					ID:        "61020304620663030064086507660367",
					Mac:       "66:77:88:99:00",
					Nickname:  "Puck",
					Major:     1,
					Minor:     1,
					Rssi:      -7,
					Timestamp: 1661445284,
				},
				DeviceUUIDMap: "61020304620663030064086507660367=puck",
				Enabled:       true,
				URL:           "http://fake.com",
				Headers:       "{\"Content-Type\": \"application/json\"}",
				Method:        "GET",
			},
			out: BlueTest{
				Response:      "{\"Response\":\"{\\\"name\\\": \\\"Device puck\\\", \\\"Minor\\\": 1, \\\"Major\\\": 1}\"}",
				CallCount:     3,
				CallSignature: "GET http://fake.com",
			},
		},
		{
			name: "GET",
			in: BlueWebhookTest{
				Type: "webhook",
				Payload: blue.Payload{
					ID:        "61020304620663030064086507660367",
					Mac:       "66:77:88:99:00",
					Nickname:  "Puck",
					Major:     1,
					Minor:     1,
					Rssi:      -7,
					Timestamp: 1661445284,
				},
				DeviceUUIDMap: "1234=notfound",
				Enabled:       true,
				URL:           "http://fake.com",
				Headers:       "{\"Content-Type\": \"application/json\"}",
				Method:        "GET",
			},
			out: BlueTest{
				Response:      "{\"Response\":\"{\\\"name\\\": \\\"Device 61020304620663030064086507660367\\\", \\\"Minor\\\": 1, \\\"Major\\\": 1}\"}",
				CallCount:     4,
				CallSignature: "GET http://fake.com",
			},
		},
	}
	for _, theT := range theTests {

		httpmock.RegisterResponder(theT.in.Method, theT.in.URL,
			func(req *http.Request) (*http.Response, error) {
				buf := new(bytes.Buffer)
				buf.ReadFrom(req.Body)
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"Response": buf.String(),
				})
			},
		)
		t.Run(theT.name, func(t *testing.T) {
			sampleConfig := blue.ParseConfig("some/file/somewhere.toml")

			sampleConfig.ConfigData.Set("webhook.url", theT.in.URL)
			sampleConfig.ConfigData.Set("webhook.headers", theT.in.Headers)
			sampleConfig.ConfigData.Set("webhook.device_uuids", theT.in.DeviceUUIDMap)
			sampleConfig.ConfigData.Set("webhook.template", theT.in.Template)
			sampleConfig.ConfigData.Set("webhook.method", theT.in.Method)
			resp, err := WebhookEmit(theT.in.Payload, sampleConfig.ConfigData.Get(theT.in.Type))
			assert.Equal(t, nil, err)
			assert.Equal(t, theT.out.Response, resp)
			assert.Equal(t, theT.out.CallCount, httpmock.GetTotalCallCount())
			info := httpmock.GetCallCountInfo()
			assert.Equal(t, 1, info[theT.out.CallSignature])
		})
	}
}
