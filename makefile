

build:
	GO_CPPFLAGS="-I/opt/homebrew/Cellar/tesseract/5.3.4_1/include -I/opt/homebrew/Cellar/leptonica/1.84.1/include" \
  CGO_LDFLAGS="-L/opt/homebrew/Cellar/tesseract/5.3.4_1/lib -L/opt/homebrew/Cellar/leptonica/1.84.1/lib -ltesseract -lleptonica" \
  GOOS=darwin GOARCH=arm64 go build -o output_file main.go




build-cpp:
	clang++ main.cpp -o ocr_cpp -std=c++11 -I/opt/homebrew/Cellar/tesseract/5.3.4_1/include -I/opt/homebrew/Cellar/leptonica/1.84.1/include -L/opt/homebrew/Cellar/tesseract/5.3.4_1/lib -L/opt/homebrew/Cellar/leptonica/1.84.1/lib -ltesseract -lleptonica
	./ocr_cpp

