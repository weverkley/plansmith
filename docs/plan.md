Project: PlanSmith (TUI Edition)


PlanSmith: Like a blacksmith, but for crafting project plans.
Mission: To provide an interactive, chat-like terminal application that "crafts" raw project ideas (from Markdown) into fully-formed, actionable Kanban boards (in Trello), with human review at every step.
________________


Phase 1: Foundations & Project Setup (The "Workshop")


Goal: Bootstrap a stateful, event-driven Go TUI application.
1. Core Technology Stack:
   * Language: Go
   * TUI Framework: charmbracelet/bubbletea (This is the core of your application).
   * TUI Components: charmbracelet/bubbles (for spinners, text inputs, lists, etc.)
   * Configuration: viper (for managing API keys)
   * Prompt Templating: Go's native text/template package.
   * Trello Client: A Go client for the Trello API.
   * AI SDKs: Go SDKs for Gemini, OpenAI, etc.
2. Project Structure (TUI-Centric):
The cmd/ folder simply launches the bubbletea program. The pkg/tui/ directory becomes the new "frontend" and brain of the application.
plansmith/
├── .goreleaser.yml
├── .gitignore
├── main.go              # (Logic) go run main.go starts the app
├── cmd/
│   └── root.go          # (Logic) Just `p := tea.NewProgram(...); p.Run()`
├── pkg/
│   ├── ai/              # AI executor abstraction (unchanged)
│   │   ├── executor.go # (Logic) Defines the AIExecutor interface
│   │   ├── factory.go   # (Logic) The NewExecutor() factory
│   │   ├── gemini.go    # (Logic) Gemini implementation
│   │   └── openai.go    # (Logic) OpenAI implementation
│   ├── smith/           # The core "PlanSmith" agent brain
│   │   ├── agent.go     # (Logic) Runs the prompt chain
│   │   └── loader.go    # (Logic) Loads prompts from /prompts
│   ├── trello/          # Trello service (unchanged)
│   │   └── client.go
│   ├── state/           # State management
│   │   └── manager.go   # (Logic) Read/Write .plansmith/ files
│   └── tui/             # <-- CRITICAL: The Bubbletea App
│       ├── model.go     # (Logic) The main bubbletea model
│       ├── states.go    # (Logic) Defines app states (menu, forging, etc)
│       ├── update.go    # (Logic) The main Update() switch statement
│       └── view.go      # (Logic) The main View() switch statement
└── prompts/             # The AI's brain (unchanged from last plan)
   ├── init/
   │   ├── 01_generate_vision.md
   │   ├── 02_generate_stories.md
   │   └── 03_generate_tasks.md
   └── add/
       └── 01_generate_bundle.md

3. Configuration (viper):
   * File: ~/.plansmith/config.yaml
   * Content:
YAML
ai:
 default_provider: "gemini"
keys:
 gemini: "YOUR_GEMINI_API_KEY"
 openai: "YOUR_OPENAI_API_KEY"
trello:
 key: "YOUR_TRELLO_API_KEY"
 token: "YOUR_TRELLO_API_TOKEN"
________________


Phase 2: The "TUI Application" Core (The "Interface")


Goal: Build the bubbletea application state machine. This replaces the draft and forge commands with an interactive UI using a chat approach.
      1. The bubbletea Model (pkg/tui/model.go):
This is the main struct for your application. It will hold all state.
      2. Application States (pkg/tui/states.go):
