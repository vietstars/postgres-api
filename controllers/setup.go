package controllers

import (
  "net/http"

  "github.com/gorilla/mux"
)

func New() http.Handler {
  router := mux.NewRouter()


  router.HandleFunc("/users", TokenVerifyMiddleWare(GetAllUsers)).Methods("GET")
  router.HandleFunc("/user/{id}", TokenVerifyMiddleWare(GetUser)).Methods("GET")
  router.HandleFunc("/user", SignUp).Methods("POST") 
  router.HandleFunc("/signIn", SignIn).Methods("POST") 

  router.HandleFunc("/quests", TokenVerifyMiddleWare(GetAllQuests)).Methods("GET")
  router.HandleFunc("/quest/{id}",TokenVerifyMiddleWare(GetQuest)).Methods("GET")
  router.HandleFunc("/quest", TokenVerifyMiddleWare(CreateQuest)).Methods("POST") 
  router.HandleFunc("/quest/{id}", TokenVerifyMiddleWare(UpdateQuest)).Methods("PUT")
  router.HandleFunc("/quest/{id}/{version}", TokenVerifyMiddleWare(DeleteQuest)).Methods("DELETE")

  return router
}
