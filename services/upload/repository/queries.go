package repository

var InsertVideo = "INSERT INTO video_uploads(file_name, user_id, status) VALUES ($1, $2, $3) RETURNING id;"

var UpdateVideoStatus = "UPDATE video_uploads SET status = $1 WHERE id = $2;"
