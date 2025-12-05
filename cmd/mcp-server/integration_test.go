package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMCPServerIntegration tests the full MCP server flow with an in-memory connection
func TestMCPServerIntegration(t *testing.T) {
	// Create in-memory transports
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	// Start server in a goroutine
	serverImpl := &mcp.Implementation{
		Name:    "gopoker-test-server",
		Version: "1.0.0",
	}

	server := mcp.NewServer(serverImpl, &mcp.ServerOptions{
		Instructions: "Test poker server",
	})

	// Define and add the tool
	tool := &mcp.Tool{
		Name:        "evaluate_poker_hand",
		Description: "Evaluate poker hands",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"cards": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
					"minItems": 5,
					"maxItems": 5,
				},
			},
			"required": []string{"cards"},
		},
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params EvaluateHandParams
		if err := json.Unmarshal(req.Params.Arguments, &params); err != nil {
			return nil, err
		}

		if len(params.Cards) != 5 {
			return nil, assert.AnError
		}

		// Parse and evaluate cards
		cardList := make([]interface{}, 5)
		for i, cardStr := range params.Cards {
			parsed, err := parseCard(cardStr)
			if err != nil {
				return nil, err
			}
			cardList[i] = parsed
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Hand evaluated successfully",
				},
			},
		}, nil
	}

	server.AddTool(tool, handler)

	// Connect server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverDone := make(chan error, 1)
	go func() {
		serverDone <- server.Run(ctx, serverTransport)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Create and connect client
	clientImpl := &mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}
	client := mcp.NewClient(clientImpl, nil)

	clientSession, err := client.Connect(ctx, clientTransport, nil)
	require.NoError(t, err)
	defer clientSession.Close()

	// Test 1: List tools
	t.Run("list_tools", func(t *testing.T) {
		result, err := clientSession.ListTools(ctx, nil)
		require.NoError(t, err)
		require.Len(t, result.Tools, 1)
		assert.Equal(t, "evaluate_poker_hand", result.Tools[0].Name)
	})

	// Test 2: Call tool with valid hand
	t.Run("evaluate_valid_hand", func(t *testing.T) {
		params := map[string]interface{}{
			"cards": []string{
				"A spades",
				"K spades",
				"Q spades",
				"J spades",
				"10 spades",
			},
		}
		result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name:      "evaluate_poker_hand",
			Arguments: params, // Pass the map directly, not the JSON bytes
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Content, 1)

		textContent, ok := result.Content[0].(*mcp.TextContent)
		require.True(t, ok)
		assert.Contains(t, textContent.Text, "evaluated successfully")
	})

	// Test 3: Call tool with invalid card count
	t.Run("evaluate_invalid_card_count", func(t *testing.T) {
		params := map[string]interface{}{
			"cards": []string{
				"A spades",
				"K spades",
			},
		}

		_, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name:      "evaluate_poker_hand",
			Arguments: params, // Pass the map directly
		})
		assert.Error(t, err)
	})

	// Test 4: Call tool with invalid card format
	t.Run("evaluate_invalid_card_format", func(t *testing.T) {
		params := map[string]interface{}{
			"cards": []string{
				"Aspades",
				"K spades",
				"Q spades",
				"J spades",
				"10 spades",
			},
		}

		_, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
			Name:      "evaluate_poker_hand",
			Arguments: params, // Pass the map directly
		})
		assert.Error(t, err)
	})

	// Close client and wait for server
	clientSession.Close()
	cancel()

	select {
	case err := <-serverDone:
		// Server should exit cleanly when context is cancelled
		assert.Error(t, err) // context.Canceled
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down")
	}
}
