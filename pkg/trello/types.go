package trello

type Plan struct {
	ProjectName   string      `json:"project_name"`
	ProductVision string      `json:"product_vision"`
	Epics         []Epic      `json:"epics"`
	UserStories   []UserStory `json:"user_stories"`
	Tasks         []Task      `json:"tasks"`
}

type Task struct {
	ID           string   `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	StoryID      string   `json:"story_id"`
	Dependencies []string `json:"dependencies"`
	Labels       []string `json:"labels"`
}

type Epic struct {
	ID   string    `json:"id"`
	Name string `json:"name"`
}

type UserStory struct {

	ID        string    `json:"id"`

	Title     string `json:"title"`

	Story     string `json:"story"`

	Priority  int    `json:"priority"`

	EpicID    string    `json:"epic_id"`

	Checklist LocalChecklist `json:"checklist"`

}



type LocalChecklist struct {

	Name  string   `json:"name"`

	Items []string `json:"items"`

}



// Checklist represents a Trello checklist.

// https://developers.trello.com/reference/#checklist-object

type Checklist struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	IDBoard    string       `json:"idBoard,omitempty"`
	IDCard     string       `json:"idCard,omitempty"`
	Pos        float32      `json:"pos,omitempty"`
	CheckItems []*CheckItem `json:"checkItems,omitempty"`
}

// CheckItem represents an item within a Trello checklist.
// https://developers.trello.com/reference/#checkitem-object
type CheckItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	State       string  `json:"state"` // "complete" or "incomplete"
	Pos         float32 `json:"pos,omitempty"`
	IDChecklist string  `json:"idChecklist,omitempty"`
}


