package inmemory

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func (c Cache) MonitorCache(meter metric.Meter) error {
	hits, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.hits")
	if err != nil {
		otel.Handle(err)
	}

	miss, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.miss")
	if err != nil {
		otel.Handle(err)
	}

	ratio, err := meter.Float64ObservableGauge("shield.internal.store.inmemory.cache.ratio")
	if err != nil {
		otel.Handle(err)
	}

	costAdded, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.cost_added")
	if err != nil {
		otel.Handle(err)
	}

	costEvicted, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.cost_evicted")
	if err != nil {
		otel.Handle(err)
	}

	getsDropped, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.gets_dropped")
	if err != nil {
		otel.Handle(err)
	}

	getsKept, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.gets_kept")
	if err != nil {
		otel.Handle(err)
	}

	keysAdded, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.keys_added")
	if err != nil {
		otel.Handle(err)
	}

	keysEvicted, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.keys_evicted")
	if err != nil {
		otel.Handle(err)
	}

	keysUpdated, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.keys_updated")
	if err != nil {
		otel.Handle(err)
	}

	setsDropped, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.sets_dropped")
	if err != nil {
		otel.Handle(err)
	}

	setsRejected, err := meter.Int64ObservableCounter("shield.internal.store.inmemory.cache.sets_rejected")
	if err != nil {
		otel.Handle(err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			o.ObserveInt64(hits, int64(c.Metrics.Hits()))
			o.ObserveInt64(miss, int64(c.Metrics.Misses()))
			o.ObserveFloat64(ratio, c.Metrics.Ratio())
			o.ObserveInt64(costAdded, int64(c.Metrics.CostAdded()))
			o.ObserveInt64(costEvicted, int64(c.Metrics.CostEvicted()))
			o.ObserveInt64(getsDropped, int64(c.Metrics.GetsDropped()))
			o.ObserveInt64(getsKept, int64(c.Metrics.GetsKept()))
			o.ObserveInt64(keysAdded, int64(c.Metrics.KeysAdded()))
			o.ObserveInt64(keysEvicted, int64(c.Metrics.KeysEvicted()))
			o.ObserveInt64(keysUpdated, int64(c.Metrics.KeysUpdated()))
			o.ObserveInt64(setsDropped, int64(c.Metrics.SetsDropped()))
			o.ObserveInt64(setsRejected, int64(c.Metrics.SetsRejected()))

			return nil
		},
		hits, miss, ratio, costAdded, costEvicted, getsDropped, getsKept,
		keysAdded, keysEvicted, keysUpdated, setsDropped, setsRejected,
	)
	if err != nil {
		return err
	}

	return nil
}
