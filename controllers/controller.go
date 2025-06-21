package controllers

import "net/http"

// Route defines the structure for a single route
type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

// Controller defines the interface for all controllers
type Controller interface {
	Routes() []Route
}