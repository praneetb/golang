package config

import (
  "fmt"
  "log"
  "io/ioutil"
  "plugin"
  "strings"
  "github.com/go-ini/ini"
)

type PluginConfig struct {
    PluginDirectory   string
    WatchDirectory    string
    Functions         []string
    NumberOfJobs      int
    NumberOfGofers    int
}

type Config struct {
  FileName     string
  PluginConfig   *PluginConfig
  PluginNames  []string
}

var counter int

func (c *Config)ReadConfig() error {
  cfg := ini.Empty()
  err := cfg.Append(c.FileName)
  if err != nil {
    fmt.Print("Error !!")
    return err
  }

  c.PluginConfig = &PluginConfig {
    NumberOfJobs:   1,
    NumberOfGofers: 1,
  }

  err = cfg.Section("PLUGIN").MapTo(c.PluginConfig)
  if err != nil {
    fmt.Print("Error !!")
    return err
  }

  fmt.Println("Plugin Directory:", c.PluginConfig.PluginDirectory)
  fmt.Println("Watch Directory:", c.PluginConfig.WatchDirectory)
  fmt.Println("Functions:", c.PluginConfig.Functions)

  return nil
}

func (c *Config)ReadPlugin() {
  var fullname []string

  filename := c.PluginConfig.PluginDirectory
  files, err := ioutil.ReadDir(filename)
    if err != nil {
        log.Fatal(err)
    }

    fullname = append(fullname, c.PluginConfig.PluginDirectory)

    for _, file := range files {
      c.PluginNames = append(c.PluginNames, strings.Join( append(fullname, file.Name()), "" ) )
      fmt.Println("PLUGIN NAMES", c.PluginNames)
    }
}

func (c *Config)CheckPlugin() {
  // Open the plugin .so file to load the symbols
  plug, err := plugin.Open(c.PluginNames[0])
  if err != nil {
    panic(err)
  }

  for i, fn := range c.PluginConfig.Functions {
    fmt.Println(i, fn)
    // look up the exported function
    f, err := plug.Lookup(fn)
    if err != nil {
      fmt.Println("Function Not Found: ", fn)
    } else {
      f.(func())()
    }
  }
  return
}

func NewConfig(config_file string) (*Config) {
  conf := new(Config)
  conf.FileName = config_file

  conf.ReadConfig()
  return conf
}
