package main

import (
   "log"
   "time"
   "vci/kml/corepermits"
   "vci/kml/vz3gis"
   "vci/kml/nedpermits"
   "strings"
   "clustering/gjson"
   "encoding/csv"
   "fmt"
   "os"
)

// func main {{{

func main() {
   start := time.Now()

	_, err := vz3gis.ParseKML()
   if err != nil {
      log.Printf("Error parsing VZ3GIS.kml file!")
   }

   _, err = corepermits.ParseKML()
   if err != nil {
      log.Printf("Error parsing CORE_PERMITS.kml file!")
   }

   _, err = nedpermits.ParseKML()

   log.Print("Successfully parsed all KML files!\n")

   //log.Printf("Getting TCP permits\n")
   //GetTCPPermits()

   elapsed := time.Since(start)
   log.Printf("KML took %s", elapsed)
} // }}}

// func GetTCPPermits {{{
//
// Returns an array of preprocessed permits
func GetTCPPermits() {
   //core := gjson.ReadCoreFile()
   //PreProcess(core)

   ned := gjson.ReadNedFile()
   PreProcess(ned)
} // }}}

// func PreProcess {{{
func PreProcess(core gjson.CoreFeatureCollection) {
   // Create a csv file
   f, err := os.Create("/Users/sabra/go/src/clustering/tcp_jobs.csv")
   if err != nil {
     fmt.Println(err)
   }
   defer f.Close()

   var record []string
   record = append(record, "Hub")
   record = append(record, "DYEA/VZ")
   record = append(record, "NFID")
   record = append(record, "Type")
   w := csv.NewWriter(f)
   w.Write(record)

   for _, p := range core.Features {
      typ := GetType(p.Properties.Type)
      if typ == "tcp" {
         var record []string
         record = append(record, p.Properties.Hub)
         record = append(record, p.Properties.Number)
         record = append(record, p.Properties.NFID)
         record = append(record, p.Properties.Type)
         w.Write(record)
      }
   }
   w.Flush()
} // }}}

// func GetPermitType {{{
//
// Returns a normalized verison of the permit type
func GetType(permit string) string {
   ptype := strings.ToLower(permit)
   if ptype == "" {
      return ""
   } else if len(ptype) == 2 {
      return "ug"
   } else if ptype == "excavation/tcp" {
      return "tcp"
   } else if ptype == "exc/tcp" {
      return "tcp"
   }

   var ptyp string
   if len(ptype) > 8 {
      ptyp = ptype[6:9]
      if ptyp == "tcp" {
         return "tcp"
      }
   } else if len(ptype) == 7 {
      ptyp = ptype[4:7]
      if ptyp == "tcp" {
         return "tcp"
      }
   }

   ptyp = ptype[0:3]

   if ptyp == "jpa" {
      return "ae"
   } else if ptyp == "tcp" {
      return "tcp"
   } else {
      return "ug"
   }
} // }}}
