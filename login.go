package myvmware

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	baseURL = "https://my.vmware.com"
	initURL = baseURL + "/web/vmware/login"
	ssoURL  = baseURL + "/vmwauth/saml/SSO"
	authURL = "https://auth.vmware.com/oam/server/auth_cred_submit?Auth-AppID=WMVMWR"
)

func Login(username, password string) (*Client, error) {
	// The default cookie jar does not return an error ever, so we can ignore
	// it for now
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: jar}

	initRes, err := httpClient.Get(initURL)
	if err != nil {
		return nil, err
	}
	defer initRes.Body.Close()

	authRes, err := httpClient.PostForm(authURL, url.Values{
		"username": {username},
		"password": {password},
	})
	if err != nil {
		return nil, err
	}
	defer authRes.Body.Close()

	samlToken, err := getSAMLToken(authRes.Body)
	if err != nil {
		return nil, err
	}

	ssoRes, err := httpClient.PostForm(ssoURL, url.Values{
		"SAMLResponse": {samlToken},
	})
	if err != nil {
		return nil, err
	}
	defer ssoRes.Body.Close()

	// TODO: Do we actually need to add the cookie for all requests?
	//
	// httpClient.Transport = &transportWithHeaders{
	// 	transport: httpClient.Transport,
	// 	headers: map[string]string{
	// 		"foo": "bar",
	// 	},
	// }

	return &Client{HttpClient: httpClient}, nil
}
