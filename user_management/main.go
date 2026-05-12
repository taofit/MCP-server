package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type App struct {
	db *sql.DB
}

type User struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

type ListUsersArgs struct{}

type GetUserArgs struct {
	ID int `json:"id"`
}

type AddUserArgs struct {
	Name string `json:"name"`
}

type DeleteUserArgs struct {
	ID int `json:"id"`
}

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *sql.DB
	var err error

	log.Println("Waiting for database to be ready...")
	for i := range 10 {
		db, err = sql.Open("postgres", psqlInfo)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Database not ready yet (attempt %d/10): %v, retrying in %d seconds...", i+1, err, 2*(i+1))
		time.Sleep(time.Duration(2*(i+1)) * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to the database after several attempts: %v", err)
	}
	defer db.Close()
	app := &App{db: db}
	app.seedDatabase()
	log.Println("Successfully connected to the database!")

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "user-management",
		Version: "0.0.1",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_users",
		Description: "list users in the database",
	}, app.listUsers)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_user",
		Description: "get a user by ID",
	}, app.getUser)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_user",
		Description: "add a new user",
	}, app.addUser)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_user",
		Description: "delete a user by ID",
	}, app.deleteUser)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

func (app *App) listUsers(ctx context.Context, req *mcp.CallToolRequest, args *ListUsersArgs) (*mcp.CallToolResult, any, error) {
	query := "SELECT id, name, created_at FROM users ORDER BY created_at DESC"
	rows, err := app.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error when querying users: %v", err)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error when querying users: %v", err),
				},
			},
		}, nil, nil
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error when scanning user row: %v", err),
					},
				},
			}, nil, nil
		}
		users = append(users, user)
	}
	userList := ""
	for _, user := range users {
		userList += fmt.Sprintf("ID: %d, Name: %s, Created At: %s\n", user.ID, user.Name, user.CreatedAt)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Found %d users, they are:\n%s", len(users), userList),
			},
		},
	}, nil, nil
}

func (app *App) getUser(ctx context.Context, req *mcp.CallToolRequest, args *GetUserArgs) (*mcp.CallToolResult, any, error) {
	var user User
	err := app.db.QueryRowContext(ctx, "SELECT id, name, created_at FROM users WHERE id = $1", args.ID).Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("User not found: %v", err)}}}, nil, nil
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("ID: %d, Name: %s, Created: %s", user.ID, user.Name, user.CreatedAt)}}}, nil, nil
}

func (app *App) addUser(ctx context.Context, req *mcp.CallToolRequest, args *AddUserArgs) (*mcp.CallToolResult, any, error) {
	_, err := app.db.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", args.Name)
	if err != nil {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error adding user: %v", err)}}}, nil, nil
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Successfully added user: %s", args.Name)}}}, nil, nil
}

func (app *App) deleteUser(ctx context.Context, req *mcp.CallToolRequest, args *DeleteUserArgs) (*mcp.CallToolResult, any, error) {
	_, err := app.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", args.ID)
	if err != nil {
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error deleting user: %v", err)}}}, nil, nil
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Successfully deleted user with ID: %d", args.ID)}}}, nil, nil
}

func (app *App) seedDatabase() {
	users := []User{
		{Name: "John Doe"},
		{Name: "Jane Smith"},
		{Name: "Bob Johnson"},
		{Name: "Alice Williams"},
		{Name: "David Lee"},
		{Name: "Eve Martin"},
		{Name: "Frank Wilson"},
		{Name: "Grace Taylor"},
		{Name: "Henry Anderson"},
		{Name: "Ivy Thomas"},
		{Name: "Jack White"},
		{Name: "Kelly Green"},
		{Name: "Leo Black"},
		{Name: "Mia Blue"},
		{Name: "Noah Red"},
	}
	ctx := context.Background()
	for _, user := range users {
		query := "INSERT INTO users (name) VALUES ($1) ON CONFLICT (name) DO NOTHING"
		_, err := app.db.ExecContext(ctx, query, user.Name)
		if err != nil {
			log.Fatalf("Error inserting into database: %v", err)
		}
	}
	log.Println("Database seeding completed!")
}
