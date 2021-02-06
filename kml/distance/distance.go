// Sabra Bilodeau
package distance

// Import necessary libraries.
import (
   "fmt"
   "math"
   "encoding/json"
   "strings"
   "io/ioutil"
   "vci/kml/corepermits"
   "vci/kml/vz3gis"
   "github.com/paulmach/go.geojson"
)

type Match struct {
   ObjectId   int64
   VzObjectId int64
   FQNID      string
   NFID       string
   CalcLength float64
   MeasLength float64
   FiberCt    int64
   DYEAs      [6]corepermits.DYEA
   Hub        string
   Address    string
   City       string
   UDFs       []string
   Point      corepermits.Point
   VzPoint    vz3gis.Point
}

// func CompareFiles {{{
//
// Compares the VZ & Core files, adding any matches to our matches array
func CompareFiles(vz [][]vz3gis.FQNID, core [][]corepermits.Permit) []Match {
   var newCore []corepermits.Permit
   var matches []Match
   var m Match
   found := false

   for _, folder := range core {
      for _, permit := range folder {
         var fqs []string
         for _, folderVz := range vz {
            for _, fq := range folderVz {
               found, m = ComparePoints(fq, permit)
               if found {
                  fqs = append(fqs, fq.FQNID)
                  m.DYEAs = permit.DYEAs
                  m.UDFs = fq.UDFs
                  matches = append(matches, m)
               }
            }
         }
         newCore = append(newCore, permit)
      }
   }
   fmt.Printf("Total matches found: %d\n", len(matches))
   //MakeJSONFile(matches)
   //MakeMatchesCSVFile(matches)
   MakeGeoJSONFile(matches)
   //MakeTrainingData(matches)
   return matches
} // }}}

func CompareFQs(vz [][]vz3gis.FQNID, core [][]corepermits.Permit) {
   var newFQs []vz3gis.FQNID
   var matches []Match
   var m Match
   found := false

   for _, folder := range vz {
      for _, fq := range folder {
         var permits []string
         for _, foldercore := range core {
            for _, permit := range foldercore {
               found, m = ComparePoints(fq, permit)
               if found {
                  for _, dy := range permit.DYEAs {
                     permits = append(permits, dy.Number)
                  }
                  m.DYEAs = permit.DYEAs
                  m.UDFs = fq.UDFs
                  matches = append(matches, m)
               }
            }
         }
         fq.DYEAs = permits
         newFQs = append(newFQs, fq)
      }
   }
   //fmt.Println(vz)
   fmt.Printf("Total matches found: %d\n", len(matches))
   //MakeJSONFile(matches)
   //MakeFQNIDCSVFile(newFQs)

}

// func ComparePoints {{{
//
func ComparePoints(vz vz3gis.FQNID, core corepermits.Permit) (bool, Match) {
   var m Match
   found := false
   aerialfq, tcp, mistype := false, false, false

   fqnid := strings.SplitAfter(vz.FQNID, ":")
   fqn := fqnid[1]
   fq := strings.Split(fqn, ":")
   fqType := fq[0]

   if fqType == "AER" {
      aerialfq = true
   } else if fqType == "AERSPAN" {
      aerialfq = true
   }

   for _, dy := range core.DYEAs {
      if len(dy.Type) < 3 {
         continue
      }
      ptype := dy.Type[0:3]
      ptype = strings.ToLower(ptype)
      switch ptype {
      case "tcp":
         tcp = true
      case "exc":
         if aerialfq {
            mistype = true
         } else {
            mistype = false
         }
      case "jpa":
         if !aerialfq {
            mistype = true
         } else {
            mistype = false
         }
      }
   }

   if tcp {
   //   fmt.Printf("FQ type: %s\n", fqType)
   //   fmt.Printf("Permt Type: TCP\n")
      return found, m
   } else if mistype {
   //   fmt.Printf("FQ type: %s\n", fqType)
   //   fmt.Printf("Permt Type: %s\n", pty)
      return found, m
   }

   for _, p := range core.Coordinates {
      if p.Matched {
         continue
      }
      for _, p_vz := range vz.Coordinates {
         if p_vz.Matched {
            continue
         }
         d := Distance(p, p_vz)
         //fmt.Printf("Distance: %f", d)
         if d < 0.00096 {
            p.Matched = true
            p_vz.Matched = true
            m.ObjectId = core.ObjectId
            m.VzObjectId = vz.ObjectId
            m.FQNID = vz.FQNID
            m.Hub = core.HubName
            m.NFID = core.NFID
            m.CalcLength = vz.CalculatedLength
            m.MeasLength = vz.MeasuredLength
            m.FiberCt = vz.FiberCount
            m.DYEAs = core.DYEAs
            m.Address = core.Address
            m.City = core.City
            m.UDFs = vz.UDFs
            m.Point = p
            m.VzPoint = p_vz
            found = true
            return found, m
         }
      }
   }
   return found, m
} // }}}

