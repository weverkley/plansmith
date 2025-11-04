package trello

import (
	"fmt"

	"github.com/adlio/trello"
)

type Client struct {
	client *trello.Client
}

func NewClient(key, token string) *Client {
	client := trello.NewClient(key, token)
	return &Client{
		client: client,
	}
}

func (c *Client) CreateProjectBoard(name string) (*trello.Board, error) {
	// Create a new board object
	board := trello.NewBoard(name)
	board.Desc = "Project board created by PlanSmith"

	// Use the client to create the board
	err := c.client.CreateBoard(&board, trello.Defaults())
	if err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}

	// Get the lists that were automatically created by Trello
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return nil, fmt.Errorf("failed to get lists: %w", err)
	}

	// Archive the default lists
	for _, list := range lists {
		err := list.Archive()
		if err != nil {
			return nil, fmt.Errorf("failed to archive list %s: %w", list.Name, err)
		}
	}

	// Create default lists
	listNames := []string{"Epics", "User Stories (Backlog)", "To Do (Ready)", "Doing", "Done"}
	for _, listName := range listNames {
		// Use the client to create the list with the correct API
		_, err := c.client.CreateList(&board, listName, trello.Defaults())
		if err != nil {
			return nil, fmt.Errorf("failed to create list %s: %w", listName, err)
		}
	}

	return &board, nil
}

func (c *Client) PopulateBoard(boardID string, plan *Plan) error {
	board, err := c.client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		return fmt.Errorf("failed to get board: %w", err)
	}

	// Get lists
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}

	// Map list names to IDs
	listMap := make(map[string]*trello.List)
	for _, list := range lists {
		listMap[list.Name] = list
	}

	// --- Label Management ---
	// Fetch existing labels
	existingLabels, err := board.GetLabels(trello.Defaults())
	if err != nil {
		return fmt.Errorf("failed to get existing labels: %w", err)
	}
	labelMap := make(map[string]*trello.Label)
	for _, label := range existingLabels {
		labelMap[label.Name] = label
	}

	// Function to get or create a label
	getOrCreateLabel := func(name, color string) (*trello.Label, error) {
		if label, ok := labelMap[name]; ok {
			return label, nil
		}
		newLabel := &trello.Label{Name: name, Color: color}
		err := board.CreateLabel(newLabel, trello.Defaults())
		if err != nil {
			return nil, fmt.Errorf("failed to create label %s: %w", name, err)
		}
		labelMap[name] = newLabel
		return newLabel, nil
	}

	// Create/get Epic labels
	epicLabels := make(map[string]*trello.Label)
	for _, epic := range plan.Epics {
		label, err := getOrCreateLabel(epic.Name, "sky") // Assign a color, e.g., "sky"
		if err != nil {
			return err
		}
		epicLabels[epic.Name] = label
	}

	// Create/get Priority labels
	priorityLabels := make(map[int]*trello.Label)
	priorityColors := map[int]string{
		10: "red", // High priority
		9:  "orange",
		8:  "yellow",
		7:  "green", // Low priority
	}
	for p := 7; p <= 10; p++ {
		labelName := fmt.Sprintf("P%d", p)
		color := priorityColors[p]
		label, err := getOrCreateLabel(labelName, color)
		if err != nil {
			return err
		}
		priorityLabels[p] = label
	}
	// --- End Label Management ---

	// Add epics to "Epics" list
	epicsList := listMap["Epics"]
	for _, epic := range plan.Epics {
		card := &trello.Card{
			Name:   epic.Name,
			Desc:   plan.ProductVision,
			IDList: epicsList.ID,
		}
		if label, ok := epicLabels[epic.Name]; ok {
			card.IDLabels = []string{label.ID}
		}
		err := c.client.CreateCard(card, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to create epic card: %w", err)
		}
	}

	// Add user stories to "User Stories (Backlog)" list
	storiesList := listMap["User Stories (Backlog)"]
	for _, story := range plan.UserStories {
		card := &trello.Card{
			Name:   story.Title,
			Desc:   story.Story,
			IDList: storiesList.ID,
		}
		// Add Epic label
		var epicName string
		for _, epic := range plan.Epics {
			if epic.ID == story.EpicID {
				epicName = epic.Name
				break
			}
		}
		if label, ok := epicLabels[epicName]; ok {
			card.IDLabels = append(card.IDLabels, label.ID)
		}
		// Add Priority label
		if label, ok := priorityLabels[story.Priority]; ok {
			card.IDLabels = append(card.IDLabels, label.ID)
		}
		err := c.client.CreateCard(card, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to create story card: %w", err)
		}
	}

	// Add tasks to "To Do (Ready)" list
	tasksList := listMap["To Do (Ready)"]
	for _, task := range plan.Tasks {
		card := &trello.Card{
			Name:   task.Title,
			Desc:   task.Description,
			IDList: tasksList.ID,
		}
		// Add labels from task.Labels field
		for _, labelName := range task.Labels {
			if label, ok := labelMap[labelName]; ok { // Check if label already exists or was created
				card.IDLabels = append(card.IDLabels, label.ID)
			} else {
				// If label doesn't exist, create it with a default color (e.g., "purple")
				newLabel, err := getOrCreateLabel(labelName, "purple")
				if err != nil {
					return err
				}
				card.IDLabels = append(card.IDLabels, newLabel.ID)
			}
		}

		// Find the corresponding UserStory to get its Epic and Priority
		var parentStory *UserStory
		for _, story := range plan.UserStories {
			if story.ID == task.StoryID {
				parentStory = &story
				break
			}
		}
		if parentStory != nil {
			// Add Epic label from parent story
			var epicName string
			for _, epic := range plan.Epics {
				if epic.ID == parentStory.EpicID {
					epicName = epic.Name
					break
				}
			}
			if label, ok := epicLabels[epicName]; ok {
				card.IDLabels = append(card.IDLabels, label.ID)
			}
			// Add Priority label from parent story
			if label, ok := priorityLabels[parentStory.Priority]; ok {
				card.IDLabels = append(card.IDLabels, label.ID)
			}
		}

		err := c.client.CreateCard(card, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to create task card: %w", err)
		}
	}

	return nil
}

func (c *Client) LinkCards(boardID string, plan *Plan) error {
	// This would implement linking cards based on dependencies
	// Implementation would involve finding cards and creating attachments/links
	// between them based on the task dependencies
	return nil
}
