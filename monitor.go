package main

import (
  "fmt"
  "strings"
  "bufio"
  "log"
  "os"
  "strconv"
  "time"
  "math"
)

const Total_Fixtures = 70
const Max_Displayed = 65
const FIXTURES_LIST = "/Users/aursdasp/Documents/Go/network/fixtures4.csv"
const OPERATORS_LIST = "/Users/aursdasp/Documents/Go/network/operators.csv"
const TEST_HISTORY = "/Users/aursdasp/Documents/Go/network/test_history.csv"
const OUTPUT_FILE = "/Users/aursdasp/Documents/Go/network/error.txt"

func printTable(records [Total_Fixtures]UUT) {

  //**********************************************************************************************************************************************
  //----------------------------------------------------------------------------------------------------------------------
  //Index  StationID   FixtureID   Serial              Status                Timestamp         Hours since update
  fmt.Println("***********************************************************************************************************************************")
  fmt.Println("Index\t StationID\t FixtureID\t Serial           \t Status              \t Timestamp     \t\t Hours since update")
  fmt.Println("***********************************************************************************************************************************")
  for i := range records {
    //output each line of the table only if the serial number is not null
    if records[i].serial != "" && i < Max_Displayed {
      record_time := time.Date(records[i].year+2000, time.Month(records[i].month), records[i].day, records[i].hour, records[i].minute, 0, 0, time.FixedZone("Local", 0))
      //record_time := time.Date(records[i].year+2000, records[i].month, records[i].day, records[i].hour, records[i].minute, 0, 0, time.FixedZone("Local", 0))
      
      duration := time.Since(record_time)
      hours_since_launch := duration.Hours()
      hours_since_launch = Round(hours_since_launch, .5, 1)
      fmt.Println(strconv.Itoa(i+1)+".)", "\t", records[i].station_ID, "\t", records[i].fixture_ID[len(records[i].fixture_ID)-4:len(records[i].fixture_ID)], "\t\t", records[i].serial, "\t", records[i].status, "\t", strconv.Itoa(records[i].month)+"/"+strconv.Itoa(records[i].day)+"/"+strconv.Itoa(records[i].year)+" "+records[i].string_hour+":"+records[i].string_minute+"\t\t", hours_since_launch)
    } else if i < Max_Displayed {
      fmt.Println(strconv.Itoa(i+1)+".)", "\t", records[i].station_ID, "\t", records[i].fixture_ID[len(records[i].fixture_ID)-4:len(records[i].fixture_ID)], "\t\t", "No data available")
    }
    if (i == 10 || i == 21 || i == 32 || i == 43 || i == 54 || i == 59 || i == 64) && i < Max_Displayed - 1 {
      fmt.Println("----------------------------------------------------------------------------------------------------------------------")
    }
  }
  fmt.Println("***********************************************************************************************************************************")
}

func Round(val float64, roundOn float64, places int ) (newVal float64) {
  var round float64
  pow := math.Pow(10, float64(places))
  digit := pow * val
  _, div := math.Modf(digit)
  if div >= roundOn {
    round = math.Ceil(digit)
  } else {
    round = math.Floor(digit)
  }
  newVal = round / pow
  return
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
  file, err := os.Create(path)
  if err != nil {
    return err
  }
  defer file.Close()

  w := bufio.NewWriter(file)
  for _, line := range lines {
    fmt.Fprintln(w, line)
  }
  return w.Flush()
}

type UUT struct {
	serial string
	fixture_ID string
	station_ID string
	timestamp string
	user string
	quad_E string
	status string
	entry int
  month int
  day int
  year int
  hour int
  string_hour string
  minute int
  string_minute string
}

func main() {

	//this_serial := "F5K60110000FWC31J"
	units_under_test := [Total_Fixtures]UUT{}

	//pull all lines from the FIXTURES LIST file and store them in var fixture_lines
	fixture_lines, err1 := readLines(FIXTURES_LIST)
  	if err1 != nil {
    	log.Fatalf("readLines: %s", err1)
  	}
  	//step through each line of the lines pulled from the FIXTURES LIST file
  	for i, line := range fixture_lines {
    	//fmt.Println(i, line)
    	elements := strings.Split(line, ",")
    	units_under_test[i].station_ID = elements[0]
    	units_under_test[i].fixture_ID = elements[1]
  	}

    //step through each of the fixtures and find out the serial and EEEE code
  	for i := 0; i < Total_Fixtures ; i++ {
  		history_lines, err2 := readLines(TEST_HISTORY)

  		if err2 != nil {
    		log.Fatalf("readLines: %s", err2)
  		}
      
      //step through each line of the test history file where j is index and line is string of each line
  		for j, line := range history_lines {
    		elements := strings.Split(line, ",")
    		current_fixture := elements[0]
    		current_serial := elements[1]
    		current_length := len(current_serial)
        current_status := elements[2]

        //normalize status length to 20
        status_length := len(current_status)
        for k := status_length; k < 20; k++ {
          current_status=current_status+" "
        } 
    		
        //collect timestamp
        current_timestamp := elements[3]
        timestamp_elements := strings.Split(current_timestamp, " ")
        date_elements := strings.Split(timestamp_elements[0], "/")
        time_elements := strings.Split(timestamp_elements[1], ":")
        current_month, _ := strconv.Atoi(date_elements[0])
        current_day, _ := strconv.Atoi(date_elements[1])
        timestamp_size := len(date_elements)
        current_year := 16
        if timestamp_size > 2 {
          current_year, _ = strconv.Atoi(date_elements[2])
        }
        
        current_hour, _ := strconv.Atoi(time_elements[0])
        var current_string_hour string
        if current_hour < 10 {
          current_string_hour = " "+strconv.Itoa(current_hour)
        } else {
          current_string_hour = time_elements[0]
        }
        current_string_minute := time_elements[1]
        current_minute, _ := strconv.Atoi(current_string_minute)

    		current_user := elements[4]
        if units_under_test[i].fixture_ID == current_fixture && current_serial != "" {
    			units_under_test[i].serial = current_serial
          //extract quad-E if the serial number is 17 characters, otherwise, make it null string
          if current_length == 17 {
            units_under_test[i].quad_E = units_under_test[i].serial[11:15]
          } else {
            units_under_test[i].quad_E = ""
          }
          units_under_test[i].timestamp = current_timestamp
          units_under_test[i].string_minute = current_string_minute
          units_under_test[i].minute = current_minute
          units_under_test[i].string_hour = current_string_hour
          units_under_test[i].hour = current_hour
          units_under_test[i].day = current_day
          units_under_test[i].month = current_month
          units_under_test[i].year = current_year

    			units_under_test[i].user = current_user
    			units_under_test[i].status = current_status
    			units_under_test[i].entry = j
        }
  		}	
    }
  	printTable(units_under_test)
}