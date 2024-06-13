package inmemory

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func (c Cache) MonitorCache(meter metric.Meter) error {
	metrics := []string{
		"cache.hits", "cache.miss", "cache.cost_added",
		"cache.cost_evicted", "cache.gets_dropped", "cache.gets_kept",
		"cache.keys_added", "cache.keys_evicted", "cache.keys_updated",
		"cache.sets_dropped", "cache.sets_rejected",
	}

	int64Intruments := map[string]metric.Int64ObservableCounter{}
	for _, m := range metrics {
		inst, err := meter.Int64ObservableCounter(m)
		if err != nil {
			otel.Handle(err)
		}
		int64Intruments[m] = inst
	}

	ratio, err := meter.Float64ObservableGauge("cache.ratio")
	if err != nil {
		otel.Handle(err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			o.ObserveInt64(int64Intruments["cache.hits"], int64(c.Metrics.Hits()))
			o.ObserveInt64(int64Intruments["cache_miss"], int64(c.Metrics.Misses()))
			o.ObserveInt64(int64Intruments["cache.cost_added"], int64(c.Metrics.CostAdded()))
			o.ObserveInt64(int64Intruments["cache.cost_evicted"], int64(c.Metrics.CostEvicted()))
			o.ObserveInt64(int64Intruments["cache.gets_dropped"], int64(c.Metrics.GetsDropped()))
			o.ObserveInt64(int64Intruments["cache.gets_kept"], int64(c.Metrics.GetsKept()))
			o.ObserveInt64(int64Intruments["cache.keys_added"], int64(c.Metrics.KeysAdded()))
			o.ObserveInt64(int64Intruments["cache.keys_evicted"], int64(c.Metrics.KeysEvicted()))
			o.ObserveInt64(int64Intruments["cache.keys_updated"], int64(c.Metrics.KeysUpdated()))
			o.ObserveInt64(int64Intruments["cache.sets_dropped"], int64(c.Metrics.SetsDropped()))
			o.ObserveInt64(int64Intruments["cache.sets_rejected"], int64(c.Metrics.SetsRejected()))
			o.ObserveFloat64(ratio, c.Metrics.Ratio())

			return nil
		},
		int64Intruments["cache.hits"], int64Intruments["cache_miss"], ratio,
		int64Intruments["cache.cost_added"], int64Intruments["cache.cost_evicted"], int64Intruments["cache.gets_dropped"],
		int64Intruments["cache.gets_kept"], int64Intruments["cache.keys_added"], int64Intruments["cache.keys_evicted"],
		int64Intruments["cache.keys_evicted"], int64Intruments["cache.sets_dropped"], int64Intruments["cache.sets_rejected"],
	)
	if err != nil {
		return err
	}

	return nil
}
