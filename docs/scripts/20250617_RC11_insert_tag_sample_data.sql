-- Insert sample data for tag system

-- 1. Insert tag positions
INSERT INTO tag_positions (title, position, description, is_active, sort_order) VALUES
('Video Theme', 'video_theme', 'Tags for video theme filtering', true, 1),
('People', 'people', 'Tags for people/character filtering', true, 2);

-- 2. Insert tag categories
INSERT INTO tag_categories (name, code, description, color, priority, is_shown, is_system_category, filter_type) VALUES
('Gender', 'gender', 'Gender classification', '#007bff', 1, true, true, 'multiple'),
('Age Range', 'age_range', 'Age range classification', '#28a745', 2, true, true, 'single');

-- 3. Insert tags
-- Gender tags (category_id = 1)
INSERT INTO tags (category_id, name, code, description, color, sort_order, usage_count, is_active, is_system_tag) VALUES
(1, 'Male', 'male', 'Male gender', '#007bff', 1, 150, true, true),
(1, 'Female', 'female', 'Female gender', '#dc3545', 2, 180, true, true),
(1, 'Female & Male', 'female_male', 'Both female and male', '#6f42c1', 3, 95, true, true);

-- Age Range tags (category_id = 2)
INSERT INTO tags (category_id, name, code, description, color, sort_order, usage_count, is_active, is_system_tag) VALUES
(2, 'Child', 'child', 'Child age range', '#28a745', 1, 75, true, true),
(2, 'Teen', 'teen', 'Teen age range', '#17a2b8', 2, 120, true, true),
(2, 'Adult', 'adult', 'Adult age range', '#ffc107', 3, 200, true, true);

-- 4. Insert position-category mappings
-- Video Theme position (id=1) -> Gender category (id=1)
INSERT INTO tag_position_categories (tag_position_id, tag_category_id, sort_order, is_visible, display_style) VALUES
(1, 1, 1, true, 'checkbox');

-- People position (id=2) -> Age Range category (id=2)
INSERT INTO tag_position_categories (tag_position_id, tag_category_id, sort_order, is_visible, display_style) VALUES
(2, 2, 1, true, 'radio'); 