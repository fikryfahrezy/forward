# Social Media Database Exercise

## Disclaimer

The author doesn't have any experience that relate to Social Media domain nor designin the system before.

## Database Schema

### Users Table

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  username VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(100) NOT NULL,
  phone VARCHAR(20) NOT NULL,
  bio VARCHAR(20) NOT NULL DEFAULT '',
  gender VARCHAR(50) NOT NULL,
  is_verified BOOLEAN DEFAULT FALSE,
  birth_date TIMESTAMPTZ,
  status VARCHAR(20) DEFAULT 'active', -- 'active', 'suspend'
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
```

### Credential Table

```sql
CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password_hash VARCHAR(255) NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Followers Table

```sql
CREATE TABLE followers (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

CREATE INDEX idx_followers_following ON followers(following_id);
```

### Posts Table

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    repost_of UUID REFERENCES posts(id) ON DELETE CASCADE,  -- NULL = original, NOT NULL = repost
    caption TEXT DEFAULT '',  -- doubles as "quote" for reposts
    reactions_count BIGINT NOT NULL DEFAULT 0,
    comments_count BIGINT NOT NULL DEFAULT 0,
    shares_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_posts_user ON posts(user_id);
CREATE INDEX idx_posts_repost ON posts(repost_of);
CREATE UNIQUE INDEX idx_unique_repost ON posts(user_id, repost_of) WHERE repost_of IS NOT NULL;  -- prevent duplicate reposts
```

### Post Attachments Table

```sql
CREATE TABLE post_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),        -- for videos
    type VARCHAR(20) NOT NULL,         -- 'image', 'video'
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_post_attachments_post ON post_attachments(post_id);
```

### Post Reactions Table

```sql
CREATE TABLE post_reactions (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    reaction VARCHAR(20) NOT NULL,  -- 'like', 'love', 'haha', 'wow', 'sad', 'angry'
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

CREATE INDEX idx_post_reactions_post ON post_reactions(post_id);
```

### Saved Posts Table

```sql
CREATE TABLE saved_posts (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);
```

### Post Comments Table

```sql
CREATE TABLE post_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_comment_id UUID REFERENCES post_comments(id) ON DELETE CASCADE,  -- for replies
    comment TEXT DEFAULT '',
    attachment_url VARCHAR(500),
    reactions_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_post_comments_post ON post_comments(post_id);
CREATE INDEX idx_post_comments_parent ON post_comments(parent_comment_id);
```

### Comment Reactions Table

```sql
CREATE TABLE comment_reactions (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment_id UUID NOT NULL REFERENCES post_comments(id) ON DELETE CASCADE,
    reaction VARCHAR(20) NOT NULL,  -- 'like', 'love', 'haha', 'wow', 'sad', 'angry'
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (comment_id, user_id)  -- optimized for "reactions on this comment"
);
```

### Post Comment Histories Table

```sql
CREATE TABLE post_comment_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_comment_id UUID NOT NULL REFERENCES post_comments(id) ON DELETE CASCADE,
    comment TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_post_comment_histories_comment ON post_comment_histories(post_comment_id);
```

### Post Hashtags Table

```sql
CREATE TABLE post_hashtags (
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    hashtag VARCHAR(100) NOT NULL,
    PRIMARY KEY (post_id, hashtag)
);

CREATE INDEX idx_post_hashtags_hashtag ON post_hashtags(hashtag);
```

### Post Mentions Table

```sql
CREATE TABLE post_mentions (
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    mentioned_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, mentioned_user_id)
);

CREATE INDEX idx_post_mentions_user ON post_mentions(mentioned_user_id);
```

### Conversations Table

```sql
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(20) NOT NULL DEFAULT 'direct',  -- 'direct', 'group'
    name VARCHAR(100),  -- NULL for direct, name for group
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Conversation Members Table

```sql
CREATE TABLE conversation_members (
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX idx_conversation_members_user ON conversation_members(user_id);
```

### Messages Table

```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT DEFAULT '',
    attachments JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_messages_conversation ON messages(conversation_id);
CREATE INDEX idx_messages_sender ON messages(sender_id);
```

### Message Histories Table

```sql
CREATE TABLE message_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    message TEXT DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_message_histories_message ON message_histories(message_id);
```

## References

- https://facebook.com
- https://instagram.com
- https://x.com