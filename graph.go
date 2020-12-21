package main

import (
	"Graph/block"
	"errors"
	"fmt"
)



type IGraph interface {
	GetNodeList() []interface{}
	GetRoot() interface{}
	AddNode(interface{}) error
	HaveNode(interface{}) (error,bool)
}

type Graph struct {
	nodeList []interface{}
	nodeIndexMap map[interface{}]int
	edgeMatrix [][]int
	root interface{}
}

func New() *Graph{
	return &Graph{
		nodeList:  []interface{}{},
		nodeIndexMap: make(map[interface{}]int),
		edgeMatrix: [][]int{},
	}
}

func (G *Graph)HaveNode(node interface{}) (error,bool){
	_,ok := G.nodeIndexMap[node]
	if !ok{
		return errors.New(fmt.Sprintf("nodeIndexMap cannot find such %v",node) ),ok
	}

	return nil,ok


}

func (G *Graph)AddNode(node interface{}) error{
	_,ok := G.HaveNode(node)
	if ok{
		return errors.New(fmt.Sprintf("Graph has such a node %v",node))
	}

	G.nodeList = append(G.nodeList,node)
	G.nodeIndexMap[node] = len(G.nodeList) -1
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

	id1 := G.nodeIndexMap[node1]
	id2 := G.nodeIndexMap[node2]

	// expand for edgeMatrix
	var maxId int
	if id1 > id2{
		maxId = id1
	}else{
		maxId = id2
	}

	lenEdgeMatrix := len(G.edgeMatrix)
	expandNum := 0
	for ;expandNum + lenEdgeMatrix < maxId+1;{
		if lenEdgeMatrix == 0{
			expandNum += 1
		}
		expandNum += lenEdgeMatrix
	}
	if expandNum > 0{
		var temp1 []int
		for j:=0;j <lenEdgeMatrix;j++{
			temp1 = append(temp1,0)
		}

		for i := 0; i < expandNum;i++{
			G.edgeMatrix = append(G.edgeMatrix,temp1)
		}

		var temp2 []int
		for j:=0;j < expandNum;j++{
			temp2 = append(temp2,0)
		}
		for i := 0; i < len(G.edgeMatrix);i++{
			G.edgeMatrix[i] = append(G.edgeMatrix[i],temp2...)
		}
	}

	// add an edge
	G.edgeMatrix[id1][id2] = 1

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

	var children []interface{}
	id := G.nodeIndexMap[node]
	for i:=0;i<len(G.edgeMatrix[id]);i++{
		if G.edgeMatrix[id][i] >0{
			children = append(children,G.nodeList[i])
		}
	}

	return children
}


var usefulIdListMap map[int]interface{}
func (G *Graph)addUsefulId(children []interface{}){
	if len(children)>0 {
		for _,node := range children{
			usefulIdListMap[G.nodeIndexMap[node]] = &struct{}{}
			//usefulIdList = append(usefulIdList, G.nodeIndexMap[node])
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
	usefulIdListMap = make(map[int]interface{})
	usefulIdListMap[G.nodeIndexMap[node]] = &struct{}{}
	//usefulIdList = []int{G.nodeIndexMap[node] }

	G.addUsefulId(G.Child(node))
	fmt.Printf("%v\n",usefulIdListMap)
}

func (G *Graph)LogicSort(){

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
	graph.SubGraph(block5)
	//fmt.Printf("%v\n",graph.edgeMatrix)
	//println(graph.nodeList[0].(*block.Block).GetName())
	//fmt.Printf("%v\n",graph.Child(block1)[1].(*block.Block))
}
