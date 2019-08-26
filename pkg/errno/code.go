package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	ErrValidation = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase   = &Errno{Code: 20002, Message: "Database error."}
	ErrToken      = &Errno{Code: 20003, Message: "Error occurred while generating token."}
	ErrBadRequest = &Errno{Code: 20004, Message: "Error occurred while payload is not bad."}

	// user errors
	ErrEncrypt              = &Errno{Code: 20101, Message: "Error occurred while encrypting the user password."}
	ErrUserNotFound         = &Errno{Code: 20102, Message: "The user was not found."}
	ErrTokenInvalid         = &Errno{Code: 20103, Message: "The token was invalid."}
	ErrPasswordIncorrect    = &Errno{Code: 20104, Message: "The password was incorrect."}
	ErrPasswordBase64Decode = &Errno{Code: 20105, Message: "The panic from password base64 string decoding."}

	// signup error
	ErrUserSignupEmailInvalid = &Errno{Code: 20201, Message: "The email from payload is invalid."}
	ErrUserExisted            = &Errno{Code: 20202, Message: "The user has existed."}

	// signin error
	ErrUserPasswordIncorrect = &Errno{Code: 20301, Message: "Password incorrect."}

	// mail error
	ErrMailSend             = &Errno{Code: 20401, Message: "Mail send failed."}
	ErrGenerateCaptchaToken = &Errno{Code: 20402, Message: "Generating captcha token failed."}

	// orm error
	ErrUserCreate = &Errno{Code: 30001, Message: "The (*UserModel)Create() method error."}
	ErrUserUpdate = &Errno{Code: 30002, Message: "The (*UserModel)Update() method error."}
)
