// +build dev

package main

import "net/http"

var assets = http.Dir("assets")
