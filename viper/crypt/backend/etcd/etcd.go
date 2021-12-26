package etcd

import (
	"context"
	"errors"
	"fmt"
	"github.com/brickman-source/golang-utilities/viper/crypt/backend"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

type Client struct {
	client    *clientv3.Client
	keysAPI   clientv3.KV
	waitIndex uint64
}

func New(machines []string) (*Client, error) {
	newClient, err := clientv3.New(clientv3.Config{
		Endpoints: machines,
	})
	if err != nil {
		return nil, fmt.Errorf("creating new etcd client for crypt.backend.Client: %v", err)
	}
	keysAPI := clientv3.NewKV(newClient)
	return &Client{client: newClient, keysAPI: keysAPI, waitIndex: 0}, nil
}

func (c *Client) Get(key string) ([]byte, error) {
	return c.GetWithContext(context.TODO(), key)
}

func (c *Client) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	resp, err := c.keysAPI.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if resp.Count < 1 {
		return nil, errors.New("no data")
	}
	return resp.Kvs[0].Value, nil
}

func addKVPairs(kvs []*mvccpb.KeyValue, list backend.KVPairs) backend.KVPairs {
	for _, value := range kvs {
		list = append(list, &backend.KVPair{Key: string(value.Key), Value: value.Value})
	}
	return list
}

func (c *Client) List(key string) (backend.KVPairs, error) {
	return c.ListWithContext(context.TODO(), key)
}

func (c *Client) ListWithContext(ctx context.Context, key string) (backend.KVPairs, error) {
	resp, err := c.keysAPI.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if resp.Count < 1 {
		return nil, errors.New("key is not a directory")
	}
	list := addKVPairs(resp.Kvs, nil)
	return list, nil
}

func (c *Client) Set(key string, value []byte) error {
	return c.SetWithContext(context.TODO(), key, value)
}

func (c *Client) SetWithContext(ctx context.Context, key string, value []byte) error {
	_, err := c.keysAPI.Put(ctx, key, string(value))
	return err
}

func (c *Client) Watch(key string, stop chan bool) <-chan *backend.Response {
	return c.WatchWithContext(context.Background(), key, stop)
}

func (c *Client) WatchWithContext(ctx context.Context, key string, stop chan bool) <-chan *backend.Response {
	respChan := make(chan *backend.Response, 0)
	go func() {
		watcher := clientv3.NewWatcher(c.client)
		ctx, cancel := context.WithCancel(ctx)
		go func() {
			<-stop
			cancel()
		}()
		watchRespChan := watcher.Watch(ctx, key)
		for watchResp := range watchRespChan {
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					respChan <- &backend.Response{Value: event.Kv.Value, Error: nil}
				case mvccpb.DELETE:
					respChan <- &backend.Response{Value: event.Kv.Value, Error: nil}
				}
			}
		}
		//
		//for {
		//	var resp *goetcd.WatchResponse
		//	var err error
		//	// if c.waitIndex == 0 {
		//	// 	resp, err = c.client.Get(key, false, false)
		//	// 	if err != nil {
		//	// 		respChan <- &backend.Response{nil, err}
		//	// 		time.Sleep(time.Second * 5)
		//	// 		continue
		//	// 	}
		//	// 	c.waitIndex = resp.EtcdIndex
		//	// 	respChan <- &backend.Response{[]byte(resp.Node.Value), nil}
		//	// }
		//	// resp, err = c.client.Watch(key, c.waitIndex+1, false, nil, stop)
		//	resp, err = watcher.Next(ctx)
		//	if err != nil {
		//		respChan <- &backend.Response{nil, err}
		//		time.Sleep(time.Second * 5)
		//		continue
		//	}
		//	c.waitIndex = resp.Node.ModifiedIndex
		//	respChan <- &backend.Response{[]byte(resp.Node.Value), nil}
		//}
	}()
	return respChan
}
