package config

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Option is etcd config option.
type Option func(o *options)

type options struct {
	ctx    context.Context
	path   string
	prefix bool
}

//  Context with registry context.
func Context(ctx context.Context) Option {
	return Option(func(o *options) {
		o.ctx = ctx
	})
}

// Path is config path
func Path(p string) Option {
	return Option(func(o *options) {
		o.path = p
	})
}

// Prefix is config prefix
func Prefix(prefix bool) Option {
	return Option(func(o *options) {
		o.prefix = prefix
	})
}

type source struct {
	client  *clientv3.Client
	options *options
}

func New(client *clientv3.Client, opts ...Option) (*source, error) {
	options := &options{
		ctx:    context.Background(),
		path:   "",
		prefix: false,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.path == "" {
		return nil, errors.New("path invalid")
	}

	return &source{
		client:  client,
		options: options,
	}, nil
}

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	var opts []clientv3.OpOption
	if s.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}

	rsp, err := s.client.Get(s.options.ctx, s.options.path, opts...)
	if err != nil {
		return nil, err
	}

	var kvs []*config.KeyValue
	for _, item := range rsp.Kvs {
		kvs = append(kvs, &config.KeyValue{
			Key:   string(item.Key),
			Value: item.Value,
		})
	}
	return kvs, nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s), nil
}
