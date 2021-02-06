// Parses the CORE_PERMITS.KML file found in the data folder,
// extracts the desired permit information from each placemark and
// saves it to a geojson file to be used for later data plotting.
package corepermits

import (
	"encoding/xml"
   "encoding/json"
	"io/ioutil"
	"fmt"
	"os"
   "strings"
   "strconv"
)

// We are specifically looking through the KML for all the <Placemark> entries.
type kmlFile struct {
	XMLName xml.Name
	Doc     kmlDocument `xml:"Document"`
	Atts    []xml.Attr  `xml:",any,attr"`

	// If this isn't empty, then our structure needs updating, as something exists beyond "<Document>"
	BeEmpty []xmlStuff
}

type kmlDocument struct {
	XMLName xml.Name

	Folders []kmlFolder `xml:"Folder>Folder"`

	// Things like the document ID, if we neeed it.
	Atts    []xml.Attr  `xml:",any,attr"`

	// Anything left over that we haven't specifically looked for above.
	Stuff   []xmlStuff  `xml:",any"`
}

type xmlStuff struct {
	XMLName xml.Name
	Value   string     `xml:",chardata"`
	Attrs   []xml.Attr `xml:",any,attr"`
	Other   []xmlStuff `xml:",any"`
}

type kmlFolder struct {
	XMLName xml.Name

	ID string `xml:"id,attr"`

	SimpleData    []kmlSimpleData `xml:"Placemark>ExtendedData>SchemaData>SimpleData"`
   Coordinates   []string        `xml:"Placemark>LineString>coordinates"`
   MultiGeom     []string        `xml:"Placemark>MultiGeometry>LineString>coordinates"`

	// As we are not 100% sure we accounted for everything within the Folder, we allow the parser to place anything without a named column in these values. below
	Value   string     `xml:",chardata"`
	Attrs   []xml.Attr `xml:",any,attr"`
	Other   []xmlStuff `xml:",any"`
}

