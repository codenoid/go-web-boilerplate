package main

import (
	"github.com/go-redis/redis"
	"github.com/qiniu/qmgo"
)

var mainDB *qmgo.Database
var sessionDB *redis.Client
var appName = "App Example"
