package controller

import (
	"fmt"
	"upload-excel-backend/model"
	// "time"
	// "fmt"
	// "io/ioutil"
	// "net/http"
	// "path/filepath"
	// "os"
	"log"
	"strings"
	// "strconv"
	// "encoding/json"

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
type OrgJSON struct {
	Number     int    `json:"number"`
	KRA        string `json:"KRA"`
	Individual string `json:"individual"`
}

type Data struct {
	Number      	string `json:"number"`
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
	criteriaFound := false
	weightageFound := false
	individualFound := false
	totalFound := false

	gradingFound := false
	grading1 := false
	grading2 := false
	grading3 := false
	grading4 := false

	kraFound := false
	taskFound := false

	// Iterate through the rows and store the values in the array
	for rowIndex, row := range rows {
		// Check if the row is within the desired range
		if rowIndex >= 0 && rowIndex <= 2 {
			for colIndex, cell := range row {
				// Check if the column is within the desired range (A to I)
				if colIndex >= 0 && colIndex <= 8 {
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

	// Create arrays to store the extracted data
	var dataTask []string
	var dataBehavior []string
	var dataOrg []string
	var dataIndividual []string

	// If Header is correct, then read the rest of the data
	if headerComplete {
		log.Println("Header of the file is Complete. Continue reading data...")

		// Flags to track the start and end of data sections
		var readA bool
		var readB bool
		var readC bool
		var readD bool

		for rowIndex, row := range rows {
			if rowIndex >= 2  {
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

							// Check if the cell value is "#REF!"
							if cell == "#REF!" {
								ctx.JSON(400, gin.H{"message": "Upload fail"})
								return
							}

							// Store the cell value in the dataTask array
							rowData = append(rowData, cell)
						}
					}

					// Print the rowData for the current row
					rowDataString := strings.Join(rowData, ", ")

					if rowDataString != "" {
						fmt.Printf("Extracted Data for Row %d (dataTask): ----- [ %s ] -----\n", rowIndex+1, rowDataString)
						// Store the rowDataString in the dataTask array
						dataTask = append(dataTask, rowDataString)
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

					// Print the rowData for the current row
					rowDataString := strings.Join(rowData, ", ")

					if rowDataString != "" {
						fmt.Printf("Extracted Data for Row %d (dataBehavior): ----- [ %s ] -----\n", rowIndex+1, rowDataString)
						// Store the rowDataString in the dataBehavior array
						dataBehavior = append(dataBehavior, rowDataString)
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

					// Print the rowData for the current row
					rowDataString := strings.Join(rowData, ", ")

					if rowDataString != "" {
						fmt.Printf("Extracted Data for Row %d (dataOrg): ----- [ %s ] -----\n", rowIndex+1, rowDataString)
						// Store the rowDataString in the dataOrg array
						dataOrg = append(dataOrg, rowDataString)
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

					// Print the rowData for the current row
					rowDataString := strings.Join(rowData, ", ")

					if rowDataString != "" {
						fmt.Printf("Extracted Data for Row %d (dataIndividual): ----- [ %s ] -----\n", rowIndex+1, rowDataString)
						// Store the rowDataString in the dataIndividual array
						dataIndividual = append(dataIndividual, rowDataString)
					}
				}
			}
		}
	} else {
		ctx.JSON(400, gin.H{"message": "Upload fail. Make sure header of the file is in correct format."})
		log.Println("Header of the file is NOT COMPLETE.")
		return
	}
	// dataTaskString := strings.Join(dataTask, "|")

	// Print the extracted data from dataTask and dataBehavior arrays
	fmt.Println("Data from dataTask array**********:\n", dataTask)
	fmt.Println("Data from dataBehavior array:\n", dataBehavior)
	fmt.Println("Data from dataOrg array:\n", dataOrg)
	fmt.Println("Data from dataIndividual array:\n", dataIndividual)

	ctx.JSON(200, gin.H{"message": "File uploaded and processed successfully"})
	// Send the JSON response back to the client (Postman)
	// ctx.JSON(200, orgDataJSON)
}


func (c *KPIController) GetKPIs(ctx *gin.Context) {
	var kpis []model.KPI
	c.db.Find(&kpis)

	ctx.JSON(200, kpis)
}
