package request

import "strings"

// The Input struct is defined by the "opa-istio-plugin":
// https://github.com/open-policy-agent/opa-istio-plugin#example-input

type Input struct {
	ParsedPath  []string    `json:"parsed_path"`
	ParsedQuery ParsedQuery `json:"parsed_query"`
	ParsedBody  ParsedBody  `json:"parsed_body"`
	Attributes  Attributes  `json:"attributes"`
}

func (input *Input) QueryArgs() map[string]string {
	res := make(map[string]string)

	splittedPath := strings.Split(input.Attributes.Request.HTTP.Path, "?")

	if len(splittedPath) == 1 {
		return res
	}

	for _, argString := range strings.Split(splittedPath[1], "&") {
		splitted := strings.Split(argString, "=")
		res[splitted[0]] = splitted[1]
	}

	return res
}

type Attributes struct {
	Source      Destination `json:"source"`
	Destination Destination `json:"destination"`
	Request     Request     `json:"request"`
}

type Destination struct {
	Address Address `json:"address"`
}

type Address struct {
	Address AddressClass `json:"Address"`
}

type AddressClass struct {
	SocketAddress SocketAddress `json:"SocketAddress"`
}

type SocketAddress struct {
	Address       string        `json:"address"`
	PortSpecifier PortSpecifier `json:"PortSpecifier"`
}

type PortSpecifier struct {
	PortValue int64 `json:"PortValue"`
}

type Request struct {
	HTTP HTTP `json:"http"`
}

type HTTP struct {
	ID       string            `json:"id"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Path     string            `json:"path"`
	Host     string            `json:"host"`
	Protocol string            `json:"protocol"`
	Body     string            `json:"body"`
}

type ParsedBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ParsedQuery struct {
	Lang []string `json:"lang"`
}
