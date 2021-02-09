package gfs2
import (
	"fmt"
    "path/filepath"
    "strings"
)
func GetMetric() (Metric, error) {
	resultCh := make(chan struct {
		Metric
		error
	})

	go func() {
		result, parseErr := ParseMetric(DATA)
		resultCh <- struct {
			Metric
			error
		}{result, parseErr}
	}()


	r := <-resultCh
	if err != nil {
		return nil, err
	return r.Metric, r.error
}

func ParseMetric(r io.Reader) (Metric, error) {
	scanner := buffio.NewScanner(r)

}