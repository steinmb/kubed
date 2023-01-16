package main

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
)

// JWTToken structure
type JWTToken struct {
	Token string `json:"token"`
}

type ca struct {
	Cert string `json:"cert"`
}

func getJWTToken(accessToken string, issuerURL string) (string, error) {
	var jwt JWTToken

	resp, _, err := gorequest.New().Get(issuerURL).
		Set("Authorization", "Bearer "+accessToken).
		EndStruct(&jwt)

	if err != nil {
		log.Warn("Failed in fetching JWT Token ", err)
		return "", err[0]
	}

	if resp != nil && resp.StatusCode != 201 {
		log.Warn("Failed in fetching JWT Token, responsecode: ", resp.StatusCode)
		return "", errors.New("failed in fetching JWT Token")
	}

	return jwt.Token, nil
}

func getCACert(issuerURL string) ([]byte, error) {
	var caInstance ca

	resp, _, err := gorequest.New().Get(issuerURL + "/ca").
		EndStruct(&caInstance)

	if err != nil {
		log.Warn("Failed in fetching CA certificate ", err)
		return nil, err[0]
	}

	if resp != nil && resp.StatusCode != 200 {
		log.Warn("Failed in fetching CA certificate, responsecode: ", resp.StatusCode)
		return nil, errors.New("failed in fetching CA certificate")
	}
	return []byte(caInstance.Cert), nil
}
