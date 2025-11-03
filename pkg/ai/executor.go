package ai

type AIExecutor interface {
	ExecutePrompt(prompt string) (string, error)
}
