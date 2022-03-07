# Compiler

[![Go](https://github.com/iskorotkov/compiler/actions/workflows/go.yml/badge.svg)](https://github.com/iskorotkov/compiler/actions/workflows/go.yml)

Простой компилятор, разработанный в рамках курса "Формальные грамматики и методы трансляции".

## Запуск

Для сборки данного компилятора требуется [Go версии 1.18 или более](https://go.dev/dl/) (Beta/RC подходят).

После установки Go 1.18 выполните в терминале:

```shell
go run cmd/compiler/main.go
```

## Архитектура

### Взаимодействие модулей

Планируется организовать взаимодействие различных модулей компилятора посредством передачи сообщений через каналы Go между горутинами. Каждый модуль запускается в отдельной горутине и работает параллельно с другими модулями, при этом он получает от предыдущего модуля данные по каналу, переданному ему в качестве параметра, а результаты своей работы передает в другой канал, который может использоваться другим модулем компилятора. Таким образом, работа компилятора организована по принципу конвейера.

![img.png](docs/img/goroutines.png)

Каждый модуль работает в отдельной горутине, и работа компилятора завершается, когда заканчивают свою работу все горутины (типично это должен быть последний модуль компилятора, осуществляющий вывод ошибок или результатов компиляции в зависимости от ее успеха).

На данный момент предлагается такой набор модулей:

[Reader (модуль ввода)](#reader) =>
[Scanner (сканер, лексический анализатор)](#scanner) =>
[Syntax analyzer (синтаксический анализатор)](#syntax-analyzer) =>
[Semantic analyzer (семантический анализатор)](#semantic-analyzer) =>
[Code generator (генератор кода)](#code-generator).

### Обработка ошибок

Любой из модулей может прервать работу конвейера при возникновении критической ошибки. В остальных случаях модуль передает ошибку дальше для обработки следующим модулем. Это возможно благодаря тому, что по каналам передаются не обязательно только результаты работы в случае успеха, а [дизъюнктивное объединение](https://ru.wikipedia.org/wiki/%D0%A2%D0%B8%D0%BF-%D1%81%D1%83%D0%BC%D0%BC%D0%B0) результата успеха и ошибки.

Требуется спроектировать вывод ошибок компиляции пользователю и механизм прерывания компиляции. Это планируется сделать после реализации лексического анализатора (т. к. тогда будет понятнее, какие ошибки и где возникают, а также будет возможность протестировать выбранный подход уже на практике).

### Reader

[Reader](internal/reader/reader.go) читает файл построчно и каждую строку разбивает на части - литералы ([Literal](internal/data/literal/literal.go)). При этом все знаки препинания и даже переход на новую строку также сохраняются как отдельные части для последующего анализа в следующих модулях.

Reader указывает для каждого литерала строку, начальный и конечный столбец (начиная индексирование с 1) для того, чтобы в последующих модулях можно было указать на возникшую ошибку на основе этих данных в литерале.

Reader использует регулярные выражения для нахождения [несловообразующих символов](https://docs.microsoft.com/ru-ru/dotnet/standard/base-types/character-classes-in-regular-expressions). Это позволяет разбить прочитанную строку по ним и передать по частям следующим модулям.

### Scanner

[Scanner](internal/analyzers/scanner/scanner.go) читает литералы ([Literal](internal/data/literal/literal.go)) из канала, распознает их и записывает соответствующие им токены ([Token](internal/data/token/token.go)) в выходной канал с токенами. При получении ошибки из Reader она передается как есть далее. При невозможности распознать токен Scanner передает дальше ошибку и продолжает обработку следующих литералов.

Распознавание числовых и булевых констант производится с помощью регулярных выражений. Строковые константы на данный момент не поддерживаются.

Распознавание пользовательских идентификаторов производится с помощью регулярного выражения.

### Syntax analyzer

### Semantic analyzer

### Code generator

## Тестирование

В процессе реализации модулей компилятора к ним также пишутся юнит-тесты для проверки их корректной работы. Часть тестов для модуля Reader написаны классическим способом, но большая часть тестов используют снапшоты для сравнения результатов выполнения.

При первом выполнении теста сохраненного снапшота нет и поэтому тест всегда проходит и сохраняет снапшот результата вызова тестируемой функции в файл. При последующих запусках снапшот читается из файла и сравнивается с результатов вызова тестируемой функции, и если есть какие-либо различия, то тест отмечается проваленным. Данный подход позволяет избежать времязатратного описания ожидаемого результата в коде, т. к. сравниваются эталонное решение из снапшота и текущий результат. Вручную проверить эталонное решение, записанное в файле, намного проще, чем описывать его в коде и менять каждый раз при изменении особенностей реализации модулей.

## CI/CD

При каждом пуше в ветку `main` GitHub собирает и тестирует компилятор. Для каждого PR GitHub также собирает и тестирует компилятор, что позволяет обнаружить ошибки в реализации до слияния ветки.
