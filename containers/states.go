package containers

import (
  "context"
  "fmt"
  "io/ioutil"
  "net/http"
  "github.com/misthios/skypod"
)


/*
Start the container
returns error
*/
func Start(ctx context.Context, name string) error{
  client, err := skypod.GetClient(ctx)
  if err !=nil {return err}

  response, err := client.Request(ctx,http.MethodPost,"/containers/" + name+ "/start",nil,nil,nil)
  defer response.Body.Close()
  if err !=nil {return err}

  switch response.StatusCode{
  case 204:
    return nil
  case 304:
    return fmt.Errorf("Container already started")
  case 404:
    return fmt.Errorf("Container not found")
  case 500:
    body,err := ioutil.ReadAll(response.Body)
    if err !=nil {return err}

    cause,err := skypod.HandleApiError(body)
    if err !=nil {return err}

		return fmt.Errorf("Failed to start container: %s",cause)
  }
  return nil 
 }

