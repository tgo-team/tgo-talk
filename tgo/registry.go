package tgo

import (
	"fmt"
	"sync"
)



const (
	newServerPrefix = "newServer:"
	newRoutePrefix = "newRoute:"
	newProtocolPrefix = "newProtocol:"
	newLogPrefix = "newLog:"
)

var registryMap map[string]interface{}


type newServerFunc func(*Context) Server
type newRouteFunc func(*Context) Route
type newProtocol func() Protocol
type newLog func(logLevel LogLevel) Log

var clientLock sync.RWMutex
var tContextLock sync.RWMutex
func init()  {
	registryMap = map[string]interface{}{}
}

// 登记server
func RegistryServer(newFunc newServerFunc)  {
	registryMap[fmt.Sprintf("%s",newServerPrefix)] = newFunc
}


// 登记协议
func RegistryProtocol(name string,newFunc newProtocol)  {
	registryMap[fmt.Sprintf("%s-%s",newProtocolPrefix,name)] = newFunc
}


//// 登记路由
//func RegistryRoute(newFunc newRouteFunc)  {
//	registryMap[fmt.Sprintf("%s",newRoutePrefix)] = newFunc
//}

func RegistryLog(newFunc newLog)  {
	registryMap[fmt.Sprintf("%s",newLogPrefix)] = newFunc
}


func NewServer(context *Context) Server  {
	key := fmt.Sprintf("%s",newServerPrefix)
	serverFuncObj := registryMap[key]
	if serverFuncObj!=nil {
		return  serverFuncObj.(newServerFunc)(context)
	}
	return nil
}

//func NewRoute(ctx *Context) Route  {
//	key := fmt.Sprintf("%s",newRoutePrefix)
//	funcObj := registryMap[key]
//	if funcObj!=nil {
//		return  funcObj.(newRouteFunc)(ctx)
//	}
//	return nil
//}

func NewProtocol(name string) Protocol  {
	key := fmt.Sprintf("%s-%s",newProtocolPrefix,name)
	funcObj := registryMap[key]
	if funcObj!=nil {
		return  funcObj.(newProtocol)()
	}
	return nil
}

func NewLog(logLevel LogLevel) Log  {
	key := fmt.Sprintf("%s",newLogPrefix)
	funcObj := registryMap[key]
	if funcObj!=nil {
		return  funcObj.(newLog)(logLevel)
	}
	return nil
}