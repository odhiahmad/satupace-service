-- Performance indexes for run-sync application
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active, is_verified);

-- Runner profiles indexes
CREATE INDEX IF NOT EXISTS idx_runner_profiles_user_id ON runner_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_runner_profiles_experience ON runner_profiles(running_experience_level);

-- Run groups indexes
CREATE INDEX IF NOT EXISTS idx_run_groups_status ON run_groups(status);
CREATE INDEX IF NOT EXISTS idx_run_groups_scheduled ON run_groups(scheduled_at) WHERE status = 'open';
CREATE INDEX IF NOT EXISTS idx_run_groups_created_by ON run_groups(created_by);
CREATE INDEX IF NOT EXISTS idx_run_groups_location ON run_groups(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_run_groups_women_only ON run_groups(is_women_only, status);

-- Run group members indexes
CREATE INDEX IF NOT EXISTS idx_run_group_members_group_id ON run_group_members(run_group_id);
CREATE INDEX IF NOT EXISTS idx_run_group_members_user_id ON run_group_members(user_id);
CREATE INDEX IF NOT EXISTS idx_run_group_members_status ON run_group_members(status);

-- Run activities indexes
CREATE INDEX IF NOT EXISTS idx_run_activities_user_id ON run_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_run_activities_created_at ON run_activities(created_at DESC);

-- Direct matches indexes
CREATE INDEX IF NOT EXISTS idx_direct_matches_user1 ON direct_matches(user_id_1);
CREATE INDEX IF NOT EXISTS idx_direct_matches_user2 ON direct_matches(user_id_2);
CREATE INDEX IF NOT EXISTS idx_direct_matches_status ON direct_matches(status);
CREATE INDEX IF NOT EXISTS idx_direct_matches_users ON direct_matches(user_id_1, user_id_2);

-- Direct chat messages indexes
CREATE INDEX IF NOT EXISTS idx_direct_chat_match_id ON direct_chat_messages(match_id);
CREATE INDEX IF NOT EXISTS idx_direct_chat_sender ON direct_chat_messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_direct_chat_created ON direct_chat_messages(match_id, created_at DESC);

-- Group chat messages indexes
CREATE INDEX IF NOT EXISTS idx_group_chat_group_id ON group_chat_messages(group_id);
CREATE INDEX IF NOT EXISTS idx_group_chat_sender ON group_chat_messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_group_chat_created ON group_chat_messages(group_id, created_at DESC);

-- User photos indexes
CREATE INDEX IF NOT EXISTS idx_user_photos_user_id ON user_photos(user_id);
CREATE INDEX IF NOT EXISTS idx_user_photos_primary ON user_photos(user_id, is_primary);

-- Safety logs indexes
CREATE INDEX IF NOT EXISTS idx_safety_logs_user_id ON safety_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_safety_logs_group_id ON safety_logs(run_group_id);
CREATE INDEX IF NOT EXISTS idx_safety_logs_created ON safety_logs(created_at DESC);
