// Sabra Bilodeau
package distance

import (
   "vci/kml/corepermits"
   "encoding/csv"
	"fmt"
	"os"
   "strings"
   "strconv"
   "regexp"
)

type DYVZ struct {
   Number         string
   Type           string
   Status         string
   Note           string
   Hub            string
   NFID           string
   FQNID          string
   Address        string
   City           string
   PermitUrl      string
   NetworkClass   string
   NetworkType    string
   ObjectId       string
   FileId         string
   JobStatus      string
   Coordinates    string
}

func ParseDYs(pms []corepermits.Permit) {
   var core map[int64]DYVZ
   core = make(map[int64]DYVZ)
   under := regexp.MustCompile(`_`)

   for _, p := range pms {
      var dy DYVZ
      dy.Hub = p.HubName
      dy.NFID = p.NFID
      dy.Address = p.Address
      dy.City = p.City
      dy.PermitUrl = p.PermitUrl
      dy.NetworkClass = p.NetworkClass
      dy.NetworkType = p.NetworkType
      dy.ObjectId = fmt.Sprintf("%d", p.ObjectId)
      dy.FileId = p.FileId
      dy.JobStatus = p.JobStatus

      var coords string
      for _, c := range p.Coordinates {
         coord := fmt.Sprintf("(%f,%f)", c.Lat, c.Lng)
         coords = coords + coord
      }
      dy.Coordinates = coords

      var fqs string
      fmt.Printf("FQNIDs: %s\n", p.FQNIDs)
      for _, fq := range p.FQNIDs {
         fqs = fqs + fq
      }
      dy.FQNID = fqs

      /* After parsing the core permits and outputting them to a CSV
      // I made note of any abnormalities in the DYEAs, and am dealing
      // with them in the if-else block */
      for _, d := range p.DYEAs {
         var idx int64
         fmt.Printf("ObjectID: %d\n", p.ObjectId)
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

            // split the VZ# & DYEA# up at '_', returning a 3 item array
            vz := under.Split(sep[0], 3)
            idx, _ = strconv.ParseInt(vz[2], 10, 64)
            core[idx] = dy

            dy.Number = sep[1]
            dyid := under.Split(sep[1], 4)
            idx, _ = strconv.ParseInt(dyid[2], 10, 64)
            core[idx] = dy
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

            // split the VZ# & DYEA# up at '_', returning a 3 item array
            vz := under.Split(sep[0], 3)
            idx, _ = strconv.ParseInt(vz[2], 10, 64)
            core[idx] = dy

            dy.Number = sep[1]
            dyid := under.Split(sep[1], 4)
            idx, _ = strconv.ParseInt(dyid[2], 10, 64)
            core[idx] = dy
            continue
         }

         if d.Number == "BEL-AF04-LOS-TCP" {
            idx = p.ObjectId            // Need to make note of this for later indexing
         } else if d.Number == "BEL-AF04" {
            idx = p.ObjectId            // Need to make note of this for later indexing
         } else if d.Number == "NOR-UG03-LOS" {
            idx = p.ObjectId            // Need to make note of this for later indexing
         } else if d.Number == "PCK-AF22-OL" {
            idx = p.ObjectId            // Need to make note of this for later indexing
         } else {
            permit := under.Split(d.Number, 4)
            fmt.Print(d.Number)
            idx, _ = strconv.ParseInt(permit[2], 10, 64)
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

         core[idx] = dy
      }
   }

   MakeCSVFile(core)
}

