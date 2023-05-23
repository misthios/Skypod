package containers

import (
  "context"
  "fmt"
  "io/ioutil"
  "net/http"
  "github.com/misthios/skypod"
)


/*
shared code for starting/stopping containers
returns error
*/
func poweraction(ctx context.Context, name string,state string) error{
  client, err := skypod.GetClient(ctx)
  if err !=nil {return err}

  response, err := client.Request(ctx,http.MethodPost,"/containers/" + name+ "/" + state,nil,nil,nil)
  defer response.Body.Close()
  if err !=nil {return err}

  switch response.StatusCode{
  case 204:
    return nil
  case 304:
    if state == "start" {return fmt.Errorf("Container already started")}
    if state == "stop" {return fmt.Errorf("Container already stopped")}
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

 /*
 Start a container
 returns error
 */
 func Start(ctx context.Context,name string) error{
   return poweraction(ctx,name,"start")
 }

 /*
 Stop a container
 returns error
 */
 func Stop(ctx context.Context,name string) error{
   return poweraction(ctx,name,"stop")
 }




