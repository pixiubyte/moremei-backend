package dify

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type TextToAudioRequest struct {
	MessageId string `json:"message_id,omitempty"`
	Text      string `json:"text"`
	User      string `json:"user"`
	Stream    bool   `json:"stream"`
}

func (c *Client) TextToAudio(ctx context.Context, req TextToAudioRequest) (data []byte, err error) {
	req.Stream = false
	var r *http.Request
	r, err = c.createPostRequest(ctx, c.buildRequestApi(TextToAudioUrl), req)
	if err != nil {
		return nil, err
	}
	resp, err := c.SetHttpRequest(r).SendRequestStream()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TextToAudio failed, status: %s, code: %d", resp.Status, resp.StatusCode)
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)

	return
}

func (c *Client) TextToAudioStream(ctx context.Context, req TextToAudioRequest) (data chan []byte, err error) {
	req.Stream = true

	req.Text = c.trimThinkingContent(req.Text)
	// #号会影响语音生成
	req.Text = regexp.MustCompile(`#{1,6} `).ReplaceAllString(req.Text, "")

	r, err := c.createPostRequest(ctx, c.buildRequestApi(TextToAudioUrl), req)
	if err != nil {
		return nil, err
	}
	resp, err := c.SetHttpRequest(r).SendRequestStream()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TextToAudio failed, status: %s, code: %d", resp.Status, resp.StatusCode)
	}
	data = make(chan []byte)
	go func() {
		defer resp.Body.Close()
		defer close(data)
		for {
			buf := make([]byte, 1024*100)
			n, err1 := resp.Body.Read(buf)
			if err1 != nil {
				if err1 == io.EOF {
					return
				}
				err = fmt.Errorf("TextToAudio failed, status: %s, code: %d, err: %v", resp.Status, resp.StatusCode, err1)
				return
			}
			data <- buf[:n]
		}
	}()

	return
}
