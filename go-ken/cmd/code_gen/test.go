package main

type ErrCode int

var ()

const (
	//:404:User not found
	ErrUserNotFound = 110001
	//:400:User already exist
	ErrUserAlreadyExist = 110002
	//:400:Secret reach the max count
	ErrReachMaxCount = 110101
	//:404:Secret not found
	ErrSecretNotFound = 110102
	//:200:OK
	ErrSuccess = 100001
	//:500:Internal server error
	ErrUnknown = 100002
	//:400:Error occurred while binding the request body to the struct
	ErrBind = 100003
	//:400:Validation failed
	ErrValidation = 100004
	//:401:Token invalid
	ErrTokenInvalid = 100005
	//:500:Database error
	ErrDatabase = 100101
	//:401:Error occurred while encrypting the user password
	ErrEncrypt = 100201
	//:401:Signature is invalid
	ErrSignatureInvalid = 100202
	//:401:Token expired
	ErrExpired = 100203
	//:401:Invalid authorization header
	ErrInvalidAuthHeader = 100204
	//:401:The Authorization header was empty
	ErrMissingHeader = 100205
	//:401:Token expired
	ErrorExpired = 100206
	//:401:Password was incorrect
	ErrPasswordIncorrect = 100207
	//:403:Permission denied
	ErrPermissionDenied = 100208
	//:500:Encoding failed due to an error with the data
	ErrEncodingFailed = 100301
	//:500:Decoding failed due to an error with the data
	ErrDecodingFailed = 100302
	//:500:Data is not valid JSON
	ErrInvalidJSON = 100303
	//:500:JSON data could not be encoded
	ErrEncodingJSON = 100304
	//:500:JSON data could not be decoded
	ErrDecodingJSON = 100305
	//:500:Data is not valid Yaml
	ErrInvalidYaml = 100306
	//:500:Yaml data could not be encoded
	ErrEncodingYaml = 100307
	//:500:Yaml data could not be decoded
	ErrDecodingYaml = 100308
)
