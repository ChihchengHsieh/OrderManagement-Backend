package main

import (
	"orderFunc/apis"
	"orderFunc/databases"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	databases.InitDB()

	router := apis.InitRouter()

	// testFunc.AllTestingRun()

	// models.UpdateMemberByID()

	router.Run()
	// type testingStruct struct {
	// 	FieldA string
	// 	FieldB string
	// }

	// testArr := []testingStruct{
	// 	testingStruct{"Hello", "Good"},
	// 	testingStruct{"Second", "SecondB"},
	// }

	// testArrJSON, err := json.Marshal(testArr)

	// if err != nil {
	// 	fmt.Printf("Error: %+v", err)

	// } else {
	// 	fmt.Println(string(testArrJSON))
	// }

	// var testMap []interface{}

	// err = json.Unmarshal(testArrJSON, &testMap)

	// if err != nil {
	// 	fmt.Printf("Error: %+v", err)

	// } else {
	// 	fmt.Println(testMap)
	// }

}
