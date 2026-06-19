package main

/*
2. Получение оригинальной ссылки
GET /links/{short_code}
Response
{
  "url": "https://example.com/some/very/long/url",
  "visits": 15
}
При каждом запросе необходимо увеличивать счётчик visits.
GET http://127.0.0.1:443/link/17
*/
/*
3. Получение списка ссылок
GET /links ? limit=9 & offset=3
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type InputType01 struct {
	Url string `json:"url"`
	//Name string    `json:"name"`
}
type InputType02 struct {
	Short_code int
}

type DbRowType struct {
	Short_code int
	Url        string
	Created_at time.Time
	Visits     int
}

type THandlerData struct {
	DBconn *pgx.Conn
	Hnd    func(DBconn *pgx.Conn, w http.ResponseWriter, r *http.Request)
}

var ctx = context.Background()

func (hd *THandlerData) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hd.Hnd(hd.DBconn, w, r)
}

func Parse_request_post01(DBconn *pgx.Conn, w http.ResponseWriter, r *http.Request) InputType01 {
	var DataStructure05 InputType01

	if r.Method == http.MethodGet {
		// Обработка GET-запроса
		fmt.Fprintf(w, "\nGET-запрос")

	} else if r.Method == http.MethodPost {
		// Обработка POST-запроса
		//fmt.Fprintf(w, "\nPOST-запрос")
		if r.Body != nil {
			// Тело запроса есть, значит, это POST
			body05, err05 := io.ReadAll(r.Body)
			if err05 != nil {
				fmt.Fprintf(w, "\nPOST.Body reading error: %v", err05)
			}
			BodyStr05 := string(body05)

			// Здесь можно распарсить данные из тела
			reader05 := strings.NewReader(BodyStr05)
			json_decoder05 := json.NewDecoder(reader05)
			err051 := json_decoder05.Decode(&DataStructure05)
			if err051 != nil {
				fmt.Fprintf(w, "\nError in json_decoder05.Decode: %v", err051)
			}
			/*
				fmt.Fprintf(w, "\nUrl: %s", InputType01.Url)
				//Url: rbc.ru/lk
			*/
		}
	}
	return (DataStructure05)
}

func Parse_request_get02(DBconn *pgx.Conn, w http.ResponseWriter, r *http.Request) InputType02 {
	var DataStructure05 InputType02

	if r.Method == http.MethodPost {
		// Обработка POST-запроса
		fmt.Fprintf(w, "\nPOST-запрос")

	} else if r.Method == http.MethodGet {
		// Обработка GET-запроса  http://127.0.0.1:443/links/29
		//fmt.Fprintf(w, "\nGET-запрос")

		id := r.PathValue("id")
		DataStructure05.Short_code, _ = strconv.Atoi(id)
		//fmt.Fprintf(w, "\nShort_code: %v", DataStructure05.Short_code)
		/*
			// example ok
			GET-запрос
			id: 29
			Short_code: 29
		*/
	}
	return (DataStructure05)
}

func Db_get_row_by_shcode(DBconn *pgx.Conn, Short_code_in int) (Row_out DbRowType) {
	rows := DBconn.QueryRow(ctx, "SELECT short_code, original_url, created_at, visits FROM t_links2 WHERE short_code = $1", Short_code_in)
	err := rows.Scan(&Row_out.Short_code, &Row_out.Url, &Row_out.Created_at, &Row_out.Visits)
	if err != nil {
		log.Fatalf("Error Db_get_row_by_shcode: %v", err)
	}
	return (Row_out)
}

func Db_get_new_short_code(DBconn *pgx.Conn) (Short_code_out int) {
	/////////////////////////////
	//read DB - find MAX short_code from DB - ok
	var short_code_max4 int
	err42 := DBconn.QueryRow(ctx, "SELECT COALESCE(MAX(short_code), 0) AS short_code FROM t_links2").Scan(&short_code_max4)
	if err42 != nil {
		log.Fatal("Error Db_get_new_short_code: %v", err42)
	}
	Short_code_out = short_code_max4 + 1
	return (Short_code_out)
}

