package thirdparty

import (
	"testing"

	"github.com/erpc/erpc/common"
	"github.com/stretchr/testify/assert"
)

func TestQuicknodeFilterParams(t *testing.T) {
	vendor := CreateQuicknodeVendor().(*QuicknodeVendor)

	t.Run("extracts both tag IDs and labels", func(t *testing.T) {
		settings := common.VendorSettings{
			"apiKey":    "test-key",
			"tagIds":    []interface{}{123, 456},
			"tagLabels": []interface{}{"production", "staging"},
		}
		result := vendor.extractFilterParams(settings)
		assert.Equal(t, []int{123, 456}, result.TagIDs)
		assert.Equal(t, []string{"production", "staging"}, result.TagLabels)
	})

	t.Run("handles empty settings", func(t *testing.T) {
		settings := common.VendorSettings{"apiKey": "test-key"}
		result := vendor.extractFilterParams(settings)
		assert.Empty(t, result.TagIDs)
		assert.Empty(t, result.TagLabels)
	})
}

func TestQuicknodeCacheKey(t *testing.T) {
	vendor := CreateQuicknodeVendor().(*QuicknodeVendor)

	t.Run("no filters returns just API key", func(t *testing.T) {
		params := &QuicknodeFilterParams{}
		cacheKey := vendor.getCacheKey("test-api-key", params)
		assert.Equal(t, "test-api-key", cacheKey)
	})

	t.Run("nil params returns just API key", func(t *testing.T) {
		cacheKey := vendor.getCacheKey("test-api-key", nil)
		assert.Equal(t, "test-api-key", cacheKey)
	})

	t.Run("with tag IDs only", func(t *testing.T) {
		params := &QuicknodeFilterParams{
			TagIDs: []int{123, 456},
		}
		cacheKey := vendor.getCacheKey("test-api-key", params)
		assert.Equal(t, "test-api-key_tid:123,456", cacheKey)
	})

	t.Run("with tag labels only", func(t *testing.T) {
		params := &QuicknodeFilterParams{
			TagLabels: []string{"production", "staging"},
		}
		cacheKey := vendor.getCacheKey("test-api-key", params)
		assert.Equal(t, "test-api-key_tl:production,staging", cacheKey)
	})

	t.Run("with both tag IDs and labels", func(t *testing.T) {
		params := &QuicknodeFilterParams{
			TagIDs:    []int{789, 123},
			TagLabels: []string{"prod", "dev"},
		}
		cacheKey := vendor.getCacheKey("test-api-key", params)
		assert.Equal(t, "test-api-key_tid:123,789_tl:dev,prod", cacheKey)
	})

	t.Run("sorts tag IDs for consistent keys", func(t *testing.T) {
		params1 := &QuicknodeFilterParams{TagIDs: []int{456, 123, 789}}
		params2 := &QuicknodeFilterParams{TagIDs: []int{789, 123, 456}}

		key1 := vendor.getCacheKey("test-api-key", params1)
		key2 := vendor.getCacheKey("test-api-key", params2)

		assert.Equal(t, key1, key2, "cache keys should be identical regardless of tag ID order")
		assert.Equal(t, "test-api-key_tid:123,456,789", key1)
	})

	t.Run("sorts tag labels for consistent keys", func(t *testing.T) {
		params1 := &QuicknodeFilterParams{TagLabels: []string{"staging", "production", "dev"}}
		params2 := &QuicknodeFilterParams{TagLabels: []string{"dev", "staging", "production"}}

		key1 := vendor.getCacheKey("test-api-key", params1)
		key2 := vendor.getCacheKey("test-api-key", params2)

		assert.Equal(t, key1, key2, "cache keys should be identical regardless of tag label order")
		assert.Equal(t, "test-api-key_tl:dev,production,staging", key1)
	})

	t.Run("different filters produce different keys", func(t *testing.T) {
		params1 := &QuicknodeFilterParams{TagLabels: []string{"project-a"}}
		params2 := &QuicknodeFilterParams{TagLabels: []string{"project-b"}}

		key1 := vendor.getCacheKey("same-api-key", params1)
		key2 := vendor.getCacheKey("same-api-key", params2)

		assert.NotEqual(t, key1, key2, "different filters should produce different cache keys")
		assert.Equal(t, "same-api-key_tl:project-a", key1)
		assert.Equal(t, "same-api-key_tl:project-b", key2)
	})
}
