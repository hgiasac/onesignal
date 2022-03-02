package onesignal

import (
	"net/http"
	"net/url"
	"strconv"
)

type MessageType string

// Describes whether to set or increase/decrease your app's iOS badge count by the ios_badgeCount specified count.
type IOSBadgeType string

// iOS: Focus Modes and Interruption Levels indicate the priority and delivery timing of a notification, to ‘interrupt’ the user.
// Up until iOS 15, Apple primarily focused on Critical notifications.
// https://documentation.onesignal.com/docs/ios-focus-modes-and-interruption-levels
type IOSInterruptionLevel string

// Push delayed option
type DelayedOption string

// Huawei Message Type
type HuaweiMsgType string

const (
	MessageTypePush  MessageType = "push"
	MessageTypeEmail MessageType = "email"
	MessageTypeSMS   MessageType = "sms"
	// None leaves the count unaffected.
	IOSBadgeTypeNone IOSBadgeType = "None"
	// SetTo directly sets the badge count to the number specified in ios_badgeCount.
	IOSBadgeTypeSetTo IOSBadgeType = "SetTo"
	// Increase adds the number specified in ios_badgeCount to the total. Use a negative number to decrease the badge count
	IOSBadgeTypeIncrease IOSBadgeType = "Increase"

	IOSInterruptionActive        IOSInterruptionLevel = "active"
	IOSInterruptionPassive       IOSInterruptionLevel = "passive"
	IOSInterruptionTimeSensitive IOSInterruptionLevel = "time_sensitive"
	IOSInterruptionCritical      IOSInterruptionLevel = "critical"

	// Deliver at a specific time-of-day in each users own timezone
	DelayedOptionTimezone DelayedOption = "timezone"
	// Same as Intelligent Delivery . (Deliver at the same time of day as each user last used your app).
	// https://documentation.onesignal.com/docs/sending-notifications#intelligent-delivery
	DelayedOptionLastActive DelayedOption = "last-active"

	HuaweiMsgTypeData    HuaweiMsgType = "data"
	HuaweiMsgTypeMessage HuaweiMsgType = "message"
)

// NotificationsService handles communication with the notification related
// methods of the OneSignal API.
type NotificationsService struct {
	client *Client
}

// Notification  represents a OneSignal notification.
type Notification struct {
	ID         string            `json:"id"`
	Successful int               `json:"successful"`
	Failed     int               `json:"failed"`
	Converted  int               `json:"converted"`
	Remaining  int               `json:"remaining"`
	QueuedAt   int               `json:"queued_at"`
	SendAfter  int               `json:"send_after"`
	URL        string            `json:"url"`
	Data       interface{}       `json:"data"`
	Canceled   bool              `json:"canceled"`
	Headings   map[string]string `json:"headings"`
	Contents   map[string]string `json:"contents"`
}

