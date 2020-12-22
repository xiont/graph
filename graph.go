package main

import (
	"Graph/block"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/emirpasic/gods/stacks/arraystack"
)



type IGraph interface {
	GetNodeList() []interface{}
	GetRoot() interface{}
	AddNode(interface{}) error
	HaveNode(interface{}) (error,bool)
}

type Graph struct {
	nodeList   []interface{}
	lookup     map[interface{}]int
	nodeMap    map[interface{}][]interface{}
	indegreeMap map[interface{}][]interface{}
	root       interface{}
}

func New() *Graph{
	return &Graph{
		nodeList:   []interface{}{},
		lookup:     make(map[interface{}]int),
		nodeMap:    make(map[interface{}][]interface{}),
		indegreeMap: make(map[interface{}][]interface{}),
	}
}

func (G *Graph)HaveNode(node interface{}) (error,bool){
	_,ok := G.lookup[node]
	if !ok{
		return errors.New(fmt.Sprintf("lookup cannot find such %v",node) ),ok
	}

	return nil,ok


}

func (G *Graph)AddNode(node interface{}) error{
	_,ok := G.HaveNode(node)
	if ok{
		return errors.New(fmt.Sprintf("Graph has such a node %v",node))
	}

	G.nodeList = append(G.nodeList,node)
	G.lookup[node] = len(G.nodeList) -1
	return nil
}

func (G *Graph)AddEdge(node1 interface{},node2 interface{}) error{
	err1,ok1 := G.HaveNode(node1)
	err2,ok2 := G.HaveNode(node2)
	if !ok1 || !ok2{
		panic(errors.New(err1.Error()+"||"+err2.Error()))
		// return  errors.New(err1.Error()+"||"+err2.Error()), false
	}
	if node1 == node2{
		//panic(errors.New(fmt.Sprintf("node %v is same with node %v",node1,node2) ))
		return  errors.New(fmt.Sprintf("node %v is same with node %v",node1,node2) )
	}

	G.nodeMap[node1] = append(G.nodeMap[node1], node2)
	return nil
}


func (G *Graph)SetRoot(node interface{}) error{
	err,ok := G.HaveNode(node)
	if !ok{
		panic(err)
		return err
	}
	G.root = node
	return nil
}


func (G *Graph)Child(node interface{}) []interface{}{
	err,ok := G.HaveNode(node)
	if !ok{
		panic(err)
	}
	children,_ := G.nodeMap[node]

	return children
}


var usefulIdList *linkedhashset.Set
func (G *Graph)addUsefulId(children []interface{}){
	if len(children)>0 {
		for _,node := range children{
			usefulIdList.Add(G.lookup[node])
			//usefulIdList = append(usefulIdList, G.lookup[node])
			G.addUsefulId(G.Child(node))
		}
	}
}


func (G *Graph)SubGraph(node interface{}){
	// also trim useless blocks and edges
	err,ok := G.HaveNode(node)
	if !ok{
		panic(err)
	}

	L := hashset.New()

	stack := arraystack.New()
	stack.Push(node)
	var tmpNode interface{}
	for !stack.Empty(){
		tmpNode, _ = stack.Pop()
		L.Add(tmpNode)
		if len(G.nodeMap[tmpNode])>0{
			for _,node := range G.nodeMap[tmpNode]{
				stack.Push(node)
			}
		}
	}

	for _,node := range G.nodeList{
		if !L.Contains(node){
			delete(G.lookup,node)
		}
	}
	G.nodeList = L.Values()
	for key,node := range G.nodeList{
		G.lookup[node] = key
	}
	for key,nodes := range G.nodeMap{
		for _,node := range nodes {
			G.indegreeMap[node] = append(G.indegreeMap[node], key)
		}
	}
	G.indegreeMap[node] = []interface{}{}

	//fmt.Printf("%v",G.indegreeMap)
	G.root = node
}

func contain(v []interface{},i interface{}) int  {
	for key,value := range v{
		if value==i{
			return key
		}
	}
	return -1
}


func (G *Graph)LogicSort( cmp func(interface{},interface{}) int) []interface{}{
	L := linkedhashset.New() // result set
	S := linkedhashset.New() // node set that in-degree = 0
	S.Add(G.root)


	for S.Size()>0 {
		node := S.Values()[0]
		L.Add(node)
		S.Remove(node)

		// add items no in-degree to S except in L
		//遍历map中的key
		for key,value := range G.indegreeMap{
			if !L.Contains(key){
				co := contain(value,node)
				if co != -1{
					G.indegreeMap[key] = append(G.indegreeMap[key][0:co],G.indegreeMap[key][co+1:]... )
				}
				if len(G.indegreeMap[key]) == 0{
					S.Add(key)
				}
			}
		}


	}

	fmt.Printf("%v",L.Values())

	return L.Values()
}


func main() {
	block1 := block.New("a")
	block2 := block.New("b")
	block3 := block.New("c")
	block4 := block.New("d")
	block5 := block.New("e")
	graph := New()
	err := graph.AddNode(block1)
	if err != nil{
		println(err.Error())
	}
	_ = graph.AddNode(block2)
	_ = graph.AddNode(block3)
	_ = graph.AddNode(block4)
	_ = graph.AddNode(block5)


	_ = graph.AddEdge(block1, block2)
	_ = graph.AddEdge(block1, block3)
	_ = graph.AddEdge(block2, block3)
	_ = graph.AddEdge(block3, block4)
	_ = graph.AddEdge(block4, block5)

	// graph.edgeMatrix[0][0] = false

	// println(graph.edgeMatrix[0][0])
	graph.SubGraph(block1)
	graph.LogicSort(func(i interface{}, i2 interface{}) int {
		return 1
	})
	//graph.LogicSort( func(block1 interface{},block2 interface{}) int{ return 1} )

	//fmt.Printf("%v,%v,%v\n",graph.nodeList,graph.lookup[graph.nodeList[1]],graph.edgeMatrix)
	//fmt.Printf("%v\n",graph.edgeMatrix)
	//println(graph.nodeList[0].(*block.Block).GetName())
	//fmt.Printf("%v\n",graph.Child(block1)[1].(*block.Block))
}
