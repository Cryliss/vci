// Package for reading GeoJSON files
package gjson

import (
   "os"
   "log"
   "io/ioutil"
   "encoding/json"
)

type CoreFeatureCollection struct {
   Type     string            `json:"type"`
   Features []CoreFeature     `json:"features"`
}

type CoreFeature struct {
   Type       string          `json:"type"`
   Geometry   CoreGeometry    `json:"geometry"`
   Properties CoreProperties  `json:"properties"`
}

type CoreGeometry struct {
   Type        string         `json:"type"`
   Coordinates [][]float64    `json:"coordinates"`
}

type CoreProperties struct {
   Addres         string   `json:"address"`
   City           string   `json:"city"`
   ConstNeeded    string   `json:"const_not_needed"`
   Created        string   `json:"created_date"`
   FileId         string   `json:"file_id"`
   FQNID          string   `json:"fqnids"`
   Hub            string   `json:"hub"`
   JobStatus      string   `json:"job_status"`
   NetworkClass   string   `json:"network_class"`
   NetworkType    string   `json:"network_type"`
   NFID           string   `json:"nfid"`
   Note           string   `json:"note"`
   Number         string   `json:"dyea_vz"`
   ObjectId       int64    `json:"object_id"`
   PermitUrl      string    `json:"permit_url"`
   Status         string    `json:"status"`
   Type           string    `json:"type"`
}
// func ReadCoreFile
func ReadCoreFile() CoreFeatureCollection {
   file, err := os.Open("/Users/sabra/go/src/clustering/kml/corepermits/data/core.geojson")
   if err != nil {
      log.Fatal(err)
   }

   defer file.Close()
   f, _ := ioutil.ReadAll(file)

   var core CoreFeatureCollection

   json.Unmarshal(f, &core)

   return core
} // }}}

// func ReadNedFile {{{
func ReadNedFile() CoreFeatureCollection {
   file, err := os.Open("/Users/sabra/go/src/clustering/kml/nedpermits/data/ned.geojson")
   if err != nil {
      log.Fatal(err)
   }

   defer file.Close()
   f, _ := ioutil.ReadAll(file)

   var core CoreFeatureCollection

   json.Unmarshal(f, &core)

   return core
} // }}}
