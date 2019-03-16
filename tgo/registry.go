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
	newStoragePrefix = "newStorage:"
	newAuthPrefix = "newAuth:"
)

var registryMap map[string]interface{}


type newServerFunc func(*Context) Server
type newRouteFunc func(*Context) Route
type newProtocol func() Protocol
type newLog func(logLevel LogLevel) Log
type newStorageFunc func(*Context) Storage

var clientLock sync.RWMutex
var tContextLock sync.RWMutex
type authFunc func(ctx *Context)
func init()  {
	registryMap = map[string]interface{}{}
}

// 登记server
func RegistryServer(newFunc newServerFunc)  {
	serverFuncObj := registryMap[fmt.Sprintf("%s",newServerPrefix)]
	var serverFuncs []newServerFunc
	if serverFuncObj==nil {
		serverFuncs = []newServerFunc{}
	}else{
		serverFuncs = serverFuncObj.([]newServerFunc)
	}
	serverFuncs = append(serverFuncs,newFunc)
	registryMap[fmt.Sprintf("%s",newServerPrefix)] = serverFuncs
}


// 登记协议
func RegistryProtocol(name string,newFunc newProtocol)  {
	registryMap[fmt.Sprintf("%s-%s",newProtocolPrefix,name)] = newFunc
}

func RegistryLog(newFunc newLog)  {
	registryMap[fmt.Sprintf("%s",newLogPrefix)] = newFunc
}

func RegistryStorage(newFunc newStorageFunc)  {
	registryMap[fmt.Sprintf("%s",newStoragePrefix)] = newFunc
}

func NewStorage(context *Context) Storage {
	key := fmt.Sprintf("%s",newStoragePrefix)
	serverFuncObj := registryMap[key]
	if serverFuncObj!=nil {
		return  serverFuncObj.(newStorageFunc)(context)
	}
	return nil
}

func GetServers(context *Context) []Server  {
	key := fmt.Sprintf("%s",newServerPrefix)
	serverFuncObj := registryMap[key]
	servers := make([]Server,0)
	if serverFuncObj!=nil {
		serverFuncs := serverFuncObj.([]newServerFunc)
		for _,serverFunc :=range serverFuncs {
			servers = append(servers,serverFunc(context))
		}
		return  servers
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
