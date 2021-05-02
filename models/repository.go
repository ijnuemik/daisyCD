package models

type Repository struct {
  UserId      uint64 `json:"userid" gorm:"primary_key"`
  Url         string `json:"url" gorm:"primary_key"`
  AccessUser  string `json:"accessuser"`
  AccessToken	string `json:"accesstoken"`
  Status      uint64 `json:"status"`
}