// Sabra Bilodeau
package conduit

import (
    "log"
    "io"
    "os"
    "path/filepath"
    "encoding/json"
    "io/ioutil"
    "regexp"
    "strings"
    "strconv"
    "vci/gjson"
    "github.com/rwcarlsen/goexif/exif"
    "github.com/rwcarlsen/goexif/tiff"
)

type ImageFile struct {
   Name           string
   NewFile        string
   Path           string
   Hub            string
   ObjectId       int64
   FQNID          string
   NFID           string
   WorkOrderName  string
   UDFs           []string
   Point          Point
}

type Point struct {
   Lat float64
   Lng float64
}

var LAT, LNG float64
// func Visit {{{
// src : https://flaviocopes.com/go-list-files/
//
// To be used to get all file paths in our Conduit_Tie_Ins directory
func Visit(files *[]string, names *[]string, nometa *[]ImageFile) filepath.WalkFunc {
   return func(path string, info os.FileInfo, err error) error {
      if err != nil {
         var nm ImageFile
         switch filepath.Ext(path) {
         case ".jpg", ".png", ".JPG":
            // TO DO -- SEND THE FILE TO OCR
            nm.Name = info.Name()
            nm.Path = path
            *nometa = append(*nometa, nm)
            log.Printf("REMINDER -- TO DO HERE -- SEND FILE TO OCR. Skipping file : %s, no metadata\n", nm.Name)
         return nil
         }
      }

      switch filepath.Ext(path) {
      case ".jpg", ".png", ".JPG":
         *files = append(*files, path)
         *names = append(*names, info.Name())
         return nil
      default:
         log.Printf("Skipping file : %s, not an image\n", info.Name())
      }
      return nil
   }
} // }}}

type Printer struct {}

// func Walk {{{
// src : https://stackoverflow.com/questions/60497938/read-exif-metadata-with-go
//
// Loops through all the files, pulling their metadata, if there is any
// then converting the parsed Lat/Lng coordinates to decimal digits.
func (p Printer) Walk(name exif.FieldName, tag *tiff.Tag) error {
   //fmt.Printf("%40s: %s\n", name, tag)

   t := tag.String()

   //             + Values    - Values
   // Latidude :     "N"         "S"
   // Longitude :    "E"         "w"
   if name == "GPSLatitudeRef" {
      if t == "S" {
         log.Printf("Latidude = S - Southern Hempishere ? ok ..\n")
         LAT = -1.0
      }
   } else if name == "GPSLongitudeRef" {
      if t == "E" {
         log.Printf("Longitude = E - Eastern Hempishere ? ok ..\n")
         LNG = 1.0
      }
   } else if name == "GPSLatitude" {
      stag := tag.String()
      // Multiplying by lat here in case uh.. for what ever reason we needed to change hemispheres ?
      LAT = ParseLatitude(stag) * LAT
      log.Printf("Latidude value : %f\n", LAT)
   } else if name == "GPSLongitude" {
      stag := tag.String()
      // Multiplying by lat here in case uh.. for what ever reason we needed to change hemispheres ?
      LNG = ParseLongitude(stag) * LNG
      log.Printf("Longitude value : %f\n", LNG)
   }
   return nil
} // }}}

// func ParseLatitude {{{
// Code assistance from: https://yourbasic.org/golang/split-string-into-slice/
//
// Parses the Latidude value of the coordinate and converts it to a decimal digit,
// like the coordinates in the KML files are
// Example - Lat: 43° 38' 19.39" N
//                d  min  sec    -
// DecimalDigit := d + (min/60) + (sec/3600)
func ParseLatitude(tag string) float64 {
   a := regexp.MustCompile(`/1","*`)
   lat := a.Split(tag, 4)

   d := lat[0][2:4]
   dFloat, _ := strconv.ParseFloat(d, 64)

   min := lat[1]
   minFloat, _ := strconv.ParseFloat(min, 64)

   seconds := lat[2]
   sec := strings.SplitAfter(seconds, "/")

   divBy := sec[1][0:len(sec[1])-2]
   divByFloat, _ := strconv.ParseFloat(divBy, 64)

   s := strings.Split(sec[0], "/")
   secFloat, _ := strconv.ParseFloat(s[0], 64)
   secFloat = secFloat/divByFloat

   //log.Printf("Latidude is: %s\n", lat)
   //log.Printf("Latidude Decimal Digits - d: %f, min: %f, s: %f\n", dFloat, minFloat, secFloat)

   return dFloat + (minFloat/60) + (secFloat/3600)
} // }}}

