// Code generated by go-bindata.
// sources:
// templates/client_rpc_go.tpl
// templates/data_go.tpl
// templates/handler_go.tpl
// templates/handler_rpc_go.tpl
// templates/service_go.tpl
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesClient_rpc_goTpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xac\x56\x5f\x6f\xdb\x36\x10\x7f\x16\x3f\xc5\x55\x40\x02\xa9\x10\xa4\x3d\x1b\xf0\x80\x2e\x49\xb1\x62\x58\x12\xa4\xde\xd3\x30\x14\xb4\x74\x92\xb8\x50\xa4\x4c\x52\x4e\x0d\x41\xdf\x7d\x38\x52\xfe\xd3\xcd\x6d\x5c\xac\x2f\x89\x78\xc7\xfb\xf3\xbb\xfb\xdd\xd1\x45\x01\x37\xba\x42\x68\x50\xa1\xe1\x0e\x2b\x58\xef\xa0\xe1\x92\x7f\xde\xe5\x70\xfb\x00\xf7\x0f\x2b\xb8\xbb\xfd\xb0\xca\x59\x51\xc0\x3b\x29\xa1\x6c\xb9\x6a\xd0\x42\x37\x58\x07\x6b\x84\x4a\x2b\x04\xa1\xa0\x1c\xac\xd3\x1d\x94\x52\xa0\x72\xe0\x5a\xee\xc0\xb6\x7a\x90\x15\xa0\x70\x2d\x1a\xc0\x6e\x8d\x15\x68\x03\x2f\x86\xf7\xe0\x5a\x61\x73\xc6\x7a\x5e\x3e\xf3\x06\x61\x1c\xf3\xc7\xf0\x79\xcf\x3b\x9c\x26\xc6\x44\xd7\x6b\xe3\x20\x61\x51\xbc\xde\x39\xb4\x31\x8b\x62\x54\xa5\xae\x84\x6a\x8a\xbf\xad\x56\x24\xa8\x3b\x47\xff\x84\x2e\x84\x1e\x9c\x90\x74\x50\xe8\x8a\xd6\xb9\x3e\x66\x91\xd4\x0d\xc4\x8d\x70\xed\xb0\xce\x4b\xdd\x15\x56\x98\xa1\xb7\xa8\x0a\xa9\x1b\x33\xd8\x98\xb1\xe8\x54\xdd\x3f\x37\x05\x1a\xa3\x8d\x0f\x76\xa2\xe0\xd6\x19\x4e\x7e\x43\x61\x8a\x0e\x9d\x11\xa5\x8d\x59\xca\x98\xdb\xf5\x3e\xfd\xf7\xc8\xdd\x60\xf0\xd1\x60\x2d\x3e\x4f\xd3\x47\x34\x5b\x51\xe2\x4d\x28\x87\x50\x0e\x4d\xcd\x4b\x84\x91\x45\x5f\xbd\xec\x55\xc1\xe2\xc3\xde\xe0\x17\x5d\xed\xa6\x89\x4d\x97\x05\x7a\xe8\x9d\xd0\xca\x82\x75\x66\x28\x1d\x8c\x64\x57\x0f\xaa\x84\xb2\xc5\xf2\xf9\x32\xe3\x44\xf7\x0e\xde\x5e\x76\x37\xbd\xf4\x22\xe1\x16\x35\x90\xeb\xe5\x12\x94\x90\x24\x88\xfc\x11\xae\x2f\x73\x31\x4e\x2c\x9a\x58\x64\xd0\x0d\x46\x91\xa7\x6f\x14\xe5\xd7\xd5\xea\xf1\x5c\xe9\x6f\x75\x62\x70\x03\x6f\x89\x1f\xf9\x13\x6e\x06\xb4\x2e\x85\x64\x7f\xb6\xbd\x56\x16\x33\xf0\x24\x48\x0f\xc5\xbb\xc7\x97\x57\x52\x4c\x58\x54\x4a\xf1\xcd\x44\x32\x16\x7d\x47\x61\x33\x96\xbe\x4a\xaa\xf1\x50\x8c\xeb\x03\x6f\x1e\x8d\xd8\x72\x37\x0f\xd1\x5c\xe2\xc5\xf7\x76\x3f\xcd\x58\x14\x39\xde\xd8\x05\xcc\x4c\xcf\x57\xbc\xb1\xe4\x2e\x8a\x25\xdf\xa1\x89\x17\x00\x10\xdb\x60\xff\x29\xcc\x7c\x9c\x79\xfd\x2c\x8c\x17\x10\x8f\x63\x3e\x87\x08\xf9\xf8\x1b\x13\xfd\xa9\x05\xca\xca\x2e\x40\xea\x26\x7f\xef\xbf\x7f\x98\x73\x16\x51\x2b\x16\xb4\x88\x32\x22\xcc\x09\x49\xce\x94\xe8\x38\x29\x1e\xf0\x17\x78\xd9\x9c\xe6\x49\x96\xa1\x87\x97\x36\x91\xbd\xce\x0a\x4a\xef\x38\xf4\x5d\x2f\xb1\x43\xe5\x38\x99\xcf\x93\x1f\x28\x98\x58\x1f\xf4\x0c\x82\x14\xaa\xf3\xa4\xfe\xf3\x2f\xda\x9b\x7b\x32\x07\xb6\xd8\xde\x9f\x61\xb1\x04\x9b\x97\x52\xe4\x61\x20\x52\x3f\x9c\xa4\x78\x73\x1c\x4e\x3a\x2e\xa1\xee\x5c\x7e\x47\x1e\xea\x24\x3e\x9f\xc0\x02\xae\xb6\xb1\x77\x9b\xb2\x68\xcf\x48\x25\xa4\x17\xcd\x23\x6b\x7b\x42\x93\xc1\x27\x8a\x1c\x56\x75\xfe\x84\xbc\x7a\x27\x65\x42\xda\x9c\xd4\x69\xb8\xe9\xbf\xf3\x1b\xa9\x2d\x26\x21\x31\x2f\xfd\xe8\xb8\x1b\xac\x7f\xab\xde\x2c\xc1\x63\x0d\xa2\x87\xdf\x7c\xba\xa2\x06\x89\x2a\xd9\xc7\x4a\xe1\x67\xf8\xc9\x2b\xa2\x19\xf0\x29\x92\x99\x47\xa1\x36\x70\x55\x2d\xe0\xca\xc6\xd9\xbf\x03\x65\xc4\x0e\xa1\x9a\xa3\x53\x42\xf8\x5f\x88\x84\xf1\xff\x44\x39\x11\x7c\xbd\x84\x5e\x76\xac\xa4\x12\x92\x5d\x42\x0e\x85\x2f\x4f\xb8\x49\x3a\x74\xad\xae\x66\x3c\x19\xd4\x8a\xf4\x87\xe3\xf6\xb8\x25\xc7\x29\xfd\x92\x49\x54\xc3\x8a\x3b\x3e\x37\x8f\x5e\xde\xfc\x77\x6e\x6c\xcb\x65\xb2\xf5\x2d\xdb\xcc\x2a\x6f\x75\xef\xe3\x91\xe1\x1c\x73\x1f\x2c\x03\xff\x8c\x87\x0b\xbc\x42\x93\x90\xd7\x34\x3d\xc1\xb6\xb9\x08\x91\x5f\x66\xb4\xab\xef\x8c\xf1\x9d\x81\xc0\xf4\x74\xae\xf3\xc8\xa2\x2d\x37\x70\xee\x59\xf0\x8d\xd9\xaf\x79\x16\x79\x2c\x7f\xa8\x6e\x46\x13\xc6\xe3\x1a\x03\xe9\x88\x4c\x18\x5a\x99\xd2\x9b\x15\xc8\x74\xec\xcd\x69\x5b\xc2\x8f\x06\x82\x76\x30\x61\x13\xfb\x27\x00\x00\xff\xff\x79\x26\x17\x66\x5b\x09\x00\x00")

