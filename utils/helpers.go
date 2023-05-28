package utils

import (
  "time"
  "encoding/json"
  "net/http"
  // "regexp"
  "os"
  "fmt"

  "github.com/golang-jwt/jwt/v5"
  "github.com/vietstars/postgres-api/models"
  "golang.org/x/crypto/bcrypt"
)

func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
  response, _ := json.Marshal(payload)

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}

func RespondNotFound(w http.ResponseWriter, message string) {
  RespondJSON(w, http.StatusNotFound, map[string]string{"error": message})
}

func RespondBadRequest(w http.ResponseWriter, message string) {
  RespondJSON(w, http.StatusBadRequest, map[string]string{"error": message})
}

func RespondUnauthorized(w http.ResponseWriter, message string) {
  RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": message})
}

func RespondServerError(w http.ResponseWriter, message string) {
  RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": message})
}

// func FindErrorType(err error) error {
//   re := regexp.MustCompile(`not found.?`)
//   if re.FindString(err.Error()) != "" {
//     return RespondNotFound(err.Error())
//   }

//   return RespondServerError(err.Error())
// }

func GenerateToken(u *models.User) (string, error) {
  tokenByte := jwt.New(jwt.SigningMethodHS256)

  now := time.Now().UTC()
  duration, _ := time.ParseDuration(os.Getenv("JWT_EXPIRED_IN"))

  claims := tokenByte.Claims.(jwt.MapClaims)

  sub,_ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%+v",u.ID,u.Email)), 0);

  claims["sub"] = string(sub)
  claims["exp"] = now.Add(duration).Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()
  claims["authId"] = u.ID
  claims["authEmail"] = u.Email

  token, err := tokenByte.SignedString([]byte(os.Getenv("JWT_SECRET")))

  return token, err
}