// func ParseLongitude {{{
//
// Parses the Longitude value of the coordinate and converts it to a decimal digit,
// like the coordinates in the KML files are
// Example - Long: 116° 14' 28.86" W
//                  d  min  sec    -
// DecimalDigit := d + (min/60) + (sec/3600)
func ParseLongitude(tag string) float64 {
   a := regexp.MustCompile(`/1","*`)
   lng := a.Split(tag, 4)

   d := lng[0][2:5]
   dFloat, _ := strconv.ParseFloat(d, 64)

   min := lng[1]
   minFloat, _ := strconv.ParseFloat(min, 64)

   seconds := lng[2]
   sec := strings.SplitAfter(seconds, "/")

   divBy := sec[1][0:len(sec[1])-2]
   divByFloat, _ := strconv.ParseFloat(divBy, 64)

   s := strings.Split(sec[0], "/")
   secFloat, _ := strconv.ParseFloat(s[0], 64)
   secFloat = secFloat/divByFloat

   //log.Printf("Longitude is: %s\n", lng)
   //log.Printf("Longitude Decimal Digits - d: %f, min: %f, s: %f\n", dFloat, minFloat, secFloat)

   return dFloat + (minFloat/60) + (secFloat/3600)
} // }}}

// func LoopConduit {{{
//
// *CURRENTLY* Loops through all the images in the Conduit_Tie_Ins folder of my
// GoogleDrive, parsing its metadata for the coordinates (in decimal digits like the KMLs),
// if there are any. Makes it call distance to find a match ;-)
func LoopConduit(vz gjson.VZFeatureCollection, core gjson.CoreFeatureCollection, HUBS, WORK_ORDER_NAMES map[string]string) {
    var files  []string
    var names  []string
    var nometa []ImageFile
    var imgs   map[int]ImageFile

    imgs = make(map[int]ImageFile)

    root := "/Volumes/GoogleDrive/.shortcut-targets-by-id/1QVO3IzpFcRhx4tYXQbNVozex4gswTNe5/Conduit Tie ins"
    err := filepath.Walk(root, Visit(&files, &names, &nometa))
    if err != nil {
        log.Fatal(err)
    }

    i := 0
    n := 0
    for _, file := range files {
      LAT, LNG = 1.0, -1.0
      if file == root {
         // First time I ran this, it crashed because I called Walk on a directory
         continue
      }
      //fmt.Println(file)

      // Open the file
      f, err := os.Open(file)
      if err != nil {
         log.Fatal(err)
      }

      // Decode file for its metadata
      df, err := exif.Decode(f)
      if err != nil {
         log.Printf("Skipping file : %s, no metadata\n", file)
         continue
      }

      f.Close()

      // Go through the files metadata and parse for LAT & LNG
      var p Printer
      var img ImageFile
      df.Walk(p)

      // If the file had metadata and we correctly parsed and
      // converted it, LAT should be > 1.
      if LAT > 1.0 {
         img.Path = file
         img.Point.Lat = LAT
         img.Point.Lng = LNG
         img.ObjectId, img.FQNID, img.UDFs, img.NFID = MatchVZ(img.Point, vz)

         if img.ObjectId == -1 {
            log.Printf("Skipping file for now ..  no FQNID#. Image: %s\n", img.Path)
            continue
         }

         nfid := strings.Split(img.NFID, ".")
         img.WorkOrderName = "LSA_N_" + nfid[0] + "_" + nfid[1] + "_" + WORK_ORDER_NAMES[img.NFID]

         // Use NFID to get the hub & work order number
         s := nfid[0:2]
         first := s[0] + s[1]
         if first == "18" {
            img.Hub = HUBS[img.NFID]
         } else {
            img.Hub = HUBS[nfid[0]]
         }

         // We need the hub to know which folder to the copy the file to,
         // so if we have no hub, we can't copy. Let's see if we can match
         // the image to a core permit to find the hub.
         if img.Hub == "" {
            change := false
            // If we had a work order number using the FQNID's NFID,
            // then let's keep using that one. If we didn't, let's update
            // the work order name to match the Core NFID#
            if WORK_ORDER_NAMES[img.NFID] == "" {
               img.WorkOrderName = "LSA_N_" + nfid[0] + "_" + nfid[1] + "_"
               change = true
            }
            img.NFID = MatchCore(img.Point, core)
            nfid := strings.Split(img.NFID, ".")
            if len(nfid) > 1 {
               s := nfid[0:2]
               first := s[0] + s[1]
               if first == "18" {
                  img.Hub = HUBS[img.NFID]
               } else {
                  img.Hub = HUBS[nfid[0]]
               }
               if change {
                  img.WorkOrderName = img.WorkOrderName + WORK_ORDER_NAMES[img.NFID]
               }
            }
         }

         log.Printf("Image: %s\n", img.Path)
         log.Printf("Matching ObjectID: %d, FQNID: %s\n", img.ObjectId, img.FQNID)
         log.Printf("NFID: %s, Hub: %s\n", img.NFID, img.Hub)

         if img.Hub == "" {
            log.Printf("Skipping the copy for now .. no hub name. Image: %s\n", img.Path)
            continue
         } else if img.NFID == "" {
            log.Printf("Skipping the copy for now .. no NFID. Image: %s\n", img.Path)
            continue
         }

         // Perform the image copy.
         err := CopyImage(img)
         if err != nil {
            log.Printf("Error copying image: %s\n", img.Path)
         }

         // Add image to our log of copied images.
         imgs[i] = img
         i = i + 1
      }
      n = n + 1
   }
   MakeJSONFile(imgs)
} // }}}

