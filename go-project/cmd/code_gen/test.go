package main

const (
	//BadRequest:404:User not found
	CodeUserNotFound = 110001
	//:400:User already exist
	CodeUserAlreadyExist = 110002
	//:400:Secret reach the max count
	CodeReachMaxCount = 110101
	//:404:Secret not found
	CodeSecretNotFound = 110102
	//:200:OK
	CodeSuccess = 100001
	//:500:Internal server error
	CodeUnknown = 100002
	//:400:error occurred while binding the request body to the struct
	CodeBind = 100003
	//:400:Validation failed
	CodeValidation = 100004
	//:401:Token invalid
	CodeTokenInvalid = 100005
	//:500:Database Error
	CodeDatabase = 100101
	//:401:Error occurred while encrypting the user password
	CodeEncrypt = 100201
	//:401:Signature is invalid
	CodeSignatureInvalid = 100202
	//:401:Invalid authorization header
	CodeInvalidAuthHeader = 100204
	//:401:The Authorization header was empty
	CodeMissingHeader = 100205
	//:401:Token expired
	CodeExpired = 100206
	//:401:Password was incorrect
	CodePasswordIncorrect = 100207
	//:403:Permission denied
	CodePermissionDenied = 100208
	//:500:Encoding failed due to an Error with the data
	CodeEncodingFailed = 100301
	//:500:Decoding failed due to an Error with the data
	CodeDecodingFailed = 100302
	//:500:Data is not valid JSON
	CodeInvalidJSON = 100303
	//:500:JSON data could not be encoded
	CodeEncodingJSON = 100304
	//:500:JSON data could not be decoded
	CodeDecodingJSON = 100305
	//:500:Data is not valid Yaml
	CodeInvalidYaml = 100306
	//:500:Yaml data could not be encoded
	CodeEncodingYaml = 100307
	//:500:Yaml data could not be decoded
	CodeDecodingYaml = 100308
)
