package main

import (
  "fmt"
  "log"
  "io/ioutil"
  "plugin"
  "runtime"
  "strings"
  "time"
  "github.com/go-ini/ini"
  elec "github.com/praneetb/golang/plgin/election"
  etcl "github.com/praneetb/golang/plgin/etcdclient"
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

func (p *Plugin)ReadConf() error {
  cfg := ini.Empty()
  err := cfg.Append("main.ini")
  if err != nil {
    fmt.Print("Error !!")
    return err
  }

  p.Conf = &PluginConf {
	  NumberOfJobs:   1,
    NumberOfGofers: 1,
  }
  err = cfg.Section("PLUGIN").MapTo(p.Conf)
  if err != nil {
    fmt.Print("Error !!")
    return err
  }

  fmt.Println("Plugin Directory:", p.Conf.PluginDirectory)
  fmt.Println("Watch Directory:", p.Conf.WatchDirectory)
  fmt.Println("Functions:", p.Conf.Functions)

  return nil
}

func (p *Plugin)ReadPlugin() {
  var fullname []string

  filename := p.Conf.PluginDirectory
  files, err := ioutil.ReadDir(filename)
    if err != nil {
        log.Fatal(err)
    }

    fullname = append(fullname, p.Conf.PluginDirectory)

    for _, file := range files {
      p.PluginNames = append(p.PluginNames, strings.Join( append(fullname, file.Name()), "" ) )
      fmt.Println("PLUGIN NAMES", p.PluginNames)
    }
}

func (p *Plugin)CheckPlugin() {
  // Open the plugin .so file to load the symbols
  plug, err := plugin.Open(p.PluginNames[0])
  if err != nil {
    panic(err)
  }

  for i, fn := range p.Conf.Functions {
    fmt.Println(i, fn)
    // look up the exported function
    f, err := plug.Lookup(fn)
    if err != nil {
      fmt.Println("Function Not Found: ", fn)
    } else {
      f.(func())()
    }
  }

}

// Callback function
func HandleMsg(key, newValue string) {
  fmt.Printf("Got Msg %s: %s\n", key, newValue)

  m := make(map[string]string)
  m[key] = newValue
  
  counter += 1
  //AddJob(counter, m)
  fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
}

func main() {
  fmt.Print("Hello\n\n")

  pl := new(Plugin)
  pl.ReadConf()

  //WatcherInit(pl.Conf.NumberOfJobs)
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
  err = etcdcl.WatchRecursive(pl.Conf.WatchDirectory, HandleMsg)
  if err != nil {
    fmt.Print("Error: ", err)
    return
  }

  //AddJob(10, nil)
  //AddJob(20, nil)
  //AddJob(30, nil)
  //AddJob(40, nil)
  //AddJob(50, nil)
  //AddJob(700, nil)
  //for ;; {
    //
  //}
  //pl.ReadPlugin()
  //pl.CheckPlugin()

  /*
  val, err1 := etcdcl.Get("message")
  if err1 == nil {
    fmt.Print("Value: ", val)
  } else {
    fmt.Print("Error: ", err1)
  }
  */

  for ; ; {
    fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
    time.Sleep(100 * time.Second)
  }
}
