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
	fmt.Printf("func (id PacketId) %s(", r.name)
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
	if prevType != "" {
		fmt.Printf(" %s) (", prevType)
	} else {
		fmt.Print(") (")
	}
	var toReturn string = "*Buffer"
	fmt.Printf("%s) {\n", toReturn)

	// Calculate the packet length.
	fmt.Print("var pktLen uint32 = 1 + 4")
	var pktlenappend string = ""
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
		case "Attrs":
			fmt.Printf(" + %s.Len()", in[0])
		case "[][2]string":
			pktlenappend += fmt.Sprintf("\nfor _, ext := range %s {\n", in[0])
			pktlenappend += "pktLen += 4 + uint32(len(ext[0]))\n"
			pktlenappend += "pktLen += 4 + uint32(len(ext[1]))\n"
			pktlenappend += "}\n"
		}
	}
	fmt.Println(pktlenappend)

	// Encode the packet header
	fmt.Println("buf := NewBuffer()")
	fmt.Println("buf.Grow(4 + pktLen)")
	fmt.Println("buf.WriteUint32(pktLen)")
	fmt.Printf("buf.WriteByte(FXP_%s)\n", strings.ToUpper(r.name))
	fmt.Println("buf.WriteUint32(uint32(id))")

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
		case "Attrs":
			fmt.Printf("buf.WriteAttrs(%s)\n", in[0])
		case "[][2]string":
			fmt.Printf("for _, ext := range %s {\n", in[0])
			fmt.Println("buf.WriteString(ext[0])")
			fmt.Println("buf.WriteString(ext[1])")
			fmt.Println("}")
		}
	}

	fmt.Println(`return buf`)

	fmt.Println("}")
}

func main() {
	scan := bufio.NewScanner(os.Stdin)

	fmt.Println("package sftp\n\n// Automatically generated file. Do not touch.\n")
	fmt.Println("// Used to create SFTP version 3 packets.")
	fmt.Println("// PacketId is its own type to help unclutter documentation.")
	fmt.Println("type PacketId uint32\n")

	var req *request
	for scan.Scan() {
		t := scan.Text()
		switch {
		case strings.HasPrefix(t, "//"):
			if req != nil {
				req.Print()
				req = nil
			}
			//fmt.Println(t)
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
