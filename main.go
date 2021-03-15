package main

import (
	"fmt"
	"time"
	"math/rand"
	"runtime"
	"strconv"
)

type countHolder struct {
	royalFlush float64
	straightFlush float64
	fourofaKind float64
	fullHouse float64
	Flush float64
	straight float64
	threeofaKind float64 
	twoPair float64 
	onePair float64
	highCard float64
}

var SUITS = [4]string{"Hearts", "Diamonds", "Clubs", "Spades"}
var VALUES = [13]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King",  "Ace"}

func main() {
	numThreads := runtime.NumCPU()
	fmt.Printf("Card Program- Number of CPUs: %d\n", numThreads)
		
	var numHands int = 2598960
	fmt.Println("Num Hands: ", numHands)
	fmt.Println("Num Hands per Thread: ", numHands/numThreads)
	var count chan countHolder
	count = make(chan countHolder)

	for i := 0; i < numThreads; i++ {
		go getProbabilities(count, numHands)
	}

	var countStruct countHolder
	for i := 0; i < numThreads; i++ {
		var temp countHolder
		temp = <- count
		countStruct.royalFlush += temp.royalFlush
		countStruct.straightFlush += temp.straightFlush
		countStruct.fourofaKind += temp.fourofaKind
		countStruct.fullHouse += temp.fullHouse
		countStruct.Flush += temp.Flush
		countStruct.straight += temp.straight
		countStruct.threeofaKind += temp.threeofaKind
		countStruct.twoPair += temp.twoPair
		countStruct.onePair += temp.onePair
		countStruct.highCard += temp.highCard
	}

 fmt.Printf("Royal Flush (count: %d): %f\n", int(countStruct.royalFlush), (countStruct.royalFlush/float64(numHands)) * 100.0)
 fmt.Printf("Straight Flush (count: %d): %f\n", int(countStruct.straightFlush), (countStruct.straightFlush/float64(numHands)) * 100.0)
 fmt.Printf("Four of a Kind (count: %d): %f\n", int(countStruct.fourofaKind), (countStruct.fourofaKind/float64(numHands)) * 100.0)
 fmt.Printf("Full House (count: %d): %f\n", int(countStruct.fullHouse), (countStruct.fullHouse/float64(numHands)) * 100.0)
 fmt.Printf("Flush (count: %d): %f\n", int(countStruct.Flush), (countStruct.Flush/float64(numHands)) * 100.0)
 fmt.Printf("Straight (count: %d): %f\n", int(countStruct.straight), (countStruct.straight/float64(numHands)) * 100.0)
 fmt.Printf("Three of a Kind (count: %d): %f\n", int(countStruct.threeofaKind), (countStruct.threeofaKind/float64(numHands)) * 100.0)
 fmt.Printf("Two Pair (count: %d): %f\n", int(countStruct.twoPair), (countStruct.twoPair/float64(numHands)) * 100.0)
 fmt.Printf("One Pair (count: %d): %f\n", int(countStruct.onePair), (countStruct.onePair/float64(numHands)) * 100.0)
 fmt.Printf("High Card (count: %d): %f\n", int(countStruct.highCard), (countStruct.highCard/float64(numHands)) * 100.0)

}

