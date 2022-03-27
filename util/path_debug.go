//go:build debug
// +build debug

package util

import (
	"github.com/5g-core/path_util"
)

var (
	UdrLogPath           = path_util.N5GCPath("n5gc/udrsslkey.log")
	UdrPemPath           = path_util.N5GCPath("n5gc/support/TLS/_debug.pem")
	UdrKeyPath           = path_util.N5GCPath("n5gc/support/TLS/_debug.key")
	DefaultUdrConfigPath = path_util.N5GCPath("n5gc/config/udrcfg.yaml")
)
