// Sabra Bilodeau
package vz3gis

import (
   "encoding/csv"
   "fmt"
   "os"
   "strconv"
)

func MakeCSVFile(fqnids [][]FQNID) {
   // Create a csv file
   f, err := os.Create("/Users/sabra/go/src/mci/kml_v2/vz3gis/data/fqnids.csv")
   if err != nil {
     fmt.Println(err)
   }
   defer f.Close()

   var record []string
   record = append(record, "FQNID")
   record = append(record, "Original FQNID")
   record = append(record, "NFID")
   record = append(record, "type")
   record = append(record, "Fiber Count")
   record = append(record, "Calculated Length")
   record = append(record, "Measured Length")
   record = append(record, "Edit Work Order ID")
   record = append(record, "udf_1")
   record = append(record, "udf_2")
   record = append(record, "udf_3")
   record = append(record, "udf_4")
   record = append(record, "udf_5")
   record = append(record, "udf_6")
   record = append(record, "udf_7")
   record = append(record, "udf_8")
   record = append(record, "Work Order Name")
   record = append(record, "object_id")
   record = append(record, "Live Maps Layer")
   record = append(record, "coordinates")
   w := csv.NewWriter(f)
   w.Write(record)


    for _, folder := range fqnids {
      for _, fq := range folder {
         var record []string
         record = append(record, fq.FQNID)
         record = append(record, fq.OriginalFQNID)
         record = append(record, fq.NFID)
         record = append(record, fq.CableType)
         record = append(record, strconv.FormatInt(fq.FiberCount,10))
         record = append(record, strconv.FormatFloat(fq.CalculatedLength, 'f', 'f', 64))
         record = append(record, strconv.FormatFloat(fq.MeasuredLength, 'f', 'f', 64))
         record = append(record, fq.EWO)
         i := 0
         for _, udf := range fq.UDFs {
            i = i + 1
            record = append(record, udf)
         }
         if i < 8 {
            for j := 8; i < j; i++ {
               record = append(record, "")
            }
         }

         record = append(record, fq.WorkOrderName)
         record = append(record, strconv.FormatInt(fq.ObjectId, 10))
         record = append(record, fq.FileID)

         var coords string
         for _, c := range fq.Coordinates {
            coord := fmt.Sprintf("%f,%f", c.Lat, c.Lng)
            coords = coords + coord
         }
         record = append(record, coords)
         w.Write(record)
      }
   }
   w.Flush()
}
