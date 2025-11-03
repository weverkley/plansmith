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

	// Add epics to "Epics" list
	epicsList := listMap["Epics"]
	for _, epic := range plan.Epics {
		// Create a new card object
		card := &trello.Card{
			Name:   epic.Name,
			Desc:   plan.ProductVision,
			IDList: epicsList.ID,
		}

		// Use the client to create the card
		err := c.client.CreateCard(card, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to create epic card: %w", err)
		}
	}

	// Add user stories to "User Stories (Backlog)" list
	storiesList := listMap["User Stories (Backlog)"]
	for _, story := range plan.UserStories {
		// Create a new card object
		card := &trello.Card{
			Name:   story.Title,
			Desc:   story.Story,
			IDList: storiesList.ID,
		}

		// Use the client to create the card
		err := c.client.CreateCard(card, trello.Defaults())
		if err != nil {
			return fmt.Errorf("failed to create story card: %w", err)
		}
	}

	// Add tasks to "To Do (Ready)" list
	tasksList := listMap["To Do (Ready)"]
	for _, task := range plan.Tasks {
		// Create a new card object
		card := &trello.Card{
			Name:   task.Title,
			Desc:   task.Description,
			IDList: tasksList.ID,
		}

		// Use the client to create the card
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
