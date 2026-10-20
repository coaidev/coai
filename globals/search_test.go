package globals

import (
	"sync"
	"testing"
)

func TestSearchSettingsStayConsistentDuringUpdates(t *testing.T) {
	previous := GetSearchSettings()
	defer SetSearchSettings(previous)
	first := SearchSettings{Provider: "parallel", Endpoint: "first", CropLength: 100}
	second := SearchSettings{Provider: "searxng", Endpoint: "second", CropLength: 200}
	SetSearchSettings(first)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			SetSearchSettings(first)
			SetSearchSettings(second)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			if settings := GetSearchSettings(); settings != first && settings != second {
				t.Errorf("mixed search settings: %+v", settings)
			}
		}
	}()
	wg.Wait()
}
