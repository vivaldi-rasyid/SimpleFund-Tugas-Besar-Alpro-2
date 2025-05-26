package main

import (
	"fmt"
	"os"
	"os/exec"
)

// variable global
const NMAX int = 100
const ADMIN int = 0
const USER int = 1

var daftarProyek [NMAX]Proyek
var jumlahProyek int = 0
var daftarDonatur [NMAX]Donatur
var jumlahDonatur int = 0

var daftarUser [NMAX]User
var jumlahUser int = 0
var currentUser User
var nextUserID int = 1

var systemMessage string

type Proyek struct {
	id             int
	nama           string
	danaDibutuhkan int
	danaTerkumpul  int
}

type Donatur struct {
	id          int
	nama        string
	totalDonasi int
}

type User struct {
	id       int
	username string
	password string
	userType int
}

func main() {
    initAdmin()
    var pilihan int
    var loggedIn bool = false
    var aplikasiAktif bool = true

    for aplikasiAktif {
        if !loggedIn {
            clear()
            fmt.Println("\n====================== SimpleFund Login Menu ======================")
            fmt.Println("1. Login")
            fmt.Println("2. Register")
            fmt.Println("3. Keluar Aplikasi")
            fmt.Print("Pilih menu: ")
            fmt.Scan(&pilihan)

            if pilihan == 1 {
                loggedIn = login()
            } else if pilihan == 2 {
                register()
            } else if pilihan == 3 {
                fmt.Println("Terima kasih telah menggunakan aplikasi 🤗")
                aplikasiAktif = false
            } else {
                systemMessage = "⛔ Pilihan tidak valid\n"
            }
        } else {
            if currentUser.userType == ADMIN {
                adminMenu()
            } else {
                userMenu()
            }
            //logout
            loggedIn = false
        }
    }
}

// inisialisasi admin
func initAdmin() {
    var adminUser User
    adminUser.id = 0
    adminUser.username = "admin"
    adminUser.password = "admin"
    adminUser.userType = ADMIN

    daftarUser[jumlahUser] = adminUser
    jumlahUser++
}

func login() bool {
    var username string
    var password string
    fmt.Println("\n======================== Login SimpleFund =========================")
    fmt.Print("Username: ")
    fmt.Scan(&username)
    fmt.Print("Password: ")
    fmt.Scan(&password)

    var i int
    for i = 0; i < jumlahUser; i++ {
        if daftarUser[i].username == username && daftarUser[i].password == password {
            currentUser = daftarUser[i]
            systemMessage = fmt.Sprintf("✅ Login berhasil! \nSelamat datang %s👋 \nID Anda adalah: %d", currentUser.username, currentUser.id)
            return true
        }
    }
    systemMessage = "⛔ Username atau password salah!"
    return false
}

// Register user baru
func register() {
    if jumlahUser >= NMAX {
        systemMessage = "💬 Kapasitas pengguna penuh!"
        return
    }

    var newUser User
    fmt.Println("\n======================== Register SimpleFund ========================")
    var isUnique bool = false
    for !isUnique {
        fmt.Print("Username baru: ")
        fmt.Scan(&newUser.username)

        isUnique = true
        // Cek username di daftar user
        var i int
        for i = 0; i < jumlahUser && isUnique; i++ {
            if daftarUser[i].username == newUser.username {
                fmt.Println("💬 Username sudah digunakan, silakan pilih username lain")
                isUnique = false
            }
        }
        // Cek username di daftar donatur (tidak boleh sama dengan nama donatur)
        for i = 0; i < jumlahDonatur && isUnique; i++ {
            if daftarDonatur[i].nama == newUser.username {
                fmt.Println("💬 Username sudah digunakan, silakan pilih username lain")
                isUnique = false
            }
        }
    }

    fmt.Print("Password: ")
    fmt.Scan(&newUser.password)

    newUser.userType = USER

    // assign ID otomatis (dimulai dari 1, bertambah secara berurutan, kalau udah dipake, skip)
    var idUnik bool = false
    for !idUnik {
        idUnik = true
        var i int
        for i = 0; i < jumlahUser; i++ {
            if daftarUser[i].id == nextUserID {
                idUnik = false
            }
        }
        for i = 0; i < jumlahDonatur; i++ {
            if daftarDonatur[i].id == nextUserID {
                idUnik = false
            }
        }
        if !idUnik {
            nextUserID++
        }
    }
    newUser.id = nextUserID

    daftarUser[jumlahUser] = newUser
    jumlahUser++

    systemMessage = fmt.Sprintf("Registrasi berhasil! 🙌\nAnda telah terdaftar dengan ID: %d\n\nSilakan login dengan username dan password Anda 🤗", newUser.id)
}

