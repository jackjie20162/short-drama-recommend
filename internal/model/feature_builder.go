package model

import (
	"strings"
	"time"
)

func BuildDramaSemanticText(d DramaContent) string {
	parts := []string{
		d.Title,
		d.Subtitle,
		d.Description,
		strings.Join(d.Genres, ", "),
		strings.Join(d.Tags, ", "),
		d.Language,
		d.Country,
	}
	return strings.TrimSpace(strings.Join(parts, " | "))
}

func BuildDramaDenseFeature(d DramaContent, now time.Time) DenseFeature {
	ageDays := 0.0
	if !d.PublishedAt.IsZero() && now.After(d.PublishedAt) {
		ageDays = now.Sub(d.PublishedAt).Hours() / 24
	}
	return DenseFeature{
		DramaPopularity: d.Popularity,
		DramaCompletion: d.CompletionRate,
		DramaPayRate:    d.PayRate,
		DramaAgeDays:    ageDays,
	}
}
