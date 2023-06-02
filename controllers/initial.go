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

  api.Handle("/quests", TokenVerifyMiddleWare(GetAllQuests)).Methods("GET", "OPTIONS")
  api.Handle("/quest/{id}",TokenVerifyMiddleWare(GetQuest)).Methods("GET", "OPTIONS")
  api.Handle("/quest", TokenVerifyMiddleWare(CreateQuest)).Methods("POST", "OPTIONS") 
  api.Handle("/quest/{id}", TokenVerifyMiddleWare(UpdateQuest)).Methods("PUT", "OPTIONS")
  api.Handle("/quest/{id}", TokenVerifyMiddleWare(DeleteQuest)).Methods("DELETE", "OPTIONS")

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
    // w.Header().Set("Access-Control-Allow-Origin", "localhost:8008")
    // w.Header().Set("Access-Control-Allow-Credentials", "true")
    // w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    // w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")

    // if r.Method == "OPTIONS" {
    //   w.WriteHeader(204)
    //   return
    // }

    fmt.Println("Token verifier")

    authorization := r.Header.Get("Authorization") 

    cookie, err := r.Cookie("authToken")

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

    if !strings.HasPrefix(authorization, "Bearer ") && err == nil {
      authorization = cookie.Value
      // utils.RespondUnauthorized(w, 
      //   "Invalid Token")
      // log.Println("authorization: ", authorization)
      // return
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