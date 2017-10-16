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
	if "handle" == r.input[0][0] {
		fmt.Printf("func (h *Handle) %s(", r.name)
	} else {
		fmt.Printf("func (c *Client) %s(", r.name)
	}
	prevType := ""
	for _, in := range r.input {
		if in[0] == "handle" {
			continue
		}
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
	if prevType != "" {
		fmt.Printf(" %s) (", prevType)
	} else {
		fmt.Print(") (")
	}
	var toReturn string
	switch r.reply {
	case "":
		toReturn = "error"
	case "Handle":
		toReturn = "Handle, error"
	case "Data":
		toReturn = "[]byte, error"
	case "Name":
		toReturn = "[]Name, error"
	case "Attrs":
		toReturn = "Attrs, error"
	case "ExtendedReply":
		toReturn = "[]byte, error"
	case "Version":
		toReturn = "[][2]string, error"
	}
	fmt.Printf("%s) {\n", toReturn)

	// Create a packet ID and buffer
	if r.input[0][0] == "handle" {
		fmt.Println("id := NextId(h.client)")
	} else {
		fmt.Println("id := NextId(c)")
	}
	fmt.Printf("buf := id.Buffer%s(", r.name)

	for i, in := range r.input {
		if i == 0 && in[0] == "handle" {
			fmt.Print("h.h")
		} else if i == 0 {
			fmt.Print(in[0])
		} else {
			fmt.Printf(", %s", in[0])
		}
	}
	fmt.Println(")")

	if r.input[0][0] == "handle" {
		fmt.Println("reply := <-h.client.send(buf)")
	} else {
		fmt.Println("reply := <-c.send(buf)")
	}

	switch r.reply {
	case "Handle":
		fmt.Println(`return parseHandleResponse(reply,c)`)
	case "Attrs":
		fmt.Println("return parseAttrsResponse(reply)")
	case "Data":
		fmt.Println("return parseDataResponse(reply)")
	case "Name":
		fmt.Println("return parseNameResponse(reply)")
	case "ExtendedReply":
		fmt.Println("return parseExtendedReplyResponse(reply)")
	case "Version":
		fmt.Println(`return parseVersionResponse(reply)`)
	default:
		fmt.Println(`return parseStatusResponse(reply)`)
	}

	fmt.Println("}")
}

func main() {
	scan := bufio.NewScanner(os.Stdin)

	fmt.Println("package sftp\n\n// Automatically generated file. Do not touch.\n")

	var req *request
	for scan.Scan() {
		t := scan.Text()
		switch {
		case strings.HasPrefix(t, "//"):
			if req != nil {
				req.Print()
				req = nil
			}
			fmt.Println(t)
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
