package game

import (
	"db_client"
	"fmt"
	"gmap"
	"playchar"
	"strconv"
	"strings"
	"time"
	"math/rand"
)

type Game struct {
	Id         int `json:"game-id,omitempty"`
	Map_master *gmap.Gmap
	Gg         *playchar.Playchar
	Status     int // -1(over), 0(saved), 1(online)
	SavedName  string
}

func MakeZeroGame(mapSize int) (*Game, error) {
	//globalGameNum++
	newmap := gmap.MakeZeroMap(mapSize)
	var id int
	newCharacter := playchar.New(100, 0, 0)
	newgame := Game{id, newmap, newCharacter, 1, ""}
	return &newgame, nil
}

func MakeAGame(mapSize int, gameName string, eventNum int, dbc *db_client.DBClient) (*Game, error) {
	//globalGameNum++
	newmap := gmap.MakeAMap(mapSize)
	newmap.MapEventRandomizator(eventNum, mapSize)
	newCharacter := playchar.New(100, 0, 0)
	var id int
	strtoexec := fmt.Sprintf("INSERT INTO games (status,map,saved_name,playchar) VALUES (%d,%s,'%s',%s) RETURNING id", 1, newmap.InsertString(), gameName, newCharacter.ToString())
	res, err := dbc.DB.Query(strtoexec)
	defer res.Close()
	if err != nil {
		return nil, fmt.Errorf("MakeAGame: failed to insert into games %s", err)
	}
	res.Next()
	err = res.Scan(&id)

	newgame := Game{id, newmap, newCharacter, 1, ""}

	return &newgame, nil
}

func UpdateTheGame(GameToUpdate *Game, dbc *db_client.DBClient) error {

	strtoexec := fmt.Sprintf("update games set map=%s, playchar=%s, status=%d, saved_name='%s'  where id=%d", GameToUpdate.Map_master.InsertString(), GameToUpdate.Gg.ToString(), GameToUpdate.Status, GameToUpdate.SavedName, GameToUpdate.Id)
	res, err := dbc.DB.Query(strtoexec)
	defer res.Close()
	if err != nil {
		return fmt.Errorf("updateAGame: failed to update into games %s", err)
	}

	return nil
}

func GetTheGame(gameId int, mapSize int, dbc *db_client.DBClient) (*Game, error) {
	loadgame, _ := MakeZeroGame(mapSize)
	loadgame.Id = gameId
	loadgame.Map_master.Params = gmap.MapParams{mapSize, mapSize}
	var uint8buf []uint8

	row, err := dbc.DB.Query("SELECT status FROM games WHERE id=$1", gameId)
	defer row.Close()
	row.Next()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	err = row.Scan(&uint8buf)
	var strbuf string
	strbuf = fmt.Sprintf("%s", uint8buf)
	loadgame.Status, _ = strconv.Atoi(strbuf)

	row, err = dbc.DB.Query("SELECT saved_name FROM games WHERE id=$1", gameId)
	row.Next()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	err = row.Scan(&uint8buf)
	strbuf = fmt.Sprintf("%s", uint8buf)
	loadgame.SavedName = strbuf

	row, err = dbc.DB.Query("SELECT map FROM games WHERE id=$1", gameId)
	row.Next()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	err = row.Scan(&uint8buf)
	strbuf = fmt.Sprintf("%s", uint8buf)
	//strbuf = strbuf[1:len(strbuf)-1]
	strbuf = strings.Replace(strbuf, "{", "", -1)
	strbuf = strings.Replace(strbuf, "}", "", -1)
	strList := strings.Split(strbuf, ",")
	loadgame.Map_master = gmap.MakeAMap(mapSize)

	//cellrow := []Cell.Cell{}
	//cell := Cell.Cell{}
	iter := 0
	i := 0
	for _, s := range strList {
		strKSL := strings.Split(s, ":") // kind and status
		strKSLi := []int{}
		for _, sksl := range strKSL {
			ks, _ := strconv.Atoi(sksl)
			strKSLi = append(strKSLi, ks)
		}
		loadgame.Map_master.Field[iter][i].Kind = strKSLi[0]
		loadgame.Map_master.Field[iter][i].Hidden = strKSLi[1]
		i++
		if i == mapSize {
			iter++
			i=0
		}
	}

	row, err = dbc.DB.Query("SELECT playchar FROM games WHERE id=$1", gameId)
	row.Next()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	err = row.Scan(&uint8buf)
	strbuf = fmt.Sprintf("%s", uint8buf)
	strList = strings.Split(strbuf, ":")
	loadgame.Gg.Healthpoints, _ = strconv.Atoi(strList[0])
	loadgame.Gg.Position.Posx, _ = strconv.Atoi(strList[1])
	loadgame.Gg.Position.Posy, _ = strconv.Atoi(strList[2])

	return loadgame, nil
}

func (game Game) MapEventRandomizator(eventsNum int, mapSize int) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i, row := range game.Map_master.Field {
		for j, _ := range row {
			if game.Map_master.Field[i][j].Kind != 0 {
				//cel.Kind = r.Perm(eventsNum)[0] + 1
				randn := 1 + randInt(0, eventsNum)
				game.Map_master.Field[i][j].Kind = randn
			}
		}
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