// func Match Coords {{{
//
// Matches a point to an FQNID
func MatchCoords(point corepermits.Point, vz [][]vz3gis.FQNID, core [][]corepermits.Permit) (int64, string, []string, string) {
   var fty string

   for _, folder := range vz {
      for _, fq := range folder {
         fqnid := strings.SplitAfter(fq.FQNID, ":")
         fqn := fqnid[1]
         f := strings.Split(fqn, ":")
         fqType := f[0]
         fty = fqType
         if fty != "BUR" {
            if fty != "UGSPAN" {
               continue
            }
            continue
         }

         points := fq.Coordinates
         for _, p := range points {
            d := Distance(point, p)
         //   fmt.Printf("Distance found is - %f\n", d)
            if d < 0.00096 {                        // 5FT
               NFID := MatchCore(point, core)
               return fq.ObjectId, fq.FQNID, fq.UDFs, NFID
            }
         }
      }
   }

   var udfs []string

   return -1, "", udfs, ""
} // }}}

// func MatchCore {{{
func MatchCore(point corepermits.Point, core [][]corepermits.Permit) string {
   for _, folder := range core {
      for _, permit := range folder {
         points := permit.Coordinates
         for _, p := range points {
            d := DistanceCore(point, p)
         //   fmt.Printf("Distance found is - %f\n", d)
            if d < 0.00096 {                       // 5FT
               return permit.NFID
            }
         }
      }
   }
   return ""
} // }}}

// func Distance {{{
//
// Calculates the distance between 2 points using the distance formula
func Distance(p1 corepermits.Point, p2 vz3gis.Point) float64 {
   x1 := p1.Lat
   y1 := p1.Lng

   x2 := p2.Lat
   y2 := p2.Lng

   x := x2 - x1
   y := y2 - y1

   d := math.Sqrt((x*x)+(y*y))
   return d
} // }}}

// func Distance {{{
//
// Calculates the distance between 2 points using the distance formula
func DistanceCore(p1 corepermits.Point, p2 corepermits.Point) float64 {
   x1 := p1.Lat
   y1 := p1.Lng

   x2 := p2.Lat
   y2 := p2.Lng

   x := x2 - x1
   y := y2 - y1

   d := math.Sqrt((x*x)+(y*y))
   return d
} // }}}

// func MakeJSONFile {{{
//
func MakeJSONFile(matches []Match) {
   fileName := "matches.json"
   file, _ := json.MarshalIndent(matches, "", "    ")
   _ = ioutil.WriteFile(fileName, file, 0644)
} // }}}


// func MakeGeoJSONFile {{{
func MakeGeoJSONFile(matches []Match) {
   fc := geojson.NewFeatureCollection()
   fileName := "matches.geojson"
   for _, m := range matches {
      var coords [][]float64
      var c []float64
      var v []float64
      lng := m.Point.Lng
      c = append(c, lng)
      lat := m.Point.Lat
      c = append(c, lat)
      coords = append(coords, c)
      coords = append(coords, c)

      lng = m.VzPoint.Lng
      v = append(v, lng)
      lat = m.VzPoint.Lat
      v = append(v, lat)

      coords = append(coords, v)

      for _, d := range m.DYEAs {
         var dy corepermits.DYEA
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
            MakeFeature(dy, coords, m, fc)

            dy.Number = sep[1]
            MakeFeature(dy, coords, m, fc)
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
            MakeFeature(dy, coords, m, fc)

            dy.Number = sep[1]
            MakeFeature(dy, coords, m, fc)
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
         MakeFeature(dy, coords, m, fc)
      }
   }
   file, _ := json.MarshalIndent(fc, "", "    ")
   _ = ioutil.WriteFile(fileName, file, 0644)
} // }}}

// func MakeFeature {{{
func MakeFeature(dy corepermits.DYEA, coords [][]float64, m Match, fc *geojson.FeatureCollection) {
   // For data mining purposes, only taking the permits that are of type exc & jpa
   typ := strings.ToLower(dy.Type)
   if len(typ) >= 3 {
      typ = typ[0:3]
   }
   if (typ != "exc" && typ != "jpa") {
      return
   } else if typ == "rr" {
      return
   }
   if dy.Status == "" {
      return
   }
   if dy.Status == "TIME OUT" {
      dy.Status = "APPROVED"
   }

   f := geojson.NewLineStringFeature(coords)
   f.SetProperty("object_id", m.ObjectId)
   f.SetProperty("vz_object_id", m.VzObjectId)
   f.SetProperty("dyea_vz", dy.Number)
   f.SetProperty("type", typ)
   f.SetProperty("status", dy.Status)
   f.SetProperty("note", dy.Note)
   f.SetProperty("fqnid", m.FQNID)
   f.SetProperty("nfid", m.NFID)
   f.SetProperty("hub", m.Hub)
   f.SetProperty("address", m.Address)
   f.SetProperty("city", strings.ToLower(m.City))
   f.SetProperty("calculatedlength", m.CalcLength)
   f.SetProperty("measuredlength", m.MeasLength)
   f.SetProperty("fiber_count", m.FiberCt)
   f.SetProperty("fqnid_udfs", m.UDFs)
   fc.AddFeature(f)
} // }}}
