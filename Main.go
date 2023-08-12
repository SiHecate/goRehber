package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func DatabaseBaglanma() (*sql.DB, error) {
	databaseBilgileri := "user=postgres password=393406 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", databaseBilgileri)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Database'e giriş başarılı.")
	return db, nil
}

func rehberOluştur(db *sql.DB, tabloAdi string) {
	SQLoluştur := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			SiraNo SERIAL PRIMARY KEY,
			isim VARCHAR(50),
			soyisim VARCHAR(50),
			telefon INT
		);
	`, tabloAdi)

	_, err := db.Exec(SQLoluştur)
	if err != nil {
		log.Fatal(err)
	}
}

func tabloAdlariCek(db *sql.DB) ([]string, error) {
	SQLlistele := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';"

	rows, err := db.Query(SQLlistele)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tabloAdlari []string
	for rows.Next() {
		var tabloAdi string

		if err := rows.Scan(&tabloAdi); err != nil {
			return nil, err
		}
		tabloAdlari = append(tabloAdlari, tabloAdi)
	}

	return tabloAdlari, nil
}

func rehberKaldir(db *sql.DB, tabloAdi string) error {
	SQLkaldir := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tabloAdi)
	_, err := db.Exec(SQLkaldir)
	if err != nil {
		return err
	}
	fmt.Printf("'%s' adlı tablo silindi.\n", tabloAdi)
	return nil
}

//INSERT INTO umutcann (isim, soyisim, telefon) VALUES ('Ad', 'Soyad', Numara);

func rehberBilgiEkle(db *sql.DB, tabloAdi string, isim string, soyisim string, telefon int) error {
	SQLbilgiEkle := fmt.Sprintf("INSERT INTO %s (isim, soyisim, telefon) VALUES ($1, $2, $3)", tabloAdi)
	_, err := db.Exec(SQLbilgiEkle, isim, soyisim, telefon)
	if err != nil {
		return err
	}
	fmt.Println("Veri başarıyla eklendi.")
	return nil
}

func rehberIcerikGöster(db *sql.DB, tabloAdi string) error {
	SQLgöster := fmt.Sprintf("SELECT * FROM %s", tabloAdi)
	rows, err := db.Query(SQLgöster)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var siraNo int
		var isim, soyisim string
		var telefon int
		if err := rows.Scan(&siraNo, &isim, &soyisim, &telefon); err != nil {
			return err
		}
		fmt.Printf("Sıra No: %d, İsim: %s, Soyisim: %s, Telefon: %d\n", siraNo, isim, soyisim, telefon)
	}
	return nil
}

func main() {
	db, err := DatabaseBaglanma()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var secim int
	for {
		fmt.Println("Menü:")
		fmt.Println("1. Rehber Oluştur")
		fmt.Println("2. Mevcut Tablo Adlarını Göster")
		fmt.Println("3. Tablo Sil")
		fmt.Println("4. İçerik")
		fmt.Println("5. Bilgi ekle")
		fmt.Println("6. Exit")
		fmt.Scan(&secim)

		switch secim {
		case 1:
			var tabloAdi string
			fmt.Print("Tablo adını girin: ")
			fmt.Scan(&tabloAdi)
			rehberOluştur(db, tabloAdi)
			fmt.Printf("'%s' adlı tablo oluşturuldu.\n", tabloAdi)
		case 2:
			tabloAdlari, err := tabloAdlariCek(db)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Mevcut Tablo Adları:")
			for _, ad := range tabloAdlari {
				fmt.Println(ad)
			}
		case 3:
			var tabloAdi string
			fmt.Print("Silmek istediğiniz tablo adını girin: ")
			fmt.Scan(&tabloAdi)
			err := rehberKaldir(db, tabloAdi)
			if err != nil {
				log.Fatal(err)
			}
		case 4:
			var tabloAdi string
			fmt.Print("İçeriğini görmek istediğiniz tablo adını girin: ")
			fmt.Scan(&tabloAdi)
			err := rehberIcerikGöster(db, tabloAdi)
			if err != nil {
				log.Fatal(err)
			}
		case 5:
			var tabloAdi string
			var isim, soyisim string
			var telefon int
			fmt.Print("Tablo adı")
			fmt.Scan(&tabloAdi)
			fmt.Print("İsim: ")
			fmt.Scan(&isim)
			fmt.Print("Soyisim: ")
			fmt.Scan(&soyisim)
			fmt.Print("Telefon: ")
			fmt.Scan(&telefon)
			err := rehberBilgiEkle(db, tabloAdi, isim, soyisim, telefon)
			if err != nil {
				log.Fatal(err)
			}
		case 6:
			fmt.Println("Programdan çıkılıyor.")
			return
		default:
			fmt.Println("Geçersiz seçim. Lütfen tekrar deneyin.")
		}
	}
}
