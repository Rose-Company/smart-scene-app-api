CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status TEXT NOT NULL DEFAULT 'pending',
    title TEXT NOT NULL,
    file_path TEXT NOT NULL,
    duration INTEGER NOT NULL,
    thumbnail_url TEXT,
    has_character_analysis BOOLEAN DEFAULT FALSE,
    character_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    metadata JSONB
);

-- 2. CHARACTERS (UUID - cần cho security)
CREATE TABLE characters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    avatar TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    gender TEXT CHECK (gender IN ('male', 'female', 'other', 'unknown')),
    character_type TEXT CHECK (character_type IN ('person', 'animal', 'cartoon', 'object', 'other')) DEFAULT 'person',
    character_image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    metadata JSONB
);

-- =============================================================================
-- HỆ THỐNG TAGS PHÂN CẤP (INT - hiệu suất cao)
-- =============================================================================

-- 3. TAG_POSITIONS (INT - vị trí hiển thị tags)
-- Giống như bảng đầu tiên trong hình của bạn
CREATE TABLE tag_positions (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,                    -- "Menu top", "Side bar left"
    position TEXT NOT NULL UNIQUE,          -- "menu_top", "side_bar_left"
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 4. TAG_CATEGORIES (INT - loại tags, giống tag_search của bạn)  
CREATE TABLE tag_categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,              -- "Gender", "Age Range", "Character Type"
    code TEXT NOT NULL UNIQUE,              -- "gender", "age_range", "character_type"
    description TEXT,
    color TEXT DEFAULT '#007bff',
    priority INTEGER DEFAULT 0,             -- Thứ tự hiển thị
    is_shown BOOLEAN DEFAULT TRUE,
    is_system_category BOOLEAN DEFAULT FALSE,
    filter_type TEXT CHECK (filter_type IN ('single', 'multiple', 'range')) DEFAULT 'single',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id),
    metadata JSONB
);

-- 5. TAGS (INT - tags cụ thể)
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES tag_categories(id),
    name TEXT NOT NULL,                     -- "Male", "Female", "Child"
    code TEXT NOT NULL,                     -- "male", "female", "child"
    description TEXT,
    color TEXT,                             -- Inherit từ category nếu NULL
    sort_order INTEGER DEFAULT 0,
    usage_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    is_system_tag BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id),
    metadata JSONB,
    UNIQUE(category_id, code)               -- Unique trong cùng category
);

-- 6. TAG_POSITION_CATEGORIES (Mapping - giống tag_position_tag_search)
CREATE TABLE tag_position_categories (
    id SERIAL PRIMARY KEY,
    tag_position_id INTEGER NOT NULL REFERENCES tag_positions(id),
    tag_category_id INTEGER NOT NULL REFERENCES tag_categories(id),
    sort_order INTEGER DEFAULT 0,
    is_visible BOOLEAN DEFAULT TRUE,
    display_style TEXT CHECK (display_style IN ('dropdown', 'checkbox', 'radio', 'chips')) DEFAULT 'checkbox',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(tag_position_id, tag_category_id)
);

-- 7. VIDEO_TAGS (Mapping videos với tags)
CREATE TABLE video_tags (
    id SERIAL PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id),
    tag_id INTEGER NOT NULL REFERENCES tags(id),
    character_id UUID REFERENCES characters(id)     
);

-- 8. CHARACTER_APPEARANCES
CREATE TABLE character_appearances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id UUID NOT NULL REFERENCES characters(id),
    video_id UUID NOT NULL REFERENCES videos(id),
    start_time DECIMAL(10,3) NOT NULL,
    end_time DECIMAL(10,3) NOT NULL,
    duration DECIMAL(10,3) DEFAULT 0,
    confidence DECIMAL(3,2) DEFAULT 0.0,
    is_confirmed BOOLEAN DEFAULT FALSE,
    color_shown TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES users(id)
);