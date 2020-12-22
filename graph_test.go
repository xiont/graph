package main

import (
	"Graph/block"
	"fmt"
	"math/rand"
	"testing"
	"time"
)


var nodes []interface{}

func Benchmark(b *testing.B) {
	graph := New()
	graph.AddNode(nodes[0])
	for i:=1;i<len(nodes);i++{
		graph.AddNode(nodes[i])
		graph.AddEdge(nodes[i-1],nodes[i])
		if(i>1){
			graph.AddEdge(nodes[i-2],nodes[i])
		}
	}
	graph.SetRoot(nodes[0])
	//graph.SubGraph(nodes[0])
	L:=graph.LogicSort(func(i interface{}, i2 interface{}) bool {
		return i.(*block.Block).GetName() > i2.(*block.Block).GetName()
	})

	print("aaa",len(L))

}

func TestMain(m *testing.M)  {
	rand.Seed(time.Now().Unix())
	nodes = make([]interface{}, 1000)

	for i := 0; i < 1000; i++ {
		nodes[i] = block.New(fmt.Sprint(rand.Int()))
	}
	m.Run()
}