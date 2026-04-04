package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/EyuAtske/Agrregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error{
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	_, err :=s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("There is no user with the provided name")
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	fmt.Println("The user name has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error{
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the register handler expects a single argument, the username")
	}
	_, err :=s.db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		fmt.Print(err)
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	arg := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
	}
	_, err = s.db.CreateUser(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}
	fmt.Println("User was created")
	return nil
}

func handlerReset(s *state, cmd command) error{
	err := s.db.DeleteUsers(context.Background())
	if err != nil{
		fmt.Println("Unable to reset users: %w", err)
		os.Exit(1)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error{
	users, err := s.db.GetUsers(context.Background())
	if err != nil{
		fmt.Println("Unable to fetch users: %w", err)
		os.Exit(1)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name)
			continue
		}
		fmt.Println(user.Name)
	}
	
	return nil
}

func handlerFetch(s *state, cmd command) error{
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the agg handler expects a single argument, time_between_reqs 1s, 1m or 1h")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to parse time duration: %w", err)
	}
	fmt.Printf("Collecting feeds every %s\n", cmd.Args[0])
	ticker := time.NewTicker(time.Duration(timeBetweenReqs))
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error{
	if len(cmd.Args) < 2 {
		return fmt.Errorf("the addfeed handler expects two arguments, the feed name and the feed url")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed: %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to follow feed: %w", err)
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error{
	feeds, err:= s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to fetch feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("- %s\n", feed.Name)
		fmt.Printf("  Url: %s\n", feed.Url)
		fmt.Printf("  Added by: %s\n", feed.UserName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error{
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the follow handler expects a single argument, the feed url")
	}
	feed_url := cmd.Args[0]
	feed , err := s.db.GetFeedByUrl(context.Background(), feed_url)
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %w", err)
	}
	feeds, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to follow feed: %w", err)
	}
	fmt.Printf("feed name: %s\n", feeds[0].FeedName)
	fmt.Printf("feed user name: %s\n", feeds[0].UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error{
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to fetch feed follows: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("- %s\n", feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error{
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the unfollow handler expects a single argument, the feed url")
	}
	feed_url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feed_url)
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %w", err)
	}
	err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to unfollow feed: %w", err)
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if cmd.Args != nil {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}