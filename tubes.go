package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const NMAX int = 30

type Makanan struct {
	Kode    int
	Nama    string
	Harga   int
	Jumlah  int
	Terjual int
}

type makananList struct {
	Data  [NMAX]Makanan
	Count int
}

func main() {
	var m makananList
	inisialisasiMakanan(&m)
	tampilkanMenuUtama(&m)
}

func inisialisasiMakanan(m *makananList) {
	file, err := os.Open("data/makanan.txt")
	if err != nil {
		fmt.Println("Gagal membuka file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if m.Count >= NMAX {
			fmt.Println("Data makanan sudah penuh, tidak dapat menambah lebih banyak data.")
			break
		}

		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) < 5 {
			continue
		}

		kode, _ := strconv.Atoi(fields[0])
		harga, _ := strconv.Atoi(fields[2])
		jumlah, _ := strconv.Atoi(fields[3])
		terjual, _ := strconv.Atoi(fields[4])

		m.Data[m.Count] = Makanan{
			Kode:    kode,
			Nama:    fields[1],
			Harga:   harga,
			Jumlah:  jumlah,
			Terjual: terjual,
		}
		m.Count++
	}
}

func tampilkanMenuUtama(m *makananList) {
	fmt.Println("\n--------------------")
	fmt.Println("Menu Utama")
	fmt.Println("--------------------")
	fmt.Println("1. Menu Pembeli")
	fmt.Println("2. Menu Pengelola")
	fmt.Println("3. Keluar")
	fmt.Print("Pilih opsi: ")

	var pilihan int
	fmt.Scan(&pilihan)

	if pilihan == 1 {
		otorisasiPembelian(m)
	} else if pilihan == 2 {
		menuPengelola(m)
	} else if pilihan == 3 {
		fmt.Println("Terima Kasih Telah Menggunakan Vending Machine Kami ^^")
		os.Exit(0)
	} else {
		fmt.Println("Pilihan tidak valid.")
	}
}

func menuPembeli(m *makananList, pin int) {
	var inputPin int

	fmt.Println("\n--------------------")
	fmt.Println("Menu Pembeli")
	fmt.Println("--------------------")
	fmt.Println("1. Beli Makanan")
	fmt.Println("2. Kembali ke Menu Utama")
	fmt.Print("Pilih opsi: ")

	var pilihan int
	fmt.Scan(&pilihan)

	if pilihan == 1 {
		beliMakanan(m, pin)
	} else if pilihan == 2 {
		fmt.Print("Masukkan pin otorisasi untuk keluar menu: ")
		fmt.Scan(&inputPin)
		if cekPin(pin, inputPin) == true {
			tampilkanMenuUtama(m)
		} else {
			menuPembeli(m, pin)
		}
	} else {
		fmt.Println("Pilihan tidak valid.")
		menuPembeli(m, pin)
	}
}

func otorisasiPembelian(m *makananList) {
	// procedure yang bertujuan agar pembeli tidak dapat keluar dari menu pembelian makanan
	var pin int

	fmt.Println("\n--------------------")
	fmt.Println("Menu Otorisasi Pembelian")
	fmt.Println("--------------------")
	fmt.Print("Buat pin berupa angka 0-9 untuk otorisasi pembelian: ")
	fmt.Scan(&pin)
	menuPembeli(m, pin)
}

func cekPin(pin, inputPin int) bool {
	if inputPin == pin {
		fmt.Println("PIN BENAR")
		return true
	} else {
		fmt.Println("PIN SALAH")
		return false
	}
}

func beliMakanan(m *makananList, pin int) {
	/*
		Dilakukan sorting agar list makanan tidak teracak setelah diurutkan di menu pengelola
	*/
	var pass, i int  // variabel untuk insertion sort array berdasarkan kode makanan
	var temp Makanan // variabel temp untuk insertion sort array berdasarkan kode makanan
	pass = 1
	for pass <= m.Count-1 {
		i = pass
		temp = m.Data[pass]
		for i > 0 && temp.Kode < m.Data[i-1].Kode {
			m.Data[i] = m.Data[i-1]
			i--
		}
		m.Data[i] = temp
		pass++
	}

	var kode int
	var uang int
	var pilihan int
	var seribuan, duaribuan, limaribuan int
	fmt.Println("\n--------------------")
	fmt.Println("Menu Pembeli")
	fmt.Println("--------------------")
	fmt.Println("Daftar Makanan:")
	for i = 0; i < m.Count; i++ {
		fmt.Printf("%d. %s - Rp%d (Stok: %d)\n", m.Data[i].Kode, m.Data[i].Nama, m.Data[i].Harga, m.Data[i].Jumlah)
	}

	fmt.Print("Masukkan kode makanan: ")
	fmt.Scan(&kode)

	i = 0
	for i < m.Count && (m.Data[i].Kode != kode) {
		i += 1
	}
	if i < m.Count && m.Data[i].Kode == kode {
		if m.Data[i].Jumlah > 0 {
			fmt.Printf("Harga %s adalah Rp%d\n", m.Data[i].Nama, m.Data[i].Harga)
			totalHarga := m.Data[i].Harga
			for totalHarga > 0 {
				fmt.Print("Mohon Masukkan Uang dengan Pecahan 1000/2000/5000/10000: ")
				fmt.Scan(&uang)
				if CekPecahan(uang) {
					if totalHarga-uang >= 0 {
						totalHarga -= uang
					} else {
						totalHarga -= uang
						kembalian := totalHarga
						CekKembalian(kembalian, &seribuan, &duaribuan, &limaribuan)
						fmt.Printf("Uang kembalian: %d pecahan seribu, %d pecahan dua ribu, %d pecahan lima ribu\n", seribuan, duaribuan, limaribuan)
					}
				} else {
					fmt.Println("Maaf uang tidak valid")
				}
				if totalHarga > 0 {
					fmt.Printf("Sisa pembayaran: Rp%d\n", totalHarga)
				}
			}
			m.Data[i].Jumlah--
			m.Data[i].Terjual++
			fmt.Println("Terima kasih telah membeli!")

			fmt.Println()
			fmt.Println("Apakah anda masih ingin membeli lagi?")
			fmt.Println("1. Ya")
			fmt.Println("2. Tidak")
			fmt.Print("Pilih opsi: ")
			fmt.Scan(&pilihan)
			if pilihan == 1 {
				beliMakanan(m, pin)
			} else {
				menuPembeli(m, pin)
			}
		} else {
			fmt.Println("Maaf, makanan sudah habis.")

			fmt.Println()
			fmt.Println("Apakah anda masih ingin membeli?")
			fmt.Println("1. Ya")
			fmt.Println("2. Tidak")
			fmt.Print("Pilih opsi: ")
			fmt.Scan(&pilihan)
			if pilihan == 1 {
				beliMakanan(m, pin)
			} else {
				menuPembeli(m, pin)
			}
		}
	} else {
		fmt.Print("Kode makanan tidak valid.")
		beliMakanan(m, pin)
	}
}

