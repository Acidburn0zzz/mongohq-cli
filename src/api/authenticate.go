package api

import (
  "net/http"
  "io/ioutil"
  "encoding/json"
  "errors"
  "bytes"
)

type AuthenticationArguments struct {
  Username string `json:"username"`
  Password string `json:"password"`
  GrantType string `json:"grant_type"`
  ClientId string `json:"client_id"`
}

func Authenticate(username, password, token string) (string, error) {
  var oauthToken string
  var authenticationError error
  var jsonResponse map[string]interface{}

  data, err := json.Marshal(AuthenticationArguments{Username: username, Password: password, GrantType: "password", ClientId: oauth_client_id})
  if err != nil {
    return oauthToken, errors.New("Error creating MongoHQ authentication request.")
  }

  request, err := http.NewRequest("POST", api_url("/oauth/token"), bytes.NewReader(data))

  if token != "" {
    request.Header.Add("X-Mongohq-Otp", token)
  }

  request.Header.Add("User-Agent", userAgent())
  request.Header.Add("Content-Type", "application/json")

  client := &http.Client{}
  response, err := client.Do(request)

  if err != nil {
    println(err.Error())
    authenticationError = errors.New("Error authenticating against MongoHQ.")
  } else if response.StatusCode >= 400 {
    if response.Header.Get("X-Mongohq-Otp") == "required; sms" {
      authenticationError = errors.New("2fa token required")
    } else if response.Header.Get("X-Mongohq-Otp") == "required; unconfigured" {
      authenticationError = errors.New("Account requires 2fa authentication.  Go to https://app.mongohq.com to configure")
    } else {
      authenticationError = errors.New("Error authenticating against MongoHQ.")
    }
  } else {
    responseBody, _ := ioutil.ReadAll(response.Body)
    _ = json.Unmarshal(responseBody, &jsonResponse)
    response.Body.Close()

    oauthToken = jsonResponse["access_token"].(string)
  }

  return oauthToken, authenticationError
}
