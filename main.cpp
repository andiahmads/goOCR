#include <iostream>
#include <leptonica/allheaders.h>
#include <sstream>
#include <string>
#include <tesseract/baseapi.h>
#include <vector>

int main() {
  // Inisialisasi Tesseract OCR
  tesseract::TessBaseAPI tess;
  if (tess.Init(NULL, "eng")) {
    std::cerr << "Error: Unable to initialize Tesseract." << std::endl;
    return 1;
  }

  // Load gambar
  Pix *image = pixRead("ktp.jpeg");
  if (!image) {
    std::cerr << "Error: Unable to load image." << std::endl;
    return 1;
  }

  // Set gambar input untuk Tesseract
  tess.SetImage(image);

  // Recognize teks dari gambar
  char *text = tess.GetUTF8Text();
  if (!text) {
    std::cerr << "Error: Unable to recognize text." << std::endl;
    pixDestroy(&image); // Hapus memori image sebelum keluar dari program
    return 1;
  }

  // Tampilkan teks yang terdeteksi
  std::cout << "x Detected Text: " << text << std::endl;

  std::string ocrText = text;

  // Hapus memori yang dialokasikan
  delete[] text;
  pixDestroy(&image);

  // Vector untuk menyimpan kata-kata hasil OCR
  std::vector<std::string> words;

  // String stream untuk memproses teks
  std::istringstream iss(ocrText);
  std::string word;

  // Memisahkan kata-kata dan menyimpannya dalam vector
  while (iss >> word) {
    words.push_back(word);
  }

  // Menampilkan kata-kata yang terpisah
  std::string findText;
  for (const auto &w : words) {
    std::cout << w.c_str() << std::endl;
    // findText = w.begin();
  }
  std::cout << findText << std::endl;

  return 0;
}
