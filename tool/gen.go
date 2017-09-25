package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type request struct {
	name  string
	reply string
	input [][]string
}

func (r *request) Print() {
	// Generate function defenition
	fmt.Printf("func (c *Conn) %s(", r.name)
	prevType := ""
	for _, in := range r.input {
		switch prevType {
		case "":
			fmt.Print(in[0])
		case in[1]:
			fmt.Printf(", %s", in[0])
		default:
			fmt.Printf(" %s, %s", prevType, in[0])
		}
		prevType = in[1]
	}
	fmt.Printf(" %s) (", prevType)
	var toReturn string
	switch r.reply {
	case "":
		toReturn = "status chan error"
	case "Handle":
		toReturn = "handle chan string, status chan error"
	case "Data":
		toReturn = "response chan []byte, status chan error"
	case "Name":
		toReturn = "response chan []FxpName, status chan error"
	case "Attrs":
		toReturn = "response chan FxpAttrs, status chan error"
	case "ExtendedReply":
		toReturn = "response chan []byte, status chan error"
	}
	fmt.Printf("%s) {\n", toReturn)

	// Create a packet ID
	fmt.Println("id := c.generatePacketId()")

	// Calculate the packet length.
	fmt.Print("var pktLen uint32 = 4 + 1 + 4")
	for _, in := range r.input {
		switch in[1] {
		case "uint32":
			fmt.Print(" + 4")
		case "uint64":
			fmt.Print(" + 8")
		case "string":
			fmt.Printf(" + 4")
			fmt.Printf(" + uint32(len(%s))", in[0])
		case "[]byte":
			fmt.Printf(" + uint32(len(%s))", in[0])
		case "FxpAttrs":
			fmt.Printf(" + %s.Len()", in[0])
		}
	}
	fmt.Println()

	// Encode the packet header
	fmt.Println("buf := NewBuffer()")
	fmt.Println("buf.Grow(4 + pktLen)")
	fmt.Println("buf.WriteUint32(pktLen)")
	fmt.Printf("buf.WriteByte(FXP_%s)\n", strings.ToUpper(r.name))
	fmt.Println("buf.WriteUint32(id)")

	// Encode the request specific parameters
	for _, in := range r.input {
		switch in[1] {
		case "uint32":
			fmt.Printf("buf.WriteUint32(%s)\n", in[0])
		case "uint64":
			fmt.Printf("buf.WriteUint64(%s)\n", in[0])
		case "string":
			fmt.Printf("buf.WriteString(%s)\n", in[0])
		case "[]byte":
			fmt.Printf("buf.Write(%s)\n", in[0])
		case "FxpAttrs":
			fmt.Printf("buf.WriteAttrs(%s)\n", in[0])
		}
	}

	switch r.reply {
	case "Handle":
		fmt.Println("handle = c.handleResponse(id)")
	case "Attrs":
		fmt.Println("response = c.attrsResponse(id)")
	case "Data":
		fmt.Println("response = c.dataResponse(id)")
	case "Name":
		fmt.Println("response = c.nameResponse(id)")
	case "ExtendedReply":
		fmt.Println("response = c.extendedReplyResponse(id)")
	}
	fmt.Println("status = c.statusResponse(id)")
	fmt.Println("c.send(id, buf)")
	fmt.Println("return")
	fmt.Println("}")
}

func main() {
	scan := bufio.NewScanner(os.Stdin)

	fmt.Println("package sftp\n\n// Automatically generated file. Do not touch.\n")

	var req *request
	for scan.Scan() {
		t := scan.Text()
		switch {
		case strings.HasPrefix(t, "  <-"):
			req.reply = strings.Trim(t, " <-\t\n")
		case strings.HasPrefix(t, "\t"):
			in := strings.Fields(t)
			in[0] = strings.ToLower(in[0])
			req.input = append(req.input, in)
		default:
			if req != nil {
				req.Print()
			}
			req = &request{name: strings.TrimSpace(t)}
		}
	}
	req.Print()
}