func CekPecahan(uang int) bool {
	if uang == 1000 || uang == 2000 || uang == 5000 || uang == 10000 {
		return true
	} else {
		return false
	}
}

func CekKembalian(kembalian int, seribuan, duaribuan, limaribuan *int) {
	var sisa int
	kembalian = kembalian * -1
	*limaribuan = kembalian / 5000
	sisa = kembalian % 5000
	*duaribuan = sisa / 2000
	sisa = sisa % 2000
	*seribuan = sisa / 1000
}

func menuPengelola(m *makananList) {
	fmt.Println("\n--------------------")
	fmt.Println("Menu Pengelola")
	fmt.Println("--------------------")
	fmt.Println("1. Cari Data Makanan Tertentu")
	fmt.Println("2. Urutkan Data Makanan Berdasarkan Ketersediaan")
	fmt.Println("3. Cari Makanan Paling Laku")
	fmt.Println("4. Kembali ke Menu Utama")
	fmt.Print("Pilih opsi: ")

	var pilihan int
	fmt.Scan(&pilihan)

	if pilihan == 1 {
		cariMakananTertentu(m)
	} else if pilihan == 2 {
		urutkanMakananBerdasarkanKetersediaanTerbanyak(m)
	} else if pilihan == 3 {
		cariMakananPalingLaku(m)
	} else if pilihan == 4 {
		tampilkanMenuUtama(m)
	} else {
		fmt.Println("Pilihan tidak valid.")
		menuPengelola(m)
	}
}

func cariMakananTertentu(m *makananList) {
	var kode int
	fmt.Println("\n--------------------")
	fmt.Println("Menu Pengelola")
	fmt.Println("--------------------")
	fmt.Print("Masukkan kode makanan yang ingin dicari: ")
	fmt.Scan(&kode)

	i := 0
	isKetemu := false
	for i < m.Count && isKetemu != true {
		if m.Data[i].Kode == kode {
			fmt.Printf("Nama: %s, Harga: Rp%d, Terjual: %d, Stok: %d\n", m.Data[i].Nama, m.Data[i].Harga, m.Data[i].Terjual, m.Data[i].Jumlah)
			isKetemu = true
		}
		i++
	}
	menuPengelola(m)
}

func urutkanMakananBerdasarkanKetersediaanTerbanyak(m *makananList) {
	/*
		Sorting menggunakan selection sort
	*/
	var pass, idx, i int
	var temp Makanan
	pass = 1
	for pass <= m.Count-1 {
		idx = pass - 1
		i = pass
		for i < m.Count {
			if (m.Data[idx].Jumlah < m.Data[i].Jumlah) || (m.Data[idx].Jumlah == m.Data[i].Jumlah && m.Data[idx].Nama > m.Data[i].Nama) {
				idx = i
			}
			i++
		}
		temp = m.Data[pass-1]
		m.Data[pass-1] = m.Data[idx]
		m.Data[idx] = temp
		pass++
	}
	fmt.Println("\n--------------------")
	fmt.Println("Menu Pengelola")
	fmt.Println("--------------------")
	fmt.Println("Daftar Makanan Berdasarkan Ketersediaan Terbanyak:")
	for i := 0; i < m.Count; i++ {
		fmt.Println(m.Data[i])
	}
	menuPengelola(m)
}

func cariMakananPalingLaku(m *makananList) {
	var max, j int

	max = -1

	for i := 0; i < m.Count; i++ {
		if m.Data[i].Terjual > max {
			max = m.Data[i].Terjual
			j = i
		}
	}

	fmt.Println("\n--------------------")
	fmt.Println("Menu Pengelola")
	fmt.Println("--------------------")
	fmt.Printf("Makanan paling laku adalah %s dengan %d terjual.\n", m.Data[j].Nama, m.Data[j].Terjual)
	menuPengelola(m)
}
