package trello

type Plan struct {
	ProjectName   string      `json:"project_name"`
	ProductVision string      `json:"product_vision"`
	Epics         []Epic      `json:"epics"`
	UserStories   []UserStory `json:"user_stories"`
	Tasks         []Task      `json:"tasks"`
}

type Epic struct {
	Name string `json:"name"`
}

type UserStory struct {
	Title    string `json:"title"`
	Story    string `json:"story"`
	Priority int    `json:"priority"`
	Epic     string `json:"epic"`
}

type Task struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	StoryTitle   string   `json:"story_title"`
	Dependencies []string `json:"dependencies"`
	Labels       []string `json:"labels"` // Added for Trello labels
}
