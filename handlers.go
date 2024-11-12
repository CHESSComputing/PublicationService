package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	srvConfig "github.com/CHESSComputing/golib/config"
	services "github.com/CHESSComputing/golib/services"
	zenodo "github.com/CHESSComputing/golib/zenodo"
	"github.com/gin-gonic/gin"
)

/*
Implementation is based on few resources:
- https://felipecrp.com/2021/01/01/uploading-to-zenodo-through-api.html
- Zenodo REST API
  https://developers.zenodo.org/
- snd, some discussion about zenodo APIs can be found here:
  https://github.com/zenodo/zenodo/issues/2168
*/

// DocParams defines binding URI parameters
type DocParams struct {
	Id       int64  `uri:"id"`
	Bucket   string `uri:"bucket"`
	FileName string `uri:"file"`
}

// DocsHandler services / end-point
func DocsHandler(c *gin.Context) {
	/*
	 curl 'https://zenodo.org/api/deposit/depositions?access_token=<KEY>'
	 curl 'https://zenodo.org/api/deposit/depositions/<123>?access_token=<KEY>'
	*/
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions?access_token=%s", zurl, token)
	var doc DocParams
	if err := c.ShouldBindUri(&doc); err == nil {
		if doc.Id != 0 {
			rurl = fmt.Sprintf("%s/deposit/depositions/%d?access_token=%s", zurl, doc.Id, token)
		}
	}
	if Verbose > 0 {
		log.Println("request", rurl)
	}
	resp, err := _httpReadRequest.Get(rurl)
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.HttpRequestError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if doc.Id == 0 {
		var results []map[string]any
		if err := dec.Decode(&results); err != nil {
			rec := services.Response("Publication", http.StatusBadRequest, services.DecodeError, err)
			c.JSON(http.StatusBadRequest, rec)
			return
		}
		c.JSON(http.StatusOK, results)
		return
	}
	var record map[string]any
	if err := dec.Decode(&record); err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.DecodeError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	c.JSON(resp.StatusCode, record)
}

// CreateHandler services /create end-point
func CreateHandler(c *gin.Context) {
	/*
	 curl --request POST 'https://zenodo.org/api/deposit/depositions?access_token=<KEY>' \
	 --header 'Content-Type: application/json'  \
	 --data-raw '{}'
	*/
	// create new deposit
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions?access_token=%s", zurl, token)
	if Verbose > 0 {
		log.Println("request", rurl)
	}
	resp, err := _httpWriteRequest.Post(rurl, "application/json", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.HttpRequestError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	log.Println("POST", rurl, string(data), err)
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.ReaderError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	var response zenodo.Response
	err = json.Unmarshal(data, &response)
	if Verbose > 1 {
		log.Println("repsonse struct", response)
	}
	if response.Status > 0 && response.Status != 200 {
		msg := fmt.Sprintf("Zenodo response %s", string(data))
		rec := services.Response("Publication", http.StatusBadRequest, services.ZenodoError, errors.New(msg))
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	if response.Status == 0 {
		c.Data(http.StatusOK, "application/json", data)
		return
	}
	c.JSON(http.StatusOK, response)
}

// AddHandler services /add end-point
func AddHandler(c *gin.Context) {
	/*
		curl --upload-file readme.md --request PUT
		'https://zenodo.org/api/files/50b47f75-c97d-47c6-af11-caa6e967c1d5/readme.md?access_token=<KEY>'
	*/
	var doc DocParams
	if err := c.ShouldBindUri(&doc); err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.BindError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}

	// create new deposit
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/files/%s/%s?access_token=%s", zurl, doc.Bucket, doc.FileName, token)
	if Verbose > 0 {
		log.Println("request", rurl)
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.ReaderError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}

	// place HTTP request to zenodo upstream server
	req, err := http.NewRequest("PUT", rurl, bytes.NewReader(body))
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.HttpRequestError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	log.Println("PUT", rurl, string(data), err)
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.ReaderError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	c.JSON(resp.StatusCode, string(data))
}

// UpdateHandler services /add end-point
func UpdateHandler(c *gin.Context) {
	/*
	   # add mandatory metadata to our publication
	   curl -v -X PUT "https://zenodo.org/api/deposit/depositions/<ID>?access_token=<TOKEN>" \
	           -H "Content-type: application/json" -d@meta1.json

	   	{
	   	    "metadata": {
	   	        "publication_type": "article",
	   	        "upload_type":"publication",
	   	        "description":"This is a test",
	   	        "keywords": ["bla", "foo"],
	   	        "title":"Test"
	   	    }
	   	}
	*/
	var doc DocParams
	if err := c.ShouldBindUri(&doc); err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.BindError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions/%d?access_token=%s", zurl, doc.Id, token)

	// read payload
	defer c.Request.Body.Close()
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.HttpRequestError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}

	// place HTTP request to zenodo upstream server
	req, err := http.NewRequest("PUT", rurl, bytes.NewReader(data))
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.HttpRequestError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	log.Println("PUT", rurl, string(data), err)
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.ReaderError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	c.JSON(resp.StatusCode, string(data))
}

// PublishHandler services /add end-point
func PublishHandler(c *gin.Context) {
	// curl -v -X POST "https://zenodo.org/api/deposit/depositions/<ID>/actions/publish?access_token=<TOKEN>"
	var doc DocParams
	if err := c.ShouldBindUri(&doc); err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.BindError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions/%d/actions/publish?access_token=%s", zurl, doc.Id, token)

	// place HTTP request to zenodo upstream server
	req, err := http.NewRequest("POST", rurl, nil)
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.HttpRequestError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	log.Println("POST", rurl, string(data), err)
	if err != nil {
		rec := services.Response("Publication", http.StatusBadRequest, services.ReaderError, err)
		c.JSON(http.StatusBadRequest, rec)
		return
	}
	c.JSON(resp.StatusCode, string(data))
}
