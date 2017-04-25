package routers

import (  
	"github.com/astaxie/beego"
	"OAuthServer/controllers"
	_"net"
	_"net/http"
	"net/rpc"
)  

func regiRpcServer() error {
	server := rpc.NewServer()
    err := server.Register(new(controllers.MsgInfo))
    if err != nil {
        beego.Error("Format of service Arith isn't correct. %s", err)
		return err
    }
	
	beego.Handler("/rpc", server)

	/*
    listener, err := net.Listen("tcp", "127.0.0.1:8003")
    if err != nil {
        beego.Error("Couldn't start listening on port 1234. Error", err)
		return err
    }
    beego.Debug("Serving RPC handler",listener)
  	err = http.Serve(listener, nil)
    if err != nil {
        beego.Error("Error serving", err)
		return err
    }*/
	return nil
}




