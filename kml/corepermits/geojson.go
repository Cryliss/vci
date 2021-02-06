// Sabra Bilodeau
package corepermits

import (
   "fmt"
   "strings"
	"io/ioutil"
   "encoding/json"
   "github.com/paulmach/go.geojson"
)

// func MakeGeoJSONFile {{{
func MakeGeoJSONFile(pms [][]Permit) {
   fc := geojson.NewFeatureCollection()
   fileName := "/Users/sabra/go/src/clustering/kml/corepermits/data/core.geojson"
   for _, folder := range pms {
      for _, p := range folder {
         var coords [][]float64
         for _, c := range p.Coordinates {
            var coord []float64
            lng := c.Lng
            coord = append(coord, lng)
            lat := c.Lat
            coord = append(coord, lat)

            coords = append(coords, coord)
         }
         if len(coords) < 3 {
            var coord []float64
            lng := p.Coordinates[len(coords)-1].Lng
            coord = append(coord, lng)

            lat := p.Coordinates[len(coords)-1].Lat
            coord = append(coord, lat)

            coords = append(coords, coord)
         }

         for _, d := range p.DYEAs {
            var dy DYEA
            if d.Number == "" {
               continue
            } else if d.Number == " " {
               continue
            } else if d.Number == "In Progress" {
               continue
            } else if d.Number == "EXPIRED: 7-MAY-20" {
               continue
            } else if d.Number == "VZ_LAN_00000203/DYEA_LSA_8107782" {
               // sep == sep[0] = VZ_LAN_00000203
               //        sep[1] = DYEA_LSA_8107782
               sep := strings.Split(d.Number, "/")

               dy.Number = sep[0]
               typ := strings.Split(d.Type, "(")
               dy.Type = typ[0]
               dy.Status = d.Status
               if len(typ) > 1 {
                  note := fmt.Sprintf("%s, %s", d.Note, typ[1])
                  dy.Note = note
               } else {
                  dy.Note = d.Note
               }
               MakeFeature(dy, coords, p, fc)

               dy.Number = sep[1]
               MakeFeature(dy, coords, p, fc)
               continue
            } else if d.Number == "DYEA_LSA_8108575/VZ_LAN_00003058" {
               sep := strings.Split(d.Number, "/")

               dy.Number = sep[0]
               typ := strings.Split(d.Type, "(")
               dy.Type = typ[0]
               dy.Status = d.Status

               if len(typ) > 1 {
                  note := fmt.Sprintf("%s, %s", d.Note, typ[1])
                  dy.Note = note
               } else {
                  dy.Note = d.Note
               }

               MakeFeature(dy, coords, p, fc)

               dy.Number = sep[1]
               MakeFeature(dy, coords, p, fc)
               continue
            }

            dy.Number = d.Number
            typ := strings.Split(d.Type, "(")
            dy.Type = typ[0]
            dy.Status = d.Status

            if len(typ) > 1 {
               note := fmt.Sprintf("%s, %s", d.Note, typ[1])
               dy.Note = note
            } else {
               dy.Note = d.Note
            }

            MakeFeature(dy, coords, p, fc)
         }
      }
   }

   file, _ := json.MarshalIndent(fc, "", "    ")
   _ = ioutil.WriteFile(fileName, file, 0644)
}// }}}

// func MakeFeature {{{
//
// Makes a gejson feature
func MakeFeature(dy DYEA, coords [][]float64, p Permit, fc *geojson.FeatureCollection) {
   /*typ := strings.ToLower(dy.Type)
   if len(typ) >= 3 {
      typ = typ[0:3]
   }

   if (typ != "exc" && typ != "jpa") {
      return
   } else if typ == "rr" {
      return
   }*/

   f := geojson.NewLineStringFeature(coords)
   f.SetProperty("object_id", p.ObjectId)
   f.SetProperty("file_id", p.FileId)
   f.SetProperty("created_date", p.Created)
   f.SetProperty("network_class", p.NetworkClass)
   f.SetProperty("network_type", p.NetworkType)
   f.SetProperty("permit_url", p.PermitUrl)
   f.SetProperty("dyea_vz", dy.Number)
   f.SetProperty("type", dy.Type)
   f.SetProperty("status", dy.Status)
   f.SetProperty("note", dy.Note)
   f.SetProperty("nfid", p.NFID)
   f.SetProperty("hub", p.HubName)
   f.SetProperty("address", p.Address)
   f.SetProperty("city", p.City)
   f.SetProperty("job_status", p.JobStatus)
   f.SetProperty("const_not_needed", p.ConstNeeded)
   fc.AddFeature(f)
} // }}}

