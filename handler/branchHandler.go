package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"preview/config"
	"preview/entity"
	"preview/helpers"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func GetBranch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to connect to the database",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to connect to the database: %v", err)
		return
	}
	defer db.Close()

	ctx := context.Background()

	query := `
		SELECT branch_id, name, location FROM branches
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to fetch branches",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to fetch branches: %v", err)
		return
	}
	defer rows.Close()

	var branch []entity.Branch
	for rows.Next() {
		b := entity.Branch{}
		err := rows.Scan(&b.ID, &b.Name, &b.Location)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		branch = append(branch, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branch)
}

func GetBranchByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to connect to the database",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to connect to the database: %v", err)
		return
	}
	defer db.Close()

	ctx := context.Background()
	var branch entity.Branch

	id := p.ByName("id")

	branchID, err := strconv.Atoi(id)
	if err != nil {
		errResponse := helpers.ErrorResponse{

			Status:     "400",
			Title:      "Bad Request",
			Detail:     "Invalid Branches ID",
			StatusCode: http.StatusBadRequest,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusBadRequest)
		return
	}

	query := `
		SELECT branch_id, name, location FROM branches
		WHERE branch_id = ?
	`

	row := db.QueryRowContext(ctx, query, branchID)
	err = row.Scan(&branch.ID, &branch.Name, &branch.Location)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errResponse := helpers.ErrorResponse{
				Status:     "404",
				Title:      "Not Found",
				Detail:     "Branches Not Found",
				StatusCode: http.StatusNotFound,
			}

			helpers.SendErrorResponse(w, errResponse, http.StatusNotFound)
			return
		}

		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to fetch branches details",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to fetch branches details: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branch)
}

func CreateBranch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to connect to the database",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to connect to the database: %v", err)
		return
	}
	defer db.Close()

	ctx := context.Background()
	var branch entity.Branch

	err = json.NewDecoder(r.Body).Decode(&branch)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "400",
			Title:      "Bad Request",
			Detail:     "Invalid request body",
			StatusCode: http.StatusBadRequest,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusBadRequest)
		return
	}

	if branch.Name == "" || branch.Location == "" {
		errResponse := helpers.ErrorResponse{
			Status:     "400",
			Title:      "Bad Request",
			Detail:     "Name and Location are required fields",
			StatusCode: http.StatusBadRequest,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO branches (name, location)
		VALUES (?, ?)
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to prepare SQL statement",
			StatusCode: http.StatusInternalServerError,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to prepare SQL statement: %v", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, branch.Name, branch.Location)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to insert into the database",
			StatusCode: http.StatusInternalServerError,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to insert into the database: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to create branches",
			StatusCode: http.StatusInternalServerError,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Println("Error creating branches: no rows affected")
		return
	}

	id, _ := result.LastInsertId()
	branch.ID = int(id)

	successResponse := helpers.SuccessResponse{
		Status:     "201",
		Title:      "Success",
		Detail:     "Branches Successfully Created",
		StatusCode: http.StatusCreated,
	}

	helpers.SendSuccessResponse(w, successResponse, http.StatusCreated)
}

func DeleteBranchByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to connect to the database",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to connect to the database: %v", err)
		return
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	branchID, err := strconv.Atoi(id)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "502",
			Title:      "Bad Gateway",
			Detail:     "Invalid Branches ID",
			StatusCode: http.StatusBadGateway,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusBadGateway)
		return
	}

	var existingBranchID int
	checkQuery := "SELECT branch_id FROM branches WHERE branch_id = ?"
	err = db.QueryRowContext(ctx, checkQuery, branchID).Scan(&existingBranchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errRespone := helpers.ErrorResponse{
				Status:     "404",
				Title:      "Not Found",
				Detail:     "Branch Not Found",
				StatusCode: http.StatusNotFound,
			}

			helpers.SendErrorResponse(w, errRespone, http.StatusNotFound)
			return
		}

		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed To Check Branches Existence",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to check branches existence: %v", err)
		return
	}

	deleteQuery := "DELETE FROM branches WHERE branch_id = ?"
	stmt, err := db.PrepareContext(ctx, deleteQuery)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to prepare SQL statement",
			StatusCode: http.StatusInternalServerError,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to prepare SQL statement: %v", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, branchID)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed To Delete Branch",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to delete branch: %v", err)
		return
	}

	successResponse := helpers.SuccessResponse{
		Status:     "200",
		Title:      "Success",
		Detail:     "Branch Deleted Successfully",
		StatusCode: http.StatusOK,
	}

	helpers.SendSuccessResponse(w, successResponse, http.StatusOK)
}

func UpdateBranchByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	db, err := config.GetDB()
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to connect to the database",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to connect to the database: %v", err)
		return
	}
	defer db.Close()

	ctx := context.Background()

	id := p.ByName("id")
	branchID, err := strconv.Atoi(id)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "502",
			Title:      "Bad Gateway",
			Detail:     "Invalid Branch ID",
			StatusCode: http.StatusBadGateway,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusBadGateway)
		return
	}

	var existingBranchID int
	checkQuery := "SELECT branch_id FROM branches WHERE branch_id = ?"
	err = db.QueryRowContext(ctx, checkQuery, branchID).Scan(&existingBranchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errRespone := helpers.ErrorResponse{
				Status:     "404",
				Title:      "Not Found",
				Detail:     "Branch Not Found",
				StatusCode: http.StatusNotFound,
			}

			helpers.SendErrorResponse(w, errRespone, http.StatusNotFound)
			return
		}

		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed To Check Branch Existence",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to check branch existence: %v", err)
		return
	}

	var updatedBranch entity.Branch
	err = json.NewDecoder(r.Body).Decode(&updatedBranch)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "400",
			Title:      "Bad Request",
			Detail:     "Invalid request body",
			StatusCode: http.StatusBadRequest,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusBadRequest)
		return
	}

	if updatedBranch.Name == "" || updatedBranch.Location == "" {
		errResponse := helpers.ErrorResponse{
			Status:     "400",
			Title:      "Bad Request",
			Detail:     "Name and Location are required fields",
			StatusCode: http.StatusBadRequest,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusBadRequest)
		return
	}

	updateQuery := "UPDATE branches SET name = ?, location = ? WHERE branch_id = ?"
	stmt, err := db.PrepareContext(ctx, updateQuery)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed to prepare SQL statement",
			StatusCode: http.StatusInternalServerError,
		}
		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to prepare SQL statement: %v", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, updatedBranch.Name, updatedBranch.Location, branchID)
	if err != nil {
		errResponse := helpers.ErrorResponse{
			Status:     "500",
			Title:      "Internal Server Error",
			Detail:     "Failed To Update Branch",
			StatusCode: http.StatusInternalServerError,
		}

		helpers.SendErrorResponse(w, errResponse, http.StatusInternalServerError)
		log.Printf("Failed to update branch: %v", err)
		return
	}

	successResponse := helpers.SuccessResponse{
		Status:     "200",
		Title:      "Success",
		Detail:     "Branch Updated Successfully",
		StatusCode: http.StatusOK,
	}

	helpers.SendSuccessResponse(w, successResponse, http.StatusOK)
}
