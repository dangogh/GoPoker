package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dangogh/GoPoker/cards"
	"github.com/dangogh/GoPoker/hand"
)

type CardInput struct {
	Suit string `json:"suit"`
	Rank string `json:"rank"`
}

type EvaluateHandParams struct {
	Cards []string `json:"cards"`
}

var stringToSuit = map[string]cards.Suit{
	"clubs":    cards.Clubs,
	"diamonds": cards.Diamonds,
	"hearts":   cards.Hearts,
	"spades":   cards.Spades,
	"♣":        cards.Clubs,
	"♦":        cards.Diamonds,
	"♥":        cards.Hearts,
	"♠":        cards.Spades,
}

var stringToRank = map[string]cards.Rank{
	"2":     cards.Two,
	"3":     cards.Three,
	"4":     cards.Four,
	"5":     cards.Five,
	"6":     cards.Six,
	"7":     cards.Seven,
	"8":     cards.Eight,
	"9":     cards.Nine,
	"10":    cards.Ten,
	"J":     cards.Jack,
	"jack":  cards.Jack,
	"Q":     cards.Queen,
	"queen": cards.Queen,
	"K":     cards.King,
	"king":  cards.King,
	"A":     cards.Ace,
	"ace":   cards.Ace,
}

func parseCard(cardStr string) (cards.Card, error) {
	parts := strings.Fields(cardStr)
	if len(parts) != 2 {
		return cards.Card{}, fmt.Errorf("invalid card format: %s (expected 'rank suit')", cardStr)
	}

	rank, ok := stringToRank[parts[0]]
	if !ok {
		return cards.Card{}, fmt.Errorf("invalid rank: %s", parts[0])
	}

	suit, ok := stringToSuit[parts[1]]
	if !ok {
		return cards.Card{}, fmt.Errorf("invalid suit: %s", parts[1])
	}

	return cards.NewCard(suit, rank), nil
}

func main() {
	transport := flag.String("transport", "stdio", "MCP transport protocol: stdio or streamable_http")
	port := flag.String("port", "8080", "Port for streamable_http transport")
	flag.Parse()

	// Create MCP server
	impl := &mcp.Implementation{
		Name:    "gopoker-mcp-server",
		Version: "1.0.0",
	}

	s := mcp.NewServer(impl, &mcp.ServerOptions{
		Instructions: "Poker hand evaluation server for 5-card draw",
	})

	// Define the tool
	tool := &mcp.Tool{
		Name:        "evaluate_poker_hand",
		Description: "Evaluate a 5-card poker hand and get recommended discards for 5-card draw. Returns the hand category (e.g., Pair, Flush, Full House) and suggests which cards to discard to improve the hand. Each card is a string with rank and suit separated by space (e.g., 'A spades', 'K hearts', '10 clubs').",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"cards": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type":        "string",
						"description": "Card as 'rank suit', e.g., 'A spades', 'K hearts', '10 clubs'",
						"pattern":     "^(2|3|4|5|6|7|8|9|10|J|Q|K|A) (clubs|diamonds|hearts|spades|♣|♦|♥|♠)$",
					},
					"minItems":    5,
					"maxItems":    5,
					"description": "Array of 5 cards, each as 'rank suit'",
				},
			},
			"required": []string{"cards"},
		},
	}

	// Handler function
	handler := func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params EvaluateHandParams
		if err := json.Unmarshal(req.Params.Arguments, &params); err != nil {
			return nil, fmt.Errorf("invalid arguments: %w", err)
		}

		if len(params.Cards) != 5 {
			return nil, fmt.Errorf("must provide exactly 5 cards, got %d", len(params.Cards))
		}

		// Parse cards
		cardList := make([]cards.Card, 5)
		for i, cardStr := range params.Cards {
			parsed, err := parseCard(cardStr)
			if err != nil {
				return nil, fmt.Errorf("card %d: %w", i, err)
			}
			cardList[i] = parsed
		}

		// Evaluate hand
		h := hand.Hand{Cards: cardList}
		eval := hand.Evaluate(h)

		// Get recommended discards
		maxDiscard := hand.ComputeMaxDiscard(h)
		discardIdxs := hand.RecommendDiscards(h, maxDiscard)

		// Format discarded cards
		discards := make([]string, len(discardIdxs))
		for i, idx := range discardIdxs {
			discards[i] = cardList[idx].String()
		}

		// Build response text
		resultText := fmt.Sprintf(`Hand Category: %s
Recommended Discards: %v
Max Discards Allowed: %d
Hand Strength: %s with ranks %v`,
			eval.Category.String(),
			discards,
			maxDiscard,
			eval.Category.String(),
			eval.Ranks,
		)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: resultText,
				},
			},
		}, nil
	}

	// Add tool to server
	s.AddTool(tool, handler)

	// Create transport based on flag
	var mcpTransport mcp.Transport
	switch *transport {
	case "stdio":
		mcpTransport = &mcp.StdioTransport{}
		log.Println("Starting GoPoker MCP server (stdio transport)")
		
	case "streamable_http":
		// For streamable_http, we need to create a custom HTTP server
		// that handles MCP protocol over HTTP with chunked transfer
		log.Printf("Starting GoPoker MCP server on port %s (streamable_http transport)", *port)
		
		http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Set headers for chunked streaming
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Transfer-Encoding", "chunked")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")

			// Create an IOTransport for this request
			transport := &mcp.IOTransport{
				Reader: r.Body,
				Writer: &httpChunkWriter{w: w},
			}

			// Connect the server for this request
			ctx := r.Context()
			_, err := s.Connect(ctx, transport, nil)
			if err != nil {
				log.Printf("Connection error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// The connection will be handled by the transport
		})

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		})

		if err := http.ListenAndServe(":"+*port, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
		return
		
	default:
		log.Fatalf("Unknown transport: %s (valid options: stdio, streamable_http)", *transport)
	}

	if err := s.Run(context.Background(), mcpTransport); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// httpChunkWriter wraps http.ResponseWriter to flush after each write
type httpChunkWriter struct {
	w http.ResponseWriter
}

func (h *httpChunkWriter) Write(p []byte) (n int, err error) {
	n, err = h.w.Write(p)
	if f, ok := h.w.(http.Flusher); ok {
		f.Flush()
	}
	return n, err
}

func (h *httpChunkWriter) Close() error {
	return nil
}
