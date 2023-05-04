package skypod

import (
        "context"
        "net/url"
	"io"
        "net"
        "net/http"
        "fmt"
)

type Client struct {
        URI *url.URL
        Client *http.Client
}

/*
Create a new api client for the unix socket
Creates a new context.Background with a key "Client" containing the client
returns the client / error
*/
func NewClient(uri string) (context.Context, error){

        parurl,err := url.Parse(uri)
        if err != nil {return nil,fmt.Errorf("Url %s is not a valid url")}

        if parurl.Scheme != "unix" {return nil,fmt.Errorf("Scheme %s is not supported")}

        hclient := &http.Client{
                Transport: &http.Transport{
                        DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
                                return (&net.Dialer{}).DialContext(ctx, "unix", parurl.Path)
                        },
                        DisableCompression: true,
                },
        }
        c := &Client{URI: parurl, Client: hclient}
        ctx := context.WithValue(context.Background(),"Client",c)

	ctx, err = CheckConnection(ctx)
	if err != nil {return nil,err}

        return ctx,nil
}

/*
Check if the context contains a client
return the client / error
*/
func GetClient(ctx context.Context) (*Client,error){
        if c, ok := ctx.Value("Client").(*Client); ok {
                return c, nil
        }
        return nil,fmt.Errorf("Client is not present in context")
}

/*
Make a request to the api
Caller is required to close the request body
returns the response / error
*/
func (c *Client) Request(ctx context.Context, method string,endpoint string,body io.Reader, headers http.Header,params url.Values) (*http.Response,error) {

	var uri string

        if v, ok := ctx.Value("Version").(string); ok {
        	uri = "http://d/v" + v + "/libpod"+endpoint
	} else {
        	uri = "http://d/libpod"+endpoint
	}

        req, err := http.NewRequestWithContext(ctx, method, uri, body)
        if err != nil {return nil, err}

	if len(params) > 0 {req.URL.RawQuery = params.Encode() }

	for key, val := range headers {
		for _, v := range val {req.Header.Add(key, v)}
	}
        req.Header.Set("User-Agent", "Skypod-agent/1.0")

        response, err := c.Client.Do(req)
        if err != nil {return nil, err}
        return response,nil
}

/*
Checks if the client is able to reach the api
If successfull it will set the api version in the current context
return error if not
*/
func CheckConnection(ctx context.Context) (context.Context,error) {
        client, err := GetClient(ctx)
        if err != nil {return nil,err}

        response,err := client.Request(ctx,http.MethodGet,"/_ping",nil,nil,nil)
        if err != nil {return nil,err}

        defer response.Body.Close()
	
	if response.StatusCode != http.StatusOK {
		return nil,fmt.Errorf("failed to check if the api is ready: received wrong status code")
	}

	//set the api-version 
	v := response.Header.Get("Libpod-API-Version")
	if v == "" { return nil,fmt.Errorf("API did not return the version")}

	ctx = context.WithValue(ctx,"Version",v)
        return ctx,nil
}



