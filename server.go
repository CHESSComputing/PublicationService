package main

// server module
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	srvConfig "github.com/CHESSComputing/golib/config"
	server "github.com/CHESSComputing/golib/server"
	services "github.com/CHESSComputing/golib/services"
	"github.com/gin-gonic/gin"
)

var _httpReadRequest, _httpWriteRequest *services.HttpRequest
var Verbose int

// helper function to setup our router
func setupRouter() *gin.Engine {
	routes := []server.Route{
		server.Route{Method: "GET", Path: "/docs", Handler: DocsHandler, Authorized: false},
		server.Route{Method: "GET", Path: "/docs/:id", Handler: DocsHandler, Authorized: false},
		server.Route{Method: "POST", Path: "/create", Handler: CreateHandler, Authorized: true, Scope: "write"},
		//         server.Route{Method: "PUT", Path: "/add/:bucket/:file", Handler: AddHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "PUT", Path: "/add/:bucket/:file", Handler: AddHandler},
		server.Route{Method: "POST", Path: "/publish", Handler: PublishHandler, Authorized: true, Scope: "write"},
	}
	r := server.Router(routes, nil, "static", srvConfig.Config.Publication.WebServer)
	return r
}

// Server defines our HTTP server
func Server() {
	Verbose = srvConfig.Config.Publication.WebServer.Verbose
	_httpReadRequest = services.NewHttpRequest("read", Verbose)
	_httpWriteRequest = services.NewHttpRequest("read", Verbose)

	// setup web router and start the service
	r := setupRouter()
	webServer := srvConfig.Config.Publication.WebServer
	server.StartServer(r, webServer)
}
