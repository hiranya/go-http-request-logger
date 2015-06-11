package main

import (
	"log"
	"fmt"
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"bytes"
	"bufio"
	"time"
)

var flag_bindport = flag.String("bindport", "9000", "The port number this reverse proxy based logger should bind to")
var flag_serverhost = flag.String("serverhost", "localhost", "The hostname of the server all requests should be routed to")
var flag_serverport = flag.String("serverport", "9191", "The port number of the server all requests should be routed to")


func main() {
	flag.Parse()
	
	u, err := url.Parse("http://"+ *flag_serverhost +":"+ *flag_serverport +"/")
	if err != nil {
		log.Fatal(err)
	}

	reverse_proxy := httputil.NewSingleHostReverseProxy(u)
	http.HandleFunc("/", handler(reverse_proxy))

	fmt.Println("HTTP interceptor started on localhost:" + *flag_bindport + " routing traffic to http://"+*flag_serverhost + ":" + *flag_serverport) // well eventually :)
	if err = http.ListenAndServe(":" + *flag_bindport, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now())
		
		dump, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(dump))

		// You can optionally capture/wrap the transport if that's necessary (for
			// instance, if the transport has been replaced by middleware). Example:
			// proxy.Transport = &myTransport{proxy.Transport}
		p.Transport = &myTransport{}

		p.ServeHTTP(w, r)
	}
}

func extractBodyFromRequestDump(dump []byte) string {
	byte_reader := bytes.NewReader(dump)
	buff_reader := bufio.NewReader(byte_reader)

	var body bytes.Buffer

	isBody := false
	for {
		v, _, err := buff_reader.ReadLine()
		if err != nil {
			break
		}

		if !isBody && len(v) == 0 { // the first empty line is the delimieter between HEADER and BODY
			isBody = true
			continue
		} else if isBody {
			body.WriteString(string(v))
		}
	}

	return body.String()
}

type myTransport struct {
	// Uncomment this if you want to capture the transport
	// CapturedTransport http.RoundTripper
}


func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)
	// or, if you captured the transport
	// response, err := t.CapturedTransport.RoundTrip(request)

	// The httputil package provides a DumpResponse() func that will copy the
	// contents of the body into a []byte and return it. It also wraps it in an
	// ioutil.NopCloser and sets up the response to be passed on to the client.
	body, err := httputil.DumpResponse(response, true)
	if err != nil {
		// copying the response body did not work
		return nil, err
	}

	// You may want to check the Content-Type header to decide how to deal with
	// the body. In this case, we're assuming it's text.
	fmt.Println("HTTP response from proxyied server : "+string(body))

	return response, err
}
