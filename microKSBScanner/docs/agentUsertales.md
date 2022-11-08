# microKSBScanner

## Описание

Сканнер состоит из двух компонент:

- Микросервиса, реализующего интерфейс активов для сервиса интеграций. Сервис интеграций обращается к этому микросервису за списком активов.
- Агента, непосредственно сканирующего компьютер на котором он установлен.

Агент при старте начинает сканировать компьютер. Полученную информацию держит в памяти. Либо по расписанию перепроверяет информацию, либо по возможности сделать детекторы изменений информаций. При наличии изменений с прошлого сканирования соединяется с Микросервисом и отправляет ему новые данные.

Микросервис получив команду на сбор данных отдает имеющуюся у него информацию от Агентов.

Регистрация Агента в Микросервисе. При старте Агента указывается адрес Микросервиса. Агент отправляет информацию о себе и записывает ответ что зарегистрирован.

## Схемы

### Агента (Agent)

```puml
object Explorer
object Storer
object Communicator
Communicator -> Explorer
Communicator -d-> Storer
```

### Микросервиса (Scanner)

```puml
object Communicator
object Storer
object Agenter
Communicator -> Agenter
Communicator -d-> Storer
```

### Общая

```puml
object Intgsrv {
    Сервис интеграции
}
object Scanner
object Agent1
object Agent2
object Agent3

Intgsrv -> Scanner
Scanner -d-> Agent1
Scanner -d-> Agent2
Scanner -d-> Agent3
```

## Прочее

Идея: сделать возможность каскадных Микросервисов по типу прокси для сложных схем организации сетей.
Или просто второй добавлять на уровне первого?

```puml
object Intgsrv {
    Сервис интеграции
}

object Scanner
object Agent1
object Agent2
object Agent3

Intgsrv -> Scanner

Scanner -d-> Agent1
Scanner -d-> Agent2
Scanner -d-> Agent3

object Scanner_A
object Agent1_A
object Agent2_A
object Agent3_A

Scanner -> Scanner_A

Scanner_A -d-> Agent1_A
Scanner_A -d-> Agent2_A
Scanner_A -d-> Agent3_A

```