// menu user biasa
func userMenu() {
    var pilihan int
    var pilihanSorting int
    var selesai bool = false
    var urutkanDefault bool = true
    var urutkanNamaProyek bool = false
    var urutkanDanaTerkumpul bool = false
    var urutkanDanaDibutuhkan bool = false
    var ProyekDicari bool = false
    var lastSearch bool = false // untuk melacak searching sebelumnya

    for !selesai {
        if !lastSearch {
            clear() // hanya clear line jika bukan setelah pencarian
        }
        lastSearch = false // reset lastSearch

        if systemMessage != "" {
            fmt.Println(systemMessage)
            systemMessage = ""
        }
        
        if ProyekDicari == false {
            fmt.Println("\n======================= Menu User SimpleFund =======================")
        }
        // pastikan selalu ada output tabel kecuali sedang searching
        if !ProyekDicari {
            if urutkanDefault {
                tampilkanProyekDefault(daftarProyek, jumlahProyek)
            } else if urutkanNamaProyek {
                tampilkanProyekUrutkanNama(daftarProyek, jumlahProyek)
            } else if urutkanDanaTerkumpul {
                tampilkanProyekUrutkanDanaTerkumpul(daftarProyek, jumlahProyek)
            } else if urutkanDanaDibutuhkan {
                tampilkanProyekUrutkanDanaDibutuhkan(daftarProyek, jumlahProyek)
            } else {
                // jika tidak ada pengurutan yang aktif, kembalikan ke default
                urutkanDefault = true
                tampilkanProyekDefault(daftarProyek, jumlahProyek)
            }
        }

        fmt.Println("\n=========================== Pilih Menu 🔎 ==========================")
        fmt.Println("1. Tambah Proyek")
        fmt.Println("2. Donasi ke Proyek")
        fmt.Println("3. Urutkan Proyek")
        fmt.Println("4. Cari Proyek")
        fmt.Println("5. Logout")
        fmt.Print("Pilih menu: ")
        fmt.Scan(&pilihan)

        if pilihan == 1 {
            tambahProyek(&daftarProyek, &jumlahProyek)
        } else if pilihan == 2 {
            donasiUser(&daftarProyek, jumlahProyek, &daftarDonatur, &jumlahDonatur)
        } else if pilihan == 3 {
            fmt.Println("============== Pilih Jenis Sorting 🔀 ==============")
            fmt.Println("1. Urutkan berdasarkan nama (Ascending)")
            fmt.Println("2. Urutkan berdasarkan dana terkumpul (Descending)")
            fmt.Println("3. Urutkan berdasarkan dana dibutuhkan (Descending)")
            fmt.Println("4. Kembali ke Urutan Default")
            fmt.Print("Pilih menu: ")
            fmt.Scan(&pilihanSorting)
            if pilihanSorting == 1 {
                urutkanDefault = false
                urutkanNamaProyek = true
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                ProyekDicari = false
            } else if pilihanSorting == 2 {
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = true
                urutkanDanaDibutuhkan = false
                ProyekDicari = false
            } else if pilihanSorting == 3 {
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = true
                ProyekDicari = false
            } else if pilihanSorting == 4 {
                urutkanDefault = true
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                ProyekDicari = false
            } else {
                fmt.Println("⛔ Pilihan tidak valid\n")
            }
        } else if pilihan == 4 {
            fmt.Println("\n============ Cari Proyek 🔎 ==============")
            fmt.Println("1. Cari berdasarkan ID")
            fmt.Println("2. Cari berdasarkan Nama")
            fmt.Print("Pilihan: ")
            var metodeCari int
            fmt.Scan(&metodeCari)

            if metodeCari == 1 {
                var dicariID int
                fmt.Print("Masukkan ID Proyek yang ingin dicari: ")
                fmt.Scan(&dicariID)
                clear() // clear line sebelum menampilkan hasil pencarian

                var indexDitemukan int
                indexDitemukan = linearSearchIDProyek(daftarProyek, jumlahProyek, dicariID)
                fmt.Println("\n======================= Menu User SimpleFund =======================")
                tampilkanProyekDicari(daftarProyek, indexDitemukan)
                lastSearch = true
                ProyekDicari = true
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                if indexDitemukan == -1 {
                    ProyekDicari = false
                    urutkanDefault = true
                }
            } else if metodeCari == 2 {
                var dicariNama string
                fmt.Print("Masukkan Nama Proyek yang ingin dicari: ")
                fmt.Scan(&dicariNama)
                clear() // clear line sebelum menampilkan hasil pencarian

                var indexDitemukan int
                indexDitemukan = linearSearchNamaProyek(daftarProyek, jumlahProyek, dicariNama)
                fmt.Println("\n======================= Menu User SimpleFund =======================")
                tampilkanProyekDicari(daftarProyek, indexDitemukan)
                lastSearch = true
                ProyekDicari = true
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                if indexDitemukan == -1 {
                    ProyekDicari = false
                    urutkanDefault = true
                }
            } else {
                systemMessage = "⛔ Pilihan tidak valid\n"
                ProyekDicari = false
                urutkanDefault = true
            }
        } else if pilihan == 5 {
            selesai = true
            systemMessage = "Logout berhasil 🙌"
        } else {
            systemMessage = "⛔ Pilihan tidak valid\n"
            ProyekDicari = false
        }
    }
}

