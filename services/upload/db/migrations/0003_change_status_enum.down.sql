ALTER TABLE video_uploads
ALTER COLUMN status TYPE text;

DROP TYPE video_status;

CREATE TYPE video_status AS ENUM (
    'uploading',
    'processing',
    'completed'
);

ALTER TABLE video_uploads
ALTER COLUMN status TYPE video_status
USING status::video_status;
