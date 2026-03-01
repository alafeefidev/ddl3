package main

import (
	// "encoding/xml"
	"fmt"
	// "os"

	"github.com/alafeefidev/ddl3/internal/parser"
)

func main() {
	// name := "C:\\Users\\Admin\\Documents\\code\\ddl3\\internal\\parser\\test\\THEEQUALIZER3Y2023M.mpd"
	// file, err := os.ReadFile(name)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var mpd smbd.Mpd
	// xml.Unmarshal(file, &mpd)

	mpd, err := parser.LoadFromUri(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	mpd.ResolveUrls(url)

	fmt.Println(mpd.Periods[0].AdaptationSets[0].Representations[0].ResolvedURL)

}
