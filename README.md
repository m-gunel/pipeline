Аналог Unix Pipeline на Go

Это решение задачи, в котором реализован аналог Unix-пайплайна с использованием каналов в Go. Я создаю конвейер обработки данных, где результаты работы каждой функции передаются как входные данные для следующей.

Описание

Задача состоит из нескольких этапов:

Написание функции ExecutePipeline, которая организует конвейерную обработку функций.
Написание нескольких функций, которые вычисляют хеш-сумму входных данных через несколько этапов.


Расчет хеш-суммы реализован следующей цепочкой:

SingleHash считает значение crc32(data)+"~"+crc32(md5(data)) ( конкатенация двух строк через ~), где data - то что пришло на вход (по сути - числа из первой функции)  
MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки), где th=0..5 ( т.е. 6 хешей на каждое входящее значение ), потом берёт конкатенацию результатов в порядке расчета (0..5), где data - то что пришло на вход (и ушло на выход из SingleHash)  
CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/), объединяет отсортированный результат через _ (символ подчеркивания) в одну строку  
crc32 считается через функцию DataSignerCrc32  
md5 считается через DataSignerMd5

В чем подвох:

DataSignerMd5 может одновременно вызываться только 1 раз, считается 10 мс. Если одновременно запустится несколько - будет перегрев на 1 сек  
DataSignerCrc32, считается 1 сек  
На все расчеты у есть 3 сек.  
Если делать в лоб, линейно - для 7 элементов это займёт почти 57 секунд, следовательно надо это как-то распараллелить
