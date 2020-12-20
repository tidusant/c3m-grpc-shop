package main

import (
	"github.com/tidusant/c3m-common/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tidusant/c3m-common/c3mcommon"
	pb "github.com/tidusant/c3m-grpc-protoc/protoc"
	rpch "github.com/tidusant/chadmin-repo/cuahang"
	"github.com/tidusant/chadmin-repo/models"

	"encoding/json"
	"fmt"
	"net"

	"os"
	"strconv"

	"context"
	"google.golang.org/grpc"
)

const (
	name string = "auth"
	ver  string = "1"
)

type service struct {
	pb.UnimplementedGRPCServicesServer
}

func (s *service) Call(ctx context.Context, in *pb.RPCRequest) (*pb.RPCResponse, error) {
	resp := &pb.RPCResponse{Data: "Hello " + in.GetAppName(), RPCName: name, Version: ver}
	rs := models.RequestResult{}
	//get input data into user session
	var usex models.UserSession
	usex.Session = in.Session
	usex.Action = in.Action

	usex.UserIP = in.UserIP
	usex.Params = in.Params
	usex.UserID, _ = primitive.ObjectIDFromHex(in.UserID)

	//check shop permission
	if in.ShopID != "" {
		shopidObj, _ := primitive.ObjectIDFromHex(in.ShopID)
		shop := rpch.GetShopById(usex.UserID, shopidObj)
		if shop.Status == 0 {
			if usex.Action != "lsi" {
				rs.Error = "Site is disable"
			}
		}
		usex.Shop = shop
	}

	if usex.Action == "cs" {
		rs = changeShop(usex)
	} else if usex.Action == "lsi" {
		rs = loadshopinfo(usex)
	} else if usex.Action == "ca" {
		rs = doCreateAlbum(usex)
	} else if usex.Action == "la" {
		rs = doLoadalbum(usex)
	} else if usex.Action == "ea" {
		rs = doEditAlbum(usex)
	} else if usex.Action == "cga" {
		rs = configGetAll(usex)
	} else if usex.Action == "cgs" {
		rs = configSave(usex)
	} else if usex.Action == "lims" {
		rs = getShopLimits(usex)
	} else {
		//unknow action
		return resp, nil
	}
	//convert RequestResult into json
	b, _ := json.Marshal(rs)
	resp.Data = string(b)
	return resp, nil

}

type ConfigViewData struct {
	ShopConfigs     models.ShopConfigs
	TemplateConfigs []ConfigItem
	BuildConfigs    models.BuildConfig
}
type ConfigItem struct {
	Key   string
	Type  string
	Value string
}

func loadshopinfo(usex models.UserSession) models.RequestResult {
	strrt := `{"Shop":`
	b, _ := json.Marshal(usex.Shop)
	strrt += string(b)

	//get langs info
	strrt += `,"Languages":[`
	for _, lang := range usex.Shop.Config.Langs {
		strrt += `{"Code":"` + lang + `","Name":"` + c3mcommon.GetLangnameByCode(lang) + `","Flag":"` + c3mcommon.Code2Flag(lang) + `"},`
	}
	if len(usex.Shop.Config.Langs) > 0 {
		strrt = strrt[:len(strrt)-1] + `]`
	}
	b, _ = json.Marshal(usex.Shop.Config)
	strrt += `,"ShopConfigs":` + string(b)

	//maxfileupload
	strrt += `,"MaxFileUpload":` + strconv.Itoa(rpch.GetShopLimitbyKey(usex.Shop.ID, "maxfileupload"))
	strrt += `,"MaxSizeUpload":` + strconv.Itoa(rpch.GetShopLimitbyKey(usex.Shop.ID, "maxsizeupload"))

	//orther shop
	otherShops := rpch.GetOtherShopById(usex.UserID, usex.Shop.ID)
	strrt += `,"Others":[`
	for _, shop := range otherShops {
		strrt += `{"Name":"` + shop.Name + `","ID":"` + shop.ID.Hex() + `"},`
	}
	if len(otherShops) > 0 {
		strrt = strrt[:len(strrt)-1] + `]`
	} else {
		strrt += `]`
	}

	//get user info
	user := rpch.GetUserInfo(usex.UserID)
	strrt += `,"User":{"Name":"` + user.Name + `"}`
	strrt += "}"
	return models.RequestResult{Status: 1, Error: "", Message: "", Data: strrt}

}

