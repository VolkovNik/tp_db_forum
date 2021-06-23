CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE UNLOGGED TABLE users
(
    nickname    CITEXT  PRIMARY KEY,
    fullname    TEXT    NOT NULL,
    about       TEXT    NOT NULL,
    email       CITEXT  NOT NULL UNIQUE
);

CREATE UNIQUE INDEX ON users (nickname, email);
CREATE UNIQUE INDEX ON users (nickname, email, about, fullname);
CREATE UNIQUE INDEX ON users (nickname DESC);

CREATE UNLOGGED TABLE forums
(
    slug    CITEXT PRIMARY KEY,
    title   TEXT    NOT NULL,
    author  CITEXT  NOT NULL,
    posts   INT DEFAULT 0,
    threads INT DEFAULT 0,
    CONSTRAINT fk_author FOREIGN KEY(author) REFERENCES users(nickname) ON DELETE CASCADE
);

CREATE UNLOGGED TABLE forums_users
(
    author CITEXT NOT NULL,
    slug   CITEXT NOT NULL,
    CONSTRAINT fk_author FOREIGN KEY(author) REFERENCES users(nickname) ON DELETE CASCADE,
    CONSTRAINT fk_slug   FOREIGN KEY(slug)   REFERENCES forums(slug)    ON DELETE CASCADE,
    PRIMARY KEY(author, slug)
);

CREATE INDEX ON forums_users (slug);
CREATE INDEX ON forums_users (author);

CREATE UNLOGGED TABLE threads
(
    id          SERIAL  PRIMARY KEY,
    title       TEXT    NOT NULL,
    author      CITEXT  NOT NULL,
    forum       CITEXT  NOT NULL,
    message     TEXT    NOT NULL,
    slug        CITEXT  UNIQUE,
    created     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    votes       INT     NOT NULL DEFAULT 0,
    CONSTRAINT fk_forum  FOREIGN KEY(forum)  REFERENCES forums(slug)    ON DELETE CASCADE,
    CONSTRAINT fk_author FOREIGN KEY(author) REFERENCES users(nickname) ON DELETE CASCADE
);

CREATE UNLOGGED TABLE posts
(
    id          SERIAL  PRIMARY KEY,
    parent      INT     NOT NULL DEFAULT 0,
    author      CITEXT  NOT NULL,
    message     TEXT    NOT NULL,
    isEdited    BOOLEAN NOT NULL DEFAULT FALSE,
    forum       CITEXT  NOT NULL,
    thread      INT     NOT NULL,
    created     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    path        INT ARRAY NOT NULL,
        CONSTRAINT fk_forum  FOREIGN KEY(forum)  REFERENCES forums(slug)    ON DELETE CASCADE,
    CONSTRAINT fk_thread FOREIGN KEY(thread) REFERENCES threads(id)     ON DELETE CASCADE,
    CONSTRAINT fk_author FOREIGN KEY(author) REFERENCES users(nickname) ON DELETE CASCADE
);

CREATE INDEX ON posts(thread);
CREATE INDEX ON posts(author);
CREATE INDEX ON posts(thread, path DESC);
CREATE INDEX ON posts(thread, path ASC);
CREATE INDEX ON posts(thread, id DESC);

CREATE UNLOGGED TABLE votes
(
    nickname    CITEXT     NOT NULL,
    thread_id   INT        NOT NULL,
    vote        SMALLINT   NOT NULL,
    PRIMARY KEY (nickname, thread_id),
    CONSTRAINT fk_nickname  FOREIGN KEY(nickname)   REFERENCES users(nickname) ON DELETE CASCADE,
    CONSTRAINT fk_thread_id FOREIGN KEY(thread_id)  REFERENCES threads(id)     ON DELETE CASCADE
);

CREATE FUNCTION declare_path() RETURNS TRIGGER AS $declare_path$
DECLARE
    temp INT ARRAY;
BEGIN
    IF new.parent ISNULL OR new.parent = 0 THEN new.path = ARRAY [new.id];
    ELSE
        SELECT path INTO temp FROM posts WHERE id = new.parent;
        new.path = temp || new.id;

    END IF;
    RETURN new;
END;
$declare_path$ LANGUAGE plpgsql;

CREATE TRIGGER update_posts_path BEFORE INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE declare_path();

CREATE FUNCTION update_after_insert_vote() RETURNS TRIGGER AS $update_after_insert_vote$
BEGIN
    UPDATE threads SET
        votes = votes + NEW.vote
    WHERE threads.id = NEW.thread_id;
    RETURN NEW;
END;
$update_after_insert_vote$ LANGUAGE plpgsql;

CREATE TRIGGER update_after_insert_vote AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE update_after_insert_vote();

CREATE FUNCTION update_after_update_vote() RETURNS TRIGGER AS $update_after_update_vote$
BEGIN
    UPDATE threads SET
        votes = votes + NEW.vote - OLD.vote
    WHERE threads.id = NEW.thread_id;
    RETURN NEW;
END;
$update_after_update_vote$ LANGUAGE plpgsql;

CREATE TRIGGER update_after_update_vote AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE update_after_update_vote();

CREATE FUNCTION update_threads_counter() RETURNS TRIGGER AS $update_threads_counter$
BEGIN
    UPDATE forums SET
                      threads = threads + 1
    WHERE forums.slug = NEW.forum;
    RETURN NEW;
END;
$update_threads_counter$ LANGUAGE plpgsql;

CREATE TRIGGER update_threads_counter AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE update_threads_counter();

CREATE FUNCTION update_posts_counter() RETURNS TRIGGER AS $update_posts_counter$
BEGIN
    UPDATE forums SET
        posts = posts + 1
    WHERE forums.slug = NEW.forum;
    RETURN NEW;
END;
$update_posts_counter$ LANGUAGE plpgsql;

CREATE TRIGGER update_posts_counter AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE update_posts_counter();

CREATE FUNCTION insert_user_to_forums_users() RETURNS TRIGGER AS $insert_user_to_forums_users$
BEGIN
    INSERT INTO forums_users (author, slug)
        VALUES (NEW.author, NEW.forum)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$insert_user_to_forums_users$ LANGUAGE plpgsql;

CREATE TRIGGER insert_user_to_forums_users AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE insert_user_to_forums_users();

CREATE TRIGGER insert_user_to_forums_users AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE insert_user_to_forums_users();

VACUUM ANALYSE;
