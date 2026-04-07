package mcp

import (
	"encoding/json"
	"testing"
)

func TestHandleInitialize(t *testing.T) {
	req := request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}

	result := handleRequest(req)
	resp, ok := result.(*response)
	if !ok {
		t.Fatal("expected *response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	resultMap, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatal("expected map result")
	}
	if resultMap["protocolVersion"] != "2024-11-05" {
		t.Errorf("protocol version = %v", resultMap["protocolVersion"])
	}

	serverInfo, ok := resultMap["serverInfo"].(map[string]any)
	if !ok {
		t.Fatal("expected serverInfo map")
	}
	if serverInfo["name"] != "devid" {
		t.Errorf("server name = %v", serverInfo["name"])
	}
}

func TestHandleToolsList(t *testing.T) {
	req := request{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	result := handleRequest(req)
	resp, ok := result.(*response)
	if !ok {
		t.Fatal("expected *response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	resultMap, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatal("expected map result")
	}

	toolsList, ok := resultMap["tools"].([]toolDef)
	if !ok {
		t.Fatal("expected tools list")
	}
	if len(toolsList) != 3 {
		t.Errorf("expected 3 tools, got %d", len(toolsList))
	}

	names := make(map[string]bool)
	for _, tool := range toolsList {
		names[tool.Name] = true
	}
	for _, expected := range []string{"get_identity", "get_snippet", "get_project"} {
		if !names[expected] {
			t.Errorf("missing tool: %s", expected)
		}
	}
}

func TestHandlePing(t *testing.T) {
	req := request{JSONRPC: "2.0", ID: 3, Method: "ping"}
	result := handleRequest(req)
	resp, ok := result.(*response)
	if !ok {
		t.Fatal("expected *response")
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
}

func TestHandleUnknownMethod(t *testing.T) {
	req := request{JSONRPC: "2.0", ID: 4, Method: "bogus/method"}
	result := handleRequest(req)
	resp, ok := result.(*response)
	if !ok {
		t.Fatal("expected *response")
	}
	if resp.Error == nil {
		t.Fatal("expected error for unknown method")
	}
	if resp.Error.Code != -32601 {
		t.Errorf("error code = %d, want -32601", resp.Error.Code)
	}
}

func TestHandleToolCallUnknown(t *testing.T) {
	params, _ := json.Marshal(callToolParams{Name: "nonexistent"})
	req := request{JSONRPC: "2.0", ID: 5, Method: "tools/call", Params: params}
	result := handleRequest(req)
	resp, ok := result.(*response)
	if !ok {
		t.Fatal("expected *response")
	}
	if resp.Error == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestHandleNotification(t *testing.T) {
	req := request{JSONRPC: "2.0", Method: "notifications/initialized"}
	result := handleRequest(req)
	if result != nil {
		t.Error("notifications should return nil")
	}
}
