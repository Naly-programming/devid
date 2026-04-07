package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/generate"
)

// MCP JSON-RPC types
type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   *rpcError  `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type notification struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// Tool definitions
type toolDef struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	InputSchema inputSchema `json:"inputSchema"`
}

type inputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]property `json:"properties,omitempty"`
}

type property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type textContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type callToolParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

var tools = []toolDef{
	{
		Name:        "get_identity",
		Description: "Get the developer's identity profile including tone, stack, conventions, and AI preferences",
		InputSchema: inputSchema{Type: "object"},
	},
	{
		Name:        "get_snippet",
		Description: "Get a compact identity snippet suitable for pasting into a conversation",
		InputSchema: inputSchema{Type: "object"},
	},
	{
		Name:        "get_project",
		Description: "Get project-specific context for a given project name",
		InputSchema: inputSchema{
			Type: "object",
			Properties: map[string]property{
				"project": {Type: "string", Description: "Project name to look up"},
			},
		},
	},
}

// Serve runs the MCP server on stdin/stdout.
func Serve() error {
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var req request
		if err := json.Unmarshal(line, &req); err != nil {
			continue
		}

		resp := handleRequest(req)
		if resp != nil {
			out, _ := json.Marshal(resp)
			fmt.Fprintf(writer, "%s\n", out)
		}
	}
}

func handleRequest(req request) any {
	switch req.Method {
	case "initialize":
		return &response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]any{
					"tools": map[string]any{},
				},
				"serverInfo": map[string]any{
					"name":    "devid",
					"version": "0.1.0",
				},
			},
		}

	case "notifications/initialized":
		return nil // No response for notifications

	case "tools/list":
		return &response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"tools": tools,
			},
		}

	case "tools/call":
		var params callToolParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return &response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &rpcError{Code: -32602, Message: "invalid params"},
			}
		}
		return handleToolCall(req.ID, params)

	case "ping":
		return &response{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{}}

	default:
		return &response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &rpcError{Code: -32601, Message: "method not found: " + req.Method},
		}
	}
}

func handleToolCall(id any, params callToolParams) *response {
	var text string
	var err error

	switch params.Name {
	case "get_identity":
		text, err = getIdentityText()
	case "get_snippet":
		text, err = getSnippetText()
	case "get_project":
		projName, _ := params.Arguments["project"].(string)
		text, err = getProjectText(projName)
	default:
		return &response{
			JSONRPC: "2.0",
			ID:      id,
			Error:   &rpcError{Code: -32602, Message: "unknown tool: " + params.Name},
		}
	}

	if err != nil {
		return &response{
			JSONRPC: "2.0",
			ID:      id,
			Result: map[string]any{
				"content": []textContent{{Type: "text", Text: "Error: " + err.Error()}},
				"isError": true,
			},
		}
	}

	return &response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]any{
			"content": []textContent{{Type: "text", Text: text}},
		},
	}
}

func getIdentityText() (string, error) {
	id, err := config.Load()
	if err != nil {
		return "", err
	}
	return generate.Render(id, generate.TargetClaudeGlobal, nil)
}

func getSnippetText() (string, error) {
	id, err := config.Load()
	if err != nil {
		return "", err
	}
	return generate.Render(id, generate.TargetSnippet, nil)
}

func getProjectText(name string) (string, error) {
	id, err := config.Load()
	if err != nil {
		return "", err
	}

	if name == "" {
		var names []string
		for _, p := range id.Projects {
			names = append(names, p.Name)
		}
		if len(names) == 0 {
			return "No projects configured. Run `devid add` to add one.", nil
		}
		return fmt.Sprintf("Available projects: %v", names), nil
	}

	for i := range id.Projects {
		if id.Projects[i].Name == name || id.Projects[i].Repo == name {
			return generate.Render(id, generate.TargetClaudeProject, &id.Projects[i])
		}
	}

	return fmt.Sprintf("Project %q not found.", name), nil
}
