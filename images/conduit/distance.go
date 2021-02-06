// Sabra Bilodeau
package conduit

// Import necessary libraries.
import (
   "math"
   "strings"
   "mci/gjson"
)

// func Match VZ {{{
//
// Matches a point to an FQNID
func MatchVZ(point Point, vz gjson.VZFeatureCollection) (int64, string, []string, string) {
   var fty string

   for _, fq := range vz.Features {
      // Check the FQNID type.
      // These are conduit tie in photos, so they should
      // ONLY match up with underground FQNIDs -
      // FIB:BUR::XXXXXX || FIB:TAIL::XXXXXX || OSP:UGSPAN::XXXXXX
      fqnid := strings.SplitAfter(fq.Properties.FQNID, ":")
      fqn := fqnid[1]
      f := strings.Split(fqn, ":")
      fqType := f[0]
      fty = fqType
      if fty != "BUR" {
         if fty != "UGSPAN" {
            if fty != "TAIL" {
               continue
            }
         }
      }

      // Now let's see if the image lines up with any point in the FQNID.
      points := fq.Geometry.Coordinates
      for _, p := range points {
         d := Distance(point, p)
      //   fmt.Printf("Distance found is - %f\n", d)
         if d < 0.00096 {                        // 5FT
            return fq.Properties.ObjectId, fq.Properties.FQNID,
                   fq.Properties.UDFs, fq.Properties.NFID
         }
      }
   }

   var udfs []string

   return -1, "", udfs, ""
} // }}}

// func MatchCore {{{
func MatchCore(point Point, core gjson.CoreFeatureCollection) string {
   for _, permit := range core.Features {
      points := permit.Geometry.Coordinates
      for _, p := range points {
         d := Distance(point, p)
      //   fmt.Printf("Distance found is - %f\n", d)
         if d < 0.00096 {                       // 5FT
            return permit.Properties.NFID
         }
      }
   }
   return ""
} // }}}

// func Distance {{{
//
// Calculates the distance between 2 points using
// the distance formula
func Distance(p1 Point, p2 []float64) float64 {
   x1 := p1.Lng
   y1 := p1.Lat

   x2 := p2[0]
   y2 := p2[1]

   x := x2 - x1
   y := y2 - y1

   d := math.Sqrt((x*x)+(y*y))
   return d
} // }}}