func changeShop(usex models.UserSession) models.RequestResult {
	shopidObj, _ := primitive.ObjectIDFromHex(usex.Params)
	shopchange := rpch.UpdateShopLogin(usex.Session, shopidObj)
	if shopchange.ID == primitive.NilObjectID {
		return models.RequestResult{Error: "Change shop fail"}

	}
	//change shopid
	usex.Shop = shopchange
	return loadshopinfo(usex)
}
func configSave(usex models.UserSession) models.RequestResult {
	//var config ConfigViewData

	//err := json.Unmarshal([]byte(usex.Params), &config)
	//if !c3mcommon.CheckError("json parse page", err) {
	//	return models.RequestResult{Error: "json parse fail"}
	//}
	//usex.Shop.Config = config.ShopConfigs
	//rpch.SaveShopConfig(usex.Shop)
	//
	//// //save template config
	//str := `{"Code":"` + usex.Shop.Theme + `","TemplateConfigs":[{}`
	//for _, conf := range config.TemplateConfigs {
	//	str += `,{"Key":"` + conf.Key + `","Value":"` + conf.Value + `"}`
	//}
	//str += `]`
	//b, _ := json.Marshal(config.BuildConfigs)
	//str += `,"BuildConfig":` + string(b) + `}`
	//
	//request := "savetemplateconfig|" + usex.Session
	//resp := c3mcommon.RequestBuildService(request, "POST", str)
	//
	//if resp.Status != 1 {
	//	return resp
	//}
	//
	//// //save build config
	//
	//// var bcf models.BuildConfig
	//// bcf = config.BuildConfigs
	//// bcf.ShopId = usex.Shop.ID
	//// rpb.SaveConfig(bcf)
	////rebuild config
	//rpch.Rebuild(usex)
	//return c3mcommon.ReturnJsonMessage("1", "", "success", "")
	return models.RequestResult{Status: 1}

}
func configGetAll(usex models.UserSession) models.RequestResult {
	//var config ConfigViewData
	//config.ShopConfigs = usex.Shop.Config
	//log.Debugf("configGetAll")
	//request := "gettemplateconfig|" + usex.Session
	//resp := c3mcommon.RequestBuildService(request, "POST", usex.Shop.Theme)
	//log.Debugf("RequestBuildService call done")
	//if resp.Status != 1 {
	//	return resp
	//}
	//var confs struct {
	//	TemplateConfigs []ConfigItem
	//	BuildConfigs    models.BuildConfig
	//}
	//json.Unmarshal([]byte(resp.Data), &confs)
	//
	//config.TemplateConfigs = confs.TemplateConfigs
	//config.BuildConfigs = confs.BuildConfigs
	//config.BuildConfigs.ID = ""
	//config.BuildConfigs.ShopId = ""
	//b, _ := json.Marshal(config)
	//
	//return c3mcommon.ReturnJsonMessage("1", "", "success", string(b))
	return models.RequestResult{Status: 1}

}
func getShopLimits(usex models.UserSession) models.RequestResult {

	limits := rpch.GetShopLimits(usex.Shop.ID)

	b, _ := json.Marshal(limits)
	return models.RequestResult{Status: 1, Error: "", Data: string(b)}

}

// func loadcat(usex models.UserSession) string {
// 	log.Debugf("loadcat begin")
// 	shop := rpch.GetShopById(usex.UserID, usex.ShopID)

