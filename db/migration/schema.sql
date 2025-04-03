-- 用户表（学生/管理员）
CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(10) CHECK (role IN ('student', 'admin')) NOT NULL,
    department VARCHAR(50) NOT NULL,  -- 所属院系
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 自习室表
CREATE TABLE room (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    department VARCHAR(50) NOT NULL,  -- 所属院系
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    qr_code TEXT NOT NULL,  -- 每日更新的二维码路径
    is_active BOOLEAN NOT NULL DEFAULT FALSE -- 默认新建自习室不启用
);

-- 座位表
CREATE TABLE seat (
    id SERIAL PRIMARY KEY,
    room_id INT REFERENCES room(id) NOT NULL,
    number VARCHAR(10) NOT NULL,
    has_socket BOOLEAN NOT NULL DEFAULT FALSE,  -- 是否有插座, 默认无插座
    is_available BOOLEAN NOT NULL DEFAULT TRUE -- 默认新建座位可用
);

-- 预约记录表
CREATE TABLE reservation (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id INT REFERENCES "user"(id) NOT NULL,
    seat_id INT REFERENCES seat(id) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) CHECK (
        status IN ('reserved', 'completed', 'canceled', 'violated')
    ) NOT NULL DEFAULT 'reserved',
    checkin_time TIMESTAMP NOT NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 违约记录表
CREATE TABLE violation (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES "user"(id) NOT NULL,
    reservation_id UUID REFERENCES reservation(id) NOT NULL,
    reason VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 索引优化
CREATE INDEX idx_reservation_time ON reservation (start_time, end_time);
CREATE INDEX idx_seat_availability ON seat (is_available, has_socket);
CREATE INDEX idx_reservation_seat ON reservation (seat_id);
CREATE INDEX idx_violation_user ON violation (user_id);
CREATE INDEX idx_reservation_expire ON reservation (status, start_time) WHERE status = 'reserved';