Your app will be a state machine. The Update and View functions will use a switch on the current state.
         * stateMainMenu: "Welcome to PlanSmith. [Start New Project] [Open Project]"
         * stateNewProjectInput: Asks for the path to the .md file.
         * stateGenerating: A "thinking" state. This will show a spinner and a status message (e.g., "Forging Vision...", "Crafting Stories...").
         * stateReviewingPlan: (This is Improvement #1). A scrollable list/viewport showing the AI-generated plan from the local JSON.
         * stateForgingTrello: An "executing" state. Shows spinners/progress as it calls the TUI API.
         * stateComplete: "Success! Board created."
         * stateProjectDashboard: The "home screen" for an existing project. Shows "Add Feature", "View Board", etc.
         * stateAddingFeature: Asks for the .md file for the new feature.
         * stateError: Displays any errors.
________________


Phase 3: State & Trello Service (The "Anvil & Hammer")


Goal: Establish the on-disk "memory" and Trello connection logic.
         1. State Management (pkg/state/manager.go):
         * plansmith.plan.json (The "Source of Truth"): This is the most important file. It's the full, nested JSON of the entire plan (Vision, Epics, Stories, Tasks). The TUI loads this to show the stateProjectDashboard and stateReviewingPlan. The add command uses it for context.
         * .plansmith/state.json (The "Config"): This is a tiny file, separate from the plan. It just stores the TrelloBoardID and TrelloBoardURL so the app knows which board this local plan is linked to.
         2. Trello Service (pkg/trello/client.go):
         * CreateProjectBoard: Creates the board and the 5 lists: Epics, User Stories (Backlog), To Do (Ready), Doing, Done.
         * PopulateBoard: Reads the plansmith.plan.json and creates all the cards.
         * LinkCards: (This is Improvement #4). After all cards are created, this function loops through the plan, finds the dependencies for each task, and uses the Trello API to create attachments/links between the cards.
________________


Phase 4: Full Prompt Definitions (The "Crafting")


Goal: Define the "chatty" 1+N+M prompt chain, as you requested, but with the "dependency" enhancement.


1. init Workflow Prompts


File: prompts/init/01_generate_vision.md


Markdown




You are a world-class Product Owner. Your purpose is to analyze a raw project document and forge it into a clear, actionable vision.

You must strictly adhere to the following rules:
1.  Analyze the user's markdown document provided below.
2.  Generate a concise and inspiring **Project Name**.
3.  Generate a **Product Vision** statement (1-3 sentences) that captures the "why" of the project.
4.  Generate a list of 3-5 high-level **Epics** (the main features).
5.  You MUST respond ONLY with a single, valid JSON object. Do not include any text, apologies, or explanations before or after the JSON.

**JSON Schema:**
{
 "project_name": "string",
 "product_vision": "string",
 "epics": ["string", "string", ...]
}

---
**Project Markdown:**
{{.MarkdownInput}}

File: prompts/init/02_generate_stories.md


Markdown




You are a world-class Product Owner. Your purpose is to break down a high-level Epic into concrete, value-driven User Stories.

You must strictly adhere to the following rules:
1.  Your overall project vision is: {{.ProductVision}}
2.  The Epic you are currently breaking down is: {{.EpicName}}
3.  Generate 3-5 User Stories for this Epic.
4.  Each User Story MUST follow the format: "As a [USER], I want [FEATURE], so that [VALUE]."
5.  Assign a priority (1-10, 10=highest) to each story based on its importance to the Epic.
6.  You MUST respond ONLY with a single, valid JSON object. Do not include any text, apologies, or explanations before or after the JSON.

**JSON Schema:**
{
 "user_stories": [
   {
     "title": "string (A short title for the story)",
     "story": "string (As a user, I want..., so that...)",
     "priority": "integer",
     "epic": "{{.EpicName}}"
   }
 ]
}

File: prompts/init/03_generate_tasks.md


Markdown




You are an expert Lead Developer. Your purpose is to take a User Story and break it down into small, concrete technical tasks for your development team.

You must strictly adhere to the following rules:
1.  The User Story you are breaking down is:
    - **Title:** {{.StoryTitle}}
    - **Story:** {{.StoryBody}}
2.  Generate 2-5 granular, technical sub-tasks required to complete this story.
3.  **CRITICAL:** For each task, identify any other tasks *in this list* that MUST be completed first. Use the exact task `title` for the dependency. If there are no dependencies, return an empty array.
4.  You MUST respond ONLY with a single, valid JSON object. Do not include any text, apologies, or explanations before or after the JSON.

**JSON Schema:**
{
 "tasks": [
    {
     "title": "string (The technical task)",
     "description": "string (A brief 1-sentence description of the task)",
     "story_title": "{{.StoryTitle}}",
     "dependencies": ["string (task title)", ...]
   }
 ]
}

________________


2. add Workflow Prompt


File: prompts/add/01_generate_bundle.md


Markdown




You are an expert Agile team (Product Owner + Lead Developer). A new feature request has come in for an existing project, and you must plan it.

You must strictly adhere to the following rules:
1.  The existing **Project Plan** (which includes all epics, stories, and tasks) is:
    ---
   {{.CurrentPlanJSON}}
   ---
2.  The new **Feature Request** (in markdown) is:
   ---
   {{.MarkdownInput}}
   ---
3.  First, act as the **Product Owner**: Generate 1-3 User Stories for this new feature. Do NOT duplicate stories from the existing plan.
4.  Second, act as the **Lead Developer**: For each new User Story, generate its technical sub-tasks, including dependencies.
5.  You MUST respond ONLY with a single, valid JSON object containing *only the new items*.

**JSON Schema:**
{
 "feature_bundle": [
   {
     "title": "string (A short title for the story)",
     "story": "string (As a user, I want..., so that...)",
     "priority": "integer",
     "tasks": [
       {
         "title": "string (The technical task)",
         "description": "string (A brief 1-sentence description of the task)",
         "dependencies": ["string (task title)", ...]
       }
     ]
   }
 ]
}

________________


Phase 5: The TUI-Driven Agentic Workflow


Goal: Orchestrate all components, driven by user interaction in the TUI in a chat style.
         1. "New Project" Workflow:
         1. User launches plansmith -> TUI enters stateMainMenu.
         2. User selects "Start New Project" -> TUI enters stateNewProjectInput.
         3. User inputs path/to/project.md.
         4. TUI enters stateGenerating and sends a tea.Cmd (a bubbletea command).
         5. The tea.Cmd calls agent.GenerateVision(). The TUI Update function gets a (visionMsg, err) back.
         6. TUI status updates: "Forging Stories for Epic 1...". It loops and calls agent.GenerateStories() for each epic.
         7. TUI status updates: "Forging Tasks...". It loops and calls agent.GenerateTasks() for each story.
         8. All AI-generated data is saved locally to plansmith.plan.json.
         9. TUI enters stateReviewingPlan, loading the plansmith.plan.json into a scrollable viewport.
         10. The user reviews the plan and hits "Approve" (or "Edit", a v2 feature).
         11. TUI enters stateForgingTrello. It sends commands to trello.CreateProjectBoard(), trello.PopulateBoard(), and trello.LinkCards().
         12. Trello board ID is saved to .plansmith/state.json.
         13. TUI enters stateComplete.
         2. "Add Feature" Workflow:
         1. User launches plansmith in a directory containing plansmith.plan.json.
         2. TUI detects the file, loads it, and enters stateProjectDashboard.
         3. User selects "Add New Feature" -> TUI enters stateAddingFeature.
         4. User inputs path/to/feature.md.
         5. TUI enters stateGenerating. It calls agent.GenerateAddBundle(), passing both the new markdown and the entire plansmith.plan.json as context (Improvement #3).
         6. The new feature data is appended to plansmith.plan.json.
         7. TUI enters stateReviewingPlan to review only the new items.
         8. User "Approves".
         9. TUI enters stateForgingTrello. It calls new Trello functions like trello.AddCardsToBoard() and trello.LinkCards(), updating the existing board.
________________


Phase 6: Testing & Deployment (QA & Shipping)


Goal: Ensure reliability and distribute the application.
         1. Testing Strategy:
         * Unit Tests: Will test the smith.loader to ensure templates are populated correctly. Will also test the state.manager.
         * Agent Logic Tests (Critical): Test the agent functions without the TUI. Mock the AIExecutor interface to return a pre-defined JSON string. This will test the entire agent workflow, JSON parsing, and state transitions work.
         * TUI Testing: Testing bubbletea apps is complex. You can use bubbletea/teatest to send virtual keypresses and snapshot the TUI's string output, but focus most of your effort on testing the agent logic separately.
         2. Deployment (.goreleaser.yml):
         * goreleaser is perfect for this. It will compile your plansmith application into a single, static binary for Windows, macOS, and Linux, and upload them to a GitHub Release.