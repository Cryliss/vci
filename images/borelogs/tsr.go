// Functions for dealing with TSR Drill Reports in the borelogs folders
package borelogs
import (
   "encoding/csv"
   "strconv"
   "fmt"
)

type DrillReport struct {
   File string
   Pipes []Pipe
}

type Pipe struct {
   Number int64
   Lat float64
   Lng float64
}

func tsr(file string) {
   records, err := readData(file)
   var dr DrillReport
   dr.File = file

    if err != nil {
        log.Fatal(err)
    }

    for _, record := range records {
      if record[12] == "" {
         continue
      }

      pipe := Pipe {
         Number: strconv.ParseInt(record[14], 64),
         Lat: strconv.ParseFloat(record[12], 64),
         Lng: strconv.ParseFlost(record[13], 64),
      }

      dr.Pipes = append(dr.Pipes, pipe)
      fmt.Printf("%d : (%f, %f)\n", pipe.Number, pipe.lng, pipe.lng)

    }
}

func readData(fileName string) ([][]string, error) {

    f, err := os.Open(fileName)

    if err != nil {
        return [][]string{}, err
    }
    defer f.Close()

    r := csv.NewReader(f)

    // skip first line
    if _, err := r.Read(); err != nil {
        return [][]string{}, err
    }

    records, err := r.ReadAll()

    if err != nil {
        return [][]string{}, err
    }

    return records, nil
}
