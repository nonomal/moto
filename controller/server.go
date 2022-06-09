package controller

import (
	"moto/config"
	"moto/utils"
	"net"
	"strings"
	"sync"
	"time"
)

func Listen(rule *config.Rule, wg *sync.WaitGroup) {
	defer wg.Done()
	//监听
	listener, err := net.Listen("tcp", rule.Listen)
	if err != nil {
		utils.Logger.Error(rule.Name + " failed to listen at " + rule.Listen)
		return
	}
	utils.Logger.Info(rule.Name + " listing at " + rule.Listen)
	for {
		//处理客户端连接
		conn, err := listener.Accept()
		if err != nil {
			utils.Logger.Error(rule.Name + " failed to accept at " + rule.Listen)
			time.Sleep(time.Second * 1)
			continue
		}
		//判断黑名单
		if len(rule.Blacklist) != 0 {
			clientIP := conn.RemoteAddr().String()
			clientIP = clientIP[0:strings.LastIndex(clientIP, ":")]
			if rule.Blacklist[clientIP] {
				utils.Logger.Info(rule.Name + " disconnected ip in blacklist: " + clientIP)
				conn.Close()
				continue
			}
		}
		//选择运行模式
		switch rule.Mode {
		case "normal":
			go HandleNormal(conn, rule)
		case "regex":
			go HandleRegexp(conn, rule)
		case "boost":
			go HandleBoost(conn, rule)
		case "roundrobin":
			go HandleRoundrobin(conn, rule)
		}
	}
}