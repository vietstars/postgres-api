package controllers

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "net/http"

  "github.com/go-playground/validator/v10"
  "github.com/gorilla/mux"
  "github.com/vietstars/postgres-api/models"
  "github.com/vietstars/postgres-api/utils"
)

var questValidate *validator.Validate

type QuestInput struct {
  Title string `json:"title" validate:"required"`
  Description string `json:"description" validate:"required"`
  Reward int `json:"reward" validate:"required"`
}

type DeleteQuestInput struct {
  Version int `json:"version" validate:"required"`
}

func CreateQuest(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  var input QuestInput 

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &input)

  questValidate = validator.New()
  err := questValidate.Struct(input)

  if err != nil {
    utils.RespondBadRequest(w,
      fmt.Sprintf("%+v\n", err))
    return 
  }

  quest, err := models.NewQuest(input.Title, input.Description, input.Reward)

  json.NewEncoder(w).Encode(quest) 

}


func GetAllQuests(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  var quests []models.Quest
  models.DB.Find(&quests)

  json.NewEncoder(w).Encode(quests)
}


func GetQuest(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  id := mux.Vars(r)["id"]
  var quest models.Quest

  if err := models.DB.Where("id = ?", id).First(&quest).Error; err != nil{
    utils.RespondNotFound(w, 
      "Quest not found")
    return
  }

  json.NewEncoder(w).Encode(quest)
}


func DeleteQuest(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  var input DeleteQuestInput 
  var quest models.Quest

  id := mux.Vars(r)["id"]

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &input)

  questValidate = validator.New()
  err := questValidate.Struct(input)

  if err != nil {
    utils.RespondBadRequest(w,
      fmt.Sprintf("%+v\n", err))
    return 
  }

  tx := models.DB.Begin()
  if err := tx.Error; err != nil {
      return
  }

  if err := models.DB.Where("id = ? And version = ?", id, input.Version).First(&quest).Error; err != nil{
    utils.RespondNotFound(w, 
      "Quest not found")
    tx.Rollback()
    return
  }

  // models.DB.Delete(&quest)
  tx.Delete(&quest)

  tx.Commit()

  // // Xóa bản ghi theo điều kiện
  // tx.Where("full_name = ?", "bob").Delete(&model.Student{})

  // // Xóa theo khóa chính
  // tx.Delete(&model.Student{}, "2")

  // // Xóa theo danh sách khóa chính
  // tx.Delete(&model.Student{}, []string{"3", "4"})

  w.WriteHeader(http.StatusNoContent)
  json.NewEncoder(w).Encode(quest)
}


func UpdateQuest(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  id := mux.Vars(r)["id"]
  var quest models.Quest

  if err := models.DB.Where("id = ?", id).First(&quest).Error; err != nil{
    utils.RespondNotFound(w, 
      "Quest not found")
    return
  }

  var input QuestInput 

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &input)

  userValidate = validator.New()
  err := userValidate.Struct(input)

  if err != nil {
    utils.RespondBadRequest(w, 
      "Validation Error")
    return 
  }
  
  quest.Title = input.Title
  quest.Description = input.Description
  quest.Reward = input.Reward

  models.DB.Save(&quest)

  json.NewEncoder(w).Encode(quest)
}
