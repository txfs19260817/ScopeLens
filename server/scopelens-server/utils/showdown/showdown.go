package showdown

import (
	"bufio"
	"github.com/sugoiuguu/showgone"
	"image/png"
	"os"
	"scopelens-server/config"
	"scopelens-server/utils/file"
	"scopelens-server/utils/showdown/utils"
	"strings"
)

// Parse Showdown text to struct
func ShowdownParser(text string) []utils.Pokemon {
	pms := []utils.Pokemon{}
	r := bufio.NewReader(strings.NewReader(text))
	for {
		poke, err := showgone.Parse(r)
		if err != nil {
			break
		}
		var tp []string
		if t, ok := utils.Poke2Types[utils.String2Filename(string(poke.Species))]; ok {
			tp = t
		}
		pm := utils.Pokemon{
			Name:    string(poke.Species),
			Type:    tp,
			Item:    string(poke.Item),
			Ability: string(poke.Ability),
			Moves:   []string{string(poke.Moves[0]), string(poke.Moves[1]), string(poke.Moves[2]), string(poke.Moves[3])},
		}
		pms = append(pms, pm)
	}
	return pms
}

// from Pokemon struct to rental team preview image
func RentalTeamMaker(text string) (string, error) {
	// parse showdown text
	pms := ShowdownParser(text)

	// Read background image from file that already exists
	bgImageFile, err := os.Open(utils.SpritePath + "bg.png")
	if err != nil {
		return "", err
	}
	defer bgImageFile.Close()

	// Since we know it is a png already, call png.Decode()
	bg, err := png.Decode(bgImageFile)
	if err != nil {
		return "", err
	}

	bg, err = utils.AppendPokemon(bg, &pms)
	if err != nil {
		return "", err
	}

	bg, err = utils.AppendItems(bg, &pms)
	if err != nil {
		return "", err
	}

	bg, err = utils.AppendInfo(bg, &pms)
	if err != nil {
		return "", err
	}

	bg, err = utils.AppendTypes(bg, &pms)
	if err != nil {
		return "", err
	}

	bg, err = utils.AppendMoves(bg, &pms)
	if err != nil {
		return "", err
	}

	// Save an image
	path := config.Path.ImageSavePath + file.Rename(text[11:45]+".png")
	err = utils.SaveImage(bg, path)
	if err != nil {
		return "", err
	}
	return path, nil
}
