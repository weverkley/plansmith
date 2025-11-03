package smith

import (
	"fmt"
	"os"
	"text/template"
)

type Loader struct {
	promptsPath string
}

func NewLoader() *Loader {
	return &Loader{
		promptsPath: "prompts",
	}
}

func (l *Loader) LoadInitPrompt(step int, data map[string]interface{}) (string, error) {
	filename := fmt.Sprintf("%02d_generate_vision.md", step)
	if step == 2 {
		filename = fmt.Sprintf("%02d_generate_stories.md", step)
	} else if step == 3 {
		filename = fmt.Sprintf("%02d_generate_tasks.md", step)
	}

	filePath := fmt.Sprintf("%s/init/%s", l.promptsPath, filename)
	return l.loadAndExecuteTemplate(filePath, data)
}

func (l *Loader) LoadAddPrompt(step int, data map[string]interface{}) (string, error) {
	filename := fmt.Sprintf("%02d_generate_bundle.md", step)
	filePath := fmt.Sprintf("%s/add/%s", l.promptsPath, filename)
	return l.loadAndExecuteTemplate(filePath, data)
}

func (l *Loader) loadAndExecuteTemplate(templatePath string, data map[string]interface{}) (string, error) {
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", templatePath, err)
	}

	tmpl, err := template.New("prompt").Parse(string(templateBytes))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result string
	writer := &stringWriter{&result}
	err = tmpl.Execute(writer, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result, nil
}

// stringWriter is a helper to write template output to a string
type stringWriter struct {
	str *string
}

func (sw *stringWriter) Write(p []byte) (n int, err error) {
	*sw.str += string(p)
	return len(p), nil
}
