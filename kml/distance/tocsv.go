// Sabra Bilodeau
package distance

import (
   "encoding/csv"
   "fmt"
   "os"
   "strconv"
   "vci/kml/vz3gis"
)

func MakeFQNIDCSVFile(fqnids []vz3gis.FQNID) {
   // Create a csv file
   f, err := os.Create("./fqnids.csv")
   if err != nil {
     fmt.Println(err)
   }
   defer f.Close()

   var record []string
   record = append(record, "FQNID")
   record = append(record, "Original FQNID")
   record = append(record, "NFID")
   record = append(record, "DYEAs")
   record = append(record, "Cable Type")
   record = append(record, "Fiber Count")
   record = append(record, "Calculated Length")
   record = append(record, "Measured Length")
   record = append(record, "udf_1")
   record = append(record, "udf_2")
   record = append(record, "udf_3")
   record = append(record, "udf_4")
   record = append(record, "udf_5")
   record = append(record, "udf_6")
   record = append(record, "udf_7")
   record = append(record, "udf_8")
   record = append(record, "Work Order Name")
   record = append(record, "Live Maps Object ID")
   record = append(record, "Live Maps Layer")
   record = append(record, "Coordinates")
   w := csv.NewWriter(f)
   w.Write(record)

   for _, fq := range fqnids {
      var record []string
      record = append(record, fq.FQNID)
      record = append(record, fq.OriginalFQNID)
      record = append(record, fq.NFID)

      var dys string
      for _, dy := range fq.DYEAs {
         dys = dys + dy + ", "
      }

      record = append(record, dys)
      record = append(record, fq.CableType)
      record = append(record, strconv.FormatInt(fq.FiberCount,10))
      record = append(record, strconv.FormatFloat(fq.CalculatedLength, 'f', 'f', 64))
      record = append(record, strconv.FormatFloat(fq.MeasuredLength, 'f', 'f', 64))

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
         coord := fmt.Sprintf("(%f,%f)", c.Lat, c.Lng)
         coords = coords + coord
      }
      record = append(record, coords)
      w.Write(record)
   }
   w.Flush()
}

func MakeMatchesCSVFile(matches []Match) {
    // Create a csv file
   f, err := os.Create("./matches.csv")
   if err != nil {
      fmt.Println(err)
   }
   defer f.Close()

   var record []string
   record = append(record, "PermitObjectId")
   record = append(record, "VZObjectId")
   record = append(record, "FQNID")
   record = append(record, "DYEA")
   record = append(record, "Type")
   record = append(record, "Status")
   record = append(record, "Note")
   record = append(record, "CoreLat")
   record = append(record, "CoreLng")
   record = append(record, "VZLat")
   record = append(record, "VZLng")
   record = append(record, "udf_1")
   record = append(record, "udf_2")
   record = append(record, "udf_3")
   record = append(record, "udf_4")
   record = append(record, "udf_5")
   record = append(record, "udf_6")
   record = append(record, "udf_7")
   record = append(record, "udf_8")
   w := csv.NewWriter(f)
   w.Write(record)

   for _, m := range matches {
      var record []string
      record = append(record, strconv.FormatInt(m.ObjectId, 10))
      record = append(record, strconv.FormatInt(m.VzObjectId, 10))
      record = append(record, m.FQNID)
      if len(m.DYEAs) >= 1 {
         if m.DYEAs[0].Number == "" {
            record = append(record, m.DYEAs[1].Number)
            record = append(record, m.DYEAs[1].Type)
            record = append(record, m.DYEAs[1].Status)
            record = append(record, m.DYEAs[1].Note)
         } else {
            record = append(record, m.DYEAs[0].Number)
            record = append(record, m.DYEAs[0].Type)
            record = append(record, m.DYEAs[0].Status)
            record = append(record, m.DYEAs[0].Note)
         }
      } else {
         record = append(record, "")
         record = append(record, "")
         record = append(record, "")
         record = append(record, "")
      }

      record = append(record, strconv.FormatFloat(m.Point.Lat, 'f', 'f', 64))
      record = append(record, strconv.FormatFloat(m.Point.Lng, 'f', 'f', 64))
      record = append(record, strconv.FormatFloat(m.VzPoint.Lat, 'f', 'f', 64))
      record = append(record, strconv.FormatFloat(m.VzPoint.Lng, 'f', 'f', 64))
      for _, udf := range m.UDFs {
         record = append(record, udf)
      }
      w.Write(record)
   }
   w.Flush()
}