func ParseDYEAs(pms [][]corepermits.Permit) {
   var core map[int64]DYVZ
   core = make(map[int64]DYVZ)
   under := regexp.MustCompile(`_`)

   for _, folder := range pms {
      for _, p := range folder {
         var dy DYVZ
         dy.Hub = p.HubName
         dy.NFID = p.NFID
         dy.Address = p.Address
         dy.City = p.City
         dy.PermitUrl = p.PermitUrl
         dy.NetworkClass = p.NetworkClass
         dy.NetworkType = p.NetworkType
         dy.ObjectId = fmt.Sprintf("%d", p.ObjectId)
         dy.FileId = p.FileId
         dy.JobStatus = p.JobStatus

         var coords string
         for _, c := range p.Coordinates {
            coord := fmt.Sprintf("(%f,%f)", c.Lat, c.Lng)
            coords = coords + coord
         }
         dy.Coordinates = coords

         var fqs string
         fmt.Printf("FQNIDs: %s\n", p.FQNIDs)
         for _, fq := range p.FQNIDs {
            fqs = fqs + fq
         }
         dy.FQNID = fqs

         /* After parsing the core permits and outputting them to a CSV
         // I made note of any abnormalities in the DYEAs, and am dealing
         // with them in the if-else block */
         for _, d := range p.DYEAs {
            var idx int64
            fmt.Printf("ObjectID: %d\n", p.ObjectId)
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

               // split the VZ# & DYEA# up at '_', returning a 3 item array
               vz := under.Split(sep[0], 3)
               idx, _ = strconv.ParseInt(vz[2], 10, 64)
               core[idx] = dy

               dy.Number = sep[1]
               dyid := under.Split(sep[1], 4)
               idx, _ = strconv.ParseInt(dyid[2], 10, 64)
               core[idx] = dy
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

               // split the VZ# & DYEA# up at '_', returning a 3 item array
               vz := under.Split(sep[0], 3)
               idx, _ = strconv.ParseInt(vz[2], 10, 64)
               core[idx] = dy

               dy.Number = sep[1]
               dyid := under.Split(sep[1], 4)
               idx, _ = strconv.ParseInt(dyid[2], 10, 64)
               core[idx] = dy
               continue
            }

            if d.Number == "BEL-AF04-LOS-TCP" {
               idx = p.ObjectId            // Need to make note of this for later indexing
            } else if d.Number == "BEL-AF04" {
               idx = p.ObjectId            // Need to make note of this for later indexing
            } else if d.Number == "NOR-UG03-LOS" {
               idx = p.ObjectId            // Need to make note of this for later indexing
            } else if d.Number == "PCK-AF22-OL" {
               idx = p.ObjectId            // Need to make note of this for later indexing
            } else {
               permit := under.Split(d.Number, 4)
               fmt.Printf("%s\n", d.Number)
               idx, _ = strconv.ParseInt(permit[2], 10, 64)
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

            core[idx] = dy
         }
      }
   }

   MakeCSVFile(core)
}

func MakeCSVFile(core map[int64]DYVZ) {
   f, err := os.Create("./core.csv")

   if err != nil {
      fmt.Println(err)
   }

   defer f.Close()

   var record []string
   record = append(record, "DYEA/VZ#")
   //record = append(record, "FQNID")
   record = append(record, "NFID")
   record = append(record, "Address")
   record = append(record, "Hub")
   record = append(record, "City")
   record = append(record, "Type")
   record = append(record, "Status")
   record = append(record, "URL")
   record = append(record, "Note")
   record = append(record, "Network Class")
   record = append(record, "Network Type")
   record = append(record, "Live Maps Object ID")
   record = append(record, "Live Maps Layer")
   record = append(record, "Live Maps Object Job Status")
   record = append(record, "Coordinates")

   w := csv.NewWriter(f)
   w.Write(record)

   for _, permit := range core {
      var record []string
      record = append(record, permit.Number)
      //record = append(record, permit.FQNID)
      record = append(record, permit.NFID)
      record = append(record, permit.Address)
      record = append(record, permit.Hub)
      record = append(record, permit.City)
      record = append(record, permit.Type)
      record = append(record, permit.Status)
      record = append(record, permit.PermitUrl)
      record = append(record, permit.Note)
      record = append(record, permit.NetworkClass)
      record = append(record, permit.NetworkType)
      record = append(record, permit.ObjectId)
      record = append(record, permit.FileId)
      record = append(record, permit.JobStatus)
      record = append(record, permit.Coordinates)
      w.Write(record)
   }
   w.Flush()
}
