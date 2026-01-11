package errors

import "net/http"

type SendableHTTPError interface {
	StatusCode() int
	Error() string
}

type TimeoutError struct {
	Message string
}

func (e *TimeoutError) Error() string   { return e.Message }
func (e *TimeoutError) StatusCode() int { return http.StatusGatewayTimeout }

type TunnelNotRespondingError struct{}

func (e *TunnelNotRespondingError) Error() string {
	return "tunnel did not respond"
}
func (e *TunnelNotRespondingError) StatusCode() int { return http.StatusBadGateway }
