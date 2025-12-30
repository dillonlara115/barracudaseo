package models

import "time"

// GA4Performance represents Google Analytics 4 performance data for a URL
type GA4Performance struct {
	URL                string    `json:"url"`
	Sessions           int64     `json:"sessions"`
	Users              int64     `json:"users"`
	PageViews          int64     `json:"page_views"`
	BounceRate         float64   `json:"bounce_rate"`
	AvgSessionDuration float64   `json:"avg_session_duration"` // in seconds
	Conversions        int64     `json:"conversions"`
	Revenue            float64   `json:"revenue"`
	LastUpdated        time.Time `json:"last_updated"`
}

// GA4Property represents a Google Analytics 4 property
type GA4Property struct {
	PropertyID   string `json:"property_id"`
	PropertyName string `json:"property_name"`
	DisplayName  string `json:"display_name"`
}

// GA4AuthState represents OAuth state for security
type GA4AuthState struct {
	State       string    `json:"state"`
	RedirectURL string    `json:"redirect_url"`
	CreatedAt   time.Time `json:"created_at"`
}
