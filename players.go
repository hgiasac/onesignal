package onesignal

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// PlayersService handles communication with the player related
// methods of the OneSignal API.
type PlayersService struct {
	client *Client
}

// Player represents a OneSignal player.
type Player struct {
	ID                string            `json:"id"`
	Playtime          int               `json:"playtime"`
	SDK               string            `json:"sdk"`
	Identifier        string            `json:"identifier"`
	SessionCount      int               `json:"session_count"`
	Language          string            `json:"language"`
	Timezone          int               `json:"timezone"`
	GameVersion       string            `json:"game_version"`
	DeviceOS          string            `json:"device_os"`
	DeviceType        int               `json:"device_type"`
	DeviceModel       string            `json:"device_model"`
	AdID              string            `json:"ad_id"`
	Tags              map[string]string `json:"tags"`
	LastActive        int               `json:"last_active"`
	AmountSpent       float32           `json:"amount_spent"`
	CreatedAt         int               `json:"created_at"`
	InvalidIdentifier bool              `json:"invalid_identifier"`
	BadgeCount        int               `json:"badge_count"`
	TestType          int               `json:"test_type,omitempty"`
	IP                string            `json:"ip,omitempty"`
	ExternalUserID    string            `json:"external_user_id,omitempty"`
}

// PlayerRequest represents a request to create/update a player.
type PlayerRequest struct {
	AppID string `json:"app_id"`
	// Required The device's platform:
	DeviceType int `json:"device_type"`
	// For Push Notifications, this is the Push Token Identifier from Google or Apple.
	// For Apple Push identifiers, you must strip all non alphanumeric characters.
	Identifier string `json:"identifier,omitempty"`
	// Only required if you have enabled Identity Verification and device_type is 11 (Email) or 14 SMS (coming soon).
	IdentifierAuthHash string `json:"identifier_auth_hash,omitempty"`
	// Language code. Typically lower case two letters, except for Chinese where it must be one of zh-Hans or zh-Hant. Example: en
	Language string `json:"language,omitempty"`
	// Number of seconds away from UTC. Example: -28800
	Timezone int `json:"timezone,omitempty"`
	// Version of your app. Example: 1.1
	GameVersion string `json:"game_version,omitempty"`
	// Device operating system version. Example: 7.0.4
	DeviceOS string `json:"device_os,omitempty"`
	// Device make and model. Example: iPhone5,1
	DeviceModel string `json:"device_model,omitempty"`
	// The ad id for the device's platform:
	// Android = Advertising Id
	// iOS = identifierForVendor
	// WP8.1 = AdvertisingId
	AdID string `json:"ad_id,omitempty"`
	// Name and version of the plugin that's calling this API method (if any)
	SDK string `json:"sdk,omitempty"`
	// Number of times the user has played the game, defaults to 1
	SessionCount int `json:"session_count,omitempty"`
	// Custom tags for the player. Only support string key value pairs.
	// Does not support arrays or other nested objects. Example: {"foo":"bar","this":"that"}
	Tags map[string]string `json:"tags,omitempty"`
	// Amount the user has spent in USD, up to two decimal places
	AmountSpent float32 `json:"amount_spent,omitempty"`
	// Set Automatically based on the date the request was made.
	// Unix timestamp in seconds indicating date and time when the device downloaded the app or subscribed to the website.
	CreatedAt int `json:"created_at,omitempty"`
	// Seconds player was running your app.
	Playtime int `json:"playtime,omitempty"`
	// Set Automatically based on the date the request was made.
	// Unix timestamp in seconds indicating date and time when the device last used the app or website.
	LastActive int `json:"last_active,omitempty"`
	// This is used in deciding whether to use your iOS Sandbox or Production push certificate when sending a push when both have been uploaded.
	// Set to the iOS provisioning profile that was used to build your app.
	// 1 = Development
	// 2 = Ad-Hoc
	// Omit this field for App Store builds.
	TestType int `json:"test_type,omitempty"`
	// 1 = subscribed
	// -2 = unsubscribed
	// iOS - These values are set each time the user opens the app from the SDK. Use the SDK function set Subscription instead.
	// Android - You may set this but you can no longer use the SDK method setSubscription later in your app as it will create synchronization issues.
	NotificationTypes string `json:"notification_types,omitempty"`
	// Longitude of the device, used for geotagging to segment on.
	Long float64 `json:"long,omitempty"`
	// Latitude of the device, used for geotagging to segment on.
	Lat float64 `json:"lat,omitempty"`
	// Country code in the ISO 3166-1 Alpha 2 format
	Country string `json:"country,omitempty"`
	// A custom user ID
	ExternalUserID string `json:"external_user_id,omitempty"`
	// Only required if you have enabled Identity Verification.
	ExternalUserIDAuthHash string `json:"external_user_id_auth_hash,omitempty"`
	// Current iOS badge count displayed on the app icon
	// NOTE: Not supported for apps created after June 2018, since badge count for apps created after this date are handled on the client.
	BadgeCount int `json:"badge_count"`
}

