package profile

var starterItems = map[string]int{
	"AthenaPickaxe:DefaultPickaxe":                               1,
	"AthenaGlider:DefaultGlider":                                 1,
	"AthenaDance:EID_DanceMoves":                                 1,
	"AthenaCharacter:CID_001_Athena_Commando_F_Default":          1,
	"AthenaCharacter:CID_002_Athena_Commando_F_Default":          1,
	"AthenaCharacter:CID_003_Athena_Commando_F_Default":          1,
	"AthenaCharacter:CID_004_Athena_Commando_F_Default":          1,
	"AthenaCharacter:CID_005_Athena_Commando_M_Default":          1,
	"AthenaCharacter:CID_006_Athena_Commando_M_Default":          1,
	"AthenaCharacter:CID_007_Athena_Commando_M_Default":          1,
	"AthenaCharacter:CID_008_Athena_Commando_M_Default":          1,
	"AthenaCharacter:CID_556_Athena_Commando_F_RebirthDefaultA":  11,
	"AthenaCharacter:CID_557_Athena_Commando_F_RebirthDefaultB":  11,
	"AthenaCharacter:CID_558_Athena_Commando_F_RebirthDefaultC":  11,
	"AthenaCharacter:CID_559_Athena_Commando_F_RebirthDefaultD":  11,
	"AthenaCharacter:CID_560_Athena_Commando_M_RebirthDefaultA":  11,
	"AthenaCharacter:CID_561_Athena_Commando_M_RebirthDefaultB":  11,
	"AthenaCharacter:CID_562_Athena_Commando_M_RebirthDefaultC":  11,
	"AthenaCharacter:CID_563_Athena_Commando_M_RebirthDefaultD":  11,
	"AthenaCharacter:CID_873_Athena_Commando_M_RebirthDefaultE":  14,
	"AthenaCharacter:CID_874_Athena_Commando_M_RebirthDefaultF":  14,
	"AthenaCharacter:CID_875_Athena_Commando_M_RebirthDefaultG":  14,
	"AthenaCharacter:CID_876_Athena_Commando_M_RebirthDefaultH":  14,
	"AthenaCharacter:CID_877_Athena_Commando_M_RebirthDefaultI":  14,
	"AthenaCharacter:CID_878_Athena_Commando_F_RebirthDefault_E": 14,
	"AthenaCharacter:CID_879_Athena_Commando_F_RebirthDefault_F": 14,
	"AthenaCharacter:CID_880_Athena_Commando_F_RebirthDefault_G": 14,
	"AthenaCharacter:CID_881_Athena_Commando_F_RebirthDefault_H": 14,
	"AthenaCharacter:CID_882_Athena_Commando_F_RebirthDefault_I": 14,
}

func availableIn(templateID string, season int) bool {
	introduced, starter := starterItems[templateID]
	return !starter || introduced <= season
}
