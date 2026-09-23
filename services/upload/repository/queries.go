package repository

var InsertQuery = "INSERT INTO video_uploads(file_name, user_id, status) VALUES ($1, $2, $3);"
