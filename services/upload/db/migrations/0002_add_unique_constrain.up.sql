ALTER TABLE video_uploads
ADD CONSTRAINT unique_user_file_name UNIQUE (user_id, file_name);
