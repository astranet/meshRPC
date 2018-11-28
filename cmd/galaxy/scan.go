package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var metricTagsStart = []byte(`metrics.Tags{`)
var metricTagsEnd = []byte(`}`)

func anyServiceIn(set map[string]map[string]bool) string {
	for k, v := range set {
		if v["service"] {
			return k
		}
	}
	return ""
}

func scanMetricTagsAll(servicesRoot string) map[string]map[string]bool {
	set := make(map[string]map[string]bool)
	if err := filepath.Walk(servicesRoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		} else if filepath.Ext(path) != ".go" {
			return nil
		}
		return scanMetricTagsFile(set, path)
	}); err != nil {
		log.Fatalln("Error scanning services:", err)
	}
	return set
}

func scanMetricTagsFile(set map[string]map[string]bool, path string) error {
	if filepath.Ext(path) != ".go" {
		return errors.New("not a Go file")
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	start := bytes.Index(data, metricTagsStart)
	if start < 0 {
		return nil
	}
	end := start + bytes.Index(data[start:], metricTagsEnd)
	if end < 0 {
		return nil
	}
	buf := data[start : end+1]
	if !bytes.Contains(buf, metricTagsStart) {
		return nil
	}
	tags, err := parseMetricTags(string(buf))
	if err != nil {
		log.Printf("Warning: failed to parse Tags of %s: %v", path, err)
		return nil
	}
	serviceTag := tags["service"]
	if len(serviceTag) == 0 {
		log.Printf("Warning: no service tag in Tags of %s", path)
		return nil
	}
	if set[serviceTag] == nil {
		set[serviceTag] = make(map[string]bool)
	}
	layerTag := tags["layer"]
	if len(layerTag) == 0 {
		log.Printf("Warning: no layer tag in Tags of %s", path)
		return nil
	}
	set[serviceTag][layerTag] = true
	return nil
}

var tagsMapRx = regexp.MustCompile(`"(\w+)":\s+"(\w+)"`)

func parseMetricTags(spec string) (map[string]string, error) {
	tags := make(map[string]string)
	matches := tagsMapRx.FindAllStringSubmatch(spec, -1)
	for _, m := range matches {
		if len(m) != 3 {
			continue
		}
		tags[m[1]] = m[2]
	}
	return tags, nil
}
