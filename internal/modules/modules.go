// Package modules defines the navigation targets of the application.
// Titles are user-facing and therefore German.
package modules

// Status tells whether an area is usable or only reserved in the navigation.
type Status string

const (
	Active  Status = "active"
	Planned Status = "planned"
)

// Module is one navigation target.
type Module struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status Status `json:"status"`
}

// All returns the navigation targets in display order.
func All() []Module {
	return []Module{
		{ID: "rhymes", Title: "Reimarchiv", Status: Active},
		{ID: "lines", Title: "Reimzeilen", Status: Planned},
		{ID: "generator", Title: "Generator", Status: Planned},
		{ID: "settings", Title: "Einstellungen", Status: Active},
	}
}
