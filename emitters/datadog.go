package emitters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/go-kit/log/level"
	"github.com/myoung34/bluey/bluey"
	"strconv"
	"strings"
	"text/template"
)

type Datadog struct {
	Enabled       bool
	DeviceUUIDMap string `json:"device_uuids"`
	StatsdHost    string `json:"statsd_host"`
	StatsdPort    int    `json:"statsd_port"`
	Tags          string `json:"tags"`
	Gauges        string `json:"gauges"`
}

func DatadogEmitWithClient(payload blue.Payload, emitterConfig interface{}, client statsd.ClientInterface) (string, error) {

	defer client.Close()

	datadog := Datadog{}
	jsonString, _ := json.Marshal(emitterConfig)
	json.Unmarshal(jsonString, &datadog)

	nickname := getNickNameFromUUIDs(datadog.DeviceUUIDMap, payload.ID)

	ddTags := make([]string, 0)
	tags := strings.Split(datadog.Tags, ",")
	for _, tag := range tags {
		tmpl, err := template.New("datadog").Parse(tag)
		if err != nil {
			level.Error(blue.Logger).Log("emitters.datadog", err)
			return "", err
		}
		var tpl bytes.Buffer
		tagTemplate := Template{
			Nickname:  nickname,
			Minor:     strconv.Itoa(int(payload.Minor)),
			Mac:       payload.Mac,
			Major:     strconv.Itoa(int(payload.Major)),
			Timestamp: payload.Timestamp,
		}
		if err := tmpl.Execute(&tpl, tagTemplate); err != nil {
			level.Error(blue.Logger).Log("emitters.datadog", err)
			return "", err
		}
		ddTags = append(ddTags, tpl.String())
	}

	level.Debug(blue.Logger).Log("emitters.datadog.tags", fmt.Sprintf("%+v", ddTags))

	gauges := strings.Split(datadog.Gauges, ",")
	for _, gauge := range gauges {
		_gauge := strings.Split(gauge, "=")

		tmpl, err := template.New("datadog").Parse(_gauge[1])
		if err != nil {
			level.Error(blue.Logger).Log("emitters.datadog.gauge", err)
			return "", err
		}
		var tpl bytes.Buffer
		tagTemplate := Template{
			Nickname:  nickname,
			Minor:     strconv.Itoa(int(payload.Minor)),
			Mac:       payload.Mac,
			Major:     strconv.Itoa(int(payload.Major)),
			Timestamp: payload.Timestamp,
		}
		if err := tmpl.Execute(&tpl, tagTemplate); err != nil {
			level.Error(blue.Logger).Log("emitters.datadog", err)
			return "", err
		}

		gaugeVal, err := strconv.Atoi(tpl.String())
		if err != nil {
			level.Error(blue.Logger).Log("emitters.datadog.gauge", err)
			return "", err
		}
		level.Debug(blue.Logger).Log("emitters.datadog.gauge", fmt.Sprintf("%s=%d tags:%+v", _gauge[0], gaugeVal, ddTags))
		client.Gauge(_gauge[0],
			float64(gaugeVal),
			tags,
			1,
		)
	}

	return "", nil
}

func DatadogEmit(payload blue.Payload, emitterConfig interface{}) (string, error) {
	datadog := Datadog{}
	jsonString, _ := json.Marshal(emitterConfig)
	json.Unmarshal(jsonString, &datadog)
	client, err := statsd.New(fmt.Sprintf("%s:%d", datadog.StatsdHost, datadog.StatsdPort))
	if err != nil {
		level.Error(blue.Logger).Log("emitters.datadog", err)
		return "", err
	}
	return DatadogEmitWithClient(payload, emitterConfig, client)
}
