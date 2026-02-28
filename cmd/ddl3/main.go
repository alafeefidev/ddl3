package main

import (
	"encoding/xml"
	"fmt"
	"os"

	smbd "github.com/alafeefidev/ddl3/internal/parser/mpd"
)

func main() {
	name := "C:\\Users\\Admin\\Documents\\code\\ddl3\\internal\\parser\\test\\THEEQUALIZER3Y2023M.mpd"
	file, err := os.ReadFile(name)
	if err != nil {
		fmt.Println(err)
	}
	var mpd smbd.Mpd
	xml.Unmarshal(file, &mpd)

	fmt.Println(*mpd.BaseUrl)

}
