package format

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type tInfo struct {
	fullName  string
	nameStart string
	nameEnd   string
}

var (
	empty = ""

	tList = []tInfo{
		// Random entries test:
		{fullName: "AAAAA LIMITED", nameStart: "AAAAA", nameEnd: "LIMITED"},
		{fullName: "AAAAA LTD", nameStart: "AAAAA", nameEnd: "LTD"},
		{fullName: "AAAAA 456 LIMITED  &^£@", nameStart: "AAAAA 456", nameEnd: "LIMITED  &^£@"},
		{fullName: "AAAAA 456     LIMITED  &^£@", nameStart: "AAAAA 456", nameEnd: "LIMITED  &^£@"},
		{fullName: "AAAAA 456 78AZ LIMITED .", nameStart: "AAAAA 456 78AZ", nameEnd: "LIMITED ."},
		{fullName: "AAAAA 456 78AZ    LIMITED.", nameStart: "AAAAA 456 78AZ", nameEnd: "LIMITED."},
		{fullName: "AAAAA 456 78AZ    LIMITED&", nameStart: "AAAAA 456 78AZ", nameEnd: "LIMITED&"},
		{fullName: "AAAAA --  LIMITED      ", nameStart: "AAAAA --", nameEnd: "LIMITED      "},
		{fullName: "AAAAA-----LIMITED  ", nameStart: "AAAAA-----LIMITED  ", nameEnd: empty},
		{fullName: "THE GIRLS' PUBLIC DAY SCHOOL TRUST (1872)", nameStart: "THE GIRLS' PUBLIC DAY SCHOOL TRUST (1872)", nameEnd: empty},
		{fullName: "THE GIRLS' PUBLIC DAY SCHOOL TRUST (1872) LTD", nameStart: "THE GIRLS' PUBLIC DAY SCHOOL TRUST (1872)", nameEnd: "LTD"},

		// Table entries test:
		{fullName: "AAAAA  123  (BBB) AEIE", nameStart: "AAAAA  123  (BBB)", nameEnd: "AEIE"},
		{fullName: "AAAAA  123  (BBB) ANGHYFYNGEDIG", nameStart: "AAAAA  123  (BBB)", nameEnd: "ANGHYFYNGEDIG"},
		{fullName: "AAAAA  123  (BBB) C.B.C", nameStart: "AAAAA  123  (BBB)", nameEnd: "C.B.C"},
		{fullName: "AAAAA  123  (BBB) C.C.C", nameStart: "AAAAA  123  (BBB)", nameEnd: "C.C.C"},
		{fullName: "AAAAA  123  (BBB) C.I.C", nameStart: "AAAAA  123  (BBB)", nameEnd: "C.I.C"},
		{fullName: "AAAAA  123  (BBB) CBC", nameStart: "AAAAA  123  (BBB)", nameEnd: "CBC"},
		{fullName: "AAAAA  123  (BBB) CBCN", nameStart: "AAAAA  123  (BBB)", nameEnd: "CBCN"},
		{fullName: "AAAAA  123  (BBB) CBP", nameStart: "AAAAA  123  (BBB)", nameEnd: "CBP"},
		{fullName: "AAAAA  123  (BBB) CCC", nameStart: "AAAAA  123  (BBB)", nameEnd: "CCC"},
		{fullName: "AAAAA  123  (BBB) CCG CYF", nameStart: "AAAAA  123  (BBB)", nameEnd: "CCG CYF"},
		{fullName: "AAAAA  123  (BBB) CCG CYFYNGEDIG", nameStart: "AAAAA  123  (BBB)", nameEnd: "CCG CYFYNGEDIG"},
		{fullName: "AAAAA  123  (BBB) CIC", nameStart: "AAAAA  123  (BBB)", nameEnd: "CIC"},
		{fullName: "AAAAA  123  (BBB) COMMUNITY INTEREST COMPANY", nameStart: "AAAAA  123  (BBB)", nameEnd: "COMMUNITY INTEREST COMPANY"},
		{fullName: "AAAAA  123  (BBB) COMMUNITY INTEREST P.L.C", nameStart: "AAAAA  123  (BBB)", nameEnd: "COMMUNITY INTEREST P.L.C"},
		{fullName: "AAAAA  123  (BBB) COMMUNITY INTEREST PLC", nameStart: "AAAAA  123  (BBB)", nameEnd: "COMMUNITY INTEREST PLC"},
		{fullName: "AAAAA  123  (BBB) COMMUNITY INTEREST PUBLIC LIMITED COMPANY", nameStart: "AAAAA  123  (BBB)", nameEnd: "COMMUNITY INTEREST PUBLIC LIMITED COMPANY"},
		{fullName: "AAAAA  123  (BBB) CWMNI BUDDIANT C.C.C", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI BUDDIANT C.C.C"},
		{fullName: "AAAAA  123  (BBB) CWMNI BUDDIANT CCC", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI BUDDIANT CCC"},
		{fullName: "AAAAA  123  (BBB) CWMNI BUDDIANT CYMUNEDOL C.C.C", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI BUDDIANT CYMUNEDOL C.C.C"},
		{fullName: "AAAAA  123  (BBB) CWMNI BUDDIANT CYMUNEDOL CCC", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI BUDDIANT CYMUNEDOL CCC"},
		{fullName: "AAAAA  123  (BBB) CWMNI BUDDIANT CYMUNEDOL CYHOEDDUS CYFYNGEDIG", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI BUDDIANT CYMUNEDOL CYHOEDDUS CYFYNGEDIG"},
		{fullName: "AAAAA  123  (BBB) CWMNI BUDDIANT CYMUNEDOL", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI BUDDIANT CYMUNEDOL"},
		{fullName: "AAAAA  123  (BBB) CWMNI BUDDSODDIA CHYFALAF NEWIDIOL", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI BUDDSODDIA CHYFALAF NEWIDIOL"},
		{fullName: "AAAAA  123  (BBB) CWMNI BUDDSODDIANT PENAGORED", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI BUDDSODDIANT PENAGORED"},
		{fullName: "AAAAA  123  (BBB) CWMNI CELL GWARCHODEDIG", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI CELL GWARCHODEDIG"},
		{fullName: "AAAAA  123  (BBB) CWMNI CYFYNGEDIG CYHOEDDUS", nameStart: "AAAAA  123  (BBB)", nameEnd: "CWMNI CYFYNGEDIG CYHOEDDUS"},
		{fullName: "AAAAA  123  (BBB) CYF", nameStart: "AAAAA  123  (BBB)", nameEnd: "CYF"},
		{fullName: "AAAAA  123  (BBB) CYFYNGEDIG", nameStart: "AAAAA  123  (BBB)", nameEnd: "CYFYNGEDIG"},
		{fullName: "AAAAA  123  (BBB) EEIG", nameStart: "AAAAA  123  (BBB)", nameEnd: "EEIG"},
		{fullName: "AAAAA  123  (BBB) EESV", nameStart: "AAAAA  123  (BBB)", nameEnd: "EESV"},
		{fullName: "AAAAA  123  (BBB) EOFG", nameStart: "AAAAA  123  (BBB)", nameEnd: "EOFG"},
		{fullName: "AAAAA  123  (BBB) EOOS", nameStart: "AAAAA  123  (BBB)", nameEnd: "EOOS"},
		{fullName: "AAAAA  123  (BBB) EUROPEAN ECONOMIC INTEREST GROUPING", nameStart: "AAAAA  123  (BBB)", nameEnd: "EUROPEAN ECONOMIC INTEREST GROUPING"},
		{fullName: "AAAAA  123  (BBB) GEIE", nameStart: "AAAAA  123  (BBB)", nameEnd: "GEIE"},
		{fullName: "AAAAA  123  (BBB) GELE", nameStart: "AAAAA  123  (BBB)", nameEnd: "GELE"},
		{fullName: "AAAAA  123  (BBB) ICVC", nameStart: "AAAAA  123  (BBB)", nameEnd: "ICVC"},
		{fullName: "AAAAA  123  (BBB) INVESTMENT COMPANY WITH VARIABLE CAPITAL", nameStart: "AAAAA  123  (BBB)", nameEnd: "INVESTMENT COMPANY WITH VARIABLE CAPITAL"},
		{fullName: "AAAAA  123  (BBB) L.P", nameStart: "AAAAA  123  (BBB)", nameEnd: "L.P"},
		{fullName: "AAAAA  123  (BBB) L.T.D", nameStart: "AAAAA  123  (BBB)", nameEnd: "L.T.D"},
		{fullName: "AAAAA  123  (BBB) LIMITED - THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "LIMITED - THE"},
		{fullName: "AAAAA  123  (BBB) LIMITED LIABILITY PARTNERSHIP", nameStart: "AAAAA  123  (BBB)", nameEnd: "LIMITED LIABILITY PARTNERSHIP"},
		{fullName: "AAAAA  123  (BBB) LIMITED PARTNERSHIP", nameStart: "AAAAA  123  (BBB)", nameEnd: "LIMITED PARTNERSHIP"},
		{fullName: "AAAAA  123  (BBB) LIMITED THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "LIMITED THE"},
		{fullName: "AAAAA  123  (BBB) LIMITED", nameStart: "AAAAA  123  (BBB)", nameEnd: "LIMITED"},
		{fullName: "AAAAA  123  (BBB) LIMITED-THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "LIMITED-THE"},
		{fullName: "AAAAA  123  (BBB) LIMITED...THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "LIMITED...THE"},
		{fullName: "AAAAA  123  (BBB) LIMITED..THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "LIMITED..THE"},
		{fullName: "AAAAA  123  (BBB) LIMITED.THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "LIMITED.THE"},
		{fullName: "AAAAA  123  (BBB) LLP", nameStart: "AAAAA  123  (BBB)", nameEnd: "LLP"},
		{fullName: "AAAAA  123  (BBB) LP", nameStart: "AAAAA  123  (BBB)", nameEnd: "LP"},
		{fullName: "AAAAA  123  (BBB) LTD", nameStart: "AAAAA  123  (BBB)", nameEnd: "LTD"},
		{fullName: "AAAAA  123  (BBB) LTD...THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "LTD...THE"},
		{fullName: "AAAAA  123  (BBB) LTD..THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "LTD..THE"},
		{fullName: "AAAAA  123  (BBB) LTD.THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "LTD.THE"},
		{fullName: "AAAAA  123  (BBB) OEIC", nameStart: "AAAAA  123  (BBB)", nameEnd: "OEIC"},
		{fullName: "AAAAA  123  (BBB) OPEN-ENDED INVESTMENT COMPANY", nameStart: "AAAAA  123  (BBB)", nameEnd: "OPEN-ENDED INVESTMENT COMPANY"},
		{fullName: "AAAAA  123  (BBB) P.L.C", nameStart: "AAAAA  123  (BBB)", nameEnd: "P.L.C"},
		{fullName: "AAAAA  123  (BBB) PAC", nameStart: "AAAAA  123  (BBB)", nameEnd: "PAC"},
		{fullName: "AAAAA  123  (BBB) PARTNERIAETH ATEBOLRWYDD CYFYNGEDIG", nameStart: "AAAAA  123  (BBB)", nameEnd: "PARTNERIAETH ATEBOLRWYDD CYFYNGEDIG"},
		{fullName: "AAAAA  123  (BBB) PARTNERIAETH CYFYNGEDIG", nameStart: "AAAAA  123  (BBB)", nameEnd: "PARTNERIAETH CYFYNGEDIG"},
		{fullName: "AAAAA  123  (BBB) PCC LIMITED", nameStart: "AAAAA  123  (BBB)", nameEnd: "PCC LIMITED"},
		{fullName: "AAAAA  123  (BBB) PCC LTD", nameStart: "AAAAA  123  (BBB)", nameEnd: "PCC LTD"},
		{fullName: "AAAAA  123  (BBB) PCC", nameStart: "AAAAA  123  (BBB)", nameEnd: "PCC"},
		{fullName: "AAAAA  123  (BBB) PLC", nameStart: "AAAAA  123  (BBB)", nameEnd: "PLC"},
		{fullName: "AAAAA  123  (BBB) PROTECTED CELL COMPANY", nameStart: "AAAAA  123  (BBB)", nameEnd: "PROTECTED CELL COMPANY"},
		{fullName: "AAAAA  123  (BBB) PUBLIC LIMITED COMPANY .THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "PUBLIC LIMITED COMPANY .THE"},
		{fullName: "AAAAA  123  (BBB) PUBLIC LIMITED COMPANY THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "PUBLIC LIMITED COMPANY THE"},
		{fullName: "AAAAA  123  (BBB) PUBLIC LIMITED COMPANY", nameStart: "AAAAA  123  (BBB)", nameEnd: "PUBLIC LIMITED COMPANY"},
		{fullName: "AAAAA  123  (BBB) PUBLIC LIMITED COMPANY.THE", nameStart: "AAAAA  123  (BBB)", nameEnd: "PUBLIC LIMITED COMPANY.THE"},
		{fullName: "AAAAA  123  (BBB) SE", nameStart: "AAAAA  123  (BBB)", nameEnd: "SE"},
		{fullName: "AAAAA  123  (BBB) UKEIG", nameStart: "AAAAA  123  (BBB)", nameEnd: "UKEIG"},
		{fullName: "AAAAA  123  (BBB) UK SOCIETAS", nameStart: "AAAAA  123  (BBB)", nameEnd: "UK SOCIETAS"},
		{fullName: "AAAAA  123  (BBB) UNLIMITED", nameStart: "AAAAA  123  (BBB)", nameEnd: "UNLIMITED"},
		{fullName: "AAAAA  123  (BBB) UNLTD", nameStart: "AAAAA  123  (BBB)", nameEnd: "UNLTD"},

		// real Companies test:
		{fullName: "BURRAGE AND BOYDE,LIMITED", nameStart: "BURRAGE AND BOYDE,LIMITED", nameEnd: empty},
		{fullName: "TRUNIN LTD", nameStart: "TRUNIN", nameEnd: "LTD"},
		{fullName: "ROSSI94 LIMITED", nameStart: "ROSSI94", nameEnd: "LIMITED"},
		{fullName: "WALKER MINERALS (BRADWELL) LTD", nameStart: "WALKER MINERALS (BRADWELL)", nameEnd: "LTD"},
		{fullName: "EXTRALEASE RESIDENTS MANAGEMENT LIMITED", nameStart: "EXTRALEASE RESIDENTS MANAGEMENT", nameEnd: "LIMITED"},
		{fullName: "ENCON AIR SYSTEMS LIMITED", nameStart: "ENCON AIR SYSTEMS", nameEnd: "LIMITED"},
		{fullName: "T.A.W. SHOPFITTING & BUILDING SERVICES LIMITED", nameStart: "T.A.W. SHOPFITTING & BUILDING SERVICES", nameEnd: "LIMITED"},
		{fullName: "30/32 SPILLER ROAD LIMITED", nameStart: "30/32 SPILLER ROAD", nameEnd: "LIMITED"},
		{fullName: "JOHN GOODWIN (1986) LIMITED", nameStart: "JOHN GOODWIN (1986)", nameEnd: "LIMITED"},
		{fullName: "F. & J. HAWKYARD (ELLAND) LIMITED", nameStart: "F. & J. HAWKYARD (ELLAND)", nameEnd: "LIMITED"},
		{fullName: "W.B.C. BOOK MANUFACTURERS LIMITED", nameStart: "W.B.C. BOOK MANUFACTURERS", nameEnd: "LIMITED"},
		{fullName: "BURFIELD CONSTRUCTION (LONDON) LIMITED", nameStart: "BURFIELD CONSTRUCTION (LONDON)", nameEnd: "LIMITED"},
		{fullName: "SOUTH WEST STAINLESS LTD", nameStart: "SOUTH WEST STAINLESS", nameEnd: "LTD"},
		{fullName: "BAR 67 LTD", nameStart: "BAR 67", nameEnd: "LTD"},
		{fullName: "SUSSEX AIR CONDITIONING & HEAT PUMPS LIMITED", nameStart: "SUSSEX AIR CONDITIONING & HEAT PUMPS", nameEnd: "LIMITED"},
		{fullName: "JADHAV MEDICAL SERVICES LIMITED", nameStart: "JADHAV MEDICAL SERVICES", nameEnd: "LIMITED"},
		{fullName: "BROOKSON (5665N) LIMITED", nameStart: "BROOKSON (5665N)", nameEnd: "LIMITED"},
		{fullName: "STONELEIGH DESIGN AND BUILD LTD", nameStart: "STONELEIGH DESIGN AND BUILD", nameEnd: "LTD"},
		{fullName: "LSIL 2013 LIMITED", nameStart: "LSIL 2013", nameEnd: "LIMITED"},
		{fullName: "LANCING CONSULTING SERVICES LIMITED", nameStart: "LANCING CONSULTING SERVICES", nameEnd: "LIMITED"},
		{fullName: "HAZEYJANE IT SERVICES LTD", nameStart: "HAZEYJANE IT SERVICES", nameEnd: "LTD"},
		{fullName: "ROBMAR (CONSULTANCY) LIMITED", nameStart: "ROBMAR (CONSULTANCY)", nameEnd: "LIMITED"},
		{fullName: "BNY MELLON MANAGED FUNDS II", nameStart: "BNY MELLON MANAGED FUNDS II", nameEnd: empty},
		{fullName: "AT TRADING (UK) LIMITED", nameStart: "AT TRADING (UK)", nameEnd: "LIMITED"},
		{fullName: "EQT NEPTUNE BENSON CO-INVESTMENT LIMITED PARTNERSHIP", nameStart: "EQT NEPTUNE BENSON CO-INVESTMENT", nameEnd: "LIMITED PARTNERSHIP"},
		{fullName: "ALDON COMPUTER GROUP", nameStart: "ALDON COMPUTER GROUP", nameEnd: empty},
		{fullName: "MORGAN STANLEY HERCULES S.A.R.L.", nameStart: "MORGAN STANLEY HERCULES S.A.R.L.", nameEnd: empty},
		{fullName: "JOHN KELLY FUELS (IRELAND)", nameStart: "JOHN KELLY FUELS (IRELAND)", nameEnd: empty},
		{fullName: "SELF-OPENING TIN BOX COMPANY,LIMITED", nameStart: "SELF-OPENING TIN BOX COMPANY,LIMITED", nameEnd: empty},
		{fullName: "SPARECO 21 LIMITED", nameStart: "SPARECO 21", nameEnd: "LIMITED"},
		{fullName: "S.K.B INFOTECH LOGISTICS LIMITED", nameStart: "S.K.B INFOTECH LOGISTICS", nameEnd: "LIMITED"},
		{fullName: "ABLE SCHOOL OF MOTORING LIMITED", nameStart: "ABLE SCHOOL OF MOTORING", nameEnd: "LIMITED"},
		{fullName: "GET 2 SET LTD", nameStart: "GET 2 SET", nameEnd: "LTD"},
		{fullName: "BEATA'S CLEANING COMPANY LIMITED", nameStart: "BEATA'S CLEANING COMPANY", nameEnd: "LIMITED"},
		{fullName: "MRS PROPERTIES (LANCASHIRE) LTD", nameStart: "MRS PROPERTIES (LANCASHIRE)", nameEnd: "LTD"},
		{fullName: "CAL'S COMPUTER SERVICES LIMITED", nameStart: "CAL'S COMPUTER SERVICES", nameEnd: "LIMITED"},
		{fullName: "FCC WREXHAM PFI (PHASE II HOLDING) LIMITED", nameStart: "FCC WREXHAM PFI (PHASE II HOLDING)", nameEnd: "LIMITED"},
		{fullName: "LABYRINTH I.T. SERVICES LIMITED", nameStart: "LABYRINTH I.T. SERVICES", nameEnd: "LIMITED"},
		{fullName: "BONYJ HEALTHCARE PROFESSIONALS RECRUITMENT SERVICES LIMITED", nameStart: "BONYJ HEALTHCARE PROFESSIONALS RECRUITMENT SERVICES", nameEnd: "LIMITED"},
		{fullName: "MUSTANG SALLY'S DINER LIMITED", nameStart: "MUSTANG SALLY'S DINER", nameEnd: "LIMITED"},
		{fullName: "SPRIGGS & SONS LLP", nameStart: "SPRIGGS & SONS", nameEnd: "LLP"},
		{fullName: "EAT & DRINK INN LIMITED", nameStart: "EAT & DRINK INN", nameEnd: "LIMITED"},
		{fullName: "MILLYS BAKERY & COFFEE SHOP LIMITED", nameStart: "MILLYS BAKERY & COFFEE SHOP", nameEnd: "LIMITED"},
		{fullName: "SIDLEY CHEMICAL CO., LTD", nameStart: "SIDLEY CHEMICAL CO.,", nameEnd: "LTD"},
		{fullName: "BEVERLEY (FB) LTD", nameStart: "BEVERLEY (FB)", nameEnd: "LTD"},
		{fullName: "LFR TECHNOLOGY UK LTD.", nameStart: "LFR TECHNOLOGY UK", nameEnd: "LTD."},
		{fullName: "NSBAA PROPERTY COMPANY LTD.", nameStart: "NSBAA PROPERTY COMPANY", nameEnd: "LTD."},
		{fullName: "BELL INSTALLATIONS & BUILDING LIMITED", nameStart: "BELL INSTALLATIONS & BUILDING", nameEnd: "LIMITED"},
		{fullName: "ANA K LTD", nameStart: "ANA K", nameEnd: "LTD"},
		{fullName: "PROCESS PRINT (LABELLING & DESIGN) LIMITED", nameStart: "PROCESS PRINT (LABELLING & DESIGN)", nameEnd: "LIMITED"},
		{fullName: "A S PREMIUM FUNDING LIMITED", nameStart: "A S PREMIUM FUNDING", nameEnd: "LIMITED"},
		{fullName: "KINGSWOOD LETTINGS & PROPERTY MANAGEMENT LIMITED", nameStart: "KINGSWOOD LETTINGS & PROPERTY MANAGEMENT", nameEnd: "LIMITED"},
		{fullName: "BUSINESS HOMES NORTH (WEST) LIMITED", nameStart: "BUSINESS HOMES NORTH (WEST)", nameEnd: "LIMITED"},
		{fullName: "RCDS (YORKSHIRE) LIMITED", nameStart: "RCDS (YORKSHIRE)", nameEnd: "LIMITED"},
		{fullName: "JEMSONS & SAT LTD", nameStart: "JEMSONS & SAT", nameEnd: "LTD"},
		{fullName: "GREEN CLEEN (BRAINTREE/CHELMSFORD) LIMITED", nameStart: "GREEN CLEEN (BRAINTREE/CHELMSFORD)", nameEnd: "LIMITED"},
		{fullName: "S.I.H.W. LIMITED", nameStart: "S.I.H.W.", nameEnd: "LIMITED"},
		{fullName: "BROWN & HUDSON EXPEDITIONS LIMITED", nameStart: "BROWN & HUDSON EXPEDITIONS", nameEnd: "LIMITED"},
		{fullName: "KRAZY K9'S LIMITED", nameStart: "KRAZY K9'S", nameEnd: "LIMITED"},
		{fullName: "HOMEWORKS PROPERTY MAINTENANCE (N.E.) LIMITED", nameStart: "HOMEWORKS PROPERTY MAINTENANCE (N.E.)", nameEnd: "LIMITED"},
		{fullName: "COURAGE HEALTH CARE (WEST MIDLANDS) LIMITED", nameStart: "COURAGE HEALTH CARE (WEST MIDLANDS)", nameEnd: "LIMITED"},
		{fullName: "THE EDGAR HARRINGTON BICYCLE CO. LIMITED", nameStart: "THE EDGAR HARRINGTON BICYCLE CO.", nameEnd: "LIMITED"},
		{fullName: "HOLYOKE L.P.", nameStart: "HOLYOKE", nameEnd: "L.P."},
		{fullName: "NOVO BANCO, S.A.", nameStart: "NOVO BANCO, S.A.", nameEnd: empty},
		{fullName: "TAN CHONG TRADING (MALAYSIA) SDN.BHD. (UK) REPRESENTATIVE OFFICE", nameStart: "TAN CHONG TRADING (MALAYSIA) SDN.BHD. (UK) REPRESENTATIVE OFFICE", nameEnd: empty},
		{fullName: "CI QUALITY LTD.", nameStart: "CI QUALITY", nameEnd: "LTD."},
		{fullName: "PARTNERS GROUP ACCESS 70 L.P.", nameStart: "PARTNERS GROUP ACCESS 70", nameEnd: "L.P."},
		{fullName: "G. COX LIFT ENGINEERS LTD", nameStart: "G. COX LIFT ENGINEERS", nameEnd: "LTD"},
		{fullName: "BRINDAL ELECTRICAL & SECURITY SOLUTIONS LTD", nameStart: "BRINDAL ELECTRICAL & SECURITY SOLUTIONS", nameEnd: "LTD"},
		{fullName: "P S B A LIMITED", nameStart: "P S B A", nameEnd: "LIMITED"},
		{fullName: "THE TUBE CREATION CO. LIMITED", nameStart: "THE TUBE CREATION CO.", nameEnd: "LIMITED"},
		{fullName: "81G LIMITED", nameStart: "81G", nameEnd: "LIMITED"},
		{fullName: "BG&J UK LTD", nameStart: "BG&J UK", nameEnd: "LTD"},
		{fullName: "W. & J. TOLMIE LIMITED", nameStart: "W. & J. TOLMIE", nameEnd: "LIMITED"},
		{fullName: "K&C CONTRACTORS LTD", nameStart: "K&C CONTRACTORS", nameEnd: "LTD"},
		{fullName: "PD 2007 LIMITED", nameStart: "PD 2007", nameEnd: "LIMITED"},
		{fullName: "HUNTLEY & COCKRILL BUILDING AND GROUNDWORKERS LIMITED", nameStart: "HUNTLEY & COCKRILL BUILDING AND GROUNDWORKERS", nameEnd: "LIMITED"},
		{fullName: "C & S DOMESTIC APPLIANCE REPAIRS LIMITED", nameStart: "C & S DOMESTIC APPLIANCE REPAIRS", nameEnd: "LIMITED"},
		{fullName: "34@HILTON PARK GOLF LTD", nameStart: "34@HILTON PARK GOLF", nameEnd: "LTD"},
		{fullName: "ASSOCIATED SUB-CONTRACTORS IN LIMITED PARTNERSHIP WITH JOHN LARMOUR", nameStart: "ASSOCIATED SUB-CONTRACTORS IN LIMITED PARTNERSHIP WITH JOHN LARMOUR", nameEnd: empty},
		{fullName: "AFRIKAANS PROTESTANT CHURCH IN ENGLAND AND WALES", nameStart: "AFRIKAANS PROTESTANT CHURCH IN ENGLAND AND WALES", nameEnd: empty},
		{fullName: "D S BOYLE & CO LTD", nameStart: "D S BOYLE & CO", nameEnd: "LTD"},
		{fullName: "ASSOCIATED SUB-CONTRACTORS IN LIMITED PARTNERSHIP WITH DANIEL FERRETTI", nameStart: "ASSOCIATED SUB-CONTRACTORS IN LIMITED PARTNERSHIP WITH DANIEL FERRETTI", nameEnd: empty},
		{fullName: "CS CAPITAL PARTNERS V, L.P.", nameStart: "CS CAPITAL PARTNERS V,", nameEnd: "L.P."},
		{fullName: "HEALTHEAST C.I.C.", nameStart: "HEALTHEAST", nameEnd: "C.I.C."},
		{fullName: "HYREX CONSULTING L.P.", nameStart: "HYREX CONSULTING", nameEnd: "L.P."},
		{fullName: "MORGAN STANLEY LYRA S.A.R.L.", nameStart: "MORGAN STANLEY LYRA S.A.R.L.", nameEnd: empty},
		{fullName: "LASALLE UK VENTURES PROPERTY 10", nameStart: "LASALLE UK VENTURES PROPERTY 10", nameEnd: empty},
		{fullName: "720 CONSULTING SERVICES LTD", nameStart: "720 CONSULTING SERVICES", nameEnd: "LTD"},
		{fullName: "V.K.N TRANSPORT LTD", nameStart: "V.K.N TRANSPORT", nameEnd: "LTD"},
		{fullName: "PENTA INVESTMENTS (2012) SP LIMITED PARTNERSHIP", nameStart: "PENTA INVESTMENTS (2012) SP", nameEnd: "LIMITED PARTNERSHIP"},
		{fullName: "LOGISTIC MASTERS LIMITED", nameStart: "LOGISTIC MASTERS", nameEnd: "LIMITED"},
	}

	// extra: every entry of companyNameEndings appended to a constant name start
	tListExtra1   = buildEndingsTable("AAA1 BBB2 C-;;^$£.--3")
	tListExtra1t1 = buildEndingsTable("AAA1 BBB2 C-;;^$£.--3", ".")
	tListExtra2   = buildEndingsTable("AAA1")
	tListExtra3   = buildEndingsTable("AAA1  ")
	tListExtra3t1 = buildEndingsTable("AAA1  ", ".")
	tListExtra3t2 = buildEndingsTable("AAA1  ", " ")
)

// buildEndingsTable returns one entry per companyNameEndings item, made of
// the given name start followed by that ending and, optionally, a trail string
// (default empty). The trail is expected to be punctuation, which is kept
// attached to the ending when the name is split.
func buildEndingsTable(nameStart string, trailString ...string) []tInfo {
	trail := ""
	if len(trailString) > 0 {
		trail = trailString[0]
	}
	table := make([]tInfo, 0, len(companyNameEndings))
	for _, ending := range companyNameEndings {
		table = append(table, tInfo{
			fullName:  nameStart + " " + ending + trail,
			nameStart: strings.TrimSpace(nameStart),
			nameEnd:   ending + trail,
		})
	}
	return table
}

func showText(t1, t2 string) string {
	return t1 + " [" + t2 + "]"
}

func TestUnitSplitCompanyNameEndings(t *testing.T) {
	f := NewFormatter()

	Convey("SplitCompanyNameEndings should split standard endings", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST LIMITED")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, "LIMITED")
	})

	Convey("SplitCompanyNameEndings should prefer the longest valid ending", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST PUBLIC LIMITED COMPANY")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, "PUBLIC LIMITED COMPANY")
	})

	Convey("SplitCompanyNameEndings should not treat partial token matches as endings", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST UNLTD")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, "UNLTD")
	})

	Convey("SplitCompanyNameEndings should preserve trailing punctuation on endings", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST PCC LTD...")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, "PCC LTD...")
	})

	Convey("SplitCompanyNameEndings should handle irregular whitespace inside the ending", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST PUBLIC  LIMITED\tCOMPANY")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, "PUBLIC LIMITED COMPANY")
	})

	Convey("SplitCompanyNameEndings should handle irregular whitespace inside the ending", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST PUBLIC  LIMITED   COMPANY")
		So(nameStart, ShouldEqual, "TEST")
		So(nameEnd, ShouldEqual, "PUBLIC LIMITED COMPANY")
	})

	Convey("SplitCompanyNameEndings should return original name when no ending is present", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("TEST HOLDINGS")
		So(nameStart, ShouldEqual, "TEST HOLDINGS")
		So(nameEnd, ShouldEqual, "")
	})

	Convey("SplitCompanyNameEndings should literally match the endings (dots matched as dots)", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("STEF LTD...THE")
		So(nameStart, ShouldEqual, "STEF")
		So(nameEnd, ShouldEqual, "LTD...THE")
	})
	Convey("SplitCompanyNameEndings should literally match the endings (spaces not matched)", t, func() {
		nameStart, nameEnd := f.SplitCompanyNameEndings("STEF LTD   THE")
		So(nameStart, ShouldEqual, "STEF LTD   THE")
		So(nameEnd, ShouldEqual, "")
	})

	allTests := make([]tInfo, 0, len(tList)+len(tListExtra1)+len(tListExtra2)+len(tListExtra3))
	allTests = append(allTests, tList...)
	allTests = append(allTests, tListExtra1...)
	allTests = append(allTests, tListExtra1t1...)
	allTests = append(allTests, tListExtra2...)
	allTests = append(allTests, tListExtra3...)
	allTests = append(allTests, tListExtra3t1...)
	allTests = append(allTests, tListExtra3t2...)

	for _, tI := range allTests {
		Convey(showText("Given a company named", tI.fullName), t, func() {
			Convey("When SplitCompanyNameEndings is called", func() {
				nameStart, nameEnd := f.SplitCompanyNameEndings(tI.fullName)
				Convey(showText("Then name should be", tI.nameStart), func() {
					So(nameStart, ShouldEqual, tI.nameStart)
					Convey(showText("And ending", tI.nameEnd), func() {
						So(nameEnd, ShouldEqual, tI.nameEnd)
					})
				})
			})
		})
	}
}
