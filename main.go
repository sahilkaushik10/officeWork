

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
    
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type leaveApplication struct {
	ID         int     `form:"id" json:"ID"`
	Name       string  `form:"name" json:"Name"`
	LeaveType  string  `form:"leaveType" json:"LeaveType"`
	LeaveFrom  string  `form:"fromDate" json:"LeaveFrom"`
	LeaveTo    string  `form:"toDate" json:"LeaveTo"`
	Team       string  `form:"team" json:"Team"`
	File       *[]byte `form:"file" json:"File"`
	Reporter   string  `form:"reporter" json:"Reporter"`
}

type notification struct {
	ReportingManager string `json:"ReportingManager"`
	LeaveID          int    `json:"LeaveID"`
	Approved         bool   `json:"Approved"`
}

var leave_applications = []leaveApplication{
	{ID: 1, Name: "Sahil Kaushik", LeaveType: "Casual Leave", LeaveFrom: "2023-05-26", LeaveTo: "2023-05-27", Team: "AWS", Reporter: "Avinashi Sharma"},
	{ID: 2, Name: "Km Saloni", LeaveType: "Sick Leave", LeaveFrom: "2023-05-28", LeaveTo: "2023-05-29", Team: "AZURE", Reporter: "Km Saloni"},
	{ID: 3, Name: "Avinashi Sharma", LeaveType: "Earned Leave", LeaveFrom: "2023-05-30", LeaveTo: "2023-05-31", Team: "AZURE", Reporter: "Pradeep Kumar"},
	{ID: 4, Name: "Pradeep Kumar Bharti", LeaveType: "Casual Leave", LeaveFrom: "2023-06-01", LeaveTo: "2023-06-02", Team: "AZURE", Reporter: "Avinashi Sharma"},
}

var db *sql.DB

func initDB() {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

    var err error

    db, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }


    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to the database.....")

	sqlFile, err := os.ReadFile("ddl.sql")

    if err != nil {

    log.Fatal(err)

}
    _, err = db.Exec(string(sqlFile))

    if err != nil {

        log.Fatal(err)

    }

	fmt.Println("DDL statements executed successfully")
}



func getLeaveApplications(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, leave_applications)
    fmt.Println("GET API is working.....")

}

func postLeaveApplications(c *gin.Context) {
    
    var newLeaveApp leaveApplication

    if err := c.Bind(&newLeaveApp); err != nil {
        return
    }

    newLeaveApp.ID = len(leave_applications) + 1

    leave_applications = append(leave_applications, newLeaveApp)

    
    _, err := db.Exec("INSERT INTO my_schema.leave_table (name, leave_type, leave_from, leave_to, team, file, reporter) VALUES ($1, $2, $3, $4, $5, $6, $7)",
        newLeaveApp.Name, newLeaveApp.LeaveType, newLeaveApp.LeaveFrom, newLeaveApp.LeaveTo, newLeaveApp.Team, newLeaveApp.File, newLeaveApp.Reporter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusCreated, newLeaveApp)

    c.JSON(http.StatusOK, gin.H{"message": "Leave application submitted successfully"})
}


func getLeaveApplicationByID(c *gin.Context) {
	id := c.Param("id")

	var leaveApp leaveApplication
	row := db.QueryRow("SELECT id, name, leave_type, leave_from, leave_to, team, file, reporter FROM my_schema.leave_table WHERE id = $1", id)

	err := row.Scan(&leaveApp.ID, &leaveApp.Name, &leaveApp.LeaveType, &leaveApp.LeaveFrom, &leaveApp.LeaveTo, &leaveApp.Team, &leaveApp.File, &leaveApp.Reporter)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave application not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, leaveApp)
}


func getSavedLeaveApplications(c *gin.Context) {
    rows, err := db.Query("SELECT id, name, leave_type, leave_from, leave_to, team ,file, reporter FROM my_schema.leave_table")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var savedLeaveApps []leaveApplication

    for rows.Next() {
        var leaveApp leaveApplication
        err := rows.Scan(&leaveApp.ID, &leaveApp.Name, &leaveApp.LeaveType, &leaveApp.LeaveFrom, &leaveApp.LeaveTo, &leaveApp.Team, &leaveApp.File, &leaveApp.Reporter)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        savedLeaveApps = append(savedLeaveApps, leaveApp)
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.IndentedJSON(http.StatusOK, savedLeaveApps)
    fmt.Println("Data fetched from database successfully......")

}



func getNotificationByID(c *gin.Context) {
	id := c.Param("id")

	var notification notification

	row := db.QueryRow("SELECT leave_id, reporting_manager, approved FROM my_schema.notifications WHERE leave_id = $1", id)

	err := row.Scan(&notification.LeaveID, &notification.ReportingManager, &notification.Approved)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, notification)
}


// ........................................................................................................

