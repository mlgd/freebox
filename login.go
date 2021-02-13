package freebox

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	endPointLogin              = "login"
	endPointLoginAuthorization = "login/authorize"
	endPointLoginSession       = "login/session"
	endPointLoginLogout        = "login/logout"

	// AuthorizationStatusUnknown :
	AuthorizationStatusUnknown AuthorizationStatus = "unknown"
	// AuthorizationStatusPending :
	AuthorizationStatusPending AuthorizationStatus = "pending"
	// AuthorizationStatusTimeout :
	AuthorizationStatusTimeout AuthorizationStatus = "timeout"
	// AuthorizationStatusGranted :
	AuthorizationStatusGranted AuthorizationStatus = "granted"
	// AuthorizationStatusDenied :
	AuthorizationStatusDenied AuthorizationStatus = "denied"

	// AppPermissionParental :
	AppPermissionParental AppPermission = "parental"
	// AppPermissionContacts :
	AppPermissionContacts AppPermission = "contacts"
	// AppPermissionExplorer :
	AppPermissionExplorer AppPermission = "explorer"
	// AppPermissionTV :
	AppPermissionTV AppPermission = "tv"
	// AppPermissionWDO :
	AppPermissionWDO AppPermission = "wdo"
	// AppPermissionDownloader :
	AppPermissionDownloader AppPermission = "downloader"
	// AppPermissionProfile :
	AppPermissionProfile AppPermission = "profile"
	// AppPermissionCamera :
	AppPermissionCamera AppPermission = "camera"
	// AppPermissionSettings :
	AppPermissionSettings AppPermission = "settings"
	// AppPermissionCalls :
	AppPermissionCalls AppPermission = "calls"
	// AppPermissionHome :
	AppPermissionHome AppPermission = "home"
	// AppPermissionPVR :
	AppPermissionPVR AppPermission = "pvr"
	// AppPermissionVM :
	AppPermissionVM AppPermission = "vm"
	// AppPermissionPlayer :
	AppPermissionPlayer AppPermission = "player"
)

// AuthorizationStatus :
type AuthorizationStatus string

// AppPermission :
type AppPermission string

// TokenRequest :
type TokenRequest struct {
	AppID      string `json:"app_id"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	DeviceName string `json:"device_name"`
}

// TokenResponse :
type TokenResponse struct {
	Success   bool                 `json:"success"`
	ErrorCode string               `json:"error_code"`
	Message   string               `json:"msg"`
	Result    *TokenResponseResult `json:"result"`
}

// TokenResponseResult :
type TokenResponseResult struct {
	AppToken string `json:"app_token"`
	TrackID  int    `jons:"track_id"`
}

// TrackAuthorizationProgressResponse :
type TrackAuthorizationProgressResponse struct {
	Success   bool                                      `json:"success"`
	ErrorCode string                                    `json:"error_code"`
	Message   string                                    `json:"msg"`
	Result    *TrackAuthorizationProgressResponseResult `json:"result"`
}

// TrackAuthorizationProgressResponseResult :
type TrackAuthorizationProgressResponseResult struct {
	Status    AuthorizationStatus `json:"status"`
	Challenge string              `json:"challenge"`
}

// GetChallengeResponse :
type GetChallengeResponse struct {
	Success   bool                        `json:"success"`
	ErrorCode string                      `json:"error_code"`
	Message   string                      `json:"msg"`
	Result    *GetChallengeResponseResult `json:"result"`
}

// GetChallengeResponseResult :
type GetChallengeResponseResult struct {
	LoggedIn     bool   `json:"logged_in"`
	Challenge    string `json:"challenge"`
	PasswordSalt string `json:"password_salt"`
	PasswordSet  bool   `json:"password_set"`
}

// OpenSessionRequest :
type OpenSessionRequest struct {
	AppID    string `json:"app_id"`
	Password string `jaons:"password"`
}

// OpenSessionResponse :
type OpenSessionResponse struct {
	Success   bool                       `json:"success"`
	ErrorCode string                     `json:"error_code"`
	Message   string                     `json:"msg"`
	Result    *OpenSessionResponseResult `json:"result"`
}

// OpenSessionResponseResult :
type OpenSessionResponseResult struct {
	SessionToken string                 `json:"session_token"`
	Challenge    string                 `json:"challenge"`
	Permissions  map[AppPermission]bool `json:"permissions"`
	PasswordSalt string                 `json:"password_salt"`
	PasswordSet  bool                   `json:"password_set"`
}

// CloseSessionResponse :
type CloseSessionResponse struct {
	Success   bool   `json:"success"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"msg"`
}

// RequestAuthorization :
func (device Device) RequestAuthorization() (response *TokenResponse, err error) {
	api := fmt.Sprintf("%s%s/", device.url(), endPointLoginAuthorization)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api, nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// TrackAuthorizationProgress :
func (device Device) TrackAuthorizationProgress(trackID int) (response *TrackAuthorizationProgressResponse, err error) {
	api := fmt.Sprintf("%s%s/%d", device.url(), endPointLoginAuthorization, trackID)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api, nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetChallenge :
func (device Device) GetChallenge() (response *GetChallengeResponse, err error) {
	api := fmt.Sprintf("%s%s/", device.url(), endPointLogin)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, api, nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// OpenSession :
func (device Device) OpenSession(appID string, challenge string) (response *OpenSessionResponse, err error) {
	api := fmt.Sprintf("%s%s/", device.url(), endPointLoginSession)

	mac := hmac.New(sha1.New, []byte(appID))
	if _, err := mac.Write([]byte(challenge)); err != nil {
		return nil, err
	}
	password := fmt.Sprintf("%02x", mac.Sum(nil))
	request := &OpenSessionRequest{
		AppID:    appID,
		Password: password,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err = ioutil.ReadAll(httpResponse.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// CloseSession :
func (device Device) CloseSession() (response *CloseSessionResponse, err error) {
	api := fmt.Sprintf("%s%s/", device.url(), endPointLoginLogout)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, api, nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	httpResponse, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}
