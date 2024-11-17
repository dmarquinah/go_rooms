package db

// Saving all queries needed by each table

// User
const INSERT_USER_STATEMENT = "Insert INTO User (email, password, created_at) VALUES (?, ?, ?)"
const GET_USER_LOGIN_STATEMENT = "SELECT user_id, email, password FROM User WHERE email = ?"
const GET_LOGGED_USER_STATEMENT = "SELECT user_id, email, created_at, user_handle FROM User WHERE user_id = ?"

// Host

const GET_HOST_LOGIN_STATEMENT = "SELECT host_id, host_username, password, is_verified, created_at FROM Host WHERE host_username = ?"
const INSERT_HOST_STATEMENT = "Insert INTO Host (host_username, password) VALUES (?, ?)"
const GET_LOGGED_HOST_STATEMENT = "SELECT host_id, host_username, is_verified, created_at, description FROM Host WHERE host_id = ?"
