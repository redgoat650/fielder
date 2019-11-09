package buzzedsheets

import "testing"

func TestParseGender(t *testing.T) {

	input := [][]string{
		[]string{"", ""},
		[]string{"Name", "Gender"},
		[]string{"Nick Wright", "Male"},
		[]string{"Patty Barf", "Female"},
		[]string{"Not a name", ""},
	}

	nameList, err := getNameList(input)
	if err != nil {
		t.Fatal(err)
	}
	genderList, err := getGenderList(input)
	if err != nil {
		t.Fatal(err)
	}
	if nameList[0] != input[2][0] || nameList[1] != input[3][0] {
		t.Fatal("Name doesn't match")
	}
	if len(nameList) > 2 {
		t.Fatal("Too many entries in name list")
	}
	if len(genderList) > 2 {
		t.Fatal("Too many entries in the gender list")
	}
}
