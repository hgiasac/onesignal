package onesignal

import (
	"net/http"
	"net/url"
	"time"
)

type APNSEnvironment string

const (
	APNSEnvSandbox    APNSEnvironment = "sandbox"
	APNSEnvProduction APNSEnvironment = "production"
)

// AppsService handles communication with the app methods of the OneSignal API.
type AppsService struct {
	client *UserClient
}

// App represents a OneSignal app.
type App struct {
	ID                               string          `json:"id"`
	Name                             string          `json:"name"`
	Players                          int             `json:"players"`
	MessagablePlayers                int             `json:"messagable_players"`
	UpdatedAt                        time.Time       `json:"updated_at"`
	CreatedAt                        time.Time       `json:"created_at"`
	GCMKey                           string          `json:"gcm_key"`
	ChromeKey                        string          `json:"chrome_key"`
	ChromeWebOrigin                  string          `json:"chrome_web_origin"`
	ChromeWebGCMSenderID             string          `json:"chrome_web_gcm_sender_id"`
	ChromeWebDefaultNotificationIcon string          `json:"chrome_web_default_notification_icon"`
	ChromeWebSubDomain               string          `json:"chrome_web_sub_domain"`
	APNSEnv                          APNSEnvironment `json:"apns_env"`
	APNSCertificates                 string          `json:"apns_certificates"`
	SafariAPNSCertificate            string          `json:"safari_apns_certificate"`
	SafariSiteOrigin                 string          `json:"safari_site_origin"`
	SafariPushID                     string          `json:"safari_push_id"`
	SafariIcon16x16                  string          `json:"safari_icon_16_16"`
	SafariIcon32x32                  string          `json:"safari_icon_32_32"`
	SafariIcon64x64                  string          `json:"safari_icon_64_64"`
	SafariIcon128x128                string          `json:"safari_icon_128_128"`
	SafariIcon256x256                string          `json:"safari_icon_256_256"`
	SiteName                         string          `json:"site_name"`
	BasicAuthKey                     string          `json:"basic_auth_key"`
}

// AppRequest represents a request to create/update an app.
type AppRequest struct {
	// Required: The name of your new app, as displayed on your apps list on the dashboard. This can be renamed later.
	Name string `json:"name"`
	// iOS: Either sandbox or production
	APNSEnv APNSEnvironment `json:"apns_env,omitempty"`
	// iOS: Your apple push notification p12 certificate file, converted to a string and Base64 encoded.
	APNSP12 string `json:"apns_p12,omitempty"`
	// iOS: Required if adding p12 certificate - Password for the apns_p12 file
	APNSP12Password string `json:"apns_p12_password,omitempty"`
	// Android: Your FCM Google Push Server Auth Key
	GCMKey string `json:"gcm_key,omitempty"`
	// Android: Your FCM Google Project number. Also know as Sender ID.
	AndroidGCMSenderID string `json:"android_gcm_sender_id,omitempty"`
	// Chrome (All Browsers except Safari) (Recommended): The URL to your website. This field is required if you wish to enable web push and specify other web push parameters.
	ChromeWebOrigin string `json:"chrome_web_origin,omitempty"`
	// Chrome (All Browsers except Safari): Your default notification icon. Should be 256x256 pixels, min 80x80.
	ChromeWebDefaultNotificationIcon string `json:"chrome_web_default_notification_icon,omitempty"`
	// Chrome (All Browsers except Safari): A subdomain of your choice in order to support Web Push on non-HTTPS websites. This field must be set in order for the chrome_web_gcm_sender_id property to be processed.
	ChromeWebSubDomain string `json:"chrome_web_sub_domain,omitempty"`
	// All Browsers (Recommended): The Site Name. Requires both chrome_web_origin and safari_site_origin to be set to add or update it.
	SiteName string `json:"site_name,omitempty"`
	// Safari (Recommended): The hostname to your website including http(s)://
	SafariSiteOrigin string `json:"safari_site_origin,omitempty"`
	// Safari: Your apple push notification p12 certificate file for Safari Push Notifications, converted to a string and Base64 encoded.
	SafariAPNSP12 string `json:"safari_apns_p12,omitempty"`
	// Safari: Password for safari_apns_p12 file
	SafariAPNSP12Password string `json:"safari_apns_p12_password,omitempty"`
	// Safari: A url for a 16x16 png notification icon. This is the only Safari icon URL you need to provide.
	SafariIcon16x16 string `json:"safari_icon_16_16,omitempty"`
	// Safari: A url for a 32x32 png notification icon. This is the only Safari icon URL you need to provide.
	SafariIcon32x32 string `json:"safari_icon_32_32,omitempty"`
	// Safari: A url for a 64x64 png notification icon. This is the only Safari icon URL you need to provide.
	SafariIcon64x64 string `json:"safari_icon_64_64,omitempty"`
	// Safari: A url for a 128x128 png notification icon. This is the only Safari icon URL you need to provide.
	SafariIcon128x128 string `json:"safari_icon_128_128,omitempty"`
	// Safari: A url for a 256x256 png notification icon. This is the only Safari icon URL you need to provide.
	SafariIcon256x256 string `json:"safari_icon_256_256,omitempty"`
}

// List the apps.
// https://documentation.onesignal.com/reference/view-apps-apps
func (s *AppsService) List() ([]App, *http.Response, error) {
	// build the URL
	u, err := url.Parse("/apps")
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var apps []App
	resp, err := s.client.Do(req, &apps)
	if err != nil {
		return nil, resp, err
	}

	return apps, resp, err
}

// Get a single app.
//
// OneSignal API docs: https://documentation.onesignal.com/reference/view-an-app
func (s *AppsService) Get(appID string) (*App, *http.Response, error) {
	// build the URL
	u, err := url.Parse("/apps/" + appID)
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	app := &App{}
	resp, err := s.client.Do(req, app)
	if err != nil {
		return nil, resp, err
	}

	return app, resp, err
}

// Create an app.
//
// OneSignal API docs: https://documentation.onesignal.com/reference/create-an-app
func (s *AppsService) Create(opt AppRequest) (*App, *http.Response, error) {
	// build the URL
	u, err := url.Parse("/apps")
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("POST", u.String(), opt)
	if err != nil {
		return nil, nil, err
	}

	app := &App{}
	resp, err := s.client.Do(req, app)
	if err != nil {
		return nil, resp, err
	}

	return app, resp, err
}

// Update an app.
//
// OneSignal API docs: https://documentation.onesignal.com/reference/update-an-app
func (s *AppsService) Update(appID string, opt AppRequest) (*App, *http.Response, error) {
	// build the URL
	u, err := url.Parse("/apps/" + appID)
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("PUT", u.String(), opt)
	if err != nil {
		return nil, nil, err
	}

	app := &App{}
	resp, err := s.client.Do(req, app)
	if err != nil {
		return nil, resp, err
	}

	return app, resp, err
}
