// Sabra Bilodeau
package vz3gis

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
	Doc kmlDocument `xml:"Document"`
	Atts []xml.Attr `xml:",any,attr"`

	// If this isn't empty, then our structure needs updating, as something exists beyond "<Document>"
	BeEmpty []xmlStuff
}

type kmlDocument struct {
	XMLName xml.Name

	Folders []kmlFolder `xml:"Folder>Folder>Folder"`

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
   Coordinates   []string `xml:"Placemark>LineString>coordinates"`

	// As we are not 100% sure we accounted for everything within the Folder, we allow the parser to place anything without a named column in these values. below
	Value   string     `xml:",chardata"`
	Attrs   []xml.Attr `xml:",any,attr"`
	Other   []xmlStuff `xml:",any"`
}

type kmlSimpleData struct {
	Name	string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type FQNID struct {
   ObjectId          int64    `json:"object_id"`
   FileID            string   `json:"file_id"`
   CalculatedLength  float64  `json:"calculated_length"`
   MeasuredLength    float64  `json:"measured_length"`
   CableType         string   `json:"cable_type"`
   FiberCount        int64    `json:"fiber_count"`
   FQNID             string   `json:"fqnid"`
   NFID              string   `json:"nfid"`
   UDFs              []string `json:"udfs"`
   OriginalFQNID     string   `json:"original_fqnid"`
   WorkOrderName     string   `json:"work_order_name"`
   Coordinates       []Point  `json:"coordinates"`
   Created           string   `json:"created_date"`
   EWO               string   `json:"edit_work_order_id"`
   DYEAs             []string `json:"dyeas"`
}

type Point struct {
   Lat      float64  `json:"lat"`
   Lng      float64  `json:"lng"`
   Matched  bool     `json:"matched"`
}

// func ParseKML {{{
//
func ParseKML() ([][]FQNID, error) {
   var fqs []FQNID
   var all [][]FQNID

	file := "/Users/sabra/go/src/clustering/kml/vz3gis/data/VZ3GIS.kml"

	fmt.Printf("Loading file: %s\n", file)

	kml := &kmlFile{}
	if err := loadXML(file, kml); err != nil {
		fmt.Printf("Error loading KML: %s", err)
	   return all, err
	}

	// If you want to see how the document is loaded internally, uncomment this line.
	//fmt.Printf("KML: %#v\n", kml)

   for _, folder := range kml.Doc.Folders {
      /* There are 8 folders in this KML, but only these 2 have data in them.
      if folder.ID != "CPD_CABLE" {
         if folder.ID != "CPD_CONDUIT" {
             continue
          }
      }*/
		fmt.Printf("\nFolder ID %s\n", folder.ID)


      i := 0
      var fq FQNID
		for _, sd := range folder.SimpleData {
		//	fmt.Printf("- ExtendedData: %s = %s\n", sd.Name, sd.Value)
         fq = getSimpleData(fq, sd)
         fq.FileID = folder.ID
         if sd.Name == "work_order_name" {
            coord := folder.Coordinates[i]
            nonewline := strings.Replace(coord, "\n", "", -1)
            notab := strings.Replace(nonewline, "\t", "", -1)
            c := strings.Split(notab, " ")
            fq.Coordinates = getCoordinates(c)

            fqs = append(fqs, fq)
            fq.UDFs = ResetUDFs()

            //fmt.Printf("- ObjectId: %s\n", fqs[i].ObjectId)
            //fmt.Printf("- Coordinates: %f\n", fqs[i].Coordinates)
            //fmt.Printf("\n")

            i++
            if i == len(folder.Coordinates) {
               break
            }
         }
		}
      fmt.Printf("Number of FQNIDs in folder: %d\n", len(fqs))
      //MakeJSONFile(fqs, folder.ID)
      all = append(all, fqs)
      fqs = ResetArray()
	}
   MakeGeoJSONFile(all)
   MakeCSVFile(all)
   //MakeSimpleGeoJSONFile(all)
   return all, nil
} // }}}

// func ResetArray {{{
//
// I noticed I was getting repeating values in my return,
// so I'm resetting the FQNID array just in case.
func ResetArray() []FQNID {
   var fqs []FQNID
   return fqs
} // }}}

// func ResetUDFs {{{
//
// I noticed I was getting repeating values in my return,
// so I'm resetting the UDFs just in case.
func ResetUDFs() []string {
   var udfs []string
   return udfs
} // }}}

// func MakeJSONFile {{{
//
// Writes the array of FQNIDs found in the current KML folder
// to a JSON file, labelled by its FolderID
func MakeJSONFile(fqs []FQNID, name string) {
   fileName := name + ".json"
   file, _ := json.MarshalIndent(fqs, "", "    ")
   _ = ioutil.WriteFile(fileName, file, 0644)
} // }}}


// func getSimpleData {{{
//
// Get the information we want from the KML Simple Data fields
// We take in the fq to ensure we don't overwrite any fields during the return
func getSimpleData(fq FQNID, sd kmlSimpleData) FQNID {
   data := fq
   switch sd.Name {
      case "objectid":
         data.ObjectId, _ = strconv.ParseInt(sd.Value, 10, 64)
      case "datecreated":
         data.Created = sd.Value
      case "calculatedlength":
         data.CalculatedLength, _ = strconv.ParseFloat(sd.Value, 64)
      case "measuredlength":
         data.MeasuredLength, _ = strconv.ParseFloat(sd.Value, 64)
      case "cabletype":
         data.CableType = sd.Value
      case "fibercount":
         // I made FiberCount a uint16 initially but it won't do it as that ???
         data.FiberCount, _ = strconv.ParseInt(sd.Value, 10, 16)
      case "fqn_id":
         data.FQNID = sd.Value
      case "site_span_nfid":
         data.NFID = sd.Value
      case "udf1", "udf2", "udf3", "udf4", "udf5", "udf6", "udf7", "udf8", "udf9", "udf10", "udf11", "udf12":
         if sd.Value != "073735-001" {
            // 073735-001 is the Capital Project Number and is frequently
            // listed in the UDFs, but I don't care about that value
            data.UDFs = append(data.UDFs, sd.Value)
         }
      case "originating_fqnid":
         data.OriginalFQNID = sd.Value
      case "editworkorderid":
         data.EWO = sd.Value
      case "work_order_name":
         data.WorkOrderName = sd.Value
   }
   return data
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

// func loadXML {{{
//
// Attempts to decode the provided file into the provided interface.
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