func Db_insert_url(DBconn *pgx.Conn, Url_in string, Short_code_in int, Visits_in int) (RC_out int) {
	p_short_code1 := Short_code_in //Get_new_short_code(DBconn)
	p_timestamp1 := time.Now()
	p_url1 := Url_in
	p_visits1 := Visits_in

	Commandtag, err := DBconn.Exec(ctx, "INSERT INTO t_links2 (short_code, original_url, created_at, visits) VALUES ($1, $2, $3, $4)", p_short_code1, p_url1, p_timestamp1, p_visits1)
	if err != nil {
		log.Fatalf("INSERT ERROR: %v", err)
	}
	if Commandtag.RowsAffected() != 0 {
		RC_out = 0
	} else {
		RC_out = 1
	}
	return RC_out
}

func Db_update_visits(DBconn *pgx.Conn, Short_code_in int, Visits_in int) (RC_out int) {
	Commandtag, err := DBconn.Exec(ctx, "UPDATE t_links2 SET visits = $1 WHERE short_code = $2", Visits_in, Short_code_in)
	if err != nil {
		log.Fatalf("UPDATE ERROR: %v", err)
	}
	if Commandtag.RowsAffected() != 0 {
		RC_out = 0
	} else {
		RC_out = 1
	}
	return RC_out
}

func Exec_post_link(DBconn *pgx.Conn, w http.ResponseWriter, r *http.Request) {
	/*
	   1. Создание короткой ссылки
	   POST /links
	   Request
	   {"url": "https://example.com/some/very/long/url"}
	   Response
	   {"short_code": "abc123"}
	   Сервис должен:
	   — сгенерировать short_code
	   — сохранить ссылку в базу данных
	   — вернуть short_code в ответе
	*/
	///////////////////////////////
	// Parse Request
	var InputStructure0101 InputType01
	InputStructure0101 = Parse_request_post01(DBconn, w, r)

	//////////////////////////////
	// Generate SHORT_CODE
	short_code := Db_get_new_short_code(DBconn)
	visits := 0

	///////////////////////////////
	// Save Url to DB
	Db_insert_url(DBconn, InputStructure0101.Url, short_code, visits)

	////////////////////////////////////////
	// Generate JSON answer
	type answer_type0108 struct {
		Short_code string `json:"short_code"`
	}
	answer0108 := answer_type0108{Short_code: strconv.Itoa(short_code)}
	json_data0108, err0108 := json.Marshal(answer0108)
	if err0108 != nil {
		fmt.Fprintf(w, "Ошибка при сериализации 0108:", err0108)
		return
	}
	fmt.Fprintf(w, "\n%s", string(json_data0108))

}

