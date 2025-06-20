-- Bảng: roles
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    public_id TEXT
);

CREATE TABLE public.action_control_list (
    id SERIAL PRIMARY KEY,
    action_id TEXT NOT NULL,
    role_id TEXT,
    user_id TEXT,
    status INT NOT NULL DEFAULT 1
);

-- Bảng: users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role UUID NOT NULL REFERENCES roles(id),
    status TEXT NOT NULL CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ
);

-- Bảng: videos
CREATE TABLE videos (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    title TEXT NOT NULL,
    filepath TEXT NOT NULL,
    duration INT NOT NULL,
    width INT,
    height INT,
    folder TEXT,
    format TEXT,
    metadata JSONB,
    thumbnail_url TEXT
);


-- Bảng: characters
CREATE TABLE characters (
    id UUID PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id),
    name TEXT NOT NULL,
    display_name TEXT,
    image_url TEXT,
    color TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    detection_type TEXT CHECK (detection_type IN ('manual', 'automatic')),
    confidence FLOAT,
    metadata JSONB
);


-- Bảng: character_appearances
CREATE TABLE character_appearances (
    id UUID PRIMARY KEY,
    character_id UUID NOT NULL REFERENCES characters(id),
    video_id UUID NOT NULL REFERENCES videos(id),
    start_time FLOAT NOT NULL,
    end_time FLOAT NOT NULL,
    start_frame INT,
    end_frame INT,
    confidence FLOAT CHECK (confidence >= 0.0 AND confidence <= 1.0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    metadata JSONB
);


-- Bảng: segments
CREATE TABLE segments (
    id UUID PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id),
    label TEXT,
    description TEXT,
    start_time FLOAT NOT NULL,
    end_time FLOAT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    tags JSONB,
    metadata JSONB
);


-- Bảng: segment_characters
CREATE TABLE segment_characters (
    id UUID PRIMARY KEY,
    segment_id UUID NOT NULL REFERENCES segments(id),
    character_id UUID NOT NULL REFERENCES characters(id),
    is_included BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id)
);
