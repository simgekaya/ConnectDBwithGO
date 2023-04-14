package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
)

type Kisi struct {
	Id    int
	Ad    string
	Soyad string
}

func main() {
	db, err := sql.Open("sqlite3", "./kisiler.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Yeni bir kişi ekleme
	err = ekle(db, Kisi{Ad: "Ali", Soyad: "Demir"})
	err = ekle(db, Kisi{Ad: "Ayça", Soyad: "Onursal"})
	if err != nil {
		panic(err)
	}

	// Kişi bilgilerini güncelleme
	err = guncelle(db, Kisi{Id: 1, Ad: "Ayşe", Soyad: "Yılmaz"})
	if err != nil {
		panic(err)
	}

	// Kişi silme
	err = sil(db, 2)
	if err != nil {
		panic(err)
	}

	// Tüm kişileri getirme
	kisiler := hepsiniGetir(db)

	jsonData, err := json.Marshal(kisiler)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("kisiler.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Veriler başarıyla JSON dosyasına aktarıldı.")
}

// Veri ekleme
func ekle(db *sql.DB, kisi Kisi) error {
	stmt, err := db.Prepare("INSERT INTO kisiler(ad, soyad) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(kisi.Ad, kisi.Soyad)
	if err != nil {
		return err
	}

	return nil
}

// Veri güncelleme
func guncelle(db *sql.DB, kisi Kisi) error {
	stmt, err := db.Prepare("UPDATE kisiler SET ad=?, soyad=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(kisi.Ad, kisi.Soyad, kisi.Id)
	if err != nil {
		return err
	}

	return nil
}

// Veri silme
func sil(db *sql.DB, id int) error {
	stmt, err := db.Prepare("DELETE FROM kisiler WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

// Tüm verileri getirme
func hepsiniGetir(db *sql.DB) []Kisi {
	rows, err := db.Query("SELECT * FROM kisiler")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	kisiler := []Kisi{}
	for rows.Next() {
		kisi := Kisi{}
		err := rows.Scan(&kisi.Id, &kisi.Ad, &kisi.Soyad)
		if err != nil {
			panic(err)
		}
		kisiler = append(kisiler, kisi)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	return kisiler

}
