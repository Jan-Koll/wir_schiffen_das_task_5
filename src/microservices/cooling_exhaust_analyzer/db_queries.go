package main

import (
	"log"
)

func getAnalysisById(id int) (result job, err error) {
	row := db.QueryRow("SELECT * FROM cooling_exhaust_analyzer WHERE job_id = $1", id)
	err = row.Scan(&result.JobId, &result.ConfigurationId, &result.Success, &result.CoolingSystem, &result.ExhaustSystem, &result.CreatedAt, &result.LastModified, &result.Status)
	return
}

func createJob(input job) (result job, err error) {
	result = input
	row := db.QueryRow(
		`INSERT INTO cooling_exhaust_analyzer (configuration_id, status, cooling_system, exhaust_system)
		 VALUES ($1, $2, $3, $4)
		 RETURNING job_id, created_at, last_modified`,
		input.ConfigurationId,
		input.Status,
		input.CoolingSystem,
		input.ExhaustSystem,
	)
	err = row.Scan(&result.JobId, &result.CreatedAt, &result.LastModified)
	if err != nil {
		result = job{}
	}
	return
}

func setStatus(input job) (result job, err error) {
	result = input
	row := db.QueryRow(
		`UPDATE cooling_exhaust_analyzer SET 
				status = $2,
				last_modified = now()
			WHERE job_id = $1
			RETURNING last_modified;`,
		input.JobId,
		input.Status,
	)
	err = row.Scan(&result.LastModified)
	if err != nil {
		result = job{}
	}
	return
}
func setResult(input job) (result job, err error) {
	result = input
	row := db.QueryRow(
		`UPDATE cooling_exhaust_analyzer SET
				status = $2,
				success = $3,
				last_modified = now()
			WHERE job_id = $1
			RETURNING last_modified;`,
		input.JobId,
		input.Status,
		input.Success,
	)
	err = row.Scan(&result.LastModified)
	if err != nil {
		result = job{}
	}
	return
}

func deleteByConfigurationId(id int) (err error) {
	_, err = db.Exec("DELETE FROM cooling_exhaust_analyzer WHERE configuration_id = $1", id)
	return
}

func getAnalysesByConfigurationId(configuration_id int) (result []job, err error) {
	rows, err := db.Query("SELECT * FROM cooling_exhaust_analyzer WHERE configuration_id = $1 ORDER BY job_id ASC", configuration_id)
	if err != nil {
		log.Println("Failed to query database\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var job job
		err = rows.Scan(&job.JobId, &job.ConfigurationId, &job.Success, &job.CoolingSystem, &job.ExhaustSystem, &job.CreatedAt, &job.LastModified, &job.Status)
		if err != nil {
			log.Println("Failed to read row\n", err)
			return
		}
		result = append(result, job)
	}
	return
}
