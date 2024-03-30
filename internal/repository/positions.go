package repository

import (
	"context"
	"errors"

	"github.com/Zhiyenbek/sp-positions-main-service/config"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type positionRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

func NewPositionRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) PositionRepository {
	return &positionRepository{
		db:     db,
		logger: logger,
		cfg:    cfg,
	}
}

func (r *positionRepository) GetAllPositions(search string, pageNum, pageSize int) ([]models.Position, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	offset := (pageNum - 1) * pageSize
	// Query to retrieve positions for the current page
	query := `
		SELECT p.public_id, p.name, p.status, c.public_id, c.description, p.description, c.name, c.logo
		FROM positions AS p
		INNER JOIN recruiters r ON p.recruiter_public_id = r.public_id
		INNER JOIN companies c ON r.company_public_id = c.public_id
		WHERE p.name ILIKE $1
		GROUP BY p.public_id, p.name, p.status, c.public_id, c.name, c.logo, c.description, p.description
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, "%"+search+"%", pageSize, offset)
	if err != nil {
		r.logger.Errorf("Error occurred while retrieving positions: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	positions := []models.Position{}
	for rows.Next() {
		position := models.Position{}
		company := &models.Company{}
		err := rows.Scan(
			&position.PublicID,
			&position.Name,
			&position.Status,
			&company.PublicID,
			&company.Description,
			&position.Description,
			&company.Name,
			&company.Logo,
		)

		position.Company = company
		if err != nil {
			r.logger.Errorf("Error occurred while scanning position rows: %v", err)
			return nil, 0, err
		}
		query := `SELECT array_agg(s.name) from skills AS s
		INNER JOIN position_skills ps ON ps.skill_id = s.id
		INNER JOIN positions p ON ps.position_id = p.id
		WHERE p.public_id = $1`
		if position.PublicID != nil {
			err = r.db.QueryRow(ctx, query, *position.PublicID).Scan(&position.Skills)
			if err != nil {
				r.logger.Errorf("Error retrieving positions by company: %v", err)
				return nil, 0, err
			}
		}
		positions = append(positions, position)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error occurred while iterating over position rows: %v", err)
		return nil, 0, err
	}

	// Query to retrieve the total count of positions
	countQuery := `
		SELECT COUNT(*) FROM positions WHERE name ILIKE $1
	`

	var totalCount int
	err = r.db.QueryRow(ctx, countQuery, "%"+search+"%").Scan(&totalCount)
	if err != nil {
		r.logger.Errorf("Error occurred while retrieving position count: %v", err)
		return nil, 0, err
	}

	return positions, totalCount, nil
}

func (r *positionRepository) GetPositionInterviews(publicID string, pageNum int, pageSize int) ([]models.InterviewResults, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `
		SELECT i.public_id, i.results
		FROM interviews i
		INNER JOIN user_interviews ui ON ui.interview_id = i.id
		INNER JOIN positions p ON p.id = ui.position_id
		WHERE p.public_id = $1
		GROUP BY i.public_id, i.results
		LIMIT $2 OFFSET $3;
	`
	offset := (pageNum - 1) * pageSize
	rows, err := r.db.Query(ctx, query, publicID, pageSize, offset)
	if err != nil {
		r.logger.Errorf("Error occurred while retrieving interview result: %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	res := make([]models.InterviewResults, 0)
	for rows.Next() {
		var resultBytes []byte
		result := models.InterviewResults{}
		err = rows.Scan(
			&result.PublicID,
			&resultBytes,
		)
		if err != nil {
			r.logger.Errorf("Error occurred while retrieving interview result: %v", err)
			return nil, 0, err
		}
		result.Result = resultBytes
		res = append(res, result)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error occurred while iterating over interview result for position rows: %v", err)
		return nil, 0, err
	}
	query = `
		SELECT COUNT(*)
		FROM interviews i
		INNER JOIN user_interviews ui ON ui.interview_id = i.id
		INNER JOIN positions p ON p.id = ui.position_id
		WHERE p.public_id = $1
	`
	var totalCount int
	err = r.db.QueryRow(ctx, query, publicID).Scan(&totalCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res, 0, nil
		}
		r.logger.Errorf("Error occurred while retrieving position count: %v", err)
		return nil, 0, err
	}
	return res, totalCount, nil

}

func (r *positionRepository) GetPosition(publicID string) (*models.Position, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `SELECT p.public_id, p.name, p.status, p.description, c.public_id, c.name, c.description, r.public_id, c.logo, array_agg(s.name)
	FROM skills s
	INNER JOIN position_skills ps ON ps.skill_id = s.id
	INNER JOIN positions p ON ps.position_id = p.id
	INNER JOIN recruiters r ON p.recruiter_public_id = r.public_id
	INNER JOIN users u ON r.public_id = u.public_id
	INNER JOIN companies c ON r.company_public_id = c.public_id
	WHERE p.public_id = $1
	GROUP BY p.public_id, p.name, p.status, p.description, c.public_id, c.name, c.description, r.public_id, c.logo`
	res := &models.Position{
		Company: &models.Company{},
	}
	err := r.db.QueryRow(ctx, query, publicID).Scan(
		&res.PublicID,
		&res.Name,
		&res.Status,
		&res.Description,
		&res.Company.PublicID,
		&res.Company.Name,
		&res.Company.Description,
		&res.RecruiterPublicID,
		&res.Company.Logo,
		&res.Skills,
	)
	if err != nil {
		r.logger.Errorf("Error occurred while getting position: %v", err)
		return res, err
	}
	return res, nil
}

func (r *positionRepository) Exists(publicID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM positions WHERE public_id = $1)`

	err := r.db.QueryRow(ctx, query, publicID).Scan(&exists)
	if err != nil {
		r.logger.Errorf("Error occurred while checking position existence: %v", err)
		return false, err
	}

	return exists, nil
}

