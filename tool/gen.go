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

	// Create a packet ID
	if r.input[0][0] == "handle" {
		fmt.Println("id := h.client.nextPacketId()")

	} else {
		fmt.Println("id := c.nextPacketId()")

	}

	// Calculate the packet length.
	fmt.Print("var pktLen uint32 = 4 + 1 + 4")
	for i, in := range r.input {
		if i == 0 && in[0] == "handle" {
			fmt.Print(" + 4 + uint32(len(h.h))")
			continue
		}
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
		case "Attrs":
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
	for i, in := range r.input {
		if i == 0 && in[0] == "handle" {
			fmt.Println("buf.WriteString(h.h)")
			continue
		}
		switch in[1] {
		case "uint32":
			fmt.Printf("buf.WriteUint32(%s)\n", in[0])
		case "uint64":
			fmt.Printf("buf.WriteUint64(%s)\n", in[0])
		case "string":
			fmt.Printf("buf.WriteString(%s)\n", in[0])
		case "[]byte":
			fmt.Printf("buf.Write(%s)\n", in[0])
		case "Attrs":
			fmt.Printf("buf.WriteAttrs(%s)\n", in[0])
		}
	}

	if r.input[0][0] == "handle" {
		fmt.Println("reply := h.client.send(buf)")
	} else {
		fmt.Println("reply := c.send(buf)")
	}
	fmt.Println("replyisnil := nil == reply")
	//fmt.Println("defer bufPool.Put(reply)")

	fmt.Println("// TODO Temporary")
	switch r.reply {
	case "Handle":
		fmt.Println(`if replyisnil {`)
		fmt.Println(`return Handle{},errors.New("Internal: Nil response channel.")`)
		fmt.Println(`}`)
		fmt.Println(`return parseHandleResponse(<-reply,c)`)
	case "Attrs":
		fmt.Println(`if replyisnil {`)
		fmt.Println(`return Attrs{},errors.New("Internal: Nil response channel.")`)
		fmt.Println(`}`)
		fmt.Println("return parseAttrsResponse(<-reply)")
	case "Data":
		fmt.Println(`if replyisnil {`)
		fmt.Println(`return nil,errors.New("Internal: Nil response channel.")`)
		fmt.Println(`}`)
		fmt.Println("return parseDataResponse(<-reply)")
	case "Name":
		fmt.Println(`if replyisnil {`)
		fmt.Println(`return nil,errors.New("Internal: Nil response channel.")`)
		fmt.Println(`}`)
		fmt.Println("return parseNameResponse(<-reply)")
	case "ExtendedReply":
		fmt.Println(`if replyisnil {`)
		fmt.Println(`return nil,errors.New("Internal: Nil response channel.")`)
		fmt.Println(`}`)
		fmt.Println("return parseExtendedReplyResponse(<-reply)")
	case "Version":
		fmt.Println(`if replyisnil {`)
		fmt.Println(`return nil,errors.New("Internal: Nil response channel.")`)
		fmt.Println(`}`)
		fmt.Println(`return parseVersionResponse(<-reply)`)
	default:
		fmt.Println(`if replyisnil {`)
		fmt.Println(`return errors.New("Internal: Nil response channel.")`)
		fmt.Println(`}`)
		fmt.Println(`return parseStatusResponse(<-reply)`)
	}

	fmt.Println("}")
}

func main() {
	scan := bufio.NewScanner(os.Stdin)

	fmt.Println("package sftp\n\n// Automatically generated file. Do not touch.\n")
	fmt.Println(`import "errors"`)

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