// func MakeGeoJSONFile {{{
func MakeSimpleGeoJSONFile(pms [][]Permit) {
   fc := geojson.NewFeatureCollection()
   fileName := "/Users/sabra/go/src/mci/kml_v2/corepermits/data/coreRose.geojson"
   for _, folder := range pms {
      for _, p := range folder {
         if p.HubName != "ROSE BOWL (LAKE AVE)" {
            continue
         }
         var coords [][]float64

         for _, c := range p.Coordinates {
            var coord []float64
            lng := c.Lng
            coord = append(coord, lng)
            lat := c.Lat
            coord = append(coord, lat)

            coords = append(coords, coord)
         }
         if len(coords) < 3 {
            var coord []float64
            lng := p.Coordinates[len(coords)-1].Lng
            coord = append(coord, lng)

            lat := p.Coordinates[len(coords)-1].Lat
            coord = append(coord, lat)

            coords = append(coords, coord)
         }

         for _, d := range p.DYEAs {
            var dy DYEA

            if d.Number == "" {
               continue
            } else if d.Number == " " {
               continue
            } else if d.Number == "In Progress" {
               continue
            } else if d.Number == "EXPIRED: 7-MAY-20" {
               continue
            } else if d.Number == "VZ_LAN_00000203/DYEA_LSA_8107782" {
               // sep == sep[0] = VZ_LAN_00000203
               //        sep[1] = DYEA_LSA_8107782
               sep := strings.Split(d.Number, "/")

               dy.Number = sep[0]
               dy.Type = GetPermitType(d)

               //dy.Status = d.Status
               /*if len(typ) > 1 {
                  note := fmt.Sprintf("%s, %s", d.Note, typ[1])
                  dy.Note = note
               } else {
                  dy.Note = d.Note
               }*/
               MakeSFeature(dy, coords, p, fc)
               dy.Number = sep[1]

               MakeSFeature(dy, coords, p, fc)
               continue
            } else if d.Number == "DYEA_LSA_8108575/VZ_LAN_00003058" {
               sep := strings.Split(d.Number, "/")

               dy.Number = sep[0]
               dy.Type = GetPermitType(d)
               //dy.Status = d.Status

               /*if len(typ) > 1 {
                  note := fmt.Sprintf("%s, %s", d.Note, typ[1])
                  dy.Note = note
               } else {
                  dy.Note = d.Note
               }*/

               MakeSFeature(dy, coords, p, fc)

               dy.Number = sep[1]
               MakeSFeature(dy, coords, p, fc)
               continue
            }

            dy.Number = d.Number
            dy.Type = GetPermitType(d)
            //dy.Status = d.Status

            /*if len(typ) > 1 {
               note := fmt.Sprintf("%s, %s", d.Note, typ[1])
               dy.Note = note
            } else {
               dy.Note = d.Note
            }*/
            MakeSFeature(dy, coords, p, fc)
         }
      }
   }

   file, _ := json.MarshalIndent(fc, "", "    ")
   _ = ioutil.WriteFile(fileName, file, 0644)
}// }}}

// func MakeFeature {{{
//
// Makes a gejson feature
func MakeSFeature(dy DYEA, coords [][]float64, p Permit, fc *geojson.FeatureCollection) {
   /*typ := strings.ToLower(dy.Type)
   if len(typ) >= 3 {
      typ = typ[0:3]
   }

   if (typ != "exc" && typ != "jpa") {
      return
   } else if typ == "rr" {
      return
   }*/

   f := geojson.NewLineStringFeature(coords)
   f.SetProperty("object_id", p.ObjectId)
   f.SetProperty("dyea_vz", dy.Number)
   f.SetProperty("type", dy.Type)
   fc.AddFeature(f)
} // }}}


// func GetPermitType {{{
//
// Returns a normalized verison of the permit type
func GetPermitType(dy DYEA) string {
   ptype := dy.Type
   ptype = strings.ToLower(ptype)
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
   ptyp = ptype[0:3]

   if ptyp == "jpa" {
      return "ae"
   } else if ptyp == "tcp" {
      return "tcp"
   } else {
      return "ug"
   }
} // }}}