func (r *positionRepository) CreatePosition(position *models.Position) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()
	var id int64
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Error occurred while starting transaction: %v", err)
		return "", err
	}

	// Insert the position into the positions table
	insertPositionQuery := `INSERT INTO positions (description, name, status, recruiter_public_id)
		VALUES ($1, $2, $3, $4) RETURNING public_id, id`
	row := tx.QueryRow(ctx, insertPositionQuery, position.Description, position.Name, position.Status, position.RecruiterPublicID)
	if err := row.Scan(&position.PublicID, &id); err != nil {
		r.logger.Errorf("Error occurred while creating position: %v", err)
		tx.Rollback(ctx)
		return "", err
	}

	// Loop over the skills array
	for _, skillName := range position.Skills {
		// Check if the skill already exists in the database
		var skillID int
		query := `
		SELECT id FROM skills WHERE name = $1
		`
		err := tx.QueryRow(ctx, query, skillName).Scan(&skillID)
		if err != nil {
			if err == pgx.ErrNoRows {
				// Skill doesn't exist, so insert it into the database
				insertQuery := `
				INSERT INTO skills (name) VALUES ($1)
				RETURNING id
				`
				err = tx.QueryRow(ctx, insertQuery, skillName).Scan(&skillID)
				if err != nil {
					r.logger.Errorf("Error inserting new skill: %v", err)
					return "", err
				}
			} else {
				r.logger.Errorf("Error checking skill existence: %v", err)
				return "", err
			}
		}

		insertQuery := `
		INSERT INTO position_skills (position_id, skill_id) VALUES (
			$1,
			$2
		)
		`
		_, err = tx.Exec(ctx, insertQuery, id, skillID)
		if err != nil {
			r.logger.Errorf("Error adding skill to position: %v", err)
			return "", nil
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Errorf("Error occurred while committing transaction: %v", err)
		return "", err
	}

	return *position.PublicID, nil
}

func (r *positionRepository) UpdatePosition(position *models.Position) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Error occurred while starting transaction: %v", err)
		return err
	}

	// Update the position in the positions table
	updatePositionQuery := `
		UPDATE positions
		SET description = COALESCE($1, description),
			name = COALESCE($2, name),
			status = COALESCE($3, status),
			recruiter_public_id = COALESCE($4, recruiter_public_id)
		WHERE public_id = $5
	`
	_, err = tx.Exec(ctx, updatePositionQuery, position.Description, position.Name, position.Status, position.RecruiterPublicID, position.PublicID)
	if err != nil {
		r.logger.Errorf("Error occurred while updating position: %v", err)
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Errorf("Error occurred while committing transaction: %v", err)
		return err
	}

	return nil
}

