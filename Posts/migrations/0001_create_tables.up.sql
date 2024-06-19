CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       name VARCHAR(100) NOT NULL
);

-- Create posts table
CREATE TABLE posts (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       title VARCHAR(100) NOT NULL,
                       content TEXT NOT NULL,
                       allow_comments BOOLEAN NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       author_id UUID NOT NULL,
                       CONSTRAINT fk_author FOREIGN KEY(author_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- Create comments table
CREATE TABLE comments (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          content TEXT NOT NULL,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          post_id UUID NOT NULL,
                          parent_id UUID,
                          author_id UUID NOT NULL,
                          CONSTRAINT fk_post FOREIGN KEY(post_id) REFERENCES posts(id) ON UPDATE CASCADE ON DELETE CASCADE,
                          CONSTRAINT fk_parent_comment FOREIGN KEY(parent_id) REFERENCES comments(id) ON UPDATE CASCADE ON DELETE CASCADE,
                          CONSTRAINT fk_author_comment FOREIGN KEY(author_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);
