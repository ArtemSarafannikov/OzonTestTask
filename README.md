# Система управления постами и комментариями на GraphQL

## Введение
Данный проект представляет собой высокопроизводительное решение для управления постами и комментариями с использованием GraphQL, обеспечивая гибкость выборки данных и поддержку реального времени через GraphQL Subscriptions. Система разработана на Go и предназначена для использования как с PostgreSQL, так и с хранилищем в оперативной памяти (in-memory).

## Функциональность сервиса
На сайте присутствует система регистрации и авторизации пользователей, просмотр, создание и редактирование постов, создание комментариев, подписка на пост для отслеживания новых комментариев, а также возможность отвечать на комментарии.

Без авторизации пользователь может:
- Просматривать список постов
- Просматривать список комментариев под постом
- Подписаться на новые комментарии под постом

С авторизацией у пользователя появляются возможности:
- Создать пост
- Редактировать созданные посты
- Оставлять комментарии под постом


### Посты

- Просмотр списка постов с возможностью пагинации и фильтром по автору с комментариями.

- Просмотр детальной информации о конкретном посте вместе с комментариями.

- Автор может редактировать свои посты (изменить заголовок, содержание, а также запретить комментарии).

- Каждый пост имеет временную метку создание я редактирования, если пост не редактировался, то метка редактирования будет null.


### Комментарии

- Вложенная структура комментариев с обработкой глубины вложенности.

- Ограничение максимальной длины текста (до 2000 символов).

- Пагинация комментариев для эффективного извлечения данных.

- Возможность ответа на комментарий.

### Подписки (GraphQL Subscriptions)

Реализована возможность подписаться на событие появления нового комментария под постом.

### Аутентификация и управление пользователями

- Регистрация пользователей с возможностью последующей авторизации.

- Генерация JWT-токенов для аутентифицированных пользователей.

- Отслеживание времени последней активности пользователя (изменяется при регистрации, авторизации, создании или редактирования поста, создании комментария комментария).

## Описание реализации
- Для написания сервера GraphQL была использована библиотека [gqlgen](https://github.com/99designs/gqlgen), которая имеет кодогенерацию по схеме.
- Для оптимизации пробелмы n + 1 были сделаны data loaders с помощью библиотеки [dataloaden](https://github.com/vektah/dataloaden), которая также имеет кодогенерацию для конкретного типа данных.
- Для оптимизации запросов в базу данных были созданы индексы для полей, по которым часто извлекаются данные.
- В коде присутствуют кастомные ошибки, пользователь видит только одну из них. В случае, если возникла какая то проблемы, пользователь не будет видеть детали проблемы, а увидит лишь `Internal server error` или другую ошибку, связанную с данными.
- Было реализовано ограничение вложенности запроса. В случае, когда запрос имеет большую вложенность данных, пользователь получает об этом уведомление.

## Запуск приложения
Для приложения написан `Dockerfile`, а также `docker-compose.yml`, в котором дополнительно поднимается контейнер с PostgreSQL и выполняется скрипт инициализации базы со всеми необходимыми таблицами.
1) Создайте конфиг в папке `config` в формате yaml
```yaml
env: "dev"

port: 8080

storage_type: "postgres"

storage:
  db_name: "post_db"
  db_address: "localhost:5432"
  db_user: "postgres"
  db_password: "postgres"
  db_sslmode: "disable"
```
В поле env можно указать `local`, `dev`, `prod`. Эти варианты запуска влияет на уровень и формат логов, а также в `dev` режиме не запускается playground.

В поле `storage_type` необходимо указать тип хранилища: `inmemory` или `postgres`.
2) Добавьте .env файл в корень проекта и укажите там значение `JWT_SECRET` (секретный ключ для генерации JWT токена).
3) Запустите сборку контейнера `docker-compose up -d --build`.
4) Чтобы открыть GraphQL Playground, перейдите на http://localhost:8080/
5) Запросы GraphQL необхоимо отправлять на `/query`

## Примеры запросов

[GraphQL схема](https://github.com/ArtemSarafannikov/OzonTestTask/blob/master/internal/graphql/schema.graphql)

- Запрос списка постов

```graphql
query {
    posts(limit: 10, offset: 0) {
        id
        author {
            username
        }
        title
        content
        allowComments
        createdAt
    }
}
```

- Запрос поста с комментариями

```graphql
query {
    post(postID: "1") {
        title
        comments(limit: 3, offset: 0) {
            id
            text
            author {
                username
            }
            replies(limit: 3, offset: 0) {
                id
                text
            }
        }
    }
}
```

- Создание поста

```graphql
mutation {
    createPost(post: { title: "Новый пост", content: "Пример поста", allowComments: true }) {
        id
        title
    }
}
```

Добавление комментария
```graphql
mutation {
    createComment(comment: { postID: "1", text: "Отличный пост!" }) {
        id
        text
    }
}
```

Подписка на новые комментарии к посту
```graphql
subscription {
    newCommentPost(postID: "1") {
        id
        text
        author {
            username
        }
    }
}
```

## Тестирование
Были написаны unit-тесты для бизнес-логики, [тестовое покрытие](https://github.com/ArtemSarafannikov/OzonTestTask/blob/master/ServiceTestCoverage.html) составляет 100% пакета `service`.
```shell
go test ./internal/service
```

В качестве проверки среднего времени ответа было проведено тестирование на 100 пользователей в течение 5 минут в программе Postman. [Результаты](https://github.com/ArtemSarafannikov/OzonTestTask/blob/master/PostmanPerformanceTest.html)
