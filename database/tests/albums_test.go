package tests

import "testing"
import "bytes"
import "os"
import "encoding/json"
import "github.com/SharedCode/parallels/database"
import "github.com/SharedCode/parallels/database/repository"

type AlbumType string

const (
	Album = "album"
	Checkpoint = "checkpoint"
)

func TestBasic(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var config, _ = database.LoadConfiguration(dir + "/config.json")
	repoSet, e := database.NewRepositorySet(config)
	if e != nil {
		t.Error(e)
	}

	// insert Album XML w/ ID 1665078 to Store
	albumID := "1665078"
	rs := repoSet.Store.Set(*repository.NewKeyValue(Album, albumID, albumXML))
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
	albumID2 := "1665079"
	rs = repoSet.Store.Set(*repository.NewKeyValue(Album, albumID2, albumXML))
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
   }
   
   ba, e := json.Marshal([]string{albumID, albumID2})
 	if e != nil {
 		t.Error(e)
    }
    
	// insert checkpoint 1 containing AlbumIDs inserted to the Entity Store.
	rs = repoSet.NavigableStore.Set(*repository.NewKeyValue(Checkpoint, "1", ba))
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}

	gr,rs := repoSet.NavigableStore.Get(Checkpoint, "1")
	if !rs.IsSuccessful() {
		t.Error(e)
   }

   var albumIDs []string
	e = json.Unmarshal(gr[0].Value, &albumIDs)
	if e != nil {
		t.Error(e)
	}
	gr,rs = repoSet.Store.Get(Album, albumIDs...)
	r := gr
	if !rs.IsSuccessful() {
		t.Error(rs.Error)
	}
	if r == nil || len(r) != 2 {
		t.Errorf("Expected 2 Albums not read.")
	}
	if r[0].Key != albumID && r[0].Key != albumID2 {
		t.Errorf("Expected 2 Albums' keys not read.")
	}
	if r[1].Key != albumID && r[1].Key != albumID2 {
		t.Errorf("Expected 2 Albums' keys not read.")
	}
   if !bytes.Equal(r[0].Value, r[1].Value) {
		t.Errorf("Expected AlbumXML retrieved from DB did not match.")
	}
}

var albumXML = []byte(`
<ALBUM ID="1665078" LANGUAGE_ID="44">
<REVISION>
   <LEVEL>7</LEVEL>
   <TAG>94746649EC653E360EEF1A9A4E6FCA08</TAG>
</REVISION>
<QUALITY>
   <LEVEL>100</LEVEL>
</QUALITY>
<COMPILATION>0</COMPILATION>
<GENRE_ID VERSION="1" ORDINAL="1">175</GENRE_ID>
<GENRE_ID VERSION="2" PRIMARY="Y">19073</GENRE_ID>
<TITLE>
   <DISPLAY>透明</DISPLAY>
</TITLE>
<REQUISITION_EXISTS>Y</REQUISITION_EXISTS>
<LOOKCNT>316</LOOKCNT>
<ARTIST ID="447651" PRIMARY="1">
   <NAME>
      <DISPLAY>梁咏琪</DISPLAY>
      <FIRSTNAME>梁咏琪</FIRSTNAME>
   </NAME>
   <ARTISTTYPE_ID>2656</ARTISTTYPE_ID>
   <ERA_ID PRIMARY="N">2648</ERA_ID>
   <ERA_ID PRIMARY="Y">2650</ERA_ID>
   <ORIGIN_ID>3984</ORIGIN_ID>
</ARTIST>
<TOC ID="2110862">
   <OFFSETS>150 22054 36278 55024 74329 91714 111740 127565 144231 159976 179770 201314 215878 233769 255313 275181 288333 309047 328449 347163</OFFSETS>
   <MEDIAID>664D9CCEA36C4F815B6B2B0693A1F339</MEDIAID>
   <TUI>
      <ID>30681883</ID>
      <TAG>F249A3004BAD8D7D071BE176891DC1F7</TAG>
   </TUI>
</TOC>
<TOC ID="3803172">
   <OFFSETS>150 22054 36278 55024 74329 91714 111740 127565 144231 159976 179770 201314 215878 233769 255313 275181 288333 309047 328449 347164</OFFSETS>
   <MEDIAID>08861F572FDD6368686B1C02B9A069AD</MEDIAID>
   <TUI>
      <ID>54355581</ID>
      <TAG>1601DE59DFD4C04A239A3EBF24E92C42</TAG>
   </TUI>
</TOC>
<TRACKS COUNT="19">
   <TRACK ID="21266227" ORDINAL="1">
      <TITLE>
         <DISPLAY>01hshui</DISPLAY>
      </TITLE>
      <TUI TOC_ID="2110862">
         <ID>30681884</ID>
         <TAG>6BFF0B2F45D6869B5D4D63E5A3C670CA</TAG>
      </TUI>
      <TUI TOC_ID="3803172">
         <ID>54355582</ID>
         <TAG>118E72C79D8A5362D736A3F53F2B1263</TAG>
      </TUI>
   </TRACK>
   <TRACK ID="21266228" ORDINAL="2">
      <TITLE>
         <DISPLAY>02tming</DISPLAY>
      </TITLE>
      <TUI TOC_ID="2110862">
         <ID>30681885</ID>
         <TAG>E08E1DC4D59335F77419C6E0420DDB2B</TAG>
      </TUI>
      <TUI TOC_ID="3803172">
         <ID>54355583</ID>
         <TAG>259908297FE1E34198657CB8CC612EAA</TAG>
      </TUI>
   </TRACK>
   <TRACK ID="21266229" ORDINAL="3">
      <TITLE>
         <DISPLAY>03lguang</DISPLAY>
      </TITLE>
      <TUI TOC_ID="2110862">
         <ID>30681886</ID>
         <TAG>03732EA48461BA5B06E6DBFA434AD115</TAG>
      </TUI>
      <TUI TOC_ID="3803172">
         <ID>54355584</ID>
         <TAG>6280C6053452F379B7AB242752C10B34</TAG>
      </TUI>
   </TRACK>
</TRACKS>
</ALBUM>`)
