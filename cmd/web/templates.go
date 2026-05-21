package main

import "github.com/kyzercmd/snipio/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
