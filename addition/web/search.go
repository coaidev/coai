package web

import (
	"chat/globals"
	"chat/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearXNGResponse struct {
	Query           string `json:"query"`
	NumberOfResults int    `json:"number_of_results"`
	Results         []struct {
		Url           string   `json:"url"`
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		PublishedDate *string  `json:"publishedDate,omitempty"`
		Thumbnail     *string  `json:"thumbnail,omitempty"`
		Engine        string   `json:"engine"`
		ParsedUrl     []string `json:"parsed_url"`
		Template      string   `json:"template"`
		Engines       []string `json:"engines"`
		Positions     []int    `json:"positions"`
		Score         float64  `json:"score"`
		Category      string   `json:"category"`
		IframeSrc     string   `json:"iframe_src,omitempty"`
	} `json:"results"`
	Answers             []interface{} `json:"answers"`
	Corrections         []interface{} `json:"corrections"`
	Infoboxes           []interface{} `json:"infoboxes"`
	Suggestions         []interface{} `json:"suggestions"`
	UnresponsiveEngines [][]string    `json:"unresponsive_engines"`
}

func formatResponse(data *SearXNGResponse) string {
	res := make([]string, 0)
	for _, item := range data.Results {
		if item.Content == "" || item.Url == "" || item.Title == "" {
			continue
		}

		res = append(res, fmt.Sprintf("%s (%s): %s", item.Title, item.Url, item.Content))
	}

	return strings.Join(res, "\n")
}

func createURLParams(query string, settings globals.SearchSettings) string {
	params := url.Values{}

	params.Add("q", query)
	params.Add("format", "json")
	params.Add("safesearch", strconv.Itoa(settings.SafeSearch))
	if len(settings.Engines) > 0 {
		params.Add("engines", settings.Engines)
	}
	if len(settings.ImageProxy) > 0 {
		params.Add("image_proxy", settings.ImageProxy)
	}

	return fmt.Sprintf("%s?%s", settings.Endpoint, params.Encode())
}

func createSearXNGRequest(query string, settings globals.SearchSettings) (*SearXNGResponse, error) {
	data, err := utils.Get(createURLParams(query, settings), nil)

	if err != nil {
		return nil, err
	}

	return utils.MapToRawStruct[SearXNGResponse](data)
}

func GenerateSearchResult(q string) (string, error) {
	return GenerateSearchResultContext(context.Background(), q)
}

func GenerateSearchResultContext(ctx context.Context, q string) (string, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return "", errors.New("search query is empty")
	}
	settings := globals.GetSearchSettings()
	provider := strings.ToLower(strings.TrimSpace(settings.Provider))
	if provider == "" || provider == "auto" {
		if strings.TrimSpace(settings.Endpoint) != "" {
			provider = "searxng"
		} else {
			provider = "parallel"
		}
	}
	var content string
	var err error
	switch provider {
	case "parallel":
		content, err = parallelSearch(ctx, q)
	case "searxng":
		if strings.TrimSpace(settings.Endpoint) == "" {
			return "", errors.New("SearXNG requires a search endpoint")
		}
		var res *SearXNGResponse
		res, err = createSearXNGRequest(q, settings)
		if err == nil {
			content = formatResponse(res)
		}
	default:
		return "", fmt.Errorf("unknown search provider: %s", provider)
	}
	if err != nil {
		globals.Warn(fmt.Sprintf("[web] failed to get search result: %s (query: %s)", err.Error(), utils.Extract(q, 20, "...")))

		content := fmt.Sprintf("search failed: %s", err.Error())
		return content, errors.New(content)
	}

	globals.Debug(fmt.Sprintf("[web] search result: %s (query: %s)", utils.Extract(content, 50, "..."), q))

	if settings.Crop {
		globals.Debug(fmt.Sprintf("[web] crop search result length %d to %d max", len(content), settings.CropLength))
		return utils.Extract(content, settings.CropLength, "..."), nil
	}
	return content, nil
}

func TestSearch(c *gin.Context) {
	// get `query` param from query
	query := c.Query("query")

	fmt.Println(query)

	res, err := GenerateSearchResultContext(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"error":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"result": res,
		})
	}
}
