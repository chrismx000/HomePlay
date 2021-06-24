package model

import (
	"gorm.io/gorm"
)

type VideoInfo struct {
	gorm.Model `json:"-"`
	VideoNo    string `json:"video" gorm:"default:'';column:videoNo"` //剧目编号
	Name       string `json:"name" gorm:"default:'';column:name"`     //名称
	ImgSrc     string `json:"imgSrc" gorm:"default:'';column:imgSrc"` //图片源地址
	Last       string `json:"last" gorm:"default:'';column:last"`     //最新集数
	Url        string `json:"url" gorm:"default:'';column:url"`       //地址
}

type VideoDetail struct {
	gorm.Model `json:"-"`
	VideoNo    string `json:"infoID" gorm:"default:'';column:videoNo"`  //剧目编号
	SourceNo   int    `json:"sourceNo" gorm:"column:sourceNo"`          //播放源编号
	Episode    string `json:"episode" gorm:"default:'';column:episode"` //播放话数
	Url        string `json:"url" gorm:"default:'';column:url"`         //播放地址
}
