package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// hämtar högsta mongodb passport _id
func GetHighestItemID(client *mongo.Client, dbName, collectionName string) (int, error) {
	collection := client.Database(dbName).Collection(collectionName)

	var result struct {
		ItemID int `bson:"itemid"`
	}

	options := options.FindOne().SetSort(bson.D{{"ItemID", -1}})

	err := collection.FindOne(context.TODO(), bson.D{}, options).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//returnerar _id som 0 om det inte finns några existerande passports
			return 0, nil
		}
		//felhantering om queryn misslyckar att hämta _id
		log.Println("Error retrieving highest itemid:", err)
		return 0, err
	}

	//returnerar högsta _id samt nil om error
	return result.ItemID, nil
}

// TODO: Ändra så att funktionen tar query parametrar istället för hårdkodad data
func Createpassport(ItemN string, OriginN string, client *mongo.Client, database, collection string, SensitiveArray []string, LinkMadeFromN []string, LinkMakesN []string) (itemID int) {
	//funktionsanrop för att hämta det nuvarande högsta mongodb passport _id i databasen
	highestItemID, err := GetHighestItemID(client, database, collection)
	if err != nil {
		log.Fatal("Error getting highest itemid:", err)
	}
	log.Println("Highest ItemID:", highestItemID)
	now := time.Now()

	//Hämtar PassPort struct i models och ger den värden
	Passport := PassPort{
		ItemID:       highestItemID + 1,
		ItemName:     ItemN,
		Origin:       OriginN,
		LinkMadeFrom: LinkMadeFromN, //Ska matas in länk från IPFS som ska stores
		LinkMakes:    LinkMakesN,    //Samma här gäller det.
		LinkEvents:   []string{},
		Sensitive:    SensitiveArray,
		CreationDate: now.Format("01-02-2006"),
	}

	//skickar det nyskapade passport till databas
	Coll := client.Database(database).Collection(collection)
	var ctx = context.TODO()
	insertResult, err := Coll.InsertOne(ctx, Passport)
	if err != nil {
		log.Fatal(err)
	}

	//ser till att vi disconnectar från databasen även om ett error skulle förekomma vid insert till databas
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	fmt.Println(insertResult)
	return highestItemID + 1

}

func sensetiveArray() (sensitiveArray []string) {
	var input string
	sensitiveArray = []string{"0", "0", "0", "0"}
	fmt.Println("Enter sensetivite value 0 not sensetive : 1 sensetive")
	fmt.Print("LinkMadeFrom: ")
	fmt.Scan(&input)
	for input != "0" && input != "1" {
		fmt.Println("Input must be 0 or 1")
		fmt.Print("LinkMadeFrom: ")
		fmt.Scan(&input)
	}
	sensitiveArray = append(sensitiveArray, input)
	fmt.Print("LinkMakes: ")
	fmt.Scan(&input)
	for input != "0" && input != "1" {
		fmt.Println("Input must be 0 or 1")
		fmt.Print("LinkMakes: ")
		fmt.Scan(&input)
	}
	sensitiveArray = append(sensitiveArray, input)
	// LinkEvents special
	sensitiveArray = append(sensitiveArray, "2")
	// Sensetive
	sensitiveArray = append(sensitiveArray, "1")
	// CreationDate
	sensitiveArray = append(sensitiveArray, "0")
	return sensitiveArray
}

func LinkMadeFrom() (LinkMadeFrom []string) {
	var CID, inputMore, linkPassport string
	fmt.Println("Press 1 to start entering CIDs for LinkMadeFrom: ")
	fmt.Scan(&inputMore)
	for inputMore == "1" {
		fmt.Println("Enter CID (Enter 0 if no more): ")
		fmt.Scan(&CID)
		if CID != "0" {
			linkPassport = passportFromCID(CID)

			LinkMadeFrom = append(LinkMadeFrom, linkPassport)
		} else {
			inputMore = "0"
		}
	}
	return LinkMadeFrom
}

// Funtion som tar in hårdkodad objectid för tillfället och gör det möjligt att lägga till event som hänt med produkten.
// Behöver lägga till där man hämtar objectid för att välja vilken produkt som det ska uppdateras för
func RemanufactureEvent(client *mongo.Client, database, collection, mongoid string, RemanEvent string) {
	Coll := client.Database(database).Collection(collection)
	id, _ := primitive.ObjectIDFromHex(mongoid)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "LinkEvents", Value: RemanEvent}}}}
	result, err := Coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Documents matched: %v\n", result.MatchedCount)
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)
}
func passportMenu(client *mongo.Client, database, collection string) (itemID int) {
	//temporär input för test ändamål, ska ändras framöver för att kunna göras via hemsida/program etc
	var i int
	fmt.Println("What do you want to do? 1: Createpassport, 2: Remanufacture events for passports")
	fmt.Scan(&i)
	switch i {
	case 1:

		//testinput av item name samt item origin
		var ItemN, OriginN string
		fmt.Println("Enter item name : ")
		fmt.Scan(&ItemN)
		fmt.Println("Enter item origin : ")
		fmt.Scan(&OriginN)
		// var LinkMadeFrom []string
		LinkMadeFrom := LinkMadeFrom()
		sensitiveArray := sensetiveArray()
		LinkMakes := []string{}

		//funktionsanrop för att skapa passport.
		//TODO: ska kunna hantera querys senare
		return Createpassport(ItemN, OriginN, client, database, collection, sensitiveArray, LinkMadeFrom, LinkMakes)
	case 2:

		//testinput för att lägga till ett remanufacture event till en produkt
		fmt.Println("Enter what has been updated on this certain product:")
		fmt.Scan("")
		reader := bufio.NewReader(os.Stdin)
		RemanEvent, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		RemanEvent = RemanEvent[:len(RemanEvent)-1]

		//funktionsanrop för att uppdatera en produkt med ett remanufacture event
		//TODO: andra variabeln som skickas med i funktionen måste bytas ut med en dynamisk variabel "objectid" senare, är hårdkodad för nuvarandet med ett _id
		//TODO: ska kunna hantera querys
		var remanafactureProductID string = "65b1282112afb84376254117"
		RemanufactureEvent(client, database, collection, remanafactureProductID, RemanEvent)

	default:
		fmt.Println("xdd")

	}
	return
}
