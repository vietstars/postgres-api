package utils

import (
  "time"
  "encoding/json"
  "net/http"
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

func RespondJSONData(w http.ResponseWriter, payload interface{}) {
  response, _ := json.Marshal(payload)

  w.Header().Set("Content-Type", "application/json")
  w.Write(response)
}

func RespondNothing(w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusNoContent)
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

func GenerateToken(u *models.UserEntity) (string, error) {
  tokenByte := jwt.New(jwt.SigningMethodHS256)

  now := time.Now().UTC()
  duration, _ := time.ParseDuration(os.Getenv("JWT_EXPIRED_IN"))

  claims := tokenByte.Claims.(jwt.MapClaims)

  sub,_ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%+v",u.UserID,u.Email)), 0);

  claims["sub"] = string(sub)
  claims["exp"] = now.Add(duration).Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()
  claims["authId"] = u.UserID
  claims["authEmail"] = u.Email

  token, err := tokenByte.SignedString([]byte(os.Getenv("JWT_SECRET")))

  return token, err
}

func SetCookie(w http.ResponseWriter, key string, val string, maxAge int, duration time.Time) {
  cookie := http.Cookie{
    Name:     key,
    Value:    val,
    Path:     "/",
    MaxAge:   maxAge,
    Expires:  duration,
    Secure:   true,
    HttpOnly: true,
    SameSite: http.SameSiteNoneMode,
  }

  // if err != nil {
  //   switch {
  //   case errors.Is(err, http.ErrNoCookie):
  //     http.Error(w, "cookie not found", http.StatusBadRequest)
  //   default:
  //     log.Println(err)
  //     http.Error(w, "server error", http.StatusInternalServerError)
  //   }
  //   return
  // }

  http.SetCookie(w, &cookie)
}
