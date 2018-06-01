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
