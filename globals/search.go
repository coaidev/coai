package globals

import "sync"

var searchMutex sync.RWMutex

// SearchSettings is copied once per request so an admin update cannot mix one
// provider's endpoint or options with another configuration.
type SearchSettings struct {
	Provider, Endpoint, Engines, ImageProxy string
	Crop                                    bool
	CropLength, SafeSearch                  int
}

func GetSearchSettings() SearchSettings {
	searchMutex.RLock()
	defer searchMutex.RUnlock()
	return SearchSettings{
		Provider: SearchProvider, Endpoint: SearchEndpoint, Engines: SearchEngines,
		ImageProxy: SearchImageProxy, Crop: SearchCrop,
		CropLength: SearchCropLength, SafeSearch: SearchSafeSearch,
	}
}

func SetSearchSettings(settings SearchSettings) {
	searchMutex.Lock()
	defer searchMutex.Unlock()
	SearchProvider, SearchEndpoint = settings.Provider, settings.Endpoint
	SearchEngines, SearchImageProxy = settings.Engines, settings.ImageProxy
	SearchCrop, SearchCropLength, SearchSafeSearch = settings.Crop, settings.CropLength, settings.SafeSearch
}
