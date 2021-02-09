package gfs2

import (
	"fmt"
    "path/filepath"
    "strings"
)
func GetMetric() (ClusterMetric, error) {
	resultCh := make(chan struct {
		ClusterMetric
		error
	})

	go func() {
		result, parseErr := ParseMetric(DATA)
		resultCh <- struct {
			ClusterMetric
			error
		}{result, parseErr}
	}()


	r := <-resultCh
	if err != nil {
		return nil, err
	return r.ClusterMetric, r.error
}

func ParseMetric(r io.Reader) (Metric, error) {
	scanner := buffio.NewScanner(r)

}