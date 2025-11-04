package trello

type Plan struct {
	ProjectName   string      `json:"project_name"`
	ProductVision string      `json:"product_vision"`
	Epics         []Epic      `json:"epics"`
	UserStories   []UserStory `json:"user_stories"`
	Tasks         []Task      `json:"tasks"`
}

type Epic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserStory struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Story    string `json:"story"`
	Priority int    `json:"priority"`
	EpicID   int    `json:"epic_id"`
}

type Task struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	StoryID      int      `json:"story_id"`
	Dependencies []string `json:"dependencies"`
	Labels       []string `json:"labels"` // Added for Trello labels
}
