package controller

import (
	"fmt"
	// "upload-excel-backend/model"
	// "time"
	// "fmt"
	// "io/ioutil"
	"net/http"
	// "path/filepath"
	// "os"
	"log"
	"strings"
	"strconv"
	// "encoding/json"
	"errors"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIController struct {
	db *gorm.DB
}

func NewKPIController(db *gorm.DB) *KPIController {
	return &KPIController{
		db: db,
	}
}
// Create a struct to represent the JSON format
type KpiFileJson struct {
	NameId             string `json:"nameId"`
	Period             string `json:"period"`
	ObjectiveType      string `json:"objType"`
	KRA                string `json:"kra"`
	Description        string `json:"desc"`
	IndividualCriteria int    `json:"individualCriteria"`
	Mark1Desc          string `json:"mark1Desc"`
	Mark2Desc          string `json:"mark2Desc"`
	Mark3Desc          string `json:"mark3Desc"`
	Mark4Desc          string `json:"mark4Desc"`
}

func (c *KPIController) PostFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Open the uploaded file
	excelFile, err := file.Open()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to open file"})
		return
	}
	defer excelFile.Close()

	// Load the Excel file using excelize
	xlsx, err := excelize.OpenReader(excelFile)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to read Excel file"})
		return
	}

	// Get the values from the first sheet
	sheetName := xlsx.GetSheetName(1)
	rows := xlsx.GetRows(sheetName)

	// Create an array to store the extracted data
	var data []string

	// Flags to track the presence of "Criteria" and "Performance Grading"
	criteriaFound,weightageFound, individualFound, totalFound := false, false, false, false

	gradingFound := false
	grading1, grading2, grading3, grading4 := false, false, false, false

	kraFound := false
	taskFound := false

	// Iterate through the rows and store the values in the array
	for rowIndex, row := range rows {
		// Check if the row is empty (all cells in the row are empty)
		if isEmptyRow(row) {
			continue // Skip processing the current row and move to the next iteration
		}
		// Check if the row is within the desired range
		if rowIndex >= 0 && rowIndex <= 2 {
			for colIndex, cell := range row {
				// Check if the column is within the desired range (A to I)
				if colIndex >= 0 && colIndex <= 9 {
					data = append(data, cell)
				}
				if colIndex == 2 && (rowIndex == 0 || rowIndex == 1) {
					if cell == "Criteria" {
						criteriaFound = true
					}
				}
				if colIndex == 2 && rowIndex == 2 {
					if cell == "TASK" {
						taskFound = true
					}
				}
				if colIndex == 1 && rowIndex == 2 {
					if cell == "KRA" {
						kraFound = true
					}
				}
				if (colIndex == 3 || colIndex == 4) && (rowIndex == 0 || rowIndex == 1) {
					if cell == "Weightage (%)" {
						weightageFound = true
					} else if cell == "Individual Criteria" {
						individualFound = true
					} else if cell == "Total" {
						totalFound = true
					}
				}
				if (colIndex == 5 || colIndex == 6 || colIndex == 7 || colIndex == 8) && (rowIndex == 0 || rowIndex == 1) {
					if cell == "Performance Grading" {
						gradingFound = true
					} else if cell == "4" {
						grading4 = true
					} else if cell == "3" {
						grading3 = true
					} else if cell == "2" {
						grading2 = true
					} else if cell == "1" {
						grading1 = true
					}
				}
			}
		}
	}

	headerComplete := false

	if kraFound && taskFound && criteriaFound && weightageFound && gradingFound && grading1 && grading2 && grading3 && grading4 && individualFound && totalFound {
		headerComplete = true
		log.Println("All found")
	} else {
		log.Println("Not all found")
	}

	if kraFound {
		log.Println("FOUND IT")
	} else {
		log.Println("NOT FOUND")
	}

	var dataTask [][]string
	var dataBehavior [][]string
	var dataOrg [][]string
	var dataIndividual [][]string
	var dataAll [][]string

	// If Header is correct, then read the rest of the data
	if headerComplete {
		log.Println("Header of the file is Complete. Continue reading data...")

		var readA bool
		var readB bool
		var readC bool
		var readD bool

		// Iterate through the rows and store the values in the array
		for rowIndex, row := range rows {
			if rowIndex >= 2 {
				if len(row) > 0 && row[0] == "A" {
					// Skip reading the current row if a cell in column A is equal to "A"
					readA = true
					continue
				}
				if readA && !readB {
					// Create an array to store the extracted data for the current row
					var rowData []string

					for _, cell := range row {
						// Check if the cell is not empty
						if cell != "" {
							// Check if the cell value is "B"
							if cell == "B" {
								// Mark the end of the "dataTask" section and start the "dataBehavior" section
								readA = false
								readB = true
								break
							}
							if cell == "#REF!" {
								ctx.JSON(400, gin.H{"message": "Upload fail"})
								return
							}
							rowData = append(rowData, cell)
						}
					}
					rowDataString := strings.Join(rowData, ", ")

					if rowDataString != "" {
						rowData = append(rowData, "Task")
						dataTask = append(dataTask, rowData)
						fmt.Printf("Extracted Data for Row %d (dataTask): ----- [ %s ] -----\n", rowIndex+1, dataTask)
					}
				} else if readB && !readC {
					// Create an array to store the extracted data for the current row
					var rowData []string

					for _, cell := range row {
						// Check if the cell is not empty
						if cell != "" {
							// Check if the cell value is "C"
							if cell == "C" {
								// Mark the end of the "dataBehavior" section and start the "dataOrg" section
								readB = false
								readC = true
								break
							}
							// Check if the cell value is "#REF!"
							if cell == "#REF!" {
								ctx.JSON(400, gin.H{"message": "Upload fail"})
								return
							}
							// Store the cell value in the dataBehavior array
							rowData = append(rowData, cell)
						}
					}
					rowDataString := strings.Join(rowData, ", ")

					if rowDataString != "" {
						rowData = append(rowData, "Behavior/Attitude")
						dataBehavior = append(dataBehavior, rowData)
						fmt.Printf("Extracted Data for Row %d (dataBehavior): ----- [ %s ] -----\n", rowIndex+1, dataBehavior)
					}
				} else if readC && !readD {
					// Create an array to store the extracted data for the current row
					var rowData []string

					for _, cell := range row {
						// Check if the cell is not empty
						if cell != "" {
							// Check if the cell value is "D"
							if cell == "D" {
								// Mark the end of the "dataOrg" section and start the "dataIndividual" section
								readC = false
								readD = true
								break
							}
							// Check if the cell value is "#REF!"
							if cell == "#REF!" {
								ctx.JSON(400, gin.H{"message": "Upload fail"})
								return
							}
							// Store the cell value in the dataOrg array
							rowData = append(rowData, cell)
						}
					}
					rowDataString := strings.Join(rowData, ", ")

					if rowDataString != "" {
						rowData = append(rowData, "Organizational Goal")
						dataOrg = append(dataOrg, rowData)
						fmt.Printf("Extracted Data for Row %d (dataOrg): ----- [ %s ] -----\n", rowIndex+1, dataOrg)
					}
				} else if readD {
					// Create an array to store the extracted data for the current row
					var rowData []string

					for _, cell := range row {
						// Check if the cell is not empty
						if cell != "" {
							// Check if the cell value is "#REF!"
							if cell == "#REF!" {
								ctx.JSON(400, gin.H{"message": "Upload fail"})
								return
							}
							// Store the cell value in the dataIndividual array
							rowData = append(rowData, cell)
						}
					}
					rowDataString := strings.Join(rowData, ", ")

					if rowDataString != "" {
						rowData = append(rowData, "Individual Goal")
						dataIndividual = append(dataIndividual, rowData)
						fmt.Printf("Extracted Data for Row %d (dataIndividual): ----- [ %s ] -----\n", rowIndex+1, dataIndividual)
					}
				}
			}
		}
		dataAll = append(dataAll, dataTask...)
		dataAll = append(dataAll, dataBehavior...)
		dataAll = append(dataAll, dataOrg...)
		dataAll = append(dataAll, dataIndividual...)
	} else {
		log.Println("Header of the file is NOT COMPLETE.")
	}

	jsonData := make([]map[string]interface{}, 0)

	for _, kpiValues := range dataAll {
		removePercent := strings.Replace(kpiValues[3], "%", "", 1)
		// Convert the number value to an integer
		numberValue, err := strconv.Atoi(removePercent)
		if err != nil {
			log.Println("Error converting number value:", err)
			continue
		}

		loginName := ctx.PostForm("loginName")
		period := ctx.PostForm("period")
    	if len(kpiValues) >= 10 {
        // Create an instance of the KpiFileJson struct
        	data := KpiFileJson{
				NameId:             loginName,
				Period:             period,
				ObjectiveType:      kpiValues[9],
				KRA:                kpiValues[1],
				Description:        kpiValues[2],
				IndividualCriteria: numberValue,
				Mark1Desc:          kpiValues[5],
				Mark2Desc:          kpiValues[6],
				Mark3Desc:          kpiValues[7],
				Mark4Desc:          kpiValues[8],
        	}

			// Convert data to a map of key-value pairs based on the KpiFileJson struct
            dataMap := map[string]interface{}{
                "nameId":             data.NameId,
                "period":             data.Period,
                "objType":            data.ObjectiveType,
                "kra":                data.KRA,
                "desc":               data.Description,
                "individualCriteria": data.IndividualCriteria,
                "mark1Desc":          data.Mark1Desc,
                "mark2Desc":          data.Mark2Desc,
                "mark3Desc":          data.Mark3Desc,
                "mark4Desc":          data.Mark4Desc,
            }

			// Save the data to the database using the SaveToDatabase function
			if err := c.SaveToDatabase(data); err != nil {
				log.Printf("Error saving data to database: %v", err)
				// return
			}
			// Append the dataMap to the jsonData slice
			jsonData = append(jsonData, dataMap)
		} else {
			// Log the data that caused the issue for debugging
			log.Printf("kpiValues does not have enough elements: %+v", kpiValues)
		}
	}
	fmt.Println("++++++Data from ALL DATA array:\n\n\n", dataAll)
	ctx.JSON(http.StatusOK, jsonData)
}

func (c *KPIController) SaveToDatabase(data KpiFileJson) error {
    // Ensure the database connection is open
    if c.db.Error != nil {
        log.Printf("Error connecting to the database: %v", c.db.Error)
        return c.db.Error
    }

    // Call the stored procedure using GORM's Raw method
    result := c.db.Exec("EXEC usp_CreateKpi ?, ?, ?, ?, ?, ?, ?, ?, ?, ?", 
        data.NameId, data.Period, data.ObjectiveType, data.KRA, data.Description, 
        data.IndividualCriteria, data.Mark1Desc, data.Mark2Desc, data.Mark3Desc, data.Mark4Desc)

    if result.Error != nil {
        // Handle the error
        log.Printf("Error calling stored procedure: %v", result.Error)
        return result.Error
    }

    // Check the number of rows affected by the execution
    rowsAffected := result.RowsAffected
    if rowsAffected == 0 {
        // If no rows were affected, return an error indicating that the data was not inserted
        return errors.New("data not inserted into database")
    }

    return nil
}

func isEmptyRow(row []string) bool {
    for _, cell := range row {
        if cell != "" {
            return false
        }
    }
    return true
}