func getProbabilities(count chan countHolder, numHands int) {
	numThreads := runtime.NumCPU()

	var seed = time.Now().UnixNano()
	var randGen = rand.New(rand.NewSource(seed))

	var countStruct countHolder
	for i := 0; i < numHands/numThreads; i++ {
		fullHand := getCards(5, randGen)

		mapBySuits := splitBySuit(fullHand)
		mapByValues := splitByValue(fullHand)


		if isRoyalFlush(mapBySuits, mapByValues) {
			countStruct.royalFlush = countStruct.royalFlush + 1
		} else if isStraightFlush(mapBySuits, mapByValues) {
			countStruct.straightFlush = countStruct.straightFlush + 1
		} else if isFourofaKind(mapBySuits, mapByValues){
			countStruct.fourofaKind = countStruct.fourofaKind + 1
		} else if isFullHouse(mapBySuits, mapByValues) {
			countStruct.fullHouse = countStruct.fullHouse + 1
		} else if isFlush(mapBySuits, mapByValues) {
			countStruct.Flush = countStruct.Flush + 1
		} else if isStraight(mapBySuits, mapByValues){
			countStruct.straight = countStruct.straight + 1
		} else if isThreeofaKind(mapBySuits, mapByValues){
			countStruct.threeofaKind = countStruct.threeofaKind + 1
		} else if isTwoPair(mapBySuits, mapByValues) {
			countStruct.twoPair = countStruct.twoPair + 1
		} else if isOnePair(mapBySuits, mapByValues) {
			countStruct.onePair = countStruct.onePair + 1
		} else{
			countStruct.highCard = countStruct.highCard + 1
		}

	}
	count <- countStruct
}
func isRoyalFlush(suits map[string][]int, values map[string][]int) bool { //check for this first
	var flag bool = true
	if len(values["10"]) == 0 {
		flag = false
	}else if len(values["Jack"]) == 0 {
		flag = false
	}else if len(values["Queen"]) == 0 {
		flag = false
	}else if len(values["King"]) == 0 {
		flag = false
	}else if len(values["Ace"]) == 0 {
		flag = false
	}
	if (len(suits["Clubs"]) != 5) && (len(suits["Diamonds"]) != 5) && (len(suits["Hearts"]) != 5) && (len(suits["Spades"]) != 5) {
		flag = false
	}
	return flag
}
func isStraightFlush(suits map[string][]int, values map[string][]int) bool { //check for this second
	var flag bool = true
	keys := make([]string, len(values))
	for i := 0; i < 9; i++ {
		keys[i] = strconv.Itoa(i + 2)
	}
	keys[9] = "Jack"
	keys[10] = "Queen"
	keys[11] = "King"
	keys[12] = "Ace"
	var cardVals []int
	for L, k := range keys {
		if len(values[k]) > 0 {
			cardVals = append(cardVals, L + 2)
		}
	}
	if len(cardVals) != 5 {
		flag = false
	} else if (cardVals[1] -  cardVals[0]) != 1 || (cardVals[2] -  cardVals[1]) != 1 || (cardVals[3] -  cardVals[2]) != 1 || (cardVals[4] -  cardVals[3]) != 1 {
		flag = false
	}
	if len(values["Ace"]) == 1 && len(values["2"]) == 1 && len(values["3"]) == 1 && len(values["4"]) == 1 && len(values["5"]) == 1 {//ace is low edge case
		flag = true
	}
	if (len(suits["Clubs"]) != 5) && (len(suits["Diamonds"]) != 5) && (len(suits["Hearts"]) != 5) && (len(suits["Spades"]) != 5) {
		flag = false
	}
	return flag
}
func isFourofaKind(suits map[string][]int, values map[string][]int) bool { 
	var flag bool = false
	for _, element := range values {
		if len(element) == 4 {
			flag = true
			break
		}
	}
	return flag
}
func isFullHouse(suits map[string][]int, values map[string][]int) bool { 
	var flag bool = false
	if isThreeofaKind(suits, values) && isOnePair(suits, values)  {
		flag = true
		}
	return flag
}
func isFlush(suits map[string][]int, values map[string][]int) bool { 
	var flag bool = false
	for _, element := range suits {
		if len(element) == 5 {
			flag = true
			break
		}
	}
	return flag
}
func isStraight(suits map[string][]int, values map[string][]int) bool { 
	var flag bool = true
	keys := make([]string, len(values))
	for i := 0; i < 9; i++ {
		keys[i] = strconv.Itoa(i + 2)
	}
	keys[9] = "Jack"
	keys[10] = "Queen"
	keys[11] = "King"
	keys[12] = "Ace"
	var cardVals []int
	for L, k := range keys {
		if len(values[k]) > 0 {
			cardVals = append(cardVals, L + 2)
		}
	}
	if len(cardVals) != 5 {
		flag = false
	} else if (cardVals[1] -  cardVals[0]) != 1 || (cardVals[2] -  cardVals[1]) != 1 || (cardVals[3] -  cardVals[2]) != 1 || (cardVals[4] -  cardVals[3]) != 1 {
		flag = false
	}
	if len(values["Ace"]) == 1 && len(values["2"]) == 1 && len(values["3"]) == 1 && len(values["4"]) == 1 && len(values["5"]) == 1 {//ace is low edge case
		flag = true
	}
	return flag
}
func isThreeofaKind(suits map[string][]int, values map[string][]int) bool { 
	var flag bool = false
	for _, element := range values {
		if len(element) == 3 {
			flag = true
			break
		}
	}
	return flag
}
func isTwoPair(suits map[string][]int, values map[string][]int) bool { 
	var flag bool = false
	var mainkey1 string = "no val"
	for key, element := range values {
		if len(element) == 2 {
			mainkey1 = key
			break
		}
	}
	for key, element := range values {
		if (len(element) == 2) && (mainkey1 != key) {
			flag = true
			break
		}
	}
	return flag
}
func isOnePair(suits map[string][]int, values map[string][]int) bool { 
	var flag bool = false
	for _, element := range values {
		if len(element) == 2 {
			flag = true
			break
		}
	}
	return flag
}
//--
func getNumericValue(card int) int {
	return card % 13
}
//--
func getStringValue(card int) string {
	return VALUES[getNumericValue(card)]
}
//--
func getSuit(card int) string {
	return SUITS[card / 13]
}
//--
func getCards(numCards int, randGen *rand.Rand) []int {
	hand := make([]int, numCards)

	cards := make([]int, 52)
	for i := 0;i < 52;i++ {
		cards[i] = i
	}

	for i := 0;i < numCards;i++ {
		randPos := randGen.Intn(len(cards))
		hand[i] = cards[randPos]
		cards[randPos] = cards[len(cards) - 1]
		cards = cards[:len(cards) - 1]
	}
	sortCards(hand)
	return hand
}
//--
func sortCards(hand []int) {	
	for i := 0;i < len(hand) - 1;i++ {
		for j := 0;j < len(hand) - 1;j++ {
			if getNumericValue(hand[j]) < getNumericValue(hand[j + 1]) {
				temp := hand[j]
				hand[j] = hand[j + 1]
				hand[j + 1] = temp
			}
		}
	}
}
//--
func printCards(cards []int) {
	//go through all the cards
	for _, card := range cards {
		//print each card
		printCard(card)
		fmt.Println()
	}
}
//--
func printCard(card int) {
	//get the suit and value
	value := getStringValue(card)
	suit := getSuit(card)

	//print the suit and value
	fmt.Printf("%s of %s", value, suit)
}
//--
func splitBySuit(hand []int) map[string][]int {
	suitMap := make(map[string][]int)

	for _, suit := range SUITS {
		suitMap[suit] = make([]int, 0, 13)		
	}

	for _, card := range hand {
		theCardsSuit := getSuit(card)
		for _, suit := range SUITS {
			if theCardsSuit == suit {
				suitMap[suit] = append(suitMap[suit], card)
				break
			}
		}
	}

	return suitMap
}
//--
func splitByValue(hand []int) map[string][]int {
	valueMap := make(map[string][]int)

	for _, value := range VALUES {
		valueMap[value] = make([]int, 0, 4)
	}

	for _, card := range hand {
		theCardsValue := getStringValue(card)
		for _, value := range VALUES {
			if theCardsValue == value {
				valueMap[value] = append(valueMap[value], card)
				break
			}
		}
	}
	return valueMap
}