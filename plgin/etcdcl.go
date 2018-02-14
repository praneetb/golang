/*
 * Copyright 2018 - Praneet Bachheti
 *
 * ETCD Client interface
 *
 */

package main

import (
  "fmt"
  "github.com/coreos/etcd/client"
  "golang.org/x/net/context"
)


type Callback func(key, newValue string)

type EtcdClient interface {
  // Get gets a value in Etcd
  Get(key string) (string, error)

  // Set sets a value in Etcd
  Set(key, value string) error

  // Recursively Watches a Dirctory for changes
  WatchRecursive(directory string, callback Callback) error
}

// IntentEtcdClient implements EtcdClient
type IntentEtcdClient struct {
    etcd client.Client
}


// Dial constructs a new EtcdClient
func Dial(etcdURI string) (EtcdClient, error) {
  cfg := client.Config{
    Endpoints: []string{etcdURI},
    //Transport: DefaultTransport,
  }

  etcd, err := client.New(cfg)
  if err != nil {
    fmt.Printf("Error connecting to ETCD: %s\n", etcdURI)
    return nil, err
  }
  return &IntentEtcdClient{etcd}, nil
}

// Get gets a value in Etcd
func (etcdClient *IntentEtcdClient) Get(key string) (string, error) {
  api := client.NewKeysAPI(etcdClient.etcd)
  response, err := api.Get(context.Background(), key, nil)
  if err != nil {
    if client.IsKeyNotFound(err) {
      return "", nil
    }
    return "", err
  }
  return response.Node.Value, nil
}

// Set sets a value in Etcd
func (etcdClient *IntentEtcdClient) Set(key, value string) error {
  api := client.NewKeysAPI(etcdClient.etcd)
  _, err := api.Set(context.Background(), key, value, nil)
  return err
}

func (etcdClient *IntentEtcdClient) WatchRecursive(directory string, callback Callback) error {
  api := client.NewKeysAPI(etcdClient.etcd)
  afterIndex := uint64(0)

  for {
    watcher := api.Watcher(directory, &client.WatcherOptions{Recursive: true, AfterIndex: afterIndex})
    response, err := watcher.Next(context.Background())
    if err != nil {
      if shouldIgnoreError(err) {
        continue
      }
      return err
    }

    afterIndex = response.Index
    callback(response.Node.Key, response.Node.Value)
  }
}

func shouldIgnoreError(err error) bool {
  switch err := err.(type) {
  default:
    return false
  case *client.Error:
    return err.Code == client.ErrorCodeEventIndexCleared
  }
}