package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/twoonefour/alist-auth/common"
	"github.com/twoonefour/alist-auth/utils"
)

var (
	dropBoxAppSecret string
	dropBoxAppId     string
)

const (
	DropboxAuthUrl = "https://api.dropboxapi.com/oauth2/token"
)

type AccessTokenReqData struct {
	Code         string `json:"code" form:"code"`
	ClientId     string `json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
	GrantType    string `json:"grant_type" form:"grant_type"`
	RedirectUri  string `json:"redirect_uri" form:"redirect_uri"`
}

type AccessTokenResData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	AccountId   string `json:"account_id"`
	Uid         string `json:"uid"`
}

func getDropBoxToken(g *gin.Context) {
	var accessTokenReqData AccessTokenReqData
	err := g.Bind(&accessTokenReqData)
	if err != nil {
		common.Error(g, err)
		return
	}

	if accessTokenReqData.GrantType == "refresh_token" {
		refreshDropBoxToken(g, accessTokenReqData)
		return
	}

	client := utils.RestyClient.R()
	client.SetFormData(map[string]string{
		"code": accessTokenReqData.Code,
		"client_id": func() string {
			if accessTokenReqData.ClientSecret == "" {
				return dropBoxAppId
			}
			return accessTokenReqData.ClientId
		}(),
		"client_secret": func() string {
			if accessTokenReqData.ClientSecret == "" {
				return dropBoxAppSecret
			}
			return accessTokenReqData.ClientSecret
		}(),
		"grant_type":   accessTokenReqData.GrantType,
		"redirect_uri": accessTokenReqData.RedirectUri,
	}).SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "application/json",
	})
	res, err := client.Execute("POST", DropboxAuthUrl)
	if err != nil {
		common.Error(g, err)
		return
	}
	if err := common.JsonBytes(g, res.Body()); err != nil {
		utils.GetLogger(g).Error(err)
		return
	}
}

func refreshDropBoxToken(g *gin.Context, accessTokenReqData AccessTokenReqData) {
	if accessTokenReqData.GrantType != "refresh_token" {
		common.ErrorStr(g, "Invalid grant type")
		return
	}
	client := utils.RestyClient.R()
	client.SetFormData(map[string]string{
		"refresh_token": accessTokenReqData.RefreshToken,
		"client_id": func() string {
			if accessTokenReqData.ClientId == "" {
				return dropBoxAppId
			}
			return accessTokenReqData.ClientId
		}(),
		"client_secret": func() string {
			if accessTokenReqData.ClientSecret == "" {
				return dropBoxAppSecret
			}
			return accessTokenReqData.ClientSecret
		}(),
		"grant_type": accessTokenReqData.GrantType,
	}).SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Accept":       "application/json",
	})
	res, err := client.Execute("POST", DropboxAuthUrl)
	if err != nil {
		common.Error(g, err)
		return
	}
	common.JsonBytes(g, res.Body())
}
