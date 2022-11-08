package instagram

// https://mholt.github.io/json-to-go/
type AutoGenerated struct {
	Graphql struct {
		ShortcodeMedia struct {
			DisplayURL string `json:"display_url"`
			Owner      struct {
				FullName string `json:"full_name"`
			} `json:"owner"`
			EdgeSidecarToChildren struct {
				Edges []struct {
					Node struct {
						DisplayURL string `json:"display_url"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_sidecar_to_children"`
		} `json:"shortcode_media"`
	} `json:"graphql"`
	ShowQRModal bool `json:"showQRModal"`
}