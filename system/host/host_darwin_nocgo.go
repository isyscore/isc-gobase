//go:build darwin && !cgo

package host

import (
	"context"
	"github.com/isyscore/isc-gobase/system/common"
)

func SensorsTemperaturesWithContext(ctx context.Context) ([]TemperatureStat, error) {
	return []TemperatureStat{}, common.ErrNotImplementedError
}
