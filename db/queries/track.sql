-- name: SaveTrack :exec
insert into track (user_id, engine, market, board_group, security, date, tracked_volume) values (
    $1, $2, $3, $4, $5, $6, $7
);
-- name: GetUserTracks :many
select * from track where user_id = $1;
-- name: DeleteTrackByUserIdAndId :exec
delete from track where user_id = $1 and id = $2;
-- name: TrackSecurityByUserIdAndId :exec
update track set is_tracked = true where user_id = $1 and id = $2;
-- name: UntrackSecurityByUserIdAndId :exec
update track set is_tracked = false where user_id = $1 and id = $2;
-- name: GetAllMustNotifiedTracks :many
select t.engine, t.market, t.board_group, t.security, t.tracked_volume, u.id from track as t inner join users as u on t.user_id = u.id where u.followed = true and u.banned = false and t.is_tracked = true;