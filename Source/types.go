package discordgoself

import (
	"encoding/json"
	"net/http"
	"time"
)

type Timestamp string

func (t Timestamp) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(t))
}

type RESTError struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody []byte
	Message      *APIErrorMessage
}

func newRestError(req *http.Request, resp *http.Response, body []byte) *RESTError {
	restErr := &RESTError{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
	}

	var msg *APIErrorMessage
	err := json.Unmarshal(body, &msg)
	if err == nil {
		restErr.Message = msg
	}

	return restErr
}

func (r RESTError) Error() string {
	return "HTTP " + r.Response.Status + ", " + string(r.ResponseBody)
}
