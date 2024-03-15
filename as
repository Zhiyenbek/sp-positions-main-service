SELECT
        c.public_id,
        c.current_position,
        c.education,
        c.resume,
        c.bio,
        u.first_name,
        u.last_name,
        array_agg(s.name) AS skills
FROM
        candidates c
JOIN
        users u ON c.public_id = u.public_id
JOIN
        candidate_skills cs ON c.id = cs.candidate_id
JOIN
        skills s ON cs.skill_id = s.id
WHERE
        (LOWER(u.first_name) LIKE LOWER('%%') OR LOWER(u.last_name) LIKE LOWER('%%'))
        AND c.id IN (
                                        SELECT DISTINCT cs2.candidate_id
                                        FROM candidate_skills cs2
                                        JOIN skills s2 on cs2.skill_id = s2.id
                                        WHERE s2.name  ANY('{"HTML"}')
                                ) GROUP BY c.id, u.first_name, u.last_name LIMIT 0  OFFSET 10;