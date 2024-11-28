package utils

import "regexp"

func ValidateEmail(email string) bool {
  const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
  re := regexp.MustCompile(emailRegex)
  return re.MatchString(email)
}

func ValidateUsername(username string) bool {
  // TODO: Validate Username isn't already taken
  return len(username) > 3
}

// Validates a user's Password strength based on the following parameters:
//   - Length must be at least 8 characters long.
//   - Must contain at least 1 number
//   - Must contain at least 1 special character
func ValidatePassword(password string) bool {
  const lengthRegex      = `.{8,}`
  const numberRegex      = `[0-9]`
  const specialCharRegex = `[!@#$%^&*(),.?":{}|<>]`

  lengthCheck  := regexp.MustCompile(lengthRegex)
  numberCheck  := regexp.MustCompile(numberRegex)
  specialCheck := regexp.MustCompile(specialCharRegex)

  return lengthCheck.MatchString(password) &&
  numberCheck.MatchString(password) &&
  specialCheck.MatchString(password)
}
