package emitters

import (
	"bytes"
	"encoding/json"
	"fmt"
	_mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log/level"
	"github.com/myoung34/bluey/bluey"
	"strconv"
	"text/template"
)

type MockClient struct {
	token _mqtt.Token
}

type MQTT struct {
	Enabled       bool
	DeviceUUIDMap string `json:"device_uuids"`
	Template      string
	URL           string
	ClientID      string `json:"client_id""`
	Username      string
	Topic         string
	Retained      bool
}

func MQTTEmit(payload blue.Payload, emitterConfig interface{}) (string, error) {
	mqtt := MQTT{}
	jsonString, _ := json.Marshal(emitterConfig)
	json.Unmarshal(jsonString, &mqtt)
	opts := _mqtt.NewClientOptions().
		AddBroker(mqtt.URL).
		SetClientID(mqtt.ClientID).
		SetUsername(mqtt.Username)

	c := _mqtt.NewClient(opts)

	return MQTTEmitWithClient(payload, emitterConfig, c)
}

func MQTTEmitWithClient(payload blue.Payload, emitterConfig interface{}, client _mqtt.Client) (string, error) {
	mqtt := MQTT{}
	jsonString, _ := json.Marshal(emitterConfig)
	json.Unmarshal(jsonString, &mqtt)

	level.Info(blue.Logger).Log("emitters.mqtt.enabled", fmt.Sprintf("%v", mqtt.Enabled))
	level.Info(blue.Logger).Log("emitters.mqtt.url", mqtt.URL)
	level.Info(blue.Logger).Log("emitters.mqtt.client_id", mqtt.ClientID)
	level.Info(blue.Logger).Log("emitters.mqtt.username", mqtt.Username)
	level.Info(blue.Logger).Log("emitters.mqtt.topic", mqtt.Topic)
	level.Info(blue.Logger).Log("emitters.mqtt.retained", mqtt.Retained)
	level.Info(blue.Logger).Log("emitters.mqtt.template", fmt.Sprintf("%+v", mqtt.Template))

	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		level.Error(blue.Logger).Log("emitters.mqtt", token.Error())
		return "", token.Error()
	}

	nickname := getNickNameFromUUIDs(mqtt.DeviceUUIDMap, payload.ID)
	bodyTemplate := Template{
		Nickname:  nickname,
		Minor:     strconv.Itoa(int(payload.Minor)),
		Mac:       payload.Mac,
		Major:     strconv.Itoa(int(payload.Major)),
		Timestamp: payload.Timestamp,
	}
	level.Info(blue.Logger).Log("emitters.webhook.payload", fmt.Sprintf("%+v", payload))

	tmpl, err := template.New("webhook").Parse(`{"name": "Device {{.Nickname}}", "Minor": {{.Minor}}, "Major": {{.Major}}}`)
	if len(mqtt.Template) > 0 {
		tmpl, err = template.New("mqtt").Parse(mqtt.Template)
	}
	if err != nil {
		level.Error(blue.Logger).Log("emitters.mqtt", err)
		return "", err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, bodyTemplate); err != nil {
		level.Error(blue.Logger).Log("emitters.mqtt", err)
		return "", err
	}

	_ = client.Publish(mqtt.Topic, 1, mqtt.Retained, tpl.Bytes())
	client.Disconnect(250)
	return tpl.String(), nil
}
