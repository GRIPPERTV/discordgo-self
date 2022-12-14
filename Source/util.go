package discordgoself

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strconv"
	"time"
)

func SnowflakeTimestamp(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(0, timestamp*1000000)
	return
}

func MultipartBodyWithJSON(data interface{}, files []*File) (requestContentType string, requestBody []byte, err error) {
	body := &bytes.Buffer{}
	bodywriter := multipart.NewWriter(body)

	payload, err := json.Marshal(data)
	if err != nil {
		return
	}

	var p io.Writer

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", "application/json")

	p, err = bodywriter.CreatePart(h)
	if err != nil {
		return
	}

	if _, err = p.Write(payload); err != nil {
		return
	}

	for i, file := range files {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file%d"; filename="%s"`, i, quoteEscaper.Replace(file.Name)))
		contentType := file.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		h.Set("Content-Type", contentType)

		p, err = bodywriter.CreatePart(h)
		if err != nil {
			return
		}

		if _, err = io.Copy(p, file.Reader); err != nil {
			return
		}
	}

	err = bodywriter.Close()
	if err != nil {
		return
	}

	return bodywriter.FormDataContentType(), body.Bytes(), nil
}
