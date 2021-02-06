// Package for handling bore log files
package borelogs

import (
   "archive/zip"
    "fmt"
    "log"
    "os"
    "io"
    "path/filepath"
    "encoding/json"
    "io/ioutil"
    "regexp"
    "strings"
    "strconv"
    "vci/kml/corepermits"
    "vci/kml/vz3gis"
    "vci/kml/distance"
    "github.com/rwcarlsen/goexif/exif"
    "github.com/rwcarlsen/goexif/tiff"
)

type ImageFile struct {
   Name           string
   Path           string
   Hub            string
   ObjectId       int64
   FQNID          string
   NFID           string
   WorkOrderName  string
   UDFs           []string
   Point          corepermits.Point
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
         return nil
      case ".zip":
         *files = append(*files, path)
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
      if tag.String() == "E" {
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

   log.Printf("Latidude is: %s\n", lat)
   log.Printf("Latidude Decimal Digits - d: %f, min: %f, s: %f\n", dFloat, minFloat, secFloat)

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

   log.Printf("Longitude is: %s\n", lng)
   log.Printf("Longitude Decimal Digits - d: %f, min: %f, s: %f\n", dFloat, minFloat, secFloat)

   return dFloat + (minFloat/60) + (secFloat/3600)
} // }}}

// func LoopConduit {{{
//
// *CURRENTLY* Loops through all the images in the Conduit_Tie_Ins folder of my
// GoogleDrive, parsing its metadata for the coordinates (in decimal digits like the KMLs),
// if there are any. Makes it call distance to find a match ;-)
func LoopBoreLogs(vz [][]vz3gis.FQNID, core [][]corepermits.Permit) {
    var files  []string
    var names  []string
    var nometa []ImageFile
    var imgs   map[int]ImageFile

    imgs = make(map[int]ImageFile)

    root := "/Volumes/GoogleDrive/My Drive/CX File Dump/BorelogsTestFiles"
    err := filepath.Walk(root, Visit(&files, &names, &nometa))

    if err != nil {
        panic(err)
    }

    i := 0
    for _, file := range files {
      if file == root {
         // First time I ran this, it crashed because I called Walk on a directory
         continue
      }

      path := file[len(file) - 4:len(file)]
      path = strings.ToLower(path)

      if path == ".zip" {
         dst := file[0:len(file)-4]
         log.Printf("Dst: %s", dst)
         blimgs, err := Unzip(file, dst)
         if err != nil {
            log.Print(err)
            continue
         }
         for _, img := range blimgs {
            imgs, i = OpenImage(img, imgs, i, vz, core)
         }
      } else if path == ".jpg" {
         imgs, i = OpenImage(file, imgs, i, vz, core)
      } else if path == ".png" {
         imgs, i = OpenImage(file, imgs, i, vz, core)
      }
      //fmt.Println(file)
   }
   MakeJSONFile(imgs)
} // }}}

// func OpenImage() {{{
func OpenImage(file string, imgs map[int]ImageFile, i int, vz [][]vz3gis.FQNID, core [][]corepermits.Permit) (map[int]ImageFile, int) {
   LAT, LNG = 1.0, -1.0

   // Open the file
   f, err := os.Open(file)
   if err != nil {
      log.Print(err)
      return imgs, i
   }

   // Decode it for the metadata
   df, err := exif.Decode(f)
   if err != nil {
      log.Printf("Skipping file : %s, no metadata\n", file)
      return imgs, i
   }

   f.Close()

   var p Printer
   var img ImageFile
   df.Walk(p)

   if LAT > 1.0 {
      img.Path = file
      img.Point.Lat = LAT
      img.Point.Lng = LNG
      img.ObjectId, img.FQNID, img.WorkOrderName, img.UDFs, img.NFID, img.Hub = distance.MatchCoords(img.Point, vz, core)
      if img.ObjectId == -1 {
         return imgs, i
      }

      log.Printf("Image: %s\n", img.Path)
      log.Printf("Matching ObjectID: %d, FQNID: %s\n", img.ObjectId, img.FQNID)

      if img.Hub == "" {
         log.Printf("Skipping the copy for now .. no hub name. Image: %s\n", img.Path)
         return imgs, i
      } else if img.NFID == "" {
         log.Printf("Skipping the copy for now .. no NFID. Image: %s\n", img.Path)
         return imgs, i
      }
      err := CopyImage(img)
      if err != nil {
         log.Printf("Error copying image: %s\n", img.Path)
      }
      imgs[i] = img
      i = i + 1
   }
   return imgs, i
}// }}}

// func MakeJSONFile {{{
//
func MakeJSONFile(imgs map[int]ImageFile) {
   fileName := "imgs.json"
   file, _ := json.MarshalIndent(imgs, "", "    ")
   _ = ioutil.WriteFile(fileName, file, 0644)
} // }}}

// func CopyImage {{{
func CopyImage(img ImageFile) error {
   src, err := os.Open(img.Path)
   if err != nil {
      return err
   }

   defer src.Close()

   path := img.Path[len(img.Path) - 4:len(img.Path)]
   fqnid := strings.SplitAfter(img.FQNID, ":")

   fqnid[0] = fqnid[0][0:3]
   fqnid[1] = fqnid[1][0:3]

   dir := "/Volumes/GoogleDrive/My Drive/InvoicePackages" + img.Hub + "/" + img.NFID + "/" + img.FQNID + "/Bore Logs/Photos"
   err = os.MkdirAll(dir, 0755)
   log.Print(err)

   var fName string
   var imgNum int64
   if img.WorkOrderName == "LSA_N_DATA_CORRECTION" {
      nfid := strings.Split(img.NFID, ".")
      fName =  nfid[0] + "_" + nfid[1] + "_FIB_BUR_" + fqnid[3] +"_BORE_LOG_PIC" + path
      imgNum = 2
      n := dir + "/" + fName
      exists := fileExists(n)
      if exists {
         for {
            name := fName[0:len(fName)-4]
            n = name + "(" + strconv.FormatInt(imgNum, 10) + ")" + path
            exists = fileExists(n)
            if !exists {
               fName = n
               break
            }
            imgNum = imgNum + 1
         }
      }
   } else {
      fName = img.WorkOrderName + "_FIB_BUR_" + fqnid[3] +"_BORE_LOG_PIC" + path
      imgNum = 2
      n := dir + "/" + fName
      exists := fileExists(n)
      if exists {
         for {
            name := fName[0:len(fName)-4]
            n = name + "(" + strconv.FormatInt(imgNum, 10) + ")" + path
            exists = fileExists(n)
            if !exists {
               fName = n
               break
            }
            imgNum = imgNum + 1
         }
      }

   }
   log.Printf("This is the new file name? %s\n", fName)

   newFile := dir + "/" + fName

   dst, err := os.Create(newFile)
   if err != nil {
      return err
   }
   defer dst.Close()

   _, err = io.Copy(dst, src)
   if err != nil {
      return err
   }
   return nil
} // }}}

// func fileExists {{{
// src: https://golangcode.com/check-if-a-file-exists/
// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
} // }}}

// func Unzip {{{
// src: https://golangcode.com/unzip-files-in-go/
//
// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {
    var filenames []string

    r, err := zip.OpenReader(src)
    if err != nil {
        return filenames, err
    }
    defer r.Close()

    for _, f := range r.File {

        // Store filename/path for returning and using later on
        fpath := filepath.Join(dest, f.Name)

        // Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return filenames, fmt.Errorf("%s: illegal file path", fpath)
        }

        filenames = append(filenames, fpath)

        if f.FileInfo().IsDir() {
            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        // Make File
        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return filenames, err
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return filenames, err
        }

        rc, err := f.Open()
        if err != nil {
            return filenames, err
        }

        _, err = io.Copy(outFile, rc)

        // Close the file without defer to close before next iteration of loop
        outFile.Close()
        rc.Close()

        if err != nil {
            return filenames, err
        }
    }
    return filenames, nil
} // }}}