// menu admin
func adminMenu() {
    var pilihan int
    var pilihanSorting int
    var urutkanDefault bool = true
    var urutkanNamaProyek bool = false
    var urutkanDanaTerkumpul bool = false
    var urutkanDanaDibutuhkan bool = false
    var urutkanDefaultDonatur bool = false
    var urutkanNamaDonatur bool = false
    var urutkanTotalDonasiDonatur bool = false
    var ProyekDicari bool = false
    var DonaturDicari bool = false
    var lastSearch bool = false // untuk melacak search sebelumnya

    var selesai bool = false
    var i int = 1

    for !selesai {
        if !lastSearch {
            clear() // hanya clear line jika bukan setelah pencarian
        }
        lastSearch = false // reset lastSearch

        if systemMessage != "" {
            fmt.Println(systemMessage)
            systemMessage = ""
        }
        if i == 1 {
            fmt.Println("\n================== Selamat Datang Admin SimpleFund =================")
        }
        if ProyekDicari == false && DonaturDicari == false {
            fmt.Println("============================ Menu Admin ============================")
        }
        // pastikan selalu ada output tabel kecuali sedang searching
        if !ProyekDicari && !DonaturDicari {
            if urutkanDefault {
                tampilkanProyekDefault(daftarProyek, jumlahProyek)
            } else if urutkanNamaProyek {
                tampilkanProyekUrutkanNama(daftarProyek, jumlahProyek)
            } else if urutkanDanaTerkumpul {
                tampilkanProyekUrutkanDanaTerkumpul(daftarProyek, jumlahProyek)
            } else if urutkanDanaDibutuhkan {
                tampilkanProyekUrutkanDanaDibutuhkan(daftarProyek, jumlahProyek)
            } else if urutkanDefaultDonatur {
                tampilkanDonaturDefault(daftarDonatur, jumlahDonatur)
            } else if urutkanNamaDonatur {
                tampilkanDonaturUrutkanNama(daftarDonatur, jumlahDonatur)
            } else if urutkanTotalDonasiDonatur {
                tampilkanDonaturUrutkanTotalDonasi(daftarDonatur, jumlahDonatur)
            } else {
                // jika tidak ada pengurutan yang aktif kembalikan ke default proyek
                urutkanDefault = true
                tampilkanProyekDefault(daftarProyek, jumlahProyek)
            }
        }

        fmt.Println("\n=========================== Pilih Menu 🔎 ==========================")
        fmt.Println("1. Tambah Proyek")
        fmt.Println("2. Tambah Donatur")
        fmt.Println("3. Donasi ke Proyek")
        if urutkanDefault || urutkanNamaProyek || urutkanDanaTerkumpul || urutkanDanaDibutuhkan {
            fmt.Println("4. Tampilkan Donatur")
        } else {
            fmt.Println("4. Tampilkan Proyek")
        }
        if urutkanDefault || urutkanNamaProyek || urutkanDanaTerkumpul || urutkanDanaDibutuhkan {
            fmt.Println("5. Urutkan Proyek")
        } else {
            fmt.Println("5. Urutkan Donatur")
        }
        fmt.Println("6. Cari Proyek")
        fmt.Println("7. Cari Donatur")
        fmt.Println("8. Edit Proyek")
        fmt.Println("9. Edit Donatur")
        fmt.Println("10. Hapus Proyek")
        fmt.Println("11. Hapus Donatur")
        fmt.Println("12. Logout")
        fmt.Print("Pilih menu: ")
        fmt.Scan(&pilihan)
        i++

        if pilihan == 1 {
            tambahProyek(&daftarProyek, &jumlahProyek)
        } else if pilihan == 2 {
            tambahDonatur(&daftarDonatur, &jumlahDonatur)
        } else if pilihan == 3 {
            donasi(&daftarProyek, jumlahProyek, &daftarDonatur, jumlahDonatur)
        } else if pilihan == 4 && (urutkanDefault || urutkanNamaProyek || urutkanDanaTerkumpul || urutkanDanaDibutuhkan) {
            urutkanDefault = false
            urutkanNamaProyek = false
            urutkanDanaTerkumpul = false
            urutkanDanaDibutuhkan = false
            urutkanDefaultDonatur = true
            urutkanNamaDonatur = false
            urutkanTotalDonasiDonatur = false
            ProyekDicari = false
            DonaturDicari = false
        } else if pilihan == 4 && !(urutkanDefault || urutkanNamaProyek || urutkanDanaTerkumpul || urutkanDanaDibutuhkan) {
            urutkanDefault = true
            urutkanNamaProyek = false
            urutkanDanaTerkumpul = false
            urutkanDanaDibutuhkan = false
            urutkanDefaultDonatur = false
            urutkanNamaDonatur = false
            urutkanTotalDonasiDonatur = false
            ProyekDicari = false
            DonaturDicari = false
        } else if pilihan == 5 && (urutkanDefault || urutkanNamaProyek || urutkanDanaTerkumpul || urutkanDanaDibutuhkan) {
            fmt.Println("\n============== Pilih Jenis Sorting 🔀 ==============")
            fmt.Println("1. Urutkan berdasarkan nama (Ascending)")
            fmt.Println("2. Urutkan berdasarkan dana terkumpul (Descending)")
            fmt.Println("3. Urutkan berdasarkan dana dibutuhkan (Descending)")
            fmt.Println("4. Kembali ke Urutan Default")
            fmt.Print("Pilih menu: ")
            fmt.Scan(&pilihanSorting)
            if pilihanSorting == 1 {
                urutkanDefault = false
                urutkanNamaProyek = true
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = false
                ProyekDicari = false
                DonaturDicari = false
            } else if pilihanSorting == 2 {
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = true
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = false
                ProyekDicari = false
                DonaturDicari = false
            } else if pilihanSorting == 3 {
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = true
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = false
                ProyekDicari = false
                DonaturDicari = false
            } else if pilihanSorting == 4 {
                urutkanDefault = true
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = false
                ProyekDicari = false
                DonaturDicari = false
            } else {
                systemMessage = "⛔ Pilihan tidak valid\n"
            }
        } else if pilihan == 5 && !(urutkanDefault || urutkanNamaProyek || urutkanDanaTerkumpul || urutkanDanaDibutuhkan) {
            fmt.Println("\n============== Pilih Jenis Sorting 🔀 ==============")
            fmt.Println("1. Urutkan berdasarkan nama (Ascending)")
            fmt.Println("2. Urutkan berdasarkan total donasi (Descending)")
            fmt.Println("3. Kembali ke Urutan Default")
            fmt.Print("Pilih menu: ")
            fmt.Scan(&pilihanSorting)
            if pilihanSorting == 1 {
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = true
                urutkanTotalDonasiDonatur = false
                ProyekDicari = false
                DonaturDicari = false
            } else if pilihanSorting == 2 {
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = true
                ProyekDicari = false
                DonaturDicari = false
            } else if pilihanSorting == 3 {
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = true
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = false
                ProyekDicari = false
                DonaturDicari = false
            } else {
                systemMessage = "⛔ Pilihan tidak valid\n"
            }
        } else if pilihan == 6 {
            fmt.Println("\n============ Cari Proyek 🔎 ==============")
            fmt.Println("1. Cari berdasarkan ID")
            fmt.Println("2. Cari berdasarkan Nama")
            fmt.Print("Pilihan: ")
            var metodeCari int
            fmt.Scan(&metodeCari)

            if metodeCari == 1 {
                var dicariID int
                fmt.Print("Masukkan ID Proyek yang ingin dicari: ")
                fmt.Scan(&dicariID)
                clear() // clear line sebelum menampilkan hasil pencarian

                var indexDitemukan int
                indexDitemukan = linearSearchIDProyek(daftarProyek, jumlahProyek, dicariID)
                if indexDitemukan != -1 {
                    fmt.Println("============================ Menu Admin ============================")
                }
                tampilkanProyekDicari(daftarProyek, indexDitemukan)
                lastSearch = true
                ProyekDicari = true
                DonaturDicari = false
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = false
                if indexDitemukan == -1 {
                    ProyekDicari = false
                    urutkanDefault = true
                }
            } else if metodeCari == 2 {
                var dicariNama string
                fmt.Print("Masukkan Nama Proyek yang ingin dicari: ")
                fmt.Scan(&dicariNama)
                clear() // clear line sebelum menampilkan hasil pencarian

                var indexDitemukan int
                indexDitemukan = linearSearchNamaProyek(daftarProyek, jumlahProyek, dicariNama)
                if indexDitemukan != -1 {
                    fmt.Println("============================ Menu Admin ============================")
                }
                tampilkanProyekDicari(daftarProyek, indexDitemukan)
                lastSearch = true
                ProyekDicari = true
                DonaturDicari = false
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = false
                if indexDitemukan == -1 {
                    ProyekDicari = false
                    urutkanDefault = true
                }
            } else {
                systemMessage = "⛔ Pilihan tidak valid\n"
                ProyekDicari = false
                DonaturDicari = false
                urutkanDefault = true
            }
        } else if pilihan == 7 {
            fmt.Println("\n============ Cari Donatur 🔎 ============")
            fmt.Println("1. Cari berdasarkan ID")
            fmt.Println("2. Cari berdasarkan Nama")
            fmt.Print("Pilihan: ")
            var metodeCari int
            fmt.Scan(&metodeCari)

            if metodeCari == 1 {
                var dicariID int
                fmt.Print("Masukkan ID Donatur yang ingin dicari: ")
                fmt.Scan(&dicariID)
                clear() // clear line sebelum menampilkan hasil pencarian

                selectionSortIDDonaturASC(&daftarDonatur, jumlahDonatur)
                var indexDitemukan int
                indexDitemukan = binarySearchIDDonatur(daftarDonatur, jumlahDonatur, dicariID)
                if indexDitemukan != -1 {
                    fmt.Println("============================ Menu Admin ============================")
                }
                tampilkanDonaturDicari(daftarDonatur, indexDitemukan)
                lastSearch = true
                ProyekDicari = false
                DonaturDicari = true
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = false
                if indexDitemukan == -1 {
                    DonaturDicari = false
                    urutkanDefaultDonatur = true
                }
            } else if metodeCari == 2 {
                var dicariNama string
                fmt.Print("Masukkan Nama Donatur yang ingin dicari: ")
                fmt.Scan(&dicariNama)
                clear() // clear line sebelum menampilkan hasil pencarian

                selectionSortNamaDonaturASC(&daftarDonatur, jumlahDonatur)
                var indexDitemukan int
                indexDitemukan = binarySearchNamaDonatur(daftarDonatur, jumlahDonatur, dicariNama)
                if indexDitemukan != -1 {
                    fmt.Println("============================ Menu Admin ============================")
                }
                tampilkanDonaturDicari(daftarDonatur, indexDitemukan)
                lastSearch = true
                ProyekDicari = false
                DonaturDicari = true
                urutkanDefault = false
                urutkanNamaProyek = false
                urutkanDanaTerkumpul = false
                urutkanDanaDibutuhkan = false
                urutkanDefaultDonatur = false
                urutkanNamaDonatur = false
                urutkanTotalDonasiDonatur = false
                if indexDitemukan == -1 {
                    DonaturDicari = false
                    urutkanDefaultDonatur = true
                }
            } else {
                systemMessage = "⛔ Pilihan tidak valid\n"
                ProyekDicari = false
                DonaturDicari = false
                urutkanDefaultDonatur = true
            }
        } else if pilihan == 8 {
            editProyek(&daftarProyek, jumlahProyek)
        } else if pilihan == 9 {
            editDonatur(&daftarDonatur, jumlahDonatur)
        } else if pilihan == 10 {
            hapusProyek(&daftarProyek, &jumlahProyek)
        } else if pilihan == 11 {
            hapusDonatur(&daftarDonatur, &jumlahDonatur)
        } else if pilihan == 12 {
            selesai = true
            systemMessage = "Logout berhasil 🙌"
        } else {
            systemMessage = "⛔ Pilihan tidak valid\n"
            ProyekDicari = false
            DonaturDicari = false
        }
    }
}

