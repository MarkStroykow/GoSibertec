package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// json [1, 2, 3, 4.5, 23]

func readjson(name string) ([]float64, error) {
	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var nums []float64
	err = json.Unmarshal(file, &nums)

	if err != nil {
		return nil, err
	}
	return nums, nil
}

func sumn(nums []float64) float64 {
	var sum float64
	for _, n := range nums {
		sum += n
	}
	return sum
}

// 4. Логирует результаты каждого шага в файл.
func logs(input interface{}) {
	file, err := os.OpenFile(`info.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	infoLog := log.New(file, ``, log.Ldate|log.Ltime)
	infoLog.Println(input)

}

func main() {
	logs(``)

	//1. Читает из файла JSON с массивом чисел.
	namejson := `text.json`
	nums, err := readjson(namejson)
	if err != nil {
		logs(fmt.Sprintf(`Can't read JSON: %v`, err))
		return
	}

	//3. Выполняет HTTP GET запрос на заданный URL и проверяет статус ответа (должен быть 200).
	url := `https://www.youtube.com`
	response, err := http.Get(url)
	if err != nil {
		logs(fmt.Sprintf(`Non-200 %v`, err))
		return
	}
	logs(fmt.Sprintf(`Nums %v`, nums))

	//2. Считает сумму всех чисел в массиве.
	res := sumn(nums)
	logs(fmt.Sprintf(`Sum %v`, res))
	logs(fmt.Sprintf(`Get %v %v`, url, response.StatusCode))
	fmt.Printf("Nums: %v\nSum: %v\n%s", nums, res, fmt.Sprintf(`Get %v %v`, url, response.StatusCode))
}
