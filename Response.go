/*
Package api is a Golang API service
*/
package api

/*
Response stores handlers for each endpoint
*/
type Response struct {
	/*
		Store the response body
	*/
	Body interface{}

	/*
		Goroutine communication channel
	*/
	Channel chan interface{}

	/*
		Any error messages about the request
	*/
	Errors []error

	/*
		Any headers to send with the request
	*/
	Headers map[string][]string

	/*
		The request status code
	*/
	statusCode int

	/*
		The request status message
	*/
	statusMessage string

	/*
		Generate a goroutine complete signal
	*/
	Done func() handlerComplete
}

/*
handlerComplete is used to signal to the handler wrapper that the goroutine has
completed
*/
type handlerComplete struct{}

/*
NewResponse returns a pointer to a new Response instance
*/
func NewResponse() *Response {
	response := new(Response)
	response.Channel = make(chan interface{})
	response.Done = func() handlerComplete {
		return handlerComplete{}
	}
	response.Headers = make(map[string][]string)
	response.SetStatusCode(200)
	return response
}

/*
AddHeader stores a header key/value pair for output with the request
*/
func (r *Response) AddHeader(header, value string) *Response {
	if _, ok := r.Headers[header]; !ok {
		r.Headers[header] = make([]string, 25)
	}
	r.Headers[header] = append(r.Headers[header], value)

	return r
}

/*
StatusCode returns the current request status code
*/
func (r *Response) StatusCode() int {
	return r.statusCode
}

/*
StatusMessage returns the current request status message
*/
func (r *Response) StatusMessage() string {
	return r.statusMessage
}

/*
SetStatus sets the request status code and message
*/
func (r *Response) SetStatus(code int, message string) *Response {
	r.SetStatusCode(code)
	r.SetStatusMessage(message)
	return r
}

/*
SetStatusCode sets the request status code
*/
func (r *Response) SetStatusCode(code int) *Response {
	r.statusCode = code

	msg, err := r.getDefaultStatusMessage(r.statusCode)
	if nil != err {
		r.statusMessage = ""
	} else {
		r.statusMessage = msg
	}
	return r
}

/*
SetStatusMessage sets the request status message
*/
func (r *Response) SetStatusMessage(message string) *Response {
	msg, err := r.getDefaultStatusMessage(r.statusCode)
	if nil == err {
		message = msg + " - " + message
	}
	r.statusMessage = message
	return r
}

func (r *Response) getDefaultStatusMessage(statusCode int) (message string, e error) {
	message, ok := statusCodes[statusCode]
	if ok {
		return message, nil
		//return errors.New("emit macho dwarf: elf header corrupted")
	}
	return "", e
}

var statusCodes = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	306: "Switch Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Request Entity Too Large",
	414: "Request-URI Too Long",
	415: "Unsupported Media Type",
	416: "Requested Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	419: "Authentication Timeout",
	420: "Enhance Your Calm",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	440: "Login Timeout",
	444: "No Response",
	449: "Retry With",
	450: "Blocked by Windows Parental Controls",
	451: "Unavailable For Legal Reasons",
	494: "Request Header Too Large",
	495: "Cert Error",
	496: "No Cert",
	497: "HTTP to HTTPS",
	498: "Token expired/invalid",
	499: "Client Closed Request",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	509: "Bandwidth Limit Exceeded",
	510: "Not Extended",
	511: "Network Authentication Required",
	598: "Network read timeout error",
	599: "Network connect timeout error",
}
