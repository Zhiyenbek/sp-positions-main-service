-- Creating tables 
-- Enable the UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the tables
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT,
    photo TEXT
);

CREATE TABLE IF NOT EXISTS candidates (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL,
    current_position TEXT,
    education TEXT,
    resume TEXT,
    bio TEXT
);

CREATE TABLE IF NOT EXISTS recruiters (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL,
    company_public_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    name TEXT,
    description TEXT
);

CREATE TABLE IF NOT EXISTS positions (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    name TEXT,
    status int DEFAULT 0,
    recruiters_public_id UUID
);

CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    name TEXT
);

CREATE TABLE IF NOT EXISTS areas (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    position_id INT,
    name TEXT
);

CREATE TABLE IF NOT EXISTS interviews (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    results JSONB
);

CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    interviews_public_id UUID,
    path TEXT
);

CREATE TABLE IF NOT EXISTS auth (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,
    login TEXT UNIQUE,
    password TEXT
);

CREATE TABLE IF NOT EXISTS position_skills (
    position_id INT,
    skill_id INT,
    PRIMARY KEY (position_id, skill_id),
    CONSTRAINT fk_position_skills_positions FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE CASCADE,
    CONSTRAINT fk_position_skills_skills FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS candidate_skills (
    candidate_id INT,
    skill_id INT,
    PRIMARY KEY (candidate_id, skill_id),
    CONSTRAINT fk_candidate_skills_candidates FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE,
    CONSTRAINT fk_candidate_skills_skills FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_interviews (
    candidate_id INT,
    position_id INT,
    interview_id INT,
    PRIMARY KEY (candidate_id, position_id, interview_id),
    CONSTRAINT fk_user_interviews_candidates FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_interviews_positions FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_interviews_interviews FOREIGN KEY (interview_id) REFERENCES interviews(id) ON DELETE CASCADE
);

-- Creating references
ALTER TABLE recruiters ADD CONSTRAINT fk_recruiters_users FOREIGN KEY (public_id) REFERENCES users(public_id) ON DELETE CASCADE;
ALTER TABLE candidates ADD CONSTRAINT fk_candidates_users FOREIGN KEY (public_id) REFERENCES users(public_id) ON DELETE CASCADE;
ALTER TABLE auth ADD CONSTRAINT fk_auth_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE positions ADD CONSTRAINT fk_positions_recruiters FOREIGN KEY (recruiters_public_id) REFERENCES recruiters(public_id) ON DELETE CASCADE;
ALTER TABLE videos ADD CONSTRAINT fk_videos_interviews FOREIGN KEY (interviews_public_id) REFERENCES interviews(public_id) ON DELETE CASCADE;



INSERT INTO users (first_name, last_name, photo)
VALUES
    ('John', 'Doe', 'path/to/photo1'),
    ('Jane', 'Smith', 'path/to/photo2'),
    ('Michael', 'Johnson', 'path/to/photo3'),
    ('Emily', 'Williams', 'path/to/photo4'),
    ('David', 'Brown', 'path/to/photo5'),
    ('Olivia', 'Jones', 'path/to/photo6'),
    ('Daniel', 'Miller','path/to/photo7'),
    ('Sophia', 'Taylor','path/to/photo8'),
    ('Matthew', 'Anderson','path/to/photo9'),
    ('Ava', 'Thomas','path/to/photo10');


INSERT INTO candidates (public_id, current_position, resume, bio, education)
SELECT public_id, 'Software Engineer', 'John Doe Resume', 'John Doe Bio',  'MTI'
FROM users
WHERE id <= 5;

INSERT INTO companies (name, description)
VALUES
    ('Company A', 'A technology company that specializes in software development.'),
    ('Company B', 'A global retail company with a focus on e-commerce.'),
    ('Company C', 'A financial services company providing investment and banking solutions.');

INSERT INTO recruiters (public_id, company_public_id)
SELECT public_id, (SELECT public_id FROM companies WHERE name = 'Company A')
FROM users
WHERE id > 5;




INSERT INTO positions (public_id, name, recruiters_public_id)
SELECT public_id, 'Software Engineer', (SELECT public_id FROM recruiters WHERE id = 1)
FROM candidates;


INSERT INTO skills (name)
VALUES
    ('Java'),
    ('Python'),
    ('JavaScript'),
    ('SQL'),
    ('HTML'),
    ('CSS'),
    ('React'),
    ('Node.js'),
    ('AWS'),
    ('Agile Methodology');

INSERT INTO areas (position_id, name)
SELECT id, 'Area ' || id
FROM positions;


INSERT INTO interviews (public_id, results)
SELECT public_id, '{"result": "Pass"}'
FROM candidates;


INSERT INTO videos (public_id, interviews_public_id, path)
SELECT public_id, (SELECT public_id FROM interviews WHERE id = 1), '/path/to/video'
FROM candidates;


INSERT INTO auth (user_id, login, password)
SELECT id, 'user' || id, 'password' || id
FROM users;


insert into position_skills values (1,2);
insert into position_skills values (2,2);
insert into position_skills values (1,3);
insert into position_skills values (3,4);

insert into candidate_skills values (1,2);
insert into candidate_skills values (2,2);
insert into candidate_skills values (1,3);
insert into candidate_skills values (3,4);

INSERT INTO user_interviews (candidate_id, position_id, interview_id)
SELECT c.id AS candidate_id, p.id AS position_id, i.id AS interview_id
FROM candidates c
CROSS JOIN positions p
CROSS JOIN interviews i
WHERE c.id <= 5;