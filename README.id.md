<p align="center"><a href="README.md">English</a> · <a href="README.zh-CN.md">简体中文</a> · <strong>Bahasa Indonesia</strong></p>

<p align="center">
  <img src="brand/icon.png" width="112" height="112" alt="Ikon MyTube: lima bilah suara bergradasi biru ke ungu di atas ubin gelap membulat" />
</p>

<h1 align="center">MyTube</h1>

<p align="center"><strong>Satu pemutar untuk YouTube Music dan musik di komputer Anda.</strong></p>

<p align="center">Pemutar YouTube Music dan musik lokal untuk Windows. YouTube Music bisa dipakai tanpa akun. Satu pencarian, satu antrean, playlist campuran, lirik tersinkron, tema, dan plugin.</p>

<p align="center"><sub>Ini terjemahan dari <a href="README.md">README berbahasa Inggris</a>. Jika ada perbedaan, yang terbaru adalah versi bahasa Inggris.</sub></p>

<p align="center">
  <a href="https://github.com/bigdaddy284/mymusic/actions/workflows/ci.yml"><img src="https://github.com/bigdaddy284/mymusic/actions/workflows/ci.yml/badge.svg" alt="CI" /></a>
  <a href="https://github.com/bigdaddy284/mymusic/actions/workflows/pages.yml"><img src="https://github.com/bigdaddy284/mymusic/actions/workflows/pages.yml/badge.svg" alt="Situs web" /></a>
  <a href="https://github.com/bigdaddy284/mymusic/releases/latest"><img src="https://img.shields.io/github/v/release/bigdaddy284/mymusic" alt="Rilis terbaru" /></a>
  <a href="https://github.com/bigdaddy284/mymusic/releases"><img src="https://img.shields.io/github/downloads/bigdaddy284/mymusic/total" alt="Unduhan" /></a>
  <a href="https://pkg.go.dev/github.com/bigdaddy284/mymusic"><img src="https://pkg.go.dev/badge/github.com/bigdaddy284/mymusic.svg" alt="Referensi Go" /></a>
  <a href="go.mod"><img src="https://img.shields.io/github/go-mod/go-version/bigdaddy284/mymusic" alt="Versi Go" /></a>
  <img src="https://img.shields.io/badge/platform-Windows-lightgrey" alt="Windows" />
  <a href="#lisensi"><img src="https://img.shields.io/badge/license-MIT%20or%20GPL--3.0-blue" alt="MIT atau GPL-3.0" /></a>
</p>

<p align="center">
  <a href="https://bigdaddy284.github.io/mymusic/">Situs web</a> ·
  <a href="https://github.com/bigdaddy284/mymusic/releases/latest">Unduh</a> ·
  <a href="https://bigdaddy284.github.io/mymusic/docs/">Dokumentasi</a> ·
  <a href="CHANGELOG.md">Catatan perubahan</a>
</p>

<p align="center">
  <img src="docs/banner.webp" alt="MyTube 0.5.9 dengan tampilan bawaannya, tata letak Blend dengan tema Mix, berdiri tegak di atas panggung gelap dan menghadap ke arah Anda: di tengah, アイドル oleh YOASOBI sedang diputar, bergaris tepi merah muda, dengan romaji di bawah setiap baris liriknya dan pemutar mini di depannya; di kiri sebuah pencarian, di kanan mix di pustaka" />
</p>

## Mengapa MyTube

Jika musik Anda terbagi antara YouTube Music dan file di komputer Anda,
biasanya Anda perlu dua aplikasi, dengan dua pencarian, dua antrean, dan
dua pustaka. MyTube menyatukan keduanya dalam satu pemutar:

- **Satu pencarian.** MyTube mencari di YouTube Music dan musik Anda
  sendiri sekaligus, dan menandai setiap hasil dengan asalnya.
- **Satu antrean dan satu pustaka.** File di komputer Anda bisa diputar
  tepat setelah lagu YouTube Music, dan satu playlist bisa berisi keduanya.
- **Tidak perlu akun.** YouTube Music langsung bisa diputar begitu aplikasi
  dibuka. Masuk ke akun hanya jika Anda ingin memakai pustaka Anda sendiri.
- **Program kecil.** MyTube adalah satu program Go yang menggambar
  jendelanya sendiri di GPU. Tidak ada browser di dalamnya yang harus
  diunduh atau ditampung di memori.
