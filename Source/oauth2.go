package discordgoself

type MembershipState int

const (
	MembershipStateInvited  MembershipState = 1
	MembershipStateAccepted MembershipState = 2
)

type TeamMember struct {
	User            *User           `json:"user"`
	TeamID          string          `json:"team_id"`
	MembershipState MembershipState `json:"membership_state"`
	Permissions     []string        `json:"permissions"`
}

type Team struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	OwnerID     string        `json:"owner_user_id"`
	Members     []*TeamMember `json:"members"`
}

type Application struct {
	ID                  string    `json:"id,omitempty"`
	Name                string    `json:"name"`
	Description         string    `json:"description,omitempty"`
	Icon                string    `json:"icon,omitempty"`
	Secret              string    `json:"secret,omitempty"`
	RedirectURIs        *[]string `json:"redirect_uris,omitempty"`
	BotRequireCodeGrant bool      `json:"bot_require_code_grant,omitempty"`
	BotPublic           bool      `json:"bot_public,omitempty"`
	RPCApplicationState int       `json:"rpc_application_state,omitempty"`
	Flags               int       `json:"flags,omitempty"`
	Owner               *User     `json:"owner"`
	Bot                 *User     `json:"bot"`
	Team                *Team     `json:"team"`
}

func (s *Session) Application(appID string) (st *Application, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointOAuth2Application(appID), nil, EndpointOAuth2Application(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) Applications() (st []*Application, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointOAuth2Applications, nil, EndpointOAuth2Applications)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ApplicationCreate(ap *Application) (st *Application, err error) {

	data := struct {
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		RedirectURIs *[]string `json:"redirect_uris,omitempty"`
	}{ap.Name, ap.Description, ap.RedirectURIs}

	body, err := s.RequestWithBucketID("POST", EndpointOAuth2Applications, data, EndpointOAuth2Applications)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ApplicationUpdate(appID string, ap *Application) (st *Application, err error) {

	data := struct {
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		RedirectURIs *[]string `json:"redirect_uris,omitempty"`
	}{ap.Name, ap.Description, ap.RedirectURIs}

	body, err := s.RequestWithBucketID("PUT", EndpointOAuth2Application(appID), data, EndpointOAuth2Application(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ApplicationDelete(appID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointOAuth2Application(appID), nil, EndpointOAuth2Application(""))
	if err != nil {
		return
	}

	return
}

type Asset struct {
	Type int    `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Session) ApplicationAssets(appID string) (ass []*Asset, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointOAuth2ApplicationAssets(appID), nil, EndpointOAuth2ApplicationAssets(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &ass)
	return
}

func (s *Session) ApplicationBotCreate(appID string) (st *User, err error) {

	body, err := s.RequestWithBucketID("POST", EndpointOAuth2ApplicationsBot(appID), nil, EndpointOAuth2ApplicationsBot(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}
