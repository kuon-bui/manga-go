package common

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SSEStream struct {
	writer gin.ResponseWriter
}

func NewSSEStream(c *gin.Context) *SSEStream {
	header := c.Writer.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	header.Set("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	return &SSEStream{writer: c.Writer}
}

func (s *SSEStream) SendComment(comment string) error {
	normalizedComment := strings.ReplaceAll(strings.ReplaceAll(comment, "\r\n", " "), "\n", " ")
	if _, err := fmt.Fprintf(s.writer, ": %s\n\n", normalizedComment); err != nil {
		return err
	}

	s.writer.Flush()

	return nil
}

func (s *SSEStream) SendEvent(eventName string, data string) error {
	if eventName != "" {
		if _, err := fmt.Fprintf(s.writer, "event: %s\n", eventName); err != nil {
			return err
		}
	}

	normalizedData := strings.ReplaceAll(data, "\r\n", "\n")
	normalizedData = strings.ReplaceAll(normalizedData, "\r", "\n")
	for line := range strings.SplitSeq(normalizedData, "\n") {
		if _, err := fmt.Fprintf(s.writer, "data: %s\n", line); err != nil {
			return err
		}
	}

	if _, err := s.writer.Write([]byte("\n")); err != nil {
		return err
	}

	s.writer.Flush()

	return nil
}

func (s *SSEStream) SendHeartbeat() error {
	return s.SendEvent("heartbeat", "{}")
}
