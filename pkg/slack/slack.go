package slack

import (
	"github.com/sirupsen/logrus"
	"fmt"
	"bytes"
	"net/http"
)

func PostBestEffort(url string, msg string, log *logrus.Entry) {
	body := fmt.Sprintf(`{"text":"%s"}`, msg)
	buf := bytes.NewReader([]byte(body))
	log.Debug("message to be posted to Slack: %s", msg)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		log.Warnf("failed to post to slack: %s", err)
		return
	}
	status := resp.StatusCode
	if status != 200 {
		log.Warnf("unexpected slack status code: %d", status)
		return
	}
}
