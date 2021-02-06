// Package for reading GeoJSON files
package gjson

import (
   "os"
   "log"
   "io/ioutil"
   "encoding/json"
)

type VZFeatureCollection struct {
   Type     string            `json:"type"`
   Features []VZFeature         `json:"features"`
}

type VZFeature struct {
   Type       string          `json:"type"`
   Geometry   VZGeometry      `json:"geometry"`
   Properties VZProperties    `json:"properties"`
}

type VZGeometry struct {
   Type        string         `json:"type"`
   Coordinates [][]float64    `json:"coordinates"`
}

type VZProperties struct {
   FQNID            string    `json:"FQNID"`
   CableType        string    `json:"cable_type"`
   CalculatedLength float64   `json:"calculated_length"`
   Created          string    `json:"created_date"`
   EditWorkOrderId  string    `json:"edit_work_order_id"`
   FiberCount       int64     `json:"fiber_count"`
   FileId           string    `json:"file_id"`
   MeasuredLength   float64   `json:"measured_length"`
   NFID             string    `json:"nfid"`
   ObjectId         int64     `json:"object_id"`
   OriginalFQNID    string    `json:"original_fqnid"`
   UDFs             []string  `json:"udfs"`
   WorkOrderName    string    `json:"work_order_name"`
}

// func ReadVZFile {{{
//
//
func ReadVZFile() VZFeatureCollection {
   file, err := os.Open("/Users/sabra/go/src/clustering/kml/vz3gis/data/vz.geojson")
   if err != nil {
      log.Fatal(err)
   }

   defer file.Close()
   f, _ := ioutil.ReadAll(file)

   var vz VZFeatureCollection

   json.Unmarshal(f, &vz)

   return vz
} // }}}