// NotificationRequest represents a request to create a notification.
type NotificationRequest struct {
	AppID string `json:"app_id"`
	// An identifier for tracking message within the OneSignal dashboard or export analytics.
	// Optional for Push. Not shown to end user.
	Name string `json:"name,omitempty"`
	// The notification's content (excluding the title), a map of language codes to text for each language.
	// Required unless content_available=true or template_id is set.
	Contents map[string]string `json:"contents,omitempty"`
	// The notification's title, a map of language codes to text for each language.
	// Each hash must have a language code string for a key,
	// mapped to the localized text you would like users to receive for that language.
	// Required for Huawei
	// Web Push requires a heading but can be omitted from request since defaults to the Site Name set in OneSignal Settings.
	Headings map[string]string      `json:"headings,omitempty"`
	Subtitle map[string]interface{} `json:"subtitle,omitempty"`
	// Indicates whether to send to all devices registered under your app's Apple iOS platform.
	IsIOS bool `json:"isIos,omitempty"`
	// Indicates whether to send to all devices registered under your app's Google Android platform.
	IsAndroid bool `json:"isAndroid,omitempty"`
	// Indicates whether to send to all devices registered under your app's Windows platform.
	IsWP_WNS bool `json:"isWP_WNS,omitempty"`
	// Indicates whether to send to all devices registered under your app's Huawei Android platform.
	IsHuawei bool `json:"isHuawei,omitempty"`
	// Indicates whether to send to all devices registered under your app's Amazon Fire platform.
	IsADM bool `json:"isAdm,omitempty"`
	// Indicates whether to send to all devices registered under your app's Google Chrome Apps & Extension platform.
	// This flag is not used for web push Please see isChromeWeb for sending to web push users.
	IsChrome bool `json:"isChrome,omitempty"`
	// Indicates whether to send to all Google Chrome, Chrome on Android,
	// and Mozilla Firefox users registered under your Chrome & Firefox web push platform.
	IsChromeWeb bool `json:"isChromeWeb,omitempty"`
	// Indicates whether to send to all Mozilla Firefox desktop users registered under your Firefox web push platform.
	IsFirefox bool `json:"isFirefox,omitempty"`
	// Does not support iOS Safari. Indicates whether to send to all Apple's Safari desktop users registered under your Safari web push platform.
	IsSafari bool `json:"isSafari,omitempty"`
	// Indicates whether to send to all subscribed web browser users, including Chrome, Firefox, and Safari.
	IsAnyWeb bool `json:"isAnyWeb,omitempty"`
	// Indicates if the message type when targeting with include_external_user_ids for cases
	// where an email, sms, and/or push subscribers have the same external user id.
	// Example: Use the string "push" to indicate you are sending a push notification or the string "email"for sending emails or "sms"for sending SMS.
	ChannelForExternalUserIDs MessageType `json:"channel_for_external_user_ids,omitempty"`
	IncludedSegments          []string    `json:"included_segments,omitempty"`
	ExcludedSegments          []string    `json:"excluded_segments,omitempty"`
	IncludeExternalUserIDs    []string    `json:"include_external_user_ids,omitempty"`
	IncludeEmailTokens        []string    `json:"include_email_tokens,omitempty"`
	IncludePhoneNumber        []string    `json:"include_phone_numbers,omitempty"`
	IncludePlayerIDs          []string    `json:"include_player_ids,omitempty"`
	IncludeIOSTokens          []string    `json:"include_ios_tokens,omitempty"`
	IncludeAndroidRegIDs      []string    `json:"include_android_reg_ids,omitempty"`
	IncludeWPURIs             []string    `json:"include_wp_uris,omitempty"`
	IncludeWPWNSURIs          []string    `json:"include_wp_wns_uris,omitempty"`
	IncludeAmazonRegIDs       []string    `json:"include_amazon_reg_ids,omitempty"`
	IncludeChromeRegIDs       []string    `json:"include_chrome_reg_ids,omitempty"`
	IncludeChromeWebRegIDs    []string    `json:"include_chrome_web_reg_ids,omitempty"`
	AppIDs                    []string    `json:"app_ids,omitempty"`
	Tags                      interface{} `json:"tags,omitempty"`

	// Describes whether to set or increase/decrease your app's iOS badge count by the ios_badgeCount specified count.
	// Can specify None, SetTo, or Increase.
	IOSBadgeType IOSBadgeType `json:"ios_badgeType,omitempty"`
	// Used with ios_badgeType, describes the value to set or amount to increase/decrease your app's iOS badge count by.
	// You can use a negative number to decrease the badge count when used with an ios_badgeType of Increase.
	IOSBadgeCount int `json:"ios_badgeCount,omitempty"`
	// Sound file that is included in your app to play instead of the default device notification sound.
	// Pass nil to disable vibration and sound for the notification.
	IOSSound string `json:"ios_sound,omitempty"`
	// Adds media attachments to notifications. Set as JSON object, key as a media id of your choice and the value as a valid local filename or URL.
	// User must press and hold on the notification to view.
	// Do not set mutable_content to download attachments. The OneSignal SDK does this automatically.
	IOSAttachments map[string]string `json:"ios_attachments,omitempty"`
	// iOS: Category APS payload, use with registerUserNotificationSettings:categories in your Objective-C / Swift code.
	// Example: calendar category which contains actions like accept and decline
	// iOS 10+ This will trigger your UNNotificationContentExtension whose ID matches this category.
	IOSCategory string `json:"ios_category,omitempty"`
	// deprecated: this field doesn't work on Android 8 (Oreo) and newer devices!
	AndroidSound string `json:"android_sound,omitempty"`
	// deprecated: this field ONLY works on EMUI 5 (Android 7 based) and older devices.
	HuaweiSound string `json:"huawei_sound,omitempty"`
	// deprecated: this field doesn't work on Android 8 (Oreo) and newer devices!
	ADMSound string `json:"adm_sound,omitempty"`
	// Sound file that is included in your app to play instead of the default device notification sound.
	WPWNSSound string `json:"wp_wns_sound,omitempty"`
	// iOS 10+, Android Only one notification with the same id will be shown on the device.
	// Use the same id to update an existing notification instead of showing a new one. Limit of 64 characters.
	CollapseID string `json:"collapse_id,omitempty"`
	// Display multiple notifications at once with different topics.
	WebPushTopic string `json:"web_push_topic,omitempty"`
	// iOS 10+	iOS can localize push notification messages on the client using special parameters such as loc-key.
	// When using the Create Notification endpoint,you must include these parameters inside of a field called apns_alert.
	// https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/CreatingtheNotificationPayload.html#//apple_ref/doc/uid/TP40008194-CH10-SW1
	APNSAlert map[string]interface{} `json:"apns_alert"`
	// A custom map of data that is passed back to your app.
	// Can use up to 2048 bytes of data.
	Data interface{} `json:"data,omitempty"`
	// iOS 8.0+, Android 4.1+, and derivatives like Amazon: Buttons to add to the notification. Icon only works for Android.
	// Buttons show in reverse order of array position i.e. Last item in array shows as first button on device.
	Buttons   interface{} `json:"buttons,omitempty"`
	IconType  string      `json:"icon_type,omitempty"`
	SmallIcon string      `json:"small_icon,omitempty"`
	LargeIcon string      `json:"large_icon,omitempty"`

	// Chrome 48+: Add action buttons to the notification. The id field is required.
	WebButtons interface{} `json:"web_buttons,omitempty"`
	// Android: Picture to display in the expanded view. Can be a drawable resource name or a URL.
	BigPicture   string `json:"big_picture,omitempty"`
	ADMSmallIcon string `json:"adm_small_icon,omitempty"`
	ADMLargeIcon string `json:"adm_large_icon,omitempty"`
	// Amazon: Picture to display in the expanded view. Can be a drawable resource name or a URL.
	ADMBigPicture string `json:"adm_big_picture,omitempty"`
	ChromeIcon    string `json:"chrome_icon,omitempty"`
	// ChromeApp: Large picture to display below the notification text. Must be a local URL.
	ChromeBigPicture string `json:"chrome_big_picture,omitempty"`
	ChromeWebIcon    string `json:"chrome_web_icon,omitempty"`
	FirefoxIcon      string `json:"firefox_icon,omitempty"`
	// The URL to open in the browser when a user clicks on the notification.
	URL string `json:"url,omitempty"`
	// Same as url but only sent to app platforms.
	// Including iOS, Android, macOS, Windows, ChromeApps, etc.
	AppURL string `json:"app_url,omitempty"`
	// Same as url but only sent to web push platforms.
	// Including Chrome, Firefox, Safari, Opera, etc.
	WebURL string `json:"web_url,omitempty"`

	// Schedule notification for future delivery. API defaults to UTC.
	SendAfter string `json:"send_after,omitempty"`
	// If send_after is used, this takes effect after the send_after time has elapsed.
	// Cannot be used if Throttling enabled. Set throttle_rate_per_minute to 0 to disable throttling if enabled to use these features.
	DelayedOption DelayedOption `json:"delayed_option,omitempty"`
	// Use with delayed_option=timezone.
	DeliveryTimeOfDay string `json:"delivery_time_of_day,omitempty"`
	// Sets the devices LED notification light if the device has one. ARGB Hex format.
	// deprecated: this field doesn't work on Android 8 (Oreo) and newer devices!
	// Android, Chrome, ChromeWeb	Delivery priority through the push server (example GCM/FCM).
	// Pass 10 for high priority or any other integer for normal priority.
	// Defaults to normal priority for Android and high for iOS.
	// For Android 6.0+ devices setting priority to high will wake the device out of doze mode.
	Priority uint `json:"priority,omitempty"`
	// iOS Set the value to voip for sending VoIP Notifications
	// This field maps to the APNS header apns-push-type.
	// Note: alert and background are automatically set by OneSignal
	// https://documentation.onesignal.com/docs/voip-notifications
	APNSPushTypeOverride string `json:"apns_push_type_override,omitempty"`
	// Time To Live - In seconds. The notification will be expired if the device does not come back online within this time.
	// The default is 259,200 seconds (3 days).
	// Max value to set is 2419200 seconds (28 days).
	TTL uint `json:"ttl"`
	// Apps with throttling enabled
	// - does not work with timezone or intelligent delivery, throttling limits will take precedence. Set to 0 if you want to use timezone or intelligent delivery.
	// - the parameter value will be used to override the default application throttling value set from the dashboard settings.
	// - parameter value 0 indicates not to apply throttling to the notification.
	// - if the parameter is not passed then the default app throttling value will be applied to the notification.
	// Apps with throttling disabled
	// - this parameter can be used to throttle delivery for the notification even though throttling is not enabled at the application level.
	// https://documentation.onesignal.com/docs/throttling
	ThrottleRatePerMinute uint `json:"throttle_rate_per_minute,omitempty"`
	// When frequency capping is enabled for the app, sending true will apply the frequency capping to the notification.
	// If the parameter is not included, the default behavior is to apply frequency capping if the setting is enabled for the app.
	// Setting the parameter to false will override the frequency capping, meaning that the notification will be sent without considering frequency capping.
	EnableFrequencyCap bool `json:"enable_frequency_cap,omitempty"`

	AndroidLEDColor string `json:"android_led_color,omitempty"`
	// Sets the devices LED notification light if the device has one. ARGB Hex format.
	// deprecated: this field doesn't work on Android 8 (Oreo) and newer devices!
	HuaweiLEDColor string `json:"huawei_led_color,omitempty"`
	// Sets the background color of the notification circle to the left of the notification text.
	// Only applies to apps targeting Android API level 21+ on Android 5.0+ devices.
	AndroidAccentColor string `json:"android_accent_color,omitempty"`
	// Accent Color used on Action Buttons and Group overflow count.
	// Uses RGB Hex value (E.g. #9900FF).
	// Defaults to device’s theme color if not set.
	HuaweiAccentColor string `json:"huawei_accent_color,omitempty"`
	// deprecated: this field doesn't work on Android 8 (Oreo) and newer devices!
	AndroidVisibility int `json:"android_visibility,omitempty"`
	// deprecated: this field ONLY works on EMUI 5 (Android 7 based) and older devices.
	HuaweiVisibility int `json:"huawei_visibility,omitempty"`
	// Sending true wakes your app from background to run custom native code (Apple interprets this as content-available=1).
	// Note: Not applicable if the app is in the "force-quit" state (i.e app was swiped away).
	// Omit the contents field to prevent displaying a visible notification.
	ContentAvailable      bool `json:"content_available,omitempty"`
	AndroidBackgroundData bool `json:"android_background_data,omitempty"`
	AmazonBackgroundData  bool `json:"amazon_background_data,omitempty"`
	// Use a template you setup on our dashboard.
	// The template_id is the UUID found in the URL when viewing a template on our dashboard.
	TemplateID string `json:"template_id,omitempty"`
	// Android: Notifications with the same group will be stacked together using Android's Notification Grouping feature.
	AndroidGroup string `json:"android_group,omitempty"`
	// Android: Summary message to display when 2+ notifications are stacked together. Default is "# new messages".
	// Include $[notif_count] in your message and it will be replaced with the current number.
	// Note: This only works for Android 6 and older. Android 7+ allows full expansion of all message.
	AndroidGroupMessage interface{} `json:"android_group_message,omitempty"`
	// Amazon: Notifications with the same group will be stacked together using Android's Notification Grouping feature.
	ADMGroup string `json:"adm_group,omitempty"`
	// Amazon: Summary message to display when 2+ notifications are stacked together. Default is "# new messages".
	// Include $[notif_count] in your message and it will be replaced with the current number. "en" (English) is required.
	ADMGroupMessage interface{} `json:"adm_group_message,omitempty"`
	// iOS 12+ This parameter is supported in iOS 12 and above. It allows you to group related notifications together.
	ThreadID string `json:"thread_id,omitempty"`
	// iOS 12+ When using thread_id to create grouped notifications in iOS 12+, you can also control the summary.
	// For example, a grouped notification can say "12 more notifications from John Doe".
	SummaryArg string `json:"summary_arg,omitempty"`
	// iOS 12+ When using thread_id, you can also control the count of the number of notifications in the group.
	// For example, if the group already has 12 notifications, and you send a new notification with summary_arg_count = 2,
	// the new total will be 14 and the summary will be "14 more notifications from summary_arg"
	SummaryArgCount int `json:"summary_arg_count,omitempty"`
	// iOS 15+ Relevance Score is a score to be set per notification to indicate how it should be displayed when grouped.
	// https://documentation.onesignal.com/docs/ios-relevance-score
	IOSRelevanceScore float32 `json:"ios_relevance_score,omitempty"`
	// iOS 15+ Focus Modes and Interruption Levels indicate the priority and delivery timing of a notification, to ‘interrupt’ the user.
	IosInterruptionLevel IOSInterruptionLevel `json:"ios_interruption_level,omitempty"`

	Filters    interface{} `json:"filters,omitempty"`
	ExternalID string      `json:"external_id,omitempty"`
	// Use to target a specific experience in your App Clip, or to target your notification to a specific window in a multi-scene App.
	// https://documentation.onesignal.com/docs/app-clip-support
	TargetContentIdentifier string `json:"target_content_identifier,omitempty"`
	// Use "data" or "message" depending on the type of notification you are sending
	// https://documentation.onesignal.com/docs/data-notifications
	HuaweiMsgType string `json:"huawei_msg_type,omitempty"`
	// Huawei: Picture to display in the expanded view. Can be a drawable resource name or a URL
	HuaweiBigPicture string `json:"huawei_big_picture,omitempty"`
	// Chrome 56+: Sets the web push notification's large image to be shown below the notification's title and text.
	ChromeWebImage string `json:"chrome_web_image,omitempty"`
	// The Android Oreo Notification Category to send the notification under.
	AndroidChannelID string `json:"android_channel_id,omitempty"`
	// The Android Oreo Notification Category to send the notification under
	HuaweiChannelID string `json:"huawei_channel_id,omitempty"`
	// email specific content
	EmailSubject string `json:"email_subject,omitempty"`
	// Required unless template_id is set.
	// The body of the email you wish to send. Typically, customers include their own HTML templates here.
	// Must include [unsubscribe_url] in an <a> tag somewhere in the email.
	// Note: any malformed HTML content will be sent to users. Please double-check your HTML is valid.
	EmailBody string `json:"email_body,omitempty"`
	// The name the email is from. If not specified, will default to "from name" set in the OneSignal Dashboard Email Settings.
	EmailFromName string `json:"email_from_name,omitempty"`
	// The email address the email is from.
	// If not specified, will default to "from email" set in the OneSignal Dashboard Email Settings.
	EmailFromAddress string `json:"email_from_address,omitempty"`
	// Default is false. If set to true the URLs included in the email will not change to link tracking URLs and will stay the same as originally set.
	// Best used for emails containing Universal Links.
	DisableEmailClickTracking bool `json:"disable_email_click_tracking,omitempty"`
	// sms specific content
	SMSFrom string `json:"sms_from,omitempty"`
	// URLs for the media files to be attached to the SMS content.
	// Limit: 10 media urls with a total max. size of 5MBs.
	SMSMediaURLs []string `json:"sms_media_urls,omitempty"`
}

