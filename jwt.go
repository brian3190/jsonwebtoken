package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "time"

  jwt "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
)

// using asymmetric crypto/RSA keys
// location of the files used
const (
    privKeyPath = "keys/app.rsa"      // openssl genrsa -out app.rsa 1024
    pubKeyPath  = "keys/app.rsa.pub"  // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

// verify key and sign key
var (
    verifyKey, signKey []byte
)

//struct User for parsing login credentials
type User struct {
    UserName string 'json:"username"'
    Password string 'json:"password"'
}

// read the key files before starting http handlers
func init() {
    var err error

    signKey, err = ioutil.ReadFile(privKeyPath)
    if err != nil {
        log.Fatal("Error reading private key")
        return
    }

    verifyKey, err = ioutil.ReadFile(pubKeyPath)
    if err != nil {
        log.Fatal("Error reading private key")
        return
    }
}

// reads the llgin credentials, checks them and creates JWT
func loginHandler(w http.ResponseWriter, r *http.Request) {
  var user User
  //decode into User struct
  err := json.NewDecoder(r.Body).Decode(&user)
  if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintln(w, "Error in request body")
      return
  }
  // validate user credentials
  if user.UserName != "user" && user.Password != "pass" {
      w.WriteHeader(http.StatusForbidden)
      fmt.Fprintln(w, "Wrong info")
      return
  }

  // create a signer for rsa 256
  t := jwt.New(jwt.GetSigningMethod("RS256"))

  // set our claims
  t.Claims["iss"] = "admin"
  t.Claims["CustomUserInfo"] = struct {
      Name string
      Role string
  }{user.UserName, "Member"}

  // set the expire time
  t.Claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
  tokenString, err := t.SignedString(signKey)
  if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintln(w, "Sorry, error while Signing Token!")
      log.Printf("Token Signing error: %v\n", err)
      return
  }
  response := Token{tokenString}
  jsonResponse(response, w)
}

// only accessible with a valid tokenString
func authHandler(w http.ResponseWriter, r *http.Request) {
    // validate the token
    token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
        // since only use one private key to sign the tokens,
        // also only use its public counterpart to verify
        return verifyKey, nil
    })

    if err != nil {
        switch err.(type) {

        case *jwt.ValidationError: // something was wrong during the validation
            vErr := err.(*jwt.ValidationError)

            switch vErr.Errors {
            case jwt.ValidationErrorExpired:
              w.WriteHeader(http.StatusUnauthorized)
              fmt.Fprintln(w, "Token Expired, get a new one.")
              return

            default:
              w.WriteHeader(http.StatusInternalServerError)
              fmt.Fprintln(w, "Error while Parsing Token!")
              log.Printf("ValidationError error: %+v\n", vErr.Errors)
              return
            }

          default: // something else went wrong
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintln(w, "Error while Parsing Token!")
            log.Printf("Token parse error: %v\n", err)
            return
        }
    }
    if token.Valid {
      response := Response{"Authorized to the system"}
      jsonResponse(response, w)
    } else {
      response := Response{"Invalid token"}
      jsonResponse(response, w)
    }
}

type Response struct {
  
}
