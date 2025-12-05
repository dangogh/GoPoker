package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

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

	// Define and add the tool
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

	handler := func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params EvaluateHandParams
		if err := json.Unmarshal(req.Params.Arguments, &params); err != nil {
			return nil, err
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

	s.AddTool(tool, handler)

	// Create transport based on flag
	switch *transport {
	case "stdio":
		mcpTransport := &mcp.StdioTransport{}
		log.Println("Starting GoPoker MCP server (stdio transport)")
		if err := s.Run(context.Background(), mcpTransport); err != nil {
			log.Fatalf("Server failed: %v", err)
		}

	case "streamable_http":
		log.Printf("Starting GoPoker MCP server on port %s (streamable_http transport)", *port)

		// Track active sessions
		sessions := &sessionManager{
			sessions: make(map[string]*sessionInfo),
		}

		http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Get or create session ID
			sessionID := r.Header.Get("Mcp-Session-Id")
			if sessionID == "" {
				// New session
				sessionID = sessions.newSession()
				log.Printf("New session: %s", sessionID)
			}

			// Set response headers
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Mcp-Session-Id", sessionID)
			w.Header().Set("Cache-Control", "no-cache")

			// Create a bidirectional pipe for this request
			pr, pw := io.Pipe()

			// Wrap the response writer to capture writes
			respWriter := &responseCapture{
				ResponseWriter: w,
				pipe:           pw,
			}

			// Create IOTransport for this request/response
			transport := &mcp.IOTransport{
				Reader: r.Body,
				Writer: respWriter,
			}

			// Handle the connection
			ctx := r.Context()
			_, err := s.Connect(ctx, transport, nil)
			if err != nil {
				log.Printf("Connection error for session %s: %v", sessionID, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Read from pipe and write to response
			io.Copy(w, pr)
		})

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		})

		// Bind to all interfaces for remote access
		addr := "0.0.0.0:" + *port
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}

	default:
		log.Fatalf("Unknown transport: %s (valid options: stdio, streamable_http)", *transport)
	}
}

// responseCapture captures writes and pipes them through for streaming
type responseCapture struct {
	http.ResponseWriter
	pipe *io.PipeWriter
	once sync.Once
}

func (r *responseCapture) Write(p []byte) (n int, err error) {
	// Write to both the HTTP response and the pipe
	n, err = r.ResponseWriter.Write(p)
	if err != nil {
		return n, err
	}

	// Flush if possible
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}

	// Also write to pipe for coordination
	r.pipe.Write(p)

	return n, err
}

func (r *responseCapture) Close() error {
	r.once.Do(func() {
		r.pipe.Close()
	})
	return nil
}

// sessionManager tracks active MCP sessions
type sessionManager struct {
	mu       sync.Mutex
	sessions map[string]*sessionInfo
	counter  int
}

type sessionInfo struct {
	id      string
	created int64
}

func (sm *sessionManager) newSession() string {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.counter++
	id := fmt.Sprintf("session-%d", sm.counter)
	sm.sessions[id] = &sessionInfo{
		id:      id,
		created: 0, // Could use time.Now().Unix()
	}

	return id
}
