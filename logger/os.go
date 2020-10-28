package logger

import "os"

var GapiHost = os.Getenv("GAPI_HOST")
var GapiClientId = os.Getenv("GAPI_CLIENT_ID")
var GapiClientSecret = os.Getenv("GAPI_CLIENT_SECRET")
var GapiUsername = os.Getenv("GAPI_USERNAME")
var GapiPassword = os.Getenv("GAPI_PASSWORD")