// PlayerListOptions specifies the parameters to the PlayersService.List method
type PlayerListOptions struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// PlayerGetOptions specifies the parameters to the PlayersService.Get method
type PlayerGetOptions struct {
	EmailAuthHash string `json:"email_auth_hash"`
}

// UpdateTagsWithExternalUserIDOptions specifies the parameters to the PlayersService.UpdateTagsWithExternalUserID method
type UpdateTagsWithExternalUserIDOptions struct {
	Tags map[string]string `json:"tags,omitempty"`
}

// PlayerListResponse wraps the standard http.Response for the
// PlayersService.List method
type PlayerListResponse struct {
	TotalCount int `json:"total_count"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	Players    []Player
}

// PlayerCreateResponse wraps the standard http.Response for the
// PlayersService.Create method
type PlayerCreateResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// PlayerOnSessionOptions specifies the parameters to the
// PlayersService.OnSession method
type PlayerOnSessionOptions struct {
	Identifier  string            `json:"identifier,omitempty"`
	Language    string            `json:"language,omitempty"`
	Timezone    int               `json:"timezone,omitempty"`
	GameVersion string            `json:"game_version,omitempty"`
	DeviceOS    string            `json:"device_os,omitempty"`
	AdID        string            `json:"ad_id,omitempty"`
	SDK         string            `json:"sdk,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// Purchase represents a purchase in the options of the
// PlayersService.OnPurchase method
type Purchase struct {
	SKU    string  `json:"sku"`
	Amount float32 `json:"amount"`
	ISO    string  `json:"iso"`
}

// PlayerOnPurchaseOptions specifies the parameters to the
// PlayersService.OnPurchase method
type PlayerOnPurchaseOptions struct {
	Purchases []Purchase `json:"purchases"`
	Existing  bool       `json:"existing,omitempty"`
}

// PlayerOnFocusOptions specifies the parameters to the
// PlayersService.OnFocus method
type PlayerOnFocusOptions struct {
	State      string `json:"state"`
	ActiveTime int    `json:"active_time"`
}

// PlayerCSVExportOptions specifies the parameters to the
// PlayersService.CSVExport method
type PlayerCSVExportOptions struct {
	ExtraFields     []string `json:"extra_fields"`
	LastActiveSince int      `json:"last_active_since"`
	SegmentName     string   `json:"segment_name"`
}

// PlayerCSVExportResponse wraps the standard http.Response for the
// PlayersService.CSVExport method
type PlayerCSVExportResponse struct {
	CSVFileURL string `json:"csv_file_url"`
}

// List the players.
//
// OneSignal API docs: https://documentation.onesignal.com/docs/players-view-devices
func (s *PlayersService) List(opt *PlayerListOptions) (*PlayerListResponse, *http.Response, error) {
	// build the URL with the query string
	u, err := url.Parse("/players")
	if err != nil {
		return nil, nil, err
	}
	q := u.Query()
	q.Set("app_id", s.client.appID)
	q.Set("limit", strconv.Itoa(opt.Limit))
	q.Set("offset", strconv.Itoa(opt.Offset))
	u.RawQuery = q.Encode()

	// create the request
	req, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	plResp := &PlayerListResponse{}
	resp, err := s.client.Do(req, plResp)
	if err != nil {
		return nil, resp, err
	}

	return plResp, resp, err
}

// Get a single player.
//
// OneSignal API docs: https://documentation.onesignal.com/reference/view-device
func (s *PlayersService) Get(playerID string, opt ...PlayerGetOptions) (*Player, *http.Response, error) {
	// build the URL
	path := fmt.Sprintf("/players/%s?app_id=%s", playerID, s.client.appID)
	u, err := url.Parse(path)
	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	q.Set("app_id", s.client.appID)
	if len(opt) > 0 {
		q.Set("email_auth_hash", opt[0].EmailAuthHash)
	}
	u.RawQuery = q.Encode()
	// create the request
	req, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	plResp := new(Player)
	resp, err := s.client.Do(req, plResp)
	if err != nil {
		return nil, resp, err
	}
	plResp.ID = playerID

	return plResp, resp, err
}

// Create a player.
//
// OneSignal API docs:
// https://documentation.onesignal.com/docs/players-add-a-device
func (s *PlayersService) Create(player PlayerRequest) (*PlayerCreateResponse, *http.Response, error) {
	// build the URL
	u, err := url.Parse("/players")
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("POST", u.String(), player)
	if err != nil {
		return nil, nil, err
	}

	plResp := &PlayerCreateResponse{}
	resp, err := s.client.Do(req, plResp)
	if err != nil {
		return nil, resp, err
	}

	return plResp, resp, err
}

// Generate a link to download a CSV list of all the players.
//
// OneSignal API docs:
// https://documentation.onesignal.com/docs/players_csv_export
func (s *PlayersService) CSVExport(opt ...PlayerCSVExportOptions) (*PlayerCSVExportResponse, *http.Response, error) {
	// build the URL with the query string
	u, err := url.Parse("/players/csv_export")
	if err != nil {
		return nil, nil, err
	}
	q := u.Query()
	q.Set("app_id", s.client.appID)
	u.RawQuery = q.Encode()

	// create the request
	var op *PlayerCSVExportOptions
	if len(opt) > 0 {
		op = &opt[0]
	}
	req, err := s.client.NewRequest("POST", u.String(), op)
	if err != nil {
		return nil, nil, err
	}

	plResp := &PlayerCSVExportResponse{}
	resp, err := s.client.Do(req, plResp)
	if err != nil {
		return nil, resp, err
	}

	return plResp, resp, err
}

// Update a player.
//
// OneSignal API docs: https://documentation.onesignal.com/reference/edit-device
func (s *PlayersService) Update(playerID string, player PlayerRequest) (*SuccessResponse, *http.Response, error) {
	// build the URL
	path := fmt.Sprintf("/players/%s", playerID)
	u, err := url.Parse(path)
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("PUT", u.String(), player)
	if err != nil {
		return nil, nil, err
	}

	plResp := &SuccessResponse{}
	resp, err := s.client.Do(req, plResp)
	if err != nil {
		return nil, resp, err
	}

	return plResp, resp, err
}

// Update an existing device's tags in one of your OneSignal apps using the External User ID.
//
// OneSignal API docs: https://documentation.onesignal.com/reference/edit-tags-with-external-user-id
func (s *PlayersService) UpdateTagsWithExternalUserID(ExternalUserID string, opt UpdateTagsWithExternalUserIDOptions) (*SuccessResponse, *http.Response, error) {
	// build the URL
	path := fmt.Sprintf("/apps/%s/users/%s", s.client.appID, ExternalUserID)
	u, err := url.Parse(path)
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("PUT", u.String(), opt)
	if err != nil {
		return nil, nil, err
	}

	plResp := &SuccessResponse{}
	resp, err := s.client.Do(req, plResp)
	if err != nil {
		return nil, resp, err
	}

	return plResp, resp, err
}
