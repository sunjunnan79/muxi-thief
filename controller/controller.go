package controller

import "github.com/google/wire"

// 修改这里捏
const BASEURL = "http://muxithief.muxixyz.com/api/v1/"

//最后居然连业务层都没用上,只用上了路由层....

var ProviderSet = wire.NewSet(
	NewAuthController)
