-- 用户表（学生/管理员）
CREATE TABLE user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(10) CHECK (role IN ('student', 'admin')) NOT NULL,
    department VARCHAR(50),  -- 所属院系
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 自习室表
CREATE TABLE room (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    department VARCHAR(50),  -- 所属院系
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    qr_code TEXT,  -- 每日更新的二维码路径
    is_active BOOLEAN DEFAULT TRUE
);

-- 座位表
CREATE TABLE seat (
    id SERIAL PRIMARY KEY,
    room_id INT REFERENCES room(id),
    number VARCHAR(10) NOT NULL,
    has_socket BOOLEAN DEFAULT FALSE,  -- 是否有插座
    is_available BOOLEAN DEFAULT TRUE
);

-- 预约记录表
CREATE TABLE reservation (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT REFERENCES user(id),
    seat_id INT REFERENCES seat(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) CHECK (
        status IN ('reserved', 'completed', 'canceled', 'violated')
    ) DEFAULT 'reserved',
    checkin_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 违约记录表
CREATE TABLE violation (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES user(id),
    reservation_id UUID REFERENCES reservation(id),
    reason VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引优化
CREATE INDEX idx_reservation_time ON reservation (start_time, end_time);
CREATE INDEX idx_seat_availability ON seat (is_available, has_socket);
CREATE INDEX idx_reservation_seat ON reservation (seat_id);
CREATE INDEX idx_violation_user ON violation (user_id);
CREATE INDEX idx_reservation_expire ON reservation (status, start_time) WHERE status = 'reserved';