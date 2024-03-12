package common

import (
	"time"

	"lazy-auth/app/model"
)

func TernaryIf[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func BuildMetaPagination(total, limit, offset *int) model.MetaPagination {
	meta := model.MetaPagination{
		Total:  0,
		Offset: 0,
	}

	if total != nil {
		meta.Total = *total
	}

	meta.Limit = meta.Total
	if limit != nil {
		meta.Limit = *limit
	}

	if offset != nil {
		meta.Offset = *offset
	}

	return meta
}

func AddTimeByDuration(durationStr string) time.Time {
	duration, _ := time.ParseDuration(durationStr)
	return time.Now().Add(duration)
}
