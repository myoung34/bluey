package emitters

import (
	"github.com/jarcoal/httpmock"
	"github.com/myoung34/bluey/bluey"
	"github.com/stretchr/testify/assert"
	"testing"
)

type BlueStdoutTest struct {
	Type          string
	Payload       blue.Payload
	DeviceUUIDMap string
	Enabled       bool
	Template      string
}

type convStdoutTest struct {
	name string
	in   BlueStdoutTest
	out  BlueTest
}

func TestStdout(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	blue.EnableLogging()

	theTests := []convStdoutTest{
		{
			name: "Tilt",
			in: BlueStdoutTest{
				Type: "stdout",
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
				Template:      "{\"color\": \"{{.Nickname}}\", \"gravity\": {{.Minor}}, \"mac\": \"{{.Mac}}\", \"temp\": {{.Major}}, \"timestamp\": \"{{.Timestamp}}\"}",
			},
			out: BlueTest{
				Response:      "{\"color\": \"RED\", \"gravity\": 1024, \"mac\": \"11:22:33:44:55\", \"temp\": 90, \"timestamp\": \"1661445284\"}",
				CallCount:     1,
				CallSignature: "",
			},
		},
		{
			name: "Puck",
			in: BlueStdoutTest{
				Type: "stdout",
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
				Template:      "{\"device\": \"{{.Nickname}}\", \"minor\": {{.Minor}}, \"mac\": \"{{.Mac}}\", \"major\": {{.Major}}, \"timestamp\": \"{{.Timestamp}}\"}",
			},
			out: BlueTest{
				Response:      "{\"device\": \"puck\", \"minor\": 1, \"mac\": \"66:77:88:99:00\", \"major\": 1, \"timestamp\": \"1661445284\"}",
				CallCount:     0,
				CallSignature: "",
			},
		},
		{
			name: "PuckDefaultTemplate",
			in: BlueStdoutTest{
				Type: "stdout",
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
			},
			out: BlueTest{
				Response:      "{\"name\": \"Device puck\", \"Minor\": 1, \"Major\": 1}",
				CallCount:     0,
				CallSignature: "",
			},
		},
		{
			name: "PuckNoNickname",
			in: BlueStdoutTest{
				Type: "stdout",
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
			},
			out: BlueTest{
				Response:      "{\"name\": \"Device 61020304620663030064086507660367\", \"Minor\": 1, \"Major\": 1}",
				CallCount:     0,
				CallSignature: "",
			},
		},
	}
	for _, theT := range theTests {
		t.Run(theT.name, func(t *testing.T) {
			sampleConfig := blue.ParseConfig("some/file/somewhere.toml")

			sampleConfig.ConfigData.Set("stdout.device_uuids", theT.in.DeviceUUIDMap)
			sampleConfig.ConfigData.Set("stdout.template", theT.in.Template)
			resp, err := StdoutEmit(theT.in.Payload, sampleConfig.ConfigData.Get(theT.in.Type))
			assert.Equal(t, nil, err)
			assert.Equal(t, theT.out.Response, resp)
		})
	}
}
