package models

type User struct {
  Id        uint64 `json:"id" gorm:"primary_key"`
  Username  string `json:"username"`
  Password  string `json:"password"`
}