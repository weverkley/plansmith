# PlanSmith

Like a blacksmith, but for crafting project plans.

PlanSmith is an interactive, chat-like terminal application that "crafts" raw project ideas (from Markdown) into fully-formed, actionable Kanban boards (in Trello), with human review at every step.

## Features

- Transform markdown project descriptions into structured plans
- Generate user stories and technical tasks with AI assistance
- Create Trello boards with pre-populated cards, labels, and checklists
- Review and approve plans before creating boards
- Add new features to existing projects
- List available AI models for the selected provider
- Chat-based interface with file picker
- Intelligent prompt guidance for focused and simplified task generation

## Installation

```bash
go install
```

## Usage

```bash
plansmith
```

## Configuration

Create a configuration file at `~/.plansmith/config.yaml`:

```yaml
ai:
  default_provider: "gemini"
  keys:
    gemini: "YOUR_GEMINI_API_KEY"
    openai: "YOUR_OPENAI_API_KEY"
    qwen: "YOUR_QWEN_API_KEY"
  models:
    gemini: "gemini-pro"
    openai: "gpt-3.5-turbo"
    qwen: "qwen-turbo"
trello:
  key: "YOUR_TRELLO_API_KEY"
  token: "YOUR_TRELLO_API_TOKEN"
```

When valid API keys are provided, the application will use real AI services and create actual Trello boards. If you don't provide API keys, the application will use mock responses for testing purposes.

Supported AI providers:
- Google Gemini
- OpenAI GPT
- Alibaba Qwen (via direct HTTP API)

You can specify which model to use for each provider in the `models` section of the configuration. If no model is specified, the application will use a default model for each provider.

## How It Works

1. Provide a markdown file describing your project
2. PlanSmith uses AI to generate a structured plan with epics, user stories, and tasks, including labels, dependencies, and optional checklists.
3. Review and approve the plan in the TUI
4. PlanSmith creates a Trello board with cards for each item, with hardcoded "Definition of Done" checklists for stories and AI-generated checklists for tasks when applicable. Card dependencies are linked.

For existing projects:
1. Load an existing project plan
2. Provide a markdown file describing a new feature
3. PlanSmith uses AI to generate stories and tasks for the new feature, including labels, dependencies, and optional checklists.
4. Review and approve the additions
5. PlanSmith updates the Trello board with the new items

While in the project dashboard, you can list the available models for your selected AI provider by pressing 'm'.

## Chat-based Interface

PlanSmith now features a chat-based interface for improved usability:

- Type commands directly in the chat interface
- Use `@` to trigger the file picker for browsing and selecting files
- See loading indicators during long-running operations
- Get real-time feedback on your actions

Available chat commands:
- `help` - Show available commands
- `new [file]` - Start a new project with an optional file path
- `open` - Open an existing project
- `list-models` - List available AI models
- `quit` - Exit the application

## Testing Trello Integration

To test the Trello integration with real API calls:

1. Obtain your Trello API key and token from https://trello.com/app-key
2. Update your `$HOME/.plansmith/config.yaml` with the actual credentials
3. Run the test program:

```bash
go run test_trello.go
```

This will create a test board on Trello with sample data to verify the integration is working correctly.

## Current Status

The application structure is complete and follows the project plan. All required packages, dependencies, and TUI interface have been created. The application can be run and both workflows are functional.

Current functionality:
1. TUI with all planned states
2. Processing markdown files into project plans (using real AI when API keys are provided)
3. Saving plans to JSON files
4. Creating Trello boards with correct lane order, labels, hardcoded story checklists, optional AI-generated task checklists, and linked card dependencies.
5. Loading existing plans (including `.json` files via file picker)
6. Adding new features to existing projects
7. Support for multiple AI providers (Gemini, OpenAI, Qwen)
8. Model selection for each AI provider
9. Chat-based interface with file picker
10. Improved prompt guidance for focused and simplified task generation.
11. Rate limiting implemented for AI generative endpoint calls.

## Development Progress

- [x] Phase 1: Foundations & Project Setup
- [x] Phase 2: TUI Application Core
- [x] Phase 3: State & Trello Service
- [x] Phase 4: Full Prompt Definitions
- [x] Phase 5: TUI-Driven Agentic Workflow
- [x] Phase 6: Testing & Deployment

The core functionality is implemented and the application can be run. All identified bugs have been addressed. The remaining work focuses on further testing and potential future enhancements.