package response

type RateLimiterResponse struct {
	Allow    bool
	Metadata map[string]interface{}
}
