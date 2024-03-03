package model

// Organization one to many Repository
type Organization struct {
	Model

	Login  string
	NodeID string
}