type kmlSimpleData struct {
	Name	string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// The information we want to extract from each Placemark
type Permit struct {
   ObjectId     int64   `json:"object_id"`
   FileId       string  `json:"file_id"`
   HubName      string  `json:"hub_name"`
   NetworkClass string  `json:"network_class"`
   NetworkType  string  `json:"network_type"`
   NFID         string  `json:"nfid"`
   PermitUrl    string  `json:"permit_url"`
   Address      string  `json:"address"`
   City         string  `json:"city"`
   DYEAs        [6]DYEA `json:"dyeas"`
   JobStatus    string  `json:"job_status"`
   Coordinates  []Point `json:"coordinates"`
   Created      string  `json:'created_date'`
   FQNIDs       []string `json:"fqnids"`
   ConstNeeded  string   `json:"const_not_needed"`
}

// The number here can expressed as either
// DYEA_LSA_XXXXXXX or
// VZ_LAN_XXXXXXXX
type DYEA struct {
   Number string  `json:"number"`
   Type   string  `json:"type"`
   Status string  `json:"status"`
   Note   string  `json:"note"`
}

type Point struct {
   Lat      float64  `json:"lat"`
   Lng      float64  `json:"lng"`
   Matched  bool     `json:"matched"`
}

// func ParseKML {{{
//
func ParseKML() ([][]Permit, error) {
   var pms []Permit
   var all [][]Permit
   total := 0
   file := "/Users/sabra/go/src/clustering/kml/corepermits/data/CORE_PERMITS.kml"

	fmt.Printf("Loading file: %s\n", file)

	kml := &kmlFile{}
	if err := loadXML(file, kml); err != nil {
		fmt.Printf("Error loading KML: %s", err)
		return all, err
	}

	// If you want to see how the document is loaded internally, uncomment this line.
	//fmt.Printf("KML: %#v\n", kml)

   for _, folder := range kml.Doc.Folders {
		fmt.Printf("\nFolder ID %s\n", folder.ID)


      /*if (folder.ID == "CORE_PERMITS_IN_PROGRESS") || (folder.ID == "CORE_PERMITS_SUBMITTED_RESUBMITTED") {
         fmt.Printf("Actually ... nevermind for right now.\n")
         continue
      }*/

      var pm Permit
      i := 0
      x := 0
		for _, sd := range folder.SimpleData {
		//	fmt.Printf("- ExtendedData: %s = %s\n", sd.Name, sd.Value)
         pm = getSimpleData(pm, sd)
         pm.FileId = folder.ID
         if sd.Name == "st_length(shape)" {
            if pm.ObjectId == 34303 || pm.ObjectId == 28884 || pm.ObjectId == 31752 || pm.ObjectId == 33548 || pm.ObjectId == 31463|| pm.ObjectId == 34436 || pm.ObjectId == 173196 || pm.ObjectId == 182012 {
               coord := folder.MultiGeom[x]
               nonewline := strings.Replace(coord, "\n", "", -1)
               notab := strings.Replace(nonewline, "\t", "", -1)
               c := strings.Split(notab, " ")
               pm.Coordinates = getCoordinates(c)
               pms = append(pms, pm)
               x++

               coord = folder.MultiGeom[x]
               nonewline = strings.Replace(coord, "\n", "", -1)
               notab = strings.Replace(nonewline, "\t", "", -1)
               c = strings.Split(notab, " ")
               pm.Coordinates = getCoordinates(c)
               pms = append(pms, pm)
               x++
            } else if pm.ObjectId == 34304 {
               coord := folder.MultiGeom[x]
               nonewline := strings.Replace(coord, "\n", "", -1)
               notab := strings.Replace(nonewline, "\t", "", -1)
               c := strings.Split(notab, " ")
               pm.Coordinates = getCoordinates(c)
               pms = append(pms, pm)
               x++

               coord = folder.MultiGeom[x]
               nonewline = strings.Replace(coord, "\n", "", -1)
               notab = strings.Replace(nonewline, "\t", "", -1)
               c = strings.Split(notab, " ")
               pm.Coordinates = getCoordinates(c)
               pms = append(pms, pm)

               x++
               coord = folder.MultiGeom[x]
               nonewline = strings.Replace(coord, "\n", "", -1)
               notab = strings.Replace(nonewline, "\t", "", -1)
               c = strings.Split(notab, " ")
               pm.Coordinates = getCoordinates(c)
               pms = append(pms, pm)
               x++
            } else if pm.ObjectId == 34306 {
               for n := 0; n < 6; n++ {
                  coord := folder.MultiGeom[x]
                  nonewline := strings.Replace(coord, "\n", "", -1)
                  notab := strings.Replace(nonewline, "\t", "", -1)
                  c := strings.Split(notab, " ")
                  pm.Coordinates = getCoordinates(c)
                  pms = append(pms, pm)
                  x++
               }
            } else if i != len(folder.Coordinates) {
               coord := folder.Coordinates[i]
               nonewline := strings.Replace(coord, "\n", "", -1)
               notab := strings.Replace(nonewline, "\t", "", -1)
               c := strings.Split(notab, " ")
               pm.Coordinates = getCoordinates(c)
               pms = append(pms, pm)
               i = i + 1
            }

            pm.DYEAs = ResetDYEAs()

            //fmt.Printf("- ObjectId: %s\n", pms[i].ObjectId)
            //fmt.Printf("- Coordinates: %f\n", pms[i].Coordinates)
            //fmt.Printf("\n")
         }
		}
      all = append(all, pms)
      total = total + len(pms)
      pms = ResetArray()
	}
   MakeGeoJSONFile(all)
   //MakeSimpleGeoJSONFile(all)
   //ParseDYEAs(all)   // Function in separatepermits.go
   fmt.Printf("Number of Core permits: %d\n", total)

   return all, nil
} // }}}

// func ResetDYEAs {{{
//
// I noticed I was getting repeating values in my return,
// so I'm resetting the DYEA array just in case.
func ResetDYEAs() [6]DYEA {
   var dy [6]DYEA
   return dy
} // }}}

// func ResetArray() {{{
//
func ResetArray() []Permit {
   var pms []Permit
   return pms
} // }}}

// func makeJSONFile {{{
//
// Writes the array of Permits found in the current KML folder
// to a JSON file, labelled by its FolderID
func makeJSONFile(pms []Permit, name string) {
   fileName := name + ".json"
   file, _ := json.MarshalIndent(pms, "", "    ")
   _ = ioutil.WriteFile(fileName, file, 0644)
} // }}}

// func getSimpleData {{{
//
// Get the information we want from the KML Simple Data fields
// We take in the permit to ensure we don't overwrite any fields
// during the return
func getSimpleData(pm Permit, sd kmlSimpleData) Permit {
   switch sd.Name {
      case "objectid":
         pm.ObjectId, _ = strconv.ParseInt(sd.Value, 10, 64)
      case "bpg_hub_name":
         pm.HubName = sd.Value
      case "bpg_network_class":
         pm.NetworkClass = sd.Value
      case "bpg_network_type":
         pm.NetworkType = sd.Value
      case "created_date":
         pm.Created = sd.Value
      case "site_nfid":
         pm.NFID = sd.Value
      case "permit_url":
         pm.PermitUrl = sd.Value
      case "address":
         pm.Address = sd.Value
      case "city":
         pm.City = sd.Value
      case "dyea_1":
         pm.DYEAs[0].Number = sd.Value
      case "permit_type_1":
         pm.DYEAs[0].Type = sd.Value
      case "permit_status_1":
         pm.DYEAs[0].Status = sd.Value
      case "permit_note_1":
         pm.DYEAs[0].Note = sd.Value
      case "dyea_2":
         pm.DYEAs[1].Number = sd.Value
      case "permit_type_2":
         pm.DYEAs[1].Type = sd.Value
      case "permit_status_2":
         pm.DYEAs[1].Status = sd.Value
      case "permit_note_2":
         pm.DYEAs[1].Note = sd.Value
      case "dyea_3":
         pm.DYEAs[2].Number = sd.Value
      case "permit_type_3":
         pm.DYEAs[2].Type = sd.Value
      case "permit_status_3":
         pm.DYEAs[2].Status = sd.Value
      case "permit_note_3":
         pm.DYEAs[2].Note = sd.Value
      case "dyea_4":
         pm.DYEAs[3].Number = sd.Value
      case "permit_type_4":
         pm.DYEAs[3].Type = sd.Value
      case "permit_status_4":
         pm.DYEAs[3].Status = sd.Value
      case "permit_note_4":
         pm.DYEAs[3].Note = sd.Value
      case "dyea_5":
         pm.DYEAs[4].Number = sd.Value
      case "permit_type_5":
         pm.DYEAs[4].Type = sd.Value
      case "permit_status_5":
         pm.DYEAs[4].Status = sd.Value
      case "permit_note_5":
         pm.DYEAs[4].Note = sd.Value
      case "dyea_6":
         pm.DYEAs[5].Number = sd.Value
      case "permit_type_6":
         pm.DYEAs[5].Type = sd.Value
      case "permit_status_6":
         pm.DYEAs[5].Status = sd.Value
      case "permit_note_6":
         pm.DYEAs[5].Note = sd.Value
      case "const_not_needed":
         pm.ConstNeeded = sd.Value
      case "job_status":
         pm.JobStatus = sd.Value
   }
   return pm
} // }}}

// func getCoordinates {{{
//
// Parses the placemarks coordinates
func getCoordinates(coords []string) []Point {
   var ls []Point

   for _, coord := range coords {
      var p Point
      //fmt.Print("Coord: \n", coord)
      c := strings.Split(coord, ",")
      //fmt.Print("C: \n", c)
      if coord == "" {
         break
      }

      p.Lng, _ = strconv.ParseFloat(c[0], 64)
      p.Lat, _ = strconv.ParseFloat(c[1], 64)
      p.Matched = false

      ls = append(ls, p)
   }

   return ls
} // }}}

// func loadXML }}}
// Attempts to decode the provided file into the provided interface.
//
// Returns an error if there were any decoding issues.
func loadXML(file string, iface interface{}) error {
	in, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("open(%s)(: %s", file, err)
	}

	defer in.Close()

	dec := xml.NewDecoder(in)

	if err = dec.Decode(iface); err != nil {
		return fmt.Errorf("Decode(%s): %s", file, err)
	}

	return nil
} // }}}
