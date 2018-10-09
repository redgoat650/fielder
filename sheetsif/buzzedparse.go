package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	fielder "github.com/redgoat650/fielder/scheduling"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("../credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := "1y-CzFJ9HgrW1M_nAn6HG8dBzU0PXVJ9hYpS62XnV1tI"

	fmt.Println(getPlayerInfo(srv, spreadsheetId))
}

func getPlayerInfo(srv *sheets.Service, ssid string) ([]*fielder.Player, error) {

	players := make([]*fielder.Player, 0)

	readRange := "Weekly Roster!A3:D23"
	resp, err := srv.Spreadsheets.Values.Get(ssid, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		panic("No name data found in range.")
	} else {
		for _, row := range resp.Values {

			// Print columns A and E, which correspond to indices 0 and 4.
			fullName := row[0]
			fmt.Printf("%s\n", fullName)

			splitname := strings.Fields(fullName.(string))

			first := ""
			last := ""
			if len(splitname) > 0 {
				first = splitname[0]
			}

			if len(splitname) > 1 {
				last = splitname[1]
			}

			genderStr := row[1]

			var gender fielder.PlayerGender
			switch genderStr {
			case "male":
				gender = fielder.MaleGender
			case "female":
				gender = fielder.FemaleGender
			}

			player := fielder.NewPlayer(first, last, gender)

			fmt.Println(row, len(row))
			if len(row) > 2 {
				player.Phone = row[2].(string)
			}

			if len(row) > 3 {
				player.Email = row[3].(string)
			}

			players = append(players, player)
		}
	}

	readRange = "Preferred Positions!A2:D22"
	resp, err = srv.Spreadsheets.Values.Get(ssid, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		panic("No pref data found in range.")
	}

	for _, row := range resp.Values {
		name := row[0]
		names := strings.Fields(name.(string))
		first := names[0]
		last := names[1]

		for _, player := range players {
			if player.FirstName == first {
				if player.LastName == last {
					//Found player - update their prefs
					prefs := row[1:]

					for i, pref := range prefs {
						mag := 3 - i
						positions := pref2pos(pref.(string))
						for _, pos := range positions {
							player.Pref[pos] = mag
						}
					}

				}
			}

		}

	}

	return players, nil

}

func pref2pos(pref string) []fielder.Position {

	switch pref {

	case "Third":
		return []fielder.Position{fielder.Third}
	case "First":
		return []fielder.Position{fielder.First}
	case "Second":
		return []fielder.Position{fielder.Second}
	case "Pitcher":
		return []fielder.Position{fielder.Pitcher}
	case "Catcher":
		return []fielder.Position{fielder.Catcher}
	case "Right Shortstop":
		return []fielder.Position{fielder.RShort}
	case "Left Shortstop":
		return []fielder.Position{fielder.LShort}
	case "Right Center":
		return []fielder.Position{fielder.RCenter}
	case "Left Center":
		return []fielder.Position{fielder.LCenter}
	case "Right Field":
		return []fielder.Position{fielder.RField}
	case "Left Field":
		return []fielder.Position{fielder.LField}

	default:
		return []fielder.Position{}
	}
}
