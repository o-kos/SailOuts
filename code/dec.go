package main

import (
	"fmt"
	"os"
	"encoding/xml"
	"strings"
)

const (
	xmlFileName string = "SalOuts.xml"
	rootName    string = "SalOut"
	childName   string = "SalOutDetail"
	rootIndex   int = 0
)

type Node struct {
	name    string
	Count   int
	csvFile *os.File
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewNode(name string) *Node {
	n := new(Node)
	n.name = name

	f, err := os.Create(n.name + ".csv")
	check(err)
	n.csvFile = f
	n.Count = -1

	return n
}

func (n *Node) process(parent string, attr []xml.Attr) int {
	if n.Count == -1 {
		var header []string
		if parent != "" {
			header = append(header, "parent")
		}
		for _, v := range attr {
			header = append(header, v.Name.Local)
		}
		n.csvFile.WriteString(strings.Join(header, ","))
		n.csvFile.WriteString("\n")
		n.Count = 0
	}

	var values []string;
	if parent != "" {
		values = append(values, parent)
	}
	for _, v := range attr {
		values = append(values, v.Value)
	}
	n.csvFile.WriteString(strings.Join(values, ","))
	n.csvFile.WriteString("\n")
	n.Count++

	return n.Count
}

func main() {
	fmt.Printf("Start parcing %s...\n", xmlFileName)

	xmlFile, err := os.Open(xmlFileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	nodes := make(map[string]*Node)
	nodes[rootName]  = NewNode(rootName)
	nodes[childName] = NewNode(childName)
	var parent string
	var inElement string
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {		//TODO Replace for "if"
		case xml.StartElement:
			inElement = se.Name.Local
			var pv string
			if node, ok := nodes[inElement]; ok {
				if inElement == rootName {
					parent = se.Attr[rootIndex].Value
				} else {
					pv = parent
				}
				node.process(pv, se.Attr)

				var c = nodes[childName].Count
				if c < 0 {
					c = 0
				}
				fmt.Printf("Processed parents: %d, childs: %d\r", nodes[rootName].Count, c)
			}
		default:
		}

	}

	fmt.Println("\nReady!")
}
