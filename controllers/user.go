package controllers

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "strconv"
  "time"
  "os"

  "github.com/go-playground/validator/v10"
  "github.com/gorilla/mux"
  "github.com/vietstars/postgres-api/dto"
  "github.com/vietstars/postgres-api/repositories"
  "github.com/vietstars/postgres-api/utils"
  "golang.org/x/crypto/bcrypt"
)

var userValidate *validator.Validate

func GetAllUsers(w http.ResponseWriter, r *http.Request){
  users, _ := repositories.GetAllUsers()

  // for i,user := range users {
  //   if user.Version == 0 {
  //     users[i].Version = 1
  //   }
  // }

  utils.RespondJSONData(w, users)
}

func GetUser(w http.ResponseWriter, r *http.Request){

  id, _ := strconv.Atoi(mux.Vars(r)["id"])
  user, err := repositories.GetUserById(uint(id)); 
  if err != nil{
    utils.RespondNotFound(w, 
      "User not found")

    return
  }

  utils.RespondJSONData(w, user)
}

func SignUp(w http.ResponseWriter, r *http.Request){
  var inputs dto.UserNew 

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &inputs)

  userValidate = validator.New()
  err := userValidate.Struct(inputs)

  if err != nil {
    utils.RespondBadRequest(w, 
      fmt.Sprintf("%+v\n", err))
    return 
  }

  user, err := repositories.NewUser(inputs)
   if err != nil {
    utils.RespondBadRequest(w, 
      "User does exist")
    return 
  }

  utils.RespondJSONData(w, user) 
}

func SignIn(w http.ResponseWriter, r *http.Request){
  var inputs dto.UserSignIn 

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &inputs)

  userValidate = validator.New()
  err := userValidate.Struct(inputs)

  if err != nil {
    utils.RespondBadRequest(w, 
      fmt.Sprintf("%+v\n", err))
    return 
  }

  user, err := repositories.GetUserByEmail(inputs.Email)

  if err != nil {
    utils.RespondNotFound(w, 
      "User does not exist")
    return 
  }

  hashedPassword := user.Password
  err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputs.Password))

  if err != nil {
    utils.RespondBadRequest(w, 
      "Password does not match")
    return
  }

  token, err := utils.GenerateToken(user)
  if err != nil {
    utils.RespondServerError(w, 
      "Error while generating token")
    return
  }

  auth := dto.Auth{user, token}
  maxAge, _ := strconv.Atoi(os.Getenv("JWT_MAXAGE")) 
  duration, _ := time.ParseDuration(os.Getenv("JWT_EXPIRED_IN"))

  utils.SetCookie(w, "authToken", token, maxAge * 60, time.Now().UTC().Add(duration))

  utils.RespondJSONData(w, auth)
}