func templatesClient_rpc_goTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesClient_rpc_goTpl,
		"templates/client_rpc_go.tpl",
	)
}

func templatesClient_rpc_goTpl() (*asset, error) {
	bytes, err := templatesClient_rpc_goTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/client_rpc_go.tpl", size: 2395, mode: os.FileMode(420), modTime: time.Unix(1555592143, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesData_goTpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xa4\x53\x41\x6e\xdb\x30\x10\x3c\x73\x5f\xb1\xd0\xa1\x90\x02\x45\xbe\x0b\xc8\x29\xad\x8f\x69\xd0\xe6\x03\x1b\x7a\xa5\x10\xa1\x28\x81\xa4\xd2\x18\x04\xff\x5e\x90\x72\x6d\xa7\x86\x63\x17\xbd\x08\xe4\x6a\x76\x39\xc3\x19\x4e\x24\x5f\xa9\x67\x0c\xa1\x79\x5c\x96\x0f\x34\x70\x8c\x00\x6a\x98\x46\xeb\xb1\x04\x21\x49\xbe\x30\x16\xbd\xf2\x2f\xf3\x73\x23\xc7\x61\x35\x91\xb7\x4a\xbe\x0e\x66\xd5\x8f\xb7\xf9\x77\x01\x42\x8f\xfd\x07\x90\x53\x76\x9e\x1c\x9b\x95\x1e\x7b\x3b\xbb\x02\x40\x1c\xff\x26\xe7\x2d\x19\xf6\xab\x9e\x34\xbd\x6f\x57\x1b\xf2\x54\x7c\x0e\x19\x38\x1d\xeb\x0a\xa8\x00\xfc\x76\xca\xac\xd7\x4c\x7e\xb6\xfc\x68\xb9\x53\xef\x31\x7e\x25\x4f\x3f\x78\x1a\x51\x19\xcf\xb6\x23\xc9\x18\x40\xa4\xd9\x4d\x2a\x43\x08\xa8\x3a\x6c\xbe\xbd\xd3\x30\x69\x5e\xcf\x46\x26\xbd\x18\x23\x88\x10\xfe\x2e\xc7\x58\x56\xc8\xd6\x8e\x16\x42\xb8\x45\x36\x9b\x18\x21\x5e\x3c\xfb\xfb\xe4\xd5\x68\x1c\x3a\x6f\x67\xe9\x31\xa4\x96\x6e\x36\x12\xe5\x0b\xcb\xd7\x8b\x7d\xe5\x38\x79\xbc\xb9\x08\xab\xae\xc0\x24\xed\xaa\xc3\x34\xf0\xee\x0e\x8d\xd2\xa9\x20\xf2\x16\xbf\x5c\xec\x0e\x11\x44\x04\x61\xd9\xcf\xd6\xa4\x21\x7b\x21\x0f\xfc\xeb\x7c\x77\x09\xa2\xd7\xe3\x33\xe9\xec\xc3\xfe\xea\xeb\xc5\x86\xfb\x1c\xa6\x9b\x1c\x9a\x26\x6f\x6a\x10\xd7\x29\xae\xa1\xfa\xcc\xf1\xb0\xa7\x9a\xa4\xa5\xd2\xa3\x55\x6f\xe4\x77\x89\x4e\xca\x53\xb1\xc5\x03\xbb\x1a\x96\xeb\x68\xff\xc1\x9a\xaa\x06\x21\x3c\xf5\xae\xc5\x5d\x1e\x9b\x27\xea\x5d\x1a\x2f\x0a\x4d\x5b\xb6\x45\x8b\x88\x45\xce\x73\x9d\xab\x8e\xed\x9b\x92\x5c\xb4\x58\x84\xd0\xfc\x5c\x76\x0b\xab\x8c\x88\xe9\xd3\x29\xd6\x1b\xd7\xa2\x1e\xfb\x66\x9d\xd7\xff\x39\x12\xc4\xe1\xc2\x5b\xdc\x2f\xeb\x64\xea\x51\x88\x4f\x2e\xea\x90\xdb\xa3\x77\x03\x59\xf2\x07\xc5\xb0\xa3\x7c\xc4\x78\x71\xf2\x0a\x2b\xe1\x5c\x18\x12\xb3\xf3\x6f\x34\x87\xaf\xb4\xf9\x80\x13\xde\x39\x1c\x67\x5f\x70\x92\xf3\x87\x7d\x6a\xb5\x3e\x61\xee\x49\xeb\xd2\x36\x49\x5b\x05\xc2\x79\xf2\x6e\x6d\xb0\xbd\xc3\x53\xe8\x93\x1a\x94\xe9\x0f\xe0\x0d\x77\x6c\x71\xd7\x52\x56\xb0\x4f\x9f\x51\x1a\x62\x12\xc1\x66\x83\x31\xfe\x0e\x00\x00\xff\xff\xb1\xcb\x07\x55\x62\x05\x00\x00")

