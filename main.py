import pytesseract
from PIL import Image
# Buka gambar
img = Image.open('ktp.jpeg')
# Lakukan OCR pada gambar
text = pytesseract.image_to_string(img)
# Tampilkan teks yang terdeteksi
print(text)

