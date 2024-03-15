
INSERT INTO users (first_name, last_name)
VALUES
    ('John', 'Doe'),
    ('Jane', 'Smith'),
    ('Michael', 'Johnson'),
    ('Emily', 'Williams'),
    ('David', 'Brown'),
    ('Olivia', 'Jones'),
    ('Daniel', 'Miller'),
    ('Sophia', 'Taylor'),
    ('Matthew', 'Anderson'),
    ('Ava', 'Thomas');


INSERT INTO candidates (public_id, current_position, resume, bio)
SELECT public_id, 'Software Engineer', 'John Doe Resume', 'John Doe Bio'
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
WHERE id > 28;




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


INSERT INTO positions_skills (position_id, skill_id)
SELECT position_id, skill_id
FROM (
    SELECT p.id AS position_id, s.id AS skill_id, ROW_NUMBER() OVER () AS rn
    FROM positions p
    CROSS JOIN skills s
) AS sub
WHERE sub.rn <= 5;


INSERT INTO candidate_skills (candidate_id, skill_id)
SELECT candidate_id, skill_id
FROM (
    SELECT c.id AS candidate_id, s.id AS skill_id, ROW_NUMBER() OVER () AS rn
    FROM candidates c
    CROSS JOIN skills s
) AS sub
WHERE sub.rn <= 5;


INSERT INTO user_interviews (candidate_id, position_id, interview_id)
SELECT c.id AS candidate_id, p.id AS position_id, i.id AS interview_id
FROM candidates c
CROSS JOIN positions p
CROSS JOIN interviews i
WHERE c.id <= 5;