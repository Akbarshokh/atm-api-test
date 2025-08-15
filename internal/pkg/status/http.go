package status

var (
	NoError              = 0
	ErrorCodeValidation  = -10
	ErrorAuthorization   = -11
	ErrorInvalidOTP      = -12
	ErrorInvalidToken    = -13
	ErrorCodeDB          = -14
	ErrorInvalidPassword = -15
)

var (
	Success = "Success"
	Failure = "Failure"
)
