# This is a website with a parser that steals information and adds it to the main page :)

---

### For start you need to install the following packages:
- gorm
- framework Echo
- redis
- driver for postgres
- goquery
- gmail-api
- telebot

and other packages that are in the go.mod file

### To start the server, you need to run the following command:
```docker compose up```

because the server is running in a docker container

Redis needs to log in, because the server uses it for caching passwords

Postgres needs to save information from web scraping

Also, you need set up your account Google for sending emails,
and you need to set up your bot for sending messages to Telegram when someone registers on the site

### The server will start at the following address:
```http://localhost:1323```

### All html pages are located in the following directory:
```/frontend/html```

#### Very easy to use and very easy to understand :)