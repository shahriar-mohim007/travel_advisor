package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"travel_advisor/domain"
	"travel_advisor/pkg/conn"

	"github.com/spf13/cobra"
)

const createDistrictsTable = `
		CREATE TABLE IF NOT EXISTS districts (
		id SERIAL PRIMARY KEY,
		division_id INT NOT NULL,
		name VARCHAR(100) NOT NULL UNIQUE,
		bn_name VARCHAR(100),
		lat DOUBLE PRECISION NOT NULL,
		long DOUBLE PRECISION NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
		);
`
const createUser = `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
`

var (
	districtsMigrationCmd = &cobra.Command{
		Use:   "districts-migration",
		Short: "Run districts migrations",
		Long:  `Run districts migrations`,
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("--------Database is connecting-------")
			err := conn.ConnectDefaultDB()
			if err != nil {
				fmt.Println("Failed to connect to database:", err)
				return
			}
		},
		Run: DistrictsMigrationsCmd,
	}
)

func init() {
	rootCmd.AddCommand(districtsMigrationCmd)
}

func DistrictsMigrationsCmd(cmd *cobra.Command, args []string) {
	db := conn.DefaultDB()
	conn.InitClient()
	client := conn.GetHTTClient()

	if res := db.Exec(createDistrictsTable); res.Error != nil {
		fmt.Println("Failed to create districts table:", res.Error)
		return
	}
	if res := db.Exec(createUser); res.Error != nil {
		fmt.Println("Failed to create users table:", res.Error)
		return
	}

	url := "https://raw.githubusercontent.com/strativ-dev/technical-screening-test/main/bd-districts.json"
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Failed to fetch districts JSON:", err)
		return
	}
	defer resp.Body.Close()

	var payload domain.DistrictResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		fmt.Println("JSON decode failed:", err)
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	stmt := `
		INSERT INTO districts (id, division_id, name, bn_name, lat, long)
		VALUES ($1, $2, $3, $4, $5, $6);
	   `

	for _, d := range payload.Districts {
		id, _ := strconv.Atoi(d.Id)
		divisionID, _ := strconv.Atoi(d.DivisionID)
		lat, _ := strconv.ParseFloat(d.Lat, 64)
		long, _ := strconv.ParseFloat(d.Long, 64)

		res := tx.Exec(stmt, id, divisionID, d.Name, d.BnName, lat, long)
		if res.Error != nil {
			fmt.Println("Failed to insert/update district:", d.Name, res.Error)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Println("Failed to commit transaction:", err)
		return
	}

	fmt.Println("Districts migration completed successfully")

}