// NotificationCreateResponse wraps the standard http.Response for the
// NotificationsService.Create method
type NotificationCreateResponse struct {
	ID         string      `json:"id"`
	Recipients int         `json:"recipients"`
	Errors     interface{} `json:"errors"`
}

// NotificationListOptions specifies the parameters to the
// NotificationsService.List method
type NotificationListOptions struct {
	AppID  string `json:"app_id"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// NotificationListResponse wraps the standard http.Response for the
// NotificationsService.List method
type NotificationListResponse struct {
	TotalCount    int `json:"total_count"`
	Offset        int `json:"offset"`
	Limit         int `json:"limit"`
	Notifications []Notification
}

// NotificationUpdateOptions specifies the parameters to the
// NotificationsService.Get method
type NotificationGetOptions struct {
	AppID string `json:"app_id"`
}

// NotificationUpdateOptions specifies the parameters to the
// NotificationsService.Update method
type NotificationUpdateOptions struct {
	AppID  string `json:"app_id"`
	Opened bool   `json:"opened"`
}

// NotificationDeleteOptions specifies the parameters to the
// NotificationsService.Delete method
type NotificationDeleteOptions struct {
	AppID string `json:"app_id"`
}

// List the notifications.
//
// OneSignal API docs:
// https://documentation.onesignal.com/docs/notifications-view-notifications
func (s *NotificationsService) List(opt *NotificationListOptions) (*NotificationListResponse, *http.Response, error) {
	// build the URL with the query string
	u, err := url.Parse("/notifications")
	if err != nil {
		return nil, nil, err
	}
	q := u.Query()
	q.Set("app_id", opt.AppID)
	q.Set("limit", strconv.Itoa(opt.Limit))
	q.Set("offset", strconv.Itoa(opt.Offset))
	u.RawQuery = q.Encode()

	// create the request
	req, err := s.client.NewRequest("GET", u.String(), nil, APP)
	if err != nil {
		return nil, nil, err
	}

	notifResp := &NotificationListResponse{}
	resp, err := s.client.Do(req, notifResp)
	if err != nil {
		return nil, resp, err
	}

	return notifResp, resp, err
}

// Get a single notification.
//
// OneSignal API docs:
// https://documentation.onesignal.com/docs/notificationsid-view-notification
func (s *NotificationsService) Get(notificationID string, opt *NotificationGetOptions) (*Notification, *http.Response, error) {
	// build the URL with the query string
	u, err := url.Parse("/notifications/" + notificationID)
	if err != nil {
		return nil, nil, err
	}
	q := u.Query()
	q.Set("app_id", opt.AppID)
	u.RawQuery = q.Encode()

	// create the request
	req, err := s.client.NewRequest("GET", u.String(), nil, APP)
	if err != nil {
		return nil, nil, err
	}

	notif := &Notification{}
	resp, err := s.client.Do(req, notif)
	if err != nil {
		return nil, resp, err
	}

	return notif, resp, err
}

// Create a notification.
//
// OneSignal API docs:
// https://documentation.onesignal.com/docs/notifications-create-notification
func (s *NotificationsService) Create(opt *NotificationRequest) (*NotificationCreateResponse, *http.Response, error) {
	// build the URL
	u, err := url.Parse("/notifications")
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("POST", u.String(), opt, APP)
	if err != nil {
		return nil, nil, err
	}

	createRes := &NotificationCreateResponse{}
	resp, err := s.client.Do(req, createRes)
	if err != nil {
		return nil, resp, err
	}

	return createRes, resp, err
}

// Update a notification.
//
// OneSignal API docs:
// https://documentation.onesignal.com/docs/notificationsid-track-open
func (s *NotificationsService) Update(notificationID string, opt *NotificationUpdateOptions) (*SuccessResponse, *http.Response, error) {
	// build the URL
	u, err := url.Parse("/notifications/" + notificationID)
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("PUT", u.String(), opt, APP)
	if err != nil {
		return nil, nil, err
	}

	updateRes := &SuccessResponse{}
	resp, err := s.client.Do(req, updateRes)
	if err != nil {
		return nil, resp, err
	}

	return updateRes, resp, err
}

// Delete a notification.
//
// OneSignal API docs:
// https://documentation.onesignal.com/docs/notificationsid-cancel-notification
func (s *NotificationsService) Delete(notificationID string, opt *NotificationDeleteOptions) (*SuccessResponse, *http.Response, error) {
	// build the URL
	u, err := url.Parse("/notifications/" + notificationID)
	if err != nil {
		return nil, nil, err
	}

	// create the request
	req, err := s.client.NewRequest("DELETE", u.String(), opt, APP)
	if err != nil {
		return nil, nil, err
	}

	deleteRes := &SuccessResponse{}
	resp, err := s.client.Do(req, deleteRes)
	if err != nil {
		return nil, resp, err
	}

	return deleteRes, resp, err
}
