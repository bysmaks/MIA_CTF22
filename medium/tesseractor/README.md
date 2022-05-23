# tesseractor
### Сложность
Средне.
### Описание
При обыске ничего найти не удалось, но подозреваемый пытался избавиться от этой книги. Прислал тебе её ксерокопию, посмотри что сможешь в ней откопать.
### Решение
пример на pytesseract:

import os
from PIL import Image
from pdf2image import convert_from_path
import pytesseract
import threading

def process(pages):

    for page_number, page_data in enumerate(pages):

        txt = pytesseract.image_to_string(page_data).encode("utf-8")
        open('output.txt','ab').write(txt)
        print(f'Processed {page_number}')

filePath = 'task.pdf'
doc = convert_from_path(filePath)
path, fileName = os.path.split(filePath)
fileBaseName, fileExtension = os.path.splitext(fileName)

open('output.txt','wb').close()

for i in range(0,len(doc),10):
    thread = threading.Thread(target=process(doc[i:i+10]))
    thread.start()

### Флаг
CTF{hOp3_yOu_r3cOgniz3_th4t}
