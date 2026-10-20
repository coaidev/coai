# Web search

Enable web search in a chat to add current sources to the model's context. New
installations use [Parallel Search MCP](https://docs.parallel.ai/integrations/mcp/search-mcp)
without an account or API key. Anonymous search uses Fast mode. When the query
consists of a single HTTP or HTTPS URL, Parallel fetches that page instead of
searching. The `web-` model prefix on the OpenAI-compatible API uses the same
search provider.

Choose a provider under **Admin → System Settings → Search** and save before
using the search test:

| Provider | Behavior |
| --- | --- |
| Automatic | Uses SearXNG when a search endpoint is configured; otherwise uses Parallel. |
| Parallel Search MCP | Uses `https://search.parallel.ai/mcp`, even if a SearXNG endpoint is saved. |
| SearXNG | Uses the configured search endpoint; an endpoint is required. |

Existing configurations without a provider keep using their saved SearXNG
endpoint. This includes the endpoint in older configuration templates; select
Parallel explicitly to switch an existing installation. A provider error never
silently switches to another provider. Turning web search off in a chat prevents
both search and page fetching.

The equivalent configuration in `config/config.yaml` is:

```yaml
system:
  search:
    provider: auto # auto, parallel, or searxng
    endpoint: "" # SearXNG endpoint; empty selects Parallel in auto mode
```

`SYSTEM_SEARCH_PROVIDER` overrides the provider. SearXNG engine selection,
safe-search and image-proxy settings apply only to SearXNG; select that provider
when those controls are required. Result cropping applies to both providers.

Parallel requests have a 30-second total deadline and honor request cancellation.
The anonymous service has usage limits. If search is unavailable, chat continues
without adding search context; the admin search test reports the error. Search
and fetched page content are temporary model context and are not saved as chat
history.
