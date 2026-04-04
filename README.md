# 🐊 Gator CLI

Gator is a command-line RSS feed aggregator built with Go and PostgreSQL. It allows users to register, follow feeds, and browse aggregated content directly from the terminal.

---

## ⚙️ Prerequisites

Before running Gator, make sure you have the following installed:

* Go (version 1.20+ recommended)
* PostgreSQL (running locally or remotely)

---

## 📦 Installation

Install the Gator CLI using `go install`:

```bash
go install github.com/EyuAtske/Agrregator@latest
```

Make sure your `$GOPATH/bin` (or `$HOME/go/bin`) is in your system `PATH`.

---

## 🗄️ Database Setup

1. Start PostgreSQL.
2. Create a database named `gator`:

```sql
CREATE DATABASE gator;
```

3. Update the connection string in `main.go` if needed:

```go
"postgres://postgres:datapost@localhost:5432/gator?sslmode=disable"
```

> ⚠️ Replace `postgres:datapost` with your actual database username and password.

---

## 🧾 Configuration Setup

Gator uses a config file to store user-specific settings.

### Create the config file:

```
~/.gatorconfig.json
```

### Example config:

```json
{
  "current_user_name": ""
}
```

This file will be updated automatically when you log in or register.

---

## ▶️ Running the Program

Run the CLI like this:

```bash
Agrregator <command> [arguments...]
```

---

## 🧠 Available Commands

### 👤 User Management

* `register <username>`
  Create a new user.

* `login <username>`
  Log in as an existing user.

* `users`
  List all registered users.

---

### 🛠️ Admin Commands

* `reset`
  Reset the database (⚠️ deletes all data).

---

### 📡 Feed Management

* `addfeed <name> <url>`
  Add a new RSS feed (requires login).

* `feeds`
  List all available feeds.

* `follow <feed_url>`
  Follow a feed.

* `unfollow <feed_url>`
  Unfollow a feed.

* `following`
  List feeds you are following.

---

### 📰 Aggregation & Browsing

* `agg`
  Fetch and store latest posts from feeds.

* `browse`
  View posts from feeds you follow.

---

## ❗ Common Errors

* **not enough arguments were provided**
  → No command was passed.

* **a second argument is required**
  → Some commands need extra arguments.

* **Database connection issues**
  → Ensure PostgreSQL is running and your connection string is correct.

---

## 🧩 Notes

* Some commands require you to be logged in (`addfeed`, `follow`, `browse`, etc.).
* All data is stored in PostgreSQL.
* The config file keeps track of the current user session.
