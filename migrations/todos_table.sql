-- Create todos table with Vietnamese text support
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    name VARCHAR(500) NOT NULL,  -- Support for Vietnamese characters
    create_time DATE NOT NULL DEFAULT CURRENT_DATE,
    importance INTEGER NOT NULL DEFAULT 0,  -- Higher number = higher importance
    flag BOOLEAN NOT NULL DEFAULT FALSE,    -- true = completed, false = not completed
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index for better performance on common queries
CREATE INDEX idx_todos_importance ON todos(importance DESC);
CREATE INDEX idx_todos_flag ON todos(flag);
CREATE INDEX idx_todos_create_time ON todos(create_time);

-- Add trigger to auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_todos_updated_at 
    BEFORE UPDATE ON todos 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();