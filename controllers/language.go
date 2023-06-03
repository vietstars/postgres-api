package controllers

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "strconv"

  "github.com/go-playground/validator/v10"
  "github.com/gorilla/mux"
  "github.com/vietstars/postgres-api/dto"
  "github.com/vietstars/postgres-api/repositories"
  "github.com/vietstars/postgres-api/utils"
)

var langValidate *validator.Validate

func GetAllLangs(w http.ResponseWriter, r *http.Request){
  langs, _ := repositories.GetAllLang()

  utils.RespondJSONData(w, langs)
}

func GetLangsByLocale(w http.ResponseWriter, r *http.Request){

  lg := mux.Vars(r)["lg"]
  langs, err := repositories.GetLangsByLocale(lg); 
  if err != nil{
    utils.RespondNotFound(w, 
      "Lang by locale not found")

    return
  }

  utils.RespondJSONData(w, langs)
}

func GetLang(w http.ResponseWriter, r *http.Request){

  id, _ := strconv.Atoi(mux.Vars(r)["id"])
  lang, err := repositories.GetLangById(uint(id)); 
  if err != nil{
    utils.RespondNotFound(w, 
      "Lang not found")

    return
  }

  utils.RespondJSONData(w, lang)
}

func CreateLang(w http.ResponseWriter, r *http.Request){

  var input dto.LangNew 

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &input)

  langValidate = validator.New()
  err := langValidate.Struct(input)

  if err != nil {
    utils.RespondBadRequest(w,
      fmt.Sprintf("%+v\n", err))
    return 
  }

  lang, err := repositories.NewLang(input)

  utils.RespondJSONData(w, lang) 
}

func DeleteLang(w http.ResponseWriter, r *http.Request){

  var input dto.LangDel 

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

  if _, err := repositories.DelLangById(uint(id), input); err != nil{
    utils.RespondNotFound(w, 
      "Lang not found")

    return
  }

  // w.Header().Set("Content-Type", "application/json")
  // w.WriteHeader(http.StatusNoContent)
  // json.NewEncoder(w).Encode(input)
  utils.RespondNothing(w)
}

func UpdateLang(w http.ResponseWriter, r *http.Request){
  var input dto.LangEdit 

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
  
  lang, err := repositories.UpdateLangById(uint(id), input);

  if err != nil {
    utils.RespondNotFound(w, 
      "Lang not found")
    return
  }

  utils.RespondJSONData(w, lang)
}