package ng

import "time"

type Skeleton struct {
	App         string      `json:"app"`
	Version     Version     `json:"version"`
	Include   []Include     `json:"include"`
	HotReload   bool        `json:"hot_reload"`
	Debug       bool        `json:"debug"`
	Config    []string      `json:"config"`
}

type Version struct {
	Auto  bool      `json:"auto"`
	Major int		`json:"major"`
	Minor int		`json:"minor"`
	Date  time.Time	`json:"date"`
}

type Include struct {
	Repo   *string `json:"repo,omitempty"`
	Branch *string `json:"branch,omitempty"`
	Local  *string `json:"local,omitempty"`
}

