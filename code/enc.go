package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"io"
)

const (
	xmlFileName string = "SalOuts.xml"
	rootName    string = "SalOut"
	childName   string = "SalOutDetail"
)

//type Node struct {
//	name    string
//	Count   int
//	csvFile *os.File
//}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//func NewNode(name string) *Node {
//	n := new(Node)
//	n.name = name
//
//	f, err := os.Create(n.name + ".csv")
//	check(err)
//	n.csvFile = f
//	n.Count = -1
//
//	return n
//}

//func (n *Node) process(parent string, attr []xml.Attr) int {
//	if n.Count == -1 {
//		var header []string
//		if parent != "" {
//			header = append(header, "parent")
//		}
//		for _, v := range attr {
//			header = append(header, v.Name.Local)
//		}
//		n.csvFile.WriteString(strings.Join(header, ","))
//		n.csvFile.WriteString("\n")
//		n.Count = 0
//	}
//
//	var values []string;
//	if parent != "" {
//		values = append(values, parent)
//	}
//	for _, v := range attr {
//		values = append(values, v.Value)
//	}
//	n.csvFile.WriteString(strings.Join(values, ","))
//	n.csvFile.WriteString("\n")
//	n.Count++
//
//	return n.Count
//}

func readLine(reader *csv.Reader) []string {
	record, err := reader.Read()
	if err != io.EOF {
		check(err)
	}
	return record
}

func record2xml(tag string, header []string, record []string) string {
	ret := fmt.Sprintf("<%s", tag)
	for i, v := range record {
		ret += fmt.Sprintf(" %s=\"%s\"", header[i], v)
	}
	return ret
}

func main() {
	fmt.Printf("Start joining %s.csv & %s.csv...\n", rootName, childName)

	csvFile, err := os.Open(rootName + ".csv")
	check(err)
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)
	header := readLine(csvReader)
	index := make(map[string]int)
	var lines []string
	for {
		record := readLine(csvReader)
		if len(record) == 0 {
			break
		}
		lines = append(lines, fmt.Sprintf("\t\t%s>\n\t\t\t<SalOutDetails>\n", record2xml(rootName, header, record)))
		index[record[0]] = len(lines) - 1
		fmt.Printf("\rProcessed parents: %d", len(lines))
	}

	csvFile, err = os.Open(childName + ".csv")
	check(err)
	csvReader = csv.NewReader(csvFile)
	header = readLine(csvReader)
	_, header = header[0], header[1:]
	childs := 0
	for {
		record := readLine(csvReader)
		if len(record) == 0 {
			break
		}

		var key string
		key, record = record[0], record[1:]
		lines[index[key]] += fmt.Sprintf("\t\t\t\t%s/>\n", record2xml(childName, header, record))
		childs++
		fmt.Printf("\rProcessed parents: %d, childs: %d", len(lines), childs)
	}

	xmlFile, err := os.Create(xmlFileName)
	check(err)
	fmt.Fprintln(xmlFile, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<ROOT>\n\t<SalOuts>")
	for _, line := range lines {
		fmt.Fprintf(xmlFile, "%s\t\t\t</SalOutDetails>\n\t\t</SalOut>\n", line)
	}
	fmt.Fprintln(xmlFile, "\t</SalOuts>\n</ROOT>")

	fmt.Printf("\nFile %s is ready!\n", xmlFileName)
}
