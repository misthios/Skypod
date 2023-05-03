package skypod

import (
        "context"
        "net/url"
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