func templatesData_goTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesData_goTpl,
		"templates/data_go.tpl",
	)
}

func templatesData_goTpl() (*asset, error) {
	bytes, err := templatesData_goTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/data_go.tpl", size: 1378, mode: os.FileMode(420), modTime: time.Unix(1551527770, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesHandler_goTpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xa4\x54\x4b\x6b\xf3\x38\x14\x5d\xeb\xfe\x8a\x3b\x82\x19\xec\x92\xaa\xfb\x40\x56\xa5\x61\x06\x86\x4e\x99\x14\xba\x56\x95\x6b\x59\xd4\x91\x8d\x24\x67\x12\x84\xff\xfb\x20\xd9\x69\xfa\xc8\xa3\xf0\x6d\x62\x59\x39\xe7\xfa\x9c\xfb\xea\xa4\x7a\x93\x9a\x30\x46\xf1\x34\x1e\x1f\xe5\x86\x86\x01\xc0\x6c\xba\xd6\x05\x2c\x80\x71\x4b\xe1\xae\x0e\xa1\xe3\x00\x8c\x6b\x13\xea\xfe\x55\xa8\x76\x73\xa7\x8d\xbd\xd5\xad\x35\x2a\x9d\x38\xb0\xd7\x5e\x7b\x2b\x35\x7e\xc4\x4c\x77\x87\xe7\xad\x6e\x39\xb0\xa6\xfd\x0c\xf2\xc6\xf5\x9d\x27\x7b\xd7\xb4\xda\xf5\xfe\xeb\x77\xa4\x0f\x4e\x26\x11\x5a\x36\x72\xb7\x4f\x28\x6d\xac\xe6\x97\x51\x1b\x0a\xce\x28\xcf\xa1\x04\x08\xfb\x2e\x7b\x5c\x92\x0c\xbd\xa3\x27\x47\x95\xd9\x0d\xc3\x9f\xd2\xae\x1b\x72\x68\x6c\x20\x57\x49\x45\x18\x21\x46\x34\x15\x8a\x87\x9d\xdc\x74\x0d\x2d\x7b\xab\x52\x46\x70\x18\x80\xc5\xf8\xf5\x7a\x18\x0a\x85\x37\xda\x58\x71\xdf\xda\x40\xbb\x50\x26\x3e\xd9\x75\xc2\x0f\x00\x5b\xe9\xce\x7f\x76\xd5\x91\xba\x20\x6a\x81\x7f\xc4\x28\xa6\xb7\x27\x67\xb6\x32\x4c\xc5\x89\xc3\x35\x47\xff\x74\xc1\xb4\xd6\xa3\x0f\xae\x57\x01\x63\xd2\x52\xf5\x56\xa1\xaa\x49\xbd\x5d\xa3\x15\x6d\x17\xf0\xe6\x1a\xaa\xbc\x0e\xc1\x08\xcc\x54\x98\xc2\x2d\x16\x68\x4d\x93\x2e\x58\x7e\xcd\xe6\x2e\x92\xe3\x00\x6c\x00\xe6\x28\xf4\xce\xa6\x18\xef\x26\x1e\xe9\xbf\xb3\xe4\x02\x98\xdf\x9e\x4a\xeb\x8a\xdc\xd6\x28\x9a\x01\xfb\x91\xbb\x19\x94\x17\x6a\x13\xdf\x75\x9d\xad\xd1\x68\x74\xfe\xf3\x8c\x97\x33\x60\x2c\x48\xed\xe7\x38\xb5\xae\x78\x96\xda\xa7\x40\x8c\x37\x72\x4f\x8e\xcf\x11\x91\xd7\x23\x93\xcf\xf2\x1f\x7e\xb4\xc5\xe7\xc8\x63\x14\x93\xc9\x51\x42\x46\x0c\xe9\xa7\x32\xd4\xac\xfd\x1c\x9b\x56\x8b\x65\x3e\xff\x7a\x54\x60\x29\xd1\x89\x8b\xe8\xb7\x6a\x96\x8a\xf5\xa1\x2f\x4f\xe5\xe4\xd8\x8e\xb9\x46\xe7\x8b\x04\x39\x0f\x9f\xd2\x00\x93\x89\x0f\x1e\xc6\x4a\x5e\x2f\x65\x92\x75\x7e\xaa\x73\x47\x15\x75\x0e\x73\x52\x74\xee\x83\xab\x73\x9f\x4c\x1d\xf4\xfe\x4b\x69\x75\x26\xec\xbd\x6c\x9a\xa2\x16\xc9\x4d\x09\xac\xb2\x7f\xb7\x1a\xe7\x8b\xec\xe1\xc5\x84\x7a\xf4\x51\x4c\xfb\x6c\xbc\xb2\x45\x2d\x46\xab\x65\x09\xcc\x07\x19\xfc\xd2\x26\xd2\xf7\xe8\xcf\x66\x63\xac\x3e\xc6\x5f\x53\x45\x0e\x27\x4a\x51\x02\x30\x72\x2e\x51\x6b\xe1\xb7\x4a\x9c\xb4\x51\xe6\x09\x4d\xb8\xdf\x8e\x13\x9a\x85\x8a\x17\xe9\x6c\x41\xce\x95\xc0\x0e\xab\x5d\x3c\xb6\xc1\x54\xfb\xc3\xed\x77\x49\x0f\xce\xb5\xee\xa8\x88\x29\xb1\x0a\x2e\x8b\x0c\xa1\x13\xab\x20\x43\xef\xff\x4a\xdb\xd6\xca\x26\x15\x9b\x5c\x66\xcc\x90\x53\x7a\xce\xf1\xf7\x2d\x9f\xe1\x14\x7e\x1c\xb1\xb1\xb1\xde\xd7\xea\xff\x01\x00\x00\xff\xff\xbf\x91\xf8\x1f\xb2\x06\x00\x00")

