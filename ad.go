package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Server struct {
	host        string
	name        string
	server_type string
	description string
	application string
	environment string
}

func errorNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func dnAdd(words []string, server *Server) {
	fmt.Println(words[1])
	// Ruby named groups don't really translate here.  Using the names here is really dumb
	r := regexp.MustCompile(`^CN=(?P<server>[\w\-.]+),OU=(?P<environment>\w+),OU=(?P<application>\w+)`)
	match := r.FindStringSubmatch(words[1])
	if len(match) == 4 {
		server.host = match[1]
		server.name = match[1]
		server.server_type = "server"
		server.application = match[3]
		server.environment = match[2]
	}
	// TODO: Add routines to clean up the server name
	// TODO: Add routines to classify the server type
	return
}

func descriptionAdd(words []string, server *Server) {
	server.description = strings.Join(words[:], " ")
}

func initServer(active *bool, server *Server) {
	*active = false
	server.host = ""
	server.name = ""
	server.server_type = ""
	server.description = ""
	server.application = ""
	server.environment = ""
}

func scanLine(line string, active *bool, server *Server) {
	// Skip comment lines
	// Extract
	words := strings.Fields(line)
	if len(words) == 0 {
		return
	}
	if words[0][0:1] == "#" {
		return
	}
	switch words[0] {
	case "dn:":
		dnAdd(words, server)
	case "description:":
		descriptionAdd(words, server)
	}
	*active = true
}

func lineTypeDN(line string) bool {
	words := strings.Fields(line)
	if len(words) == 0 {
		return false
	}
	return words[0] == "dn:"
}

func outputServer(server *Server) {
	fmt.Println(server)
}

// Extract server attributes from their computer objects
func main() {

	// Source AD_SRC
	// Input  AD listing
/*
# extended LDIF
#
# LDAPv3
# base <ou=computers,ou=unix,dc=nordstrom,dc=net> with scope subtree
# filter: samaccountname=*
# requesting: samaccountname description
# with pagedResults control: size=2000000
#

# y0319t729, Test, POS, Computers, UNIX, nordstrom.net
dn: CN=y0319t729,OU=Test,OU=POS,OU=Computers,OU=UNIX,DC=nordstrom,DC=net
sAMAccountName: Y0319T729$

# y0319t349, Test, POS, Computers, UNIX, nordstrom.net
dn: CN=y0319t349,OU=Test,OU=POS,OU=Computers,OU=UNIX,DC=nordstrom,DC=net
description: MPOS hoi dev apache server
sAMAccountName: Y0319T349$

# y0319t79, Test, POS, Computers, UNIX, nordstrom.net
dn: CN=y0319t79,OU=Test,OU=POS,OU=Computers,OU=UNIX,DC=nordstrom,DC=net
description: M Merch Search /Price Check
sAMAccountName: Y0319T79$

# y0319t1238, Test, POS, Computers, UNIX, nordstrom.net
dn: CN=y0319t1238,OU=Test,OU=POS,OU=Computers,OU=UNIX,DC=nordstrom,DC=net
description: POS Test DB Server - WO#45968
sAMAccountName: Y0319T1238$
*/

	// Parse multiple lines to extract the server info

	//Output - call DB update routine to add the entries
	//   SERVER Entries - Source, name_type, host, name, description, application, environment
	//   Console Entries - Source, name_type, host, name, description, application, environment
	// name is the reported name, blank and trailing periods trimmed
	// host removes trailing domain information and trailing type information
	file, err := os.Open("../sample/active_directory.in")
	errorNil(err)
	defer file.Close()
	active := false
	server := Server{}
	initServer(&active, &server)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if active && lineTypeDN(line) {
			outputServer(&server)
			initServer(&active, &server)
		}
		scanLine(line, &active, &server)
	}
	if active {
		outputServer(&server)
	}
}
