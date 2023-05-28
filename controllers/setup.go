package controllers

import (
  "net/http"
  "strings"
  "fmt"
  "os"

  "github.com/gorilla/mux"
  "github.com/golang-jwt/jwt/v5"
  "github.com/vietstars/postgres-api/utils"
  "golang.org/x/crypto/bcrypt"
)

func New() http.Handler {
  router := mux.NewRouter()

  router.HandleFunc("/users", GetAllUsers).Methods("GET")
  router.HandleFunc("/user/{id}", TokenVerifyMiddleWare(GetUser)).Methods("GET")
  router.HandleFunc("/user", SignUp).Methods("POST") 
  router.HandleFunc("/signIn", SignIn).Methods("POST") 

  router.HandleFunc("/quests", TokenVerifyMiddleWare(GetAllQuests)).Methods("GET")
  router.HandleFunc("/quest/{id}",TokenVerifyMiddleWare(GetQuest)).Methods("GET")
  router.HandleFunc("/quest", TokenVerifyMiddleWare(CreateQuest)).Methods("POST") 
  router.HandleFunc("/quest/{id}", TokenVerifyMiddleWare(UpdateQuest)).Methods("PUT")
  router.HandleFunc("/quest/{id}", TokenVerifyMiddleWare(DeleteQuest)).Methods("DELETE")

  router.HandleFunc("/langs", TokenVerifyMiddleWare(GetAllLangs)).Methods("GET")
  router.HandleFunc("/langs/{lg}",TokenVerifyMiddleWare(GetLangsByLocale)).Methods("GET")
  router.HandleFunc("/lang/{id}",TokenVerifyMiddleWare(GetLang)).Methods("GET")
  router.HandleFunc("/lang", TokenVerifyMiddleWare(CreateLang)).Methods("POST") 
  router.HandleFunc("/lang/{id}", TokenVerifyMiddleWare(DeleteLang)).Methods("DELETE")

  return router
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
  var authToken string

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    fmt.Println("Token verifier")
    authorization := r.Header.Get("Authorization") 

    if !strings.HasPrefix(authorization, "Bearer ") {
      utils.RespondUnauthorized(w, 
        "Invalid Token")

      return
    }

    authToken = strings.TrimPrefix(authorization, "Bearer ")

    if authToken == "" {
      utils.RespondUnauthorized(w, 
        "You are not logged in")

      return
    }

    tokenByte, err := jwt.Parse(authToken, func(jwtToken *jwt.Token) (interface{}, error) {
      if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {

        return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
      }

      return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
      utils.RespondUnauthorized(w,  
        fmt.Sprintf("invalidate token: %v", err))

      return
    }

    claims, ok := tokenByte.Claims.(jwt.MapClaims)

    if !ok || !tokenByte.Valid {
      utils.RespondUnauthorized(w,
        "Invalid token claim")

      return
    }

    err = bcrypt.CompareHashAndPassword([]byte(fmt.Sprintf("%v", claims["sub"])), 
      []byte(fmt.Sprintf("%+v", claims["authId"], claims["authEmail"])))

    if err != nil {
      utils.RespondUnauthorized(w,
        "the user belonging to this token no logger exists")

      return
    }

    next.ServeHTTP(w, r)
  })
}