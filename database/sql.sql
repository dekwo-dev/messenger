CREATE TABLE comments (
    commentsid smallserial NOT NULL UNIQUE PRIMARY KEY,
    commentsauthor varchar(32) NOT NULL DEFAULT '',
    commentscontent varchar(256) NOT NULL DEFAULT '',
    commentstime timestamptz NOT NULL DEFAULT now()
)

CREATE TYPE operation AS ENUM ('DELETE', 'INSERT', 'UPDATE');

CREATE TABLE comments_changes_queue (
    queueid integer NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
    queuetime timestamptz NOT NULL DEFAULT now(),
    operation operation NOT NULL,
    commentsid smallint NOT NULL,
    commentsauthor varchar(32) NOT NULL,
    commentscontent varchar(256) NOT NULL,
    commentstime timestamptz NOT NULL
);

ALTER TABLE table_name OWNER TO owner_name;

CREATE OR REPLACE FUNCTION enqueue_comments_changes()
RETURNS TRIGGER AS $comments_changes_queue$
DECLARE enqueue_count integer; BEGIN
IF (TG_OP = 'DELETE') THEN
    INSERT INTO comments_changes_queue
    (operation, commentsid, commentsauthor, commentscontent, commentstime)
    SELECT 'DELETE', o.* FROM old_table o;
    SELECT INTO enqueue_count count(*) FROM old_table;
ELSIF (TG_OP = 'INSERT') THEN
    INSERT INTO comments_changes_queue
    (operation, commentsid, commentsauthor, commentscontent, commentstime)
    SELECT 'INSERT', n.* FROM new_table n;
    SELECT INTO enqueue_count count(*) FROM new_table;
ELSIF (TG_OP = 'UPDATE') THEN
    INSERT INTO comments_changes_queue
    (operation, commentsid, commentsauthor, commentscontent, commentstime)
    SELECT 'UPDATE', n.* FROM new_table n;
    SELECT INTO enqueue_count count(*) FROM new_table;
END IF;
PERFORM pg_notify('enqueue_comments_changes', 'enqueue:' || enqueue_count);
RETURN NULL;
END;
$comments_changes_queue$ LANGUAGE plpgsql;

CREATE TRIGGER on_comments_ins AFTER DELETE ON comments
REFERENCING OLD TABLE AS old_table
FOR EACH STATEMENT EXECUTE FUNCTION enqueue_comments_changes();

CREATE TRIGGER on_comments_ins AFTER INSERT ON comments
REFERENCING NEW TABLE AS new_table
FOR EACH STATEMENT EXECUTE FUNCTION enqueue_comments_changes();

CREATE TRIGGER on_comments_upd AFTER UPDATE ON comments
REFERENCING NEW TABLE AS new_table OLD TABLE as old_table
FOR EACH STATEMENT EXECUTE FUNCTION enqueue_comments_changes();