func Exec_get_links(DBconn *pgx.Conn, w http.ResponseWriter, r *http.Request) {
	/*
		2. Получение оригинальной ссылки
		GET /links/{short_code}
		Response
		{
		  "url": "https://example.com/some/very/long/url",
		  "visits": 15
		}
		При каждом запросе необходимо увеличивать счётчик visits.
	*/
	///////////////////////////////
	// Parse Request
	var InputStructure0201 InputType02
	InputStructure0201 = Parse_request_get02(DBconn, w, r)
	//fmt.Fprintf(w, "\nInputStructure0201.Short_code= %v", InputStructure0201.Short_code)

	///////////////////////////
	//READ DATABASE with short_code
	Row_0204 := Db_get_row_by_shcode(DBconn, InputStructure0201.Short_code)
	var new_visits int
	new_visits = Row_0204.Visits + 1

	///////////////////////////
	//READ DATABASE with short_code
	/*
		var visits_0203 int
		var url_0203 string
		//var created_at_0203 time.Time
		rows_0203 := DBconn.QueryRow(ctx, "SELECT original_url, visits FROM t_links2 WHERE short_code = $1", InputStructure0201.Short_code)
		err_0203 := rows_0203.Scan(&url_0203, &visits_0203)
		if err_0203 != nil {
			log.Fatalf("Ошибка выполнения запроса 0203: %v", err_0203)
		}
		//fmt.Fprintf(w, "\nUrl= %s", url_0203)
		//fmt.Fprintf(w, "\nVisits= %d", visits_0203)
		//fmt.Fprintf(w, "\nCreated_at= %v", created_at_0203)
		var new_visits int
		new_visits = visits_0203 + 1
		//fmt.Fprintf(w, "\nNew_visits= %d", new_visits)
	*/
	////////////////////////////////
	// UPDATE DATABASE with new visits
	Db_update_visits(DBconn, InputStructure0201.Short_code, new_visits)
	//rc := update_visits_to_db(DBconn, InputStructure0201.Short_code, new_visits)
	//fmt.Fprintf(w, "\nrc=%d", rc)

	////////////////////////////////////////
	// Generate JSON answer
	//	{
	//	  "url": "https://example.com/some/very/long/url",
	//	  "visits": 15
	//	}
	type answer_type0208 struct {
		Url    string `json:"url"`
		Visits int    `json:"visits"`
	}
	//answer0208 := answer_type0208{Url: url_0203, Visits: new_visits}
	answer0208 := answer_type0208{Url: Row_0204.Url, Visits: new_visits}
	json_data0208, err0208 := json.Marshal(answer0208)
	if err0208 != nil {
		fmt.Fprintf(w, "Ошибка при сериализации 0208:", err0208)
		return
	}
	fmt.Fprintf(w, "\n%s", string(json_data0208))
	/*
		//Example
		GET http://127.0.0.1:443/links/30
		InputStructure0201.Short_code= 30
		Url= https://dzen.ru/pictures
		Visits= 2
		Created_at= 2026-06-18 15:41:36.577092 +0300 MSK
		New_visits= 3
		rc=0
		{"url":"https://dzen.ru/pictures","visits":3}
	*/
}

func main() {

	///////////////////////////////////
	//POSTGRES DATABASE CONNECTION
	db_link_psw := os.Getenv("DB_LINK_PSW")
	if db_link_psw == "" {
		fmt.Println("Ошибка: переменная окружения DB_LINK_PSW не задана")
		return
	}
	db_link_usr := os.Getenv("DB_LINK_USR")
	if db_link_usr == "" {
		fmt.Println("Ошибка: переменная окружения DB_LINK_USR не задана")
		return
	}
	// Строка подключения
	connStr := "postgres://" + db_link_usr + ":" + db_link_psw + "@localhost:5432/links"

	// Set connection
	conn1, err1 := pgx.Connect(ctx, connStr)
	if err1 != nil {
		log.Fatalf("\nНе удалось подключиться к БД: %v", err1)
	} else {
		//DBconn = conn1
		fmt.Println("\nConnected to DB 'links' ok")
	}
	defer conn1.Close(context.Background())

	////////////////////////////////////////////
	// Handler 01 - 'POST /link'
	my_addr := "127.0.0.1:443"
	my_path01 := "/link"
	handl01 := THandlerData{DBconn: conn1, Hnd: Exec_post_link}
	http.HandleFunc(my_path01, handl01.ServeHTTP)
	fmt.Println("\nСервис запущен " + my_addr + my_path01)

	////////////////////////////////////////////
	// Handler 02 - 'GET /links'
	my_path02 := "/links/{id}"
	handl02 := THandlerData{DBconn: conn1, Hnd: Exec_get_links}
	http.HandleFunc(my_path02, handl02.ServeHTTP)
	fmt.Println("\nСервис запущен " + my_addr + my_path02)

	///////////////////////////////////////////
	// RUN SERVER
	//my_addr := "127.0.0.1:443"
	//fmt.Println("\nСервер запущен " + my_addr + my_path)
	err2 := http.ListenAndServe(my_addr, nil)
	if err2 != nil {
		log.Fatalf("\nError starting the server %s : %v", my_addr, err2)
	}
}
