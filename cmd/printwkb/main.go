package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx"
	"github.com/terranodo/tegola"
	"github.com/terranodo/tegola/basic"
	"github.com/terranodo/tegola/mvt"
	"github.com/terranodo/tegola/wkb"
)

// 3857
var road1 = "\001\002\000\000\000\014\000\000\000\307N\300\223p\006j\301|&\344k'\224UAZ\000\237\361o\006j\301\316S\330\015o\224UA\371CjJo\006j\301\272\215\222Zw\224UA\300\375{Wm\006j\301wb\020E\201\224UA\311%RAk\006j\301\315\264\332\336\207\224UA\254\332\313\343h\006j\301\367\203|\217\215\224UA\200\337\333^f\006j\3012O{u\222\224UAF\3330\241c\006j\301w\220J@\226\224UAh\211d\333^\006j\301\277}\005{\232\224UAa\346@\354U\006j\301\261\037\022\236\237\224UA7y-\201T\006j\301\326c\252+\240\224UA\266u\367_R\006j\301\222\377\025\257\240\224UA"
var road2 = "\001\002\000\000\000$\000\000\000\262\327\207\322\376\006j\301\241\240+\257\331\224UA\321\027\035\012\374\006j\301\350\370\3377\335\224UA\313\375'\317\367\006j\301\365\201\312k\340\224UA\2600\337\377\364\006j\301\307i\234%\343\224UA\215\355h\205\352\006j\301\002|\337\332\355\224UA\032\243\317\324\345\006j\301\203\346~\034\363\224UA\327\217 \002\344\006j\301>\360>g\365\224UAv\354\275/\342\006j\301f5\214\002\370\224UA\237\340Ar\340\006j\3018)\306\015\373\224UA\245\212\332\315\336\006j\301\013pG\304\376\224UA\220g*Z\332\006j\301\010m\223X\012\225UA\260~N1\327\006j\301\326\302r\306\021\225UA&\034\3662\323\006j\301\250\216.;\032\225UA\213\304\006k\321\006j\301\217\017\005\365\034\225UA<\006\234\326\313\006j\301\034\273\223z#\225UA\337u\302\374\311\006j\301\235\235\303i%\225UA\351\362d\024\307\006j\301N\214\000\226'\225UAdlI\033\277\006j\301W5\007\030,\225UA\311\233\312*\272\006j\301\007\347\352/.\225UA\225(\346\363\266\006j\301\260\342?\030/\225UA\277\274\254\342\260\006j\301\212\313\232\0130\225UA\315e\252\307\254\006j\301\307\240\307\0250\225UA\251\242/\270\251\006j\301\237w\332\245/\225UA\274\0367Y\244\006j\301\253\371\376\236.\225UA\033\015\020\252\242\006j\301,eqN.\225UAQ\354\324\006\234\006j\301\320\337\274F-\225UA\225\026\230\321\224\006j\301-\355F\262+\225UA\307\271b\325\215\006j\301\242\022\345\214*\225UA\334\001\010\360\202\006j\301\323\363\247`(\225UA\234\343g\367{\006j\301S\200\371c'\225UA\206\324\333\362y\006j\301\035\205\345\364&\225UA\253<G\024u\006j\301#\217\317\307$\225UA\021:$zr\006j\301\232\351fp#\225UA=U[@i\006j\301$6\263\361\035\225UA\316V[tf\006j\301\324\307jg\034\225UA\203\324\3528a\006j\301\366\363\273\213\032\225UA"
var tile1 = tegola.BoundingBox{
	Minx: -1.3644926791343994e+07,
	Miny: 5.657563084768066e+06,
	Maxx: -1.3644315295117797e+07,
	Maxy: 5.65695158854187e+06,
}
var water1 = "\001\002\000\000\000\004\000\000\000\246o\223Y%\035\017A\262f\023]\225@TA\306a\235E)\034\017A\025*\017\\\237@TA\374\267\2702\002\033\017A0R\013Q\243@TA\351\311\341\020\034\032\017A\366\226\"\373\244@TA"
var water2 = "\001\002\000\000\0006\000\000\000\351\311\341\020\034\032\017A\366\226\"\373\244@TA\215}\250OX\030\017A(>\031!\251@TA\022\306\031\340]\027\017A\230\313\263>\257@TA\\\200\241\305\321\024\017AK\300\020V\256@TA\372\023\335\031\252\023\017AS\330\217\366\255@TA\225\376\336p\351\022\017A<S\017\246\256@TA\300\217Ia\272\021\017A5\005\325H\270@TA\224\027\311\365:\021\017AS3\\\345\272@TAD'\375\307\210\020\017A\002\323\004\216\272@TA<Z\020\357k\017\017A8\204\257\342\264@TA\325\230\177\025E\016\017A\243\020fk\261@TA\277\330}z4\015\017AA<\033\370\262@TA\262\330\001\343Z\014\017A\213\370\314\361\270@TAK\246\251s\237\013\017A\207\027\211\267\267@TA\025;\244\346>\012\017A\000\017;\233\250@TA\\M\301\214l\011\017A\215\214\251)\245@TA\203\"\001S\210\007\017A\207\247\261\325\247@TA_\261\015G\360\004\017A\315D\355\322\236@TA\367Kni\001\004\017Ax\225z\262\234@TA\262\321WS)\377\016A,]\216\013\240@TA7\304uHf\374\016A\321F\303\277\252@TAz\352\011*\213\372\016A\346F\214\276\256@TA\330\341\362W7\370\016A\000]C\215\255@TAQ\"d$\305\365\016A\203\254\274_\243@TA\015\014\037\215w\364\016A3n\025v\244@TA}\240%8\003\363\016A~r\017\032\252@TA\204\361\021\305k\360\016A\221\374\213\365\257@TA\003\250\031\307\234\357\016A\216\353U\206\261@TAW/\326\325\266\355\016A\201$\244\226\251@TA\271\211\254\\s\353\016AA\303\223\002\240@TANX\177\004\315\351\016A\320g\030\232\237@TA\234\027z3\233\346\016A\274\343\326Z\243@TA\300\202_\202\344\344\016A\217\311\364U\253@TA\246c\203\027\274\343\016A\371'\260\300\255@TAAIFo\377\337\016A\030 \312f\252@TAP\264\3275\323\336\016A\310\214_\273\251@TA\372L\272\032V\336\016A\205iz\033\244@TA\271\002\003\376\277\335\016A\224.\305\305\243@TA$\356\207[)\335\016A\232\036u3\245@TARX\3223p\333\016Ay \344\356\244@TA\211\241>\373\331\331\016A\001`\330\030\256@TA\200\3635^2\330\016A\275\004r@\262@TA\015\376n\027`\327\016A\006\224\337\316\256@TA\325\225\302\275\314\326\016A\005\310\241\233\244@TA@/d\301\026\326\016A\000\337\305p\240@TA\015\254w\267\025\324\016A+T s\245@TA}E\251[r\322\016Ab\005c9\253@TAN\013\251\027T\321\016A\375\335\313%\253@TA\\Z\233\317T\316\016A:\252P\365\241@TA\\\220\337\032\241\314\016A\007\314n\012\237@TA^Py'\021\313\016AOF\246\234\242@TA\306D\365\030\267\311\016Am\275v\342\256@TAfD]\252e\307\016A\336_$\244\265@TA=\015\216\212\260\305\016A\277\314\3633\301@TA"
var water3 = "\001\002\000\000\000s\000\000\000\242\330tJ!\235\017A\202\237\245\323\342ATA5\020C\013\350\232\017A\"\233\313}\367ATA\262;n\336+\230\017A)\036C\\\005BTAz\317QC\357\225\017A.:\275\306\006BTAa9\264\213\016\224\017A\001m\263\337\005BTA.\344\233\246\346\222\017A\014\002\336e\373ATA=\034.O\020\221\017A\220E\342,\362ATAu\232k\372'\217\017A\246\227\215i\356ATA\256\217\262B\201\215\017Am\332\271\310\355ATA\300\034\021\031\252\213\017A\315\255,\256\347ATA\357\022\363\327\005\212\017A\212\220Vn\335ATA\"Ll\256Z\211\017A\202\220\321\256\321ATA\365\326Dv\250\211\017AI\255r8\305ATA:\2350)Q\212\017A\357\276\"\231\272ATA\177#n\230\340\211\017AE\022z\311\254ATA\300\345\241\271V\211\017A\000#\207\311\240ATA\215%\274\365\275\211\017A\207\306\307?\222ATAk\316\335\202~\211\017An\036\324\311\206ATA\232\004n\205\363\210\017A\024\277Qx\177ATAX\355\265\206\007\211\017A\000]\011\363qATAN\314\023\216\307\210\017A\011\207\362\325hATAG;e\267\032\207\017A\244\241\276\236_ATA\351\363\016\237\236\205\017A\320\310\242SZATA\306)\203\312\233\203\017A# \231\015]ATA\341\256Ik\355\202\017AF\024\034\015^ATA0U\242e\363\200\017AD'\244\263^ATA\225\315\351\300\206\177\017A\316m\367\323]ATA\354\023\205$\351}\017A\267\342\247\322YATA\013\243\351\007f|\017A_y\234\224OATA\007\276\025\262\250z\017A\232\031o]FATA<\216\356fVy\017A\302\263\371\203?ATAd\\\036\307`y\017A\330^B\2677ATA\260\023\251\317\273x\017A\335\343YK4ATAr\334\205'\221w\017A\007\326\360\3764ATA\242&\035\006\020u\017A\261l\004\372>ATA\372\347\325\361\250r\017A\253a\177\212DATA\357\213h^\227q\017A\341\237\2315DATAy\233\204\300\014q\017Az\240\327\022;ATA\200{agWo\017AM\376\315\3452ATA\233\250J\367em\017A3\3216\2022ATA\257\374\013Oul\017A\344\226\263:3ATA&\220/3uk\017A\224\237$\206/ATA\337y\242K\236j\017Am='\336,ATA&\350\370=\215i\017A\365<4t*ATA\207+'e\351h\017A\303\326\003\337\"ATA\001]\004\237\307g\017A\234s\020=!ATA\332\034\366\232Lg\017A\336q\220\276\033ATA\011\254\036f,f\017A\257UXc\024ATA5e\253\024<d\017A0\205z\326\017ATAu\223{\322uc\017A\353\303\022\253\014ATA\307M\357\256\260b\017A\332W\204\230\005ATA\306\276\266DEb\017A\240\360.\001\004ATAR\346I\304ma\017A\206n}\362\003ATA++\257:\313_\017A\216R\301\365\003ATA\350\226\363\347y^\017AA\001\016\302\375@TA\216\203\373K\015^\017A\232\326o\342\367@TA\336\276>\027\031^\017A\343\271C!\353@TA&\007\374^G^\017A\367\003\252\206\334@TA\215SUa0^\017A@E\3567\323@TAJ\312\276\276\323]\017AP]0\315\320@TA\320\015\005V\010]\017A \350O\277\320@TAs\272\026&\371[\017A\371\355\353\361\320@TA\004\352\275r\322Z\017Ag2\342\260\311@TA\313\226\237b\022Z\017A[\274+\003\277@TA\256^5\314ZX\017A\001\010q\351\272@TAW\364v-\207W\017AR\024\177\017\272@TA\2467\341l\200V\017A\275)\333\205\272@TA\230Y\242D\351U\017A\002uU+\265@TA\344\303,\234lU\017A\011B\034\326\253@TA\304\340\376\222\241T\017Ak\217\331s\252@TAG\005\005\230\254Q\017A\377\240@'\265@TA\344\271\256\235tO\017A,\307\205\021\266@TA2\245\363{\241M\017A\032g\251\357\271@TA\332\262\261\341\255K\017A)&\217\275\270@TA\365\253\340\233\352J\017A<\205\342\004\272@TA\036\277\214\200.J\017A/\2730Y\300@TA\324\347\226\316\224I\017A(+>\323\304@TAM)\327XuH\017A,\274B\\\302@TAP\357j\034\011H\017A\254D\272\343\272@TA\254s\272\356\034G\017A\254t\340\327\266@TAIO[\236\034@\017ATs\3119\247@TA^Q9\234\014?\017A\360\313GW\252@TAI\234\377\370?=\017A\276c(\353\265@TAQ\234\327\351\333;\017A\234\261\374Z\266@TAw\231\226&\347:\017A\270\216\022\222\262@TAT\012\216\234,:\017A?VT\015\263@TA\316\031\346\026?9\017Ae\200\243\015\264@TA\346\245\031\326R8\017A \223\212E\260@TA\267i\217\034\2727\017A\334L\320\007\261@TA%\311\304\343\2276\017A\320,\347\354\270@TA\020rN\002\0074\017A\205.\027\177\311@TA6j\342EA3\017A\200\312\261\022\324@TA\007]\214q\3432\017A\264c1q\326@TAN\030\233\304 2\017A\326\277\017\230\325@TAy\317U\255,1\017A\274;\354j\317@TA\301\257\212C\2140\017A5\307H\373\314@TA\357\335N\311\331/\017A\001\312P\020\317@TA\212\002\005\255\016/\017A}c\2155\316@TAN\013\330\230\206.\017A\372\006H\011\320@TA_\362\235\365l-\017A\240\225S\"\327@TA\001\237gu\324,\017A\337\365\207\030\327@TA!G\0224W,\017AE\330{'\320@TA\2741\024\213\226+\017AN#\271!\310@TAG\220\272\023^*\017A\265\324\225\313\303@TA\227\017cU\307(\017A\256-f\364\303@TA\362\351\251l\357&\017A\011\302+\377\271@TARIc\234\363$\017A\031\343n\274\267@TA\024\332\025\335M#\017A\201\275\000\267\260@TA\342:\346\032.\"\017Az\376F\013\257@TA\303\256\306dS!\017A$)\035h\251@TAJ\235\320n\233\037\017AD\035%\241\246@TA\262\304K~\353\036\017A]\242\215$\237@TAZ\356\250\362\312\036\017Au\246$\025\232@TA\271\275n\000<\036\017A\026kt\312\225@TA\246o\223Y%\035\017A\262f\023]\225@TA"
var tile2 = tegola.BoundingBox{
	Minx: 254382.43009765446,
	Miny: 5310233.228288574,
	Maxx: 256828.4150024396,
	Maxy: 5307787.243383789,
}

