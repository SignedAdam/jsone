package main

import "fmt"

const systemPrompt = `You are a structured data extraction tool. You receive text input and output valid JSON.
Output ONLY valid JSON. No markdown, no explanation, no preamble, no code fences.

Rules:
- Output must be valid JSON
- For tabular data: use an array of objects
- For key-value data: use an object
- For lists: use an array
- Use appropriate types: numbers as numbers, booleans as booleans, not strings
- Preserve all data from the input unless the instruction filters it
- If the input has clear column headers, use them as keys
- For grep/find output, parse the standard format (file:line:content)
- Keep keys lowercase, use underscores for spaces`

func buildPrompt(input, instruction string, truncated bool) string {
	var prompt string

	if instruction != "" {
		prompt = fmt.Sprintf("Convert this input to JSON. Instruction: %s\n\nInput:\n%s", instruction, input)
	} else {
		prompt = fmt.Sprintf("Convert this input to the most obvious JSON structure.\n\nInput:\n%s", input)
	}

	if truncated {
		prompt += "\n\n(Note: input was truncated. Process what's available.)"
	}

	return prompt
}