/*list
linearsearch : linearSearchNamaProyek, linearSearchIDProyek
binarysearch : binarySearchNamaDonatur, binarySearchIDDonatur
insertionSort : insertionSortNamaProyekASC, insertionSortDanaTerkumpulProyekDSC, insertionSortDanaDibutuhkanProyekDSC
selectionSort : selectionSortNamaDonaturASC, selectionSortTotalDonasiDonaturDSC,
*/

func linearSearchNamaProyek(data [NMAX]Proyek, n int, dicari string) int {
    var i int
    for i = 0; i < n; i++ {
        if data[i].nama == dicari {
            return i
        }
    }
    return -1
}

func linearSearchIDProyek(data [NMAX]Proyek, n int, dicari int) int {
    var i int
    for i = 0; i < n; i++ {
        if data[i].id == dicari {
            return i
        }
    }
    return -1
}

func binarySearchNamaDonatur(data [NMAX]Donatur, n int, dicari string) int {
    var kiri int = 0
    var kanan int = n - 1

    for kiri <= kanan {
        var tengah int = (kiri + kanan) / 2

        if data[tengah].nama == dicari {
            return tengah
        } else {
            if data[tengah].nama < dicari {
                kiri = tengah + 1
            } else {
                kanan = tengah - 1
            }
        }
    }

    return -1
}

func binarySearchIDDonatur(data [NMAX]Donatur, n int, dicari int) int {
    var kiri int = 0
    var kanan int = n - 1

    for kiri <= kanan {
        var tengah int = (kiri + kanan) / 2

        if data[tengah].id == dicari {
            return tengah
        } else {
            if data[tengah].id < dicari {
                kiri = tengah + 1
            } else {
                kanan = tengah - 1
            }
        }
    }

    return -1
}

// Fungsi insertion sort untuk proyek berdasarkan nama (ascending)
func insertionSortNamaProyekASC(data *[NMAX]Proyek, jumlah int) {
    var i, j int
    for i = 1; i < jumlah; i++ {
        var key Proyek = data[i]
        j = i - 1
        for j >= 0 && data[j].nama > key.nama {
            data[j+1] = data[j]
            j--
        }
        data[j+1] = key
    }
}

// Fungsi insertion sort untuk Dana terkumpul proyek (descending)
func insertionSortDanaTerkumpulProyekDSC(data *[NMAX]Proyek, jumlah int) {
    var i, j int
    for i = 1; i < jumlah; i++ {
        var key Proyek = data[i]
        j = i - 1
        for j >= 0 && data[j].danaTerkumpul < key.danaTerkumpul {
            data[j+1] = data[j]
            j--
        }
        data[j+1] = key
    }
}

// Fungsi insertion sort untuk Dana dibutuhkan proyek (descending)
func insertionSortDanaDibutuhkanProyekDSC(data *[NMAX]Proyek, jumlah int) {
    var i, j int
    for i = 1; i < jumlah; i++ {
        var key Proyek = data[i]
        j = i - 1
        for j >= 0 && data[j].danaDibutuhkan < key.danaDibutuhkan {
            data[j+1] = data[j]
            j--
        }
        data[j+1] = key
    }
}

// Fungsi selection sort untuk donatur berdasarkan totalDonasi (descending)
func selectionSortTotalDonasiDonaturDSC(data *[NMAX]Donatur, jumlah int) {
    var i, j, maxIdx int
    for i = 0; i < jumlah-1; i++ {
        maxIdx = i
        for j = i + 1; j < jumlah; j++ {
            if data[j].totalDonasi > data[maxIdx].totalDonasi {
                maxIdx = j
            }
        }
        data[i], data[maxIdx] = data[maxIdx], data[i]
    }
}