- **Bebas diubah.** Ada sepuluh tema bawaan, Anda bisa menambahkan tema dan
  plugin sendiri, dan setiap perintah bisa diberi tombol pilihan Anda
  sendiri.
- **Privat.** MyTube tidak punya akun, server, analitik, atau telemetri
  sendiri.
- **Gratis dan open source,** di bawah lisensi MIT atau GPL 3.0.

## Fitur

### Cari dan putar

- **Satu pencarian.** Hasil dari YouTube Music dan musik di komputer Anda
  tampil berdampingan, masing-masing ditandai dengan asalnya.
  Hasil YouTube Music menunjukkan berapa kali setiap lagu diputar, sehingga
  versi aslinya mudah dibedakan dari versi cover.
- **Satu antrean.** Campur lagu YouTube Music dan file Anda sendiri, acak
  dan ulangi semuanya sekaligus, dan biarkan radio melanjutkan saat antrean habis.
- **Explore (Jelajahi).** Rilis baru, lagu yang sedang tren, tangga lagu,
  mood dan genre, serta lagu yang terkait dengan lagu yang sedang diputar.
- **Lirik tersinkron** yang mengikuti lagu, atau teks biasa jika Anda lebih
  suka.
- **Pemutar mini** yang tetap berada di atas jendela-jendela lain.

### Pustaka Anda

- **Mix:** playlist yang berisi lagu YouTube Music dan file Anda sendiri.
  Impor playlist YouTube Music ke dalam mix, lalu cari dan urutkan lagu di dalam mix mana
  pun.
- **Lagu yang disukai** dari YouTube Music dan komputer Anda dalam satu daftar.
- **Riwayat dan statistik mendengarkan,** dari hari ke hari.
- **Akun Anda.** Masuk ke YouTube Music lewat halaman milik Google sendiri
  untuk mendapatkan feed beranda, lagu yang disukai, playlist, album, dan
  artis Anda, baik sebagai akun Anda maupun salah satu channel brand milik
  akun itu. Lagu yang Anda sukai, album dan playlist yang Anda simpan, serta
  artis yang Anda ikuti juga tercatat di YouTube Music. Lagu yang Anda putar masuk ke
  riwayat YouTube Music Anda, yang menjadi dasar rekomendasinya.

### Suara

- Crossfade, kenyaringan yang setara antarlagu, dan perataan volume.
- Equalizer sepuluh band dengan preset.
- Visualizer.
- Cache offline untuk lagu-lagu yang pernah Anda putar.

### Di desktop Anda

- Tombol media dan kontrol media Windows.
- Tombol pemutar di thumbnail taskbar Windows.
- Kontrol di tray, notifikasi setiap kali lagu baru diputar, dan berjalan
  otomatis saat login.
- Tautan `mytube://`.
- **Last.fm:** setiap lagu yang Anda dengar tercatat di profil Anda, dan
  scrobble yang dibuat saat offline dikirim begitu Anda kembali online.
- **Discord:** profil Anda menampilkan lagu yang sedang diputar, lengkap
  dengan sampul dan sisa waktunya.
- **Pembaruan:** MyTube memeriksa rilis baru saat dibuka dan setiap
  beberapa jam, memeriksa tanda tangan digitalnya, lalu memasangnya saat
  Anda memintanya, atau secara otomatis jika Anda memilih opsi itu.

### Sesuaikan sesuka Anda

- **Tema.** Sepuluh tema bawaan dengan efek yang bisa diatur, satu tema
  untuk mode terang dan satu untuk mode gelap, serta warna aksen pilihan
  Anda sendiri atau yang diambil dari sampul lagu yang sedang diputar.
- **Plugin** menambahkan halaman, perintah, aksi lagu, dan widget lewat API
  bertipe, dan semuanya tampil dengan tema Anda.
- **Mengutamakan keyboard.** Pintasan untuk semua hal yang umum, tombol
  pilihan Anda sendiri untuk perintah apa pun (termasuk perintah dari
  plugin), dan palet perintah (<kbd>Ctrl</kbd> <kbd>K</kbd>) untuk sisanya.
- **File pengaturan.** Pengaturan Anda tersimpan sebagai `settings.json` di
  folder data Anda (`%AppData%\MyTube`), jadi bisa Anda simpan bersama dotfiles Anda. Simpan file itu saat
  MyTube berjalan, dan perubahannya langsung berlaku; file ini bisa dibuka
  lewat Settings › Desktop › Settings file (Pengaturan › Desktop › File
  pengaturan).