func templatesHandler_goTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesHandler_goTpl,
		"templates/handler_go.tpl",
	)
}

func templatesHandler_goTpl() (*asset, error) {
	bytes, err := templatesHandler_goTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/handler_go.tpl", size: 1714, mode: os.FileMode(420), modTime: time.Unix(1543793167, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesHandler_rpc_goTpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xac\x54\x4d\x4f\xe3\x3c\x10\x3e\x7b\x7e\xc5\xc8\x87\x57\xed\xab\x6e\x72\x8f\xc4\x61\x97\x0f\x2d\x17\xa8\x80\xfb\xe2\x3a\x53\xc7\x4b\x62\x47\xf6\xa4\x50\x45\xf9\xef\x2b\xa7\x01\x5a\x09\x4a\x57\xda\x4b\x6b\xcf\xc7\x33\xf3\xcc\x33\x4e\x9e\xe3\xb9\x2f\x09\x0d\x39\x0a\x8a\xa9\xc4\xd5\x16\x8d\xaa\xd5\xcb\x36\xc3\x8b\x5b\xbc\xb9\x7d\xc0\xcb\x8b\xeb\x87\x0c\xf2\x1c\xbf\xd7\x35\xea\x4a\x39\x43\x11\x9b\x2e\x32\xae\x08\x4b\xef\x08\xad\x43\xdd\x45\xf6\x0d\xea\xda\x92\x63\xe4\x4a\x31\xc6\xca\x77\x75\x89\x64\xb9\xa2\x80\xd4\xac\xa8\x44\x1f\xf0\x39\xa8\x16\xb9\xb2\x31\x03\x68\x95\x7e\x52\x86\xb0\xef\xb3\xe5\xee\x78\xa3\x1a\x1a\x06\x00\xdb\xb4\x3e\x30\xce\x40\x48\x72\xda\x97\xd6\x99\xfc\x77\xf4\x4e\x82\x90\x8e\x38\xaf\x98\x5b\x09\x20\x56\x9d\x89\x4e\x19\x94\xc6\x72\xd5\xad\x32\xed\x9b\x7c\xb2\xbd\xfe\x7f\x33\x5e\x82\xa8\xfd\x61\x50\xb4\xa1\x6b\x23\xb9\xbc\xf6\x26\x74\x31\x61\xed\xbb\x55\xe4\xa0\x52\xa1\xdd\x30\x52\x94\xb1\xce\xc8\xe3\x51\x0d\x71\xb0\x3a\x4a\x98\x03\xf0\xb6\x1d\x89\x5d\x91\xe2\x2e\xd0\x32\xd0\xda\xbe\x0c\xc3\xdd\xf2\xfc\xa7\x72\x65\x4d\x01\xad\x63\x0a\x6b\xa5\x09\x7b\xe8\xfb\xec\xdd\x73\xfd\xea\xf8\xe1\xcb\xed\x30\xc0\x00\xb0\x51\xe1\x28\xd8\x7d\x4b\xfa\x78\xb5\x33\xfc\xef\xa0\xc8\x32\xd8\x8d\xe2\x69\xe0\xfd\x70\x42\xc3\xb7\x2d\x5b\xef\x22\x46\x0e\x9d\x66\xec\x53\x5f\xeb\xce\x69\xd4\x15\xe9\xa7\x13\x32\x67\xbe\x65\xfc\xff\x84\xc0\xf9\x49\x51\xd8\x83\xb0\x6b\x4c\xa0\x67\x67\xe8\x6c\x9d\x0c\x62\xbc\x8e\x5c\xbf\xca\xef\x07\x10\x03\x88\x40\xdc\x05\x97\x60\xde\x08\xdd\xd0\xf3\xb1\xfc\x19\x88\xb8\xf9\x68\xdc\xf7\x14\x36\x56\xd3\x02\xc4\xa9\x4c\x17\x30\x3f\x2e\x5b\xff\xd6\xe0\x31\xf9\x76\xbc\x8b\xbf\x52\x62\xbe\x00\x21\x58\x99\x58\xe0\xb4\xb7\xd9\x83\x32\x31\x61\x09\x59\xab\x2d\x05\x59\x20\xa2\x0c\xad\xfe\x55\xed\xb2\xe5\x62\x74\xc6\x1d\x4d\x59\xa0\xec\xfb\x6c\x22\xbd\xeb\x64\x8c\x18\xd2\xcf\xda\x52\x5d\xc6\x02\x6b\x6f\xb2\xab\xf1\xfc\x6f\x90\x41\xa4\xe1\x17\x18\x37\x7a\x91\xf4\xdb\xdb\xdc\x4f\xa6\xf3\xbe\xb0\xa3\x6c\x9f\xeb\x06\xe3\x38\x0e\xa6\x01\x13\x8f\x3d\x1a\x3b\x71\x4f\x52\x37\x35\x77\xf8\xb2\x9b\xb6\xa6\x86\x1c\xab\xe4\x9f\x9e\xf7\x67\x2f\xef\x32\x04\x1f\xee\x28\xb6\xde\x45\xda\x23\x31\xda\xd3\xdd\x3a\x83\x8f\xe9\xab\x58\x48\x4a\x36\xf9\xf8\xb6\xc1\xee\xa3\x0d\x3e\x00\x9c\x51\x08\x38\xa6\x7d\xb4\x81\x87\xb5\xdf\x97\xf0\x8b\xc8\x24\xf1\x68\x28\x12\x74\x36\x1e\x67\xf3\x49\xa7\x3f\x01\x00\x00\xff\xff\x79\xc2\xc0\xa5\x6d\x06\x00\x00")

