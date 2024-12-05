package models

var UserTable = struct {
	TableName string
	Firstname string
	Surname   string
	Password1 string
	Email     string
}{
	TableName: "UserSQLModel",
	Firstname: "firstname",
	Surname:   "surname",
	Password1: "paswoord1",
	Email:     "email",
}

var BuildingTable = struct {
	TableName  string
	Gid        string
	Shape__are string
	Shape__len string
	Globalid   string
	Creationda string
	Creator    string
	Editdate   string
	Editor     string
	Num_storey string
	Building_t string
	Ghanapost_ string
	Plot_numbe string
	Developmen string
	Name       string
	Parcel_id  string
	Exact_use  string
	Building_u string
	Remarks    string
	Other_info string
	Other_in_1 string
	Geom       string
}{
	TableName:  "building",
	Gid:        "gid",
	Shape__are: "shape__are",
	Shape__len: "shape__len",
	Globalid:   "globalid",
	Creationda: "creationda",
	Creator:    "creator",
	Editdate:   "editdate",
	Editor:     "editor",
	Num_storey: "num_storey",
	Building_t: "building_t",
	Ghanapost_: "ghanapost_",
	Plot_numbe: "plot_numbe",
	Developmen: "developmen",
	Name:       "name",
	Parcel_id:  "parcel_id",
	Exact_use:  "exact_use",
	Building_u: "building_u",
	Remarks:    "remarks",
	Other_info: "other_info",
	Other_in_1: "other_in_1",
	Geom:       "geom",
}

var OtherPolygonTable = struct {
	TableName  string
	Gid        string
	Shape__are string
	Shape__len string
	Globalid   string
	Creationda string
	Creator    string
	Editdate   string
	Editor     string
	Usage1     string
	Ghanapostg string
	Developmen string
	Name       string
	Structure_ string
	Parcel_id  string
	Exact_use  string
	Remarks    string
	Other_info string
	Other_in_1 string
	Geom       string
	Street_nam string
	Mixed_usag string
}{
	TableName:  "other_polygon_structure",
	Gid:        "gid",
	Shape__are: "shape__are",
	Shape__len: "shape__len",
	Globalid:   "globalid",
	Creationda: "creationda",
	Creator:    "creator",
	Editdate:   "editdate",
	Editor:     "editor",
	Usage1:     "usage1",
	Structure_: "structure_",
	Ghanapostg: "ghanapostg",
	Developmen: "developmen",
	Name:       "name",
	Parcel_id:  "parcel_id",
	Exact_use:  "exact_use",
	Remarks:    "remarks",
	Other_info: "other_info",
	Other_in_1: "other_in_1",
	Street_nam: "street_nam",
	Geom:       "geom",
	Mixed_usag: "mixed_usag",
}