## Sekilas tampilan

<p align="center">
  <img src="docs/look-listen.webp" alt="Halaman-halaman MyTube dengan tata letak Blend, berdiri tegak di atas panggung gelap dan menghadap ke arah Anda: di tengah, アイドル oleh YOASOBI sedang diputar, bergaris tepi merah muda, setiap baris liriknya dengan romaji dan terjemahan bahasa Inggris di bawahnya; di kiri album THE BOOK 3, di kanan halaman podcast yang diikuti" />
</p>

<p align="center"><sub><b>Dengarkan.</b> Di tengah, アイドル oleh YOASOBI sedang diputar: setiap baris menyala saat dinyanyikan, dengan romaji dan terjemahan di bawahnya. Di kiri, sebuah album; di kanan, podcast yang Anda ikuti.</sub></p>

<p align="center">
  <img src="docs/look-yours.webp" alt="Bagian lain MyTube, berdiri tegak di atas panggung gelap dan menghadap ke arah Anda: di tengah, Settings › Appearance (Pengaturan › Tampilan) dengan keempat tata letak, bergaris tepi merah muda; di kiri, tema Material yang diwarnai sampul lagu yang diputar; di kanan, tema komunitas di marketplace" />
</p>

<p align="center"><sub><b>Sesuaikan sesuka Anda.</b> Di tengah, keempat tata letak: Blend, yang bawaan, Expressive, Crate, dan Tiles, di atas tema-temanya. Di kiri, tema Material yang mengambil warnanya dari sampul lagu yang diputar. Di kanan, tema komunitas di marketplace.</sub></p>

## Pemasangan

