package discordgoself

import (
	"errors"
	"fmt"
	"runtime"
	"net/http"
	"time"
)

const VERSION = "0.1.0"

var ErrMFA = errors.New("account has 2FA enabled")

func New(args ...interface{}) (s *Session, err error) {
	if args == nil {
		return
	}

	s = &Session{
		State:                  NewState(),
		Ratelimiter:            NewRatelimiter(),
		StateEnabled:           true,
		ShouldReconnectOnError: true,
		MaxRestRetries:         3,
		Client:                 &http.Client{Timeout: (20 * time.Second)},
		UserAgent:              "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
		sequence:               new(int64),
		LastHeartbeatAck:       time.Now().UTC(),
	}
	
	s.Identify.Token = ""
	s.Identify.Compress = false
	s.Identify.LargeThreshold = 250
	s.Identify.Properties = IdentifyProperties{
		OS:              runtime.GOOS,
		Browser:         "DiscordGo-self v" + VERSION,
		Device:	         "",
		Referrer:        "",
		ReferringDomain: "",
	}

	var auth, pass string

	for _, arg := range args {

		switch v := arg.(type) {

		case []string:
			if len(v) > 3 {
				err = fmt.Errorf("too many string parameters provided")
				return
			}

			if len(v) > 0 {
				auth = v[0]
			}

			if len(v) > 1 {
				pass = v[1]
			}

			if len(v) > 2 {
				s.Identify.Token = v[2]
			}

		case string:
			if auth == "" {
				auth = v
			} else if pass == "" {
				pass = v
			} else if s.Identify.Token == "" {
				s.Identify.Token = v
			} else {
				err = fmt.Errorf("too many string parameters provided")
				return
			}
		default:
			err = fmt.Errorf("unsupported parameter type provided")
			return
		}
	}

	if pass == "" {
		s.Identify.Token = auth
	} else {
		err = s.Login(auth, pass)

		if err != nil || s.Identify.Token == "" {
			if s.MFA {
				err = ErrMFA
			} else {
				err = fmt.Errorf("Unable to fetch discord authentication token. %v", err)
			}
			return
		}
	}

	return
}
