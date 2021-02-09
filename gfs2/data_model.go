package gfs2

type ClusterMetric map[string]FSMetric

type FSMetric map[string]Metric

type Metric struct {
	Value uint64
	State string
}