// 	strrt := "["
// 	log.Debugf("load cats:%v", shop.ShopCats)
// 	catinfstr := ""
// 	for _, cat := range shop.ShopCats {
// 		catlangs := ""
// 		for lang, catinf := range cat.Langs {
// 			catlangs += """ + lang + "":{"name":"" + catinf.Slug + "","slug":"" + catinf.Name + "","description":"" + catinf.Description + ""},"
// 		}
// 		catlangs = catlangs[:len(catlangs)-1]
// 		catinfstr += "{"code":"" + cat.Code + "","langs":{" + catlangs + "}},"
// 	}
// 	if catinfstr == "" {
// 		strrt += "{}]"
// 	} else {
// 		strrt += catinfstr[:len(catinfstr)-1] + "]"
// 	}

// 	return c3mcommon.ReturnJsonMessage("1", "", "success", strrt)

// }

func doCreateAlbum(usex models.UserSession) models.RequestResult {
	//albumname := usex.Params
	//if albumname == "" {
	//	return models.RequestResult{Error:"album's name empty"}
	//
	//}
	////get config
	//
	//if usex.Shop.ID == "" {
	//	return models.RequestResult{Error:"shop not found"}
	//
	//}
	//
	//// if usex.Shop.Config.Level == 0 {
	//// 	return c3mcommon.ReturnJsonMessage("0", "config error", "", "")
	//
	//// }
	//// if usex.Shop.Config.MaxAlbum <= len(usex.Shop.Albums) {
	//// 	return c3mcommon.ReturnJsonMessage("2", "album count limited", "", "")
	//// }
	//
	//slug := strings.ToLower(mystring.Camelize(mystring.Asciify(albumname)))
	//albumslug := slug
	//if slug == "" {
	//	return models.RequestResult{Error:"album's slug empty"}
	//
	//}
	//
	////save albumname
	//var album models.ShopAlbum
	//album.Slug = albumslug
	//album.Name = albumname
	//album.ShopID = usex.Shop.ID
	//album.UserId = usex.UserID
	//album = rpch.SaveAlbum(album)
	//b, _ := json.Marshal(album)
	//
	//return c3mcommon.ReturnJsonMessage("1", "", "success", string(b))
	return models.RequestResult{Status: 1}

}
func doLoadalbum(usex models.UserSession) models.RequestResult {

	////get albums
	//albums := rpch.LoadAllShopAlbums(usex.Shop.ID)
	//if len(albums) == 0 {
	//	//create
	//	var album models.ShopAlbum
	//	album.Slug = "default"
	//	album.Name = "Default"
	//	album.ShopID = usex.Shop.ID
	//	album.UserId = usex.UserID
	//	album = rpch.SaveAlbum(album)
	//	albums = append(albums, album)
	//}
	//
	//b, err := json.Marshal(albums)
	//c3mcommon.CheckError("json parse doLoadalbum", err)
	//return c3mcommon.ReturnJsonMessage("1", "", "", string(b))

	return models.RequestResult{Status: 1}

}
func doEditAlbum(usex models.UserSession) models.RequestResult {
	////log.Debugf("update album ")
	//var newitem models.ShopAlbum
	//log.Debugf("Unmarshal %s", usex.Params)
	//err := json.Unmarshal([]byte(usex.Params), &newitem)
	//if !c3mcommon.CheckError("json parse page", err) {
	//	return c3mcommon.ReturnJsonMessage("0", "json parse ShopAlbum fail", "", "")
	//}
	//newitem.ShopID = usex.Shop.ID
	//newitem.UserId = usex.UserID
	//rpch.SaveAlbum(newitem)
	////log.Debugf("update album false %s", albumname)
	//return c3mcommon.ReturnJsonMessage("0", "album not found", "", "")
	return models.RequestResult{Status: 1}
}
func main() {
	//default port for service
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8902"
	}
	//open service and listen
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Printf("listening on %s\n", port)
	pb.RegisterGRPCServicesServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve : %v", err)
	}
	fmt.Print("exit")
}

//repush