func templatesHandler_rpc_goTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesHandler_rpc_goTpl,
		"templates/handler_rpc_go.tpl",
	)
}

func templatesHandler_rpc_goTpl() (*asset, error) {
	bytes, err := templatesHandler_rpc_goTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/handler_rpc_go.tpl", size: 1645, mode: os.FileMode(420), modTime: time.Unix(1543449660, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesService_goTpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xa4\x54\xc1\x6e\xdb\x30\x0c\x3d\x4b\x5f\x41\xf8\x30\xd8\x43\x62\xdf\x0d\xf4\xb4\x35\xc7\x2e\xd8\xfa\x03\xac\x42\xab\x42\x65\xc9\x90\xe4\x2e\x81\xa0\x7f\x1f\xe4\x38\x69\x83\x26\x73\xb1\x5d\x12\x89\x7a\xa4\xf8\x1e\x9f\x3c\xa0\x78\x41\x49\x10\x63\xbd\x3d\x2e\x1f\xb0\xa7\x94\x38\x57\xfd\x60\x5d\x80\x92\x33\x81\xe2\x99\xa0\x90\x2a\x3c\x8f\x4f\xb5\xb0\x7d\x33\x60\x70\x4a\xbc\xf4\xa6\x91\x76\x3d\x1d\x17\x9c\xb3\xa7\x51\x7a\x83\xf2\x02\x39\xc7\x4e\xff\x6b\x69\x0b\xce\xb4\xbd\x04\x79\xe5\xc6\xc1\x93\x69\xb4\x95\x6e\xf4\xb9\xd6\xfb\x63\xf4\xc1\xa1\xa1\xd0\x48\xd4\xb8\x3f\x34\x3d\xe5\xdb\x7d\xc1\x2b\xce\x9b\x46\xda\x56\x92\x21\x87\x81\xe0\x08\x00\xda\x0f\xd6\x13\xac\x0f\xb0\xde\x5e\x61\x16\x0e\xc3\x44\x78\x43\x18\x46\x47\x5b\x47\x9d\xda\xa7\xf4\x8b\xdc\xab\x12\x04\xca\x04\x72\x1d\x0a\x82\xc8\x63\x54\x1d\xd4\xf7\x7b\xec\x07\x4d\x9b\xd1\x88\xb9\x06\x8b\xf1\x63\xb4\xac\x80\x9c\xb3\x8e\xc7\x48\x66\x97\x12\x5f\xbc\xeb\xc7\x10\x94\x35\x1e\x7c\x70\xa3\x08\x10\x73\x46\x37\x1a\x01\xe2\x99\xc4\xcb\x52\x5a\x69\x87\x00\x5f\x97\x50\xd5\x32\x04\x22\x67\xaa\x83\x5c\xee\xee\x0e\x8c\xd2\x39\xc0\xa6\x2d\x7c\x59\x4a\x8e\x89\xb3\xc4\x99\xa3\x30\x3a\x93\x6b\x9c\x49\x3c\xd0\xef\x9b\xc9\x65\x96\x30\x6b\xfb\x93\x06\x7b\x6f\xf0\x49\xd3\x2e\x25\x47\x83\xbd\xa2\xd6\x77\x0c\x98\x71\xab\x93\xb0\xec\x53\xcc\x57\xbc\xfa\xcb\x98\xe3\xb9\xe7\x4c\x71\x8e\x6e\x9d\x7a\xc5\x30\x3b\x65\x16\xa1\xfd\xfc\x34\xaa\x15\x67\x2c\xa0\xf4\x2d\xcc\x26\xad\x1f\x51\xfa\x5c\x88\x15\x1a\x0f\xe4\x8a\x16\x00\x0a\x7f\xcc\x2c\x56\xd3\xc1\x69\xd7\x42\xf1\xd6\xc9\xb1\x85\x09\x91\xf2\x4f\xa7\x48\xef\x7c\x0b\xda\xca\x7a\x33\xad\xff\xbf\x2a\x67\xb7\x86\xd0\x82\xbb\xd0\xfb\xc2\xcb\xd7\xb4\x7a\xb3\xf0\x3f\xcc\xf5\x7c\x4d\x96\xee\x42\x39\x3e\xf3\x7e\x47\xfb\x38\xfc\xe5\xe9\xe7\x8e\x6f\x3d\xdf\xc9\x9e\xa5\x9f\x8a\x5c\x63\x33\xf9\xe6\xe6\xfb\xce\x24\x4f\x4d\xe6\xf6\x5d\xc8\x98\x6f\xa8\x75\xe9\xeb\x4c\xa1\xe2\xcc\x07\x0c\x7e\x63\xa0\xbd\x83\x8f\xd0\x47\xd5\x2b\x23\xdf\xc0\x3b\xea\xc8\xc1\x9c\x52\x56\xfc\xec\x4c\xa3\x34\x4f\xa7\xcf\xc9\x9f\x00\x00\x00\xff\xff\xf5\x0a\x11\xf9\xab\x05\x00\x00")

func templatesService_goTplBytes() ([]byte, error) {
	return bindataRead(
		_templatesService_goTpl,
		"templates/service_go.tpl",
	)
}

func templatesService_goTpl() (*asset, error) {
	bytes, err := templatesService_goTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/service_go.tpl", size: 1451, mode: os.FileMode(420), modTime: time.Unix(1543449663, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/client_rpc_go.tpl": templatesClient_rpc_goTpl,
	"templates/data_go.tpl": templatesData_goTpl,
	"templates/handler_go.tpl": templatesHandler_goTpl,
	"templates/handler_rpc_go.tpl": templatesHandler_rpc_goTpl,
	"templates/service_go.tpl": templatesService_goTpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"client_rpc_go.tpl": &bintree{templatesClient_rpc_goTpl, map[string]*bintree{}},
		"data_go.tpl": &bintree{templatesData_goTpl, map[string]*bintree{}},
		"handler_go.tpl": &bintree{templatesHandler_goTpl, map[string]*bintree{}},
		"handler_rpc_go.tpl": &bintree{templatesHandler_rpc_goTpl, map[string]*bintree{}},
		"service_go.tpl": &bintree{templatesService_goTpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