// func MakeJSONFile {{{
//
func MakeJSONFile(imgs map[int]ImageFile) {
   fileName := "imgs.json"
   file, _ := json.MarshalIndent(imgs, "", "    ")
   _ = ioutil.WriteFile(fileName, file, 0644)
} // }}}

// func CopyImage {{{
func CopyImage(img ImageFile) error {
   // Open the original image
   src, err := os.Open(img.Path)
   if err != nil {
      return err
   }
   defer src.Close()

   // Create the new file path
   dir := "/Volumes/GoogleDrive/My Drive/InvoicePackages/" + img.Hub + "/" + img.NFID + "/" + img.FQNID + "/AsBuilts/Photos"
   err = os.MkdirAll(dir, 0755)
   log.Print(err)

   // Create the new image name
   path := img.Path[len(img.Path) - 4:len(img.Path)] // returns the file extension (.jpg)
   fqnid := strings.SplitAfter(img.FQNID, ":")

   fqnid[0] = fqnid[0][0:3]       // First 3 letters, like "FIB" or "OSP"
   if len(fqnid[1]) > 5 {         // This means our FQNID# is OSP:UGSPAN::
      fqnid[1] = fqnid[1][0:6]
   } else if len(fqnid[1]) == 5 { // This means our FQNID# is FIB:TAIL::
      fqnid[1] = fqnid[1][0:4]
   } else {
      fqnid[1] = fqnid[1][0:3]
   }

   var fName string
   var imgNum int64

   fName = img.WorkOrderName + "_" + fqnid[0] + "_" + fqnid[1] + "_" + fqnid[3] +"_CONDUIT_TIE_IN" + path
   imgNum = 2
   n := dir + "/" + fName

   // Check if this file name exists already or not
   exists := FileExists(n)
   if exists {
      // Okay so it does - let's add (2) and see if that works
      // if it doesn't let's keep adding to the image number
      // until we get one that does not exist.
      for {
         name := fName[0:len(fName)-4]
         n = dir + "/" + name + "(" + strconv.FormatInt(imgNum, 10) + ")" + path
         exists = FileExists(n)
         if !exists {

            break
         }
         imgNum = imgNum + 1
      }
   }
   fName = n
   img.NewFile = fName
   log.Printf("This is the new file name? %s\n", fName)

   // Create the destination file
   dst, err := os.Create(fName)
   if err != nil {
      return err
   }
   defer dst.Close()

   // Copy the contents of the src file to the dst file
   _, err = io.Copy(dst, src)
   if err != nil {
      return err
   }
   return nil
} // }}}

// func FileExists {{{
//
// Checks if a file exists or not
func FileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
} // }}}
