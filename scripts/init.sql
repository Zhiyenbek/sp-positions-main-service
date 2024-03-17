-- Creating tables 
-- Enable the UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the tables
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT,
    email TEXT,
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
    logo TEXT,
    description TEXT
);

CREATE TABLE IF NOT EXISTS positions (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT uuid_generate_v4() NOT NULL,
    description TEXT,
    name TEXT,
    status int DEFAULT 0,
    recruiter_public_id UUID NOT NULL
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
ALTER TABLE positions ADD CONSTRAINT fk_positions_recruiters FOREIGN KEY (recruiter_public_id) REFERENCES recruiters(public_id) ON DELETE CASCADE;
ALTER TABLE videos ADD CONSTRAINT fk_videos_interviews FOREIGN KEY (interviews_public_id) REFERENCES interviews(public_id) ON DELETE CASCADE;



INSERT INTO users (first_name, last_name, photo, email)
VALUES
    ('John', 'Doe', 'path/to/photo1', 'example@mail.com'),
    ('Jane', 'Smith', 'path/to/photo2', 'example@mail.com'),
    ('Michael', 'Johnson', 'path/to/photo3', 'example@mail.com'),
    ('Emily', 'Williams', 'path/to/photo4', 'example@mail.com'),
    ('David', 'Brown', 'path/to/photo5', 'example@mail.com'),
    ('Olivia', 'Jones', 'path/to/photo6', 'example@mail.com'),
    ('Daniel', 'Miller','path/to/photo7', 'example@mail.com'),
    ('Sophia', 'Taylor','path/to/photo8', 'example@mail.com'),
    ('Matthew', 'Anderson','path/to/photo9', 'example@mail.com'),
    ('Ava', 'Thomas','path/to/photo10', 'example@mail.com');


INSERT INTO candidates (public_id, current_position, resume, bio, education)
SELECT public_id, 'Software Engineer', 'John Doe Resume', 'John Doe Bio',  'MTI'
FROM users
WHERE id <= 5;

INSERT INTO companies (name, description, logo)
VALUES
    ('Company A', 'A technology company that specializes in software development.','path/to/logo1'),
    ('Company B', 'A global retail company with a focus on e-commerce.','path/to/logo2'),
    ('Company C', 'A financial services company providing investment and banking solutions.','path/to/logo3');

INSERT INTO recruiters (public_id, company_public_id)
SELECT public_id, (SELECT public_id FROM companies WHERE name = 'Company A')
FROM users
WHERE id > 5;




INSERT INTO positions (public_id, name, recruiter_public_id, description)
SELECT public_id, 'Software Engineer', (SELECT public_id FROM recruiters WHERE id = 1), 'This position is awesome'
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
SELECT public_id, '
{
  "questions": [
    {
      "question": "What is your experience with object-oriented programming?",
      "evaluation": "Good",
      "score": 8,
      "video_link": "https://example.com/video1",
      "emotion_results": [
        {
          "emotion": "Happiness",
          "exact_time": 24.5,
          "duration": 10.2
        },
        {
          "emotion": "Neutral",
          "exact_time": 36.2,
          "duration": 5.7
        }
      ]
    },
    {
      "question": "Describe a challenging project you have worked on.",
      "evaluation": "Excellent performance with exceptional problem-solving skills",
      "score": 9,
      "video_link": "https://example.com/video2",
      "emotion_results": [
        {
          "emotion": "Confidence",
          "exact_time": 45.8,
          "duration": 8.5
        },
        {
          "emotion": "Determination",
          "exact_time": 56.3,
          "duration": 7.1
        }
      ]
    }
  ],
  "score": 17,
  "video": "https://example.com/interview_video"
}'
FROM candidates;


INSERT INTO videos (public_id, interviews_public_id, path)
SELECT public_id, (SELECT public_id FROM interviews WHERE id = 1), '/path/to/video'
FROM candidates;


INSERT INTO auth (user_id, login, password)
SELECT id, 'user' || id, '$2a$12$TPhE59oXJf8TBvbDRiBghu7jcgVppHgYPLmZr7ePf9rjNwVWJJDuO'
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

