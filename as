SELECT questions.public_id, questions.name, questions.read_duration, questions.answer_duration
FROM questions
JOIN positions ON questions.position_id = positions.id
JOIN user_interviews ON user_interviews.position_id = positions.id
JOIN interviews ON interviews.id = user_interviews.interview_id
WHERE interviews.public_id = '38e4fa88-b01d-4ea9-a59f-4e14aa5570e9';