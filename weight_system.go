package main

import (
	"math"
	"strconv"
)

type behavior string

var alp float64

const bet float64 = 0.75
const (
	nreject     behavior = "reject"
	nsuccess    behavior = "success"
	ndisconnect behavior = "disconnect"
)

func (p *pbft) updateWeight(behavior behavior, nodeid string) {
	num, _ := strconv.ParseInt(nodeid[1:], 10, 0)
	if behavior == "success" {
		alp = math.Pow(1-p.nodeWeight[num], 1)
		p.nodeWeight[num] = p.nodeWeight[num] + alp - alp*p.nodeWeight[num] //指数平滑法
		//println("共识成功")
		return
	}
	if behavior == "reject" {
		p.nodeWeight[num] = p.nodeWeight[num] * bet
		return
	}
	if behavior == "disconnect" {
		p.nodeWeight[num] *= math.Pow(math.E, -alp)
		return
	}
	println("行为未知更新失败")
}
