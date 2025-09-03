package seeders

import (
	"log"
	"time"

	"loka-kasir/entity"

	"gorm.io/gorm"
)

func SeedBusinessTypes(db *gorm.DB) {
	businessTypes := []entity.BusinessType{
		{
			Name:        "Makanan dan Minuman",
			Code:        "FOOD_BEVERAGE",
			Description: "Restoran, kafe, warung makan, bakery, katering, dan bisnis kuliner lainnya.",
		},
		{
			Name:        "Minimarket/Toko Kelontong/Warung",
			Code:        "MINIMARKET",
			Description: "Usaha ritel skala kecil seperti minimarket, toko kelontong, warung sembako.",
		},
		{
			Name:        "Online Shop",
			Code:        "ONLINE_SHOP",
			Description: "Toko daring yang berjualan melalui website, marketplace, atau media sosial.",
		},
		{
			Name:        "Laundry",
			Code:        "LAUNDRY",
			Description: "Bisnis jasa cuci pakaian, dry cleaning, dan setrika.",
		},
		{
			Name:        "Persewaan/Rental/Studio Foto",
			Code:        "RENTAL",
			Description: "Bisnis sewa kendaraan, alat, properti, atau jasa studio foto.",
		},
		{
			Name:        "Jasa/Service Lainnya",
			Code:        "SERVICE",
			Description: "Bisnis jasa umum seperti bengkel, barbershop, salon, konsultan, dan lain-lain.",
		},
		{
			Name:        "Kesehatan/Apotek/Klinik",
			Code:        "HEALTHCARE",
			Description: "Bisnis layanan kesehatan, apotek, klinik, atau toko alat medis.",
		},
		{
			Name:        "Fashion/Pakaian/Butik",
			Code:        "FASHION",
			Description: "Bisnis pakaian, butik, distro, sepatu, dan aksesori fashion.",
		},
		{
			Name:        "Toko Elektronik/Gadget",
			Code:        "ELECTRONICS",
			Description: "Bisnis ritel gadget, komputer, smartphone, dan aksesoris elektronik.",
		},
		{
			Name:        "Toko Buku/ATK",
			Code:        "BOOKS_STATIONERY",
			Description: "Bisnis penjualan buku, alat tulis kantor, fotokopi, dan perlengkapan belajar.",
		},
		{
			Name:        "Pertanian/Peternakan",
			Code:        "AGRICULTURE",
			Description: "Usaha di bidang pertanian, perkebunan, peternakan, dan hasil bumi.",
		},
		{
			Name:        "Perikanan/Maritim",
			Code:        "FISHERY",
			Description: "Usaha perikanan tangkap, budidaya, hasil laut, dan maritim.",
		},
		{
			Name:        "Industri/Manufaktur",
			Code:        "MANUFACTURING",
			Description: "Pabrik, industri pengolahan, produksi barang skala besar.",
		},
		{
			Name:        "Pendidikan/Kursus/Bimbel",
			Code:        "EDUCATION",
			Description: "Sekolah, kursus, bimbingan belajar, dan pelatihan keterampilan.",
		},
		{
			Name:        "Hiburan/Tempat Wisata",
			Code:        "ENTERTAINMENT",
			Description: "Bisnis hiburan, bioskop, tempat wisata, wahana permainan, dan event.",
		},
		{
			Name:        "Transportasi/Logistik",
			Code:        "TRANSPORT",
			Description: "Jasa transportasi, ekspedisi, pengiriman barang, dan logistik.",
		},
		{
			Name:        "Kecantikan/Salon/Spa",
			Code:        "BEAUTY",
			Description: "Bisnis salon, spa, klinik kecantikan, dan perawatan tubuh.",
		},
	}

	for _, bt := range businessTypes {
		var existing entity.BusinessType
		err := db.Where("code = ?", bt.Code).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				bt.CreatedAt = time.Now()
				bt.UpdatedAt = time.Now()
				if err := db.Create(&bt).Error; err != nil {
					log.Printf("Gagal menambahkan business type %s: %v", bt.Name, err)
				} else {
					log.Printf("Business type %s berhasil ditambahkan", bt.Name)
				}
			} else {
				log.Printf("Gagal memeriksa business type %s: %v", bt.Name, err)
			}
		} else {
			log.Printf("Business type %s sudah ada, skip seeding", bt.Name)
		}
	}
}
