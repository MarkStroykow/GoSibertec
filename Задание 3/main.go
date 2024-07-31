package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// json [1, 2, 3, 4.5, 23]

type Conf struct {
	URL string `json:"url"`
}

func readjson(typef, filename string) ([]float64, error) {
	var data []byte
	var err error

	//2.2. Принимать через аргументы командной строки параметр для определения источника данных (файл или stdin).
	switch typef {
	case `file`:
		data, err = os.ReadFile(filename + `.json`)
		if err != nil {
			return nil, err
		}
	case `stdin`:
		fmt.Println(`[1, 2, 3] as JSON`)
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		data = []byte(text)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(`invalid type: %s`, typef)
	}

	var nums []float64
	err = json.Unmarshal(data, &nums)
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

// 1.4. Логирует результаты каждого шага в файл.
func logs(input interface{}, logname string) {
	file, err := os.OpenFile(logname+`.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	infoLog := log.New(file, ``, log.Ldate|log.Ltime)
	infoLog.Println(input)

}

// 2.4. Поддерживать конфигурацию через файл настроек или переменные окружения (например, URL для HTTP запроса).
func readConfig(file string) (Conf, error) {
	var data []byte
	var err error
	data, err = os.ReadFile(file + `.json`)
	if err != nil {
		return Conf{}, fmt.Errorf(`failed read config %v`, err)
	}

	var c Conf
	err = json.Unmarshal(data, &c)
	if err != nil {
		return Conf{}, fmt.Errorf(`faild unmarshal %v`, err)
	}

	return c, nil
}

func tryGET(url string) (int, error) {
	responce, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer responce.Body.Close()

	if responce.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(`non-200: %d`, responce.StatusCode)
	}

	return responce.StatusCode, nil
}

func output(input, namefile string) error {
	file, err := os.OpenFile(namefile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf(`error opening file: %v`, err)
	}

	defer file.Close()

	_, err = file.WriteString(input + "\n")
	if err != nil {
		return fmt.Errorf(`err writing file: %v`, err)
	}

	return nil
}

func main() {
	Type := flag.String("t", "file", "file/stdin")
	filename := flag.String("n", "text", "name")
	Logname := flag.String("l", "info", "logfile")
	Outputfile := flag.String("o", "output", "logfile")
	Config := flag.String("c", "config", "config(url)")
	flag.Parse()

	logs(fmt.Sprintf(`%s %s.json %s.log %s %s.json`, *Type, *filename, *Logname, *Outputfile, *Config), *Logname)
	//1.1. Читает из файла JSON с массивом чисел.
	//2.1. Читать данные не только из файла, но и из стандартного ввода.
	nums, err := readjson(*Type, *filename)
	if err != nil {
		logs(fmt.Sprintf(`cant read. %v`, err), *Logname)
		return
	}

	conf, err := readConfig(*Config)
	if err != nil {
		logs(fmt.Sprintf(`failed read %v`, err), *Logname)
		return
	}

	//1.3. Выполняет HTTP GET запрос на заданный URL и проверяет статус ответа (должен быть 200).
	code, err := tryGET(conf.URL)
	if err != nil {
		logs(fmt.Sprintf(`%v`, err), *Logname)
		return
	}

	//1.2. Считает сумму всех чисел в массиве.
	res := sumn(nums)

	logs(fmt.Sprintf(`Nums %v`, nums), *Logname)
	logs(fmt.Sprintf(`Sum %v`, res), *Logname)
	logs(fmt.Sprintf(`Get %v with code %d`, conf.URL, code), *Logname)

	//2.3. Сохранять результат работы в указанный пользователем файл, а не только в стандартный вывод.
	err = output(fmt.Sprintf(`Nums %v`, nums), *Outputfile)
	if err != nil {
		logs(fmt.Sprintf(`%v`, err), *Logname)
		return
	}

	err = output(fmt.Sprintf(`Sum %v`, res), *Outputfile)
	if err != nil {
		logs(fmt.Sprintf(`%v`, err), *Logname)
		return
	}

	err = output(fmt.Sprintf(`Get %v with code %d`, conf.URL, code), *Outputfile)
	if err != nil {
		logs(fmt.Sprintf(`%v`, err), *Logname)
		return
	}

}
