package db

// Saving all queries needed by each table

// User
const INSERT_USER_STATEMENT = "Insert INTO User (email, password, created_at) VALUES (?, ?, ?)"
const GET_USER_LOGIN_STATEMENT = "SELECT user_id, email, password FROM User WHERE email = ?"
const GET_LOGGED_USER_STATEMENT = "SELECT user_id, email, created_at, user_handle FROM User WHERE user_id = ?"
