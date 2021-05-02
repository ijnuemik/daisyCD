package models

type Application struct {
	ReleaseName	  	  string `json:"releaseName" gorm:"primary_key"`
	ReleaseNamespace  string `json:"releaseNamespace"`
	ChartPath         string `json:"chartPath"`
}