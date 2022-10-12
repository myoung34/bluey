package emitters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/myoung34/bluey/bluey"
	"io"
	"net/http"
	"strconv"
	"text/template"
)

type Webhook struct {
	Enabled       bool
	DeviceUUIDMap string `json:"device_uuids"`
	URL           string
	Headers       string
	Template      string
	Method        string
}

func WebhookEmit(payload blue.Payload, emitterConfig interface{}) (string, error) {
	webhook := Webhook{}
	jsonString, _ := json.Marshal(emitterConfig)
	json.Unmarshal(jsonString, &webhook)

	level.Info(blue.Logger).Log("emitters.webhook.url", webhook.URL)
	level.Info(blue.Logger).Log("emitters.webhook.enabled", fmt.Sprintf("%v", webhook.Enabled))
	level.Info(blue.Logger).Log("emitters.webhook.headers", fmt.Sprintf("%+v", webhook.Headers))
	level.Info(blue.Logger).Log("emitters.webhook.template", fmt.Sprintf("%+v", webhook.Template))
	level.Info(blue.Logger).Log("emitters.webhook.method", fmt.Sprintf("%+v", webhook.Method))

	client := &http.Client{}

	nickname := getNickNameFromUUIDs(webhook.DeviceUUIDMap, payload.ID)

	bodyTemplate := Template{
		Nickname:  nickname,
		Minor:     strconv.Itoa(int(payload.Minor)),
		Mac:       payload.Mac,
		Major:     strconv.Itoa(int(payload.Major)),
		Timestamp: payload.Timestamp,
	}
	level.Info(blue.Logger).Log("emitters.webhook.payload", fmt.Sprintf("%+v", payload))

	tmpl, err := template.New("webhook").Parse(`{"name": "Device {{.Nickname}}", "Minor": {{.Minor}}, "Major": {{.Major}}}`)
	if len(webhook.Template) > 0 {
		tmpl, err = template.New("webhook").Parse(webhook.Template)
	}
	if err != nil {
		level.Error(blue.Logger).Log("emitters.webhook", err)
		return "", err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, bodyTemplate); err != nil {
		level.Error(blue.Logger).Log("emitters.webhook", err)
		return "", err
	}
	bodyReader := bytes.NewReader(tpl.Bytes())

	level.Info(blue.Logger).Log("emitters.webhook.rendered_template", fmt.Sprintf("%+v", tpl.String()))
	// Set up the request
	req, err := http.NewRequest(webhook.Method, webhook.URL, bodyReader)
	if err != nil {
		level.Error(blue.Logger).Log("emitters.webhook", err)
		return "", err
	}

	// Parse the headers and add them
	var result map[string]string
	json.Unmarshal([]byte(webhook.Headers), &result)
	for key, value := range result {
		req.Header.Add(key, value)
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		level.Error(blue.Logger).Log("emitters.webhook", err)
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		level.Error(blue.Logger).Log("emitters.webhook", err)
		return "", err
	}

	return string(respBody), nil
}