func (r *positionRepository) CreateSkillsForPosition(positionPublicID string, skills []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Error beginning transaction: %v", err)
		return err
	}

	// Loop over the skills array
	for _, skillName := range skills {
		// Check if the skill already exists in the database
		var skillID int
		query := `
			SELECT id FROM skills WHERE name = $1
		`
		err := tx.QueryRow(ctx, query, skillName).Scan(&skillID)
		if err != nil {
			if err == pgx.ErrNoRows {
				// Skill doesn't exist, so insert it into the database
				insertQuery := `
					INSERT INTO skills (name) VALUES ($1)
					RETURNING id
				`
				err = tx.QueryRow(ctx, insertQuery, skillName).Scan(&skillID)
				if err != nil {
					tx.Rollback(ctx)
					r.logger.Errorf("Error inserting new skill: %v", err)
					return err
				}
			} else {
				tx.Rollback(ctx)
				r.logger.Errorf("Error checking skill existence: %v", err)
				return err
			}
		}

		// Associate the skill with the position
		insertQuery := `
			INSERT INTO position_skills (position_id, skill_id) VALUES (
				(SELECT id FROM positions WHERE public_id = $1),
				$2
			) ON CONFLICT DO NOTHING
		`
		_, err = tx.Exec(ctx, insertQuery, positionPublicID, skillID)
		if err != nil {
			tx.Rollback(ctx)
			r.logger.Errorf("Error adding skill to position: %v", err)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (r *positionRepository) DeleteSkillsFromPosition(positionPublicID string, skills []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Error beginning transaction: %v", err)
		return err
	}

	// Loop over the skills array
	for _, skillName := range skills {
		// Get the skill ID
		var skillID int
		query := `
			SELECT id FROM skills WHERE name = $1
		`
		err := tx.QueryRow(ctx, query, skillName).Scan(&skillID)
		if err != nil {
			if errors.Is(pgx.ErrNoRows, err) {
				continue
			}
			r.logger.Errorf("Error retrieving skill ID: %v", err)
			tx.Rollback(ctx)
			return err
		}

		// Delete the skill from the position
		deleteQuery := `
			DELETE FROM position_skills
			WHERE position_id = (SELECT id FROM positions WHERE public_id = $1)
			AND skill_id = $2
		`
		_, err = tx.Exec(ctx, deleteQuery, positionPublicID, skillID)
		if err != nil {
			r.logger.Errorf("Error deleting skill from position: %v", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func (r *positionRepository) GetPositionsByCompany(companyID string, pageNum int, pageSize int, search string) ([]models.Position, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()
	countQuery := `
		SELECT
			COUNT(*)
		FROM
			positions p
		INNER JOIN
			recruiters r ON p.recruiter_public_id = r.public_id
		WHERE
			r.company_public_id = $1
			AND (p.name ILIKE $2 OR p.description ILIKE $2)
	`

	query := `
		SELECT
			p.public_id,
			p.description,
			p.name,
			p.status,
			p.recruiter_public_id
		FROM
			positions p
		INNER JOIN
			recruiters r ON p.recruiter_public_id = r.public_id
		WHERE
			r.company_public_id = $1
			AND (p.name ILIKE $2 OR p.description ILIKE $2)
		GROUP BY 
			p.id, 
			p.public_id,
			p.description,
			p.name,
			p.status,
			p.recruiter_public_id
		ORDER BY
			p.id ASC
		LIMIT $3 OFFSET $4
	`

	// Calculate the offset based on the page number and page size
	offset := (pageNum - 1) * pageSize

	// Format the search query by adding wildcard characters
	searchQuery := "%" + search + "%"

	// Retrieve the count of positions
	var count int
	err := r.db.QueryRow(ctx, countQuery, companyID, searchQuery).Scan(&count)
	if err != nil {
		r.logger.Errorf("Error retrieving positions by company: %v", err)
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, query, companyID, searchQuery, pageSize, offset)
	if err != nil {
		r.logger.Errorf("Error retrieving positions by company: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	positions := []models.Position{}
	for rows.Next() {
		var position models.Position
		err := rows.Scan(
			&position.PublicID,
			&position.Description,
			&position.Name,
			&position.Status,
			&position.RecruiterPublicID,
		)
		if err != nil {
			r.logger.Errorf("Error retrieving positions by company: %v", err)
			return nil, 0, err
		}
		query := `SELECT array_agg(s.name) from skills AS s
		INNER JOIN position_skills ps ON ps.skill_id = s.id
		INNER JOIN positions p ON ps.position_id = p.id
		WHERE p.public_id = $1`
		if position.PublicID != nil {
			err = r.db.QueryRow(ctx, query, *position.PublicID).Scan(&position.Skills)
			if err != nil {
				r.logger.Errorf("Error retrieving positions by company: %v", err)
				return nil, 0, err
			}
		}

		positions = append(positions, position)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error retrieving positions by company: %v", err)
		return nil, 0, err
	}

	return positions, count, nil
}

func (r *positionRepository) GetPositionsByRecruiter(recruiterID string, pageNum int, pageSize int, search string) ([]models.Position, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()
	countQuery := `
		SELECT
			COUNT(*)
		FROM
			positions p
		INNER JOIN
			recruiters r ON p.recruiter_public_id = r.public_id
		WHERE
			r.public_id = $1
			AND (p.name ILIKE $2 OR p.description ILIKE $2)
	`

	query := `
		SELECT
			p.public_id,
			p.description,
			p.name,
			p.status,
			p.recruiter_public_id
		FROM
			positions p
		INNER JOIN
			recruiters r ON p.recruiter_public_id = r.public_id
		WHERE
			r.public_id = $1
			AND (p.name ILIKE $2 OR p.description ILIKE $2)
		GROUP BY 
			p.id, 
			p.public_id,
			p.description,
			p.name,
			p.status,
			p.recruiter_public_id
		ORDER BY
			p.id ASC
		LIMIT $3 OFFSET $4
	`

	// Calculate the offset based on the page number and page size
	offset := (pageNum - 1) * pageSize

	// Format the search query by adding wildcard characters
	searchQuery := "%" + search + "%"

	// Retrieve the count of positions
	var count int
	err := r.db.QueryRow(ctx, countQuery, recruiterID, searchQuery).Scan(&count)
	if err != nil {
		r.logger.Errorf("Error retrieving positions: %v", err)
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, query, recruiterID, searchQuery, pageSize, offset)
	if err != nil {
		r.logger.Errorf("Error retrieving positions: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	positions := []models.Position{}
	for rows.Next() {
		var position models.Position
		err := rows.Scan(
			&position.PublicID,
			&position.Description,
			&position.Name,
			&position.Status,
			&position.RecruiterPublicID,
		)
		if err != nil {
			r.logger.Errorf("Error retrieving positions: %v", err)
			return nil, 0, err
		}
		query := `SELECT array_agg(s.name) from skills AS s
		INNER JOIN position_skills ps ON ps.skill_id = s.id
		INNER JOIN positions p ON ps.position_id = p.id
		WHERE p.public_id = $1`
		if position.PublicID != nil {
			err = r.db.QueryRow(ctx, query, *position.PublicID).Scan(&position.Skills)
			if err != nil {
				r.logger.Errorf("Error retrieving positions: %v", err)
				return nil, 0, err
			}
		}

		positions = append(positions, position)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error retrieving positions: %v", err)
		return nil, 0, err
	}

	return positions, count, nil
}

func (r *positionRepository) AddQuestionsToPosition(positionPublicID string, questions []*models.Question) ([]*models.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Error occurred while starting transaction: %v", err)
		return nil, err
	}

	// Retrieve the position ID from the positions table based on the public ID
	var positionID int
	getPositionIDQuery := `SELECT id FROM positions WHERE public_id = $1`
	err = tx.QueryRow(ctx, getPositionIDQuery, positionPublicID).Scan(&positionID)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			tx.Rollback(ctx)
			return nil, models.ErrPositionNotFound
		}
		r.logger.Errorf("Error retrieving position ID: %v", err)
		tx.Rollback(ctx)
		return nil, err
	}

	for _, question := range questions {
		insertQuery := `
		INSERT INTO questions (name, position_public_id, position_id, read_duration, answer_duration)
		VALUES ($1, $2, $3, $4, $5) RETURNING public_id
		`
		err = tx.QueryRow(
			ctx,
			insertQuery,
			question.Name,
			positionPublicID,
			positionID,
			question.ReadDuration,
			question.AnswerDuration).Scan(&question.PublicID)
		if err != nil {
			r.logger.Errorf("Error adding question to position: %v", err)
			tx.Rollback(ctx)
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Errorf("Error occurred while committing transaction: %v", err)
		return nil, err
	}

	return questions, nil
}