var geomField int

func print(srid int, geostr string, tile tegola.BoundingBox) {

	rd1 := strings.NewReader(geostr)
	geo, err := wkb.Decode(rd1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding goe(road1): %v", err)
		os.Exit(2)
	}
	gwm, err := basic.ToWebMercator(srid, geo)
	geol := gwm.AsMultiPolygon()
	c := mvt.NewCursor(tile, 4096)
	g := c.ScaleGeo(geol)
	cmin, cmax := c.MinMax()
	fmt.Printf("Rec: %v,%v\n", cmin, cmax)
	fmt.Printf("Scle GEO: %T: %#[1]v\n", g)
	cg, err := c.ClipGeo(g)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Clip GEO:%T: %#[1]v\n", cg)

}

func main() {

	/*
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer file.Close()
		csvFile := csv.NewReader(file)
		// We are going to assume the first record is the names of the columns. We only care for the geometry field
		fieldNames, err := csvFile.ReadAll()
		if err != nil {
			panic(err)
		}
	*/
	if len(os.Args) < 5 {
		panic("Need x, y and z values.")

	}
	srid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		srid = 4326
	}
	z, _ := strconv.Atoi(os.Args[2])
	x, _ := strconv.Atoi(os.Args[3])
	y, _ := strconv.Atoi(os.Args[4])
	tile := tegola.Tile{
		X: x,
		Y: y,
		Z: z,
	}
	bbox := tile.BoundingBox()
	minGeo, err := basic.FromWebMercator(srid, &basic.Point{bbox.Minx, bbox.Miny})
	if err != nil {
		panic(err)
	}
	maxGeo, err := basic.FromWebMercator(srid, &basic.Point{bbox.Maxx, bbox.Maxy})
	if err != nil {
		panic(err)
	}
	minPt := minGeo.AsPoint()
	maxPt := maxGeo.AsPoint()
	config := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "bonn_osm",
			User:     "gdey",
			Password: "redfox",
		},
	}

	pool, err := pgx.NewConnPool(config)
	if err != nil {
		panic(fmt.Sprintf("Failed while creating connection pool: %v", err))
	}
	/*
		var bbox = tegola.BoundingBox{
			Minx: 7.097167967762075,
			Miny: 50.75035930616591,
			Maxx: 7.119140624009016,
			Maxy: 50.736455131807304,
		}
	*/
	sql := fmt.Sprintf(
		`SELECT  ST_AsBinary("wkb_geometry") AS "geometry" from forest WHERE "wkb_geometry" && ST_MakeEnvelope(%v,%v,%v,%v,%v)`,
		minPt.X(),
		minPt.Y(),
		maxPt.X(),
		maxPt.Y(),
		srid,
	)
	fmt.Println("Running:", sql)
	rows, err := pool.Query(sql)
	if err != nil {
		panic(fmt.Sprintf("Got the following error (%v) running this sql (%v)", err, sql))
	}
	defer rows.Close()
	//	fetch rows FieldDescriptions. this gives us the OID for the data types returned to aid in decoding
	var geobytes []byte
	var ok bool
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			panic(fmt.Sprintf("Got an error trying to run SQL: %v ; %v", sql, err))
		}
		println("Vals:", vals)

		if geobytes, ok = vals[0].([]byte); !ok {
			panic("Was unable to convert geometry field into bytes.")
		}
		print(srid, string(geobytes), bbox)
	}

	//print(water3, tile2)

}
