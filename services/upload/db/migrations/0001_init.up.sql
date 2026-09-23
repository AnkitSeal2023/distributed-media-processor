CREATE TYPE video_status AS ENUM('uploading', 'processing', 'completed');

CREATE TABLE video_uploads(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	file_name VARCHAR(255) NOT NULL,
	user_id UUID NOT NULL,
	status video_status NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL
);
