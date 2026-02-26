-- name: CheckRoomMembershipByRoomId :one
SELECT EXISTS (
    SELECT 1
    FROM rooms r
    LEFT JOIN room_memberships rm
        ON rm.room_id = r.id
        AND rm.user_id = $2
    WHERE r.id = $1
        AND (r.owner_id = $2 OR rm.user_id IS NOT NULL)
);

-- name: CheckRoomOwnershipByRoomId :one
SELECT EXISTS (
    SELECT 1
    FROM rooms
    WHERE id = $1
        AND owner_id = $2
);

-- name: CheckRoomMembershipBySessionId :one
SELECT EXISTS (
    SELECT 1
    FROM rooms r
    LEFT JOIN room_memberships rm
        ON r.id = rm.room_id
        AND rm.user_id = $2
    JOIN sessions s
        ON r.id = s.room_id
    WHERE s.id = $1
        AND (r.owner_id = $2 OR rm.user_id IS NOT NULL)
);

-- name: CheckRoomOwnershipBySessionId :one
SELECT EXISTS (
    SELECT 1
    FROM rooms r
    JOIN sessions s
        ON r.id = s.room_id
    WHERE s.id = $1
        AND r.owner_id = $2
);

-- name: CheckSessionMembershipOnly :one
SELECT EXISTS (
    SELECT 1
    FROM sessions s
    LEFT JOIN session_memberships sm
        ON s.id = sm.session_id
        AND sm.user_id = $2
    WHERE s.id = $1
        AND sm.user_id IS NOT NULL
);

-- name: CheckSessionMembership :one
SELECT EXISTS (
SELECT 1
FROM sessions s
WHERE s.id = $1
  AND (
      s.owner_id = $2
      OR EXISTS (
          SELECT 1
          FROM session_memberships sm
          WHERE sm.session_id = s.id
            AND sm.user_id = $2
      )
  )
);

-- name: CheckSessionActive :one
SELECT EXISTS (
    SELECT 1
    FROM sessions
    WHERE id = $1
        AND is_active = TRUE
);

-- name: CheckQuestionBelongsToSession :one
SELECT EXISTS (
    SELECT 1
    FROM questions
    WHERE id = $1
        AND session_id = $2
);

-- name: CheckCanDeleteQuestion :one
SELECT EXISTS (
    SELECT 1
    FROM questions q
    JOIN sessions s
        ON q.session_id = s.id
    WHERE q.id = $1
        AND (
            q.user_id = $2 OR
            s.owner_id = $2
        )
);

-- name: CheckCanDeleteQuestionLike :one
SELECT EXISTS (
    SELECT 1
    FROM question_likes
    WHERE user_id = $1
);