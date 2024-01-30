package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	srvConfig "github.com/CHESSComputing/golib/config"
	services "github.com/CHESSComputing/golib/services"
	"github.com/gin-gonic/gin"
)

func DataHandler(c *gin.Context) {
	zurl := srvConfig.Config.Publication.Zenodo.URL
	token := srvConfig.Config.Publication.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions?access_token=%s", zurl, token)
	if Verbose > 0 {
		log.Println("request", rurl)
	}
	resp, err := _httpReadRequest.Get(rurl)
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.ReaderError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	defer resp.Body.Close()
	//     data, err := io.ReadAll(resp.Body)
	//     if err != nil {
	//         rec := services.Response("Publication", http.StatusBadRequest, services.ReaderError, err)
	//         c.JSON(http.StatusBadRequest, rec)
	//         return
	//     }
	dec := json.NewDecoder(resp.Body)
	var results []map[string]any
	if err := dec.Decode(&results); err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.ReaderError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	c.JSON(http.StatusOK, results)
}
func PublishHandler(c *gin.Context) {
}