// Fungsi selection sort untuk donatur berdasarkan nama (ascending)
func selectionSortNamaDonaturASC(data *[NMAX]Donatur, jumlah int) {
    var i, j, minIdx int
    for i = 0; i < jumlah-1; i++ {
        minIdx = i
        for j = i + 1; j < jumlah; j++ {
            if data[j].nama < data[minIdx].nama {
                minIdx = j
            }
        }
        data[i], data[minIdx] = data[minIdx], data[i]
    }
}

func selectionSortIDDonaturASC(data *[NMAX]Donatur, jumlah int) {
    var i, j, minIdx int
    for i = 0; i < jumlah-1; i++ {
        minIdx = i
        for j = i + 1; j < jumlah; j++ {
            if data[j].id < data[minIdx].id {
                minIdx = j
            }
        }
        data[i], data[minIdx] = data[minIdx], data[i]
    }
}

func cariProyekByID(id int, data [NMAX]Proyek, jumlah int) int {
    var i int
    for i = 0; i < jumlah; i++ {
        if data[i].id == id {
            return i
        }
    }
    return -1
}

func cariDonaturByID(id int, data [NMAX]Donatur, jumlah int) int {
    var i int
    for i = 0; i < jumlah; i++ {
        if data[i].id == id {
            return i
        }
    }
    return -1
}

