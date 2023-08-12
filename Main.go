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
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tabloAdi)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Printf("'%s' adlı tablo silindi.\n", tabloAdi)
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
		fmt.Println("4. Çıkış")
		fmt.Print("Seçiminizi yapın (1/2/3/4): ")
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
			fmt.Println("Programdan çıkılıyor...")
			return
		default:
			fmt.Println("Geçersiz seçim. Lütfen tekrar deneyin.")
		}
	}
}
