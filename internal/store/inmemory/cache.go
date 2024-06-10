package inmemory

import (
	"errors"

	"github.com/dgraph-io/ristretto"
)

var ErrParsing = errors.New("parsing error")

type Config struct {
	NumCounters  int64 `yaml:"num_counters" mapstructure:"num_counters"  default:"10000000"`
	MaxCost      int64 `yaml:"max_cost" mapstructure:"max_cost"  default:"1073741824"`
	BufferItems  int64 `yaml:"buffer_items" mapstructure:"buffer_items"  default:"64"`
	Metrics      bool  `yaml:"metrics" mapstructure:"metrics"  default:"true"`
	TTLInSeconds int   `yaml:"ttl_in_seconds" mapstructure:"ttl_in_seconds"  default:"3600"`
}

type Cache struct {
	*ristretto.Cache
	config Config
}

func NewCache(cfg Config) (Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.NumCounters,
		MaxCost:     cfg.MaxCost,
		BufferItems: cfg.BufferItems,
		Metrics:     cfg.Metrics,
	})
	if err != nil {
		return Cache{}, err
	}

	return Cache{
		Cache:  cache,
		config: cfg,
	}, nil
}
