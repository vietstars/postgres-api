package controllers

import (
  "net/http"
  "strings"
  "fmt"
  "os"
  // "log"

  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  "github.com/golang-jwt/jwt/v5"
  "github.com/vietstars/postgres-api/utils"
  "golang.org/x/crypto/bcrypt"
)

func Routers() http.Handler {
  router := mux.NewRouter()

  header := handlers.AllowedHeaders([]string{"X-Requested-With", "X-CSRF-Token", "Content-Type", "Authorization"})
  methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "HEAD", "OPTIONS"})
  origins := handlers.AllowedOrigins([]string{"localhost:8008"})

  api := router.PathPrefix("/api").Subrouter()

  router.HandleFunc("/users", GetAllUsers).Methods("GET", "OPTIONS")
  router.HandleFunc("/user", SignUp).Methods("POST", "OPTIONS") 
  router.HandleFunc("/signIn", SignIn).Methods("POST", "OPTIONS") 
  api.Handle("/user/{id}", TokenVerifyMiddleWare(GetUser)).Methods("GET", "OPTIONS")

  api.Handle("/langs", TokenVerifyMiddleWare(GetAllLangs)).Methods("GET", "OPTIONS")
  api.Handle("/langs/{lg}",TokenVerifyMiddleWare(GetLangsByLocale)).Methods("GET", "OPTIONS")
  api.Handle("/lang/{id}",TokenVerifyMiddleWare(GetLang)).Methods("GET", "OPTIONS")
  api.Handle("/lang", TokenVerifyMiddleWare(CreateLang)).Methods("POST", "OPTIONS") 
  api.Handle("/lang/{id}", TokenVerifyMiddleWare(UpdateLang)).Methods("PUT", "OPTIONS")
  api.Handle("/lang/{id}", TokenVerifyMiddleWare(DeleteLang)).Methods("DELETE", "OPTIONS")

  // return router
  return handlers.CORS(header, methods, origins)(router)
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
  var authToken string

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Token verifier")

    authorization := r.Header.Get("Authorization") 

    cookie, err := r.Cookie("authToken")

    if !strings.HasPrefix(authorization, "Bearer ") && err == nil {
      authorization = cookie.Value
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