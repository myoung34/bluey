package emitters

import (
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/myoung34/bluey/bluey"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatadogEmit(t *testing.T) {
	blue.EnableLogging()

	sampleConfig := blue.ParseConfig("some/file/somewhere.toml")
	sampleConfig.ConfigData.Set("datadog.enabled", true)
	sampleConfig.ConfigData.Set("datadog.statsd_host", "testing")
	sampleConfig.ConfigData.Set("datadog.statsd_port", "8125")
	sampleConfig.ConfigData.Set("datadog.tags", "color:{{.Nickname}},mac:{{.Mac}}")
	sampleConfig.ConfigData.Set("datadog.gauges", "bluey.temperature={{.Major}},bluey.gravity={{.Minor}}")

	payload := blue.Payload{
		ID:        "61020304620663030064086507660367",
		Mac:       "66:77:88:99:00",
		Nickname:  "BLACK",
		Major:     65,
		Minor:     1098,
		Rssi:      -7,
		Timestamp: 1661445284,
	}
	resp, err := DatadogEmitWithClient(payload, sampleConfig.ConfigData.Get("datadog"), &statsd.NoOpClient{})
	assert.Equal(t, nil, err)
	assert.Equal(t, "", resp)

	resp, err = DatadogEmit(payload, sampleConfig.ConfigData.Get("datadog"))
	assert.NotNil(t, err, "This should have failed DNS lookup")
	assert.Equal(t, resp, "")
}