MyTube berjalan di Windows 10 dan 11. Unduh dari
[rilis terbaru](https://github.com/bigdaddy284/mymusic/releases/latest):

| Windows | Unduhan |
|---|---|
| 10 dan 11 (x64, Arm), installer | [`mytube-<version>-windows-<arch>-setup.exe`](https://github.com/bigdaddy284/mymusic/releases/latest) |
| Portabel, tanpa instalasi | [`mytube-<version>-windows-<arch>-portable.zip`](https://github.com/bigdaddy284/mymusic/releases/latest) |

`<arch>` adalah `amd64` untuk x64 dan `arm64` untuk Arm. Setiap rilis
menyertakan `SHA256SUMS` untuk memeriksa file yang Anda unduh. Anda tidak
perlu file bernama `update-…`, `….delta`, atau `-windows-<arch>.tar.gz`:
file-file itu dipakai aplikasi yang sudah terpasang untuk memperbarui
dirinya.

Installer dan zip portabel bisa memperbarui dirinya sendiri. Untuk zip
portabel, ekstrak di mana saja lalu buka `MyTube.exe`.

### Pertama kali dijalankan

- Installer ini belum ditandatangani, jadi SmartScreen mungkin memberi
  peringatan: pilih More info (Info selengkapnya), lalu Run anyway (Tetap
  jalankan).
- Masuk ke YouTube Music memerlukan
  WebView2 Runtime, yang sudah ada di Windows 11. Di Windows 10, pasang
  dari [Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/).
- MyTube menemukan sendiri audio sebagian besar lagu YouTube Music. Untuk
  sisanya, seperti lagu yang dibuat untuk anak-anak, MyTube memakai
  [yt-dlp](https://github.com/yt-dlp/yt-dlp), yang dipasang MyTube untuk
  Anda dari Settings › Playback (Pengaturan › Pemutaran) dan selalu
  diperbaruinya. yt-dlp memerlukan runtime JavaScript:
  [Deno](https://deno.com), Node.js, Bun, atau QuickJS.

### Akun (semuanya opsional)

- **YouTube Music:** pilih Sign in to YouTube Music (Masuk ke YouTube Music)
  di Settings › Accounts (Pengaturan › Akun), lalu masuk ke Google di jendela
  yang terbuka. Apa yang Anda ketik hanya dikirim ke Google. MyTube
  menyimpan sesi itu dalam keadaan terenkripsi, dengan kunci di Credential
  Manager (Manajer Kredensial) Windows.
- **Last.fm:** buat akun API di
  [last.fm/api/account/create](https://www.last.fm/api/account/create)
  (nama apa saja, tanpa callback URL), tempelkan API key dan shared secret
  dari akun itu di Settings › Accounts, lalu setujui MyTube di browser
  Anda.
- **Discord:** selama aplikasi Discord berjalan, profil Anda menampilkan apa
  yang sedang Anda dengarkan. Matikan di Settings › Desktop (Pengaturan ›
  Desktop).

Ada yang tidak berfungsi? [docs/troubleshooting.md](docs/troubleshooting.md)
membahas masalah yang umum terjadi dan masalah yang sudah diketahui.

## Membangun dari kode sumber

Anda memerlukan [Go](https://go.dev/dl/) 1.27.3 atau yang lebih baru,
[Bun](https://bun.sh) 1.2 atau yang lebih baru, yang menjalankan CLI MyGo,
dan Git. MyGo tidak memerlukan compiler C, jadi Linux juga bisa
membangun aplikasi Windows secara silang, seperti saat membuat rilis.

```powershell
winget install GoLang.Go
winget install Oven-sh.Bun
winget install Git.Git
```

Buka terminal baru, lalu:

```powershell
git clone https://github.com/bigdaddy284/mymusic
cd mytube
bun install
bun run dev
```

Fitur masuk ke YouTube Music memerlukan WebView2 Runtime, yang
sudah ada di Windows 11. Di Windows 10, pasang dari
[Microsoft](https://developer.microsoft.com/microsoft-edge/webview2/).

### Nix

Dengan [Nix](https://nixos.org) di Linux, `nix develop` menyediakan
alat-alat yang dipakai untuk membangun rilis, pada versi yang dikunci oleh
`flake.lock`: Go, Bun, NSIS, osslsigncode, GoReleaser, dan Zola.

### Pengembangan

`bun run dev` membangun aplikasi, menjalankannya, dan membangunnya ulang
setiap kali Anda mengedit kode. Datanya terpisah dari MyTube yang
terpasang, dengan nama `MyTube Dev`.

```sh
bun run build:plugins          # the plugins' JavaScript, from their TypeScript
go test ./...                  # the app's tests
bun run typecheck              # the plugin SDK and the example plugins
bun run build                  # the installer and the app, in build/
bun run build -- -platform windows/amd64,windows/arm64
```

Di shell Nix, buat file unduhan dengan
`.github/scripts/build.sh windows/amd64,windows/arm64`, seperti saat membuat
rilis, dan `.github/scripts/package.sh` membuat zip portabelnya.

[Situs web](https://bigdaddy284.github.io/mymusic/) ada di `site/`, dibangun
dengan [Zola](https://www.getzola.org), dan panduan di `docs/` menjadi
halaman `/docs/` di situs itu. `go run ./tools/sitedocs && zola --root site serve`
menampilkan pratinjaunya. Rilis dibangun di shell Nix dan diterbitkan oleh
[GoReleaser](https://goreleaser.com) (`.goreleaser.yaml`).

[docs/architecture](docs/architecture/README.md) menjelaskan bagaimana
bagian-bagiannya saling terhubung.

### Struktur kode

Aplikasinya ada di package root, dengan package-package di `internal/` di
bawahnya. Package di luar `internal/` boleh diimpor oleh siapa saja, dan
didokumentasikan di [pkg.go.dev](https://pkg.go.dev/github.com/bigdaddy284/mymusic).

```
*.go                 the app and its windows: the pages, the player bar, the
                     mini player, the tray, media keys, deep links
fonts/               the typefaces the windows draw with
internal/app         the services the windows call, and the desktop integration
internal/player      the queue, played by the audio engine or on a cast device
internal/audio       MyTube's audio engine, and the system's audio output
internal/plugins     plugins, each in a JavaScript engine of its own
internal/youtube     YouTube Music's catalog, and your account there
internal/stream      audio streams, from YouTube or yt-dlp, opened for the engine
internal/match       the songs of YouTube Music's videos
internal/library     liked songs, mixes, history and stats
internal/lastfm      scrobbling to Last.fm
internal/discord     what plays, on your Discord profile
internal/diag        the log, without secrets, and the reports of crashes
internal/addons      themes and plugins
builtin/themes       the themes MyTube ships with
public/              the page of the picture of music videos
sdk/                 the types of the plugin API
examples/            example themes and plugins
```

Package yang bisa Anda impor:

```
opus/                an Opus decoder in pure Go, a port of libopus's
aac/                 an AAC-LC decoder in pure Go
webm/                the audio of WebM files: packets, times and seeking
mp4/                 the audio of MP4 files: packets, times and seeking
resample/            sample rate conversion, with Kaiser-windowed sincs
lyrics/              synced lyrics from LRCLIB, and the LRC format
brand/               MyTube's logo, for code that draws it
```

## Tema dan plugin

Tema dan plugin dimuat dari folder add-on Anda, dan Anda bisa membuatnya
sendiri. [docs/addons.md](docs/addons.md) menjelaskan cara menambahkannya.

### Membuat tema

Tema adalah file `theme.json` yang mengatur warna, font, dan bentuk sudut.
[docs/themes.md](docs/themes.md) mencantumkan setiap token dan memberi
contoh tema yang lengkap. Untuk mencoba tema Anda, letakkan foldernya di
folder tema, yang bisa dibuka lewat Settings › Appearance › open themes
folder (Pengaturan › Tampilan › buka folder tema), lalu pilih tema itu di
Settings › Appearance.

### Membuat plugin

Plugin adalah modul kecil, ditulis dalam TypeScript atau JavaScript, yang
menambahkan halaman, perintah, aksi lagu, dan widget.
[docs/plugins.md](docs/plugins.md) mendokumentasikan API-nya, beserta
[definisi tipenya](sdk/index.d.ts) untuk editor Anda. Plugin
[sleep timer](examples/plugins/sleep-timer) dan
[spectrum](examples/plugins/spectrum) di `examples/` adalah contoh yang benar-benar
berfungsi.

## Privasi

MyTube tidak punya akun, server, analitik, atau telemetri sendiri.
MyTube berkomunikasi dengan:

- **YouTube dan Google:** katalog YouTube Music, stream audionya, dan
  tingkat kenyaringan setiap lagu, yang diminta tanpa login, serta gambar
  sampul; sebagian stream lewat yt-dlp. Jika Anda masuk: pustaka Anda di
  sana, dan lagu yang Anda putar, yang masuk ke riwayat YouTube Music Anda,
  kecuali Anda mematikannya di Settings › Accounts (Pengaturan › Akun).
- **LRCLIB:** judul, artis, album, dan durasi lagu yang sedang diputar,
  untuk mencari liriknya.
- **GitHub:** pemeriksaan versi baru saat aplikasi dibuka dan setiap
  beberapa jam, yang bisa Anda matikan di Settings › About (Pengaturan ›
  Tentang); dan unduhan yt-dlp.
- **Last.fm,** setelah Anda menghubungkannya: lagu yang Anda putar.
- **Discord,** kecuali Anda mematikannya: lagu yang sedang diputar, dikirim
  ke aplikasi Discord di komputer Anda, yang kemudian menampilkannya di
  profil Anda.

Yang disimpan MyTube di komputer Anda:

- Pengaturan dan pustaka Anda (lagu yang disukai, mix, riwayat, dan
  pencarian terakhir), di folder yang disebutkan di Settings › About.
- Lagu yang pernah Anda putar, di cache, sebanyak yang Anda izinkan di
  Settings › Playback (Pengaturan › Pemutaran), beserta tingkat
  kenyaringannya menurut pengukuran YouTube.
- Sampul yang pernah Anda lihat, hingga 256 MB.
- Data login Anda, termasuk sesi YouTube Music, di Credential Manager
  (Manajer Kredensial). Jika Credential Manager tidak bisa menyimpannya,
  data itu disimpan di file yang dienkripsi untuk akun Windows Anda.
- Log, tanpa data login Anda.

Halaman riwayat di pustaka bisa menghapus daftar lagu yang pernah Anda
putar, Settings › Playback bisa menghapus cache dan sampul, dan keluar dari
sebuah layanan akan menghapus data login layanan tersebut.

Visualizer menganalisis audio yang diputar MyTube di memori dan tidak
pernah merekam atau mengirimkannya.

## Kredit

MyTube dibangun di atas proyek-proyek berikut:

- **[MyGo](https://github.com/egoist/mygo)**, framework desktop yang dipakai
  untuk membangun MyTube. MyTube memakai UI native MyGo: MyGo menggambar
  sendiri setiap jendela, dengan Go, tanpa webview, mulai dari halaman dan
  bilah pemutar hingga pemutar mini. CLI MyGo mengemas aplikasi ini,
  lengkap dengan installer dan pembaruan bertanda tangan yang dipakai
  aplikasi untuk memperbarui dirinya.
- **[Go](https://go.dev)**, bahasa pemrograman aplikasi ini, bersama
  [purego](https://github.com/ebitengine/purego), yang memanggil library
  sistem tanpa cgo, serta [golang.org/x/image](https://pkg.go.dev/golang.org/x/image)
  dan [golang.org/x/text](https://pkg.go.dev/golang.org/x/text).
- **[Sobek](https://github.com/grafana/sobek)**, mesin JavaScript tempat
  setiap plugin berjalan.
- **[Nix](https://nixos.org)**, yang mengunci versi alat-alat untuk
  membangun rilis.
- **[GoReleaser](https://goreleaser.com)**, yang menerbitkan rilis.
- **[Bun](https://bun.sh)**, yang menjalankan proses build, dan
  **[TypeScript](https://www.typescriptlang.org)**, yang mengompilasi plugin
  dan memeriksa tipe pada SDK plugin.
- **[Zola](https://www.getzola.org)**, yang membangun situs web.
- **[NSIS](https://nsis.sourceforge.io)**, yang membuat installer Windows.
- **[yt-dlp](https://github.com/yt-dlp/yt-dlp)**, yang menemukan audio
  lagu-lagu yang tidak bisa ditemukan MyTube sendiri.
- **[LRCLIB](https://lrclib.net)**, database terbuka berisi lirik
  tersinkron.
- **[libopus](https://opus-codec.org)**, karena decoder Opus MyTube adalah
  port dari decodernya, dan **[SharpJaad](https://github.com/jimm98y/SharpJaad)**,
  port C# dari [JAAD](https://sourceforge.net/projects/jaadec/) yang
  menjadi dasar decoder AAC.
- **[FFmpeg](https://ffmpeg.org)**, yang menjadi pembanding dalam pengujian
  decoder.
- **[Selawik](https://github.com/microsoft/Selawik)** dan
  **[Quicksand](https://github.com/andrew-paglinawan/QuicksandFamily)**,
  jenis huruf yang dipakai di jendela aplikasi.
- **[Kevin MacLeod](https://incompetech.com)**, yang musiknya, di bawah
  lisensi [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/),
  diputar di situs web.

## Berkontribusi

Perbaikan, fitur, tema, plugin, dan laporan bug, semuanya diterima dengan
senang hati.

- **Laporkan masalah** dari Settings › About › report a problem
  (Pengaturan › Tentang › laporkan masalah). Ini akan membuka issue GitHub baru yang
  sudah terisi versi, sistem, tema, dan plugin Anda. Atau
  [buka issue](https://github.com/bigdaddy284/mymusic/issues/new) sendiri.
- **Kirim perubahan** sebagai pull request.
  [CONTRIBUTING.md](CONTRIBUTING.md) menjelaskan pemeriksaan yang perlu
  dijalankan, gaya commit, dan entri changelog yang diperlukan untuk
  perubahan yang akan disadari pengguna.

[CONTRIBUTORS.md](CONTRIBUTORS.md) berterima kasih kepada orang-orang yang
kode, ide, dan laporannya ikut membentuk MyTube.

## Lisensi

MyTube memiliki lisensi ganda: gunakan di bawah [lisensi MIT](LICENSE-MIT)
atau [GNU General Public License, versi 3](LICENSE-GPL), mana pun yang Anda
pilih. Kontribusi diterima di bawah kedua lisensi tersebut. Jenis huruf di
`fonts` berada di bawah SIL Open Font License:
[lisensi Selawik](fonts/OFL.txt) dan
[lisensi Quicksand](fonts/quicksand-OFL.txt). Decoder Opus adalah port dari
libopus, di bawah [lisensi BSD-nya](opus/LICENSE.libopus), dan decoder AAC
berasal dari SharpJaad, di bawah [lisensi MIT-nya](aac/LICENSE.sharpjaad).

MyTube tidak berafiliasi dengan, tidak didukung secara resmi oleh, dan
tidak disponsori oleh Google. MyTube memutar apa yang disediakan YouTube
Music untuk penggunanya, melalui antarmukanya sendiri; gunakan MyTube
sesuai dengan yang diizinkan oleh ketentuan layanan tersebut. YouTube Music
adalah merek dagang milik pemiliknya.
