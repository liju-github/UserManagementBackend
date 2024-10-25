package models

const (
	InvalidInput                       = "Invalid input"
	UserAlreadyExists                  = "User already exists"
	UserDoesntExist                    = "user doesnt exist"
	LoginSuccessful                    = "Login successful"
	UserIsBlocked                      = "User is Blocked"
	LogoutSuccessful                   = "Logout successful"
	EmailVerifiedSuccessfully          = "Email verified successfully"
	VerificationEmailResent            = "Verification email resent"
	PasswordResetEmailSent             = "Password reset email sent"
	PasswordResetSuccessfully          = "Password reset successfully"
	ProfileUpdatedSuccessfully         = "Profile updated successfully"
	ProfilePictureUploadedSuccessfully = "Profile picture uploaded successfully"
	SignupSuccessful                   = "User signed up successfully!"
	MinPasswordLength                  = 8
	MaxPasswordLength                  = 72
	ErrRequiredFieldsEmpty             = "required fields cannot be empty"
	ErrInvalidEmailFormat              = "invalid email format"
	ErrNegativeAge                     = "age must be positive"
	ErrPasswordComplexity              = "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
	ErrPasswordLength                  = "password must be between %d and %d characters"
	InvalidID = "Unauthorized or invalid user ID"
)
