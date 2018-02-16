package main

import (
  "fmt"
  "runtime"
  "time"
  "github.com/praneetb/golang/plgin/config"
  elec "github.com/praneetb/golang/plgin/election"
  etcl "github.com/praneetb/golang/plgin/etcdclient"
  "github.com/praneetb/golang/plgin/watch"
)

type PluginConf struct {
    PluginDirectory   string
    WatchDirectory    string
    Functions         []string
    NumberOfJobs      int
    NumberOfGofers    int
}

type Plugin struct {
  Name        string
  Conf        *PluginConf
  PluginNames []string
}

var counter int

// Callback function
func HandleMsg(index uint64, key, newValue string) {
  fmt.Printf("Index: %d, Got Msg %s: %s\n", key, newValue, index)

  m := make(map[string]string)
  m[key] = newValue

  counter += 1
  watch.AddJob(counter, m)
  fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
}

func main() {
  fmt.Print("Hello\n\n")

  conf := new(config.Config)
  //conf.ReadPlugin()
  //conf.CheckPlugin()

  //WatcherInit(conf.Conf.NumberOfJobs)
  //InitDispatcher(pl.Conf.NumberOfGofers)
  //RunDispatcher()

  etcdcl, err := etcl.Dial("http://127.0.0.1:2379")
  if err != nil {
    fmt.Print("Error: ", err)
    return
  }

  // Master Election
  elec.NewMember(etcdcl)

  // Create Work Queues
  // etcl.NewEtcdQueues(etcdcl)

  //  Watch the Configured etcd directory for messages
  err = etcdcl.WatchRecursive(conf.PluginConfig.WatchDirectory, HandleMsg)
  if err != nil {
    fmt.Print("Error: ", err)
    return
  }

  for ; ; {
    fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
    time.Sleep(100 * time.Second)
  }
}
