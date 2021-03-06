package controllers

import (
  //"fmt"
  "api"
  "encoding/json"
  "io/ioutil"
  "os"
  "errors"
  "github.com/codegangsta/cli"
  "code.google.com/p/gopass"
  "fmt"
)

var credentialPath = os.Getenv("HOME") + "/.mongohq"
var credentialFile = credentialPath + "/credentials"

var Email, OauthToken string

func login() (string, string, error) {
  fmt.Println("Enter your MongoHQ credentials.")
  username := prompt("Email")
  password, err := gopass.GetPass("Password (typing will be hidden): ")

  if err != nil {
    return "", "", errors.New("Error returning password.  We may not be compliant with your system yet.  Please send us a message telling us about your system to support@mongohq.com.")
  }

  oauthToken, err := api.Authenticate(username, password, "")
  return processAuthenticationResponse(username, password, oauthToken, err) 
}

func processAuthenticationResponse(username, password, oauthToken string, err error) (string, string, error) { 
  if err != nil {
    if err.Error() == "2fa token required" {
      twoFactorToken := prompt("2fa token")
      oauthToken, err := api.Authenticate(username, password, twoFactorToken)
      return processAuthenticationResponse(username, password, oauthToken, err) 
    } else {
      return username, "", err
    }
  } else {
    err = storeCredentials(username, oauthToken)

    if err != nil {
      return username, oauthToken, err
    } else {
      fmt.Println("\nAuthentication complete.\n\n")
      return username, oauthToken, nil
    }
  }
}

func storeCredentials(username, oauth string) (error) {
  credentials := make(map[string]interface{})
  credentials["email"] = username
  credentials["oauth_token"] = oauth

  jsonText, _ := json.Marshal(credentials)

  err := os.MkdirAll(credentialPath, 0700)

  if err != nil {
    return errors.New("Error creating directory " + credentialPath)
  }

  err = ioutil.WriteFile(credentialFile, jsonText, 0500)

  if err != nil {
    err = errors.New("Error writing credentials to " + credentialFile)
  }

  return err
}

func readCredentialFile() (jsonResponse map[string]interface{}, err error) {
  if _, err := os.Stat(credentialFile); os.IsNotExist(err) { // check if file exists
    return nil, errors.New("Credential file does not exist.")
  } else {
    jsonText, err := ioutil.ReadFile(credentialFile)
    _ = json.Unmarshal(jsonText, &jsonResponse)

    return jsonResponse, err
  }
}

func RequireAuth(*cli.Context) (err error) {
  for !verifyAuth() {}
  return err
}

func Logout() {
  os.Remove(credentialFile)
}

func verifyAuth() (bool) {
  userMap, err := readCredentialFile()
  if err != nil {
    username, oauthToken, err := login()

    if err != nil {
      fmt.Println("\n"+err.Error()+"\n")
      return false
    } else {
      Email = username
      OauthToken = oauthToken
      return true
    }
  } else {
     Email = userMap["email"].(string)
     OauthToken = userMap["oauth_token"].(string)
     return true
  }
}