// tampilkan proyek dicari
func tampilkanProyekDicari(data [NMAX]Proyek, indexDitemukan int) {
    if indexDitemukan != -1 {
        fmt.Println("\n🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
    }

    if indexDitemukan == -1 {
        fmt.Println("⛔ Proyek Yang Dicari Tidak Ditemukan\n")
    } else {
        fmt.Println("======================================= ✅ Proyek Yang Dicari Ditemukan ======================================== ")
    }

    if indexDitemukan != -1 {
        fmt.Printf("%-4s | %-10s | %-20s | %-15s | %-15s | %-30s\n", "No", "ID Proyek", "Nama Proyek", "Dana Terkumpul", "Dana Dibutuhkan", "Status Proyek")
        fmt.Println("----------------------------------------------------------------------------------------------------------------")
    }

    if indexDitemukan != -1 {
        var status string
        if data[indexDitemukan].danaTerkumpul >= data[indexDitemukan].danaDibutuhkan {
            status = "Dana Sudah Mencukupi"
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s\n", indexDitemukan+1, data[indexDitemukan].id, data[indexDitemukan].nama, data[indexDitemukan].danaTerkumpul, data[indexDitemukan].danaDibutuhkan, status)
        } else {
            status = "Kurang Rp"
            var sisaKebutuhan int
            sisaKebutuhan = data[indexDitemukan].danaDibutuhkan - data[indexDitemukan].danaTerkumpul
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s%d\n", indexDitemukan+1, data[indexDitemukan].id, data[indexDitemukan].nama, data[indexDitemukan].danaTerkumpul, data[indexDitemukan].danaDibutuhkan, status, sisaKebutuhan)
        }
    }
    if indexDitemukan != -1 {
        fmt.Println("🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
    }
}

// tampilkan proyek default
func tampilkanProyekDefault(data [NMAX]Proyek, jumlah int) {
    fmt.Println("\n🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
    fmt.Println("================================================ Daftar Proyek =================================================")
    fmt.Printf("%-4s | %-10s | %-20s | %-15s | %-15s | %-30s\n", "No", "ID Proyek", "Nama Proyek", "Dana Terkumpul", "Dana Dibutuhkan", "Status Proyek")
    fmt.Println("----------------------------------------------------------------------------------------------------------------")

    if jumlah == 0 {
        fmt.Printf("%-4s | %-10s | %-20s | %-15s | %-15s | %-30s\n", "-", "-", "-", "-", "-", "-")
    }

    var i int
    for i = 0; i < jumlah; i++ {
        var status string
        if data[i].danaTerkumpul >= data[i].danaDibutuhkan {
            status = "Dana Sudah Mencukupi"
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s\n", i+1, data[i].id, data[i].nama, data[i].danaTerkumpul, data[i].danaDibutuhkan, status)
        } else {
            status = "Kurang Rp"
            var sisaKebutuhan int
            sisaKebutuhan = data[i].danaDibutuhkan - data[i].danaTerkumpul
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s%d\n", i+1, data[i].id, data[i].nama, data[i].danaTerkumpul, data[i].danaDibutuhkan, status, sisaKebutuhan)
        }
    }
    fmt.Println("🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
}

// tampilkan proyek udah diurut nama
func tampilkanProyekUrutkanNama(data [NMAX]Proyek, jumlah int) {
    insertionSortNamaProyekASC(&data, jumlah)
    fmt.Println("\n🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
    fmt.Println("===================================== Daftar Proyek (🔼 Ascending Nama) ========================================")
    fmt.Printf("%-4s | %-10s | %-20s | %-15s | %-15s | %-30s\n", "No", "ID Proyek", "Nama Proyek", "Dana Terkumpul", "Dana Dibutuhkan", "Status Proyek")
    fmt.Println("----------------------------------------------------------------------------------------------------------------")

    if jumlah == 0 {
        fmt.Printf("%-4s | %-10s | %-20s | %-15s | %-15s | %-30s\n", "-", "-", "-", "-", "-", "-")
    }

    var i int
    for i = 0; i < jumlah; i++ {
        var status string
        if data[i].danaTerkumpul >= data[i].danaDibutuhkan {
            status = "Dana Sudah Mencukupi"
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s\n", i+1, data[i].id, data[i].nama, data[i].danaTerkumpul, data[i].danaDibutuhkan, status)
        } else {
            status = "Kurang Rp"
            var sisaKebutuhan int
            sisaKebutuhan = data[i].danaDibutuhkan - data[i].danaTerkumpul
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s%d\n", i+1, data[i].id, data[i].nama, data[i].danaTerkumpul, data[i].danaDibutuhkan, status, sisaKebutuhan)
        }
    }
    fmt.Println("🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
}
// tampilkan proyek udah urut dana terkumpul
func tampilkanProyekUrutkanDanaTerkumpul(data [NMAX]Proyek, jumlah int) {
    insertionSortDanaTerkumpulProyekDSC(&data, jumlah)
    fmt.Println("\n🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
    fmt.Println("================================= Daftar Proyek (🔽 Descending Dana Terkumpul) =================================")
    fmt.Printf("%-4s | %-10s | %-20s | %-15s | %-15s | %-30s\n", "No", "ID Proyek", "Nama Proyek", "Dana Terkumpul", "Dana Dibutuhkan", "Status Proyek")
    fmt.Println("----------------------------------------------------------------------------------------------------------------")

    if jumlah == 0 {
        fmt.Printf("%-4s | %-10s | %-20s | %-15s | %-15s | %-30s\n", "-", "-", "-", "-", "-", "-")
    }

    var i int
    for i = 0; i < jumlah; i++ {
        var status string
        if data[i].danaTerkumpul >= data[i].danaDibutuhkan {
            status = "Dana Sudah Mencukupi"
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s\n", i+1, data[i].id, data[i].nama, data[i].danaTerkumpul, data[i].danaDibutuhkan, status)
        } else {
            status = "Kurang Rp"
            var sisaKebutuhan int
            sisaKebutuhan = data[i].danaDibutuhkan - data[i].danaTerkumpul
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s%d\n", i+1, data[i].id, data[i].nama, data[i].danaTerkumpul, data[i].danaDibutuhkan, status, sisaKebutuhan)
        }
    }
    fmt.Println("🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
}

// tampilkan proyek udah urut dana dibutuhkan
func tampilkanProyekUrutkanDanaDibutuhkan(data [NMAX]Proyek, jumlah int) {
    insertionSortDanaDibutuhkanProyekDSC(&data, jumlah)
    fmt.Println("\n🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
    fmt.Println("================================= Daftar Proyek (🔽 Descending Dana Dibutuhkan) ================================")
    fmt.Printf("%-4s | %-10s | %-20s | %-15s | %-15s | %-30s\n", "No", "ID Proyek", "Nama Proyek", "Dana Terkumpul", "Dana Dibutuhkan", "Status Proyek")
    fmt.Println("----------------------------------------------------------------------------------------------------------------")

    if jumlah == 0 {
        fmt.Printf("%-4s | %-10s | %-20s | %-15s | %-15s | %-30s\n", "-", "-", "-", "-", "-", "-")
    }

    var i int
    for i = 0; i < jumlah; i++ {
        var status string
        if data[i].danaTerkumpul >= data[i].danaDibutuhkan {
            status = "Dana Sudah Mencukupi"
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s\n", i+1, data[i].id, data[i].nama, data[i].danaTerkumpul, data[i].danaDibutuhkan, status)
        } else {
            status = "Kurang Rp"
            var sisaKebutuhan int
            sisaKebutuhan = data[i].danaDibutuhkan - data[i].danaTerkumpul
            fmt.Printf("%-4d | %-10d | %-20s | Rp%-13d | Rp%-13d | %s%d\n", i+1, data[i].id, data[i].nama, data[i].danaTerkumpul, data[i].danaDibutuhkan, status, sisaKebutuhan)
        }
    }
    fmt.Println("🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦🟦")
}

// tampilkan donatur
func tampilkanDonaturDicari(data [NMAX]Donatur, indexDitemukan int) {
    if indexDitemukan != -1 {
        fmt.Println("\n🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩")
    }
    if indexDitemukan == -1 {
        fmt.Println("⛔ Donatur Yang Dicari Tidak Ditemukan\n")
    } else {
        fmt.Println("======================================= ✅ Donatur Yang Dicari Ditemukan =======================================")
    }
    if indexDitemukan != -1 {
        fmt.Printf("%-4s | %-10s | %-20s | %-25s\n", "No", "ID Donatur", "Nama Donatur", "Total Donasi")
        fmt.Println("----------------------------------------------------------------------------------------------------------------")
    }

    if indexDitemukan != -1 {
        fmt.Printf("%-4d | %-10d | %-20s | %-25d\n", indexDitemukan+1, data[indexDitemukan].id, data[indexDitemukan].nama, data[indexDitemukan].totalDonasi)
    }

    if indexDitemukan != -1 {
        fmt.Println("🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩")
    }
}

func tampilkanDonaturDefault(data [NMAX]Donatur, jumlah int) {
    fmt.Println("\n🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩")
    fmt.Println("================================================ Daftar Donatur ================================================")
    fmt.Printf("%-4s | %-10s | %-20s | %-25s\n", "No", "ID Donatur", "Nama Donatur", "Total Donasi")
    fmt.Println("----------------------------------------------------------------------------------------------------------------")

    if jumlah == 0 {
        fmt.Printf("%-4s | %-10s | %-20s | %-25s\n", "-", "-", "-", "-")
    }

    var i int
    for i = 0; i < jumlah; i++ {
        fmt.Printf("%-4d | %-10d | %-20s | %-25d\n", i+1, data[i].id, data[i].nama, data[i].totalDonasi)
    }
    fmt.Println("🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩")
}

func tampilkanDonaturUrutkanNama(data [NMAX]Donatur, jumlah int) {
    selectionSortNamaDonaturASC(&data, jumlah)
    fmt.Println("\n🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩")
    fmt.Println("====================================== Daftar Donatur (🔼 Ascending Nama) ======================================")
    fmt.Printf("%-4s | %-10s | %-20s | %-25s\n", "No", "ID Donatur", "Nama Donatur", "Total Donasi")
    fmt.Println("----------------------------------------------------------------------------------------------------------------")

    if jumlah == 0 {
        fmt.Printf("%-4s | %-10s | %-20s | %-25s\n", "-", "-", "-", "-")
    }

    var i int
    for i = 0; i < jumlah; i++ {
        fmt.Printf("%-4d | %-10d | %-20s | %-25d\n", i+1, data[i].id, data[i].nama, data[i].totalDonasi)
    }
    fmt.Println("🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩")
}

func tampilkanDonaturUrutkanTotalDonasi(data [NMAX]Donatur, jumlah int) {
    selectionSortTotalDonasiDonaturDSC(&data, jumlah)
    fmt.Println("\n🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩")
    fmt.Println("================================== Daftar Donatur (🔽 Descending Total Donasi) =================================")
    fmt.Printf("%-4s | %-10s | %-20s | %-25s\n", "No", "ID Donatur", "Nama Donatur", "Total Donasi")
    fmt.Println("----------------------------------------------------------------------------------------------------------------")

    if jumlah == 0 {
        fmt.Printf("%-4s | %-10s | %-20s | %-25s\n", "-", "-", "-", "-")
    }

    var i int
    for i = 0; i < jumlah; i++ {
        fmt.Printf("%-4d | %-10d | %-20s | %-25d\n", i+1, data[i].id, data[i].nama, data[i].totalDonasi)
    }
    fmt.Println("🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩")
}

func tambahProyek(data *[NMAX]Proyek, jumlah *int) {
    var n int
    fmt.Print("Berapa proyek yang ingin ditambahkan? ")
    fmt.Scan(&n)

    var i int
    for i = 0; i < n && *jumlah < NMAX; i++ {
        var p Proyek
        fmt.Println("\n============ ➕ Tambahkan Proyek ➕ ==============")

        // Validasi ID proyek unik
        var ulang bool = true
        for ulang {
            fmt.Print("Masukkan ID Proyek: ")
            fmt.Scan(&p.id)
            var duplikat bool = false
            var j int
            for j = 0; j < *jumlah && !duplikat; j++ {
                if data[j].id == p.id {
                    fmt.Println("💬 ID proyek sudah digunakan, silakan masukkan ID lain.")
                    duplikat = true
                }
            }
            ulang = duplikat
        }

        // Validasi nama proyek unik
        var ulangNama bool = true
        for ulangNama {
            fmt.Print("Nama Proyek: ")
            fmt.Scan(&p.nama)
            var duplikatNama bool = false
            var j int
            for j = 0; j < *jumlah && !duplikatNama; j++ {
                if data[j].nama == p.nama {
                    fmt.Println("💬 Nama proyek sudah digunakan, silakan masukkan nama lain.")
                    duplikatNama = true
                }
            }
            ulangNama = duplikatNama
        }

        fmt.Print("Dana Dibutuhkan: ")
        fmt.Scan(&p.danaDibutuhkan)

        p.danaTerkumpul = 0
        data[*jumlah] = p
        *jumlah++
        systemMessage = "\n✅ Proyek Berhasil Ditambahkan\n"
    }

    if *jumlah >= NMAX {
        systemMessage = "\nKapasitas proyek penuh."
    }
}

func tambahDonatur(data *[NMAX]Donatur, jumlah *int) {
    var n int
    fmt.Print("Berapa donatur yang ingin ditambahkan? ")
    fmt.Scan(&n)

    var i int
    for i = 0; i < n && *jumlah < NMAX; i++ {
        var d Donatur
        fmt.Println("\n============ ➕ Tambahkan Donatur ➕ ==============")

        // Validasi ID donatur unik
        var ulang bool = true
        for ulang {
            fmt.Print("ID Donatur: ")
            fmt.Scan(&d.id)
            var duplikat bool = false
            var j int
            // Cek di daftar donatur
            for j = 0; j < *jumlah && !duplikat; j++ {
                if data[j].id == d.id {
                    fmt.Println("💬 ID donatur sudah digunakan di daftar donatur, silakan masukkan ID lain.")
                    duplikat = true
                }
            }
            // Cek di daftar user
            for j = 0; j < jumlahUser && !duplikat; j++ {
                if daftarUser[j].id == d.id {
                    fmt.Println("💬 ID donatur sudah digunakan di daftar user, silakan masukkan ID lain.")
                    duplikat = true
                }
            }
            ulang = duplikat
        }

        // Validasi nama donatur unik
        var ulangNama bool = true
        for ulangNama {
            fmt.Print("Nama Donatur: ")
            fmt.Scan(&d.nama)
            var duplikatNama bool = false
            var j int
            // Tidak boleh sama dengan username user
            for j = 0; j < jumlahUser && !duplikatNama; j++ {
                if daftarUser[j].username == d.nama {
                    fmt.Println("💬 Nama donatur sudah digunakan di daftar user, silakan masukkan nama lain.")
                    duplikatNama = true
                }
            }
            // Tidak boleh sama dengan nama donatur lain
            for j = 0; j < *jumlah && !duplikatNama; j++ {
                if data[j].nama == d.nama {
                    fmt.Println("💬 Nama donatur sudah digunakan di daftar donatur, silakan masukkan nama lain.")
                    duplikatNama = true
                }
            }
            ulangNama = duplikatNama
        }

        d.totalDonasi = 0
        data[*jumlah] = d
        *jumlah++
        systemMessage = "✅ Donatur Berhasil Ditambahkan\n"
    }

    if *jumlah >= NMAX {
        systemMessage = "⛔ Kapasitas donatur penuh\n"
    }
}

func donasi(proyek *[NMAX]Proyek, jumlahProyek int, donatur *[NMAX]Donatur, jumlahDonatur int) {
    var idProyek, idDonatur, nominal int
    fmt.Println("\n--- Donasi ke Proyek ---")
    fmt.Print("Masukkan ID proyek: ")
    fmt.Scan(&idProyek)
    fmt.Print("Masukkan ID donatur: ")
    fmt.Scan(&idDonatur)
    fmt.Print("Masukkan jumlah donasi: ")
    fmt.Scan(&nominal)

    if nominal <= 0 {
        systemMessage = "⛔ Nominal donasi harus lebih dari 0\n"
        return
    }

    var indeksProyek int
    var indeksDonatur int
    indeksProyek = cariProyekByID(idProyek, *proyek, jumlahProyek)
    indeksDonatur = cariDonaturByID(idDonatur, *donatur, jumlahDonatur)

    if indeksProyek == -1 || indeksDonatur == -1 {
        systemMessage = "⛔ ID Proyek atau Donatur tidak ditemukan\n"
        return
    }

    var sisaKebutuhan int
    sisaKebutuhan = proyek[indeksProyek].danaDibutuhkan - proyek[indeksProyek].danaTerkumpul
    if nominal > sisaKebutuhan {
        systemMessage = fmt.Sprintf("⛔ Donasi melebihi kebutuhan proyek, maksimal yang bisa didonasikan: %d\n", sisaKebutuhan)
        return
    }

    proyek[indeksProyek].danaTerkumpul += nominal
    donatur[indeksDonatur].totalDonasi += nominal
    systemMessage = "✅ Donasi berhasil ditambahkan\n"
}

// fungsi donasi khusus untuk user biasa
func donasiUser(proyek *[NMAX]Proyek, jumlahProyek int, donatur *[NMAX]Donatur, jumlahDonatur *int) {
    var idProyek, nominal int
    fmt.Println("\n--- Donasi ke Proyek ---")
    fmt.Print("Masukkan ID proyek: ")
    fmt.Scan(&idProyek)
    fmt.Print("Masukkan jumlah donasi: ")
    fmt.Scan(&nominal)

    if nominal <= 0 {
        systemMessage = "⛔ Nominal donasi harus lebih dari 0\n"
        return
    }

    var indeksProyek int
    indeksProyek = cariProyekByID(idProyek, *proyek, jumlahProyek)
    if indeksProyek == -1 {
        systemMessage = "⛔ ID proyek tidak ditemukan\n"
        return
    }

    // cari donatur berdasarkan username user yang sedang login
    var indeksDonatur int = -1
    var i int = 0
    var ditemukanDonatur bool = false

    for i < *jumlahDonatur && !ditemukanDonatur {
        if donatur[i].nama == currentUser.username {
            indeksDonatur = i
            ditemukanDonatur = true
        }
        i++
    }

    // jika donatur belum ada, buat donatur baru secara otomatis
    if indeksDonatur == -1 {
        if *jumlahDonatur >= NMAX {
            systemMessage = "⛔ Kapasitas donatur penuh, hubungi admin untuk bantuan\n"
            return
        }

        var newID int = currentUser.id
        var d Donatur
        d.id = newID
        d.nama = currentUser.username
        d.totalDonasi = 0

        donatur[*jumlahDonatur] = d
        indeksDonatur = *jumlahDonatur
        *jumlahDonatur++
    }

    var sisaKebutuhan int
    sisaKebutuhan = proyek[indeksProyek].danaDibutuhkan - proyek[indeksProyek].danaTerkumpul
    if nominal > sisaKebutuhan {
        systemMessage = fmt.Sprintf("⛔ Donasi melebihi kebutuhan proyek, maksimal yang bisa didonasikan: %d\n", sisaKebutuhan)
        return
    }

    proyek[indeksProyek].danaTerkumpul += nominal
    donatur[indeksDonatur].totalDonasi += nominal
    systemMessage = fmt.Sprintf("✅ Donasi berhasil ditambahkan.\nAnda (%s) dengan ID %d telah mendonasikan Rp%d ke proyek %s.\n", currentUser.username, currentUser.id, nominal, proyek[indeksProyek].nama)
}

func editProyek(data *[NMAX]Proyek, jumlah int) {
    var id int
    fmt.Println("\n============ Edit Proyek 📑 ==============")
    fmt.Print("Masukkan ID proyek yang ingin diedit: ")
    fmt.Scan(&id)
    var indeks int
    indeks = cariProyekByID(id, *data, jumlah)

    if indeks == -1 {
        systemMessage = "⛔ Proyek dengan ID tersebut tidak ditemukan\n"
        return
    }

    fmt.Println("✅ Proyek ditemukan, masukkan data baru 📝")

    // Validasi nama proyek unik
    var ulangNama bool = true
    var namaBaru string
    for ulangNama {
        fmt.Print("Nama Proyek: ")
        fmt.Scan(&namaBaru)
        var duplikatNama bool = false
        var j int
        for j = 0; j < jumlah && !duplikatNama; j++ {
            if j != indeks && data[j].nama == namaBaru {
                fmt.Println("💬 Nama proyek sudah digunakan, silakan masukkan nama lain.")
                duplikatNama = true
            }
        }
        ulangNama = duplikatNama
    }
    data[indeks].nama = namaBaru

    fmt.Print("Dana Dibutuhkan: ")
    fmt.Scan(&data[indeks].danaDibutuhkan)
    systemMessage = "📝 Data proyek berhasil diperbarui\n"
}

func editDonatur(data *[NMAX]Donatur, jumlah int) {
    var id int
    fmt.Println("\n============ Edit Donatur 📑 ==============")
    fmt.Print("Masukkan ID donatur yang ingin diedit: ")
    fmt.Scan(&id)
    var indeks int
    indeks = cariDonaturByID(id, *data, jumlah)

    if indeks == -1 {
        systemMessage = "⛔ Donatur dengan ID tersebut tidak ditemukan\n"
        return
    }

    fmt.Println("✅ Donatur ditemukan, masukkan data baru 📝")

    // Validasi nama donatur unik dan tidak sama dengan username user
    var ulangNama bool = true
    var namaBaru string
    for ulangNama {
        fmt.Print("Nama Donatur: ")
        fmt.Scan(&namaBaru)
        var duplikatNama bool = false
        var j int
        for j = 0; j < jumlahUser && !duplikatNama; j++ {
            if daftarUser[j].username == namaBaru {
                fmt.Println("💬 Nama donatur sudah digunakan di daftar user, silakan masukkan nama lain.")
                duplikatNama = true
            }
        }
        for j = 0; j < jumlah && !duplikatNama; j++ {
            if j != indeks && data[j].nama == namaBaru {
                fmt.Println("💬 Nama donatur sudah digunakan di daftar donatur, silakan masukkan nama lain.")
                duplikatNama = true
            }
        }
        ulangNama = duplikatNama
    }
    data[indeks].nama = namaBaru

    systemMessage = "📝 Data donatur berhasil diperbarui\n"
}

func hapusProyek(data *[NMAX]Proyek, jumlah *int) {
    var id int
    fmt.Println("\n============ Hapus Proyek 📑 ==============")
    fmt.Print("Masukkan ID proyek yang ingin dihapus: ")
    fmt.Scan(&id)
    var indeks int
    indeks = cariProyekByID(id, *data, *jumlah)

    if indeks == -1 {
        systemMessage = "⛔ Proyek dengan ID tersebut tidak ditemukan\n"
        return
    }

    var i int
    for i = indeks; i < *jumlah-1; i++ {
        data[i] = data[i+1]
    }
    *jumlah--
    systemMessage = "📑 Proyek berhasil dihapus\n"
}

func hapusDonatur(data *[NMAX]Donatur, jumlah *int) {
    var id int
    fmt.Println("\n============ Hapus Donatur 📑 ==============")
    fmt.Print("Masukkan ID donatur yang ingin dihapus: ")
    fmt.Scan(&id)
    var indeks int
    indeks = cariDonaturByID(id, *data, *jumlah)

    if indeks == -1 {
        systemMessage = "⛔ Donatur dengan ID tersebut tidak ditemukan\n"
        return
    }

    var i int
    for i = indeks; i < *jumlah-1; i++ {
        data[i] = data[i+1]
    }
    *jumlah--
    systemMessage = "📑 Donatur berhasil dihapus\n"
}

func clear() {
    cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()

    if systemMessage != "" {
        fmt.Println(systemMessage)
        systemMessage = ""
    }
}