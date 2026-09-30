package repository

var updateVideoStatusToCompleted = "UPDATE video_uploads vu SET status='completed' WHERE vu.id=$1;"
var updateVideoStatusToFailed = "UPDATE video_uploads vu SET status='failed' WHERE vu.id=$1;"

var updateRequeuesWithFailure = `
	UPDATE video_uploads
	SET
		requeues = requeues + 1,
		error = $2,
		status = CASE
			WHEN requeues >= 3 THEN 'failed'::video_status
			ELSE status
		END,
		failed_at = CASE
			WHEN requeues >= 3 THEN now()
			ELSE failed_at
		END,
		updated_at = now()
	WHERE id = $1
	RETURNING requeues;
`

var getRetryCount = "SELECT retry_count FROM videos WHERE id = $1;"