func get_top_5_employees_leave_2023(c *gin.Context) {
    type leave_3 struct {

        Name           string   `form:"name" json:"name"`
        TotalLeaveDays     string   `form:"total_leave_days" json:"total_leave_days"`
		// Rank			string		`form:"rank" json:"rank"`
    }

    rows, err := db.Query("SELECT name, total_leave_days FROM my_schema.top_5_employees_leave_2023")

    if err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return

    }

    defer rows.Close()

    var savedLeaveApps []leave_3

    for rows.Next() {

        var leaveApp leave_3

        err := rows.Scan(&leaveApp.Name, &leaveApp.TotalLeaveDays)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        savedLeaveApps = append(savedLeaveApps, leaveApp)

    }

    if err := rows.Err(); err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

        return

    }
    c.IndentedJSON(http.StatusOK, savedLeaveApps)
    fmt.Println("Data loaded successfully for the KPI")

}


func get_leave_counts_by_manager_2023(c *gin.Context) {
    type leave_4 struct {

        Reporter           string   `form:"Name" json:"reporter"`
        EmployeeCount     string   `form:"EmployeeCount" json:"employeeCount"`
    }

    rows, err := db.Query("SELECT reporter, employee_count FROM my_schema.leave_counts_by_manager_2023")

    if err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return

    }

    defer rows.Close()

    var savedLeaveApps []leave_4

    for rows.Next() {

        var leaveApp leave_4

        err := rows.Scan(&leaveApp.Reporter, &leaveApp.EmployeeCount)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        savedLeaveApps = append(savedLeaveApps, leaveApp)

    }

    if err := rows.Err(); err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

        return

    }
    c.IndentedJSON(http.StatusOK, savedLeaveApps)
    fmt.Println("Data loaded successfully for the KPI")
}

func get_leave_type_distribution_top_2_teams_2022(c *gin.Context) {
    type leave_6 struct {
        Team        string `form:"team" json:"team"`
        LeaveType   string `form:"leave_type" json:"leave_type"`
        TotalLeaveDays  string `form:"total_leave_days" json:"total_leave_days"`
    }
    

    rows, err := db.Query("SELECT team, leave_type, total_leave_days FROM my_schema.leave_type_distribution_top_2_teams WHERE team='IT'")
    // row1, err := db.Query("SELECT team, leave_type, total_leave_days FROM my_schema.leave_type_distribution_top_2_teams WHERE team=''")

    if err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return

    }

    defer rows.Close()

    var savedLeaveApps []leave_6

    for rows.Next() {

        var leaveApp leave_6

        err := rows.Scan(&leaveApp.Team, &leaveApp.LeaveType, &leaveApp.TotalLeaveDays)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        savedLeaveApps = append(savedLeaveApps, leaveApp)

    }

    if err := rows.Err(); err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

        return
    }
    c.IndentedJSON(http.StatusOK, savedLeaveApps)
    fmt.Println("Data loaded successfully for the KPI")
}

func get_leave_type_distribution_top_2_teams_2022_2(c *gin.Context) {
    type leave_6 struct {
        Team        string `form:"team" json:"team"`
        LeaveType   string `form:"leave_type" json:"leave_type"`
        TotalLeaveDays  string `form:"total_leave_days" json:"total_leave_days"`
    }
    

    rows, err := db.Query("SELECT team, leave_type, total_leave_days FROM my_schema.leave_type_distribution_top_2_teams WHERE team='Data engineering'")

    if err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return

    }

    defer rows.Close()

    var savedLeaveApps []leave_6

    for rows.Next() {

        var leaveApp leave_6

        err := rows.Scan(&leaveApp.Team, &leaveApp.LeaveType, &leaveApp.TotalLeaveDays)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        savedLeaveApps = append(savedLeaveApps, leaveApp)

    }

    if err := rows.Err(); err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

        return

    }
    c.IndentedJSON(http.StatusOK, savedLeaveApps)
    fmt.Println("Data loaded successfully for the KPI")
}



func main() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "sahil@22")
	os.Setenv("DB_NAME", "my_database")

	initDB()
	defer db.Close()

	leaveRouter := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	leaveRouter.Use(cors.New(config))

	leaveRouter.GET("/leave-applications", getLeaveApplications)
	leaveRouter.POST("/leave-applications", postLeaveApplications)
	leaveRouter.GET("/saved-leave-applications", getSavedLeaveApplications)
	leaveRouter.GET("/leave-applications/:id", getLeaveApplicationByID)
	leaveRouter.GET("/notifications/:id", getNotificationByID)

	// ........views...........
    leaveRouter.GET("/view3", get_top_5_employees_leave_2023)
    leaveRouter.GET("/view4", get_leave_counts_by_manager_2023)
    leaveRouter.GET("/view6", get_leave_type_distribution_top_2_teams_2022)
    leaveRouter.GET("/view6Two", get_leave_type_distribution_top_2_teams_2022_2)
    

	leaveRouter.Run(":8000")
}





