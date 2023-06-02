package controllers

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "strconv"

  "github.com/go-playground/validator/v10"
  "github.com/gorilla/mux"
  "github.com/vietstars/postgres-api/repositories"
  "github.com/vietstars/postgres-api/utils"
)

var langValidate *validator.Validate

type LangInput struct {
  Locale string `json:"lg" validate:"required"`
  Group string `json:"group" validate:"required"`
  Key string `json:"key" validate:"required"`
  Val string `json:"val" validate:"required"`
}

type DelLangInput struct {
  Version int `json:"version" validate:"required"`
  ForceDel bool `json:"force_del" default:"false"`
}

type UpdateLangInput struct {
  Locale string `json:"lg" validate:"required"`
  Group string `json:"group" validate:"required"`
  Key string `json:"key" validate:"required"`
  Val string `json:"val" validate:"required"`
  Version int `json:"version" validate:"required"`
}

func GetAllLangs(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")
  langs, _ := repositories.GetAllLang()

  json.NewEncoder(w).Encode(langs)
}

func GetLangsByLocale(w http.ResponseWriter, r *http.Request){
   w.Header().Set("Content-Type", "application/json")

  lg := mux.Vars(r)["lg"]
  langs, err := repositories.GetLangsByLocale(lg); 
  if err != nil{
    utils.RespondNotFound(w, 
      "Lang by locale not found")

    return
  }

  json.NewEncoder(w).Encode(langs)
}

func GetLang(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  id, _ := strconv.Atoi(mux.Vars(r)["id"])
  lang, err := repositories.GetLangById(uint(id)); 
  if err != nil{
    utils.RespondNotFound(w, 
      "Lang not found")

    return
  }

  json.NewEncoder(w).Encode(lang)
}

func CreateLang(w http.ResponseWriter, r *http.Request){
   w.Header().Set("Content-Type", "application/json")

  var input LangInput 

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &input)

  langValidate = validator.New()
  err := langValidate.Struct(input)

  if err != nil {
    utils.RespondBadRequest(w,
      fmt.Sprintf("%+v\n", err))
    return 
  }

  lang, err := repositories.NewLang(input.Locale, input.Group, input.Key, input.Val)

  json.NewEncoder(w).Encode(lang) 
}

func DeleteLang(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  var input DelLangInput 

  id, _ := strconv.Atoi(mux.Vars(r)["id"])
  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &input)
  langValidate = validator.New()
  err := langValidate.Struct(input)

  if err != nil {
    utils.RespondBadRequest(w,
      fmt.Sprintf("%+v\n", err))
    return 
  }

  if _, err := repositories.DelLangById(uint(id), uint(input.Version), input.ForceDel); err != nil{
    utils.RespondNotFound(w, 
      "Lang not found")

    return
  }

  w.WriteHeader(http.StatusNoContent)
  json.NewEncoder(w).Encode(input)
}

func UpdateLang(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")
  var input UpdateLangInput 

  id, _ := strconv.Atoi(mux.Vars(r)["id"])

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &input)
  langValidate = validator.New()
  err := langValidate.Struct(input)

  if err != nil {
    utils.RespondBadRequest(w,
      fmt.Sprintf("%+v\n", err))
    return 
  }
  
  lang, err := repositories.UpdateLangById(uint(id), uint(input.Version), input.Locale, input.Group, input.Key, input.Val);

  if err != nil {
    utils.RespondNotFound(w, 
      "Lang not found")
    return
  }

  json.NewEncoder(w).Encode(lang)
}