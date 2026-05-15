package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "xiaohongshu-mcp"
	serverVersion = "0.1.0"

	// defaultSearchLimit is the default number of results returned by search.
	// Increased from 10 to 20 for more comprehensive results by default.
	defaultSearchLimit = 20

	// maxSearchLimit caps the number of results to avoid overly large responses.
	maxSearchLimit = 50
)

func main() {
	// Initialize MCP server for Xiaohongshu (Little Red Book)
	s := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithToolCapabilities(true),
	)

	// Register search tool
	searchTool := mcp.NewTool(
		"xiaohongshu_search",
		mcp.WithDescription("Search for notes/posts on Xiaohongshu (Little Red Book)"),
		mcp.WithString(
			"keyword",
			mcp.Required(),
			mcp.Description("The keyword to search for on Xiaohongshu"),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description(fmt.Sprintf("Maximum number of results to return (default: %d, max: %d)", defaultSearchLimit, maxSearchLimit)),
		),
	)
	s.AddTool(searchTool, searchHandler)

	// Register get note detail tool
	getNoteTool := mcp.NewTool(
		"xiaohongshu_get_note",
		mcp.WithDescription("Get detailed information about a specific Xiaohongshu note"),
		mcp.WithString(
			"note_id",
			mcp.Required(),
			mcp.Description("The ID of the Xiaohongshu note"),
		),
	)
	s.AddTool(getNoteTool, getNoteHandler)

	// Start stdio server
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}

// searchHandler handles the xiaohongshu_search tool call
func searchHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keyword, ok := req.Params.Arguments["keyword"].(string)
	if !ok || keyword == "" {
		return mcp.NewToolResultError("keyword is required and must be a string"), nil
	}

	limit := defaultSearchLimit
	if l, ok := req.Params.Arguments["limit"].(float64); ok && l > 0 {
		limit = int(l)
		if limit > maxSearchLimit {
			limit = maxSearchLimit
		}
	}

	results, err := searchXiaohongshu(ctx, keyword, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// getNoteHandler handles the xiaohongshu_get_note tool call
func getNoteHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	noteID, ok := req.Params.Arguments["note_id"].(string)
	if !ok || noteID == "" {
		return mcp.NewToolResultError("note_id is required and must be a string"), nil
	}

	note, err := getNoteDetail(ctx, noteID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get note: %v", err)), nil
	}

	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal note: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}
